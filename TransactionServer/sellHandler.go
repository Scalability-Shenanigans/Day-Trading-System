package main

import (
	"TransactionServer/db"
	"TransactionServer/log"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Sell struct {
	User           string  `json:"user"`
	Stock          string  `json:"stock"`
	Amount         float64 `json:"amount"`
	TransactionNum int     `json:"transactionNum"`
}

type CancelSetSell struct {
	User           string  `json:"user"`
	Stock          string  `json:"stock"`
	Amount         float64 `json:"amount"`
	TransactionNum int     `json:"transactionNum"`
}

type CommitSell struct {
	User           string `json:"user"`
	TransactionNum int    `json:"transactionNum"`
}

type CancelSell struct {
	User           string `json:"user"`
	TransactionNum int    `json:"transactionNum"`
}

func sellHandler(w http.ResponseWriter, r *http.Request) {
	var sell Sell
	err := json.NewDecoder(r.Body).Decode(&sell)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}

	cmd := &log.UserCommand{
		Timestamp:      time.Now().UnixMilli(),
		Server:         "localhost",
		TransactionNum: int64(sell.TransactionNum),
		Command:        "SELL",
		Username:       sell.User,
		Funds:          sell.Amount,
	}
	sysEvent := &log.SystemEvent{
		Timestamp:      time.Now().UnixMilli(),
		Server:         "localhost",
		TransactionNum: int64(sell.TransactionNum),
		Command:        "SELL",
		Username:       sell.User,
		Funds:          sell.Amount,
	}
	log.CreateUserCommandsLog(cmd)
	log.CreateSystemEventLog(sysEvent)

	quote := GetQuote(sell.Stock, sell.User, sell.TransactionNum)

	transaction := db.Transaction{
		User:   sell.User,
		Stock:  sell.Stock,
		Amount: int(sell.Amount / quote), // Amount = no. of shares to be sold
		Price:  quote,
		Is_Buy: false,
	}

	db.CreateTransaction(transaction)
}

func commitSell(w http.ResponseWriter, r *http.Request) {
	var commitSell CommitSell
	err := json.NewDecoder(r.Body).Decode(&commitSell)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}
	user := commitSell.User
	transaction := db.ConsumeLastSellTransaction(user)

	if transaction.Transaction_ID == -1 {
		errorEvent := &log.ErrorEvent{
			Timestamp:    time.Now().UnixMilli(),
			Server:       "localhost",
			Command:      "commitSell",
			Username:     commitSell.User,
			ErrorMessage: "Error: no sell to commit",
		}
		log.CreateErrorEventLog(errorEvent)
		return
	}

	cmd := &log.UserCommand{
		Timestamp:      time.Now().UnixMilli(),
		Server:         "localhost",
		TransactionNum: int64(commitSell.TransactionNum),
		Command:        "COMMIT_SELL",
		Username:       commitSell.User,
	}
	log.CreateUserCommandsLog(cmd)

	// transaction.Amount is no. of shares, transaction.Price is the selling price for one share
	transactionCost := float64(transaction.Amount) * transaction.Price
	// update how much stock they hold after selling
	if db.UpdateStockHolding(user, transaction.Stock, -1*transaction.Amount) {
		if db.UpdateBalance(transactionCost, user, int64(commitSell.TransactionNum)) { // update account balance after selling
			fmt.Println("Transaction Commited")
		}
	}
}

func cancelSell(w http.ResponseWriter, r *http.Request) {
	var cancelSell CancelSell
	err := json.NewDecoder(r.Body).Decode(&cancelSell)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}

	cmd := &log.UserCommand{
		Timestamp:      time.Now().UnixMilli(),
		Server:         "localhost",
		TransactionNum: int64(cancelSell.TransactionNum),
		Command:        "CANCEL_SELL",
		Username:       cancelSell.User,
	}
	log.CreateUserCommandsLog(cmd)

	//consumes the last transaction but does nothing with it so its effectively cancelled
	db.ConsumeLastSellTransaction(cancelSell.User)
}

func setSellAmountHandler(w http.ResponseWriter, r *http.Request) {
	var sellAmountOrder db.SellAmountOrder
	err := json.NewDecoder(r.Body).Decode(&sellAmountOrder)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}

	cmd := &log.UserCommand{
		Timestamp:      time.Now().UnixMilli(),
		Server:         "localhost",
		TransactionNum: int64(sellAmountOrder.TransactionNum),
		Command:        "SET_SELL_AMOUNT",
		Username:       sellAmountOrder.User,
		Funds:          sellAmountOrder.Amount,
	}
	log.CreateUserCommandsLog(cmd)

	fmt.Println(sellAmountOrder)
	// Note: SellAmountOrder.Amount is not the no. of shares, it is the dollar amount user specified in command
	// At this stage we don't know how many shares to sell because we don't have trigger price yet
	db.CreateSellAmountOrder(sellAmountOrder)
}

func setSellTriggerHandler(w http.ResponseWriter, r *http.Request) {
	var triggerOrder TriggerOrder
	err := json.NewDecoder(r.Body).Decode(&triggerOrder)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}
	fmt.Println(triggerOrder)

	cmd := &log.UserCommand{
		Timestamp:      time.Now().UnixMilli(),
		Server:         "localhost",
		TransactionNum: int64(triggerOrder.TransactionNum),
		Command:        "SET_SELL_TRIGGER",
		Username:       triggerOrder.User,
		Funds:          triggerOrder.Price,
	}
	log.CreateUserCommandsLog(cmd)

	// check mongodb for sell Amount object with same user and stock
	found, sellAmountOrder := db.FindSellAmountOrder(triggerOrder.User, triggerOrder.Stock)

	if found {
		fmt.Println("Found SellAmountOrder")
		fmt.Println(sellAmountOrder)

		// calculate max no. of shares which can be sold based on trigger price
		var no_of_shares_to_sell = int(sellAmountOrder.Amount / triggerOrder.Price)

		// check user account to see if they have enough shares of the stock to sell and decrement stock holding
		if db.UpdateStockHolding(sellAmountOrder.User, sellAmountOrder.Stock, -1*no_of_shares_to_sell) {
			fmt.Println("Creating SellAmountOrder")
			// add TriggeredSellAmountOrder to db for PollingService to act on
			var triggeredSellAmountOrder db.TriggeredSellAmountOrder
			triggeredSellAmountOrder.User = sellAmountOrder.User
			triggeredSellAmountOrder.Stock = sellAmountOrder.Stock
			triggeredSellAmountOrder.Num_of_shares = no_of_shares_to_sell
			triggeredSellAmountOrder.Price = triggerOrder.Price
			db.CreateTriggeredSellAmountOrder(triggeredSellAmountOrder)
		}
	}
}

// In this function, I guess we also need to update Stockholding for user
// i.e add back those shares to stockholding which were reserved to be sold
func cancelSetSell(w http.ResponseWriter, r *http.Request) {
	var cancelSetSell CancelSetSell
	err := json.NewDecoder(r.Body).Decode(&cancelSetSell)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}
	fmt.Println(cancelSetSell)

	cmd := &log.UserCommand{
		Timestamp:      time.Now().UnixMilli(),
		Server:         "localhost",
		TransactionNum: int64(cancelSetSell.TransactionNum),
		Command:        "CANCEL_SET_SELL",
		Username:       cancelSetSell.User,
		Funds:          cancelSetSell.Amount,
	}
	log.CreateUserCommandsLog(cmd)

	db.DeleteSellAmountOrder(cancelSetSell.User, cancelSetSell.Stock)
}

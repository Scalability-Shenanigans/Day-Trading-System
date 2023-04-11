package main

import (
	"TransactionServer/db"
	"TransactionServer/log"
	"TransactionServer/middleware"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Sell struct {
	User   string  `json:"user"`
	Stock  string  `json:"stock"`
	Amount float64 `json:"amount"`
}

type SellResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type CancelSetSell struct {
	User   string  `json:"user"`
	Stock  string  `json:"stock"`
	Amount float64 `json:"amount"`
}

type CommitSell struct {
	User string `json:"user"`
}

type CommitSellResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Stock   string `json:"stock,omitempty"`
}

type CancelSell struct {
	User string `json:"user"`
}

func sellHandler(w http.ResponseWriter, r *http.Request) {
	transactionNumber := middleware.GetTransactionNumberFromContext(r)

	var sell Sell
	var response SellResponse

	err := json.NewDecoder(r.Body).Decode(&sell)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}

	cmd := &log.UserCommand{
		Timestamp:      time.Now().UnixMilli(),
		Server:         "localhost",
		TransactionNum: int64(transactionNumber),
		Command:        "SELL",
		Username:       sell.User,
		Funds:          sell.Amount,
	}
	sysEvent := &log.SystemEvent{
		Timestamp:      time.Now().UnixMilli(),
		Server:         "localhost",
		TransactionNum: int64(transactionNumber),
		Command:        "SELL",
		Username:       sell.User,
		Funds:          sell.Amount,
	}
	log.CreateUserCommandsLog(cmd)

	quote := GetQuote(sell.Stock, sell.User, int(transactionNumber))

	if !db.CanSellStock(sell.User, sell.Stock, int(sell.Amount/quote)) {
		// Log an error message and return an appropriate HTTP response
		response.Status = "failure"
		response.Message = "not enough to sell"
		errorEvent := &log.ErrorEvent{
			Timestamp:      time.Now().UnixMilli(),
			Server:         "localhost",
			TransactionNum: int64(transactionNumber),
			Command:        "COMMIT_SELL",
			Username:       sell.User,
			ErrorMessage:   "Error: not enough to sell",
		}
		log.CreateErrorEventLog(errorEvent)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		return
	}
	response.Status = "success"

	log.CreateSystemEventLog(sysEvent)

	transaction := db.Transaction{
		User:      sell.User,
		Stock:     sell.Stock,
		Amount:    int(sell.Amount / quote),
		Price:     quote,
		Is_Buy:    false,
		Timestamp: time.Now().UnixMilli(),
	}

	db.CreatePendingTransaction(transaction)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func commitSell(w http.ResponseWriter, r *http.Request) {
	transactionNumber := middleware.GetTransactionNumberFromContext(r)

	var commitSell CommitSell
	var response CommitSellResponse

	err := json.NewDecoder(r.Body).Decode(&commitSell)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}
	user := commitSell.User
	transaction := db.ConsumeLastSellTransaction(user)
	db.CreateFinishedTransaction(transaction)

	if transaction.Transaction_ID == -1 {
		errorEvent := &log.ErrorEvent{
			Timestamp:      time.Now().UnixMilli(),
			Server:         "localhost",
			TransactionNum: int64(transactionNumber),
			Command:        "COMMIT_SELL",
			Username:       commitSell.User,
			ErrorMessage:   "Error: no sell to commit",
		}
		log.CreateErrorEventLog(errorEvent)
		return
	}

	cmd := &log.UserCommand{
		Timestamp:      time.Now().UnixMilli(),
		Server:         "localhost",
		TransactionNum: int64(transactionNumber),
		Command:        "COMMIT_SELL",
		Username:       commitSell.User,
	}
	log.CreateUserCommandsLog(cmd)

	// transaction.Amount is no. of shares, transaction.Price is the selling price for one share
	transactionCost := float64(transaction.Amount) * transaction.Price
	if transactionCost != 0 {
		// update how much stock they hold after selling
		if db.UpdateStockHolding(user, transaction.Stock, -1*transaction.Amount, int64(transactionNumber)) {
			if db.UpdateBalance(transactionCost, user, int64(transactionNumber)) { // update account balance after selling
				fmt.Println("Transaction Commited")
				response.Status = "success"
				response.Message = "Sell committed successfully"
				response.Stock = transaction.Stock
			}
		}
	} else {
		errorEvent := &log.ErrorEvent{
			Timestamp:      time.Now().UnixMilli(),
			Server:         "localhost",
			TransactionNum: int64(transactionNumber),
			Command:        "COMMIT_SELL",
			Username:       commitSell.User,
		}
		log.CreateErrorEventLog(errorEvent)

		response.Status = "failure"
		response.Message = "Sell commit failure"
		response.Stock = transaction.Stock
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func cancelSell(w http.ResponseWriter, r *http.Request) {
	transactionNumber := middleware.GetTransactionNumberFromContext(r)

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
		TransactionNum: int64(transactionNumber),
		Command:        "CANCEL_SELL",
		Username:       cancelSell.User,
	}
	log.CreateUserCommandsLog(cmd)

	//consumes the last transaction but does nothing with it so its effectively cancelled
	db.ConsumeLastSellTransaction(cancelSell.User)
}

func setSellAmountHandler(w http.ResponseWriter, r *http.Request) {
	transactionNumber := middleware.GetTransactionNumberFromContext(r)

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
		TransactionNum: int64(transactionNumber),
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
	transactionNumber := middleware.GetTransactionNumberFromContext(r)

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
		TransactionNum: int64(transactionNumber),
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
		if db.UpdateStockHolding(sellAmountOrder.User, sellAmountOrder.Stock, -1*no_of_shares_to_sell, int64(transactionNumber)) {
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
	transactionNumber := middleware.GetTransactionNumberFromContext(r)

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
		TransactionNum: int64(transactionNumber),
		Command:        "CANCEL_SET_SELL",
		Username:       cancelSetSell.User,
		Funds:          cancelSetSell.Amount,
	}
	log.CreateUserCommandsLog(cmd)

	db.DeleteSellAmountOrder(cancelSetSell.User, cancelSetSell.Stock)
}

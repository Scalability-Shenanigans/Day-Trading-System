package main

import (
	"TransactionServer/db"
	"TransactionServer/log"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TriggerOrder struct {
	User  string  `json:"user"`
	Stock string  `json:"stock"`
	Price float64 `json:"price"`
}

type Buy struct {
	User           string  `json:"user"`
	Stock          string  `json:"stock"`
	Amount         float64 `json:"amount"`
	TransactionNum int     `json:"transactionNum"`
}

func buyHandler(w http.ResponseWriter, r *http.Request) {
	//get stock price
	//add transaction to pending transactions collection
	//	store user, transactionId, stock name, price, amount, anything else we need
	//we can either check if a user has enough funds here or during the commit
	//	(both probably work commit just means we do more processing before failing)
	var buy Buy
	err := json.NewDecoder(r.Body).Decode(&buy)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}

	cmd := &log.UserCommand{
		Timestamp:      time.Now().UnixNano(),
		Server:         "localhost",
		TransactionNum: int64(buy.TransactionNum),
		Command:        "BUY",
		Username:       buy.User,
		Funds:          buy.Amount,
	}

	log.CreateUserCommandsLog(cmd)

	quote := GetQuote(buy.Stock, buy.User)

	transaction := db.Transaction{
		User:   buy.User,
		Stock:  buy.Stock,
		Amount: int(buy.Amount / quote),
		Price:  quote,
		Is_Buy: true,
	}

	db.CreateTransaction(transaction)
}

func commitBuy(w http.ResponseWriter, r *http.Request) {
	//read the last transaction from the db
	//consume it (delete it after its done)
	//update users account in db
	// balance, holdings, ???
	user := r.URL.Query().Get("user")
	transaction := db.ConsumeLastTransaction(user)

	transactionCost := float64(transaction.Amount) * transaction.Price

	if db.UpdateBalance(transactionCost*-1.0, user, 0) {
		if db.UpdateStockHolding(user, transaction.Stock, transaction.Amount) {
			fmt.Println("Transaction Commited")
		}
	}
}

func cancelBuy(w http.ResponseWriter, r *http.Request) {
	var buy Buy
	err := json.NewDecoder(r.Body).Decode(&buy)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}
	//consumes the last transaction but does nothing with it so its effectively cancelled
	db.ConsumeLastTransaction(buy.User)

}

func setBuyAmountHandler(w http.ResponseWriter, r *http.Request) {
	var buyAmountOrder db.BuyAmountOrder
	err := json.NewDecoder(r.Body).Decode(&buyAmountOrder)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}
	fmt.Println(buyAmountOrder)

	cmd := &log.UserCommand{
		Timestamp:      time.Now().UnixNano(),
		Server:         "localhost",
		TransactionNum: int64(buyAmountOrder.TransactionNum),
		Command:        "SET_BUY_AMOUNT",
		Username:       buyAmountOrder.User,
		Funds:          buyAmountOrder.Amount,
	}
	log.CreateUserCommandsLog(cmd)

	db.CreateBuyAmountOrder(buyAmountOrder) // Add buyAmountOrder to db
}

func setBuyTriggerHandler(w http.ResponseWriter, r *http.Request) {
	var triggerOrder TriggerOrder
	err := json.NewDecoder(r.Body).Decode(&triggerOrder)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}
	fmt.Println(triggerOrder)

	// check mongodb for buyAmount object with same user and stock
	found, buyAmountOrder := db.FindBuyAmountOrder(triggerOrder.User, triggerOrder.Stock)

	if found {
		fmt.Println("Found BuyAmountOrder")
		fmt.Println(buyAmountOrder)

		// check user account to see if they have enough funds and decrement Account balance if they do
		if db.UpdateBalance((buyAmountOrder.Amount * triggerOrder.Price * -1), buyAmountOrder.User, 0) {
			fmt.Println("Creating BuyAmountOrder")
			// add TriggeredBuyAmountOrder to db for PollingService to act on
			var triggeredBuyAmountOrder db.TriggeredBuyAmountOrder
			triggeredBuyAmountOrder.User = buyAmountOrder.User
			triggeredBuyAmountOrder.Stock = buyAmountOrder.Stock
			triggeredBuyAmountOrder.Amount = buyAmountOrder.Amount
			triggeredBuyAmountOrder.Price = triggerOrder.Price
			db.CreateTriggeredBuyAmountOrder(triggeredBuyAmountOrder)
		}
	}
}

func cancelSetBuy(http.ResponseWriter, *http.Request) {
	//undo everything you did in the setBuyAmount one
}

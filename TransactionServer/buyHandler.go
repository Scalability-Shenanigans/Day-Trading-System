package main

import (
	"TransactionServer/db"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type TriggerOrder struct {
	User  string  `json:"user"`
	Stock string  `json:"stock"`
	Price float64 `json:"price"`
}

func buyHandler(w http.ResponseWriter, r *http.Request) {
	//get stock price
	//add transaction to pending transactions collection
	//	store user, transactionId, stock name, price, amount, anything else we need
	//we can either check if a user has enough funds here or during the commit
	//	(both probably work commit just means we do more processing before failing)
	user := r.URL.Query().Get("user")
	stock := r.URL.Query().Get("stock")
	amount, err := strconv.Atoi(r.URL.Query().Get("amount"))
	if err != nil {
		fmt.Println(err)
		return
	}

	quote, err := strconv.Atoi(GetQuote(stock, user))
	if err != nil {
		fmt.Println(err)
		return
	}

	transaction := db.Transaction{
		User:   user,
		Stock:  stock,
		Amount: amount,
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

	transactionCost := transaction.Amount * transaction.Price

	if db.UpdateBalance(transactionCost*-1, user) {
		if db.UpdateStockHolding(user, transaction.Stock, transaction.Amount) {
			fmt.Println("Transaction Commited")
		}
	}
}

func cancelBuy(w http.ResponseWriter, r *http.Request) {
	//cancel the buy
	//delete it from db and everywhere else
	//do we reserve funds when a buy is added but not commited?
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

	// check user account to see if they have enough funds and decrement Account balance if they do
	if db.UpdateBalance(int(buyAmountOrder.Amount*-1), buyAmountOrder.User) {
		fmt.Println("Creating BuyAmountOrder")
		db.CreateBuyAmountOrder(buyAmountOrder) // Add buyAmountOrder to db
	}
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
		fmt.Println(buyAmountOrder)
		fmt.Println("Found BuyAmountOrder")
		//check stock price
		//if price is <= trigger price
		//buy x amount of the stock -> update user stock holdings
		// else send to polling service -> user, stock, amount, trigger price
	}
}

func cancelSetBuy(http.ResponseWriter, *http.Request) {
	//undo everything you did in the setBuyAmount one
}

package main

import (
	"TransactionServer/db"
	"encoding/json"
	"fmt"
	"net/http"
)

// type TriggerOrder struct {
// 	User  string  `json:"user"`
// 	Stock string  `json:"stock"`
// 	Price float64 `json:"price"`
// }

type Sell struct {
	User   string  `json:"user"`
	Stock  string  `json:"stock"`
	Amount float64 `json:"amount"`
}

func sellHandler(w http.ResponseWriter, r *http.Request) {

	var sell Sell

	err := json.NewDecoder(r.Body).Decode(&sell)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}

	quote := GetQuote(sell.Stock, sell.User)

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

	user := r.URL.Query().Get("user")
	transaction := db.ConsumeLastTransaction(user)
	// transaction.Amount is no. of shares, transaction.Price is the selling price for one share
	transactionCost := float64(transaction.Amount) * transaction.Price
	// Update user's account balance after they sold stock
	if db.UpdateBalance(transactionCost, user, 0) {
		if db.UpdateStockHolding(user, transaction.Stock, -1*transaction.Amount) { // update how much stock they hold after selling
			fmt.Println("Transaction Commited")
		}
	}
}

func cancelSell(w http.ResponseWriter, r *http.Request) {
	//cancel the buy
	//delete it from db and everywhere else
	//do we reserve funds when a buy is added but not commited?
}

func setSellAmountHandler(w http.ResponseWriter, r *http.Request) {
	var sellAmountOrder db.SellAmountOrder
	err := json.NewDecoder(r.Body).Decode(&sellAmountOrder)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}
	fmt.Println(sellAmountOrder)

	db.CreateSellAmountOrder(sellAmountOrder) // Add buyAmountOrder to db

}

func setSellTriggerHandler(w http.ResponseWriter, r *http.Request) {

}

func cancelSetSell(http.ResponseWriter, *http.Request) {
	//undo everything you did in the setBuyAmount one
}

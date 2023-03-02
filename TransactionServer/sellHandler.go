package main

import (
	"TransactionServer/db"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

	quote:= GetQuote(sell.Stock, sell.User)

	transaction := db.Transaction{
		User:   sell.User,
		Stock:  sell.Stock,
		Amount: int(sell.Amount),
		Price:  quote,
		Is_Buy: false,
	}

	db.CreateTransaction(transaction)

}

func commitSell(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	transaction := db.ConsumeLastTransaction(user)

	// since Amount is not no. of shares, it is baically the sell amount
	// After selling, add the sell amount to user's acct balance
	if db.UpdateBalance(transaction.Amount, user) {
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

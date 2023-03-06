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

type CancelSetSell struct {
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
	var triggerOrder TriggerOrder
	err := json.NewDecoder(r.Body).Decode(&triggerOrder)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}
	fmt.Println(triggerOrder)

	// check mongodb for sell Amount object with same user and stock
	found, sellAmountOrder := db.FindSellAmountOrder(triggerOrder.User, triggerOrder.Stock)

	if found {
		fmt.Println("Found SellAmountOrder")
		fmt.Println(sellAmountOrder)

		// check user account to see if they have enough funds and decrement Account balance if they do
		if db.UpdateStockHolding(sellAmountOrder.User, sellAmountOrder.Stock, -1*int(sellAmountOrder.Amount)) {
			fmt.Println("Triggering SellAmountOrder")
			// add TriggeredBuyAmountOrder to db for PollingService to act on
			var triggeredSellAmountOrder db.TriggeredSellAmountOrder
			triggeredSellAmountOrder.User = sellAmountOrder.User
			triggeredSellAmountOrder.Stock = sellAmountOrder.Stock
			triggeredSellAmountOrder.Amount = sellAmountOrder.Amount
			triggeredSellAmountOrder.Price = triggerOrder.Price
			db.CreateTriggeredSellAmountOrder(triggeredSellAmountOrder)
		}
	}
}

func cancelSetSell(w http.ResponseWriter, r *http.Request) {
	var cancelSetSell CancelSetSell
	err := json.NewDecoder(r.Body).Decode(&cancelSetSell)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}
	fmt.Println(cancelSetSell)

	db.DeleteSellAmountOrder(cancelSetSell.User, cancelSetSell.Stock)
}

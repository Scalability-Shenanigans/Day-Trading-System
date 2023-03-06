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
	// update how much stock they hold after selling
	if db.UpdateStockHolding(user, transaction.Stock, -1*transaction.Amount) {
		if db.UpdateBalance(transactionCost, user, 0) { // update account balance after selling
			fmt.Println("Transaction Commited")
		}
	}
}

func cancelSell(w http.ResponseWriter, r *http.Request) {
	var sell Sell
	err := json.NewDecoder(r.Body).Decode(&sell)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}
	//consumes the last transaction but does nothing with it so its effectively cancelled
	db.ConsumeLastTransaction(sell.User)
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
	// Note: SellAmountOrder.Amount is not the no. of shares, it is the dollar amount user specified in command
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

	// check mongodb for sellAmount object with same user and stock
	found, sellAmountOrder := db.FindSellAmountOrder(triggerOrder.User, triggerOrder.Stock)

	if found {
		fmt.Println("Found SellAmountOrder")
		fmt.Println(sellAmountOrder)

		var no_of_shares_to_sell = int(sellAmountOrder.Amount / triggerOrder.Price)

		// check user account to see if they have enough shares of the stock to sell and decrement stock holding
		if db.UpdateStockHolding(sellAmountOrder.User, sellAmountOrder.Stock, -1*no_of_shares_to_sell) {
			fmt.Println("Creating SellAmountOrder")
			// add TriggeredSellAmountOrder to db for PollingService to act on
			var triggeredSellAmountOrder db.TriggeredSellAmountOrder
			triggeredSellAmountOrder.User = sellAmountOrder.User
			triggeredSellAmountOrder.Stock = sellAmountOrder.Stock
			triggeredSellAmountOrder.Amount = no_of_shares_to_sell
			triggeredSellAmountOrder.Price = triggerOrder.Price
			db.CreateTriggeredSellAmountOrder(triggeredSellAmountOrder)
		}
	}
}

func cancelSetSell(http.ResponseWriter, *http.Request) {
	//undo everything you did in the setBuyAmount one
}

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

type Sell struct {
	User   string  `json:"user"`
	Stock  string  `json:"stock"`
	Amount float64 `json:"amount"`
}


func sellHandler(w http.ResponseWriter, r *http.Request) {

	var sell Buy

	err := json.NewDecoder(r.Body).Decode(&sell)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}

	quote, err := strconv.Atoi(GetQuote(sell.Stock, sell.User))
	if err != nil {
		fmt.Println(err)
		return
	}

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

	
	
}

func cancelSell(w http.ResponseWriter, r *http.Request) {
	//cancel the buy
	//delete it from db and everywhere else
	//do we reserve funds when a buy is added but not commited?
}

func setSellAmountHandler(w http.ResponseWriter, r *http.Request) {
	
}

func setSellTriggerHandler(w http.ResponseWriter, r *http.Request) {
	
}

func cancelSetSell(http.ResponseWriter, *http.Request) {
	//undo everything you did in the setBuyAmount one
}

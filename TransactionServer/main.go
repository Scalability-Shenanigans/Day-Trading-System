package main

import (
	"TransactionServer/db"
	"encoding/json"
	"fmt"
	"net/http"
)

type AddFunds struct {
	User   string  `json:"user"`
	Amount float64 `json:"amount"`
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	//check if user exists in db
	//if not create user
	//add whatever the funds amount is
	var addFunds AddFunds
	err := json.NewDecoder(r.Body).Decode(&addFunds)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}
	if db.UpdateBalance(int(addFunds.Amount), addFunds.User) {
		return
	}
	db.CreateAccount(addFunds.User, int(addFunds.Amount))
}

func createUser() {
	//create user in db
}

func main() {
	port := ":8080"
	db.InitConnection()
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/buy", buyHandler)
	http.HandleFunc("/commitBuy", commitBuy)
	http.HandleFunc("/cancelBuy", cancelBuy)
	http.HandleFunc("/setBuyAmount", setBuyAmountHandler)
	http.HandleFunc("/setBuyTrigger", setBuyTriggerHandler)
	http.ListenAndServe(port, nil)
}

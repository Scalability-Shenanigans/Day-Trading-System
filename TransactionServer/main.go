package main

import (
	"TransactionServer/db"
	"TransactionServer/log"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AddFunds struct {
	User           string  `json:"user"`
	Amount         float64 `json:"amount"`
	TransactionNum int     `json:"transactionNum"`
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	//check if user exists in db
	//if not create user
	//add whatever the funds amount is
	var addFunds AddFunds

	// var addUserCommand UserCommand;

	err := json.NewDecoder(r.Body).Decode(&addFunds)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}

	cmd := &log.UserCommand{
		Timestamp:      time.Now().UnixNano(),
		Server:         "localhost",
		TransactionNum: int64(addFunds.TransactionNum),
		Command:        "ADD",
		Username:       addFunds.User,
		Funds:          addFunds.Amount,
	}

	log.CreateUserCommandsLog(cmd)

	if db.UpdateBalance(addFunds.Amount, addFunds.User, int64(addFunds.TransactionNum)) {
		return
	}
	fmt.Println("Creating an account for user")
	db.CreateAccount(addFunds.User, addFunds.Amount, int64(addFunds.TransactionNum))
}

func createUser() {
	//create user in db
}

func main() {
	port := ":8080"
	db.InitConnection()
	log.InitLogDBConnection()
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/buy", buyHandler)
	http.HandleFunc("/commitBuy", commitBuy)
	http.HandleFunc("/cancelBuy", cancelBuy)
	http.HandleFunc("/setBuyAmount", setBuyAmountHandler)
	http.HandleFunc("/setBuyTrigger", setBuyTriggerHandler)
	http.HandleFunc("/dumplog", log.DumplogHandler)
	http.ListenAndServe(port, nil)
}

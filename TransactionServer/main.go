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
		Timestamp:      time.Now().UnixMilli(),
		Server:         "localhost",
		TransactionNum: int64(addFunds.TransactionNum),
		Command:        "ADD",
		Username:       addFunds.User,
		Funds:          addFunds.Amount,
	}

	transaction := &log.AccountTransaction{
		Timestamp:      time.Now().UnixMilli(),
		Server:         "localhost",
		TransactionNum: int64(addFunds.TransactionNum),
		Action:         "add",
		Username:       addFunds.User,
		Funds:          addFunds.Amount,
	}

	log.CreateUserCommandsLog(cmd)
	log.CreateAccountTransactionLog(transaction)

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
	http.HandleFunc("/cancelSetBuy", cancelSetBuy)
	http.HandleFunc("/setBuyAmount", setBuyAmountHandler)
	http.HandleFunc("/setBuyTrigger", setBuyTriggerHandler)
	http.HandleFunc("/sell", sellHandler)
	http.HandleFunc("/commitSell", commitSell)
	http.HandleFunc("/cancelSell", cancelSell)
	http.HandleFunc("/cancelSetSell", cancelSetSell)
	http.HandleFunc("/setSellAmount", setSellAmountHandler)
	http.HandleFunc("/setSellTrigger", setSellTriggerHandler)
	http.HandleFunc("/dumplog", log.DumplogHandler)
	http.HandleFunc("/quote", quoteHandler)
	http.HandleFunc("/display", quoteHandler)
	http.ListenAndServe(port, nil)
}

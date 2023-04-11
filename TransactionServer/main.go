package main

import (
	"TransactionServer/db"
	"TransactionServer/log"
	"TransactionServer/middleware"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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

	transactionNumber := middleware.GetTransactionNumberFromContext(r)

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
		TransactionNum: int64(transactionNumber),
		Command:        "ADD",
		Username:       addFunds.User,
		Funds:          addFunds.Amount,
	}

	log.CreateUserCommandsLog(cmd)

	if db.UpdateBalance(addFunds.Amount, addFunds.User, int64(transactionNumber)) {
		return
	}
	fmt.Println("Creating an account for user")
	db.CreateAccount(addFunds.User, addFunds.Amount, int64(transactionNumber))
}

func createUser() {
	//create user in db
}

func main() {
	port := ":8080"
	mux := http.NewServeMux()
	db.InitConnection()
	log.InitLogDBConnection()

	// Register handlers with mux
	mux.HandleFunc("/add", addHandler)
	mux.HandleFunc("/buy", buyHandler)
	mux.HandleFunc("/commitBuy", commitBuy)
	mux.HandleFunc("/cancelBuy", cancelBuy)
	mux.HandleFunc("/cancelSetBuy", cancelSetBuy)
	mux.HandleFunc("/setBuyAmount", setBuyAmountHandler)
	mux.HandleFunc("/setBuyTrigger", setBuyTriggerHandler)
	mux.HandleFunc("/sell", sellHandler)
	mux.HandleFunc("/commitSell", commitSell)
	mux.HandleFunc("/cancelSell", cancelSell)
	mux.HandleFunc("/cancelSetSell", cancelSetSell)
	mux.HandleFunc("/setSellAmount", setSellAmountHandler)
	mux.HandleFunc("/setSellTrigger", setSellTriggerHandler)
	mux.HandleFunc("/dumplog", log.DumplogHandler)
	mux.HandleFunc("/quote", quoteHandler)
	mux.HandleFunc("/display", displayHandler)

	http.HandleFunc("/dbwipe", db.DBWiper)

	http.Handle("/", middleware.TransactionNumberMiddleware(mux))
	http.ListenAndServe(port, nil)
}

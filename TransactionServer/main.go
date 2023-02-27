package main

import (
	"TransactionServer/db"
	"net/http"
	"strconv"
)

func addHandler(w http.ResponseWriter, r *http.Request) {
	//check if user exists in db
	//if not create user
	//add whatever the funds amount is
	user := r.URL.Query().Get("user")
	funds, err := strconv.Atoi(r.URL.Query().Get("funds"))
	if err != nil {
		return
	}
	//update balance makes a new account if it doesnt exist so create user is redundant atm
	db.UpdateBalance(funds, user)
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
	http.HandleFunc("/setBuyAmount", setBuyAmount)
	http.ListenAndServe(port, nil)
}

package main

import (
	"TransactionServer/db"
	"net/http"
)

func addHandler(http.ResponseWriter, *http.Request) {
	//check if user exists in db
	//if not create user
	//add whatever the funds amount is
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

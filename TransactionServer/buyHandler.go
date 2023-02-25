package main

import (
	"net/http"
)

func buyHandler(w http.ResponseWriter, r *http.Request) {
	//get stock price
	//add transaction to pending transactions collection
	//	store user, transactionId, stock name, price, amount, anything else we need
	//we can either check if a user has enough funds here or during the commit
	//	(both probably work commit just means we do more processing before failing)
}

func commitBuy(w http.ResponseWriter, r *http.Request) {
	//read the transaction from the db
	//consume it (delete it after its done)
	//update users account in db
	// balance, holdings, ???
}

func cancelBuy(w http.ResponseWriter, r *http.Request) {
	//cancel the buy
	//delete it from db and everywhere else
	//do we reserve funds when a buy is added but not commited?
}

func setBuyAmount(w http.ResponseWriter, r *http.Request) {
	//check price
	//if price is <= buy order
	//buy x amount of the stock
	//if it isnt't add the buy order to db?
	// how are we dealing with buy orders?
	// periodically polling?
	// if we do this i think we make a seperate service that just handles buy orders on its own so the main tx service doesn't waste time
	// spawn a thread each time a buy order is placed?
	// sockets will be a terrible time if we do this i think
}

func cancelSetBuy(http.ResponseWriter, *http.Request) {
	//undo everything you did in the setBuyAmount one
}

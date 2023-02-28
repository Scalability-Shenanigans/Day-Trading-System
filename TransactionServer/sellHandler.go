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

func sellHandler(w http.ResponseWriter, r *http.Request) {

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

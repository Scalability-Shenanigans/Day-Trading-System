package main

import (
	"TransactionServer/db"
	"TransactionServer/log"
	"TransactionServer/middleware"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type TriggerOrder struct {
	User  string  `json:"user"`
	Stock string  `json:"stock"`
	Price float64 `json:"price"`
}

type Buy struct {
	User   string  `json:"user"`
	Stock  string  `json:"stock"`
	Amount float64 `json:"amount"`
}

type CancelBuy struct {
	User string `json:"user"`
}

type CancelSetBuy struct {
	User  string `json:"user"`
	Stock string `json:"stock"`
}

type CommitBuy struct {
	User string `json:"user"`
}

type CommitBuyResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Stock   string `json:"stock,omitempty"`
}

func buyHandler(w http.ResponseWriter, r *http.Request) {
	//get stock price
	//add transaction to pending transactions collection
	//	store user, transactionId, stock name, price, amount, anything else we need
	//we can either check if a user has enough funds here or during the commit
	//	(both probably work commit just means we do more processing before failing)
	var buy Buy

	transactionNumber := middleware.GetTransactionNumberFromContext(r)

	err := json.NewDecoder(r.Body).Decode(&buy)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}

	cmd := &log.UserCommand{
		Timestamp:      time.Now().UnixMilli(),
		Server:         "localhost",
		TransactionNum: int64(transactionNumber),
		Command:        "BUY",
		Username:       buy.User,
		Funds:          buy.Amount,
	}
	sysEvent := &log.SystemEvent{
		Timestamp:      time.Now().UnixMilli(),
		Server:         "localhost",
		TransactionNum: int64(transactionNumber),
		Command:        "BUY",
		Username:       buy.User,
		Funds:          buy.Amount,
	}

	log.CreateUserCommandsLog(cmd)
	log.CreateSystemEventLog(sysEvent)

	quote := GetQuote(buy.Stock, buy.User, int(transactionNumber))

	transaction := db.Transaction{
		User:   buy.User,
		Stock:  buy.Stock,
		Amount: int(buy.Amount / quote),
		Price:  quote,
		Is_Buy: true,
	}

	db.CreatePendingTransaction(transaction)
}

func commitBuy(w http.ResponseWriter, r *http.Request) {
	//read the last transaction from the db
	//consume it (delete it after its done)
	//update users account in db
	// balance, holdings, ???
	transactionNumber := middleware.GetTransactionNumberFromContext(r)

	var commitBuy CommitBuy
	var response CommitBuyResponse

	err := json.NewDecoder(r.Body).Decode(&commitBuy)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}

	user := commitBuy.User
	transaction := db.ConsumeLastBuyTransaction(user)
	db.CreateFinishedTransaction(transaction)

	if transaction.Transaction_ID == -1 {
		errorEvent := &log.ErrorEvent{
			Timestamp:      time.Now().UnixMilli(),
			Server:         "localhost",
			TransactionNum: int64(transactionNumber),
			Command:        "COMMIT_BUY",
			Username:       commitBuy.User,
			ErrorMessage:   "Error: no buy to commit",
		}
		log.CreateErrorEventLog(errorEvent)
		return
	}

	cmd := &log.UserCommand{
		Timestamp:      time.Now().UnixMilli(),
		Server:         "localhost",
		TransactionNum: int64(transactionNumber),
		Command:        "COMMIT_BUY",
		Username:       commitBuy.User,
		Funds:          float64(transaction.Amount),
	}
	log.CreateUserCommandsLog(cmd)

	transactionCost := float64(transaction.Amount) * transaction.Price

	fmt.Println("transactionCost is " + strconv.FormatFloat(transactionCost, 'E', -1, 64))

	if transactionCost != 0 {
		if db.UpdateBalance(transactionCost*-1.0, user, int64(transactionNumber)) {
			if db.UpdateStockHolding(user, transaction.Stock, transaction.Amount, int64(transactionNumber)) {
				fmt.Println("Transaction Commited")
				response.Status = "success"
				response.Message = "Buy committed successfully"
				response.Stock = transaction.Stock
			}
		}
	} else {
		errorEvent := &log.ErrorEvent{
			Timestamp:      time.Now().UnixMilli(),
			Server:         "localhost",
			TransactionNum: int64(transactionNumber),
			Command:        "COMMIT_BUY",
			Username:       commitBuy.User,
		}
		log.CreateErrorEventLog(errorEvent)

		response.Status = "failure"
		response.Message = "Buy commit failure"
		response.Stock = transaction.Stock
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func cancelBuy(w http.ResponseWriter, r *http.Request) {
	transactionNumber := middleware.GetTransactionNumberFromContext(r)

	var cancelBuy CancelBuy
	err := json.NewDecoder(r.Body).Decode(&cancelBuy)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}

	cmd := &log.UserCommand{
		Timestamp:      time.Now().UnixMilli(),
		Server:         "localhost",
		TransactionNum: int64(transactionNumber),
		Command:        "CANCEL_BUY",
		Username:       cancelBuy.User,
	}
	log.CreateUserCommandsLog(cmd)

	//consumes the last transaction but does nothing with it so its effectively cancelled
	db.ConsumeLastBuyTransaction(cancelBuy.User)

}

func setBuyAmountHandler(w http.ResponseWriter, r *http.Request) {
	transactionNumber := middleware.GetTransactionNumberFromContext(r)

	var buyAmountOrder db.BuyAmountOrder
	err := json.NewDecoder(r.Body).Decode(&buyAmountOrder)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}
	fmt.Println(buyAmountOrder)

	cmd := &log.UserCommand{
		Timestamp:      time.Now().UnixMilli(),
		Server:         "localhost",
		TransactionNum: int64(transactionNumber),
		Command:        "SET_BUY_AMOUNT",
		Username:       buyAmountOrder.User,
		Funds:          buyAmountOrder.Amount,
	}
	log.CreateUserCommandsLog(cmd)

	db.CreateBuyAmountOrder(buyAmountOrder) // Add buyAmountOrder to db
}

func setBuyTriggerHandler(w http.ResponseWriter, r *http.Request) {
	transactionNumber := middleware.GetTransactionNumberFromContext(r)

	var triggerOrder TriggerOrder
	err := json.NewDecoder(r.Body).Decode(&triggerOrder)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}
	fmt.Println(triggerOrder)

	cmd := &log.UserCommand{
		Timestamp:      time.Now().UnixMilli(),
		Server:         "localhost",
		TransactionNum: int64(transactionNumber),
		Command:        "SET_BUY_TRIGGER",
		Username:       triggerOrder.User,
		Funds:          triggerOrder.Price,
	}
	log.CreateUserCommandsLog(cmd)

	// check mongodb for buyAmount object with same user and stock
	found, buyAmountOrder := db.FindBuyAmountOrder(triggerOrder.User, triggerOrder.Stock)

	if found {
		fmt.Println("Found BuyAmountOrder")
		fmt.Println(buyAmountOrder)

		// check user account to see if they have enough funds and decrement Account balance if they do
		if db.UpdateBalance((buyAmountOrder.Amount * triggerOrder.Price * -1), buyAmountOrder.User, int64(transactionNumber)) {
			fmt.Println("Creating BuyAmountOrder")
			// add TriggeredBuyAmountOrder to db for PollingService to act on
			var triggeredBuyAmountOrder db.TriggeredBuyAmountOrder
			triggeredBuyAmountOrder.User = buyAmountOrder.User
			triggeredBuyAmountOrder.Stock = buyAmountOrder.Stock
			triggeredBuyAmountOrder.Amount = buyAmountOrder.Amount
			triggeredBuyAmountOrder.Price = triggerOrder.Price
			db.CreateTriggeredBuyAmountOrder(triggeredBuyAmountOrder)
		}
	}
}

func cancelSetBuy(w http.ResponseWriter, r *http.Request) {
	transactionNumber := middleware.GetTransactionNumberFromContext(r)

	var cancelSetBuy CancelSetBuy
	err := json.NewDecoder(r.Body).Decode(&cancelSetBuy)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}
	fmt.Println(cancelSetBuy)

	cmd := &log.UserCommand{
		Timestamp:      time.Now().UnixMilli(),
		Server:         "localhost",
		TransactionNum: int64(transactionNumber),
		Command:        "CANCEL_SET_BUY",
		Username:       cancelSetBuy.User,
	}
	log.CreateUserCommandsLog(cmd)

	db.DeleteBuyAmountOrder(cancelSetBuy.User, cancelSetBuy.Stock)
}

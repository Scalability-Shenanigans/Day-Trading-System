package main

import (
	"TransactionServer/db"
	"TransactionServer/log"
	"TransactionServer/middleware"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/cors"
)

type AddFunds struct {
	User   string  `json:"user"`
	Amount float64 `json:"amount"`
}

type AddResponse struct {
	Balance float64 `json:"balance,omitempty"`
}

type GetBalance struct {
	User string `json:"user"`
}

type GetStocks struct {
	User string `json:"user"`
}

type GetBalanceResponse struct {
	Balance float64 `json:"balance"`
}

type GetStocksResponse struct {
	Stocks []db.StockHolding `json:"stock_holding"`
}

type TransactionsByUserRequest struct {
	User string `json:"user"`
}

type GetAllUserTransactions struct {
	User string `json:"user"`
}

type TransactionsByUserResponse struct {
	PendingTransactions  []db.Transaction `json:"pending_transactions"`
	FinishedTransactions []db.Transaction `json:"finished_transactions"`
}

func allTransactionsByUserHandler(w http.ResponseWriter, r *http.Request) {
	var req TransactionsByUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	pendingTransactions, err := db.GetPendingTransactionsByUser(req.User)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	finishedTransactions, err := db.GetFinishedTransactionsByUser(req.User)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp := TransactionsByUserResponse{
		PendingTransactions:  pendingTransactions,
		FinishedTransactions: finishedTransactions,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	//check if user exists in db
	//if not create user
	//add whatever the funds amount is
	var addFunds AddFunds
	var response AddResponse

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
		response.Balance = db.GetBalance(addFunds.User)
	} else {
		fmt.Println("Creating an account for user")
		db.CreateAccount(addFunds.User, addFunds.Amount, int64(transactionNumber))

		response.Balance = addFunds.Amount
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getBalanceHandler(w http.ResponseWriter, r *http.Request) {
	var getBalance GetBalance
	var response GetBalanceResponse

	err := json.NewDecoder(r.Body).Decode(&getBalance)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}

	response.Balance = db.GetBalance(getBalance.User)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getStocksHandler(w http.ResponseWriter, r *http.Request) {
	var getStocks GetStocks
	var response GetStocksResponse

	err := json.NewDecoder(r.Body).Decode(&getStocks)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}

	userStockHolding, err := db.GetStockHoldings(getStocks.User)

	if err != nil {
		fmt.Println("cannot find stocks of user")
		return
	}

	response.Stocks = userStockHolding

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func createUser() {
	//create user in db
}

func main() {
	port := ":8080"
	mux := http.NewServeMux()
	db.InitConnection()
	log.InitLogDBConnection()

	// Wrap your mux with the CORS middleware
	corsWrapper := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

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

	http.Handle("/", corsWrapper.Handler(middleware.TransactionNumberMiddleware(mux)))

	// frontend specific endpoints
	http.Handle("/getBalance", corsWrapper.Handler(http.HandlerFunc(getBalanceHandler)))
	http.Handle("/allTransactionsByUser", corsWrapper.Handler(http.HandlerFunc(allTransactionsByUserHandler)))
	http.Handle("/stocks", corsWrapper.Handler(http.HandlerFunc(getStocksHandler)))

	http.ListenAndServe(port, nil)
}

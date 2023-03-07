package main

import (
	"TransactionServer/log"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	SERVER_HOST = "quoteserve.seng.uvic.ca"
	SERVER_PORT = "4444"
	SERVER_TYPE = "tcp"
)

type TransactionResult struct {
	Price     float64
	Symbol    string
	Username  string
	TimeStamp int
	Key       string
}
type Quote struct {
	User  string `json:"user"`
	Stock string `json:"stock"`
}

func quoteHandler(w http.ResponseWriter, r *http.Request) {
	//get stock price
	//add user command to log
	var quote Quote
	err := json.NewDecoder(r.Body).Decode(&quote)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}

	GetQuote(quote.Stock, quote.User)

	cmd := &log.UserCommand{
		Timestamp:      time.Now().UnixMilli(),
		Server:         "localhost",
		TransactionNum: 0, //for now
		Command:        "QUOTE",
		Username:       quote.User,
		StockSymbol:    quote.Stock,
	}
	log.CreateUserCommandsLog(cmd)

}

func displayHandler(w http.ResponseWriter, r *http.Request) {
	//just add the command to logs for now
	user := r.URL.Query().Get("user")

	cmd := &log.UserCommand{
		Timestamp:      time.Now().UnixMilli(),
		Server:         "localhost",
		TransactionNum: 0, //for now
		Command:        "DISPLAY_SUMMARY",
		Username:       user,
	}
	log.CreateUserCommandsLog(cmd)

}

func GetQuote(stock string, user string) float64 {
	command := stock + " " + user + " \n"
	requestTime := time.Now().UnixMilli()
	result := SendRequest(command)
	quoteServer := &log.QuoteServer{
		Timestamp:       requestTime,
		Server:          "localhost",
		TransactionNum:  0, //for now
		QuoteServerTime: int64(result.TimeStamp),
		Username:        user,
		StockSymbol:     result.Symbol,
		Price:           result.Price,
		Cryptokey:       result.Key,
	}
	log.CreateQuoteServerLog(quoteServer)
	return result.Price
}

func TransactionServerRequest(stock string, user string) *TransactionResult {
	command := stock + " " + user + " \n"
	return SendRequest(command)
}

func SendRequest(command string) *TransactionResult {
	//Establish Connection
	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		panic(err)
	}

	///Send Command
	_, err = connection.Write([]byte(command))
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	defer connection.Close()

	// Process Result
	result := strings.Split(strings.Split(string(buffer[:mLen]), "\n")[0], ",")
	Amount, err := strconv.ParseFloat(result[0], 64)
	time, err := strconv.Atoi(result[2])
	transactionResult := TransactionResult{Amount, result[1], result[2], time, result[4]}

	//Return Result
	return &transactionResult
}

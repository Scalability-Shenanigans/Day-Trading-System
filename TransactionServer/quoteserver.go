package main

import (
	"TransactionServer/cache"
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
	User           string `json:"user"`
	Stock          string `json:"stock"`
	TransactionNum int    `json:"transactionNum"`
}

type DisplaySummary struct {
	User           string `json:"user"`
	TransactionNum int    `json:"transactionNum"`
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

	GetQuote(quote.Stock, quote.User, quote.TransactionNum)

	cmd := &log.UserCommand{
		Timestamp:      time.Now().UnixMilli(),
		Server:         "localhost",
		TransactionNum: int64(quote.TransactionNum), //for now
		Command:        "QUOTE",
		Username:       quote.User,
		StockSymbol:    quote.Stock,
	}
	log.CreateUserCommandsLog(cmd)

}

func displayHandler(w http.ResponseWriter, r *http.Request) {
	//just add the command to logs for now
	var displaySummary DisplaySummary
	err := json.NewDecoder(r.Body).Decode(&displaySummary)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}

	user := displaySummary.User
	transactionNum := displaySummary.TransactionNum

	cmd := &log.UserCommand{
		Timestamp:      time.Now().UnixMilli(),
		Server:         "localhost",
		TransactionNum: int64(transactionNum),
		Command:        "DISPLAY_SUMMARY",
		Username:       user,
	}
	log.CreateUserCommandsLog(cmd)

}

func GetQuote(stock string, user string, transactionNum int) float64 {
	redisClient := cache.NewRedisClient()

	// First check cache
	ok, value := redisClient.Get(stock)
	if ok {
		fmt.Println("Cache hit for " + stock)
		return value // cache hit
	}

	command := stock + " " + user + " \n"
	requestTime := time.Now().UnixMilli()
	result := SendRequest(command)
	quoteServer := &log.QuoteServer{
		Timestamp:       requestTime,
		Server:          "localhost",
		TransactionNum:  int64(transactionNum),
		QuoteServerTime: int64(result.TimeStamp),
		Username:        user,
		StockSymbol:     result.Symbol,
		Price:           result.Price,
		Cryptokey:       result.Key,
	}
	log.CreateQuoteServerLog(quoteServer)

	// cache the price
	redisClient.Set(stock, strconv.FormatFloat(result.Price, 'E', -1, 64), 60*time.Second)
	return result.Price
}

func SendRequest(command string) *TransactionResult {
	serverPorts := [20]string{"4441", "4442", "4443", "4444", "4445", "4446", "4447", "4448", "4449", "4450", "4451", "4452", "4453", "4454", "4455", "4456", "4457", "4458", "4459", "4460"}
	for {
		//Establish Connection
		for _, port := range serverPorts {
			connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+port)
			if err != nil {
				fmt.Println("Error connecting to server:", err)
				time.Sleep(time.Second) // Wait for a second before retrying
				continue                // Retry connection
			}

			//Send Command
			_, err = connection.Write([]byte(command))
			if err != nil {
				fmt.Println("Error sending command:", err)
				err := connection.Close()
				if err != nil {
					continue
				}
				time.Sleep(time.Second) // Wait for a second before retrying
				continue                // Retry connection
			}

			// Read Response
			buffer := make([]byte, 1024)
			mLen, err := connection.Read(buffer)
			if err != nil {
				fmt.Println("Error reading response:", err)
				err := connection.Close()
				if err != nil {
					continue
				}
				time.Sleep(time.Second) // Wait for a second before retrying
				continue                // Retry connection
			}

			// Process Result
			result := strings.Split(strings.Split(string(buffer[:mLen]), "\n")[0], ",")
			Amount, err := strconv.ParseFloat(result[0], 64)
			timeStamp, err := strconv.Atoi(result[2])
			transactionResult := TransactionResult{Amount, result[1], result[2], timeStamp, result[4]}

			// Close connection and return result
			connection.Close()
			return &transactionResult
		}
	}
}

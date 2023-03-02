package main

import (
	"fmt"
	"net"
	"strings"
	"stringcov"
)

const (
	SERVER_HOST = "quoteserve.seng.uvic.ca"
	SERVER_PORT = "4444"
	SERVER_TYPE = "tcp"
)

type TransactionResult struct {
	Price 			float64  
	Symbol          string 
	Username        bool   
	TimeStamp       int    
	Key         	string    
}

func GetQuote(stock string, user string) float64 {
	command := stock + " " + user + " \n"
	result  := SendRequest(command)
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
	result := strings.Split(strings.Split(string(buffer[:mLen]), "\n")[0],",")
	Amount, err := strconv.ParseFloat(result[0], 64)
	time,err :=strconv.Atoi(result[2])
	transactionResult := TransactionResult{Amount,result[1],result[2],time,result[4]}

	//Return Result
	return transactionResult
}

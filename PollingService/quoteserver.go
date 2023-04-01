package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

const (
	SERVER_HOST = "quoteserver"
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

func TransactionServerRequest(stock string, user string) TransactionResult {
	command := stock + " " + user + " \n"
	return SendRequest(command)
}

func SendRequest(command string) TransactionResult {
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

	// fmt.Println("the command is " + command)
	command_split := strings.Split(command, " ")
	symbol := command_split[0]
	username := command_split[1]

	// Process Result
	result := strings.Split(strings.Split(string(buffer[:mLen]), "\n")[0], ",")
	Amount, err := strconv.ParseFloat(result[0], 64)
	time, err := strconv.Atoi(result[1])
	transactionResult := TransactionResult{Amount, symbol, username, time, result[2]}

	//Return Result
	return transactionResult
}

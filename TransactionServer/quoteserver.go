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
	Price int  
	Symbol          string 
	Username         bool   
	TimeStamp         int    
	Key         string    
}

func TransactionServerRequest(command string) string {
	return SendRequest(command)
}

func SendRequest(command string) *TransactionResult {
	//establish connection
	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		panic(err)
	}
	///send some data
	_, err = connection.Write([]byte("S asidhoias\n"))
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	defer connection.Close()
	result := strings.Split(strings.Split(string(buffer[:mLen]), "\n")[0],",")
	Amount, err := strconv.ParseFloat(result[0], 64)
	time,err :=strconv.Atoi(result[2])
	transactionResult := TransactionResult{Amount,result[1],result[2],time,result[4]}
	return transactionResult
}

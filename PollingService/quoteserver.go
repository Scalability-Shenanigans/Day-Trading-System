package main

import (
	"fmt"
	"net"
)

const (
	SERVER_HOST = "quoteserve.seng.uvic.ca"
	SERVER_PORT = "4444"
	SERVER_TYPE = "tcp"
)

func TransactionServerRequest(command string) string {
	return SendRequest(command)
}

func SendRequest(command string) string {
	//establish connection
	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		panic(err)
	}
	///send some data
	_, err = connection.Write([]byte(command))
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	defer connection.Close()
	return string(buffer[:mLen])
}

func GetQuote(user string, stock string) string {
	//
	return "32"
}

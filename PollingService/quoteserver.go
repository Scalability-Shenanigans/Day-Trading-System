package main

import (
	"fmt"
	"net"
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

func TransactionServerRequest(stock string, user string) TransactionResult {
	command := stock + " " + user + " \n"
	return SendRequest(command)
}

func SendRequest(command string) TransactionResult {
	transactionResult := TransactionResult{}
	serverPorts := [20]string{"4441", "4442", "4443", "4444", "4445", "4446", "4447", "4448", "4449", "4450", "4451", "4452", "4453", "4454", "4455", "4456", "4457", "4458", "4459", "4460"}
	for _, port := range serverPorts {
		//Establish Connection
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

		buffer := make([]byte, 1024)
		mLen, err := connection.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
		defer connection.Close()

		// Process Result
		result := strings.Split(strings.Split(string(buffer[:mLen]), "\n")[0], ",")
		Amount, err := strconv.ParseFloat(result[0], 64)
		timeStamp, err := strconv.Atoi(result[2])
		transactionResult := TransactionResult{Amount, result[1], result[2], timeStamp, result[4]}

		//Return Result
		return transactionResult
	}
	return transactionResult
}

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const SERVER_URL = "http://localhost:5100/"

func sendRequest(endpoint string, data map[string]interface{}) (string, error) {
	payload, _ := json.Marshal(data)
	response, err := http.Post(SERVER_URL+endpoint, "application/json", bytes.NewBuffer(payload))

	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		return "Success", nil
	}
	return "", fmt.Errorf("request failed with status code: %d", response.StatusCode)
}

func lineProcessor(line string) {
	regex := regexp.MustCompile(`\[(\d*)\]`)
	matches := regex.FindStringSubmatch(line)

	if len(matches) < 2 {
		fmt.Println("Invalid input")
		return
	}

	transactionNum, _ := strconv.Atoi(matches[1])
	fmt.Println(line)

	lineSplit := strings.Split(line, " ")
	args := strings.Split(lineSplit[1], ",")

	command := args[0]
	data := map[string]interface{}{
		"transactionNum": transactionNum,
	}

	switch command {
	case "BUY":
		amount, _ := strconv.ParseFloat(args[3], 64)
		data["user"] = args[1]
		data["stock"] = args[2]
		data["amount"] = amount
		fmt.Println(data)
		fmt.Println(sendRequest("buy", data))
	case "COMMIT_BUY":
		data["user"] = args[1]
		fmt.Println(data)
		fmt.Println(sendRequest("commitBuy", data))
	case "CANCEL_BUY":
		data["user"] = args[1]
		fmt.Println(data)
		fmt.Println(sendRequest("cancelBuy", data))
	case "CANCEL_SET_BUY":
		data["user"] = args[1]
		data["stock"] = args[2]
		fmt.Println(data)
		fmt.Println(sendRequest("cancelSetBuy", data))
	case "SET_BUY_AMOUNT":
		amount, _ := strconv.ParseFloat(args[3], 64)
		data["user"] = args[1]
		data["stock"] = args[2]
		data["amount"] = amount
		fmt.Println(data)
		fmt.Println(sendRequest("setBuyAmount", data))
	case "SET_BUY_TRIGGER":
		amount, _ := strconv.ParseFloat(args[3], 64)
		data["user"] = args[1]
		data["stock"] = args[2]
		data["amount"] = amount
		fmt.Println(data)
		fmt.Println(sendRequest("setBuyTrigger", data))
	case "ADD":
		amount, _ := strconv.ParseFloat(args[2], 64)
		data["user"] = args[1]
		data["amount"] = amount
		fmt.Println(data)
		fmt.Println(sendRequest("add", data))
	case "DUMPLOG":
		data["filename"] = args[1]
		fmt.Println(data)
		fmt.Println("ended at: ", time.Now())
		fmt.Println(sendRequest("dumplog", data))
	case "SELL":
		amount, _ := strconv.ParseFloat(args[3], 64)
		data["user"] = args[1]
		data["stock"] = args[2]
		data["amount"] = amount
		fmt.Println(data)
		fmt.Println(sendRequest("sell", data))
	case "COMMIT_SELL":
		data["user"] = args[1]
		fmt.Println(data)
		fmt.Println(sendRequest("commitSell", data))
	case "CANCEL_SET_SELL":
		data["user"] = args[1]
		data["stock"] = args[2]
		fmt.Println(data)
		fmt.Println(sendRequest("cancelSetSell", data))
	case "CANCEL_SELL":
		data["user"] = args[1]
		fmt.Println(sendRequest("cancelSell", data))
	case "SET_SELL_AMOUNT":
		amount, _ := strconv.ParseFloat(args[3], 64)
		data["user"] = args[1]
		data["stock"] = args[2]
		data["amount"] = amount
		fmt.Println(data)
		fmt.Println(sendRequest("setSellAmount", data))
	case "SET_SELL_TRIGGER":
		amount, _ := strconv.ParseFloat(args[3], 64)
		data["user"] = args[1]
		data["stock"] = args[2]
		data["amount"] = amount
		fmt.Println(data)
		fmt.Println(sendRequest("setSellTrigger", data))
	case "QUOTE":
		data["user"] = args[1]
		data["stock"] = args[2]
		fmt.Println(data)
		fmt.Println(sendRequest("quote", data))
	case "DISPLAY_SUMMARY":
		data["user"] = args[1]
		fmt.Println(data)
		fmt.Println(sendRequest("display", data))
	default:
		fmt.Println("Invalid command")
		return
	}
}

func main() {
	var choice string
	fmt.Print("mode: manual, dumplog, dbwipe, or automatic? ")
	fmt.Scanln(&choice)

	switch choice {
	case "manual":
		for {
			var line string
			fmt.Print("enter command: ")
			fmt.Scanln(&line)
			if line == "quit" {
				break
			} else {
				lineProcessor(line)
			}
		}
	case "dumplog":
		data := map[string]interface{}{
			"filename":       "logfile.xml",
			"transactionNum": 100,
		}
		fmt.Println(data)
		resp, err := sendRequest("dumplog", data)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println(resp)
		}
	case "dbwipe":
		resp, err := sendRequest("dbwipe", map[string]interface{}{})
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println(resp)
		}
	case "automatic":
		startTime := time.Now()
		fmt.Println("Starting at:", startTime)
		file, err := os.Open("workload_files/user100.txt")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			lineProcessor(line)
		}
		endTime := time.Now()
		difference := endTime.Sub(startTime)
		fmt.Println("The difference is", difference)
	default:
		fmt.Println("Invalid choice")
	}
}

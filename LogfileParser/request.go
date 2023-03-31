package main

import (
	"net/http"
	"net/url"
	"strings"
)

const (
	TRANSACTION_SERVER_URL = "http://localhost:8080/"
)

func DUMPLOG() {
	PostRequest("dumplog", url.Values{
		"filename":       {"logfile.xml"},
		"transactionNum": {"10000000000000000"},
	})
}
func ExecuteCommands(commands []string) {
	for _, command := range commands {
		HandleRequest(command)
	}
}
func HandleRequest(commandString string) {
	args := strings.Split(string(commandString), ",")
	command := args[1]
	switch {
	case command == "BUY":
		PostRequest("buy", url.Values{
			"user":           {args[2]},
			"stock":          {args[3]},
			"amount":         {args[4]},
			"transactionNum": {args[0]},
		})
	case command == "COMMIT_BUY":
		PostRequest("commitBuy", url.Values{
			"user":           {args[2]},
			"transactionNum": {args[0]},
		})
	case command == "CANCEL_BUY":
		PostRequest("cancelBuy", url.Values{
			"user":           {args[2]},
			"transactionNum": {args[0]},
		})
	case command == "CANCEL_SET_BUY":
		PostRequest("cancelSetBuy", url.Values{
			"user":           {args[2]},
			"stock":          {args[3]},
			"transactionNum": {args[0]},
		})
	case command == "SET_BUY_AMOUNT":
		PostRequest("setBuyAmount", url.Values{
			"user":           {args[2]},
			"stock":          {args[3]},
			"amount":         {args[4]},
			"transactionNum": {args[0]},
		})
	case command == "SET_BUY_TRIGGER":
		PostRequest("setBuyTrigger", url.Values{
			"user":           {args[2]},
			"stock":          {args[3]},
			"amount":         {args[4]},
			"transactionNum": {args[0]},
		})
	case command == "ADD":
		PostRequest("add", url.Values{
			"user":           {args[2]},
			"amount":         {args[3]},
			"transactionNum": {args[0]},
		})
	case command == "SELL":
		PostRequest("sell", url.Values{
			"user":           {args[2]},
			"stock":          {args[3]},
			"amount":         {args[4]},
			"transactionNum": {args[0]},
		})
	case command == "COMMIT_SELL":
		PostRequest("commitSell", url.Values{
			"user":           {args[2]},
			"transactionNum": {args[0]},
		})
	case command == "CANCEL_SELL":
		PostRequest("cancelSell", url.Values{
			"user":           {args[2]},
			"transactionNum": {args[0]},
		})
	case command == "CANCEL_SET_SELL":
		PostRequest("cancelSetSell", url.Values{
			"user":           {args[2]},
			"stock":          {args[3]},
			"transactionNum": {args[0]},
		})
	case command == "SET_SELL_AMOUNT":
		PostRequest("setSellAmount", url.Values{
			"user":           {args[2]},
			"stock":          {args[3]},
			"amount":         {args[4]},
			"transactionNum": {args[0]},
		})
	case command == "SET_SELL_TRIGGER":
		PostRequest("setSellTrigger", url.Values{
			"user":           {args[2]},
			"stock":          {args[3]},
			"amount":         {args[4]},
			"transactionNum": {args[0]},
		})
	case command == "QUOTE":
		PostRequest("quote", url.Values{
			"user":           {args[2]},
			"stock":          {args[3]},
			"transactionNum": {args[0]},
		})
	case command == "DISPLAY_SUMMARY":
		PostRequest("display", url.Values{
			"user":           {args[2]},
			"transactionNum": {args[0]},
		})
	}
}
func PostRequest(command string, body url.Values) {
	resp, err := http.PostForm(TRANSACTION_SERVER_URL+command,
		body)

	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}

	if err != nil {
		panic(err)
	}
}

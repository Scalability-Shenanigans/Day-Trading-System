package main

import "strings"

type TradingCommand struct {
	Command     string
	User        string
	stockSymbol string
	fileName    string
	Amount      float64
}

func CreateCommand(commandString string) TradingCommand {
	command := strings.Split(commandString, ",")
	tradingCommand := TradingCommand{Command: command[0]}
	return tradingCommand
}

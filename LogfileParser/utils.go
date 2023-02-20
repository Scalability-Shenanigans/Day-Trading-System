package main

import (
	"fmt"
	"os"
	"strings"
)

func CheckIfArgumentsValid(arguments []string) bool {
	if len(arguments) == 1 {
		fmt.Println("No file was provided")
		return true
	}
	return false
}

func ReadContents(arguments []string) []string {
	fileName := arguments[1]
	contents, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("File reading error", err)
		panic(err)
	}
	return GetCommandStrings(strings.Split(string(contents), "\n"))
}

func GetCommandStrings(contents []string) []string {
	commands := []string{}
	for _, v := range contents {
		var commandString string
		var commandNumber int
		_, err := fmt.Sscanf(v, "[%d] %s", &commandNumber, &commandString)
		if err != nil {
			panic(err.Error() + "\n Error while parsing commands from logfile")
		}
		commands = append(commands, commandString)
	}
	return commands
}

// func GetCommands(contents []string) []TradingCommand {
// 	commands := []TradingCommand{}
// 	fmt.Println("Contents of file:\n ")
// 	for i := 0; i < len(contents); i++ {
// 		var commandString string
// 		var commandNumber int
// 		_, err := fmt.Sscanf(contents[i], "[%d] %s", &commandNumber, &commandString)
// 		if err != nil {
// 			panic(err.Error() + "\n Error while parsing commands from logfile")
// 		}
// 		commands = append(commands, CreateCommand(commandString))
// 	}

// 	return commands
// }

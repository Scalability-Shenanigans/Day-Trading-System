package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func CheckIfArgumentsValid(arguments []string) bool {
	if len(arguments) == 1 {
		fmt.Println("No file was provided")
		return true
	}
	return false
}

func ReadContents(arguments []string) string {
	fileName := arguments[1]
	contents, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("File reading error", err)
		panic(err)
	}
	return string(contents)
}

func ProcessContents(contents string) []string {
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
		commands = append(commands, strconv.Itoa(commandNumber)+","+commandString)
	}
	return commands
}

func MakeUserMap(content []string) map[string][]string {
	userMap := make(map[string][]string)
	for _, command := range content {
		arguments := strings.Split(string(command), ",")
		if arguments[1] == "DUMPLOG" {
			continue
		}
		user := arguments[2]
		userMap[user] = append(userMap[user], command)
	}
	return userMap
}

package utils

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
		return nil
	}
	return strings.Split(string(contents), "\n")
}

func GetCommands(contents []string) []string {
	commands := []string{}
	fmt.Println("Contents of file:\n ")
	for i := 0; i < len(contents); i++ {
		commands = append(commands, strings.Split(contents[i], "]")[0])
	}
	return commands
}

package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if checkIfArgumentsValid() {
		return
	}
	contents := readContents()
	commands := getCommands(contents)
	fmt.Println(commands)
}

func checkIfArgumentsValid() bool {
	if len(os.Args) == 1 {
		fmt.Println("No file was provided")
		return true
	}
	return false
}

func readContents() []string {
	fileName := os.Args[1]
	contents, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("File reading error", err)
		return nil
	}
	return strings.Split(string(contents), "\n")
}

func getCommands(contents []string) []string {

	fmt.Println("Contents of file:\n ")
	for i := 0; i < len(contents); i++ {
		fmt.Println(contents[i])
	}
}

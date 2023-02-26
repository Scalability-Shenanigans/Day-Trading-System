package main

import (
	"os"
)

func main() {
	if CheckIfArgumentsValid(os.Args) {
		return
	}
	contents := ReadContents(os.Args)
	ProduceCommands(contents)
}

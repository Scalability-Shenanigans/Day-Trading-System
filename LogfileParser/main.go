package main

import (
	"fmt"
	"os"
)

func main() {
	if CheckIfArgumentsValid(os.Args) {
		return
	}
	contents := ReadContents(os.Args)
	fmt.Println(contents)
}

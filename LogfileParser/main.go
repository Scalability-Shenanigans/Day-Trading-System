package main

import (
	"os"
	"sync"
)

var wg sync.WaitGroup

func main() {
	if CheckIfArgumentsValid(os.Args) {
		return
	}
	userMap := MakeUserMap(ProcessContents(ReadContents(os.Args)))
	for _, command := range userMap {
		wg.Add(1)
		go ExecuteCommands(command)
	}
	wg.Wait()
	DUMPLOG()
}

package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

type Command struct {
	XMLName        xml.Name `xml:"userCommand"`
	TimeStamp      int      `xml:"timestamp"`
	TransactionNum int      `xml:"transactionNum"`
	Command        string   `xml:"command"`
	UserName       string   `xml:"username"`
	Funds          int      `xml:"funds"`
}

type CommandArray struct {
	Commands []Command
}

type Company struct {
	XMLName      xml.Name `xml:"log"`
	userCommands CommandArray
}

func (s *CommandArray) AddCommand(timestamp int, transactionNum int, command string, username string, funds int) {
	staffRecord := Command{TimeStamp: timestamp, TransactionNum: transactionNum, Command: command, UserName: username, Funds: funds}

	s.Commands = append(s.Commands, staffRecord)
}

func Kundi() {

	v := Company{}

	// put a for loop here to add more data
	// this example will just add 2 rows of data.

	v.userCommands.AddStaff(103, 1234, "ADD", "ADADS", 150)

	// sanity check - display on screen
	xmlString, err := xml.MarshalIndent(v, "", "    ")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%s \n", string(xmlString))

	// everything ok now, write to file.
	filename := "newstaffs.xml"
	file, _ := os.Create(filename)

	xmlWriter := io.Writer(file)

	enc := xml.NewEncoder(xmlWriter)
	enc.Indent("  ", "    ")
	if err := enc.Encode(v); err != nil {
		fmt.Printf("error: %v\n", err)
	}

}
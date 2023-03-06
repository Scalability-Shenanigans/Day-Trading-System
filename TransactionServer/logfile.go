package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)
type DebugEvent struct {
	TimeStamp      int      	`xml:"timestamp"`
	Server         string       `xml:"server"`
	TransactionNum int      	`xml:"transactionNum"`
	Command        string   	`xml:"command"`
	UserName       string       `xml:"username,omitempty"`
	StockSymbol	   string       `xml:"stocksymbol,omitempty"`
	FileName	   string       `xml:"filename,omitempty"`
	Funds          float64      `xml:"funds,omitempty"`
	DebugMessage   string		`xml:"debugMessage,omitempty"`
}

type ErrorEvent struct {
	TimeStamp      int      	`xml:"timestamp"`
	Server         string       `xml:"server"`
	TransactionNum int      	`xml:"transactionNum"`
	Command        string   	`xml:"command"`
	UserName       string       `xml:"username,omitempty"`
	StockSymbol	   string       `xml:"stocksymbol,omitempty"`
	FileName	   string       `xml:"filename,omitempty"`
	Funds          float64      `xml:"funds,omitempty"`
	ErrorMessage   string		`xml:"errorMessage,omitempty"`
}

type SystemEvent struct {
	TimeStamp      int      	`xml:"timestamp"`
	Server         string       `xml:"server"`
	TransactionNum int      	`xml:"transactionNum"`
	Command        string   	`xml:"command"`
	UserName       string       `xml:"username"`
	StockSymbol	   string       `xml:"stocksymbol,omitempty"`
	FileName	   string       `xml:"filename,omitempty"`
	Funds          float64      `xml:"funds,omitempty"`
}

type QuoteServer struct {
	TimeStamp      int      	`xml:"timestamp"`
	Server         string       `xml:"server"`
	TransactionNum int      	`xml:"transactionNum"`
	Price          float64      `xml:"price"`
	StockSymbol	   string       `xml:"stocksymbol"`
	UserName       string       `xml:"username"`
	CryptoKey      string   	`xml:"cryptokey"`
	
}

type AccountTransaction struct {	
	TimeStamp      int      	`xml:"timestamp"`
	Server         string       `xml:"server"`
	TransactionNum int      	`xml:"transactionNum"`
	Action         string   	`xml:"action"`
	UserName       string       `xml:"username"`
	Funds          float64      `xml:"funds"`
}

type UserCommand struct {	
	TimeStamp      int      	`xml:"timestamp"`
	Server         string       `xml:"server"`
	TransactionNum int      	`xml:"transactionNum"`
	Command        string   	`xml:"command"`
	UserName       string       `xml:"username"`
	StockSymbol	   string       `xml:"stocksymbol,omitempty"`
	FileName	   string       `xml:"filename,omitempty"`
	Funds          float64      `xml:"funds,omitempty"`
}

type Log struct {
	XMLName      	    xml.Name 				`xml:"log"`
	UserCommands 	    []UserCommand  			`xml:"userCommand,omitempty"`
	AccountTransactions []AccountTransaction	`xml:"accountTransaction,omitempty"`
	QuoteServers		[]QuoteServer			`xml:"quoteServe,omitemptyr"`
	SystemEvents        []SystemEvent 			`xml:"systemEvent,omitempty"`
	ErrorEvents         []ErrorEvent			`xml:"errorEvent,omitempty"`
	DebugEvents         []DebugEvent 			`xml:"debugEvent,omitempty"`
}

func (l *Log) AddUserCommand(userCommand UserCommand) {
	l.UserCommands = append(l.UserCommands, userCommand)
}

func (l *Log) AddAccountTransaction(accountTransaction AccountTransaction) {
	l.AccountTransactions = append(l.AccountTransactions, accountTransaction)
}

func (l *Log) AddQuoteServer(quoteServer QuoteServer) {
	l.QuoteServers = append(l.QuoteServers, quoteServer)
}

func (l *Log) AddSystemEvent(systemEvent SystemEvent) {
	l.SystemEvents = append(l.SystemEvents, systemEvent)
}

func (l *Log) AddDebugEvent(debugEvent DebugEvent) {
	l.DebugEvents = append(l.DebugEvents, debugEvent)
}

func (l *Log) AddErrorEvent(errorEvent ErrorEvent) {
	l.ErrorEvents = append(l.ErrorEvents, errorEvent)
}

func WriteXML(filename string) {

	logfile := Log{}

	// sanity check - display on screen
	// xmlString, err := xml.MarshalIndent(logfile, "", "    ")

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Printf("%s \n", string(xmlString))

	// Create File.
	file, _ := os.Create(filename)
	xmlWriter := io.Writer(file)
	enc := xml.NewEncoder(xmlWriter)
	enc.Indent("  ", "    ")

	// Write to file.
	if err := enc.Encode(logfile); err != nil {
		fmt.Printf("error: %v\n", err)
	}

}
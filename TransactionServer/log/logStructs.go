package log

import "encoding/xml"

type UserCommand struct {
	XMLName        xml.Name `xml:"userCommand"`
	Timestamp      int64    `xml:"timestamp"`
	Server         string   `xml:"server"`
	TransactionNum int64    `xml:"transactionNum"`
	Command        string   `xml:"command,omitempty"`
	Username       string   `xml:"username,omitempty"`
	StockSymbol    string   `xml:"stockSymbol,omitempty"`
	Filename       string   `xml:"filename,omitempty"`
	Funds          float64  `xml:"funds,omitempty"`
}

type QuoteServer struct {
	XMLName         xml.Name `xml:"quoteServer"`
	Timestamp       int64    `xml:"timestamp"`
	Server          string   `xml:"server"`
	TransactionNum  int64    `xml:"transactionNum"`
	Price           float64  `xml:"price"`
	StockSymbol     string   `xml:"stockSymbol"`
	Username        string   `xml:"username"`
	QuoteServerTime int64    `xml:"quoteServerTime"`
	Cryptokey       string   `xml:"cryptokey"`
}

type AccountTransaction struct {
	XMLName        xml.Name `xml:"accountTransaction"`
	Timestamp      int64    `xml:"timestamp"`
	Server         string   `xml:"server"`
	TransactionNum int64    `xml:"transactionNum"`
	Action         string   `xml:"action,omitempty"`
	Username       string   `xml:"username,omitempty"`
	Funds          float64  `xml:"funds,omitempty"`
}

type SystemEvent struct {
	XMLName        xml.Name `xml:"systemEvent"`
	Timestamp      int64    `xml:"timestamp"`
	Server         string   `xml:"server"`
	TransactionNum int64    `xml:"transactionNum"`
	Command        string   `xml:"command"`
	Username       string   `xml:"username"`
	StockSymbol    string   `xml:"stockSymbol,omitempty"`
	Filename       string   `xml:"filename,omitempty"`
	Funds          float64  `xml:"funds,omitempty"`
}

type ErrorEvent struct {
	XMLName        xml.Name `xml:"errorEvent"`
	Timestamp      int64    `xml:"timestamp"`
	Server         string   `xml:"server"`
	TransactionNum int64    `xml:"transactionNum"`
	Command        string   `xml:"command"`
	Username       string   `xml:"username,omitempty"`
	StockSymbol    string   `xml:"stockSymbol,omitempty"`
	Filename       string   `xml:"filename,omitempty"`
	Funds          float64  `xml:"funds,omitempty"`
	ErrorMessage   string   `xml:"errorMessage,omitempty"`
}

type Debug struct {
	XMLName        xml.Name `xml:"debug"`
	Timestamp      int64    `xml:"timestamp"`
	Server         string   `xml:"server"`
	TransactionNum int64    `xml:"transactionNum"`
	Command        string   `xml:"command"`
	Username       string   `xml:"username,omitempty"`
	StockSymbol    string   `xml:"stockSymbol,omitempty"`
	Filename       string   `xml:"filename,omitempty"`
	Funds          float64  `xml:"funds,omitempty"`
	DebugMessage   string   `xml:"debugMessage,omitempty"`
}

type Log struct {
	XMLName             xml.Name             `xml:"log"`
	UserCommands        []UserCommand        `xml:"userCommand,omitempty"`
	AccountTransactions []AccountTransaction `xml:"accountTransaction,omitempty"`
	QuoteServers        []QuoteServer        `xml:"quoteServer,omitempty"`
	SystemEvents        []SystemEvent        `xml:"systemEvent,omitempty"`
	ErrorEvents         []ErrorEvent         `xml:"errorEvent,omitempty"`
	DebugEvents         []Debug              `xml:"debug,omitempty"`
}

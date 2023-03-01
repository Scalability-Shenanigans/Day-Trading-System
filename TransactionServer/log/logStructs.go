package log

import "encoding/xml"

type UserCommand struct {
	XMLName     xml.Name `xml:"userCommand"`
	Command     string   `xml:"command"`
	Username    string   `xml:"username"`
	StockSymbol string   `xml:"stockSymbol"`
	Filename    string   `xml:"filename"`
	Funds       float64  `xml:"funds"`
}

type QuoteServer struct {
	XMLName         xml.Name `xml:"quoteServer"`
	Price           float64  `xml:"price"`
	StockSymbol     string   `xml:"stockSymbol"`
	Username        string   `xml:"username"`
	QuoteServerTime int64    `xml:"quoteServerTime"`
	Cryptokey       string   `xml:"cryptokey"`
}

type AccountTransaction struct {
	XMLName  xml.Name `xml:"accountTransaction"`
	Action   string   `xml:"action"`
	Username string   `xml:"username"`
	Funds    float64  `xml:"funds"`
}

type SystemEvent struct {
	XMLName     xml.Name `xml:"systemEvent"`
	Command     string   `xml:"command"`
	Username    string   `xml:"username"`
	StockSymbol string   `xml:"stockSymbol,omitempty"`
	Filename    string   `xml:"filename,omitempty"`
	Funds       float64  `xml:"funds,omitempty"`
}

type ErrorEvent struct {
	XMLName      xml.Name `xml:"errorEvent"`
	Command      string   `xml:"command"`
	Username     string   `xml:"username,omitempty"`
	StockSymbol  string   `xml:"stockSymbol,omitempty"`
	Filename     string   `xml:"filename,omitempty"`
	Funds        float64  `xml:"funds,omitempty"`
	ErrorMessage string   `xml:"errorMessage,omitempty"`
}

type Debug struct {
	XMLName      xml.Name `xml:"debug"`
	Command      string   `xml:"command"`
	Username     string   `xml:"username,omitempty"`
	StockSymbol  string   `xml:"stockSymbol,omitempty"`
	Filename     string   `xml:"filename,omitempty"`
	Funds        float64  `xml:"funds,omitempty"`
	DebugMessage string   `xml:"debugMessage,omitempty"`
}

type Log struct {
	Timestamp          int64               `xml:"timestamp"`
	Server             string              `xml:"server"`
	TransactionNum     int64               `xml:"transactionNum"`
	UserCommand        *UserCommand        `xml:"userCommand,omitempty"`
	QuoteServer        *QuoteServer        `xml:"quoteServer,omitempty"`
	AccountTransaction *AccountTransaction `xml:"accountTransaction,omitempty"`
	SystemEvent        *SystemEvent        `xml:"systemEvent,omitempty"`
	ErrorEvent         *ErrorEvent         `xml:"errorEvent,omitempty"`
	Debug              *Debug              `xml:"debug,omitempty"`
}

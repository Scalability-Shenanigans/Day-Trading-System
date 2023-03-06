package log

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Dumplog struct {
	Filename       string `json:"filename"`
	TransactionNum int    `json:"transactionNum"`
}

// var mongoURI = "mongodb://localhost:5000" //for local testing
var mongoURI = "mongodb://localhost:27017" //use this for when everything is containerized
var client *mongo.Client
var logs *mongo.Collection

func InitLogDBConnection() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))

	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("logs connected to mongodb")

	db := client.Database("DayTrading")
	logs = db.Collection("Logs")
}

func CreateUserCommandsLog(cmd *UserCommand) {

	cmd.XMLName = xml.Name{Local: "UserCommand"}

	res, err := logs.InsertOne(context.TODO(), cmd)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res.InsertedID)
	}

}

func CreateAccountTransactionLog(cmd *AccountTransaction) {

	cmd.XMLName = xml.Name{Local: "AccountTransaction"}

	res, err := logs.InsertOne(context.TODO(), cmd)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res.InsertedID)
	}
}

func DumplogHandler(w http.ResponseWriter, r *http.Request) {
	var dumplog Dumplog

	err := json.NewDecoder(r.Body).Decode(&dumplog)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Bad Request")
		return
	}

	cursor, err := logs.Find(context.TODO(), bson.D{})

	if err != nil {
		fmt.Println(err)
	}

	var userCommands []UserCommand
	var accountTransactions []AccountTransaction
	var quoteServers []QuoteServer
	var systemEvents []SystemEvent
	var errorEvents []ErrorEvent
	var debugs []Debug

	for cursor.Next(context.Background()) {
		var log bson.M
		err := cursor.Decode(&log)
		if err != nil {
			fmt.Println(err)
		}

		xmlname := log["xmlname"].(bson.M)
		local := xmlname["local"].(string)

		switch local {
		case "UserCommand":
			var userCommand UserCommand
			data, _ := bson.Marshal(log)
			err := bson.Unmarshal(data, &userCommand)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("UserCommand: %+v\n", userCommand)
			userCommands = append(userCommands, userCommand)
		case "AccountTransaction":
			var accountTransaction AccountTransaction
			data, _ := bson.Marshal(log)
			err := bson.Unmarshal(data, &accountTransaction)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("AccountTransaction: %+v\n", accountTransaction)
			accountTransactions = append(accountTransactions, accountTransaction)
		case "SystemEvent":
			var systemEvent SystemEvent
			data, _ := bson.Marshal(log)
			err := bson.Unmarshal(data, &systemEvent)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("SystemEvent: %+v\n", systemEvent)
			systemEvents = append(systemEvents, systemEvent)
		case "ErrorEvent":
			var errorEvent ErrorEvent
			data, _ := bson.Marshal(log)
			err := bson.Unmarshal(data, &errorEvent)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("ErrorEvent: %+v\n", errorEvent)
			errorEvents = append(errorEvents, errorEvent)
		case "Debug":
			var debug Debug
			data, _ := bson.Marshal(log)
			err := bson.Unmarshal(data, &debug)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("Debug: %+v\n", debug)
			debugs = append(debugs, debug)
		case "QuoteServer":
			var quoteServer QuoteServer
			data, _ := bson.Marshal(log)
			err := bson.Unmarshal(data, &quoteServer)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("QuoteServer: %+v\n", quoteServer)
			quoteServers = append(quoteServers, quoteServer)
		default:
			fmt.Println("Unknown document type")
		}
	}

	if err := cursor.Err(); err != nil {
		fmt.Println(err)
	}

	allLogs := Log{
		UserCommands:        userCommands,
		AccountTransactions: accountTransactions,
	}

	// create the XML file
	file, err := os.Create(dumplog.Filename + ".xml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString(xml.Header)

	// write the XML to the file
	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	if err := encoder.Encode(allLogs); err != nil {
		panic(err)
	}

}

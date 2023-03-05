package log

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

func CreateUserCommandsLog(cmd *UserCommand, transactionNum int64) {

	log := &Log{
		Timestamp:      time.Now().UnixNano(),
		Server:         "localhost",
		TransactionNum: transactionNum,
		UserCommand:    cmd,
	}

	res, err := logs.InsertOne(context.TODO(), log)

	if err != nil {
		fmt.Println("Failed to insert UserCommand log")
	} else {
		fmt.Println(res.InsertedID)
	}

}

func CreateAccountTransactionLog(cmd *AccountTransaction, transactionNum int64) {

	log := &Log{
		Timestamp:          time.Now().UnixNano(),
		Server:             "localhost",
		TransactionNum:     transactionNum,
		AccountTransaction: cmd,
	}

	res, err := logs.InsertOne(context.TODO(), log)

	if err != nil {
		fmt.Println("Failed to insert AccountTransaction log")
	} else {
		fmt.Println(res.InsertedID)
	}
}

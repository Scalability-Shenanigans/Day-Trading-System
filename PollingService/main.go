package main

import (
	"PollingService/db"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var mongoURI = "mongodb://localhost:5000" //for local testing
var mongoURI = "mongodb://localhost:27017" //use this for when everything is containerized
var client *mongo.Client
var accounts *mongo.Collection
var triggeredBuyAmountOrders *mongo.Collection

func InitConnection() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))

	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("connected to mongodb")

	db := client.Database("DayTrading")
	accounts = db.Collection("Accounts")
	triggeredBuyAmountOrders = db.Collection("TriggeredBuyAmountOrders")
}

func UpdateStockHolding(user string, stock string, amount int) bool {
	filter := bson.M{"user": user}
	var result db.Account
	accounts.FindOne(context.TODO(), filter).Decode(&result)

	stockToChange := db.StockHolding{
		Stock:  stock,
		Amount: 0,
	}

	var insertNewEntry = true

	// go through existing stocks to see if there is an existing holding for the stock
	for index, value := range result.Stocks {
		if value.Stock == stock {
			stockToChange = value
			insertNewEntry = false

			if stockToChange.Amount+amount < 0 {
				fmt.Println("not enough stock for change")
				return false
			}

			result.Stocks[index].Amount += amount
		}
	}

	// everything below here only matters if its a new stock (except for the return)
	// if stock is decreasing make sure enough is held
	if stockToChange.Amount+amount < 0 {
		fmt.Println("Not enough stock held for change")
		return false
	}
	stockToChange.Amount += amount

	// if this is a new stock to the portfolio add a new entry
	if insertNewEntry {
		result.Stocks = append(result.Stocks, stockToChange)
	}

	accounts.ReplaceOne(context.TODO(), filter, result)
	return true
}

func pollingFunction() {
	for {
		fmt.Println("Polling active triggeredBuyAmountOrders")

		cur, err := triggeredBuyAmountOrders.Find(context.Background(), bson.D{})
		fmt.Println("Fetching active triggeredBuyAmountOrders")
		if err != nil {
			log.Fatal(err)
		}
		for cur.Next(context.Background()) {
			var triggeredBuyAmountOrder db.TriggeredBuyAmountOrder
			if err := cur.Decode(&triggeredBuyAmountOrder); err != nil {
				fmt.Println("Error decoding triggeredBuyAmountOrder")
				log.Fatal(err)
			}
			fmt.Println(triggeredBuyAmountOrder)

			results := TransactionServerRequest(cur.Stock+" "+cur.User)
			curStockPrice := results.Price

			if triggeredBuyAmountOrder.Price >= curStockPrice { // If trigger price is >= stock price execute order

				UpdateStockHolding(triggeredBuyAmountOrder.User, triggeredBuyAmountOrder.Stock,
					int(triggeredBuyAmountOrder.Amount*triggeredBuyAmountOrder.Price))

				filter := bson.M{"user": triggeredBuyAmountOrder.User, "stock": triggeredBuyAmountOrder.Stock}
				_, err = triggeredBuyAmountOrders.DeleteOne(context.TODO(), filter)
				if err != nil {
					log.Fatal(err)
				}

				// Need to log transaction here for the log dump

				fmt.Printf("Completeted and Deleted TriggeredBuyAmountOrder for User %s, Stock %s\n",
					triggeredBuyAmountOrder.User, triggeredBuyAmountOrder.Stock)
			}
		}

		time.Sleep(time.Minute)
	}
}

func main() {
	InitConnection()
	pollingFunction()
}

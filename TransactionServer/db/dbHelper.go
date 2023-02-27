package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoURI = "mongodb://localhost:5000" //for local testing
// var mongoURI = "mongodb://db:27017" //use this for when everything is containerized
var client *mongo.Client
var accounts *mongo.Collection
var transactions *mongo.Collection
var buyOrders *mongo.Collection
var sellOrders *mongo.Collection

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
	transactions = db.Collection("PendingTransactions")
	buyOrders = db.Collection("BuyOrders")
	sellOrders = db.Collection("SellOrders")
}

func CreateAccount(user string, initial_balance int) {
	new_account := Account{
		User:    user,
		Balance: initial_balance,
	}
	res, err := accounts.InsertOne(context.TODO(), new_account)

	if err != nil {
		fmt.Println("Failed to create account")
	} else {
		fmt.Println(res.InsertedID)
	}

}

func UpdateBalance(amount int, user string) bool {
	filter := bson.M{"user": user}
	var result Account
	err := accounts.FindOne(context.TODO(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		result = Account{
			User: user,
		}
	}
	if amount < 0 && (result.Balance+amount) < 0 {
		fmt.Println("ERROR: funds will go below 0")
		return false
	}

	result.Balance += amount
	opts := options.Replace().SetUpsert(true)
	accounts.ReplaceOne(context.TODO(), filter, result, opts)
	fmt.Println("new balance set")
	return true

}

func UpdateStockHolding(user string, stock string, amount int) bool {
	filter := bson.M{"user": user}
	var result Account
	accounts.FindOne(context.TODO(), filter).Decode(&result)

	stockToChange := StockHolding{
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

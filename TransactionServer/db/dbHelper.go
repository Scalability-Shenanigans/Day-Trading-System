package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var mongoURI = "mongodb://localhost:5000" //for local testing
var mongoURI = "mongodb://localhost:27017" //use this for when everything is containerized
var client *mongo.Client
var accounts *mongo.Collection
var transactions *mongo.Collection
var buyOrders *mongo.Collection
var buyAmountOrders *mongo.Collection
var sellAmountOrders *mongo.Collection
var triggeredBuyAmountOrders *mongo.Collection
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
	buyAmountOrders = db.Collection("BuyAmountOrders")
	triggeredBuyAmountOrders = db.Collection("TriggeredBuyAmountOrders")
	sellOrders = db.Collection("SellOrders")
	logs = db.Collection("Logs")
}

func CreateAccount(user string, initialBalance float64) {
	newAccount := Account{
		User:    user,
		Balance: initialBalance,
	}
	res, err := accounts.InsertOne(context.TODO(), newAccount)

	if err != nil {
		fmt.Println("Failed to create account")
	} else {
		fmt.Println(res.InsertedID)
	}

}

func UpdateBalance(amount float64, user string) bool {
	filter := bson.M{"user": user}
	var result Account
	err := accounts.FindOne(context.TODO(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Println("ERROR: No account found for that user")
		return false
		/*
			result = Account{
				User: user,
			}
		*/
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

func CreateTransaction(transaction Transaction) {
	res, err := transactions.InsertOne(context.TODO(), transaction)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res.InsertedID)
}

func CreateBuyAmountOrder(buyAmountOrder BuyAmountOrder) {
	res, err := buyAmountOrders.InsertOne(context.TODO(), buyAmountOrder)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res.InsertedID)
}

func CreateSellAmountOrder(sellAmountOrder SellAmountOrder) {
	res, err := sellAmountOrders.InsertOne(context.TODO(), sellAmountOrder)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res.InsertedID)
}

func CreateTriggeredBuyAmountOrder(triggeredBuyAmountOrder TriggeredBuyAmountOrder) {
	res, err := triggeredBuyAmountOrders.InsertOne(context.TODO(), triggeredBuyAmountOrder)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Added TriggeredBuyAmountOrder to DB")
	fmt.Println(res.InsertedID)
}

func FindBuyAmountOrder(user string, stock string) (found bool, order BuyAmountOrder) {
	filter := bson.M{"user": user, "stock": stock}
	var buyAmountOrder BuyAmountOrder
	err := buyAmountOrders.FindOneAndDelete(context.TODO(), filter).Decode(&buyAmountOrder)
	if err == mongo.ErrNoDocuments {
		fmt.Println("No BuyAmountOrder for found for this user")
		return false, buyAmountOrder
	}
	return true, buyAmountOrder
}

func ConsumeLastTransaction(user string) Transaction {
	opts := options.FindOneAndDelete().SetSort(bson.M{"$natural": -1})
	filter := bson.M{"user": user}
	var transaction Transaction
	err := transactions.FindOneAndDelete(context.TODO(), filter, opts).Decode(&transaction)
	if err != nil {
		fmt.Println(err)
	}
	return transaction
}

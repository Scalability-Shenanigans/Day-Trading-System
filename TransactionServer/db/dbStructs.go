package db

type StockHolding struct {
	Stock  string `bson:"stock"`
	Amount int    `bson:"amount"`
}

type BuyOrder struct {
	Stock  string `bson:"stock"`
	Price  int    `bson:"price"`
	Amount int    `bson:"amount"`
}

type SellOrder struct {
	Stock  string `bson:"stock"`
	Price  int    `bson:"price"`
	Amount int    `bson:"amount"`
}

type Account struct {
	User        string         `bson:"user"`
	Balance     int            `bson:"balance"`
	Stocks      []StockHolding `bson:"stocks,omitempty"`
	Buy_Orders  []BuyOrder     `bson:"buy_orders,omitempty"`
	Sell_Orders []SellOrder    `bson:"sell_orders,omitempty"`
}

type Transaction struct {
	Transaction_ID int    `bson:"transaction_id"`
	Stock          string `bson:"stock"`
	Is_Buy         bool   `bson:"is_buy"`
	Amount         int    `bson:"amount"`
	User           string `bson:"user"`
}

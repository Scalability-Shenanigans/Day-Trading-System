db = new Mongo().getDB("DayTrading")

db.createCollection('Accounts')
db.createCollection('BuyOrders')
db.createCollection('BuyAmountOrders')
db.createCollection('SellOrders')
db.createCollection('PendingTransactions')
db.createCollection('FinishedTransactions')



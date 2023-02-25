db = new Mongo().getDB("DayTrading")

db.createCollection('Accounts')
db.createCollection('BuyOrders')
db.createCollection('SellOrders')
db.createCollection('PendingTransactions')



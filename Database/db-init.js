db = new Mongo().getDB("DayTrading")

db.createCollection('Accounts')
db.createCollection('PendingTransactions')



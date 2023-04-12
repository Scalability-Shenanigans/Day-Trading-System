# Day-Trading-System
This is a  solution for the contract issued by DayTrading Inc. i.e. to develop a prototype end-to-end solution to support their day trading clients and services. The project is developed as part of the SENG 468 Software Scalability course and is developed by a group of 5 people 
Objectives
The project will be a prototype for the system that is to be developed.The objective of the project can be outlined in the following points:
Support a large number of remote clients
Must be a  centralized transaction processing system.
Each client must log in through a web browser 
Perform several stock trading and account management activities 

## Functionality
The following stock trading and account management must be implemented in the prototype system:
* View their account
* Add money to their account
* Get a stock Quote
* Buy a number of shares in a stock
* Sell a number of shares in a stock they own
* Set an automated sell point for a stock
* Set a automated buy point for a stock
* Review their complete list of transactions
* Cancel a specified transaction prior to its being committed
* Commit a transaction

The application also has a containerized frontend developed using React.

## Running the application
To run the application please follow the following steps:
* Clone the repository
* Run ```docker-compose up --build``` in the root directory of the project to build and run the containers
* To access the frontend go to ```localhost:3000```
* To run the the userworkloads use the following commands after ```docker-compose up --build```
    * Go to ```LogfileParser``` directory and run ```go build -o /main``` to build a go build file for the parser 
    * Run ```./main``` to run the parser
    * Type ```automatic``` to run the automatic workload
        * This will prompt you to enter the number of users you want to run the workload for. Type it in the command line and press enter
    * Type ```manual``` to run the manual workload to run a single command manually from the workload file
* API endpoints can be accessed at ```localhost:5100``` and here are some of the endpoints available 
    * ```/add```: adds money to the account
	* ```/buy``` : set buy command
	* ```/commitBuy``` : commit buy command
	* ```/cancelBuy``` : cancel buy command
	* ```/cancelSetBuy``` : cancel set buy command
	* ```/setBuyAmount``` : set buy amount command
	* ```/setBuyTrigger``` : set buy trigger command
	* ```/sell``` : set sell command
	* ```/commitSell``` : commit sell command
	* ```/cancelSell``` : cancel sell command
	* ```/cancelSetSell``` : cancel set sell command
	* ```/setSellAmount```: set sell amount command
	* ```/setSellTrigger```: set sell trigger command
	* ```/dumplog```: dump log command
	* ```/quote```:  get quote command
	* ```/display``` : display command
    * all of these endpoints are POST endpoints and can be accessed using a POST request and the body of the request must contain appropriate parameters such as ```username```, ```stockSymbol```, ```amount```, ```trigger```, ```transactionNum``` etc according to the project requirements.

 ## Endpoints
These are the endpoints used in the application:
* ```localhost:5100``` : The API server load balanced by Nginx
* ```localhost:3000``` : The frontend server
* ```localhost:27017``` : The MongoDB database server
* ```localhost:6379``` : The Redis server
* ```localhost:4444``` : Local Quotes server

## Collections in MongoDB
The following collections are used in the application:
* ```Accounts``` : To maintain the accounts of the users
* ```BuyOrders``` : To maintain Buy orders
* ```BuyAmountOrders``` : To maintain Buy Amount orders
* ```TriggeredBuyAmountOrders``` : Used by Polling service to maintain triggered Buy/Sell Amount orders 
* ```SellOrders``` :  To maintain Sell orders
* ```PendingTransactions``` : To maintain pending transactions in the server
* ```FinishedTransactions``` : To maintain finished transactions in the server
* ```Logs``` : To maintain logs of the transactions used by the dumplog command



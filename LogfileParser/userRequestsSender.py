import requests
import re
import datetime

TRANSACTION_SERVER_URL = "http://localhost:5100/"


def send_request(endpoint, data=None):
    # Send the request and get the response
    response = requests.post(TRANSACTION_SERVER_URL+endpoint, json=data)

    # Check if the request was successful
    if response.status_code == requests.codes.ok:
        return "Success"
    else:
        # Raise an exception if the request was not successful
        response.raise_for_status()


def line_processor(line):
    transactionNum = int(re.search('\[(\d*)\]', line).group(1))

    print(line)
    line = line.split()[1:]
    args = line[0].split(",")
    print(args)

    command = args[0]

    if command == 'BUY':
        data = {
            "user": args[1],
            "stock": args[2],
            "amount": float(args[3]),
            "transactionNum": transactionNum
        }
        print(data)
        print(send_request("buy", data))
    elif command == 'COMMIT_BUY':
        data = {
            "user": args[1],
            "transactionNum": transactionNum
        }
        print(data)
        print(send_request("commitBuy", data))
    elif command == 'CANCEL_BUY':
        data = {
            "user": args[1],
            "transactionNum": transactionNum
        }
        print(data)
        print(send_request("cancelBuy", data))
    elif command == 'CANCEL_SET_BUY':
        data = {
            "user": args[1],
            "stock": args[2],
            "transactionNum": transactionNum
        }
        print(data)
        print(send_request("cancelSetBuy", data))
    elif command == 'SET_BUY_AMOUNT':
        data = {
            "user": args[1],
            "stock": args[2],
            "amount": float(args[3]),
            "transactionNum": transactionNum
        }
        print(data)
        print(send_request("setBuyAmount", data))
    elif command == 'SET_BUY_TRIGGER':
        data = {
            "user": args[1],
            "stock": args[2],
            "amount": float(args[3]),
            "transactionNum": transactionNum
        }
        print(data)
        print(send_request("setBuyTrigger", data))
    elif command == 'ADD':
        data = {
            "user": args[1],
            "amount": float(args[2]),
            "transactionNum": transactionNum
        }
        print(data)
        print(send_request("add", data))
    elif command == 'DUMPLOG':
        data = {
            "filename": args[1],
            "transactionNum": transactionNum
        }
        print(data)
        print("ended at: ", datetime.datetime.now())
        print(send_request("dumplog", data))
    elif command == 'SELL':
        data = {
            "user": args[1],
            "stock": args[2],
            "amount": float(args[3]),
            "transactionNum": transactionNum
        }
        print(data)
        print(send_request("sell", data))
    elif command == 'COMMIT_SELL':
        data = {
            "user": args[1],
            "transactionNum": transactionNum
        }
        print(data)
        print(send_request("commitSell", data))
    elif command == 'CANCEL_SET_SELL':
        data = {
            "user": args[1],
            "stock": args[2],
            "transactionNum": transactionNum
        }
        print(data)
        print(send_request("cancelSetSell", data))
    elif command == 'CANCEL_SELL':
        data = {
            "user": args[1],
            "transactionNum": transactionNum
        }
        print(send_request("cancelSell", data))
    elif command == 'SET_SELL_AMOUNT':
        data = {
            "user": args[1],
            "stock": args[2],
            "amount": float(args[3]),
            "transactionNum": transactionNum
        }
        print(data)
        print(send_request("setSellAmount", data))
    elif command == 'SET_SELL_TRIGGER':
        data = {
            "user": args[1],
            "stock": args[2],
            "amount": float(args[3]),
            "transactionNum": transactionNum
        }
        print(data)
        print(send_request("setSellTrigger", data))
    elif command == 'QUOTE':
        data = {
            "user": args[1],
            "stock": args[2],
            "transactionNum": transactionNum
        }
        print(data)
        print(send_request("quote", data))
    elif command == 'DISPLAY_SUMMARY':
        data = {
            "user": args[1],
            "transactionNum": transactionNum
        }
        print(data)
        print(send_request("display", data))


def main():
    choice = input("mode: manual, dumplog, dbwipe, or automatic? ")
    if choice == "manual":
        while True:
            line = input("enter command: ")
            # transactionNum = int(re.search('\[(\d*)\]', line).group(1))
            if line == "quit":
                break
            else:
                line_processor(line)
    elif choice == "dumplog":
        data = {
            "filename": "logfile.xml",
            "transactionNum": 100
        }
        print(data)
        print(send_request("dumplog", data))
    elif choice == "dbwipe":  # command I added for wiping all the collections
        print(send_request("dbwipe"))
    elif choice == "automatic":
        startTime = datetime.datetime.now()
        print("Starting at: ", startTime)
        with open("final.txt", "r") as f:
            for line in f:
                line_processor(line)
        endTime = datetime.datetime.now()
        difference = endTime - startTime

        hours = difference.seconds // 3600
        minutes = (difference.seconds % 3600) // 60
        seconds = difference.seconds % 60

        print(
            f"The difference is {hours} hours, {minutes} minutes, and {seconds} seconds.")


if __name__ == "__main__":
    main()

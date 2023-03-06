import requests
import re

TRANSACTION_SERVER_URL = "http://localhost:8080/"


def send_request(endpoint, data=None):
    # Send the request and get the response
    response = requests.post(TRANSACTION_SERVER_URL+endpoint, json=data)

    # Check if the request was successful
    if response.status_code == requests.codes.ok:
        return "Success"
    else:
        # Raise an exception if the request was not successful
        response.raise_for_status()

def main():
    while True:
        line = input("enter command: ")
        transactionNum = int(re.search('\[(\d*)\]', line).group(1))
        if line == "quit":
            break
        else:
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
                }
                print(data)
                print(send_request("commitBuy", data))
            elif command == 'CANCEL_BUY':
                data = {
                    "user": args[1],
                }
                print(data)
                print(send_request("cancelBuy", data))
            elif command == 'CANCEL_SET_BUY':
                data = {
                    "user": args[1],
                    "stock": args[2]
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
                    "user": args[1]
                }
                print(data)
                print(send_request("commitSell", data))
            elif command == 'CANCEL_SET_SELL':
                data = {
                    "user": args[1],
                    "stock": args[2]
                }
                print(data)
                print(send_request("cancelSetSell", data))
            elif command == 'CANCEL_SELL':
                print(data)
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
                    "stock": args[2]
                }
                print(data)
                print(send_request("quote", data))
            elif command == 'DISPLAY_SUMMARY':
                data = {
                    "user": args[1]
                }
                print(data)
                print(send_request("display", data))





if __name__ == "__main__":
    main()

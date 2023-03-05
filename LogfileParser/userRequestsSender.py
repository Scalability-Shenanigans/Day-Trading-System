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
                    "amount": float(args[3])
                }
                print(data)
            elif command == 'ADD':
                data = {
                    "user": args[1],
                    "amount": float(args[2]),
                    "transactionNum": transactionNum
                }
                print(data)
                print(send_request("add", data))



if __name__ == "__main__":
    main()

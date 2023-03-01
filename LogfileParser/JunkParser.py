import requests

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
    with open('user1.txt', 'r') as f:
        for line in f:
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
                print(send_request("buy", data))
            elif command == 'ADD':
                data = {
                    "user": args[1],
                    "amount": float(args[2])
                }
                print(send_request("add", data))
            elif command == 'QUOTE':
                pass
            elif command == 'COMMIT_BUY':
                pass
            elif command == 'SET_BUY_AMOUNT':
                data = {
                    "user": args[1],
                    "stock": args[2],
                    "amount": float(args[3])
                }
                print(send_request("setBuyAmount", data))
            elif command == 'SET_BUY_TRIGGER':
                data = {
                    "user": args[1],
                    "stock": args[2],
                    "price": float(args[3])
                }
                print(send_request("setBuyAmount", data))


if __name__ == "__main__":
    main()

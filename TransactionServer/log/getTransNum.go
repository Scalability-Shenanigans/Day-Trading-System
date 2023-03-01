package main

import (
	"fmt"
	"os"
)

// no Redis setup currently, so store transaction number in a file and everytime the
// method is called, return the transaction number, increment it, and write it back
// to the file
func getTransactionNumber() (int64, error) {
	file, err := os.OpenFile("transNumber.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var transactionNumber int64
	if _, err := fmt.Fscanf(file, "%d", &transactionNumber); err != nil {
		return 0, err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return 0, err
	}
	_, err = fmt.Fprintf(file, "%d", transactionNumber+1)
	if err != nil {
		return 0, err
	}

	return transactionNumber, nil
}

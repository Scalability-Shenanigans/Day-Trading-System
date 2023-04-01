package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"net"
	"time"
)

func generateQuote() (float64, int64, string) {
	price := randFloat64(100)
	timestamp := time.Now().Unix()
	crypto := generateCryptoKey(44)
	return price, timestamp, crypto
}

func randFloat64(max float64) float64 {
	r, _ := rand.Int(rand.Reader, big.NewInt(int64(max*1e6)))
	return float64(r.Int64()) / 1e6
}

func generateCryptoKey(length int) string {
	bytes := make([]byte, (length+3)/4*3)
	_, _ = rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}

func sendQuote(conn net.Conn) {
	price, timestamp, crypto := generateQuote()
	response := fmt.Sprintf("%f,%d,%s", price, timestamp, crypto)
	defer conn.Close()

	_, err := conn.Write([]byte(response))
	if err != nil {
		log.Printf("Error sending quote: %v", err)
		return
	}

	log.Printf("Sent quote: %s", response)
}

func main() {
	listener, err := net.Listen("tcp", ":4444")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	defer listener.Close()
	log.Printf("Listening on %s", listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go sendQuote(conn)
	}
}

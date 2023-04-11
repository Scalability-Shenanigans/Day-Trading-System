package middleware

import (
	"TransactionServer/cache"
	"context"
	"fmt"
	"net/http"
)

type contextKey string

const TransactionNumberKey contextKey = "transactionNumber"

var redisClient = cache.NewRedisClient()

func getNextTransactionNumber(r *cache.RedisClient) (uint64, error) {
	ctx := context.Background()
	count, err := r.Client.Incr(ctx, "transactionCounter").Result()
	if err != nil {
		return 0, err
	}
	return uint64(count), nil
}

func TransactionNumberMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the next transaction number
		transactionNumber, err := getNextTransactionNumber(redisClient)
		if err != nil {
			http.Error(w, "Failed to get transaction number", http.StatusInternalServerError)
			return
		}

		fmt.Printf("New transaction number: %d\n", transactionNumber) // Add logging

		ctx := context.WithValue(r.Context(), TransactionNumberKey, transactionNumber)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetTransactionNumberFromContext(r *http.Request) uint64 {
	return r.Context().Value(TransactionNumberKey).(uint64)
}

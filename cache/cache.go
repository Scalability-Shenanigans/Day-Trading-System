package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisAddress = "localhost:6379"

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient() *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr: RedisAddress,
	})

	return &RedisClient{
		Client: client,
	}
}

func (r *RedisClient) Set(key string, value string, ttl time.Duration) error {
	ctx := context.Background()

	return r.Client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisClient) Get(key string) (bool, float64) {
	ctx := context.Background()
	value, err := r.Client.Get(ctx, key).Result()
	if err == nil && value != "" {
		val, _ := strconv.ParseFloat(value, 64)
		return true, val
	}
	return false, 0
}

package initializer

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// ReddisClient initializes and returns a Redis client.
var ReddisClient = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("RedisAddr"),
	Password: os.Getenv("RedisPass"),
	DB:       0,
})

func SetRedis(key string, value any, expirationTime time.Duration) error {
	if err := ReddisClient.Set(context.Background(), key, value, expirationTime).Err(); err != nil {
		return err
	}
	return nil
}

func GetRedis(key string) (string, error) {
	jsonData, err := ReddisClient.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}
	return jsonData, nil
}

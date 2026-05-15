package db

import (
	"context"
	"log/slog"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	address := os.Getenv("REDIS_ADDR")
	
	RedisClient = redis.NewClient(&redis.Options{
		Addr: address,
		DB: 0,
	})

	// Test koneksi
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		slog.Error("Failed to connect to Redis", "error", err)

		// Tidak panik - app tetap jalan tanpa cache
		RedisClient = nil
		return
	}

	slog.Info("Redis connected succesfully")
}
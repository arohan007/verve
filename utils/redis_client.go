package utils

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

// RedisClient is a global Redis client
var RedisClient *redis.Client

// InitializeRedis initializes the Redis client
func InitializeRedis(redisAddr string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr, // Redis server address
		Password: "",        // No password set
		DB:       0,         // Use default DB
	})

	// Test the connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis")
	RedisClient = client
	return RedisClient
}

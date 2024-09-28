package reddis

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
)

func NewReddisClient(ctx context.Context) (*redis.Client, error) {
	// Connect to Redis
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Replace with your Redis server address
		Password: "",               // No password for local development
		DB:       0,                // Default DB
	})

	// Ping the Redis server to check the connection
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
		return nil, err
	}
	fmt.Println("Connected to Redis:", pong)
	return client, nil
}

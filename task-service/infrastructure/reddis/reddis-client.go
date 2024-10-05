package reddis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"task-service/internal/model"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
)

func NewReddisClient(ctx context.Context) (*redis.Client, error) {
	// Connect to Redis
	host := fmt.Sprintf("localhost:%s", os.Getenv("REDIS_PORT"))
	client := redis.NewClient(&redis.Options{
		Addr:     host, // Replace with your Redis server address
		Password: "",   // No password for local development
		DB:       0,    // Default DB
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

func GetLoginInfoFromRedis(ctx context.Context, jwtToken string) (loginInfo model.LoginCacheData, err error) {
	loginJson, err := RedisClient.Get(ctx, jwtToken).Result()
	if err != nil {
		return loginInfo, fmt.Errorf("failed get login info from redis")
	}
	err = json.Unmarshal([]byte(loginJson), &loginInfo)
	if err != nil {
		return loginInfo, fmt.Errorf("failed convert data login info from json")
	}
	return loginInfo, nil
}

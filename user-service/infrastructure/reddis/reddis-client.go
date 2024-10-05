package reddis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"user-service/internal/model"

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

func GetUserInfoFromRedis(ctx context.Context, userId string) (user model.User, err error) {
	userJson, err := RedisClient.Get(ctx, userId).Result()
	if err != nil {
		return user, fmt.Errorf("failed get user info from redis")
	}
	err = json.Unmarshal([]byte(userJson), &user)
	if err != nil {
		return user, fmt.Errorf("failed convert data user info from json")
	}
	return user, nil
}

func SetUserInfoToRedis(ctx context.Context, user model.User) error {
	userJson, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed convert user info to json")
	}

	// send user info to reddis data
	err = RedisClient.Set(ctx, user.Id, userJson, model.UserSessionTime).Err() // Set waktu kadaluarsa 30 menit
	if err != nil {
		return fmt.Errorf("error saving login info to redis. err = %s", err.Error())
	}
	return nil
}

func SetLoginInfoToRedis(ctx context.Context, tokenKey string, loginInfo model.LoginCacheData) error {
	loginJson, err := json.Marshal(loginInfo)
	if err != nil {
		return fmt.Errorf("failed convert login info to json")
	}

	// send login info to reddis data
	err = RedisClient.Set(ctx, tokenKey, loginJson, model.UserSessionTime).Err() // Set waktu kadaluarsa 30 menit
	if err != nil {
		return fmt.Errorf("error saving login info to redis. err = %s", err.Error())
	}
	return nil
}

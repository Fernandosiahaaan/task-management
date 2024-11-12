package reddis

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
	"user-service/internal/model"

	redistrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/redis/go-redis.v9"

	"github.com/redis/go-redis/v9"
)

const (
	PrefixKeyLoginInfo = "user-service:jwt"
	PrefixKeyUserInfo  = "user-service:user"
)

type RedisCln struct {
	Redis  *redis.Client
	Ctx    context.Context
	Cancel context.CancelFunc
}

func NewReddisClient(ctx context.Context) (*RedisCln, error) {
	// Connect to Redis
	ctxRedis, cancelRedis := context.WithCancel(ctx)
	host := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	var opts *redis.Options = &redis.Options{
		Addr:        host, // Replace with your Redis server address
		Password:    "",   // No password for local development
		DB:          0,    // Default DB
		DialTimeout: 10 * time.Second,
	}
	client := redis.NewClient(opts)

	// Ping the Redis server to check the connection
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	c := redistrace.NewClient(opts)
	c.Set(ctx, "test_key", "test_value", 0)

	var redis *RedisCln = &RedisCln{
		Redis:  client,
		Ctx:    ctxRedis,
		Cancel: cancelRedis,
	}
	fmt.Println("Connected to Redis:", pong)
	return redis, nil
}

func (r *RedisCln) GetLoginInfo(jwtToken string) (loginInfo model.LoginCacheData, err error) {
	keyLoginInfo := fmt.Sprintf("%s:%s", PrefixKeyLoginInfo, jwtToken)
	loginJson, err := r.Redis.Get(r.Ctx, keyLoginInfo).Result()
	if err != nil {
		return loginInfo, fmt.Errorf("failed get login info from redis")
	}
	err = json.Unmarshal([]byte(loginJson), &loginInfo)
	if err != nil {
		return loginInfo, fmt.Errorf("failed convert data login info from json")
	}
	return loginInfo, nil
}

func (r *RedisCln) SetLoginInfo(ctx context.Context, jwtToken string, loginInfo model.LoginCacheData) error {
	loginJson, err := json.Marshal(loginInfo)
	if err != nil {
		return fmt.Errorf("failed convert login info to json")
	}

	// send login info to reddis data
	keyLoginInfo := fmt.Sprintf("%s:%s", PrefixKeyLoginInfo, jwtToken)
	err = r.Redis.Set(ctx, keyLoginInfo, loginJson, model.UserSessionTime).Err() // Set waktu kadaluarsa 30 menit
	if err != nil {
		return fmt.Errorf("error saving login info to redis. err = %s", err.Error())
	}
	return nil
}

func (r *RedisCln) DeleteLoginInfo(jwtToken string) error {
	keyLoginInfo := fmt.Sprintf("%s:%s", PrefixKeyLoginInfo, jwtToken)
	return r.Redis.Del(r.Ctx, keyLoginInfo).Err()
}

func (r *RedisCln) GetUserInfo(userId string) (user *model.User, err error) {
	userInfo := fmt.Sprintf("%s:%s", PrefixKeyUserInfo, userId)
	userJson, err := r.Redis.Get(r.Ctx, userInfo).Result()
	if err != nil {
		return user, fmt.Errorf("failed get user info from redis")
	}
	err = json.Unmarshal([]byte(userJson), &user)
	if err != nil {
		return user, fmt.Errorf("failed convert data user info from json")
	}
	return user, nil
}

func (r *RedisCln) SaveUserInfo(user model.User) error {
	userJson, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed convert user info to json")
	}

	// send user info to reddis data
	keyUserInfo := fmt.Sprintf("%s:%s", PrefixKeyUserInfo, user.Id)
	err = r.Redis.Set(r.Ctx, keyUserInfo, userJson, model.UserSessionTime).Err() // Set waktu kadaluarsa 30 menit
	if err != nil {
		return fmt.Errorf("error saving login info to redis. err = %s", err.Error())
	}
	return nil
}

func (r *RedisCln) DeleteUserInfo(userId string) error {
	keyLoginInfo := fmt.Sprintf("%s:%s", PrefixKeyUserInfo, userId)
	return r.Redis.Del(r.Ctx, keyLoginInfo).Err()
}

func (r *RedisCln) Close() {
	r.Redis.Close()
	r.Cancel()
}

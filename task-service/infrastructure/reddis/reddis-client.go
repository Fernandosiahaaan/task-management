package reddis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"task-service/internal/model"

	"github.com/redis/go-redis/v9"
	redistrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/redis/go-redis.v9"
)

const (
	PrefixKeyLoginInfo = "user-service:jwt"
	PrefixKeyUserInfo  = "user-service:user"
	PrefixKeyTaskInfo  = "task-service:task"
)

type Redis struct {
	Redis  *redis.Client
	Ctx    context.Context
	Cancel context.CancelFunc
}

func NewReddisClient(ctx context.Context) (*Redis, error) {
	// Connect to Redis
	ctxRedis, cancelRedis := context.WithCancel(ctx)
	host := fmt.Sprintf("localhost:%s", os.Getenv("REDIS_PORT"))
	var opts *redis.Options = &redis.Options{
		Addr:     host, // Replace with your Redis server address
		Password: "",   // No password for local development
		DB:       0,    // Default DB
	}
	client := redis.NewClient(opts)

	// Ping the Redis server to check the connection
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
		return nil, err
	}

	c := redistrace.NewClient(opts)
	c.Set(ctx, "test_key", "test_value", 0)

	var redis *Redis = &Redis{
		Redis:  client,
		Ctx:    ctxRedis,
		Cancel: cancelRedis,
	}
	fmt.Println("Connected to Redis:", pong)
	return redis, nil
}

func (r *Redis) GetLoginInfoFromRedis(jwtToken string) (loginInfo model.LoginCacheData, err error) {
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

func (r *Redis) SetTaskInfoToRedis(task *model.Task) error {
	taskJson, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed convert user info to json")
	}

	// send user info to reddis data
	keyTaskInfo := fmt.Sprintf("%s:%d", PrefixKeyTaskInfo, task.Id)
	err = r.Redis.Set(r.Ctx, keyTaskInfo, taskJson, model.SessionTime).Err() // Set waktu kadaluarsa 30 menit
	if err != nil {
		return fmt.Errorf("error saving login info to redis. err = %s", err.Error())
	}
	return nil
}

func (r *Redis) GetTaskInfoFromRedis(taskId int64) (taskInfo *model.Task, err error) {
	keyTaskInfo := fmt.Sprintf("%s:%d", PrefixKeyTaskInfo, taskId)
	taskJson, err := r.Redis.Get(r.Ctx, keyTaskInfo).Result()
	if err != nil {
		return nil, fmt.Errorf("failed get login info from redis")
	}
	err = json.Unmarshal([]byte(taskJson), &taskInfo)
	if err != nil {
		return nil, fmt.Errorf("failed convert data login info from json")
	}
	return taskInfo, nil
}

func (r *Redis) Close() {
	r.Redis.Close()
	r.Cancel()
}

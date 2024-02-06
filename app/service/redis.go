package service

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

const (
	redisHost     = "localhost"
	redisPort     = 6379
	redisPassword = ""
)

func InitRedis() (*redis.Client, error) {
	options := redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisHost, redisPort),
		Password: redisPassword,
	}
	client := redis.NewClient(&options)
	if err := client.Ping(context.Background()).Err(); err != nil {
		// TODO
	}
	return client, nil
}

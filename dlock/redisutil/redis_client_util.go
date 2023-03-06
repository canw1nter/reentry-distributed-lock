package redisutil

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var (
	ctx context.Context
	rdb *redis.Client
)

func init() {
	rdb, ctx = redis.NewClient(&redis.Options{
		Addr:     "192.168.182.100:6379",
		Password: "redis",
		DB:       0,
	}), context.Background()
}

func GetRedisClient() (*redis.Client, context.Context) {
	return rdb, ctx
}

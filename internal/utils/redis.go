package utils

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func MustGetRedis(connectionString string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: connectionString,
	})
	ctx, cancelF := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelF()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	return client
}

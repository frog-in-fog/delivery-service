package cmd

import (
	"context"
	"fmt"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisConnection(config *config.Config) error {
	ctx := context.Background()

	RedisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisURL,
		Password: "",
		DB:       0,
	})

	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		return err
	}

	err := RedisClient.Set(ctx, "test", "How to Refresh Access Tokens the Right Way in Golang", 0).Err()
	if err != nil {
		return err
	}

	fmt.Println("✅ Redis client connected successfully...")
	return nil
}

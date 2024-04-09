package service

import (
	"context"

	"github.com/frog-in-fog/delivery-system/auth-service/internal/storage/redis"
)

func (s *authService) LogoutUser(userId string) error {
	errDelAccessToken := redis.RedisClient.Del(context.Background(), userId+accessSuffix).Err()
	if errDelAccessToken != nil {
		return errDelAccessToken
	}
	errDelRefreshToken := redis.RedisClient.Del(context.Background(), userId+refreshSuffix).Err()
	if errDelRefreshToken != nil {
		return errDelRefreshToken
	}

	return nil
}

package service

import (
	"context"
	"errors"

	"github.com/frog-in-fog/delivery-system/auth-service/internal/config"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/storage/redis"
	"github.com/frog-in-fog/delivery-system/auth-service/pkg/tokens"
)

func (s *authService) TokenPair(accessToken string, cfg *config.Config) (string, error) {
	tokenDetails, err := tokens.ValidateToken(accessToken, cfg.AccessTokenPublicKey)
	if err != nil {
		if errors.Is(err, tokens.ErrInvalidToken) {
			refreshToken, err := redis.RedisClient.Get(context.Background(), tokenDetails.UserID+refreshSuffix).Result()
			if err != nil {
				return "", err
			}
			tokenPair, err := s.RefreshAccessToken(refreshToken, cfg)
			if err != nil {
				return "", err
			}
			refreshedAccessToken := tokenPair["access_token"]
			return refreshedAccessToken, nil
		}
		return "", err
	}

	//userId := tokenDetails.UserID
	//
	// check if user exists in db
	// _, err = s.userStorage.GetUserById(context.Background(), userId)
	// if err != nil {
	// 	if errors.Is(err, sql.ErrNoRows) {
	// 		return ErrUserNotFound
	// 	}
	// 	return err
	// }

	return "allowed", nil
}

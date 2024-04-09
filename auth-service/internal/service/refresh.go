package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/frog-in-fog/delivery-system/auth-service/internal/config"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/storage/redis"
	"github.com/frog-in-fog/delivery-system/auth-service/pkg/tokens"
)

func (s *authService) RefreshAccessToken(refreshToken string, cfg *config.Config) (map[string]string, error) {
	tokenClaims, err := tokens.ValidateToken(refreshToken, cfg.RefreshTokenPublicKey)
	if err != nil {
		return nil, err
	}

	user, err := s.userStorage.GetUserById(context.Background(), tokenClaims.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	accessTokenDetails, err := tokens.CreateToken(user.ID, cfg.AccessTokenExpiresIn, cfg.AccessTokenPrivateKey)
	if err != nil {
		return nil, err
	}

	errAccess := redis.RedisClient.Set(context.TODO(), user.ID+accessSuffix, *accessTokenDetails.Token, time.Unix(*accessTokenDetails.ExpiresIn, 0).Sub(time.Now())).Err()
	if errAccess != nil {
		return nil, errAccess
	}

	refreshTokenDetails, err := tokens.CreateToken(user.ID, cfg.RefreshTokenExpiresIn, cfg.RefreshTokenPrivateKey)
	if err != nil {
		return nil, err
	}

	errRefresh := redis.RedisClient.Set(context.TODO(), user.ID+refreshSuffix, *refreshTokenDetails.Token, time.Unix(*refreshTokenDetails.ExpiresIn, 0).Sub(time.Now())).Err()
	if err != nil {
		return nil, errRefresh
	}

	tokenPair := make(map[string]string)
	tokenPair["access_token"] = *accessTokenDetails.Token
	tokenPair["refresh_token"] = *refreshTokenDetails.Token

	return tokenPair, nil
}

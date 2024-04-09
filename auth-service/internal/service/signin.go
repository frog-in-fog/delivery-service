package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/frog-in-fog/delivery-system/auth-service/internal/config"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/models"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/storage/redis"
	"github.com/frog-in-fog/delivery-system/auth-service/pkg/tokens"
	"golang.org/x/crypto/bcrypt"
)

func (s *authService) SignInUser(newUser *models.User, cfg *config.Config) (map[string]string, error) {
	user, err := s.userStorage.GetUserByEmail(context.Background(), newUser.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(newUser.PasswordHash)); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessTokenDetails, err := tokens.CreateToken(user.ID, cfg.AccessTokenExpiresIn, cfg.AccessTokenPrivateKey)
	if err != nil {
		return nil, err
	}

	refreshTokenDetails, err := tokens.CreateToken(user.ID, cfg.RefreshTokenExpiresIn, cfg.RefreshTokenPrivateKey)
	if err != nil {
		return nil, err
	}
	errAccess := redis.RedisClient.Set(context.TODO(), user.ID+accessSuffix, *accessTokenDetails.Token, time.Unix(*accessTokenDetails.ExpiresIn, 0).Sub(time.Now())).Err()
	if errAccess != nil {
		return nil, errAccess
	}

	errRefresh := redis.RedisClient.Set(context.TODO(), user.ID+refreshSuffix, *refreshTokenDetails.Token, time.Unix(*refreshTokenDetails.ExpiresIn, 0).Sub(time.Now())).Err()
	if errRefresh != nil {
		return nil, errRefresh
	}

	tokenPair := make(map[string]string)
	tokenPair["access_token"] = *accessTokenDetails.Token
	tokenPair["refresh_token"] = *refreshTokenDetails.Token

	return tokenPair, nil

}

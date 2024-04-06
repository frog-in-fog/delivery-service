package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/config"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/models"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/storage"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/storage/redis"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/storage/sqlite"
	"github.com/frog-in-fog/delivery-system/auth-service/pkg/tokens"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type AuthUsecase interface {
	SignUpUser(user *models.User) error
	SignInUser(user *models.User, cfg *config.Config) (map[string]string, error)
	RefreshAccessToken() error
	LogoutUser() error
}

type authService struct {
	userStorage storage.UserStorage
}

func NewAuthService(userStorage storage.UserStorage) AuthUsecase {
	return &authService{userStorage: userStorage}
}

func (s authService) SignUpUser(user *models.User) error {
	if err := s.userStorage.CreateUser(context.Background(), user); err != nil {
		if errors.Is(err, sqlite.ErrUserAlreadyExists) {
			return sqlite.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (s authService) SignInUser(newUser *models.User, cfg *config.Config) (map[string]string, error) {
	user, err := s.userStorage.GetUserByEmail(context.Background(), newUser.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sqlite.ErrUserNotFound
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

	errAccess := redis.RedisClient.Set(context.TODO(), accessTokenDetails.TokenUuid, user.ID, time.Unix(*accessTokenDetails.ExpiresIn, 0).Sub(time.Now())).Err()
	if errAccess != nil {
		return nil, errAccess
	}

	errRefresh := redis.RedisClient.Set(context.TODO(), refreshTokenDetails.TokenUuid, user.ID, time.Unix(*refreshTokenDetails.ExpiresIn, 0).Sub(time.Now())).Err()
	if errRefresh != nil {
		return nil, errRefresh
	}

	tokenPair := make(map[string]string)
	tokenPair["access_token"] = *accessTokenDetails.Token
	tokenPair["refresh_token"] = *refreshTokenDetails.Token

	return tokenPair, nil

}

func (s authService) RefreshAccessToken() error {
	return errors.New("refresh access token works")
}

func (s authService) LogoutUser() error {
	return errors.New("logout works")
}

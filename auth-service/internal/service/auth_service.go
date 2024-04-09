package service

import (
	"errors"

	"github.com/frog-in-fog/delivery-system/auth-service/internal/config"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/models"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/storage"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("Invalid token")
)

const (
	accessSuffix  = ":access"
	refreshSuffix = ":refresh"
)

type AuthUsecase interface {
	SignUpUser(user *models.User) error
	SignInUser(user *models.User, cfg *config.Config) (map[string]string, error)
	RefreshAccessToken(refreshToken string, cfg *config.Config) (map[string]string, error)
	LogoutUser(userId string) error
	TokenPair(accessToken string, cfg *config.Config) (string, error)
}

type authService struct {
	userStorage storage.UserStorage
}

func NewAuthService(userStorage storage.UserStorage) AuthUsecase {
	return &authService{userStorage: userStorage}
}

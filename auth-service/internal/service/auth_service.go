package service

import (
	"context"
	"errors"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/models"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/storage"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/storage/sqlite"
)

type AuthUsecase interface {
	SignUpUser(user *models.User) error
	SignInUser(user *models.User) error
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

func (s authService) SignInUser(user *models.User) error {
	return errors.New("sign in works")
}

func (s authService) RefreshAccessToken() error {
	return errors.New("refresh access token works")
}

func (s authService) LogoutUser() error {
	return errors.New("logout works")
}

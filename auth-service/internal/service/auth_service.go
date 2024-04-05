package service

import (
	"errors"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/models"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/storage"
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
	return errors.New("sign up works")
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

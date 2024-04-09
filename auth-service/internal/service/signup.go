package service

import (
	"context"
	"errors"

	"github.com/frog-in-fog/delivery-system/auth-service/internal/models"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/storage/sqlite"
)

func (s *authService) SignUpUser(user *models.User) error {
	if err := s.userStorage.CreateUser(context.Background(), user); err != nil {
		if errors.Is(err, sqlite.ErrUserAlreadyExists) {
			return sqlite.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

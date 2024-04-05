package storage

import (
	"context"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/models"
)

type UserStorage interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserById(ctx context.Context, id string) (*models.User, error)
}

package cmd

import (
	"github.com/frog-in-fog/delivery-system/auth-service/internal/config"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/handlers"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/service"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/storage"
)

type Handlers struct {
	AuthHandler handlers.AuthHandler
}

func InitHttpHandlers(cfg *config.Config, userStorage storage.UserStorage) *Handlers {
	// service
	authService := service.NewAuthService(userStorage)

	// handlers
	authHandlers := handlers.NewAuthHandler(authService)

	return &Handlers{AuthHandler: authHandlers}

}

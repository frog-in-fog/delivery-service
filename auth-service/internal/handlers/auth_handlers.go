package handlers

import (
	"net/http"

	"github.com/frog-in-fog/delivery-system/auth-service/internal/config"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/service"
)

type AuthHandler interface {
	SignUpUser(w http.ResponseWriter, r *http.Request)
	SignInUser(w http.ResponseWriter, r *http.Request)
	LogoutUser(w http.ResponseWriter, r *http.Request)
	TokenPair(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	authService service.AuthUsecase
	cfg         *config.Config
}

func NewAuthHandler(authService service.AuthUsecase, cfg *config.Config) AuthHandler {
	return &authHandler{authService: authService, cfg: cfg}
}

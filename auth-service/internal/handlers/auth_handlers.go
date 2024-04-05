package handlers

import (
	"github.com/frog-in-fog/delivery-system/auth-service/internal/service"
	"net/http"
)

type AuthHandler interface {
	SignUpUser(w http.ResponseWriter, r *http.Request)
	SignInUser(w http.ResponseWriter, r *http.Request)
	LogoutUser(w http.ResponseWriter, r *http.Request)
	RefreshAccessToken(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	authService service.AuthUsecase
}

func NewAuthHandler(authService service.AuthUsecase) AuthHandler {
	return &authHandler{authService: authService}
}

func (h *authHandler) SignUpUser(w http.ResponseWriter, r *http.Request) {

}

func (h *authHandler) SignInUser(w http.ResponseWriter, r *http.Request) {

}

func (h *authHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {

}

func (h *authHandler) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {

}

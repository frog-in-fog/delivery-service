package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/frog-in-fog/delivery-system/auth-service/internal/config"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/models"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/models/dto"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/service"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/storage/sqlite"
	"github.com/frog-in-fog/delivery-system/auth-service/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

func (h *authHandler) SignUpUser(w http.ResponseWriter, r *http.Request) {
	var payload dto.SignUpInput

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.RenderJSON(w, err)
		return
	}
	if err = json.Unmarshal(body, &payload); err != nil {
		utils.RenderJSON(w, err)
		return
	}

	validationErrs := dto.ValidateStruct(payload)
	if validationErrs != nil {
		utils.RenderJSON(w, validationErrs)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.RenderJSON(w, err)
		return
	}

	newUser := models.User{
		ID:           uuid.New().String(),
		Email:        payload.Email,
		PasswordHash: string(passwordHash),
	}

	if err = h.authService.SignUpUser(&newUser); err != nil {
		if errors.Is(err, sqlite.ErrUserAlreadyExists) {
			utils.RenderJSON(w, "User with such email already exists!")
			return
		}
		utils.RenderJSON(w, err)
		return
	}

	utils.RenderJSON(w, "User signed up successfully!")
	return
}

func (h *authHandler) SignInUser(w http.ResponseWriter, r *http.Request) {
	var payload dto.SignInInput
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.RenderJSON(w, err)
		return
	}
	if err = json.Unmarshal(body, &payload); err != nil {
		utils.RenderJSON(w, err)
		return
	}

	validationErrs := dto.ValidateStruct(payload)
	if validationErrs != nil {
		utils.RenderJSON(w, validationErrs)
		return
	}

	newUser := models.User{
		Email:        payload.Email,
		PasswordHash: payload.Password,
	}

	_, err = h.authService.SignInUser(&newUser, h.cfg)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			utils.RenderJSON(w, service.ErrInvalidCredentials)
			return
		}
		utils.RenderJSON(w, err)
		return
	}

	loggedInCookie := http.Cookie{
		Name:     "logged_in",
		Value:    "true",
		Path:     "/",
		MaxAge:   h.cfg.AccessTokenMaxAge * 60,
		Secure:   false,
		HttpOnly: false,
		Domain:   "localhost",
	}

	http.SetCookie(w, &loggedInCookie)

	utils.RenderJSON(w, "Logged in")
	return
}

func (h *authHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	var payload dto.LogoutInput
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.RenderJSON(w, err)
		return
	}

	if err = json.Unmarshal(body, &payload); err != nil {
		utils.RenderJSON(w, err)
		return
	}

	userId := payload.UserId

	if err = h.authService.LogoutUser(userId); err != nil {
		utils.RenderJSON(w, err)
		return
	}

	expired := time.Now().Add(-time.Hour * 24)
	loggedOutCookie := http.Cookie{
		Name:    "logged_in",
		Value:   "false",
		Expires: expired,
	}

	http.SetCookie(w, &loggedOutCookie)
}

type TokenPairResp struct {
	Data string `json:"data"`
}

func (h *authHandler) TokenPair(w http.ResponseWriter, r *http.Request) {
	accessTokenHeader := r.Header.Get("Authorization")
	accessToken := strings.TrimPrefix(accessTokenHeader, "Bearer ")
	if len(accessToken) == 0 {
		resp := TokenPairResp{
			Data: "empty access token",
		}
		utils.RenderJSON(w, resp)
		return
	}

	res, err := h.authService.TokenPair(accessToken, h.cfg)
	if err != nil {
		resp := TokenPairResp{
			Data: err.Error(),
		}
		utils.RenderJSON(w, resp)
		return
	}

	if len(res) != 0 {
		if res == "allowed" {
			resp := TokenPairResp{
				Data: "allowed",
			}
			utils.RenderJSON(w, resp)
			return
		} else {
			resp := TokenPairResp{
				Data: res,
			}
			utils.RenderJSON(w, resp)
			return
		}
	}

}

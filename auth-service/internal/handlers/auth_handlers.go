package handlers

import (
	"encoding/json"
	"errors"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/config"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/models"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/models/dto"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/service"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/storage/sqlite"
	"github.com/frog-in-fog/delivery-system/auth-service/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
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
	cfg         *config.Config
}

func NewAuthHandler(authService service.AuthUsecase, cfg *config.Config) AuthHandler {
	return &authHandler{authService: authService, cfg: cfg}
}

func (h *authHandler) SignUpUser(w http.ResponseWriter, r *http.Request) {
	var payload *dto.SignUpInput

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
	var payload *dto.SignInInput
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

	tokenPair, err := h.authService.SignInUser(&newUser, h.cfg)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			utils.RenderJSON(w, service.ErrInvalidCredentials)
			return
		}
		utils.RenderJSON(w, err)
		return
	}

	accessToken := tokenPair["access_token"]
	refreshToken := tokenPair["refresh_token"]

	accessCookie := http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   h.cfg.AccessTokenMaxAge * 60,
		Secure:   false,
		HttpOnly: true,
		Domain:   "localhost",
	}

	refreshCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   h.cfg.RefreshTokenMaxAge * 60,
		Secure:   false,
		HttpOnly: true,
		Domain:   "localhost",
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

	http.SetCookie(w, &accessCookie)
	http.SetCookie(w, &refreshCookie)
	http.SetCookie(w, &loggedInCookie)

	utils.RenderJSON(w, "Logged in")
	return
}

func (h *authHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	//cookie, err := r.Cookie("refresh_token")
	//if err != nil {
	//	if errors.Is(err, http.ErrNoCookie) {
	//		log.Println("Cookie refresh_token not found")
	//		utils.RenderJSON(w, "could not refresh access token")
	//		return
	//	}
	//	utils.RenderJSON(w, "could not refresh access token")
	//}
	//
	//refreshToken := cookie.Value
}

func (h *authHandler) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			log.Println("Cookie refresh_token not found")
			utils.RenderJSON(w, "could not refresh access token")
			return
		}
		utils.RenderJSON(w, "could not refresh access token")
	}

	refreshToken := cookie.Value

	refreshedTokenPair, err := h.authService.RefreshAccessToken(refreshToken, h.cfg)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			utils.RenderJSON(w, service.ErrUserNotFound)
			return
		}
		utils.RenderJSON(w, err)
		return
	}

	refreshedAccessToken := refreshedTokenPair["access_token"]
	refreshedRefreshToken := refreshedTokenPair["refresh_token"]

	log.Println("Access token after: ", refreshedAccessToken)
	log.Println("Refresh token after: ", refreshedRefreshToken)

	accessCookie := http.Cookie{
		Name:     "access_token",
		Value:    refreshedAccessToken,
		Path:     "/",
		MaxAge:   h.cfg.AccessTokenMaxAge * 60,
		Secure:   false,
		HttpOnly: true,
		Domain:   "localhost",
	}

	refreshCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshedRefreshToken,
		Path:     "/",
		MaxAge:   h.cfg.RefreshTokenMaxAge * 60,
		Secure:   false,
		HttpOnly: true,
		Domain:   "localhost",
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

	http.SetCookie(w, &accessCookie)
	http.SetCookie(w, &refreshCookie)
	http.SetCookie(w, &loggedInCookie)

	utils.RenderJSON(w, "Tokens refreshed!")
	return
}

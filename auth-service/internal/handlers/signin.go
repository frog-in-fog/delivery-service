package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/frog-in-fog/delivery-system/auth-service/internal/models"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/models/dto"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/service"
	"github.com/frog-in-fog/delivery-system/auth-service/utils"
)

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

package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/frog-in-fog/delivery-system/auth-service/internal/models"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/models/dto"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/storage/sqlite"
	"github.com/frog-in-fog/delivery-system/auth-service/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

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

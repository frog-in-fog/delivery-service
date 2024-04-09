package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/frog-in-fog/delivery-system/auth-service/internal/models/dto"
	"github.com/frog-in-fog/delivery-system/auth-service/utils"
)

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

package handlers

import (
	"net/http"
	"strings"

	"github.com/frog-in-fog/delivery-system/auth-service/internal/models/dto"
	"github.com/frog-in-fog/delivery-system/auth-service/utils"
)

func (h *authHandler) TokenPair(w http.ResponseWriter, r *http.Request) {
	accessTokenHeader := r.Header.Get("Authorization")
	accessToken := strings.TrimPrefix(accessTokenHeader, "Bearer ")
	if len(accessToken) == 0 {
		resp := dto.OneLineResp{
			Data: "empty access token",
		}
		utils.RenderJSON(w, resp)
		return
	}

	res, err := h.authService.TokenPair(accessToken, h.cfg)
	if err != nil {
		resp := dto.OneLineResp{
			Data: err.Error(),
		}
		utils.RenderJSON(w, resp)
		return
	}

	if len(res) != 0 {
		if res == "allowed" {
			resp := dto.OneLineResp{
				Data: "allowed",
			}
			utils.RenderJSON(w, resp)
			return
		} else {
			resp := dto.OneLineResp{
				Data: res,
			}
			utils.RenderJSON(w, resp)
			return
		}
	}

}

package handlers

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func TestLogger(w http.ResponseWriter, r *http.Request) {
	payload := Response{
		Message: "Logger info",
		Error:   "Some error",
	}

	jsonResponse, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

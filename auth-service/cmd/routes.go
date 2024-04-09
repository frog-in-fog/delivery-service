package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handlers) InitRoutes() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/api/v0/register", h.AuthHandler.SignUpUser)
	router.HandleFunc("/api/v0/login", h.AuthHandler.SignInUser)
	router.HandleFunc("/api/v0/logout", h.AuthHandler.LogoutUser)
	router.HandleFunc("/api/v0/tokens", h.AuthHandler.TokenPair)

	return router
}

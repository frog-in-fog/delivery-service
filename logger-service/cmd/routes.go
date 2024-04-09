package main

import (
	"net/http"

	"github.com/frog-in-fog/delivery-system/logger-service/internal/handlers"
	"github.com/gorilla/mux"
)

func InitRoutes() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/logger", handlers.TestLogger)

	return router
}

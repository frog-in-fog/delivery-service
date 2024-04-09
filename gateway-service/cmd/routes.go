package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func InitRoutes() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/", TestHandler)

	return router
}

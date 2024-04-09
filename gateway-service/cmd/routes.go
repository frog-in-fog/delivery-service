package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func InitRoutes() http.Handler {
	router := mux.NewRouter()

	//router.HandleFunc("/login", LoginPage)
	router.HandleFunc("/login", LoginHandler)
	router.HandleFunc("/logger", Authenticate(Proxy("/logger", "http://host.docker.internal:8081")))

	return router
}

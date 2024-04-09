package main

import (
	"net/http"

	"github.com/frog-in-fog/delivery-system/gateway-service/handlers"
	"github.com/gorilla/mux"
)

func InitRoutes() http.Handler {
	router := mux.NewRouter()

	//router.HandleFunc("/login", LoginPage)
	router.HandleFunc("/login", handlers.LoginHandler)
	router.HandleFunc("/logger", handlers.Authenticate(handlers.Proxy("/logger", "http://host.docker.internal:8081")))

	return router
}

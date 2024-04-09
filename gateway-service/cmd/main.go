package main

import (
	"log"
	"net/http"
)

const (
	port = "8080"
)

func main() {
	router := InitRoutes()
	log.Printf("Gateway service is running on port: %s", port)
	http.ListenAndServe(":8080", router)
}

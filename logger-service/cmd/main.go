package main

import (
	"log"
	"net/http"
)

const (
	port = "8081"
)

func main() {
	router := InitRoutes()
	log.Printf("Logger service is running on port: %s", port)
	http.ListenAndServe(":8081", router)
}

package main

import (
	"context"
	"fmt"
	"github.com/frog-in-fog/delivery-system/auth-service/cmd"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/config"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/storage/sqlite"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	server *http.Server
}

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Error parsing env variables: %v", err)
	}

	// migrate database
	//if err := cmd.Migrate(&cfg); err != nil {
	//	log.Fatalf("error migrating db: %v", err)
	//}

	// init sqlite and redis
	sqliteStorage, err := sqlite.NewSQLiteStorage(cfg.StoragePath)
	if err != nil {
		log.Fatalf("error creating storage: %v", err)
	}

	if err = cmd.NewRedisConnection(&cfg); err != nil {
		log.Fatalf("error connecting to redis: %v", err)
	}

	// init http handlers
	httpHandlers := cmd.InitHttpHandlers(&cfg, sqliteStorage)
	// launch http server
	server := new(Server)

	go func() {
		if err = server.run(cfg.WebPort, httpHandlers.InitRoutes()); err != nil {
			log.Fatalf("error occured running http server: %s", err.Error())
		}
	}()

	log.Println("Application aaa started on port: ", cfg.WebPort)

	gracefulShutdown(*server)

	log.Println("Application stopped")
}

func (s *Server) run(port string, handler http.Handler) error {
	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return s.server.ListenAndServe()
}

func (s *Server) shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func gracefulShutdown(server Server) {
	q := make(chan os.Signal, 1)
	signal.Notify(q, syscall.SIGTERM, syscall.SIGINT)
	<-q

	if err := server.shutdown(context.Background()); err != nil {
		log.Printf("error occured on server shutting down: %s", err.Error())
	}
}

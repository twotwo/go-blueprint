package server

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"rest_api/app/database"
	"rest_api/app/utils"
)

type Server struct {
	port int
	db   database.Service
}

func NewServer() *http.Server {
	cfg, err := utils.GetConfigInstance()
	if err != nil {
		panic(err)
	}
	NewServer := &Server{
		port: cfg.ServicePort,
		db:   database.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

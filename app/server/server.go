package server

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/twotwo/go-blueprint/app/database"
	"github.com/twotwo/go-blueprint/app/utils"
)

type Server struct {
	port int
	auth func(next http.Handler) http.Handler
	db   database.Service
}

func NewServer() *http.Server {
	cfg, err := utils.GetConfigInstance()
	if err != nil {
		panic(err)
	}
	NewServer := &Server{
		port: cfg.ServicePort,
		auth: BasicAuth("example", map[string]string{"admin": "admin"}),
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

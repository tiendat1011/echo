package server

import (
	"fmt"
	"net/http"
	"test/cmd/web"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	productStore *web.ProductStore
	port         int
}

func NewServer() *http.Server {
	port := 8080
	NewServer := &Server{
		port:         port,
		productStore: web.NewProductStore(),
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

package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/geophpherie/boot-dev-rss-feed-agg/internal/database"
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	db *database.Queries
}

func New() *http.Server {
	NewServer := &Server{
		db: database.GetQueries(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%v", os.Getenv("PORT")),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

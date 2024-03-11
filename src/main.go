package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/geophpherie/boot-dev-rss-feed-agg/src/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load("../.env")

	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DB_CONN_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Unable to open database.")
	}

	dbQueries := database.New(db)

	apiConfig := apiConfig{db: dbQueries}

	serveMux := http.NewServeMux()
	corsMux := middlewareCors(serveMux)

	serveMux.HandleFunc("GET /v1/readiness", handleReadiness)
	serveMux.HandleFunc("GET /v1/err", handleErr)
	serveMux.HandleFunc("POST /v1/users", apiConfig.createUser)
	serveMux.HandleFunc("GET /v1/users", apiConfig.middlewareAuth(apiConfig.getUser))
	serveMux.HandleFunc("POST /v1/feeds", apiConfig.middlewareAuth(apiConfig.createFeed))

	server := http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	fmt.Printf("SERVICE ON localhost:%v\n", port)
	log.Fatal(server.ListenAndServe())
}

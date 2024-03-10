package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../.env")

	port := os.Getenv("PORT")

	serveMux := http.NewServeMux()
	corsMux := middlewareCors(serveMux)

	serveMux.HandleFunc("GET /v1/readiness", handleReadiness)
	serveMux.HandleFunc("GET /v1/err", handleErr)

	server := http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	fmt.Printf("SERVICE ON localhost:%v\n", port)
	log.Fatal(server.ListenAndServe())
}

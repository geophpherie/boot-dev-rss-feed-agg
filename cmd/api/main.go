package main

import (
	"fmt"
	"log"

	"github.com/geophpherie/boot-dev-rss-feed-agg/internal/server"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	server := server.New()
	fmt.Printf("SERVICE ON %v\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}

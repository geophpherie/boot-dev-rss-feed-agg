package main

import (
	"fmt"
	"log"

	"github.com/geophpherie/boot-dev-rss-feed-agg/internal/database"
	"github.com/geophpherie/boot-dev-rss-feed-agg/internal/rss"
	"github.com/geophpherie/boot-dev-rss-feed-agg/internal/server"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	db := *database.GetQueries()

	server := server.New(&db)

	go rss.ScrapeFeeds(&db)

	fmt.Printf("SERVICE %v\n", server.Addr)
	log.Fatal(server.ListenAndServe())

}

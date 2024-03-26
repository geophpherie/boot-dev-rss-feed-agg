package main

import (
	"context"
	"fmt"

	"github.com/geophpherie/boot-dev-rss-feed-agg/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db := database.GetQueries()

	res, _ := db.GetNextFeedsToFetch(context.Background(), 2)

	for _, result := range res {
		fmt.Println(result)
		fmt.Println(result.LastFetchedAt.Time)
	}

	fmt.Println()
}

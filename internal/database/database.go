package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func GetQueries() *Queries {
	dbUrl := os.Getenv("DB_CONN_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	return New(db)
}

package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// type service struct {
// 	db *sql.DB
// }

// var (
// 	database = os.Getenv("DB_DATABASE")
// 	password = os.Getenv("DB_PASSWORD")
// 	username = os.Getenv("DB_USERNAME")
// 	port     = os.Getenv("DB_PORT")
// 	host     = os.Getenv("DB_HOST")
// )

func GetNew() *Queries {
	// connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)

	dbUrl := os.Getenv("DB_CONN_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := New(db)
	// s := &service{db: db}
	return dbQueries
}

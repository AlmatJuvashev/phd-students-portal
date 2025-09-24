package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// MustOpen connects to Postgres or exits on error.
func MustOpen(url string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	return db
}

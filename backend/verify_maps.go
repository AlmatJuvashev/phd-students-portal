package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load(".env")
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

    var mapCount int
    db.Get(&mapCount, "SELECT count(*) FROM program_versions")
    fmt.Printf("Program Versions (Maps): %d\n", mapCount)

    var defCount int
    db.Get(&defCount, "SELECT count(*) FROM program_version_node_definitions")
    fmt.Printf("Journey Node Definitions: %d\n", defCount)
    
    // Check if any map exists
    var mapIDs []string
    db.Select(&mapIDs, "SELECT id FROM program_versions")
    fmt.Println("Map IDs:", mapIDs)
}

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

	var count int
	err = db.Get(&count, "SELECT count(*) FROM courses WHERE tenant_id = (SELECT id FROM tenants WHERE slug='demo')")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Courses in Demo Tenant: %d\n", count)

    var nodeCount int
    err = db.Get(&nodeCount, "SELECT count(*) FROM journey_node_definitions") 
    fmt.Printf("Journey Nodes (Total): %d\n", nodeCount)
    
    // Check if courses have localized titles (starts with {)
    var titles []string
    err = db.Select(&titles, "SELECT title FROM courses WHERE tenant_id = (SELECT id FROM tenants WHERE slug='demo')")
    fmt.Println("Course Titles:", titles)
}

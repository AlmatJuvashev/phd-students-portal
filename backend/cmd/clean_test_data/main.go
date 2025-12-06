package main

import (
	"log"
	"os"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/db"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("../.env")
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	conn := db.MustOpen(dbURL)
	defer conn.Close()

	// Delete node instances and journey states for test users
	testUsers := []string{"ts5251", "ta2087"}
	
	// Delete node instances
	query1 := `DELETE FROM node_instances WHERE user_id IN (SELECT id FROM users WHERE username IN (?))`
	query1, args1, err := sqlx.In(query1, testUsers)
	if err != nil {
		log.Fatal(err)
	}
	query1 = conn.Rebind(query1)
	if _, err := conn.Exec(query1, args1...); err != nil {
		log.Printf("Error deleting node instances: %v", err)
	}

	// Delete journey states
	query2 := `DELETE FROM journey_states WHERE user_id IN (SELECT id FROM users WHERE username IN (?))`
	query2, args2, err := sqlx.In(query2, testUsers)
	if err != nil {
		log.Fatal(err)
	}
	query2 = conn.Rebind(query2)
	if _, err := conn.Exec(query2, args2...); err != nil {
		log.Printf("Error deleting journey states: %v", err)
	}
	if err != nil {
		log.Printf("Error cleaning test data: %v", err)
	} else {
		log.Println("Test data cleaned successfully")
	}
}

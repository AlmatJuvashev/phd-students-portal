package main

import (
	"flag"
	"log"
	"os"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/db"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/joho/godotenv"
)

func main() {
	// Load env from parent directory if present (for local dev)
	_ = godotenv.Load("../.env")

	playbookPath := flag.String("path", "../frontend/src/playbooks/playbook.json", "path to playbook.json")
	flag.Parse()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	conn := db.MustOpen(dbURL)
	defer conn.Close()

	log.Printf("Seeding playbook from %s...", *playbookPath)
	_, err := playbook.EnsureActive(conn, *playbookPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Playbook seeded successfully")
}

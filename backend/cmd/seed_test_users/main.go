package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/auth"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("warning: .env not found, relying on existing env vars")
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	conn, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Create admin user
	adminPassword := "meadow-pluto-pioneer48"
	adminHash, _ := auth.HashPassword(adminPassword)
	_, err = conn.Exec(`INSERT INTO users (username, email, first_name, last_name, role, password_hash, is_active)
		VALUES ('ta2087', 'ta2087@test.kaznmu.kz', 'Test', 'Admin', 'admin', $1, true)
		ON CONFLICT (username) DO UPDATE SET password_hash = $1, is_active = true`, adminHash)
	if err != nil {
		log.Printf("Admin insert/update error: %v", err)
	} else {
		fmt.Println("Admin created/updated: ta2087 / meadow-pluto-pioneer48")
	}

	// Create student user
	studentPassword := "pioneer-canvas-silver52"
	studentHash, _ := auth.HashPassword(studentPassword)
	_, err = conn.Exec(`INSERT INTO users (username, email, first_name, last_name, role, password_hash, is_active)
		VALUES ('ts5251', 'ts5251@test.kaznmu.kz', 'Test', 'Student', 'student', $1, true)
		ON CONFLICT (username) DO UPDATE SET password_hash = $1, is_active = true`, studentHash)
	if err != nil {
		log.Printf("Student insert/update error: %v", err)
	} else {
		fmt.Println("Student created/updated: ts5251 / pioneer-canvas-silver52")
	}

	fmt.Println("Done!")
}

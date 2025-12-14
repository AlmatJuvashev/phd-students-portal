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
	// Try loading env files
	if err := godotenv.Load(".env"); err != nil {
		if err := godotenv.Load("../../.env"); err != nil {
			log.Println("warning: .env not found, relying on existing env vars")
		}
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	// Connect to DB
	conn, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Tenant B constant IDs
	tenantBID := "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
	adminBID := "bbbbbbbb-bbbb-bbbb-bbbb-000000000001"
	studentBID := "bbbbbbbb-bbbb-bbbb-bbbb-000000000002"
	
	defaultPass := "demopassword123!"
	hashedPass, _ := auth.HashPassword(defaultPass)

	log.Println("Seeding Tenant B...")

	// 1. Ensure Tenant B
	_, err = conn.Exec(`
		INSERT INTO tenants (id, slug, name, tenant_type, is_active, enabled_services, primary_color, secondary_color, app_name)
		VALUES ($1, 'tenant-b', 'Tenant B University', 'university', true, ARRAY['chat'], '#FF5733', '#C70039', 'Tenant B Portal')
		ON CONFLICT (id) DO UPDATE SET is_active = true
	`, tenantBID)
	if err != nil {
		log.Printf("Failed to seed Tenant B: %v", err) // might exist with different slug?
	}

	// 2. Ensure Users (Admin B, Student B)
	ensureUser(conn, adminBID, "admin.b", "admin.b@test.com", "Admin", "B", "admin", hashedPass, tenantBID)
	ensureUser(conn, studentBID, "student.b", "student.b@test.com", "Student", "B", "student", hashedPass, tenantBID)

	fmt.Println("Multitenancy Seeding Completed!")
}

func ensureUser(db *sqlx.DB, id, username, email, first, last, role, hash, tenantID string) {
	_, err := db.Exec(`
		INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, true)
		ON CONFLICT (id) DO UPDATE SET role = $6, password_hash = $7
	`, id, username, email, first, last, role, hash)
	if err != nil {
		// Try insert without ID if ID conflict on username?
		// But we hardcoded IDs.
		// If username conflict with different ID, update ID? No, can't update PK.
		// Assuming fresh setup or consistent IDs.
		// But what if 'admin.b' exists with diff ID?
		// We'll trust the ID.
		log.Printf("Failed to ensure user %s: %v", username, err)
	}

	// Membership
	_, err = db.Exec(`
		INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary)
		VALUES ($1, $2, $3, true)
		ON CONFLICT (user_id, tenant_id) DO NOTHING
	`, id, tenantID, role)
	if err != nil {
		log.Printf("Failed to ensure membership for %s: %v", username, err)
	}
}

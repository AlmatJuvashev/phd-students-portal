package main

import (
	"context"
	"log"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found")
	}
	cfg := config.MustLoad()

	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	ctx := context.Background()

	// 1. Seed Permissions
	permissions := []string{
		"course.view", "course.create", "course.edit", "course.delete",
		"grade.view", "grade.edit",
		"user.view", "user.create", "user.edit", "user.delete",
		"role.manage", "*",
	}

	for _, perm := range permissions {
		// Pass perm twice for slug and description
		_, err := db.ExecContext(ctx, "INSERT INTO permissions (slug, description) VALUES ($1, $2) ON CONFLICT (slug) DO NOTHING", perm, perm)
		if err != nil {
			log.Printf("Error seeding permission %s: %v", perm, err)
		}
	}
	log.Println("Permissions seeded.")

	// 2. Seed Roles
	roles := map[string][]string{
		"Superadmin": {"*"}, // Wildcard
		"Instructor": {"course.view", "course.edit", "grade.view", "grade.edit"},
		"Student":    {"course.view", "grade.view"},
		"Admin":      {"user.view", "user.create", "user.edit", "course.view", "course.create"},
	}

	for roleName, perms := range roles {
		// Create Role
		var roleID uuid.UUID
		err := db.QueryRowContext(ctx, "INSERT INTO roles (name, is_system_role) VALUES ($1, true) ON CONFLICT (name, tenant_id) DO UPDATE SET name=EXCLUDED.name RETURNING id", roleName).Scan(&roleID)
		if err != nil {
			log.Printf("Error creating role %s: %v", roleName, err)
			continue
		}

		// Assign Permissions
		for _, p := range perms {
			_, err := db.ExecContext(ctx, "INSERT INTO role_permissions (role_id, permission_slug) VALUES ($1, $2) ON CONFLICT DO NOTHING", roleID, p)
			if err != nil {
				log.Printf("Error assigning permission %s to role %s: %v", p, roleName, err)
			}
		}
	}
	log.Println("Roles seeded.")

	// 3. Backfill Users
	// Cast enum to text to compare with empty string safely
	rows, err := db.QueryxContext(ctx, "SELECT id, role::text FROM users WHERE role IS NOT NULL")
	if err != nil {
		log.Fatalf("Failed to fetch users: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u struct {
			ID   uuid.UUID `db:"id"`
			Role string    `db:"role"`
		}
		if err := rows.StructScan(&u); err != nil {
			continue
		}

		// Map legacy role string to new Role Name
		// Simple mapping: capitalize first letter, or direct map
		var targetRoleName string
		switch u.Role {
		case "superadmin": targetRoleName = "Superadmin"
		case "admin": targetRoleName = "Admin"
		case "instructor", "advisor": targetRoleName = "Instructor" // Advisor -> Instructor permissions for now
		case "student": targetRoleName = "Student"
		default: 
			log.Printf("Unknown legacy role %s for user %s, skipping", u.Role, u.ID)
			continue
		}

		// Get Role ID
		var roleID uuid.UUID
		err := db.QueryRowContext(ctx, "SELECT id FROM roles WHERE name=$1", targetRoleName).Scan(&roleID)
		if err != nil {
			log.Printf("Role %s not found for user %s", targetRoleName, u.ID)
			continue
		}

		// Assign Global Role
		_, err = db.ExecContext(ctx, 
			"INSERT INTO user_context_roles (user_id, role_id, context_type, context_id) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING",
			u.ID, roleID, models.ContextGlobal, uuid.Nil,
		)
		if err != nil {
			log.Printf("Failed to assign role to user %s: %v", u.ID, err)
		}
	}
	log.Println("User backfill complete.")
}

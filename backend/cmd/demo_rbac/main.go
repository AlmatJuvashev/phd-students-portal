package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// 1. Setup
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5435/phd?sslmode=disable"
	}
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	rbacRepo := repository.NewSQLRBACRepository(db)
	authzSvc := services.NewAuthzService(rbacRepo)
	ctx := context.Background()

	// 2. Create Dummy Data (In-Memory UUIDs for simulation, or insert if needed)
	// For this demo, we will insert real data to test the SQL queries, then clean up?
	// Actually, let's just use the existing repositories to Create dummy users/courses if they don't exist,
	// or just generate random IDs and Insert mocked relations directly into rbac tables for speed.
	
	// Let's create proper data to be safe.
	userRepo := repository.NewSQLUserRepository(db)
	
	// Users
	adminID := createOrGetUser(ctx, userRepo, "admin_demo", "Admin User")
	taID := createOrGetUser(ctx, userRepo, "ta_demo", "Teaching Assistant")
	studentID := createOrGetUser(ctx, userRepo, "student_demo", "Regular Student")

	// Contexts
	mathCourseID := uuid.New()
	histCourseID := uuid.New()

	fmt.Printf("\n=== RBAC Context-Awareness Demo ===\n")
	fmt.Printf("Users:\n  Global Admin: %s\n  Math TA:      %s\n  Student:      %s\n", adminID, taID, studentID)
	fmt.Printf("Contexts:\n  Math 101:     %s\n  History 202:  %s\n", mathCourseID, histCourseID)

	// 3. Assign Roles
	// Clean up previous demo roles first
	cleanRoles(db, adminID, taID, studentID)

	assignRole(ctx, db, adminID, "Superadmin", models.ContextGlobal, uuid.Nil)
	assignRole(ctx, db, taID, "Instructor", models.ContextCourse, mathCourseID) // TA only for Math
	assignRole(ctx, db, studentID, "Student", models.ContextCourse, histCourseID) // Student only for Hist

	fmt.Println("\n[Configuration]")
	fmt.Println("  * Admin   -> Superadmin (Global)")
	fmt.Println("  * TA      -> Instructor (Math 101 ONLY)")
	fmt.Println("  * Student -> Student    (History 202 ONLY)")

	// 4. Test Matrix
	scenarios := []struct {
		User     string
		UserID   uuid.UUID
		Action   string
		Context  string
		TargetID uuid.UUID
		Expect   bool
	}{
		// Admin Checks
		{"Admin", adminID, "course.edit", "Math 101", mathCourseID, true},
		{"Admin", adminID, "course.edit", "History 202", histCourseID, true},
		
		// TA Checks
		{"TA", taID, "course.edit", "Math 101", mathCourseID, true},
		{"TA", taID, "course.edit", "History 202", histCourseID, false}, // Should fail!
		
		// Student Checks
		{"Student", studentID, "course.view", "Math 101", mathCourseID, false}, // Should fail!
		{"Student", studentID, "course.view", "History 202", histCourseID, true},
		{"Student", studentID, "course.edit", "History 202", histCourseID, false}, // Student can view but not edit
	}

	fmt.Println("\n[Permission Check Matrix]")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "User\tAction\tTarget\tExpected\tResult\tStatus")
	fmt.Fprintln(w, "----\t------\t------\t--------\t------\t------")

	for _, s := range scenarios {
		allowed, err := authzSvc.HasPermission(ctx, s.UserID, s.Action, models.ContextCourse, s.TargetID)
		if err != nil {
			log.Printf("Error checking permission: %v", err)
			continue
		}
		
		status := "✅ PASS"
		if allowed != s.Expect {
			status = "❌ FAIL"
		}
		
		fmt.Fprintf(w, "%s\t%s\t%s\t%v\t%v\t%s\n", s.User, s.Action, s.Context, s.Expect, allowed, status)
	}
	w.Flush()
	fmt.Println("\nDemo Complete.")
}

func createOrGetUser(ctx context.Context, repo repository.UserRepository, username, name string) uuid.UUID {
	u, err := repo.GetByEmail(ctx, username+"@demo.com")
	if err == nil {
		id, _ := uuid.Parse(u.ID)
		return id
	}
	// Create
	userid := uuid.New()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5435/phd?sslmode=disable"
	}
	db, _ := sqlx.Connect("postgres", dbURL)
	_, err = db.Exec("INSERT INTO users (id, username, email, first_name, last_name, password_hash, is_active, role) VALUES ($1, $2, $3, $4, 'Demo', 'hash', true, 'student')", userid, username, username+"@demo.com", name)
	if err != nil {
		log.Fatal(err)
	}
	return userid
}

func cleanRoles(db *sqlx.DB, ids ...uuid.UUID) {
	query, args, _ := sqlx.In("DELETE FROM user_context_roles WHERE user_id IN (?)", ids)
	query = db.Rebind(query)
	db.Exec(query, args...)
}

func assignRole(ctx context.Context, db *sqlx.DB, userID uuid.UUID, roleName string, matchType string, contextID uuid.UUID) {
	var roleID uuid.UUID
	err := db.QueryRow("SELECT id FROM roles WHERE name=$1", roleName).Scan(&roleID)
	if err != nil {
		log.Fatalf("Role %s not found: %v", roleName, err)
	}

	_, err = db.ExecContext(ctx, `
		INSERT INTO user_context_roles (user_id, role_id, context_type, context_id)
		VALUES ($1, $2, $3, $4)
	`, userID, roleID, matchType, contextID)
	if err != nil {
		log.Fatalf("Failed to assign role: %v", err)
	}
}

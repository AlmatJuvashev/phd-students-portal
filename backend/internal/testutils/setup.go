package testutils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// SetupTestDB connects to the test database, applies migrations, and returns the connection.
// It assumes TEST_DATABASE_URL is set, or defaults to a local test DB.
// It returns a cleanup function that should be deferred.
func SetupTestDB() (*sqlx.DB, func()) {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		// Default to localhost test DB if not set
		// IMPORTANT: Uses phd_test (not phd) to avoid wiping demo data
		dbURL = "postgres://postgres:postgres@localhost:5435/phd_test?sslmode=disable"
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to test DB: %v", err)
	}
	db.SetMaxOpenConns(1)

	// Find project root to locate migrations
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	// basepath is .../backend/internal/testutils
	// migrations are in .../backend/db/migrations
	migrationsPath := filepath.Join(basepath, "../../db/migrations")

	// Run migrations
	m, err := migrate.New(
		"file://"+migrationsPath,
		dbURL,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	// Close migrate instance to release connection
	srcErr, dbErr := m.Close()
	if srcErr != nil {
		log.Printf("Migrate source close error: %v", srcErr)
	}
	if dbErr != nil {
		log.Printf("Migrate db close error: %v", dbErr)
	}

	// Clean DB on start to ensure clean slate
	cleanupDB(db)

	return db, func() {
		cleanupDB(db)
		db.Close()
	}
}

func cleanupDB(db *sqlx.DB) {
	// Cleanup logic
	// Truncating tables.
	tables := []string{
		"users", "tenants", "user_tenant_memberships",
		"programs", "specialties", "cohorts", "departments", "specialty_programs",
		"node_instances", "node_instance_form_revisions", "node_instance_slots", "node_instance_slot_attachments", "node_outcomes", "node_events",
		"chat_rooms", "chat_room_members", "chat_messages", "chat_room_read_status",
		"notifications",
		"events", "event_attendees",
		"checklist_modules", "checklist_steps", "student_steps",
		"documents", "document_versions", "comments",
		"playbook_versions", "playbook_active_version",
		"contacts",
		"journey_states", "student_advisors", "node_deadlines",
		"profile_submissions", "profile_audit_log", "email_verification_tokens", "rate_limit_events",
	}
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			// Log but don't fail, maybe table doesn't exist yet
			log.Printf("Failed to truncate table %s: %v", table, err)
		}
	}
	// Truncate node_state_transitions and admin_notifications
	db.Exec("TRUNCATE TABLE node_state_transitions CASCADE")
	db.Exec("TRUNCATE TABLE admin_notifications CASCADE")

	// Seed default transitions (from migration 0006)
	db.Exec(`INSERT INTO node_state_transitions(from_state, to_state, allowed_roles) VALUES
		('active','submitted', ARRAY['student']),
		('submitted','needs_fixes', ARRAY['advisor','secretary','chair','admin']),
		('submitted','done', ARRAY['advisor','secretary','chair','admin']),
		('needs_fixes','submitted', ARRAY['student']),
		('done','submitted', ARRAY['admin'])
		ON CONFLICT DO NOTHING`)

	// Seed default tenant for tests
	db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type, is_active) 
		VALUES ('00000000-0000-0000-0000-000000000001', 'default-test', 'Default Test Tenant', 'university', true)
		ON CONFLICT DO NOTHING`)
}

func GetTestConfig() config.AppConfig {
	return config.AppConfig{
		RedisURL:        "redis://localhost:6379",
		Port:            "8080",
		Env:             "test",
		JWTSecret:       "test-secret",
		JWTExpDays:      1,
		DatabaseURL:     os.Getenv("TEST_DATABASE_URL"),
		UploadDir:       "./test_uploads",
		FileUploadMaxMB: 10,
		SMTPHost:        "localhost",
		SMTPPort:        "1025",
		FrontendBase:    "http://localhost:3000",
		S3Endpoint:      "http://localhost:9000",
		S3Bucket:        "test-bucket",
		ServerURL:       "http://localhost:8080",
	}
}

// CreateTestUser creates a test user and returns their ID.
func CreateTestUser(t *testing.T, db *sqlx.DB, username, role string) string {
	t.Helper()
	id := uuid.New().String()
	_, err := db.Exec(`
		INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active)
		VALUES ($1, $2, $3, 'Test', 'User', $4, 'testhash', true)
		ON CONFLICT (username) DO UPDATE SET id = EXCLUDED.id RETURNING id
	`, id, username, username+"@test.com", role)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return id
}

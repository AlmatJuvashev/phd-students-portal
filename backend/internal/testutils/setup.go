package testutils

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"sync"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	testDB *sqlx.DB
	testOnce sync.Once
)

// SetupTestDB connects to the test database, applies migrations, and returns the connection.
// It assumes TEST_DATABASE_URL is set, or defaults to a local test DB.
// It returns a cleanup function that should be deferred.
func SetupTestDB() (*sqlx.DB, func()) {
	// Check if running in CI or explicitly requested to use Testcontainers
	useContainer := os.Getenv("USE_TESTCONTAINERS") == "true"
	
	if useContainer {
		ctx := context.Background()
		dbUser := "postgres"
		dbPassword := "postgres"
		dbName := "phd_test"

		pgContainer, err := postgres.Run(ctx,
			"postgres:15-alpine",
			postgres.WithDatabase(dbName),
			postgres.WithUsername(dbUser),
			postgres.WithPassword(dbPassword),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2).
					WithStartupTimeout(5*time.Minute)),
		)
		if err != nil {
			log.Fatalf("[SetupTestDB] Failed to start postgres container: %v", err)
		}

		connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
		if err != nil {
			log.Fatalf("[SetupTestDB] Failed to get connection string: %v", err)
		}

		db, err := sqlx.Connect("postgres", connStr)
		if err != nil {
			log.Fatalf("[SetupTestDB] Failed to connect to container DB: %v", err)
		}

		// Run migrations
		// We need to construct file source url
		_, b, _, _ := runtime.Caller(0)
		migrationsPath := filepath.Join(filepath.Dir(b), "../../db/migrations")
		
		m, err := migrate.New("file://"+migrationsPath, connStr)
		if err != nil {
			log.Fatalf("[SetupTestDB] Failed to init migrate: %v", err)
		}
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("[SetupTestDB] Migration failed: %v", err)
		}
		m.Close()

		testDB = db
		
		return db, func() {
			db.Close()
			if err := pgContainer.Terminate(ctx); err != nil {
				log.Printf("[SetupTestDB] Failed to terminate container: %v", err)
			}
		}
	}

	// Fallback to local DB (Original Logic) or if useContainer is false
	// ... (Existing implementation for local DB)
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5435/phd_test?sslmode=disable"
	}

	// SAFETY CHECK: Parse URL to ensure we are NOT connecting to production 'phd' DB
	u, err := url.Parse(dbURL)
	if err != nil {
		log.Fatalf("[SetupTestDB] Invalid TEST_DATABASE_URL: %v", err)
	}
	dbPath := strings.TrimPrefix(u.Path, "/")
	if dbPath == "phd" {
		log.Fatalf("[SetupTestDB] CRITICAL SAFETY BLOCK: Attempting to run tests against production 'phd' database! Aborting immediately.")
	}
	if !strings.Contains(dbPath, "test") {
		log.Printf("[SetupTestDB] WARNING: Running tests against database '%s' which does not contain 'test'. This is risky.", dbPath)
	}

	// Identify package name for database isolation (capture BEFORE testOnce.Do)
	_, filename, _, _ := runtime.Caller(1)
	pkgName := filepath.Base(filepath.Dir(filename))
	log.Printf("[SetupTestDB] Caller file: %s, pkgName: %s", filename, pkgName)
	targetDB := "phd_test_" + pkgName
	targetDB = strings.ReplaceAll(targetDB, "-", "_")
	targetDB = strings.ReplaceAll(targetDB, ".", "_")

	testOnce.Do(func() {
		ctx := context.Background()
		// 1. Connect to 'postgres' to manage databases
		adminDSN := strings.ReplaceAll(dbURL, "phd_test", "postgres")
		adminDB, err := sqlx.Connect("postgres", adminDSN)
		if err != nil {
			log.Fatalf("[SetupTestDB] Failed to connect to admin DB: %v", err)
		}
		defer adminDB.Close()

		// Use a single connection for session-level advisory lock
		conn, err := adminDB.Connx(ctx)
		if err != nil {
			log.Fatalf("[SetupTestDB] Failed to get admin conn: %v", err)
		}
		defer conn.Close()

		// Global lock for DB management (advisory lock on ID 123456)
		if _, err := conn.ExecContext(ctx, "SELECT pg_advisory_lock(123456)"); err != nil {
			log.Fatalf("[SetupTestDB] Failed to acquire advisory lock: %v", err)
		}
		defer conn.ExecContext(ctx, "SELECT pg_advisory_unlock(123456)")

		// 2. Ensure base template exists and is migrated
		baseDB := "phd_test_base"

		// Check for base template
		var baseExists bool
		conn.GetContext(ctx, &baseExists, "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1)", baseDB)

		if baseExists {
			log.Printf("[SetupTestDB] Dropping stale base template: %s", baseDB)
			_, err := conn.ExecContext(ctx, fmt.Sprintf("SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '%s' AND pid <> pg_backend_pid()", baseDB))
			if err != nil {
				log.Printf("[SetupTestDB] Warning: Failed to terminate connections to %s: %v", baseDB, err)
			}
			_, err = conn.ExecContext(ctx, fmt.Sprintf("DROP DATABASE %s", baseDB))
			if err != nil {
				log.Fatalf("[SetupTestDB] Failed to drop base template %s: %v", baseDB, err)
			}
			baseExists = false
		}
		
		if !baseExists {
			log.Printf("[SetupTestDB] Creating base template database: %s", baseDB)
			_, err := conn.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s", baseDB))
			if err != nil {
				log.Fatalf("[SetupTestDB] Failed to create base template %s: %v", baseDB, err)
			}
		}

		// Always migrate base template to ensure it's up to date
		baseURL := strings.ReplaceAll(dbURL, "phd_test", baseDB)
		_, b, _, _ := runtime.Caller(0)
		migrationsPath := filepath.Join(filepath.Dir(b), "../../db/migrations")
		
		m, err := migrate.New("file://"+migrationsPath, baseURL)
		if err != nil {
			log.Fatalf("[SetupTestDB] Failed to migrate base template: %v", err)
		}
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("[SetupTestDB] Migration failed on base template: %v", err)
		}
		m.Close()

		// 3. Create target DB from template
		// Always drop if exists to ensure fresh clone from updated base
		var targetExists bool
		conn.GetContext(ctx, &targetExists, "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1)", targetDB)
		
		if targetExists {
			log.Printf("[SetupTestDB] Dropping stale database %s", targetDB)
			// Terminate connections
			conn.ExecContext(ctx, fmt.Sprintf("SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '%s' AND pid <> pg_backend_pid()", targetDB))
			conn.ExecContext(ctx, fmt.Sprintf("DROP DATABASE %s", targetDB))
		}

		log.Printf("[SetupTestDB] Cloning database %s from template %s", targetDB, baseDB)
		// Ensure no active connections to template (Postgres requirement for cloning) or wait/retry?
		// We can terminate connections on baseDB too just in case
		conn.ExecContext(ctx, fmt.Sprintf("SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '%s' AND pid <> pg_backend_pid()", baseDB))
		_, err = conn.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s TEMPLATE %s", targetDB, baseDB))
		if err != nil {
			log.Fatalf("[SetupTestDB] Failed to clone database %s: %v", targetDB, err)
		}

		// 4. Connect to the isolated database
		pkgDSN := strings.ReplaceAll(dbURL, "phd_test", targetDB)
		packageDB, err := sqlx.Connect("postgres", pkgDSN)
		if err != nil {
			log.Fatalf("[SetupTestDB] Failed to connect to isolated DB %s: %v", targetDB, err)
		}
		packageDB.SetMaxOpenConns(20)
		packageDB.SetMaxIdleConns(10)
		testDB = packageDB
	})

	// Clean DB inside the isolated database
	cleanupDB(testDB, "")

	return testDB, func() {
		cleanupDB(testDB, "")
	}
}

func cleanupDB(db *sqlx.DB, schema string) {
	// If schema is provided, ensure search_path is set (just in case)
	if schema != "" {
		db.Exec(fmt.Sprintf("SET search_path TO %s, public", schema))
	}
	// Cleanup logic - order matters for foreign key constraints
	// Clean child tables first, then parent tables
	tables := []string{
		// Child tables first (those with foreign keys)
		"node_instance_slot_attachments", "node_instance_slots", "node_instance_form_revisions", "node_outcomes", "node_events",
		"node_instances",
		"journey_states", "node_deadlines",
		"student_advisors",
		"chat_room_read_status", "chat_messages", "chat_room_members", "chat_rooms",
		"event_attendees", "events",
		"student_steps", "checklist_steps", "checklist_modules",
		"document_versions", "documents",
		"comments",
		"notifications",
		"admin_notifications",
		"user_tenant_memberships",
		"profile_submissions", "profile_audit_log", "email_verification_tokens", "rate_limit_events",
		"specialty_programs",
		"playbook_active_version", "playbook_versions",
		"contacts",
		// Parent tables last
		"users",
		"programs", "specialties", "cohorts", "departments",
		"tenants",
	}
	
	// Use a single TRUNCATE command for all tables - much faster than a loop
	allTables := strings.Join(tables, ", ")
	_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", allTables))
	if err != nil {
		log.Printf("[cleanupDB] Failed to truncate tables: %v", err)
	}
	
	// Truncate node_state_transitions separately (it's often handled differently)
	db.Exec("TRUNCATE TABLE node_state_transitions CASCADE")

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


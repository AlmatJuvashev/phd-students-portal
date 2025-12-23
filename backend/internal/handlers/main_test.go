package handlers_test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Ensure we use the dedicated test database
	if os.Getenv("TEST_DATABASE_URL") == "" {
		os.Setenv("TEST_DATABASE_URL", "postgres://postgres:postgres@localhost:5435/phd_test?sslmode=disable")
	}

	// Set S3 env vars for document tests
	os.Setenv("S3_BUCKET", "test-bucket")
	os.Setenv("S3_ACCESS_KEY", "test-key")
	os.Setenv("S3_SECRET_KEY", "test-secret")
	os.Setenv("S3_REGION", "us-east-1")

	code := m.Run()
	os.Exit(code)
}


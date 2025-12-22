package repository

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Ensure we use the test database if not explicitly set
	if os.Getenv("TEST_DATABASE_URL") == "" {
		os.Setenv("TEST_DATABASE_URL", "postgres://postgres:postgres@localhost:5435/phd_test?sslmode=disable")
	}

	code := m.Run()
	os.Exit(code)
}

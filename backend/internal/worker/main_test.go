package worker

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Ensure we use the dedicated test database
	if os.Getenv("TEST_DATABASE_URL") == "" {
		os.Setenv("TEST_DATABASE_URL", "postgres://postgres:postgres@localhost:5435/phd_test?sslmode=disable")
	}

	code := m.Run()
	os.Exit(code)
}

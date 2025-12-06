package testutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupTestDB(t *testing.T) {
	db, teardown := SetupTestDB()
	defer teardown()

	assert.NotNil(t, db)
	err := db.Ping()
	assert.NoError(t, err)

	// Check if a table exists (e.g., users)
	var exists bool
	err = db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users')").Scan(&exists)
	assert.NoError(t, err)
	assert.True(t, exists, "users table should exist after migrations")
}

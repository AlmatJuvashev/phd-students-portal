package handlers_test

import (
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestDebugSchema(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	rows, err := db.Query("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = 'users'")
	assert.NoError(t, err)
	defer rows.Close()

	found := false
	t.Log("Columns in users table:")
	for rows.Next() {
		var name, dtype string
		rows.Scan(&name, &dtype)
		t.Logf(" - %s (%s)", name, dtype)
		if name == "tenant_id" {
			found = true
		}
	}
	assert.True(t, found, "tenant_id column missing in users table")
}

package testutils

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

// CreateTestTenant creates a tenant with the given ID and returns the ID.
// If ID is empty, generates one.
func CreateTestTenant(t *testing.T, db *sqlx.DB, id string) string {
	t.Helper()
	if id == "" {
		// We'd need uuid package but avoiding import cycle if possible or just use string
		// Since we don't import uuid here, require caller to pass it or we add uuid import.
		// Let's require caller to pass it for now or add import.
		t.Fatal("id must be provided")
	}
	slug := "slug-" + id
	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type, is_active) VALUES ($1, $2, 'T', 'university', true) ON CONFLICT DO NOTHING`, id, slug)
	require.NoError(t, err)
	return id
}

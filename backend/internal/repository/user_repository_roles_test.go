package repository

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLUserRepository_UserRoles_Integration(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()

	repo := NewSQLUserRepository(db)

	// 1. Create User
	userID := testutils.CreateTestUser(t, db, "multi_role_user", string(models.RoleStudent))

	// 2. Insert roles manually into user_roles table (Simulating Admin Action)
	// We use direct DB execution because repo doesn't have "AddUserRole" method yet (it wasn't required for Phase 1 Auth Service logic, but table existence is required)
	_, err := db.Exec("INSERT INTO user_roles (user_id, role) VALUES ($1, $2)", userID, "student")
	require.NoError(t, err, "failed to insert into user_roles - table might not exist")

	_, err = db.Exec("INSERT INTO user_roles (user_id, role) VALUES ($1, $2)", userID, "instructor")
	require.NoError(t, err)

	// 3. Test GetUserRoles
	roles, err := repo.GetUserRoles(context.Background(), userID)
	require.NoError(t, err)
	
	assert.Len(t, roles, 2)
	assert.Contains(t, roles, "student")
	assert.Contains(t, roles, "instructor")
}

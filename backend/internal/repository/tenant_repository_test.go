package repository

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLTenantRepository_CRUD(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()

	repo := NewSQLTenantRepository(db)

	domain := "test.local"
	appName := "Test App"
	tenant := &models.Tenant{
		Slug: "test-tenant",
		Name: "Test University",
		TenantType: "university",
		Domain: &domain,
		AppName: &appName,
	}

	// 1. Create
	id, err := repo.Create(context.Background(), tenant)
	require.NoError(t, err)
	assert.NotEmpty(t, id)

	// 2. GetByID
	fetched, err := repo.GetByID(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, tenant.Name, fetched.Name)
	assert.Equal(t, tenant.Slug, fetched.Slug)

	// 3. GetBySlug
	fetchedSlug, err := repo.GetBySlug(context.Background(), "test-tenant")
	require.NoError(t, err)
	assert.Equal(t, id, fetchedSlug.ID)

	// 4. Update
	updated, err := repo.Update(context.Background(), id, map[string]interface{}{
		"name": "Updated University",
	})
	require.NoError(t, err)
	assert.Equal(t, "Updated University", updated.Name)

	// 5. Exists
	exists, err := repo.Exists(context.Background(), id)
	require.NoError(t, err)
	assert.True(t, exists)

	// 6. Delete
	err = repo.Delete(context.Background(), id)
	require.NoError(t, err)
	
	// Should be inactive now
	fetchedDel, err := repo.GetByID(context.Background(), id)
	require.NoError(t, err)
	assert.False(t, fetchedDel.IsActive)
	
	// Exists still returns true because row exists? No, Exists uses SELECT 1.
	// But Delete sets is_active=false. Logic might differ.
	// user_repo Exists checks username. tenant_repo Exists checks ID.
	// Row still exists.
}

func TestSQLTenantRepository_Membership(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()

	repo := NewSQLTenantRepository(db)
	userRepo := NewSQLUserRepository(db)

	// Setup Tenant
	tID, err := repo.Create(context.Background(), &models.Tenant{Slug: "t1", Name: "T1", TenantType: "university"})
	require.NoError(t, err)

	// Setup User
	u := &models.User{Username: "u1", Email: "u1@test.com", Role: "student"}
	uID, err := userRepo.Create(context.Background(), u)
	require.NoError(t, err)

	// 1. Add User
	err = repo.AddUserToTenant(context.Background(), uID, tID, "student", true)
	require.NoError(t, err)

	// 2. Get Membership
	m, err := repo.GetUserMembership(context.Background(), uID, tID)
	require.NoError(t, err)
	assert.Equal(t, "student", m.Role)
	assert.True(t, m.IsPrimary)

	// 3. Get Role
	msgRole, err := repo.GetRole(context.Background(), uID, tID)
	require.NoError(t, err)
	assert.Equal(t, "student", msgRole)

	// 4. ListForUser
	list, err := repo.ListForUser(context.Background(), uID)
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "T1", list[0].TenantName)

	// 5. Remove
	err = repo.RemoveUser(context.Background(), uID, tID)
	require.NoError(t, err)
	
	m2, err := repo.GetUserMembership(context.Background(), uID, tID)
	require.NoError(t, err)
	assert.Nil(t, m2)
}

func TestSQLTenantRepository_Stats(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()

	repo := NewSQLTenantRepository(db)
	tID, _ := repo.Create(context.Background(), &models.Tenant{Slug: "t2", Name: "T2"})

	stats, err := repo.GetWithStats(context.Background(), tID)
	require.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, 0, stats.UserCount)

	// Add membership manually. Must create user first to satisfy FK.
	userRepo := NewSQLUserRepository(db)
	uid, err := userRepo.Create(context.Background(), &models.User{Username: "statuser", Email: "stat@example.com", Role: "student"})
	require.NoError(t, err)

	_, err = db.Exec("INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'student', true)", uid, tID)
	require.NoError(t, err)
	
	// ListAllWithStats
	all, err := repo.ListAllWithStats(context.Background())
	require.NoError(t, err)
	found := false
	for _, v := range all {
		if v.ID == tID {
			found = true
			assert.Equal(t, 1, v.UserCount)
		}
	}
	assert.True(t, found)
}

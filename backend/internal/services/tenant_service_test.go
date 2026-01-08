package services_test

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)



func TestTenantService_GetTenants(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLTenantRepository(db)
	svc := services.NewTenantService(repo)
	ctx := context.Background()

	// 1. GetTenantBySlug
	// Insert default tenant if not exists, or handle missing
	_, err := db.Exec(`INSERT INTO tenants (id, name, slug, is_active) 
		VALUES ('00000000-0000-0000-0000-000000000001', 'Default Tenant', 'default', true) 
		ON CONFLICT (id) DO UPDATE SET slug = 'default', is_active = true`)
	require.NoError(t, err)

	slug := "default"
	tenant, err := svc.GetTenantBySlug(ctx, slug)
	require.NoError(t, err)
	require.NotNil(t, tenant, "Tenant should not be nil")
	assert.Equal(t, "default", tenant.Slug)
	assert.True(t, tenant.IsActive)

	// 2. GetTenantByID
	id := "00000000-0000-0000-0000-000000000001"
	tenantById, err := svc.GetTenantByID(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, tenantById, "Tenant should not be nil")
	assert.Equal(t, id, tenantById.ID)

	// 3. ListAllWithStats
	tenants, err := svc.ListAllWithStats(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(tenants), 1)
}

func TestTenantService_UserMembership(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLTenantRepository(db)
	svc := services.NewTenantService(repo)
	ctx := context.Background()

	userID := "10000000-0000-0000-0000-000000000001"
	tenantID := "00000000-0000-0000-0000-000000000001"
	_, err := db.Exec(`INSERT INTO tenants (id, name, slug, is_active) 
		VALUES ($1, 'Default Tenant', 'default', true) 
		ON CONFLICT (id) DO NOTHING`, tenantID)
	require.NoError(t, err)

	// Create user
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'tenantuser', 'tenant@ex.com', 'Tenant', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	// 4. AddUserToTenant
	err = svc.AddUserToTenant(ctx, userID, tenantID, string(models.RoleStudent), true)
	require.NoError(t, err)

	// 5. GetUserMembershipInTenant
	membership, err := svc.GetUserMembershipInTenant(ctx, userID, tenantID)
	require.NoError(t, err)
	assert.Equal(t, string(models.RoleStudent), membership.Role)
	assert.True(t, membership.IsPrimary)

	// 6. GetUserTenants
	memberships, err := svc.GetUserTenants(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, memberships, 1)

	// 7. GetPrimaryTenant
	primary, err := svc.GetPrimaryTenant(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, tenantID, primary.ID)

	// 8. CanAccessTenant
	canAccess, err := svc.CanAccessTenant(ctx, userID, tenantID, false)
	require.NoError(t, err)
	assert.True(t, canAccess)
	
	canAccess, err = svc.CanAccessTenant(ctx, userID, "99999999-9999-9999-9999-999999999999", false)
	assert.NoError(t, err)
	assert.False(t, canAccess)

	// 9. GetUserRoleInTenant
	role, err := svc.GetUserRoleInTenant(ctx, userID, tenantID)
	require.NoError(t, err)
	assert.Equal(t, string(models.RoleStudent), role)

	// 10. RemoveUserFromTenant
	err = svc.RemoveUserFromTenant(ctx, userID, tenantID)
	require.NoError(t, err)

	membership, err = svc.GetUserMembershipInTenant(ctx, userID, tenantID)
	require.NoError(t, err)
	assert.Nil(t, membership) // Should not be found (nil result)
}

func TestTenantService_Management(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLTenantRepository(db)
	svc := services.NewTenantService(repo)
	ctx := context.Background()

	// 1. Create
	newTenant := &models.Tenant{
		Slug: "new-tenant",
		Name: "New Tenant",
		TenantType: "university",
		AppName: services.ToPtr("App"),
	}
	id, err := svc.Create(ctx, newTenant)
	require.NoError(t, err)
	assert.NotEmpty(t, id)

	// 2. Update
	updated, err := svc.Update(ctx, id, map[string]interface{}{"name": "Renamed Tenant"})
	require.NoError(t, err)
	assert.Equal(t, "Renamed Tenant", updated.Name)

	// 3. UpdateServices
	name, err := svc.UpdateServices(ctx, id, []string{"chat"})
	require.NoError(t, err)
	assert.Equal(t, "Renamed Tenant", name)

	fetched, err := svc.GetTenantByID(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, []string{"chat"}, []string(fetched.EnabledServices))

	// 4. UpdateLogo
	err = svc.UpdateLogo(ctx, id, "http://logo.png")
	require.NoError(t, err)

	// 5. Delete
	err = svc.Delete(ctx, id)
	require.NoError(t, err)

	fetched, err = svc.GetTenantByID(ctx, id)
	require.NoError(t, err)
	assert.False(t, fetched.IsActive)
}

func TestTenantService_ProtectedTenant(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLTenantRepository(db)
	svc := services.NewTenantService(repo)
	ctx := context.Background()

	platformID := "00000000-0000-0000-0000-000000000000"

	// 1. Try Update
	_, err := svc.Update(ctx, platformID, map[string]interface{}{"name": "Hacked"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reserved system resource")

	// 2. Try Delete
	err = svc.Delete(ctx, platformID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reserved system resource")
}

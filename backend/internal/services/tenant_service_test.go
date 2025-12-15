package services_test

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTenantService_GetTenantBySlug(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()

	svc := services.NewTenantService(db)
	ctx := context.Background()

	// 1. Setup Data
	// Insert a test tenant
	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, is_active) 
		VALUES ('11111111-1111-1111-1111-111111111111', 'test-slug', 'Test University', true)`)
	require.NoError(t, err)

	// 2. Test Cases
	tests := []struct {
		name      string
		slug      string
		wantErr   bool
		wantName  string
	}{
		{
			name:     "Existing Slug",
			slug:     "test-slug",
			wantErr:  false,
			wantName: "Test University",
		},
		{
			name:    "Non-Existent Slug",
			slug:    "fake-slug",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tenant, err := svc.GetTenantBySlug(ctx, tc.slug)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tenant)
				assert.Equal(t, tc.wantName, tenant.Name)
			}
		})
	}
}

func TestTenantService_CanAccessTenant(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()

	svc := services.NewTenantService(db)
	ctx := context.Background()

	// 1. Setup Data
	tenantID := "22222222-2222-2222-2222-222222222222"
	userID := "33333333-3333-3333-3333-333333333333"

	// Create Tenant
	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, is_active) VALUES ($1, 'access-test', 'Access Univ', true)`, tenantID)
	require.NoError(t, err)

	// Create User
	// Note: We need a user in the users table first because of FK constraints
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active, is_superadmin) 
		VALUES ($1, 'user1', 'user1@example.com', 'Test', 'User', 'student', 'hash', true, false)`, userID)
	require.NoError(t, err)

	// Create Membership
	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) 
		VALUES ($1, $2, 'student', true)`, userID, tenantID)
	require.NoError(t, err)


	// 2. Test Cases
	tests := []struct {
		name         string
		userID       string
		tenantID     string
		isSuperadmin bool
		wantAccess   bool
	}{
		{
			name:         "Member can access",
			userID:       userID,
			tenantID:     tenantID,
			isSuperadmin: false,
			wantAccess:   true,
		},
		{
			name:         "Superadmin can access anything",
			userID:       "99999999-9999-9999-9999-999999999999", // Random ID
			tenantID:     tenantID,
			isSuperadmin: true,
			wantAccess:   true,
		},
		{
			name:         "Non-member cannot access",
			userID:       "99999999-9999-9999-9999-999999999999", // Random ID
			tenantID:     tenantID,
			isSuperadmin: false,
			wantAccess:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			hasAccess, err := svc.CanAccessTenant(ctx, tc.userID, tc.tenantID, tc.isSuperadmin)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantAccess, hasAccess)
		})
	}
}

func TestTenantService_AddUserToTenant(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()

	svc := services.NewTenantService(db)
	ctx := context.Background()

	// Setup
	tenantID := "44444444-4444-4444-4444-444444444444"
	userID := "55555555-5555-5555-5555-555555555555"

	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, is_active) VALUES ($1, 'add-test', 'Add Univ', true)`, tenantID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'user2', 'user2@example.com', 'Test', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	// Test Adding User
	err = svc.AddUserToTenant(ctx, userID, tenantID, models.RoleStudent, true)
	assert.NoError(t, err)

	// Verify Membership
	role, err := svc.GetUserRoleInTenant(ctx, userID, tenantID)
	assert.NoError(t, err)
	assert.Equal(t, models.RoleStudent, role)

	// Test Update Role (Idempotency)
	err = svc.AddUserToTenant(ctx, userID, tenantID, models.RoleAdvisor, true)
	assert.NoError(t, err)

	role, err = svc.GetUserRoleInTenant(ctx, userID, tenantID)
	assert.NoError(t, err)
	assert.Equal(t, models.RoleAdvisor, role)
}

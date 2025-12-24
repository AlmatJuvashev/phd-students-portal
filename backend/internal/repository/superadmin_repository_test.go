package repository

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/auth"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSuperAdminRepository_AdminManagement(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()

	repo := NewSQLSuperAdminRepository(db)

	t.Run("CreateAdmin", func(t *testing.T) {
		hash, _ := auth.HashPassword("testpass123")
		params := models.CreateAdminParams{
			Username:     "newadmin",
			Email:        "newadmin@test.com",
			FirstName:    "New",
			LastName:     "Admin",
			Role:         "admin",
			PasswordHash: hash,
		}

		adminID, err := repo.CreateAdmin(context.Background(), params)
		require.NoError(t, err)
		assert.NotEmpty(t, adminID)

		// Verify admin was created
		admin, _, err := repo.GetAdmin(context.Background(), adminID)
		require.NoError(t, err)
		assert.Equal(t, "newadmin", admin.Username)
		assert.Equal(t, "newadmin@test.com", admin.Email)
		assert.True(t, admin.IsActive)
	})

	t.Run("UpdateAdmin", func(t *testing.T) {
		// Create admin first
		hash, _ := auth.HashPassword("testpass123")
		createParams := models.CreateAdminParams{
			Username:     "updateadmin",
			Email:        "update@test.com",
			FirstName:    "Update",
			LastName:     "Test",
			Role:         "admin",
			PasswordHash: hash,
		}
		adminID, err := repo.CreateAdmin(context.Background(), createParams)
		require.NoError(t, err)

		// Update admin
		updateParams := models.UpdateAdminParams{
			FirstName: strPtr("Updated"),
			LastName:  strPtr("Name"),
			Email:     strPtr("updated@test.com"),
		}
		_, err = repo.UpdateAdmin(context.Background(), adminID, updateParams)
		require.NoError(t, err)

		// Verify update
		admin, _, err := repo.GetAdmin(context.Background(), adminID)
		require.NoError(t, err)
		assert.Equal(t, "Updated", admin.FirstName)
		assert.Equal(t, "Name", admin.LastName)
		assert.Equal(t, "updated@test.com", admin.Email)
	})

	t.Run("ResetPassword", func(t *testing.T) {
		// Create admin
		hash, _ := auth.HashPassword("oldpass123")
		createParams := models.CreateAdminParams{
			Username:     "resetadmin",
			Email:        "reset@test.com",
			FirstName:    "Reset",
			LastName:     "Test",
			Role:         "admin",
			PasswordHash: hash,
		}
		adminID, err := repo.CreateAdmin(context.Background(), createParams)
		require.NoError(t, err)

		// Reset password
		newHash, _ := auth.HashPassword("newpass123")
		_, err = repo.ResetPassword(context.Background(), adminID, newHash)
		require.NoError(t, err)

		// Verify password was changed (we can't check the hash directly, but we can verify no error)
		admin, _, err := repo.GetAdmin(context.Background(), adminID)
		require.NoError(t, err)
		assert.NotEmpty(t, admin.ID)
	})

	t.Run("DeleteAdmin", func(t *testing.T) {
		// Create admin
		hash, _ := auth.HashPassword("testpass123")
		createParams := models.CreateAdminParams{
			Username:     "deleteadmin",
			Email:        "delete@test.com",
			FirstName:    "Delete",
			LastName:     "Test",
			Role:         "admin",
			PasswordHash: hash,
		}
		adminID, err := repo.CreateAdmin(context.Background(), createParams)
		require.NoError(t, err)

		// Delete admin
		_, err = repo.DeleteAdmin(context.Background(), adminID)
		require.NoError(t, err)

		// Verify admin is inactive
		admin, _, err := repo.GetAdmin(context.Background(), adminID)
		require.NoError(t, err)
		assert.False(t, admin.IsActive, "Admin should be marked as inactive after deletion")
	})

	t.Run("ListAdmins", func(t *testing.T) {
		// Create multiple admins
		hash, _ := auth.HashPassword("testpass123")
		for i := 1; i <= 3; i++ {
			params := models.CreateAdminParams{
				Username:     "listadmin" + string(rune('0'+i)),
				Email:        "listadmin" + string(rune('0'+i)) + "@test.com",
				FirstName:    "List",
				LastName:     "Admin" + string(rune('0'+i)),
				Role:         "admin",
				PasswordHash: hash,
			}
			_, err := repo.CreateAdmin(context.Background(), params)
			require.NoError(t, err)
		}

		// List all admins
		admins, err := repo.ListAdmins(context.Background(), "")
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(admins), 3, "Should have at least 3 admins")
	})
}

func TestSuperAdminRepository_Settings(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()

	repo := NewSQLSuperAdminRepository(db)
	userRepo := NewSQLUserRepository(db)

	// Create a test user for UpdatedBy field
	testUserID, _ := userRepo.Create(context.Background(), &models.User{
		Username: "settingsuser",
		Email:    "settings@test.com",
		Role:     "admin",
	})

	t.Run("UpdateSetting", func(t *testing.T) {
		params := models.UpdateSettingParams{
			Value:       strPtr("test_value"),
			Description: strPtr("Test setting"),
			Category:    strPtr("test"),
			UpdatedBy:   testUserID,
		}

		setting, err := repo.UpdateSetting(context.Background(), "test_key", params)
		require.NoError(t, err)
		assert.Equal(t, "test_key", setting.Key)
		// Value is stored as JSON, so we need to unmarshal to compare
		var value string
		json.Unmarshal(setting.Value, &value)
		assert.Equal(t, "test_value", value)
		assert.Equal(t, "test", setting.Category)
	})

	t.Run("GetSetting", func(t *testing.T) {
		// Create setting first
		params := models.UpdateSettingParams{
			Value:       strPtr("get_value"),
			Description: strPtr("Get test"),
			Category:    strPtr("test"),
			UpdatedBy:   testUserID,
		}
		_, err := repo.UpdateSetting(context.Background(), "get_key", params)
		require.NoError(t, err)

		// Get setting
		setting, err := repo.GetSetting(context.Background(), "get_key")
		require.NoError(t, err)
		assert.Equal(t, "get_key", setting.Key)
		// Value is stored as JSON
		var value string
		json.Unmarshal(setting.Value, &value)
		assert.Equal(t, "get_value", value)
	})

	t.Run("ListSettings", func(t *testing.T) {
		// Create multiple settings in same category
		for i := 1; i <= 3; i++ {
			params := models.UpdateSettingParams{
				Value:       strPtr("value" + string(rune('0'+i))),
				Description: strPtr("Setting " + string(rune('0'+i))),
				Category:    strPtr("list_test"),
				UpdatedBy:   testUserID,
			}
			_, err := repo.UpdateSetting(context.Background(), "list_key"+string(rune('0'+i)), params)
			require.NoError(t, err)
		}

		// List settings by category
		settings, err := repo.ListSettings(context.Background(), "list_test")
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(settings), 3, "Should have at least 3 settings in list_test category")
	})

	t.Run("DeleteSetting", func(t *testing.T) {
		// Create setting
		params := models.UpdateSettingParams{
			Value:       strPtr("delete_value"),
			Description: strPtr("Delete test"),
			Category:    strPtr("test"),
			UpdatedBy:   testUserID,
		}
		_, err := repo.UpdateSetting(context.Background(), "delete_key", params)
		require.NoError(t, err)

		// Delete setting
		err = repo.DeleteSetting(context.Background(), "delete_key")
		require.NoError(t, err)

		// Verify setting is deleted (GetSetting returns nil, nil for not found)
		setting, err := repo.GetSetting(context.Background(), "delete_key")
		assert.NoError(t, err, "GetSetting should not error for missing key")
		assert.Nil(t, setting, "Setting should be nil after deletion")
	})

	t.Run("GetCategories", func(t *testing.T) {
		// Create settings in different categories
		categories := []string{"cat1", "cat2", "cat3"}
		for _, cat := range categories {
			params := models.UpdateSettingParams{
				Value:       strPtr("value"),
				Description: strPtr("Test"),
				Category:    strPtr(cat),
				UpdatedBy:   testUserID,
			}
			_, err := repo.UpdateSetting(context.Background(), "key_"+cat, params)
			require.NoError(t, err)
		}

		// Get all categories
		cats, err := repo.GetCategories(context.Background())
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(cats), 3, "Should have at least 3 categories")
	})
}

func TestSuperAdminRepository_ActivityLogs(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()

	repo := NewSQLSuperAdminRepository(db)
	userRepo := NewSQLUserRepository(db)
	tenantRepo := NewSQLTenantRepository(db)

	// Setup test data
	tID, _ := tenantRepo.Create(context.Background(), &models.Tenant{Slug: "log-test", Name: "Log Test"})
	uID, _ := userRepo.Create(context.Background(), &models.User{Username: "loguser", Email: "log@test.com", Role: "admin"})

	t.Run("LogActivity", func(t *testing.T) {
		params := models.ActivityLogParams{
			TenantID:   &tID,
			UserID:     &uID,
			Action:     "create",
			EntityType: "user",
			EntityID:   uID,
			Description: "Test log entry",
			IPAddress:  "127.0.0.1",
			Metadata:   map[string]interface{}{"test": "data", "status": "active"},
		}

		err := repo.LogActivity(context.Background(), params)
		require.NoError(t, err)
	})

	t.Run("ListLogs", func(t *testing.T) {
		// Create multiple log entries
		for i := 1; i <= 5; i++ {
			params := models.ActivityLogParams{
				TenantID:   &tID,
				UserID:     &uID,
				Action:     "test_action_" + strconv.Itoa(i),
				EntityType: "test_entity_" + strconv.Itoa(i),
				EntityID:   uID,
				IPAddress:  "127.0.0.1",
				Metadata:   map[string]interface{}{"index": i},
			}
			err := repo.LogActivity(context.Background(), params)
			require.NoError(t, err)
		}

		// List logs with filter
		filter := LogFilter{
			TenantID: tID,
		}
		pagination := Pagination{
			Limit:  10,
			Offset: 0,
		}

		logs, total, err := repo.ListLogs(context.Background(), filter, pagination)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, total, 5, "Should have at least 5 log entries")
		assert.GreaterOrEqual(t, len(logs), 5, "Should return at least 5 logs")
	})

	t.Run("GetLogStats", func(t *testing.T) {
		stats, err := repo.GetLogStats(context.Background())
		require.NoError(t, err)
		assert.NotNil(t, stats)
		assert.GreaterOrEqual(t, stats.TotalLogs, 0)
	})

	t.Run("GetActions", func(t *testing.T) {
		actions, err := repo.GetActions(context.Background())
		require.NoError(t, err)
		assert.NotEmpty(t, actions, "Should have at least one action type")
		assert.Contains(t, actions, "create")
	})

	t.Run("GetEntityTypes", func(t *testing.T) {
		entityTypes, err := repo.GetEntityTypes(context.Background())
		require.NoError(t, err)
		assert.NotEmpty(t, entityTypes, "Should have at least one entity type")
		assert.Contains(t, entityTypes, "user")
	})
}

// Helper function to create string pointers
func strPtr(s string) *string {
	return &s
}

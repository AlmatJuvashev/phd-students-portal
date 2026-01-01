package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSQLSuperAdminRepository_ListAdmins_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLSuperAdminRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "username", "email", "first_name", "last_name", 
			"role", "is_active", "is_superadmin", 
			"tenant_id", "tenant_name", "tenant_slug",
			"created_at", "updated_at",
		}).AddRow(
			"adm-1", "admin1", "admin@ex.com", "Admin", "One",
			"admin", true, false,
			"t-1", "Tenant 1", "t1",
			time.Now(), time.Now(),
		)

		// Regex to match the complex query. matches any char including newlines.
		mock.ExpectQuery(`SELECT [\s\S]* FROM users u [\s\S]* WHERE [\s\S]* ORDER BY u.username`).
			WillReturnRows(rows)

		admins, err := repo.ListAdmins(context.Background(), "")
		assert.NoError(t, err)
		if assert.Len(t, admins, 1) {
			assert.Equal(t, "admin1", admins[0].Username)
		}
	})

	t.Run("WithTenantFilter", func(t *testing.T) {
		mock.ExpectQuery(`SELECT [\s\S]* FROM users u [\s\S]* WHERE [\s\S]* AND utm.tenant_id = \$1 ORDER BY u.username`).
			WithArgs("t-1").
			WillReturnRows(sqlmock.NewRows([]string{}))

		admins, err := repo.ListAdmins(context.Background(), "t-1")
		assert.NoError(t, err)
		assert.Len(t, admins, 0)
	})
}

func TestSQLSuperAdminRepository_CreateAdmin_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLSuperAdminRepository(sqlx.NewDb(db, "sqlmock"))

	params := models.CreateAdminParams{
		Username:     "newadmin",
		Email:        "new@ex.com",
		PasswordHash: "hash",
		FirstName:    "New",
		LastName:     "Admin",
		Role:         "admin",
		IsSuperadmin: false,
		TenantIDs:    []string{"t-1", "t-2"},
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		
		// Insert User
		mock.ExpectQuery(`INSERT INTO users`).
			WithArgs(params.Username, params.Email, params.PasswordHash, params.FirstName, params.LastName, params.Role, params.IsSuperadmin).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("u-1"))

		// Insert Memberships
		// First tenant (primary)
		mock.ExpectExec(`INSERT INTO user_tenant_memberships`).
			WithArgs("u-1", "t-1", params.Role, true).
			WillReturnResult(sqlmock.NewResult(1, 1))
		
		// Second tenant (not primary)
		mock.ExpectExec(`INSERT INTO user_tenant_memberships`).
			WithArgs("u-1", "t-2", params.Role, false).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		id, err := repo.CreateAdmin(context.Background(), params)
		assert.NoError(t, err)
		assert.Equal(t, "u-1", id)
	})

	t.Run("RollbackOnUserFail", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO users`).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback() // Implicitly handled by BeginTxx/defer? No, defer tx.Rollback() does it if not commited.

		_, err := repo.CreateAdmin(context.Background(), params)
		assert.Error(t, err)
	})
}

func TestSQLSuperAdminRepository_UpdateAdmin_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLSuperAdminRepository(sqlx.NewDb(db, "sqlmock"))

	id := "u-1"
	role := "admin"
	params := models.UpdateAdminParams{
		FirstName: func() *string { s := "Updated"; return &s }(),
		Role:      &role,
		TenantIDs: []string{"t-3"},
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()

		// Update User
		mock.ExpectQuery(`UPDATE users SET`).
			WithArgs(id, nil, params.FirstName, nil, params.Role, nil, nil). // Matches arg order in query
			WillReturnRows(sqlmock.NewRows([]string{"username"}).AddRow("admin1"))

		// Update Tenants
		mock.ExpectExec(`DELETE FROM user_tenant_memberships WHERE user_id = \$1`).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 2))

		mock.ExpectExec(`INSERT INTO user_tenant_memberships`).
			WithArgs(id, "t-3", role, true).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		username, err := repo.UpdateAdmin(context.Background(), id, params)
		assert.NoError(t, err)
		assert.Equal(t, "admin1", username)
	})
}

func TestSQLSuperAdminRepository_LogActivity_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLSuperAdminRepository(sqlx.NewDb(db, "sqlmock"))

	S := func(s string) *string { return &s }

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(`INSERT INTO activity_logs`).
			WithArgs("u-1", "t-1", "login", "user", "u-1", "logged in", "127.0.0.1", "mozilla", []byte(nil)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.LogActivity(context.Background(), models.ActivityLogParams{
			UserID: S("u-1"), TenantID: S("t-1"), Action: "login", EntityType: "user", EntityID: "u-1", Description: "logged in", IPAddress: "127.0.0.1", UserAgent: "mozilla",
		})
		assert.NoError(t, err)
	})
}

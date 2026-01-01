package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestSQLTenantRepository_GetByID_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLTenantRepository(sqlxDB)

	tenantID := "tenant-1"

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "slug", "name", "tenant_type", "is_active"}).
			AddRow(tenantID, "test-slug", "Test Tenant", "university", true)

		mock.ExpectQuery(`SELECT (.+) FROM tenants WHERE id = \$1`).
			WithArgs(tenantID).
			WillReturnRows(rows)

		tenant, err := repo.GetByID(context.Background(), tenantID)

		assert.NoError(t, err)
		assert.NotNil(t, tenant)
		assert.Equal(t, tenantID, tenant.ID)
		assert.Equal(t, "test-slug", tenant.Slug)
	})
}

func TestSQLTenantRepository_ListAllWithStats_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLTenantRepository(sqlxDB)

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows([]string{"id", "slug", "name", "tenant_type", "domain", "logo_url", "app_name", "primary_color", "secondary_color", "enabled_services", "is_active", "created_at", "updated_at", "user_count", "admin_count"}).
			AddRow("t-1", "s1", "N1", "U", "d1", "l1", "a1", "#000", "#111", nil, true, now, now, 10, 2)

		mock.ExpectQuery(`SELECT t.id, t.slug, t.name, (.+) FROM tenants t LEFT JOIN`).
			WillReturnRows(rows)

		stats, err := repo.ListAllWithStats(context.Background())

		assert.NoError(t, err)
		if assert.Len(t, stats, 1) {
			assert.Equal(t, 10, stats[0].UserCount)
		}
	})
}

func TestSQLTenantRepository_Update_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLTenantRepository(sqlxDB)

	id := "t-1"
	updates := map[string]interface{}{
		"name": "Updated Name",
	}

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		mock.ExpectQuery(`UPDATE tenants SET`).
			WithArgs("Updated Name", id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "slug", "name", "tenant_type", "domain", "logo_url", "app_name", "primary_color", "secondary_color", "enabled_services", "is_active", "created_at", "updated_at"}).
				AddRow(id, "slug", "Updated Name", "U", "", "", "", "", "", nil, true, now, now))

		tenant, err := repo.Update(context.Background(), id, updates)
		assert.NoError(t, err)
		assert.NotNil(t, tenant)
		assert.Equal(t, "Updated Name", tenant.Name)
	})
}

func TestSQLTenantRepository_UpdateServices_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLTenantRepository(sqlx.NewDb(db, "sqlmock"))

	id := "t-1"
	services := []string{"chat", "calendar"}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(`UPDATE tenants SET enabled_services = \$2, updated_at = now\(\) WHERE id = \$1 RETURNING name`).
			WithArgs(id, pq.Array(services)).
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Tenant One"))

		name, err := repo.UpdateServices(context.Background(), id, services)
		assert.NoError(t, err)
		assert.Equal(t, "Tenant One", name)
	})
}

func TestSQLTenantRepository_UpdateLogo_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLTenantRepository(sqlx.NewDb(db, "sqlmock"))

	id := "t-1"
	url := "http://example.com/logo.png"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(`UPDATE tenants SET logo_url = \$2, updated_at = now\(\) WHERE id = \$1`).
			WithArgs(id, url).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateLogo(context.Background(), id, url)
		assert.NoError(t, err)
	})
}

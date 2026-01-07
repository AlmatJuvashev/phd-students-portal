package repository

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
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

func TestSQLTenantRepository_GetBySlug_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLTenantRepository(sqlx.NewDb(db, "sqlmock"))

	slug := "test-slug"
	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "slug", "name"}).AddRow("t1", slug, "Name")
		mock.ExpectQuery(`SELECT (.+) FROM tenants WHERE slug = \$1`).WithArgs(slug).WillReturnRows(rows)
		
		res, err := repo.GetBySlug(context.Background(), slug)
		assert.NoError(t, err)
		assert.Equal(t, slug, res.Slug)
	})
}

func TestSQLTenantRepository_Create_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLTenantRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		mock.ExpectQuery(`INSERT INTO tenants`).
			WithArgs("slug", "name", "type", "dom", "app", "#000", "#111").
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow("new-id", now, now))

		id, err := repo.Create(context.Background(), &models.Tenant{
			Slug: "slug", Name: "name", TenantType: "type", Domain: toPtr("dom"), AppName: toPtr("app"), PrimaryColor: toPtr("#000"), SecondaryColor: toPtr("#111"),
		})
		assert.NoError(t, err)
		assert.Equal(t, "new-id", id)
	})
}

func TestSQLTenantRepository_Delete_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLTenantRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(`UPDATE tenants SET is_active = false`).
			WithArgs("t1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Delete(context.Background(), "t1")
		assert.NoError(t, err)
	})
}

func TestSQLTenantRepository_Membership_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLTenantRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("AddUser", func(t *testing.T) {
		mock.ExpectExec(`INSERT INTO user_tenant_memberships`).
			WithArgs("u1", "t1", "student", true).
			WillReturnResult(sqlmock.NewResult(1, 1))
		err := repo.AddUserToTenant(context.Background(), "u1", "t1", "student", true)
		assert.NoError(t, err)
	})

	t.Run("RemoveUser", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM user_tenant_memberships`).
			WithArgs("u1", "t1").
			WillReturnResult(sqlmock.NewResult(1, 1))
		err := repo.RemoveUser(context.Background(), "u1", "t1")
		assert.NoError(t, err)
	})

	t.Run("GetRole", func(t *testing.T) {
		mock.ExpectQuery(`SELECT role FROM user_tenant_memberships`).
			WithArgs("u1", "t1").
			WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("admin"))
		role, err := repo.GetRole(context.Background(), "u1", "t1")
		assert.NoError(t, err)
		assert.Equal(t, "admin", role)
	})
}

func TestSQLTenantRepository_Exists_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLTenantRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("True", func(t *testing.T) {
		mock.ExpectQuery(`SELECT EXISTS`).WithArgs("t1").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
		exists, err := repo.Exists(context.Background(), "t1")
		assert.NoError(t, err)
		assert.True(t, exists)
	})
}

func TestSQLTenantRepository_ListForUser_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLTenantRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"tenant_id", "tenant_name", "tenant_slug", "role", "is_primary"}).
			AddRow("t1", "Home", "home", "student", true)
		mock.ExpectQuery(`SELECT utm.tenant_id`).WithArgs("u1").WillReturnRows(rows)
		
		list, err := repo.ListForUser(context.Background(), "u1")
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})
}

package repository

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSQLContactRepository_ListPublic(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLContactRepository(sqlxDB)

	tenantID := "t1"

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "tenant_id", "name", "title", "email", "phone", "sort_order", "is_active", "created_at", "updated_at"}).
			AddRow("c1", tenantID, []byte(`{"en":"John"}`), []byte(`{"en":"Admin"}`), "j@ex.com", "123", 1, true, "2023-01-01T00:00:00Z", "2023-01-01T00:00:00Z")

		mock.ExpectQuery("SELECT (.+) FROM contacts WHERE tenant_id = \\$1 AND is_active = true").
			WithArgs(tenantID).
			WillReturnRows(rows)

		contacts, err := repo.ListPublic(context.Background(), tenantID)
		assert.NoError(t, err)
		assert.Len(t, contacts, 1)
		assert.Equal(t, "John", contacts[0].Name["en"])
	})
}

func TestSQLContactRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLContactRepository(sqlxDB)

	tenantID := "t1"
	contact := models.Contact{
		Name:      models.LocalizedMap{"en": "John"},
		Title:     models.LocalizedMap{"en": "Admin"},
		SortOrder: 1,
		IsActive:  true,
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO contacts").
			WithArgs(tenantID, contact.Name, contact.Title, nil, nil, contact.SortOrder, contact.IsActive).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("c1"))

		id, err := repo.Create(context.Background(), tenantID, contact)
		assert.NoError(t, err)
		assert.Equal(t, "c1", id)
	})
}

func TestSQLContactRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLContactRepository(sqlxDB)

	tenantID := "t1"
	id := "c1"

	t.Run("Success", func(t *testing.T) {
		updates := map[string]interface{}{
			"is_active": false,
		}

		// sort_order is not in updates, name is not in updates.
		// Query will have updated_at and is_active.
		mock.ExpectExec("UPDATE contacts SET updated_at = now\\(\\), is_active = \\$1 WHERE id = \\$2 AND tenant_id = \\$3").
			WithArgs(false, id, tenantID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Update(context.Background(), tenantID, id, updates)
		assert.NoError(t, err)
	})

	t.Run("With Email Null", func(t *testing.T) {
		updates := map[string]interface{}{
			"email": "",
		}

		mock.ExpectExec("UPDATE contacts SET updated_at = now\\(\\), email = \\$1 WHERE id = \\$2 AND tenant_id = \\$3").
			WithArgs(nil, id, tenantID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Update(context.Background(), tenantID, id, updates)
		assert.NoError(t, err)
	})
}

func TestSQLContactRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLContactRepository(sqlxDB)

	tenantID := "t1"
	id := "c1"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("UPDATE contacts SET is_active = false").
			WithArgs(id, tenantID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(context.Background(), tenantID, id)
		assert.NoError(t, err)
	})
}

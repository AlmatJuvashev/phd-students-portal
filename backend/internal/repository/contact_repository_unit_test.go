package repository

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func newTestContactRepo(t *testing.T) (*SQLContactRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLContactRepository(sqlxDB)
	return repo, mock, func() { db.Close() }
}

func TestSQLContactRepository_ListPublic(t *testing.T) {
	repo, mock, teardown := newTestContactRepo(t)
	defer teardown()
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "tenant_id", "name", "title", "email", "phone", "sort_order", "is_active", "created_at", "updated_at"}).
		AddRow("c1", "t1", []byte(`{}`), []byte(`{}`), "test@example.com", "123", 1, true, time.Now().Format(time.RFC3339), time.Now().Format(time.RFC3339))

	mock.ExpectQuery(`SELECT id, tenant_id, .+ FROM contacts WHERE tenant_id = \$1 AND is_active = true ORDER BY sort_order, created_at`).
		WithArgs("t1").
		WillReturnRows(rows)

	list, err := repo.ListPublic(ctx, "t1")
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "c1", list[0].ID)
}

func TestSQLContactRepository_ListAdmin(t *testing.T) {
	repo, mock, teardown := newTestContactRepo(t)
	defer teardown()
	ctx := context.Background()

	t.Run("IncludesInactive", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "tenant_id", "name", "title", "email", "phone", "sort_order", "is_active", "created_at", "updated_at"}).
			AddRow("c1", "t1", []byte(`{}`), []byte(`{}`), "test@example.com", "123", 1, false, time.Now().Format(time.RFC3339), time.Now().Format(time.RFC3339))

		mock.ExpectQuery(`SELECT id, tenant_id, .+ FROM contacts WHERE tenant_id = \$1 ORDER BY sort_order, created_at`).
			WithArgs("t1").
			WillReturnRows(rows)

		list, err := repo.ListAdmin(ctx, "t1", true)
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})

	t.Run("ActiveOnly", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, tenant_id, .+ FROM contacts WHERE tenant_id = \$1 AND is_active = true ORDER BY sort_order, created_at`).
			WithArgs("t1").
			WillReturnRows(sqlmock.NewRows([]string{}))

		list, err := repo.ListAdmin(ctx, "t1", false)
		assert.NoError(t, err)
		assert.Len(t, list, 0)
	})
}

func TestSQLContactRepository_Create(t *testing.T) {
	repo, mock, teardown := newTestContactRepo(t)
	defer teardown()
	ctx := context.Background()

	email := "new@example.com"
	contact := models.Contact{
		Name:      models.LocalizedMap{"en": "New Contact"},
		Title:     models.LocalizedMap{"en": "Manager"},
		Email:     &email,
		SortOrder: 2,
		IsActive:  true,
	}

	mock.ExpectQuery(`INSERT INTO contacts`).
		WithArgs("t1", contact.Name, contact.Title, "new@example.com", nil, 2, true).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("new-id"))

	id, err := repo.Create(ctx, "t1", contact)
	assert.NoError(t, err)
	assert.Equal(t, "new-id", id)
}

func TestSQLContactRepository_Update(t *testing.T) {
	repo, mock, teardown := newTestContactRepo(t)
	defer teardown()
	ctx := context.Background()

	updates := map[string]interface{}{
		"email": "updated@example.com",
		"phone": nil,
	}

	// Because map iteration order is random, we need to be careful or rely on implementation sorting keys
	// The implementation sorts keys, so "email" comes before "phone"
	// Also "updated_at = now()" is added first in implementation.
	mock.ExpectExec(`UPDATE contacts SET updated_at = now\(\), email = \$1, phone = \$2 WHERE id = \$3 AND tenant_id = \$4`).
		WithArgs("updated@example.com", nil, "c1", "t1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Update(ctx, "t1", "c1", updates)
	assert.NoError(t, err)
}

func TestSQLContactRepository_Delete(t *testing.T) {
	repo, mock, teardown := newTestContactRepo(t)
	defer teardown()
	ctx := context.Background()

	mock.ExpectExec(`UPDATE contacts SET is_active = false, updated_at = now\(\) WHERE id = \$1 AND tenant_id = \$2`).
		WithArgs("c1", "t1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Delete(ctx, "t1", "c1")
	assert.NoError(t, err)
}

func TestHelpers(t *testing.T) {
	s := "test"
	assert.Equal(t, "test", contactNullableString(s))
	assert.Nil(t, contactNullableString(""))
	assert.Nil(t, contactNullableString("   "))

	assert.Equal(t, "test", contactNullablePtr(&s))
	assert.Nil(t, contactNullablePtr(nil))
}

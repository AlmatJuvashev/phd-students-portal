package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSQLSearchRepository_SearchUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLSearchRepository(sqlxDB)

	query := "test"
	limit := 10

	t.Run("Unauthorized Role", func(t *testing.T) {
		results, err := repo.SearchUsers(context.Background(), query, "student", "u1", limit)
		assert.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("Admin Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "role"}).
			AddRow("u1", "First", "Last", "u1@ex.com", "student")

		mock.ExpectQuery("SELECT id, first_name, last_name, email, role FROM users WHERE").
			WithArgs("%"+query+"%", limit).
			WillReturnRows(rows)

		results, err := repo.SearchUsers(context.Background(), query, "admin", "admin1", limit)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "student", results[0].Type)
		assert.Equal(t, "First Last", results[0].Title)
	})

	t.Run("Advisor Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "role"}).
			AddRow("u1", "First", "Last", "u1@ex.com", "student")

		mock.ExpectQuery("SELECT id, first_name, last_name, email, role FROM users WHERE").
			WithArgs("%"+query+"%", "adv1", limit).
			WillReturnRows(rows)

		results, err := repo.SearchUsers(context.Background(), query, "advisor", "adv1", limit)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("db error"))
		results, err := repo.SearchUsers(context.Background(), query, "superadmin", "sa1", limit)
		assert.Error(t, err)
		assert.Nil(t, results)
	})
}

func TestSQLSearchRepository_SearchDocuments(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLSearchRepository(sqlxDB)

	query := "doc"
	limit := 10

	t.Run("Admin Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "filename", "first_name", "last_name", "node_id"}).
			AddRow("d1", "file.pdf", "John", "Doe", "node1")

		mock.ExpectQuery("SELECT a.id, a.filename, u.first_name, u.last_name, ni.node_id FROM node_instance_slot_attachments a").
			WithArgs("%"+query+"%", limit).
			WillReturnRows(rows)

		results, err := repo.SearchDocuments(context.Background(), query, "admin", "admin1", limit)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "document", results[0].Type)
		assert.Equal(t, "file.pdf", results[0].Title)
	})

	t.Run("Student Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "filename", "first_name", "last_name", "node_id"}).
			AddRow("d1", "file.pdf", "John", "Doe", "node1")

		mock.ExpectQuery("SELECT (.+) FROM node_instance_slot_attachments a (.+) ni.user_id = \\$2").
			WithArgs("%"+query+"%", "s1", limit).
			WillReturnRows(rows)

		results, err := repo.SearchDocuments(context.Background(), query, "student", "s1", limit)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
	})

	t.Run("Advisor Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "filename", "first_name", "last_name", "node_id"}).
			AddRow("d1", "file.pdf", "John", "Doe", "node1")

		mock.ExpectQuery("SELECT (.+) FROM node_instance_slot_attachments a (.+) advisor_id = \\$2").
			WithArgs("%"+query+"%", "adv1", limit).
			WillReturnRows(rows)

		results, err := repo.SearchDocuments(context.Background(), query, "advisor", "adv1", limit)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
	})
}

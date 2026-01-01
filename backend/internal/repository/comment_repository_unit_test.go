package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSQLCommentRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLCommentRepository(sqlxDB)

	comment := models.Comment{
		TenantID:   "t1",
		DocumentID: "d1",
		UserID:     "u1",
		Content:    "test comment",
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO comments").
			WithArgs(comment.TenantID, comment.DocumentID, comment.UserID, comment.Content, nil).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("c1"))

		id, err := repo.Create(context.Background(), comment)
		assert.NoError(t, err)
		assert.Equal(t, "c1", id)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO comments").
			WillReturnError(fmt.Errorf("db error"))

		id, err := repo.Create(context.Background(), comment)
		assert.Error(t, err)
		assert.Empty(t, id)
	})
}

func TestSQLCommentRepository_GetByDocumentID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLCommentRepository(sqlxDB)

	tenantID := "t1"
	docID := "d1"

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "tenant_id", "document_id", "user_id", "content", "parent_id", "created_at"}).
			AddRow("c1", tenantID, docID, "u1", "comment 1", nil, "2023-01-01 10:00:00").
			AddRow("c2", tenantID, docID, "u2", "comment 2", "c1", "2023-01-01 11:00:00")

		mock.ExpectQuery("SELECT (.+) FROM comments WHERE document_id = \\$1 AND tenant_id = \\$2").
			WithArgs(docID, tenantID).
			WillReturnRows(rows)

		comments, err := repo.GetByDocumentID(context.Background(), tenantID, docID)
		assert.NoError(t, err)
		assert.Len(t, comments, 2)
		assert.Equal(t, "c1", comments[0].ID)
		assert.Equal(t, "c2", comments[1].ID)
		assert.Equal(t, "c1", *comments[1].ParentID)
	})

	t.Run("Empty", func(t *testing.T) {
		mock.ExpectQuery("SELECT").
			WithArgs(docID, tenantID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		comments, err := repo.GetByDocumentID(context.Background(), tenantID, docID)
		assert.NoError(t, err)
		assert.NotNil(t, comments)
		assert.Empty(t, comments)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT").
			WillReturnError(fmt.Errorf("db error"))

		comments, err := repo.GetByDocumentID(context.Background(), tenantID, docID)
		assert.Error(t, err)
		assert.Nil(t, comments)
	})
}

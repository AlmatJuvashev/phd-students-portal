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

func TestSQLDocumentRepository_Create_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLDocumentRepository(sqlxDB)

	doc := &models.Document{
		TenantID: "tenant-1",
		UserID:   "user-1",
		Kind:     "passport",
		Title:    "My Passport",
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(`INSERT INTO documents`).
			WithArgs(doc.TenantID, doc.UserID, doc.Kind, doc.Title).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("doc-1"))

		id, err := repo.Create(context.Background(), doc)

		assert.NoError(t, err)
		assert.Equal(t, "doc-1", id)
	})
}

func TestSQLDocumentRepository_GetByID_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLDocumentRepository(sqlxDB)

	docID := "doc-1"

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows([]string{"id", "tenant_id", "user_id", "kind", "title", "current_version_id", "created_at", "updated_at"}).
			AddRow(docID, "t1", "u1", "kind1", "Title 1", "v1", now, now)

		mock.ExpectQuery(`SELECT \* FROM documents WHERE id = \$1`).
			WithArgs(docID).
			WillReturnRows(rows)

		doc, err := repo.GetByID(context.Background(), docID)

		assert.NoError(t, err)
		assert.NotNil(t, doc)
		assert.Equal(t, docID, doc.ID)
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM documents`).
			WithArgs(docID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		doc, err := repo.GetByID(context.Background(), docID)

		assert.Error(t, err)
		assert.Nil(t, doc)
		assert.Equal(t, ErrNotFound, err)
	})
}

func TestSQLDocumentRepository_ListByUserID_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLDocumentRepository(sqlx.NewDb(db, "sqlmock"))

	userID := "u1"

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "title", "created_at", "updated_at"}).
			AddRow("d1", userID, "Doc 1", time.Now(), time.Now()).
			AddRow("d2", userID, "Doc 2", time.Now(), time.Now())

		mock.ExpectQuery(`SELECT \* FROM documents WHERE user_id = \$1 ORDER BY created_at DESC`).
			WithArgs(userID).
			WillReturnRows(rows)

		docs, err := repo.ListByUserID(context.Background(), userID)
		assert.NoError(t, err)
		assert.Len(t, docs, 2)
	})
}

func TestSQLDocumentRepository_Delete_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLDocumentRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM documents WHERE id = \$1`).
			WithArgs("doc-1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Delete(context.Background(), "doc-1")
		assert.NoError(t, err)
	})
}

func TestSQLDocumentRepository_CreateVersion_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLDocumentRepository(sqlx.NewDb(db, "sqlmock"))

	ver := &models.DocumentVersion{
		TenantID:    "t1",
		DocumentID:  "d1",
		StoragePath: "s3/path",
		MimeType:    "pdf",
		SizeBytes:   1024,
		UploadedBy:  "u1",
		Bucket:      sql.NullString{String: "bucket", Valid: true},
		ObjectKey:   sql.NullString{String: "key", Valid: true},
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(`INSERT INTO document_versions`).
			WithArgs(ver.TenantID, ver.DocumentID, ver.StoragePath, ver.MimeType, ver.SizeBytes, ver.UploadedBy, ver.Bucket, ver.ObjectKey).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("v1"))

		id, err := repo.CreateVersion(context.Background(), ver)
		assert.NoError(t, err)
		assert.Equal(t, "v1", id)
	})
}

func TestSQLDocumentRepository_GetVersion_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLDocumentRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "mime_type"}).AddRow("v1", "pdf")
		mock.ExpectQuery(`SELECT .* FROM document_versions WHERE id = \$1`).
			WithArgs("v1").
			WillReturnRows(rows)

		ver, err := repo.GetVersion(context.Background(), "v1")
		assert.NoError(t, err)
		assert.Equal(t, "pdf", ver.MimeType)
	})
}

func TestSQLDocumentRepository_GetVersionsByDocumentID_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLDocumentRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "document_id"}).
			AddRow("v2", "d1").
			AddRow("v1", "d1")
		
		mock.ExpectQuery(`SELECT .* FROM document_versions WHERE document_id = \$1 ORDER BY created_at DESC`).
			WithArgs("d1").
			WillReturnRows(rows)
		
		vers, err := repo.GetVersionsByDocumentID(context.Background(), "d1")
		assert.NoError(t, err)
		assert.Len(t, vers, 2)
	})
}

func TestSQLDocumentRepository_GetLatestVersion_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLDocumentRepository(sqlxDB)

	docID := "doc-1"

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows([]string{"id", "tenant_id", "document_id", "storage_path", "mime_type", "size_bytes", "uploaded_by", "bucket", "object_key", "note", "created_at"}).
			AddRow("v1", "t1", docID, "path1", "pdf", 100, "u1", "b1", "k1", "note1", now)

		mock.ExpectQuery(`SELECT (.+) FROM document_versions WHERE document_id = \$1 ORDER BY created_at DESC LIMIT 1`).
			WithArgs(docID).
			WillReturnRows(rows)

		ver, err := repo.GetLatestVersion(context.Background(), docID)

		assert.NoError(t, err)
		assert.NotNil(t, ver)
		assert.Equal(t, "v1", ver.ID)
	})
}

func TestSQLDocumentRepository_SetCurrentVersion_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLDocumentRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(`UPDATE documents SET current_version_id = \$1, updated_at = NOW\(\) WHERE id = \$2`).
			WithArgs("v1", "d1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.SetCurrentVersion(context.Background(), "d1", "v1")
		assert.NoError(t, err)
	})
}

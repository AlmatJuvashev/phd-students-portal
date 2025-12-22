package repository

import (
	"context"
	"database/sql"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type DocumentRepository interface {
	Create(ctx context.Context, doc *models.Document) (string, error)
	GetByID(ctx context.Context, id string) (*models.Document, error)
	ListByUserID(ctx context.Context, userID string) ([]models.Document, error)
	Delete(ctx context.Context, id string) error
	
	CreateVersion(ctx context.Context, ver *models.DocumentVersion) (string, error)
	GetVersion(ctx context.Context, id string) (*models.DocumentVersion, error)
	GetVersionsByDocumentID(ctx context.Context, docID string) ([]models.DocumentVersion, error)
	GetLatestVersion(ctx context.Context, docID string) (*models.DocumentVersion, error)
	
	SetCurrentVersion(ctx context.Context, docID, verID string) error
}

type SQLDocumentRepository struct {
	db *sqlx.DB
}

func NewSQLDocumentRepository(db *sqlx.DB) *SQLDocumentRepository {
	return &SQLDocumentRepository{db: db}
}

func (r *SQLDocumentRepository) Create(ctx context.Context, doc *models.Document) (string, error) {
	var id string
	query := `INSERT INTO documents (tenant_id, user_id, kind, title, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id`
	// Fallback for missing tenantID if not provided (should be enforced by service/handler)
	err := r.db.QueryRowContext(ctx, query, doc.TenantID, doc.UserID, doc.Kind, doc.Title).Scan(&id)
	return id, err
}

func (r *SQLDocumentRepository) GetByID(ctx context.Context, id string) (*models.Document, error) {
	var doc models.Document
	query := `SELECT * FROM documents WHERE id = $1`
	err := r.db.GetContext(ctx, &doc, query, id)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &doc, err
}

func (r *SQLDocumentRepository) ListByUserID(ctx context.Context, userID string) ([]models.Document, error) {
	var docs []models.Document
	query := `SELECT * FROM documents WHERE user_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &docs, query, userID)
	return docs, err
}

func (r *SQLDocumentRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM documents WHERE id = $1`, id)
	return err
}

func (r *SQLDocumentRepository) CreateVersion(ctx context.Context, ver *models.DocumentVersion) (string, error) {
	var id string
	query := `INSERT INTO document_versions (tenant_id, document_id, storage_path, mime_type, size_bytes, uploaded_by, bucket, object_key, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW()) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, 
		ver.TenantID, ver.DocumentID, ver.StoragePath, ver.MimeType, ver.SizeBytes, ver.UploadedBy, ver.Bucket, ver.ObjectKey).Scan(&id)
	return id, err
}

func (r *SQLDocumentRepository) GetVersion(ctx context.Context, id string) (*models.DocumentVersion, error) {
	var ver models.DocumentVersion
	versionBaseSelect := `SELECT id, tenant_id, document_id, storage_path, mime_type, size_bytes, uploaded_by, bucket, object_key, note, created_at FROM document_versions`
	query := versionBaseSelect + ` WHERE id = $1`
	err := r.db.GetContext(ctx, &ver, query, id)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &ver, err
}

func (r *SQLDocumentRepository) GetVersionsByDocumentID(ctx context.Context, docID string) ([]models.DocumentVersion, error) {
	var vers []models.DocumentVersion
	// Explicit select to avoid scan errors on unused columns (etag, checksum, etc)
	versionBaseSelect := `SELECT id, tenant_id, document_id, storage_path, mime_type, size_bytes, uploaded_by, bucket, object_key, note, created_at FROM document_versions`

	query := versionBaseSelect + ` WHERE document_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &vers, query, docID)
	return vers, err
}

func (r *SQLDocumentRepository) GetLatestVersion(ctx context.Context, docID string) (*models.DocumentVersion, error) {
	var ver models.DocumentVersion
	versionBaseSelect := `SELECT id, tenant_id, document_id, storage_path, mime_type, size_bytes, uploaded_by, bucket, object_key, note, created_at FROM document_versions`
	
	query := versionBaseSelect + ` WHERE document_id = $1 ORDER BY created_at DESC LIMIT 1`
	err := r.db.GetContext(ctx, &ver, query, docID)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &ver, err
}

func (r *SQLDocumentRepository) SetCurrentVersion(ctx context.Context, docID, verID string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE documents SET current_version_id = $1, updated_at = NOW() WHERE id = $2`, verID, docID)
	return err
}

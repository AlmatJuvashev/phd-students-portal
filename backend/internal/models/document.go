package models

import (
	"database/sql"
	"time"
)

type Document struct {
	ID               string    `db:"id" json:"id"`
	TenantID         string    `db:"tenant_id" json:"tenant_id"`
	UserID           string    `db:"user_id" json:"user_id"`
	Title            string    `db:"title" json:"title"`
	Kind             string    `db:"kind" json:"kind"` // e.g., "thesis", "report", "admin"
	CurrentVersionID *string   `db:"current_version_id" json:"current_version_id"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

type DocumentVersion struct {
	ID          string         `db:"id" json:"id"`
	TenantID    string         `db:"tenant_id" json:"tenant_id"`
	DocumentID  string         `db:"document_id" json:"document_id"`
	StoragePath string         `db:"storage_path" json:"storage_path"` // Local path or S3 key fallback
	MimeType    string         `db:"mime_type" json:"mime_type"`
	SizeBytes   int64          `db:"size_bytes" json:"size_bytes"`
	UploadedBy  string         `db:"uploaded_by" json:"uploaded_by"`
	Bucket      sql.NullString `db:"bucket" json:"bucket"`         // S3 Bucket
	ObjectKey   sql.NullString `db:"object_key" json:"object_key"` // S3 Key
	Note        sql.NullString `db:"note" json:"note"`             // Review/upload note
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
}

// DTOs for List/Get responses often need joined data
type DocumentWithVersion struct {
	Document
	CurrentVersion *DocumentVersion `json:"current_version,omitempty"`
}

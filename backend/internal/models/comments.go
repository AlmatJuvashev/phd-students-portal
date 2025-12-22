package models

type Comment struct {
	ID        string  `db:"id" json:"id"`
	TenantID  string  `db:"tenant_id" json:"tenant_id"`
	DocumentID string  `db:"document_id" json:"document_id"`
	UserID    string  `db:"user_id" json:"user_id"`
	Content   string  `db:"content" json:"content"`
	ParentID  *string `db:"parent_id" json:"parent_id"`
	CreatedAt string  `db:"created_at" json:"created_at"`
}

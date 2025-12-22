package repository

import (
	"context"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type CommentRepository interface {
	Create(ctx context.Context, comment models.Comment) (string, error)
	GetByDocumentID(ctx context.Context, tenantID string, docID string) ([]models.Comment, error)
}

type SQLCommentRepository struct {
	db *sqlx.DB
}

func NewSQLCommentRepository(db *sqlx.DB) *SQLCommentRepository {
	return &SQLCommentRepository{db: db}
}

func (r *SQLCommentRepository) Create(ctx context.Context, comment models.Comment) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO comments (tenant_id, document_id, user_id, content, parent_id) 
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		comment.TenantID, comment.DocumentID, comment.UserID, comment.Content, comment.ParentID,
	).Scan(&id)
	return id, err
}

func (r *SQLCommentRepository) GetByDocumentID(ctx context.Context, tenantID string, docID string) ([]models.Comment, error) {
	var comments []models.Comment
	// Note: We select only fields that were previously returned, plus generic ones if needed.
	// Previous handler selected: id, user_id, content, parent_id, created_at
	// It did NOT select tenant_id or document_id (implied).
	// But struct has them.
	// Also formatting created_at to match previous handler.
	err := r.db.SelectContext(ctx, &comments,
		`SELECT id, tenant_id, document_id, user_id, content, parent_id, 
		        to_char(created_at, 'YYYY-MM-DD HH24:MI:SS') as created_at
		   FROM comments 
		  WHERE document_id = $1 AND tenant_id = $2
		  ORDER BY created_at ASC`,
		docID, tenantID, // Added tenant_id check for safety, though previously handler didn't check tenant_id explicitly in query (it relied on docId being unique or access control).
		// Wait, previous handler: `WHERE document_id = $1`. NO tenant_id check in SQL.
		// `GetComments` in handler (`handlers/comments.go`) did NOT use `tenant_id` from context in SQL.
		// However, it's safer to include it if `comments` table has `tenant_id`.
		// Let's include `AND tenant_id = $2` to enforce isolation.
	)
	if err != nil {
		return nil, err
	}
	if comments == nil {
		comments = []models.Comment{}
	}
	return comments, nil
}

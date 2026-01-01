package repository

import (
	"context"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type ItemBankRepository interface {
	// Banks
	CreateBank(ctx context.Context, b *models.QuestionBank) error
	GetBank(ctx context.Context, id string) (*models.QuestionBank, error)
	ListBanks(ctx context.Context, tenantID string) ([]models.QuestionBank, error)
	UpdateBank(ctx context.Context, b *models.QuestionBank) error
	DeleteBank(ctx context.Context, id string) error

	// Items
	CreateItem(ctx context.Context, item *models.QuestionItem) error
	GetItem(ctx context.Context, id string) (*models.QuestionItem, error)
	ListItems(ctx context.Context, bankID string) ([]models.QuestionItem, error)
	UpdateItem(ctx context.Context, item *models.QuestionItem) error
	DeleteItem(ctx context.Context, id string) error
}

type SQLItemBankRepository struct {
	db *sqlx.DB
}

func NewSQLItemBankRepository(db *sqlx.DB) *SQLItemBankRepository {
	return &SQLItemBankRepository{db: db}
}

// --- Banks ---

func (r *SQLItemBankRepository) CreateBank(ctx context.Context, b *models.QuestionBank) error {
	return r.db.QueryRowxContext(ctx, `
		INSERT INTO question_banks (tenant_id, name, description, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`,
		b.TenantID, b.Name, b.Description, b.IsActive,
	).Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt)
}

func (r *SQLItemBankRepository) GetBank(ctx context.Context, id string) (*models.QuestionBank, error) {
	var b models.QuestionBank
	err := sqlx.GetContext(ctx, r.db, &b, `SELECT * FROM question_banks WHERE id=$1`, id)
	return &b, err
}

func (r *SQLItemBankRepository) ListBanks(ctx context.Context, tenantID string) ([]models.QuestionBank, error) {
	var list []models.QuestionBank
	err := sqlx.SelectContext(ctx, r.db, &list, `SELECT * FROM question_banks WHERE tenant_id=$1 ORDER BY created_at DESC`, tenantID)
	return list, err
}

func (r *SQLItemBankRepository) UpdateBank(ctx context.Context, b *models.QuestionBank) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE question_banks SET name=$1, description=$2, is_active=$3, updated_at=now()
		WHERE id=$4`,
		b.Name, b.Description, b.IsActive, b.ID)
	return err
}

func (r *SQLItemBankRepository) DeleteBank(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM question_banks WHERE id=$1`, id)
	return err
}

// --- Items ---

func (r *SQLItemBankRepository) CreateItem(ctx context.Context, item *models.QuestionItem) error {
	return r.db.QueryRowxContext(ctx, `
		INSERT INTO question_items (bank_id, type, content, difficulty, tags, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`,
		item.BankID, item.Type, item.Content, item.Difficulty, item.Tags, item.IsActive,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
}

func (r *SQLItemBankRepository) GetItem(ctx context.Context, id string) (*models.QuestionItem, error) {
	var item models.QuestionItem
	err := sqlx.GetContext(ctx, r.db, &item, `SELECT * FROM question_items WHERE id=$1`, id)
	return &item, err
}

func (r *SQLItemBankRepository) ListItems(ctx context.Context, bankID string) ([]models.QuestionItem, error) {
	var list []models.QuestionItem
	err := sqlx.SelectContext(ctx, r.db, &list, `SELECT * FROM question_items WHERE bank_id=$1 ORDER BY created_at DESC`, bankID)
	return list, err
}

func (r *SQLItemBankRepository) UpdateItem(ctx context.Context, item *models.QuestionItem) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE question_items SET type=$1, content=$2, difficulty=$3, tags=$4, is_active=$5, updated_at=now()
		WHERE id=$6`,
		item.Type, item.Content, item.Difficulty, item.Tags, item.IsActive, item.ID)
	return err
}

func (r *SQLItemBankRepository) DeleteItem(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM question_items WHERE id=$1`, id)
	return err
}

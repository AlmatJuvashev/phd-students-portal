package models

import (
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/lib/pq"
)

// QuestionBank organizes a collection of questions.
// e.g., "Anatomy 101 Midterm Pool", "General Knowledge"
type QuestionBank struct {
	ID          string    `db:"id" json:"id"`
	TenantID    string    `db:"tenant_id" json:"tenant_id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	IsActive    bool      `db:"is_active" json:"is_active"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// QuestionItem represents a single reusable question.
type QuestionItem struct {
	ID        string         `db:"id" json:"id"`
	BankID    string         `db:"bank_id" json:"bank_id"`
	Type      string         `db:"type" json:"type"` // multiple_choice, true_false, essay
	Content   types.JSONText `db:"content" json:"content"` // JSON: { "text": "...", "options": [...], "answer": "..." }
	Difficulty int           `db:"difficulty" json:"difficulty"` // 1-5
	Tags      pq.StringArray `db:"tags" json:"tags"`
	IsActive  bool           `db:"is_active" json:"is_active"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt time.Time      `db:"updated_at" json:"updated_at"`
}

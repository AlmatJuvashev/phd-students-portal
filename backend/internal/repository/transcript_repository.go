package repository

import (
	"context"
	"fmt"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type TranscriptRepository interface {
	GetStudentGrades(ctx context.Context, studentID string) ([]models.TermGrade, error)
	// Additional methods for mass-inserting grades could go here provided later
}

type SQLTranscriptRepository struct {
	db *sqlx.DB
}

func NewSQLTranscriptRepository(db *sqlx.DB) *SQLTranscriptRepository {
	return &SQLTranscriptRepository{db: db}
}

func (r *SQLTranscriptRepository) GetStudentGrades(ctx context.Context, studentID string) ([]models.TermGrade, error) {
	var grades []models.TermGrade
	// Join with Academic Terms to get Term Name if needed, but TermGrade struct has TermID.
	// We might want to ORDER BY term date.
	// But `term_grades` only has `term_id`.
	// The Service will grouping by TermID.
	// It's better if we fetch ordered by Term StartDate.
	// So we need to join academic_terms.
	
	query := `
		SELECT tg.* 
		FROM term_grades tg
		JOIN academic_terms at ON tg.term_id = at.id
		WHERE tg.student_id = $1
		ORDER BY at.start_date ASC
	`
	
	err := r.db.SelectContext(ctx, &grades, query, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student grades: %w", err)
	}
	return grades, nil
}

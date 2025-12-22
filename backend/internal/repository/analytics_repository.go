package repository

import (
	"context"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type AnalyticsRepository interface {
	GetStudentsByStage(ctx context.Context) ([]models.StudentStageStats, error)
	GetAdvisorLoad(ctx context.Context) ([]models.AdvisorLoadStats, error)
	GetOverdueTasks(ctx context.Context) ([]models.OverdueTaskStats, error)
}

type SQLAnalyticsRepository struct {
	db *sqlx.DB
}

func NewSQLAnalyticsRepository(db *sqlx.DB) *SQLAnalyticsRepository {
	return &SQLAnalyticsRepository{db: db}
}

func (r *SQLAnalyticsRepository) GetStudentsByStage(ctx context.Context) ([]models.StudentStageStats, error) {
	query := `
		SELECT 
			COALESCE(js.state, 'Not Started') as stage, 
			COUNT(*) as count
		FROM users u
		LEFT JOIN journey_states js ON u.id = js.user_id
		WHERE u.role = 'student'
		GROUP BY js.state`
	
	var stats []models.StudentStageStats
	err := r.db.SelectContext(ctx, &stats, query)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (r *SQLAnalyticsRepository) GetAdvisorLoad(ctx context.Context) ([]models.AdvisorLoadStats, error) {
	query := `
		SELECT 
			COALESCE(u.first_name || ' ' || u.last_name, 'Unknown') as advisor_name,
			COUNT(sa.student_id) as student_count
		FROM users u
		JOIN student_advisors sa ON u.id = sa.advisor_id
		WHERE u.role = 'advisor'
		GROUP BY u.id, u.first_name, u.last_name
		ORDER BY student_count DESC`
	
	var stats []models.AdvisorLoadStats
	err := r.db.SelectContext(ctx, &stats, query)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (r *SQLAnalyticsRepository) GetOverdueTasks(ctx context.Context) ([]models.OverdueTaskStats, error) {
	query := `
		SELECT 
			nd.node_id,
			COUNT(*) as count
		FROM node_deadlines nd
		LEFT JOIN node_instances ni ON nd.user_id = ni.user_id AND nd.node_id = ni.node_id
		WHERE nd.due_at < NOW() 
		AND (ni.state IS NULL OR ni.state != 'done')
		GROUP BY nd.node_id
		ORDER BY count DESC`
	
	var stats []models.OverdueTaskStats
	err := r.db.SelectContext(ctx, &stats, query)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

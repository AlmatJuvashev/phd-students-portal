package services

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type AnalyticsService struct {
	db *sqlx.DB
}

func NewAnalyticsService(db *sqlx.DB) *AnalyticsService {
	return &AnalyticsService{db: db}
}

type StudentStageStats struct {
	Stage string `db:"stage" json:"stage"`
	Count int    `db:"count" json:"count"`
}

func (s *AnalyticsService) GetStudentsByStage(ctx context.Context) ([]StudentStageStats, error) {
	query := `
		SELECT 
			COALESCE(js.state, 'Not Started') as stage, 
			COUNT(*) as count
		FROM users u
		LEFT JOIN journey_states js ON u.id = js.user_id
		WHERE u.role = 'student'
		GROUP BY js.state`
	
	var stats []StudentStageStats
	err := s.db.SelectContext(ctx, &stats, query)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

type AdvisorLoadStats struct {
	AdvisorName string `db:"advisor_name" json:"advisor_name"`
	StudentCount int   `db:"student_count" json:"student_count"`
}

func (s *AnalyticsService) GetAdvisorLoad(ctx context.Context) ([]AdvisorLoadStats, error) {
	query := `
		SELECT 
			COALESCE(u.first_name || ' ' || u.last_name, 'Unknown') as advisor_name,
			COUNT(sa.student_id) as student_count
		FROM users u
		JOIN student_advisors sa ON u.id = sa.advisor_id
		WHERE u.role = 'advisor'
		GROUP BY u.id, u.first_name, u.last_name
		ORDER BY student_count DESC`
	
	var stats []AdvisorLoadStats
	err := s.db.SelectContext(ctx, &stats, query)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

type OverdueTaskStats struct {
	NodeID string `db:"node_id" json:"node_id"`
	Count  int    `db:"count" json:"count"`
}

func (s *AnalyticsService) GetOverdueTasks(ctx context.Context) ([]OverdueTaskStats, error) {
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
	
	var stats []OverdueTaskStats
	err := s.db.SelectContext(ctx, &stats, query)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

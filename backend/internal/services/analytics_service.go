package services

import (
	"context"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type AnalyticsService struct {
	repo repository.AnalyticsRepository
}

func NewAnalyticsService(repo repository.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{repo: repo}
}

func (s *AnalyticsService) GetStudentsByStage(ctx context.Context) ([]models.StudentStageStats, error) {
	return s.repo.GetStudentsByStage(ctx)
}

func (s *AnalyticsService) GetAdvisorLoad(ctx context.Context) ([]models.AdvisorLoadStats, error) {
	return s.repo.GetAdvisorLoad(ctx)
}

func (s *AnalyticsService) GetOverdueTasks(ctx context.Context) ([]models.OverdueTaskStats, error) {
	return s.repo.GetOverdueTasks(ctx)
}

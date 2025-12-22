package services

import (
	"context"
	"encoding/json"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type ChecklistService struct {
	repo repository.ChecklistRepository
}

func NewChecklistService(repo repository.ChecklistRepository) *ChecklistService {
	return &ChecklistService{repo: repo}
}

func (s *ChecklistService) GetModules(ctx context.Context) ([]models.ChecklistModule, error) {
	return s.repo.ListModules(ctx)
}

func (s *ChecklistService) GetStepsByModule(ctx context.Context, moduleCode string) ([]models.ChecklistStep, error) {
	return s.repo.ListStepsByModule(ctx, moduleCode)
}

func (s *ChecklistService) GetStudentSteps(ctx context.Context, userID string) ([]struct {
	StepID string `db:"step_id" json:"step_id"`
	Status string `db:"status" json:"status"`
}, error) {
	return s.repo.ListStudentSteps(ctx, userID)
}

func (s *ChecklistService) UpdateStudentStep(ctx context.Context, userID, stepID, status string, data json.RawMessage) error {
	return s.repo.UpsertStudentStep(ctx, userID, stepID, status, data)
}

func (s *ChecklistService) GetAdvisorInbox(ctx context.Context) ([]models.AdvisorInboxItem, error) {
	return s.repo.GetAdvisorInbox(ctx)
}

func (s *ChecklistService) ApproveStep(ctx context.Context, userID, stepID, authorID, comment string, mentions []string) error {
	if err := s.repo.ApproveStep(ctx, userID, stepID); err != nil {
		return err
	}
	if comment != "" {
		return s.repo.AddCommentToLatestDocument(ctx, userID, comment, authorID, mentions)
	}
	return nil
}

func (s *ChecklistService) ReturnStep(ctx context.Context, userID, stepID, authorID, comment string, mentions []string) error {
	if err := s.repo.ReturnStep(ctx, userID, stepID); err != nil {
		return err
	}
	if comment != "" {
		return s.repo.AddCommentToLatestDocument(ctx, userID, comment, authorID, mentions)
	}
	return nil
}

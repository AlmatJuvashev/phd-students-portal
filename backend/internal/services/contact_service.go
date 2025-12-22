package services

import (
	"context"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type ContactService struct {
	repo repository.ContactRepository
}

func NewContactService(repo repository.ContactRepository) *ContactService {
	return &ContactService{repo: repo}
}

func (s *ContactService) ListPublic(ctx context.Context, tenantID string) ([]models.Contact, error) {
	return s.repo.ListPublic(ctx, tenantID)
}

func (s *ContactService) ListAdmin(ctx context.Context, tenantID string, includeInactive bool) ([]models.Contact, error) {
	return s.repo.ListAdmin(ctx, tenantID, includeInactive)
}

func (s *ContactService) Create(ctx context.Context, tenantID string, contact models.Contact) (string, error) {
	return s.repo.Create(ctx, tenantID, contact)
}

func (s *ContactService) Update(ctx context.Context, tenantID string, id string, updates map[string]interface{}) error {
	// Business logic: check if updates empty?
	if len(updates) == 0 {
		return nil
	}
	return s.repo.Update(ctx, tenantID, id, updates)
}

func (s *ContactService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}

package services

import (
	"context"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type ItemBankService struct {
	repo repository.ItemBankRepository
}

func NewItemBankService(repo repository.ItemBankRepository) *ItemBankService {
	return &ItemBankService{repo: repo}
}

// --- Banks ---

func (s *ItemBankService) CreateBank(ctx context.Context, b *models.QuestionBank) error {
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()
	if b.Name == "" {
		b.Name = "Untitled Bank"
	}
	b.IsActive = true
	return s.repo.CreateBank(ctx, b)
}

func (s *ItemBankService) ListBanks(ctx context.Context, tenantID string) ([]models.QuestionBank, error) {
	return s.repo.ListBanks(ctx, tenantID)
}

func (s *ItemBankService) GetBank(ctx context.Context, id string) (*models.QuestionBank, error) {
	return s.repo.GetBank(ctx, id)
}

// --- Items ---

func (s *ItemBankService) CreateItem(ctx context.Context, item *models.QuestionItem) error {
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	item.IsActive = true
	return s.repo.CreateItem(ctx, item)
}

func (s *ItemBankService) ListItems(ctx context.Context, bankID string) ([]models.QuestionItem, error) {
	return s.repo.ListItems(ctx, bankID)
}

func (s *ItemBankService) UpdateItem(ctx context.Context, item *models.QuestionItem) error {
	item.UpdatedAt = time.Now()
	return s.repo.UpdateItem(ctx, item)
}

func (s *ItemBankService) DeleteItem(ctx context.Context, id string) error {
	return s.repo.DeleteItem(ctx, id)
}

package services

import (
	"context"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type ItemBankService struct {
	repo repository.AssessmentRepository
}

func NewItemBankService(repo repository.AssessmentRepository) *ItemBankService {
	return &ItemBankService{repo: repo}
}

// --- Banks ---

func (s *ItemBankService) CreateBank(ctx context.Context, b *models.QuestionBank) error {
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()
	_, err := s.repo.CreateQuestionBank(ctx, *b) // Repo takes value, returns pointer
	// Ideally update b with returned result ID
	return err
}

func (s *ItemBankService) ListBanks(ctx context.Context, tenantID string) ([]models.QuestionBank, error) {
	return s.repo.ListQuestionBanks(ctx, tenantID)
}

func (s *ItemBankService) GetBank(ctx context.Context, id string) (*models.QuestionBank, error) {
	return s.repo.GetQuestionBank(ctx, id)
}

// --- Items ---

func (s *ItemBankService) CreateItem(ctx context.Context, item *models.Question) error {
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	_, err := s.repo.CreateQuestion(ctx, *item)
	return err
}

func (s *ItemBankService) ListItems(ctx context.Context, bankID string) ([]models.Question, error) {
	return s.repo.ListQuestionsByBank(ctx, bankID)
}

// Update/Delete not yet in AssessmentRepository, skipping for now or adding TODO
func (s *ItemBankService) UpdateItem(ctx context.Context, item *models.Question) error {
	item.UpdatedAt = time.Now()
	return s.repo.UpdateQuestion(ctx, *item)
}

func (s *ItemBankService) DeleteItem(ctx context.Context, id string) error {
	return s.repo.DeleteQuestion(ctx, id)
}

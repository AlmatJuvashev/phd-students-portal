package services

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockItemBankRepo
type MockItemBankRepo struct {
	mock.Mock
}
func (m *MockItemBankRepo) CreateBank(ctx context.Context, b *models.QuestionBank) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}
func (m *MockItemBankRepo) GetBank(ctx context.Context, id string) (*models.QuestionBank, error) { return nil, nil }
func (m *MockItemBankRepo) ListBanks(ctx context.Context, tenantID string) ([]models.QuestionBank, error) { return nil, nil }
func (m *MockItemBankRepo) UpdateBank(ctx context.Context, b *models.QuestionBank) error { return nil }
func (m *MockItemBankRepo) DeleteBank(ctx context.Context, id string) error { return nil }

func (m *MockItemBankRepo) CreateItem(ctx context.Context, i *models.QuestionItem) error {
	args := m.Called(ctx, i)
	return args.Error(0)
}
func (m *MockItemBankRepo) GetItem(ctx context.Context, id string) (*models.QuestionItem, error) { return nil, nil }
func (m *MockItemBankRepo) ListItems(ctx context.Context, bankID string) ([]models.QuestionItem, error) { return nil, nil }
func (m *MockItemBankRepo) UpdateItem(ctx context.Context, i *models.QuestionItem) error { return nil }
func (m *MockItemBankRepo) DeleteItem(ctx context.Context, id string) error { return nil }


func TestItemBankService_CreateBank(t *testing.T) {
	mockRepo := new(MockItemBankRepo)
	svc := NewItemBankService(mockRepo)
	ctx := context.Background()

	// 1. Success Case
	mockRepo.On("CreateBank", ctx, mock.MatchedBy(func(b *models.QuestionBank) bool {
		return b.Name == "Anatomy" && b.IsActive == true
	})).Return(nil)

	err := svc.CreateBank(ctx, &models.QuestionBank{Name: "Anatomy", TenantID: "t1"})
	assert.NoError(t, err)

	// 2. Default Name Case
	mockRepo.On("CreateBank", ctx, mock.MatchedBy(func(b *models.QuestionBank) bool {
		return b.Name == "Untitled Bank"
	})).Return(nil)

	err = svc.CreateBank(ctx, &models.QuestionBank{TenantID: "t1"}) // No name
	assert.NoError(t, err)
}

func TestItemBankService_CreateItem(t *testing.T) {
	mockRepo := new(MockItemBankRepo)
	svc := NewItemBankService(mockRepo)
	ctx := context.Background()

	mockRepo.On("CreateItem", ctx, mock.MatchedBy(func(i *models.QuestionItem) bool {
		return i.Type == "multiple_choice" && i.IsActive == true
	})).Return(nil)

	err := svc.CreateItem(ctx, &models.QuestionItem{BankID: "b1", Type: "multiple_choice"})
	assert.NoError(t, err)
}

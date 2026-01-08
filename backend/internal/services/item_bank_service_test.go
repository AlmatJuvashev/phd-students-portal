package services_test

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func strPtr(s string) *string {
	return &s
}

func TestItemBankService_Banks(t *testing.T) {
	ctx := context.Background()

	t.Run("CreateBank", func(t *testing.T) {
		repo := new(services.MockAssessmentRepository)
		svc := services.NewItemBankService(repo)
		
		bank := &models.QuestionBank{
			TenantID: "t1",
			Title:    "Math Bank",
			Subject:  strPtr("Math"),
		}
		
		repo.On("CreateQuestionBank", ctx, mock.MatchedBy(func(b models.QuestionBank) bool {
			return b.Title == "Math Bank" && !b.CreatedAt.IsZero()
		})).Return(&models.QuestionBank{ID: "b1", Title: "Math Bank"}, nil)

		err := svc.CreateBank(ctx, bank)
		assert.NoError(t, err)
		assert.Equal(t, "b1", bank.ID)
	})

	t.Run("ListBanks", func(t *testing.T) {
		repo := new(services.MockAssessmentRepository)
		svc := services.NewItemBankService(repo)
		
		expected := []models.QuestionBank{{ID: "b1", Title: "Bank 1"}}
		repo.On("ListQuestionBanks", ctx, "t1").Return(expected, nil)

		banks, err := svc.ListBanks(ctx, "t1")
		assert.NoError(t, err)
		assert.Len(t, banks, 1)
		assert.Equal(t, "Bank 1", banks[0].Title)
	})

	t.Run("GetBank", func(t *testing.T) {
		repo := new(services.MockAssessmentRepository)
		svc := services.NewItemBankService(repo)
		
		repo.On("GetQuestionBank", ctx, "b1").Return(&models.QuestionBank{ID: "b1"}, nil)

		bank, err := svc.GetBank(ctx, "b1")
		assert.NoError(t, err)
		assert.Equal(t, "b1", bank.ID)
	})

	t.Run("UpdateBank", func(t *testing.T) {
		repo := new(services.MockAssessmentRepository)
		svc := services.NewItemBankService(repo)
		
		bank := &models.QuestionBank{ID: "b1", Title: "Updated"}
		repo.On("UpdateQuestionBank", ctx, mock.MatchedBy(func(b models.QuestionBank) bool {
			return b.ID == "b1" && !b.UpdatedAt.IsZero()
		})).Return(nil)

		err := svc.UpdateBank(ctx, bank)
		assert.NoError(t, err)
	})

	t.Run("DeleteBank", func(t *testing.T) {
		repo := new(services.MockAssessmentRepository)
		svc := services.NewItemBankService(repo)
		
		repo.On("DeleteQuestionBank", ctx, "b1").Return(nil)

		err := svc.DeleteBank(ctx, "b1")
		assert.NoError(t, err)
	})
}

func TestItemBankService_Items(t *testing.T) {
	ctx := context.Background()

	t.Run("CreateItem", func(t *testing.T) {
		repo := new(services.MockAssessmentRepository)
		svc := services.NewItemBankService(repo)
		
		item := &models.Question{
			BankID: "b1",
			Stem:   "What is 2+2?",
			Type:   models.QuestionTypeMCQ,
		}
		
		repo.On("CreateQuestion", ctx, mock.MatchedBy(func(q models.Question) bool {
			return q.Stem == "What is 2+2?" && !q.CreatedAt.IsZero()
		})).Return(&models.Question{ID: "q1", Stem: "What is 2+2?"}, nil)

		err := svc.CreateItem(ctx, item)
		assert.NoError(t, err)
		assert.Equal(t, "q1", item.ID)
	})

	t.Run("ListItems", func(t *testing.T) {
		repo := new(services.MockAssessmentRepository)
		svc := services.NewItemBankService(repo)
		
		expected := []models.Question{{ID: "q1"}}
		repo.On("ListQuestionsByBank", ctx, "b1").Return(expected, nil)

		items, err := svc.ListItems(ctx, "b1")
		assert.NoError(t, err)
		assert.Len(t, items, 1)
	})

	t.Run("GetItem", func(t *testing.T) {
		repo := new(services.MockAssessmentRepository)
		svc := services.NewItemBankService(repo)
		
		repo.On("GetQuestion", ctx, "q1").Return(&models.Question{ID: "q1"}, nil)

		item, err := svc.GetItem(ctx, "q1")
		assert.NoError(t, err)
		assert.Equal(t, "q1", item.ID)
	})

	t.Run("UpdateItem", func(t *testing.T) {
		repo := new(services.MockAssessmentRepository)
		svc := services.NewItemBankService(repo)
		
		item := &models.Question{ID: "q1", Stem: "Updated"}
		repo.On("UpdateQuestion", ctx, mock.MatchedBy(func(q models.Question) bool {
			return q.ID == "q1" && !q.UpdatedAt.IsZero()
		})).Return(nil)

		err := svc.UpdateItem(ctx, item)
		assert.NoError(t, err)
	})

	t.Run("DeleteItem", func(t *testing.T) {
		repo := new(services.MockAssessmentRepository)
		svc := services.NewItemBankService(repo)
		
		repo.On("DeleteQuestion", ctx, "q1").Return(nil)

		err := svc.DeleteItem(ctx, "q1")
		assert.NoError(t, err)
	})
}

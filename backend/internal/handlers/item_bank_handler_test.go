package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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
func (m *MockItemBankRepo) CreateItem(ctx context.Context, i *models.QuestionItem) error { return nil }
func (m *MockItemBankRepo) GetItem(ctx context.Context, id string) (*models.QuestionItem, error) { return nil, nil }
func (m *MockItemBankRepo) ListItems(ctx context.Context, bankID string) ([]models.QuestionItem, error) { return nil, nil }
func (m *MockItemBankRepo) UpdateItem(ctx context.Context, i *models.QuestionItem) error { return nil }
func (m *MockItemBankRepo) DeleteItem(ctx context.Context, id string) error { return nil }

func TestItemBankHandler_CreateBank(t *testing.T) {
	mockRepo := new(MockItemBankRepo)
	svc := services.NewItemBankService(mockRepo)
	h := handlers.NewItemBankHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", "t-1")
	})
	r.POST("/item-banks/banks", h.CreateBank)

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("CreateBank", mock.Anything, mock.MatchedBy(func(b *models.QuestionBank) bool {
			return b.Name == "My Bank" && b.TenantID == "t-1"
		})).Return(nil)

		reqBody := map[string]interface{}{"name": "My Bank"}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/item-banks/banks", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

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

// Re-using MockGradingRepo from services package would be ideal, but due to package boundary (handlers_test vs services),
// we might need to redefine or import deeply.
// For simplicity in this test file, I'll define a local mock repo since interface is in repository package.

type MockGradingRepo struct {
	mock.Mock
}
func (m *MockGradingRepo) CreateSchema(ctx context.Context, s *models.GradingSchema) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}
func (m *MockGradingRepo) GetSchema(ctx context.Context, id string) (*models.GradingSchema, error) { return nil, nil }
func (m *MockGradingRepo) ListSchemas(ctx context.Context, tenantID string) ([]models.GradingSchema, error) { 
	args := m.Called(ctx, tenantID)
	return args.Get(0).([]models.GradingSchema), args.Error(1)
}
func (m *MockGradingRepo) GetDefaultSchema(ctx context.Context, tenantID string) (*models.GradingSchema, error) { return nil, nil }
func (m *MockGradingRepo) UpdateSchema(ctx context.Context, s *models.GradingSchema) error { return nil }
func (m *MockGradingRepo) DeleteSchema(ctx context.Context, id string) error { return nil }

func (m *MockGradingRepo) CreateEntry(ctx context.Context, e *models.GradebookEntry) error { return nil }
func (m *MockGradingRepo) GetEntry(ctx context.Context, id string) (*models.GradebookEntry, error) { return nil, nil }
func (m *MockGradingRepo) GetEntryByActivity(ctx context.Context, o, a, s string) (*models.GradebookEntry, error) { return nil, nil }
func (m *MockGradingRepo) ListEntries(ctx context.Context, oID string) ([]models.GradebookEntry, error) { return nil, nil }
func (m *MockGradingRepo) ListStudentEntries(ctx context.Context, sID string) ([]models.GradebookEntry, error) { return nil, nil }

func TestGradingHandler_CreateSchema(t *testing.T) {
	mockRepo := new(MockGradingRepo)
	svc := services.NewGradingService(mockRepo)
	h := handlers.NewGradingHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	// Mock middleware setting tenant
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", "t-1")
	})
	r.POST("/grading/schemas", h.CreateSchema)

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("CreateSchema", mock.Anything, mock.MatchedBy(func(s *models.GradingSchema) bool {
			return s.Name == "New Schema" && s.TenantID == "t-1"
		})).Return(nil)

		reqBody := map[string]interface{}{"name": "New Schema"}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/grading/schemas", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

func TestGradingHandler_ListSchemas(t *testing.T) {
	mockRepo := new(MockGradingRepo)
	svc := services.NewGradingService(mockRepo)
	h := handlers.NewGradingHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", "t-1")
	})
	r.GET("/grading/schemas", h.ListSchemas)

	t.Run("Success", func(t *testing.T) {
		expected := []models.GradingSchema{{ID: "s1", Name: "Schema 1"}}
		mockRepo.On("ListSchemas", mock.Anything, "t-1").Return(expected, nil)

		req, _ := http.NewRequest("GET", "/grading/schemas", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Schema 1")
	})
}

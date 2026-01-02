package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type LocalMockGradingRepo struct {
	mock.Mock
}

func (m *LocalMockGradingRepo) CreateSchema(ctx context.Context, s *models.GradingSchema) error {
	return m.Called(ctx, s).Error(0)
}
func (m *LocalMockGradingRepo) GetSchema(ctx context.Context, id string) (*models.GradingSchema, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.GradingSchema), args.Error(1)
}
func (m *LocalMockGradingRepo) ListSchemas(ctx context.Context, tenantID string) ([]models.GradingSchema, error) {
	args := m.Called(ctx, tenantID)
	return args.Get(0).([]models.GradingSchema), args.Error(1)
}
func (m *LocalMockGradingRepo) GetDefaultSchema(ctx context.Context, tenantID string) (*models.GradingSchema, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.GradingSchema), args.Error(1)
}
func (m *LocalMockGradingRepo) UpdateSchema(ctx context.Context, s *models.GradingSchema) error {
	return m.Called(ctx, s).Error(0)
}
func (m *LocalMockGradingRepo) DeleteSchema(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *LocalMockGradingRepo) CreateEntry(ctx context.Context, e *models.GradebookEntry) error {
	return m.Called(ctx, e).Error(0)
}
func (m *LocalMockGradingRepo) GetEntry(ctx context.Context, id string) (*models.GradebookEntry, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.GradebookEntry), args.Error(1)
}
func (m *LocalMockGradingRepo) GetEntryByActivity(ctx context.Context, o, a, s string) (*models.GradebookEntry, error) {
	args := m.Called(ctx, o, a, s)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.GradebookEntry), args.Error(1)
}
func (m *LocalMockGradingRepo) ListEntries(ctx context.Context, o string) ([]models.GradebookEntry, error) {
	args := m.Called(ctx, o)
	return args.Get(0).([]models.GradebookEntry), args.Error(1)
}
func (m *LocalMockGradingRepo) ListStudentEntries(ctx context.Context, s string) ([]models.GradebookEntry, error) {
	args := m.Called(ctx, s)
	return args.Get(0).([]models.GradebookEntry), args.Error(1)
}

func TestGradingHandler_SubmitGrade(t *testing.T) {
	mockRepo := new(LocalMockGradingRepo)
	svc := services.NewGradingService(mockRepo)
	h := NewGradingHandler(svc)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	entry := models.GradebookEntry{
		CourseOfferingID: "off-1",
		ActivityID:       "act-1",
		StudentID:        "stud-1",
		Score:            85,
		MaxScore:         100,
	}
	body, _ := json.Marshal(entry)
	c.Request, _ = http.NewRequest("POST", "/grading/entries", bytes.NewBuffer(body))
	c.Set("userID", "inst-1")
	c.Set("tenant_id", "t1")

	// Expect GetDefaultSchema call
	mockRepo.On("GetDefaultSchema", mock.Anything, "t1").Return(&models.GradingSchema{
		Name:  "Default",
		Scale: []byte(`[{"min": 90, "grade": "A"}, {"min": 0, "grade": "F"}]`),
	}, nil)
	// Expect CreateEntry call
	mockRepo.On("CreateEntry", mock.Anything, mock.AnythingOfType("*models.GradebookEntry")).Return(nil)

	h.SubmitGrade(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestGradingHandler_CreateSchema(t *testing.T) {
	mockRepo := new(LocalMockGradingRepo)
	svc := services.NewGradingService(mockRepo)
	h := NewGradingHandler(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	schema := models.GradingSchema{
		Name: "New Schema",
	}
	body, _ := json.Marshal(schema)
	c.Request, _ = http.NewRequest("POST", "/grading/schemas", bytes.NewBuffer(body))
	c.Set("tenant_id", "t1")

	mockRepo.On("CreateSchema", mock.Anything, mock.MatchedBy(func(s *models.GradingSchema) bool {
		return s.Name == "New Schema" && s.TenantID == "t1"
	})).Return(nil)

	h.CreateSchema(c)
	assert.Equal(t, http.StatusCreated, w.Code)
}

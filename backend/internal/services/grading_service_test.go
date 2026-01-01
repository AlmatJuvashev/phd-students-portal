package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGradingRepo
type MockGradingRepo struct {
	mock.Mock
}
func (m *MockGradingRepo) CreateSchema(ctx context.Context, s *models.GradingSchema) error { return nil }
func (m *MockGradingRepo) GetSchema(ctx context.Context, id string) (*models.GradingSchema, error) { return nil, nil }
func (m *MockGradingRepo) ListSchemas(ctx context.Context, tenantID string) ([]models.GradingSchema, error) { return nil, nil }
func (m *MockGradingRepo) GetDefaultSchema(ctx context.Context, tenantID string) (*models.GradingSchema, error) {
	args := m.Called(ctx, tenantID)
	return args.Get(0).(*models.GradingSchema), args.Error(1)
}
func (m *MockGradingRepo) UpdateSchema(ctx context.Context, s *models.GradingSchema) error { return nil }
func (m *MockGradingRepo) DeleteSchema(ctx context.Context, id string) error { return nil }
func (m *MockGradingRepo) CreateEntry(ctx context.Context, e *models.GradebookEntry) error { 
	args := m.Called(ctx, e)
	return args.Error(0)
}
func (m *MockGradingRepo) GetEntry(ctx context.Context, id string) (*models.GradebookEntry, error) { return nil, nil }
func (m *MockGradingRepo) GetEntryByActivity(ctx context.Context, o, a, s string) (*models.GradebookEntry, error) { return nil, nil }
func (m *MockGradingRepo) ListEntries(ctx context.Context, oID string) ([]models.GradebookEntry, error) { return nil, nil }
func (m *MockGradingRepo) ListStudentEntries(ctx context.Context, sID string) ([]models.GradebookEntry, error) { return nil, nil }


func TestGradingService_SubmitGrade(t *testing.T) {
	mockRepo := new(MockGradingRepo)
	svc := NewGradingService(mockRepo)
	ctx := context.Background()

	// 1. Setup Mock Schema (US Letter)
	rules := []map[string]interface{}{
		{"min": 90, "grade": "A"},
		{"min": 80, "grade": "B"},
		{"min": 70, "grade": "C"},
		{"min": 60, "grade": "D"},
		{"min": 0,  "grade": "F"},
	}
	rulesBytes, _ := json.Marshal(rules)
	schema := &models.GradingSchema{
		ID: "schema-1", 
		Name: "US Letter", 
		Scale: types.JSONText(rulesBytes),
	}

	mockRepo.On("GetDefaultSchema", ctx, "tenant-1").Return(schema, nil)
	mockRepo.On("CreateEntry", ctx, mock.Anything).Return(nil)

	// 2. Test Case: 85/100 -> B
	entry := &models.GradebookEntry{
		CourseOfferingID: "off-1",
		ActivityID: "act-1",
		StudentID: "stud-1",
		Score: 85,
		MaxScore: 100,
	}

	err := svc.SubmitGrade(ctx, entry, "tenant-1")
	assert.NoError(t, err)
	assert.Equal(t, "B", entry.Grade)

	// 3. Test Case: 95/100 -> A
	entryA := &models.GradebookEntry{
		CourseOfferingID: "off-1",
		ActivityID: "act-1",
		StudentID: "stud-1",
		Score: 95,
		MaxScore: 100,
	}
	err = svc.SubmitGrade(ctx, entryA, "tenant-1")
	assert.NoError(t, err)
	assert.Equal(t, "A", entryA.Grade)
	
	// 4. Test Case: 50/100 -> F
	entryF := &models.GradebookEntry{
		CourseOfferingID: "off-1",
		ActivityID: "act-1",
		StudentID: "stud-1",
		Score: 50,
		MaxScore: 100,
	}
	err = svc.SubmitGrade(ctx, entryF, "tenant-1")
	assert.NoError(t, err)
	assert.Equal(t, "F", entryF.Grade)
}

func TestGradingService_CreateSchema(t *testing.T) {
	mockRepo := new(MockGradingRepo)
	svc := NewGradingService(mockRepo)
	ctx := context.Background()

	// 1. Valid Input
	mockRepo.On("CreateSchema", ctx, mock.MatchedBy(func(s *models.GradingSchema) bool {
		return s.Name == "KPI 100"
	})).Return(nil)

	err := svc.CreateSchema(ctx, &models.GradingSchema{Name: "KPI 100", TenantID: "t1"})
	assert.NoError(t, err)

	// 2. Missing Name (Service validation)
	err = svc.CreateSchema(ctx, &models.GradingSchema{TenantID: "t1"}) // Name empty
	assert.Error(t, err)
	assert.Equal(t, "name is required", err.Error())
}

package services

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCurriculumRepo for service testing
type MockCurriculumRepo struct {
	mock.Mock
}

func (m *MockCurriculumRepo) CreateProgram(ctx context.Context, p *models.Program) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}
func (m *MockCurriculumRepo) GetProgram(ctx context.Context, id string) (*models.Program, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Program), args.Error(1)
}
func (m *MockCurriculumRepo) ListPrograms(ctx context.Context, tenantID string) ([]models.Program, error) {
	args := m.Called(ctx, tenantID)
	return args.Get(0).([]models.Program), args.Error(1)
}
func (m *MockCurriculumRepo) UpdateProgram(ctx context.Context, p *models.Program) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}
func (m *MockCurriculumRepo) DeleteProgram(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
// Course mocks
func (m *MockCurriculumRepo) CreateCourse(ctx context.Context, c *models.Course) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}
func (m *MockCurriculumRepo) GetCourse(ctx context.Context, id string) (*models.Course, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Course), args.Error(1)
}
func (m *MockCurriculumRepo) ListCourses(ctx context.Context, tenantID string, programID *string) ([]models.Course, error) {
	args := m.Called(ctx, tenantID, programID)
	return args.Get(0).([]models.Course), args.Error(1)
}
func (m *MockCurriculumRepo) UpdateCourse(ctx context.Context, c *models.Course) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}
func (m *MockCurriculumRepo) DeleteCourse(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
// Journey mocks
func (m *MockCurriculumRepo) CreateJourneyMap(ctx context.Context, jm *models.JourneyMap) error {
	args := m.Called(ctx, jm)
	return args.Error(0)
}
func (m *MockCurriculumRepo) GetJourneyMapByProgram(ctx context.Context, programID string) (*models.JourneyMap, error) {
	args := m.Called(ctx, programID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.JourneyMap), args.Error(1)
}
func (m *MockCurriculumRepo) CreateNodeDefinition(ctx context.Context, nd *models.JourneyNodeDefinition) error {
	args := m.Called(ctx, nd)
	return args.Error(0)
}
func (m *MockCurriculumRepo) GetNodeDefinitions(ctx context.Context, jMapID string) ([]models.JourneyNodeDefinition, error) {
	args := m.Called(ctx, jMapID)
	return args.Get(0).([]models.JourneyNodeDefinition), args.Error(1)
}
func (m *MockCurriculumRepo) DeleteNodeDefinition(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockCurriculumRepo) CreateCohort(ctx context.Context, c *models.Cohort) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}
func (m *MockCurriculumRepo) ListCohorts(ctx context.Context, pID string) ([]models.Cohort, error) {
	args := m.Called(ctx, pID)
	return args.Get(0).([]models.Cohort), args.Error(1)
}

// Ensure mock implements interface
var _ repository.CurriculumRepository = (*MockCurriculumRepo)(nil)

func TestCurriculumService_CreateProgram(t *testing.T) {
	mockRepo := new(MockCurriculumRepo)
	svc := NewCurriculumService(mockRepo)
	ctx := context.Background()

	// 1. Success
	p := &models.Program{
		TenantID: "t1",
		Code:     "CS101",
		Title:    "Comp Sci",
	}
	mockRepo.On("CreateProgram", ctx, p).Return(nil)
	err := svc.CreateProgram(ctx, p)
	assert.NoError(t, err)
	assert.NotZero(t, p.CreatedAt)

	// 2. Validation Error
	pInvalid := &models.Program{TenantID: "t1"} // Missing code/title
	err = svc.CreateProgram(ctx, pInvalid)
	assert.Error(t, err)
	assert.Equal(t, "title is required", err.Error())
}

func TestCurriculumService_ListPrograms(t *testing.T) {
	mockRepo := new(MockCurriculumRepo)
	svc := NewCurriculumService(mockRepo)
	ctx := context.Background()

	mockRepo.On("ListPrograms", ctx, "t1").Return([]models.Program{{ID: "p1"}}, nil)
	list, err := svc.ListPrograms(ctx, "t1")
	assert.NoError(t, err)
	assert.Len(t, list, 1)
}

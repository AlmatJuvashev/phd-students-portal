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
	if args.Get(0) == nil { return nil, args.Error(1) }
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
func (m *MockCurriculumRepo) UpdateJourneyMap(ctx context.Context, jm *models.JourneyMap) error {
	args := m.Called(ctx, jm)
	return args.Error(0)
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
func (m *MockCurriculumRepo) GetNodeDefinition(ctx context.Context, id string) (*models.JourneyNodeDefinition, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.JourneyNodeDefinition), args.Error(1)
}
func (m *MockCurriculumRepo) UpdateNodeDefinition(ctx context.Context, nd *models.JourneyNodeDefinition) error {
	args := m.Called(ctx, nd)
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
func (m *MockCurriculumRepo) SetCourseRequirement(ctx context.Context, req *models.CourseRequirement) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}
func (m *MockCurriculumRepo) GetCourseRequirements(ctx context.Context, courseID string) ([]models.CourseRequirement, error) {
	args := m.Called(ctx, courseID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.CourseRequirement), args.Error(1)
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
	// 3. Success with default name
	p2 := &models.Program{TenantID: "t1", Code: "P2", Title: "Title 2"}
	mockRepo.On("CreateProgram", ctx, p2).Return(nil)
	err = svc.CreateProgram(ctx, p2)
	assert.NoError(t, err)
	assert.Equal(t, "P2", p2.Name)
}

func TestCurriculumService_Programs(t *testing.T) {
	mockRepo := new(MockCurriculumRepo)
	svc := NewCurriculumService(mockRepo)
	ctx := context.Background()

	t.Run("GetProgram", func(t *testing.T) {
		expected := &models.Program{ID: "p1"}
		mockRepo.On("GetProgram", ctx, "p1").Return(expected, nil)
		res, err := svc.GetProgram(ctx, "p1")
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("ListPrograms", func(t *testing.T) {
		expected := []models.Program{{ID: "p1"}}
		mockRepo.On("ListPrograms", ctx, "t1").Return(expected, nil)
		res, err := svc.ListPrograms(ctx, "t1")
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("UpdateProgram", func(t *testing.T) {
		p := &models.Program{ID: "p1"}
		mockRepo.On("UpdateProgram", ctx, p).Return(nil)
		err := svc.UpdateProgram(ctx, p)
		assert.NoError(t, err)
		assert.NotZero(t, p.UpdatedAt)
	})

	t.Run("DeleteProgram", func(t *testing.T) {
		mockRepo.On("DeleteProgram", ctx, "p1").Return(nil)
		err := svc.DeleteProgram(ctx, "p1")
		assert.NoError(t, err)
	})
}

func TestCurriculumService_Courses(t *testing.T) {
	mockRepo := new(MockCurriculumRepo)
	svc := NewCurriculumService(mockRepo)
	ctx := context.Background()

	t.Run("CreateCourse Success", func(t *testing.T) {
		c := &models.Course{TenantID: "t1", Title: "C1"}
		mockRepo.On("CreateCourse", ctx, c).Return(nil)
		err := svc.CreateCourse(ctx, c)
		assert.NoError(t, err)
		assert.NotZero(t, c.CreatedAt)
	})

	t.Run("CreateCourse Validation", func(t *testing.T) {
		c := &models.Course{TenantID: "t1"} // missing title
		err := svc.CreateCourse(ctx, c)
		assert.Error(t, err)
		assert.Equal(t, "title is required", err.Error())
	})

	t.Run("GetCourse", func(t *testing.T) {
		expected := &models.Course{ID: "c1"}
		mockRepo.On("GetCourse", ctx, "c1").Return(expected, nil)
		res, err := svc.GetCourse(ctx, "c1")
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("ListCourses", func(t *testing.T) {
		expected := []models.Course{{ID: "c1"}}
		progID := "p1"
		mockRepo.On("ListCourses", ctx, "t1", &progID).Return(expected, nil)
		res, err := svc.ListCourses(ctx, "t1", &progID)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("UpdateCourse", func(t *testing.T) {
		c := &models.Course{ID: "c1"}
		mockRepo.On("UpdateCourse", ctx, c).Return(nil)
		err := svc.UpdateCourse(ctx, c)
		assert.NoError(t, err)
		assert.NotZero(t, c.UpdatedAt)
	})

	t.Run("DeleteCourse", func(t *testing.T) {
		mockRepo.On("DeleteCourse", ctx, "c1").Return(nil)
		err := svc.DeleteCourse(ctx, "c1")
		assert.NoError(t, err)
	})
}

func TestCurriculumService_Journey(t *testing.T) {
	mockRepo := new(MockCurriculumRepo)
	svc := NewCurriculumService(mockRepo)
	ctx := context.Background()

	t.Run("CreateJourneyMap", func(t *testing.T) {
		jm := &models.JourneyMap{ProgramID: "p1"}
		mockRepo.On("CreateJourneyMap", ctx, jm).Return(nil)
		err := svc.CreateJourneyMap(ctx, jm)
		assert.NoError(t, err)
		assert.Equal(t, "{}", jm.Config)
		assert.NotZero(t, jm.CreatedAt)
	})

	t.Run("GetJourneyMap", func(t *testing.T) {
		expected := &models.JourneyMap{ID: "jm1"}
		mockRepo.On("GetJourneyMapByProgram", ctx, "p1").Return(expected, nil)
		res, err := svc.GetJourneyMap(ctx, "p1")
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("CreateNodeDefinition", func(t *testing.T) {
		nd := &models.JourneyNodeDefinition{Slug: "n1"}
		mockRepo.On("CreateNodeDefinition", ctx, nd).Return(nil)
		err := svc.CreateNodeDefinition(ctx, nd)
		assert.NoError(t, err)
		assert.NotZero(t, nd.CreatedAt)
	})

	t.Run("GetNodeDefinitions", func(t *testing.T) {
		expected := []models.JourneyNodeDefinition{{ID: "n1"}}
		mockRepo.On("GetNodeDefinitions", ctx, "jm1").Return(expected, nil)
		res, err := svc.GetNodeDefinitions(ctx, "jm1")
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}

func TestCurriculumService_Cohorts(t *testing.T) {
	mockRepo := new(MockCurriculumRepo)
	svc := NewCurriculumService(mockRepo)
	ctx := context.Background()

	t.Run("CreateCohort", func(t *testing.T) {
		c := &models.Cohort{Name: "2024"}
		mockRepo.On("CreateCohort", ctx, c).Return(nil)
		err := svc.CreateCohort(ctx, c)
		assert.NoError(t, err)
		assert.NotZero(t, c.CreatedAt)
	})

	t.Run("ListCohorts", func(t *testing.T) {
		expected := []models.Cohort{{ID: "ch1"}}
		mockRepo.On("ListCohorts", ctx, "p1").Return(expected, nil)
		res, err := svc.ListCohorts(ctx, "p1")
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}

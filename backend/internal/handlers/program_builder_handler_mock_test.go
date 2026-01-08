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

type mockCurriculumRepoForBuilder struct {
	mock.Mock
}

func (m *mockCurriculumRepoForBuilder) CreateProgram(ctx context.Context, p *models.Program) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}
func (m *mockCurriculumRepoForBuilder) GetProgram(ctx context.Context, id string) (*models.Program, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Program), args.Error(1)
}
func (m *mockCurriculumRepoForBuilder) ListPrograms(ctx context.Context, tenantID string) ([]models.Program, error) {
	args := m.Called(ctx, tenantID)
	return args.Get(0).([]models.Program), args.Error(1)
}
func (m *mockCurriculumRepoForBuilder) UpdateProgram(ctx context.Context, p *models.Program) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}
func (m *mockCurriculumRepoForBuilder) DeleteProgram(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockCurriculumRepoForBuilder) CreateCourse(ctx context.Context, c *models.Course) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}
func (m *mockCurriculumRepoForBuilder) GetCourse(ctx context.Context, id string) (*models.Course, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Course), args.Error(1)
}
func (m *mockCurriculumRepoForBuilder) ListCourses(ctx context.Context, tenantID string, programID *string) ([]models.Course, error) {
	args := m.Called(ctx, tenantID, programID)
	return args.Get(0).([]models.Course), args.Error(1)
}
func (m *mockCurriculumRepoForBuilder) UpdateCourse(ctx context.Context, c *models.Course) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}
func (m *mockCurriculumRepoForBuilder) DeleteCourse(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockCurriculumRepoForBuilder) CreateJourneyMap(ctx context.Context, jm *models.JourneyMap) error {
	args := m.Called(ctx, jm)
	return args.Error(0)
}
func (m *mockCurriculumRepoForBuilder) GetJourneyMapByProgram(ctx context.Context, programID string) (*models.JourneyMap, error) {
	args := m.Called(ctx, programID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.JourneyMap), args.Error(1)
}
func (m *mockCurriculumRepoForBuilder) UpdateJourneyMap(ctx context.Context, jm *models.JourneyMap) error {
	args := m.Called(ctx, jm)
	return args.Error(0)
}

func (m *mockCurriculumRepoForBuilder) CreateNodeDefinition(ctx context.Context, nd *models.JourneyNodeDefinition) error {
	args := m.Called(ctx, nd)
	return args.Error(0)
}
func (m *mockCurriculumRepoForBuilder) GetNodeDefinitions(ctx context.Context, journeyMapID string) ([]models.JourneyNodeDefinition, error) {
	args := m.Called(ctx, journeyMapID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.JourneyNodeDefinition), args.Error(1)
}
func (m *mockCurriculumRepoForBuilder) GetNodeDefinition(ctx context.Context, id string) (*models.JourneyNodeDefinition, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.JourneyNodeDefinition), args.Error(1)
}
func (m *mockCurriculumRepoForBuilder) UpdateNodeDefinition(ctx context.Context, nd *models.JourneyNodeDefinition) error {
	args := m.Called(ctx, nd)
	return args.Error(0)
}
func (m *mockCurriculumRepoForBuilder) DeleteNodeDefinition(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockCurriculumRepoForBuilder) CreateCohort(ctx context.Context, c *models.Cohort) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}
func (m *mockCurriculumRepoForBuilder) ListCohorts(ctx context.Context, programID string) ([]models.Cohort, error) {
	args := m.Called(ctx, programID)
	return args.Get(0).([]models.Cohort), args.Error(1)
}
func (m *mockCurriculumRepoForBuilder) GetCohort(ctx context.Context, id string) (*models.Cohort, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Cohort), args.Error(1)
}
func (m *mockCurriculumRepoForBuilder) UpdateCohort(ctx context.Context, c *models.Cohort) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}
func (m *mockCurriculumRepoForBuilder) DeleteCohort(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockCurriculumRepoForBuilder) SetCourseRequirement(ctx context.Context, req *models.CourseRequirement) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}
func (m *mockCurriculumRepoForBuilder) GetCourseRequirements(ctx context.Context, courseID string) ([]models.CourseRequirement, error) {
	args := m.Called(ctx, courseID)
	return args.Get(0).([]models.CourseRequirement), args.Error(1)
}

func TestProgramBuilderHandler_JourneyMap(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mockCurriculumRepoForBuilder)
	svc := services.NewProgramBuilderService(mockRepo)
	h := NewProgramBuilderHandler(svc)

	t.Run("GetJourneyMap", func(t *testing.T) {
		mockRepo.On("GetJourneyMapByProgram", mock.Anything, "p1").Return(&models.JourneyMap{ID: "m1", ProgramID: "p1"}, nil)
		mockRepo.On("GetNodeDefinitions", mock.Anything, "m1").Return([]models.JourneyNodeDefinition{}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "p1"}}
		c.Request = httptest.NewRequest("GET", "/programs/p1/builder/map", nil)

		h.GetJourneyMap(c)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "m1")
	})

	t.Run("UpdateJourneyMap", func(t *testing.T) {
		mockRepo.On("GetJourneyMapByProgram", mock.Anything, "p1").Return(&models.JourneyMap{ID: "m1", ProgramID: "p1"}, nil)
		mockRepo.On("GetNodeDefinitions", mock.Anything, "m1").Return([]models.JourneyNodeDefinition{}, nil)
		mockRepo.On("UpdateJourneyMap", mock.Anything, mock.Anything).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "p1"}}
		body, _ := json.Marshal(map[string]any{"title": "New Title"})
		c.Request = httptest.NewRequest("PUT", "/programs/p1/builder/map", bytes.NewBuffer(body))

		h.UpdateJourneyMap(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestProgramBuilderHandler_Nodes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mockCurriculumRepoForBuilder)
	svc := services.NewProgramBuilderService(mockRepo)
	h := NewProgramBuilderHandler(svc)

	t.Run("UpdateNode", func(t *testing.T) {
		mockRepo.On("GetNodeDefinition", mock.Anything, "n1").Return(&models.JourneyNodeDefinition{ID: "n1", Slug: "old"}, nil)
		mockRepo.On("UpdateNodeDefinition", mock.Anything, mock.Anything).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "nodeId", Value: "n1"}}
		body, _ := json.Marshal(map[string]any{"slug": "new-slug"})
		c.Request = httptest.NewRequest("PUT", "/builder/nodes/n1", bytes.NewBuffer(body))

	h.UpdateNode(c)
	assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestProgramBuilderHandler_ErrorPaths(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("GetJourneyMap_ServiceError", func(t *testing.T) {
		mockRepo := new(mockCurriculumRepoForBuilder)
		svc := services.NewProgramBuilderService(mockRepo)
		h := NewProgramBuilderHandler(svc)
		
		mockRepo.On("GetJourneyMapByProgram", mock.Anything, "p1").Return(nil, assert.AnError)
		
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "p1"}}
		c.Request = httptest.NewRequest("GET", "/programs/p1/builder/map", nil)
		
		h.GetJourneyMap(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("UpdateJourneyMap_ServiceError", func(t *testing.T) {
		mockRepo := new(mockCurriculumRepoForBuilder)
		svc := services.NewProgramBuilderService(mockRepo)
		h := NewProgramBuilderHandler(svc)

		mockRepo.On("GetJourneyMapByProgram", mock.Anything, "p1").Return(&models.JourneyMap{ID: "m1", ProgramID: "p1"}, nil)
		mockRepo.On("GetNodeDefinitions", mock.Anything, "m1").Return([]models.JourneyNodeDefinition{}, nil)
		mockRepo.On("UpdateJourneyMap", mock.Anything, mock.Anything).Return(assert.AnError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "p1"}}
		body, _ := json.Marshal(map[string]any{"title": "Fail"})
		c.Request = httptest.NewRequest("PUT", "/programs/p1/builder/map", bytes.NewBuffer(body))

		h.UpdateJourneyMap(c)
		assert.Equal(t, http.StatusBadRequest, w.Code) 
	})

	t.Run("GetNodes_ServiceError", func(t *testing.T) {
		mockRepo := new(mockCurriculumRepoForBuilder)
		svc := services.NewProgramBuilderService(mockRepo)
		h := NewProgramBuilderHandler(svc)

		mockRepo.On("GetJourneyMapByProgram", mock.Anything, "p1").Return(&models.JourneyMap{ID: "m1"}, nil)
		mockRepo.On("GetNodeDefinitions", mock.Anything, "m1").Return(nil, assert.AnError) // Fetch fail

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "p1"}}
		c.Request = httptest.NewRequest("GET", "/programs/p1/builder/nodes", nil)

		h.GetNodes(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("CreateNode_ServiceError", func(t *testing.T) {
		mockRepo := new(mockCurriculumRepoForBuilder)
		svc := services.NewProgramBuilderService(mockRepo)
		h := NewProgramBuilderHandler(svc)

		mockRepo.On("GetJourneyMapByProgram", mock.Anything, "p1").Return(&models.JourneyMap{ID: "m1"}, nil)
		mockRepo.On("CreateNodeDefinition", mock.Anything, mock.Anything).Return(assert.AnError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "p1"}}
		body, _ := json.Marshal(map[string]any{"slug": "n1", "type": "step"})
		c.Request = httptest.NewRequest("POST", "/programs/p1/builder/nodes", bytes.NewBuffer(body))

		h.CreateNode(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("UpdateNode_ServiceError", func(t *testing.T) {
		mockRepo := new(mockCurriculumRepoForBuilder)
		svc := services.NewProgramBuilderService(mockRepo)
		h := NewProgramBuilderHandler(svc)

		mockRepo.On("GetNodeDefinition", mock.Anything, "n1").Return(nil, assert.AnError)
		
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "nodeId", Value: "n1"}}
		body, _ := json.Marshal(map[string]any{"slug": "n1-ud"})
		c.Request = httptest.NewRequest("PUT", "/programs/p1/builder/nodes/n1", bytes.NewBuffer(body))

		h.UpdateNode(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

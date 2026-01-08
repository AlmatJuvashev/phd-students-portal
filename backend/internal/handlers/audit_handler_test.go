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

type mockAuditRepo struct {
	mock.Mock
}

func (m *mockAuditRepo) ListLearningOutcomes(ctx context.Context, tenantID string, programID, courseID *string) ([]models.LearningOutcome, error) {
	args := m.Called(ctx, tenantID, programID, courseID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.LearningOutcome), args.Error(1)
}
func (m *mockAuditRepo) GetLearningOutcome(ctx context.Context, id string) (*models.LearningOutcome, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.LearningOutcome), args.Error(1)
}
func (m *mockAuditRepo) CreateLearningOutcome(ctx context.Context, outcome *models.LearningOutcome) error {
	args := m.Called(ctx, outcome)
	return args.Error(0)
}
func (m *mockAuditRepo) UpdateLearningOutcome(ctx context.Context, outcome *models.LearningOutcome) error {
	args := m.Called(ctx, outcome)
	return args.Error(0)
}
func (m *mockAuditRepo) DeleteLearningOutcome(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *mockAuditRepo) LinkOutcomeToAssessment(ctx context.Context, outcomeID, nodeDefID string, weight float64) error {
	args := m.Called(ctx, outcomeID, nodeDefID, weight)
	return args.Error(0)
}
func (m *mockAuditRepo) GetOutcomeAssessments(ctx context.Context, outcomeID string) ([]models.OutcomeAssessment, error) {
	args := m.Called(ctx, outcomeID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.OutcomeAssessment), args.Error(1)
}
func (m *mockAuditRepo) LogCurriculumChange(ctx context.Context, log *models.CurriculumChangeLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}
func (m *mockAuditRepo) ListCurriculumChanges(ctx context.Context, filter models.AuditReportFilter) ([]models.CurriculumChangeLog, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.CurriculumChangeLog), args.Error(1)
}

type mockCurriculumRepo struct {
	mock.Mock
}

func (m *mockCurriculumRepo) CreateProgram(ctx context.Context, p *models.Program) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}
func (m *mockCurriculumRepo) GetProgram(ctx context.Context, id string) (*models.Program, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Program), args.Error(1)
}
func (m *mockCurriculumRepo) ListPrograms(ctx context.Context, tenantID string) ([]models.Program, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.Program), args.Error(1)
}
func (m *mockCurriculumRepo) UpdateProgram(ctx context.Context, p *models.Program) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}
func (m *mockCurriculumRepo) DeleteProgram(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *mockCurriculumRepo) CreateCourse(ctx context.Context, c *models.Course) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}
func (m *mockCurriculumRepo) GetCourse(ctx context.Context, id string) (*models.Course, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Course), args.Error(1)
}
func (m *mockCurriculumRepo) ListCourses(ctx context.Context, tenantID string, programID *string) ([]models.Course, error) {
	args := m.Called(ctx, tenantID, programID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.Course), args.Error(1)
}
func (m *mockCurriculumRepo) UpdateCourse(ctx context.Context, c *models.Course) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}
func (m *mockCurriculumRepo) DeleteCourse(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *mockCurriculumRepo) CreateJourneyMap(ctx context.Context, jm *models.JourneyMap) error {
	args := m.Called(ctx, jm)
	return args.Error(0)
}
func (m *mockCurriculumRepo) GetJourneyMapByProgram(ctx context.Context, programID string) (*models.JourneyMap, error) {
	args := m.Called(ctx, programID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.JourneyMap), args.Error(1)
}
func (m *mockCurriculumRepo) UpdateJourneyMap(ctx context.Context, jm *models.JourneyMap) error {
	args := m.Called(ctx, jm)
	return args.Error(0)
}
func (m *mockCurriculumRepo) CreateNodeDefinition(ctx context.Context, nd *models.JourneyNodeDefinition) error {
	args := m.Called(ctx, nd)
	return args.Error(0)
}
func (m *mockCurriculumRepo) GetNodeDefinitions(ctx context.Context, journeyMapID string) ([]models.JourneyNodeDefinition, error) {
	args := m.Called(ctx, journeyMapID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.JourneyNodeDefinition), args.Error(1)
}
func (m *mockCurriculumRepo) GetNodeDefinition(ctx context.Context, id string) (*models.JourneyNodeDefinition, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.JourneyNodeDefinition), args.Error(1)
}
func (m *mockCurriculumRepo) UpdateNodeDefinition(ctx context.Context, nd *models.JourneyNodeDefinition) error {
	args := m.Called(ctx, nd)
	return args.Error(0)
}
func (m *mockCurriculumRepo) DeleteNodeDefinition(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *mockCurriculumRepo) CreateCohort(ctx context.Context, c *models.Cohort) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}
func (m *mockCurriculumRepo) ListCohorts(ctx context.Context, programID string) ([]models.Cohort, error) {
	args := m.Called(ctx, programID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.Cohort), args.Error(1)
}
func (m *mockCurriculumRepo) SetCourseRequirement(ctx context.Context, req *models.CourseRequirement) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}
func (m *mockCurriculumRepo) GetCourseRequirements(ctx context.Context, courseID string) ([]models.CourseRequirement, error) {
	args := m.Called(ctx, courseID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.CourseRequirement), args.Error(1)
}

func TestAuditHandler_ListPrograms(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(mockAuditRepo)
	curRepo := new(mockCurriculumRepo)
	svc := services.NewAuditService(repo, curRepo)
	curSvc := services.NewCurriculumService(curRepo)
	h := NewAuditHandler(svc, curSvc)

	curRepo.On("ListPrograms", mock.Anything, "t1").Return([]models.Program{{ID: "p1"}}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/audit/programs", nil)
	c.Set("tenant_id", "t1")

	h.ListPrograms(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuditHandler_ProgramSummaryReport(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(mockAuditRepo)
	curRepo := new(mockCurriculumRepo)
	svc := services.NewAuditService(repo, curRepo)
	curSvc := services.NewCurriculumService(curRepo)
	h := NewAuditHandler(svc, curSvc)

	curRepo.On("GetProgram", mock.Anything, "p1").Return(&models.Program{ID: "p1"}, nil)
	curRepo.On("ListCourses", mock.Anything, "t1", mock.Anything).Return([]models.Course{}, nil)
	repo.On("ListLearningOutcomes", mock.Anything, "t1", mock.Anything, mock.Anything).Return([]models.LearningOutcome{}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/audit/report?program_id=p1", nil)
	c.Set("tenant_id", "t1")

	h.ProgramSummaryReport(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuditHandler_GetProgram(t *testing.T) {
	gin.SetMode(gin.TestMode)
	curRepo := new(mockCurriculumRepo)
	svc := services.NewAuditService(nil, curRepo)
	curSvc := services.NewCurriculumService(curRepo)
	h := NewAuditHandler(svc, curSvc)

	curRepo.On("GetProgram", mock.Anything, "p1").Return(&models.Program{ID: "p1"}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "p1"}}
	c.Request, _ = http.NewRequest("GET", "/audit/programs/p1", nil)

	h.GetProgram(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuditHandler_ListCourses(t *testing.T) {
	gin.SetMode(gin.TestMode)
	curRepo := new(mockCurriculumRepo)
	svc := services.NewAuditService(nil, curRepo)
	curSvc := services.NewCurriculumService(curRepo)
	h := NewAuditHandler(svc, curSvc)

	curRepo.On("ListCourses", mock.Anything, "t1", mock.Anything).Return([]models.Course{{ID: "c1"}}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/audit/courses", nil)
	c.Set("tenant_id", "t1")

	h.ListCourses(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuditHandler_GetCourse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	curRepo := new(mockCurriculumRepo)
	svc := services.NewAuditService(nil, curRepo)
	curSvc := services.NewCurriculumService(curRepo)
	h := NewAuditHandler(svc, curSvc)

	curRepo.On("GetCourse", mock.Anything, "c1").Return(&models.Course{ID: "c1"}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "c1"}}
	c.Request, _ = http.NewRequest("GET", "/audit/courses/c1", nil)

	h.GetCourse(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuditHandler_ListOutcomes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(mockAuditRepo)
	h := NewAuditHandler(services.NewAuditService(repo, nil), nil)

	repo.On("ListLearningOutcomes", mock.Anything, "t1", mock.Anything, mock.Anything).Return([]models.LearningOutcome{{ID: "lo1"}}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/audit/outcomes", nil)
	c.Set("tenant_id", "t1")

	h.ListOutcomes(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuditHandler_ListChangeLog(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(mockAuditRepo)
	h := NewAuditHandler(services.NewAuditService(repo, nil), nil)

	repo.On("ListCurriculumChanges", mock.Anything, mock.Anything).Return([]models.CurriculumChangeLog{{ID: "l1"}}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/audit/changelog", nil)

	h.ListChangeLog(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuditHandler_CreateOutcome(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(mockAuditRepo)
	h := NewAuditHandler(services.NewAuditService(repo, nil), nil)

	repo.On("CreateLearningOutcome", mock.Anything, mock.Anything).Return(nil)
	repo.On("LogCurriculumChange", mock.Anything, mock.Anything).Return(nil)

	body, _ := json.Marshal(models.LearningOutcome{Code: "PLO1", Title: "Outcome 1"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/audit/outcomes", bytes.NewBuffer(body))
	c.Set("tenant_id", "t1")
	c.Set("userID", "u1")

	h.CreateOutcome(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestAuditHandler_UpdateOutcome(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(mockAuditRepo)
	h := NewAuditHandler(services.NewAuditService(repo, nil), nil)

	repo.On("GetLearningOutcome", mock.Anything, "lo1").Return(&models.LearningOutcome{ID: "lo1"}, nil)
	repo.On("UpdateLearningOutcome", mock.Anything, mock.Anything).Return(nil)
	repo.On("LogCurriculumChange", mock.Anything, mock.Anything).Return(nil)

	body, _ := json.Marshal(models.LearningOutcome{Code: "PLO1_v2", Title: "Outcome 1 v2"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "lo1"}}
	c.Request, _ = http.NewRequest("PUT", "/audit/outcomes/lo1", bytes.NewBuffer(body))
	c.Set("tenant_id", "t1")
	c.Set("userID", "u1")

	h.UpdateOutcome(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuditHandler_DeleteOutcome(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(mockAuditRepo)
	h := NewAuditHandler(services.NewAuditService(repo, nil), nil)

	repo.On("GetLearningOutcome", mock.Anything, "lo1").Return(&models.LearningOutcome{ID: "lo1"}, nil)
	repo.On("DeleteLearningOutcome", mock.Anything, "lo1").Return(nil)
	repo.On("LogCurriculumChange", mock.Anything, mock.Anything).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "lo1"}}
	c.Request, _ = http.NewRequest("DELETE", "/audit/outcomes/lo1", nil)
	c.Set("tenant_id", "t1")
	c.Set("userID", "u1")

	h.DeleteOutcome(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

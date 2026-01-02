package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/scheduler/solver"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// -- Mock Orchestrator --
type MockSchedulerService struct {
	mock.Mock
}

func (m *MockSchedulerService) CreateTerm(ctx context.Context, term *models.AcademicTerm) error {
	args := m.Called(ctx, term)
	return args.Error(0)
}
func (m *MockSchedulerService) ListTerms(ctx context.Context, tenantID string) ([]models.AcademicTerm, error) {
	args := m.Called(ctx, tenantID)
	return args.Get(0).([]models.AcademicTerm), args.Error(1)
}
func (m *MockSchedulerService) GetTerm(ctx context.Context, id string) (*models.AcademicTerm, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.AcademicTerm), args.Error(1)
}
func (m *MockSchedulerService) CreateOffering(ctx context.Context, offering *models.CourseOffering) error {
	args := m.Called(ctx, offering)
	return args.Error(0)
}
func (m *MockSchedulerService) ScheduleSession(ctx context.Context, session *models.ClassSession) ([]string, error) {
	args := m.Called(ctx, session)
	return args.Get(0).([]string), args.Error(1)
}
func (m *MockSchedulerService) ListSessions(ctx context.Context, offeringID string, start, end time.Time) ([]models.ClassSession, error) {
	args := m.Called(ctx, offeringID, start, end)
	return args.Get(0).([]models.ClassSession), args.Error(1)
}
func (m *MockSchedulerService) AutoSchedule(ctx context.Context, tenantID, termID string, config *solver.SolverConfig) (*solver.Solution, error) {
	args := m.Called(ctx, tenantID, termID, config)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*solver.Solution), args.Error(1)
}

// Ensure Mock implements Interface
var _ SchedulerOrchestrator = (*MockSchedulerService)(nil)

func TestSchedulerHandler_CreateTerm(t *testing.T) {
	mockSvc := new(MockSchedulerService)
	h := NewSchedulerHandler(mockSvc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	term := models.AcademicTerm{Name: "Fall 2026", Code: "2026-FA", StartDate: time.Now(), EndDate: time.Now().AddDate(0,4,0)}
	body, _ := json.Marshal(term)
	c.Request, _ = http.NewRequest("POST", "/terms", bytes.NewBuffer(body))
    c.Set("tenant_id", "t1")

	mockSvc.On("CreateTerm", mock.Anything, mock.MatchedBy(func(a *models.AcademicTerm) bool {
		return a.Name == "Fall 2026" && a.Code == "2026-FA" && a.TenantID == "t1"
	})).Return(nil)

	h.CreateTerm(c)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestSchedulerHandler_ListTerms(t *testing.T) {
	mockSvc := new(MockSchedulerService)
	h := NewSchedulerHandler(mockSvc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/terms", nil)
	c.Set("tenant_id", "t1")

	expected := []models.AcademicTerm{{Name: "Spring 2025"}}
	mockSvc.On("ListTerms", mock.Anything, "t1").Return(expected, nil)

	h.ListTerms(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSchedulerHandler_CreateOffering(t *testing.T) {
	mockSvc := new(MockSchedulerService)
	h := NewSchedulerHandler(mockSvc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	offering := models.CourseOffering{TermID: "term-1", CourseID: "course-1"}
	body, _ := json.Marshal(offering)
	c.Request, _ = http.NewRequest("POST", "/offerings", bytes.NewBuffer(body))

	mockSvc.On("CreateOffering", mock.Anything, mock.MatchedBy(func(o *models.CourseOffering) bool {
		return o.TermID == "term-1"
	})).Return(nil)

	h.CreateOffering(c)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestSchedulerHandler_CreateSession(t *testing.T) {
	mockSvc := new(MockSchedulerService)
	h := NewSchedulerHandler(mockSvc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	session := models.ClassSession{
		CourseOfferingID: "off-1", 
		Date:             time.Now(),
		StartTime:        "14:00",
		EndTime:          "15:00",
	}
	body, _ := json.Marshal(session)
	c.Request, _ = http.NewRequest("POST", "/sessions", bytes.NewBuffer(body))
	
	mockSvc.On("ScheduleSession", mock.Anything, mock.AnythingOfType("*models.ClassSession")).Return([]string{}, nil)

	h.CreateSession(c)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestSchedulerHandler_ListSessions(t *testing.T) {
	mockSvc := new(MockSchedulerService)
	h := NewSchedulerHandler(mockSvc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	c.Request, _ = http.NewRequest("GET", "/sessions?offering_id=off-1&start=2025-01-01T00:00:00Z&end=2025-02-01T00:00:00Z", nil)
	
	mockSvc.On("ListSessions", mock.Anything, "off-1", mock.Anything, mock.Anything).Return([]models.ClassSession{{ID: "sess-1"}}, nil)

	h.ListSessions(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSchedulerHandler_AutoSchedule(t *testing.T) {
	mockSvc := new(MockSchedulerService)
	h := NewSchedulerHandler(mockSvc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// User Request
	reqBody := map[string]interface{}{
		"term_id": "term-1",
		"config": map[string]interface{}{
			"max_iterations": 100,
		},
	}
	body, _ := json.Marshal(reqBody)
	c.Request, _ = http.NewRequest("POST", "/autoschedule", bytes.NewBuffer(body))
	c.Set("tenant_id", "t1") // Middleware

	// Expected Solution
	sol := &solver.Solution{
		Score: 100.0,
		IsValid: true,
	}

	mockSvc.On("AutoSchedule", mock.Anything, "t1", "term-1", mock.MatchedBy(func(cfg *solver.SolverConfig) bool {
		return cfg.MaxIterations == 100
	})).Return(sol, nil)

	h.AutoSchedule(c)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var resp solver.Solution
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 100.0, resp.Score)
	assert.True(t, resp.IsValid)
	
	// Error Case
	mockSvc.On("AutoSchedule", mock.Anything, "t1", "term-error", mock.Anything).Return(nil, assert.AnError)
	
	reqErrBody := map[string]interface{}{"term_id": "term-error"}
	bodyErr, _ := json.Marshal(reqErrBody)
	cErr := &gin.Context{Request: httptest.NewRequest("POST", "/autoschedule", bytes.NewBuffer(bodyErr))}
	cErr.Set("tenant_id", "t1")
	wErr := httptest.NewRecorder()
	// Re-init context properly
	cErr, _ = gin.CreateTestContext(wErr)
	cErr.Request = httptest.NewRequest("POST", "/autoschedule", bytes.NewBuffer(bodyErr))
	cErr.Set("tenant_id", "t1")
	
	h.AutoSchedule(cErr)
	assert.Equal(t, http.StatusInternalServerError, wErr.Code)
}

func TestSchedulerHandler_ListTerms_Fallback(t *testing.T) {
	mockSvc := new(MockSchedulerService)
	h := NewSchedulerHandler(mockSvc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// No tenant_id in context, but query param
	c.Request, _ = http.NewRequest("GET", "/terms?tenant_id=t-fallback", nil)
	
	mockSvc.On("ListTerms", mock.Anything, "t-fallback").Return([]models.AcademicTerm{}, nil) // Empty list ok

	h.ListTerms(c)
	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

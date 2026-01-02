package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupSchedulerHandler() (*SchedulerHandler, *HMockSchedulerRepo) {
	mockRepo := new(HMockSchedulerRepo)
	// Pass nil for unused dependencies (Resource, Curriculum, User, Mailer)
	// assuming tests don't hit code paths requiring them.
	svc := services.NewSchedulerService(mockRepo, nil, nil, nil, nil)
	return NewSchedulerHandler(svc), mockRepo
}

func TestSchedulerHandler_CreateTerm(t *testing.T) {
	h, mockRepo := setupSchedulerHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	term := models.AcademicTerm{Name: "Fall 2026", Code: "2026-FA", StartDate: time.Now(), EndDate: time.Now().AddDate(0,4,0)}
	body, _ := json.Marshal(term)
	c.Request, _ = http.NewRequest("POST", "/terms", bytes.NewBuffer(body))
    c.Set("tenant_id", "t1")

	mockRepo.On("CreateTerm", mock.Anything, mock.MatchedBy(func(a *models.AcademicTerm) bool {
		return a.Name == "Fall 2026" && a.Code == "2026-FA" && a.TenantID == "t1"
	})).Return(nil)

	h.CreateTerm(c)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestSchedulerHandler_ListTerms(t *testing.T) {
	h, mockRepo := setupSchedulerHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/terms", nil)
	c.Set("tenant_id", "t1")

	expected := []models.AcademicTerm{{Name: "Spring 2025"}}
	mockRepo.On("ListTerms", mock.Anything, "t1").Return(expected, nil)

	h.ListTerms(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSchedulerHandler_CreateOffering(t *testing.T) {
	h, mockRepo := setupSchedulerHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	offering := models.CourseOffering{TermID: "term-1", CourseID: "course-1"}
	body, _ := json.Marshal(offering)
	c.Request, _ = http.NewRequest("POST", "/offerings", bytes.NewBuffer(body))

	mockRepo.On("CreateOffering", mock.Anything, mock.AnythingOfType("*models.CourseOffering")).Return(nil)

	h.CreateOffering(c)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestSchedulerHandler_CreateSession(t *testing.T) {
	h, mockRepo := setupSchedulerHandler()
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
	
	// CreateSession logic calls: GetOffering -> CheckConflicts -> CreateSession
	// Must mock GetOffering
	mockRepo.On("GetOffering", mock.Anything, "off-1").Return(&models.CourseOffering{ID: "off-1", DeliveryFormat: "IN_PERSON"}, nil)

	// And CreateSession
	mockRepo.On("CreateSession", mock.Anything, mock.AnythingOfType("*models.ClassSession")).Return(nil)

	h.CreateSession(c)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestSchedulerHandler_ListSessions(t *testing.T) {
	h, mockRepo := setupSchedulerHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	c.Request, _ = http.NewRequest("GET", "/sessions?offering_id=off-1&start=2025-01-01T00:00:00Z&end=2025-02-01T00:00:00Z", nil)
	
	mockRepo.On("ListSessions", mock.Anything, "off-1", mock.Anything, mock.Anything).Return([]models.ClassSession{{ID: "sess-1"}}, nil)

	h.ListSessions(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

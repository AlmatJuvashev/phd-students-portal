package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// -- Duplicate Mocks for Handler Test (Self-contained) --

type HMockTranscriptRepo struct {
	mock.Mock
}

func (m *HMockTranscriptRepo) GetStudentGrades(ctx context.Context, studentID string) ([]models.TermGrade, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.TermGrade), args.Error(1)
}

type HMockSchedulerRepo struct {
	mock.Mock
}

func (m *HMockSchedulerRepo) GetTerm(ctx context.Context, id string) (*models.AcademicTerm, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.AcademicTerm), args.Error(1)
}

// Satisfy Interface stubbing (minimal)
func (m *HMockSchedulerRepo) CreateTerm(ctx context.Context, term *models.AcademicTerm) error { return nil }
func (m *HMockSchedulerRepo) ListTerms(ctx context.Context, tenantID string) ([]models.AcademicTerm, error) { return nil, nil }
func (m *HMockSchedulerRepo) UpdateTerm(ctx context.Context, term *models.AcademicTerm) error { return nil }
func (m *HMockSchedulerRepo) DeleteTerm(ctx context.Context, id string) error { return nil }
func (m *HMockSchedulerRepo) CreateOffering(ctx context.Context, offering *models.CourseOffering) error { return nil }
func (m *HMockSchedulerRepo) GetOffering(ctx context.Context, id string) (*models.CourseOffering, error) { return nil, nil }
func (m *HMockSchedulerRepo) ListOfferings(ctx context.Context, tenantID string, termID string) ([]models.CourseOffering, error) { return nil, nil }
func (m *HMockSchedulerRepo) ListOfferingsByInstructor(ctx context.Context, instructorID string, termID string) ([]models.CourseOffering, error) { return nil, nil }
func (m *HMockSchedulerRepo) UpdateOffering(ctx context.Context, offering *models.CourseOffering) error { return nil }
func (m *HMockSchedulerRepo) AddStaff(ctx context.Context, staff *models.CourseStaff) error { return nil }
func (m *HMockSchedulerRepo) ListStaff(ctx context.Context, offeringID string) ([]models.CourseStaff, error) { return nil, nil }
func (m *HMockSchedulerRepo) RemoveStaff(ctx context.Context, id string) error { return nil }
func (m *HMockSchedulerRepo) CreateSession(ctx context.Context, session *models.ClassSession) error { return nil }
func (m *HMockSchedulerRepo) ListSessions(ctx context.Context, offeringID string, startDate, endDate time.Time) ([]models.ClassSession, error) { return nil, nil }
func (m *HMockSchedulerRepo) ListSessionsByRoom(ctx context.Context, roomID string, startDate, endDate time.Time) ([]models.ClassSession, error) { return nil, nil }
func (m *HMockSchedulerRepo) ListSessionsByInstructor(ctx context.Context, instructorID string, startDate, endDate time.Time) ([]models.ClassSession, error) { return nil, nil }
func (m *HMockSchedulerRepo) ListSessionsForTerm(ctx context.Context, termID string) ([]models.ClassSession, error) { return nil, nil }
func (m *HMockSchedulerRepo) UpdateSession(ctx context.Context, session *models.ClassSession) error { return nil }
func (m *HMockSchedulerRepo) DeleteSession(ctx context.Context, id string) error { return nil }
func (m *HMockSchedulerRepo) AddCohortToOffering(ctx context.Context, offeringID, cohortID string) error { return nil }
func (m *HMockSchedulerRepo) GetOfferingCohorts(ctx context.Context, offeringID string) ([]string, error) { return nil, nil }
func (m *HMockSchedulerRepo) ListSessionsForCohorts(ctx context.Context, cohortIDs []string, startTime, endTime time.Time) ([]models.ClassSession, error) { return nil, nil }


func TestTranscriptHandler_GetStudentTranscript(t *testing.T) {
	// Setup Mocks
	mockTR := new(HMockTranscriptRepo)
	mockSR := new(HMockSchedulerRepo)
	
	// Setup Service
	service := services.NewTranscriptService(mockTR, mockSR)
	
	// Setup Handler
	handler := handlers.NewTranscriptHandler(service)

	// Setup Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", "student-123") // Mock Auth Middleware
	})
	router.GET("/transcript", handler.GetStudentTranscript)

	// Mock Data
	g1 := models.TermGrade{TermID: "t1", CourseCode: "A1", Credits: 3, GradePoints: 4, Percentage: 100}
	term1 := &models.AcademicTerm{ID: "t1", Name: "Term 1"}

	// Expectations
	mockTR.On("GetStudentGrades", mock.Anything, "student-123").Return([]models.TermGrade{g1}, nil)
	mockSR.On("GetTerm", mock.Anything, "t1").Return(term1, nil)

	// Requests
	req, _ := http.NewRequest(http.MethodGet, "/transcript", nil)
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response models.Transcript
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "student-123", response.StudentID)
	assert.Equal(t, float32(4.0), response.CumulativeGPA)

	// Case 2: Service Error
	mockTR.On("GetStudentGrades", mock.Anything, "student-error").Return(nil, assert.AnError)
	
	reqErr, _ := http.NewRequest(http.MethodGet, "/transcript", nil)
	wErr := httptest.NewRecorder()
	// Create context with Error Student
	c, _ := gin.CreateTestContext(wErr)
	c.Request = reqErr
	c.Set("userID", "student-error")
	
	handler.GetStudentTranscript(c)
	assert.Equal(t, http.StatusInternalServerError, wErr.Code)

	// Case 3: Missing UserID
	wAuthErr := httptest.NewRecorder()
	cAuth, _ := gin.CreateTestContext(wAuthErr)
	cAuth.Request = reqErr
	// No UserID set
	handler.GetStudentTranscript(cAuth)
	assert.Equal(t, http.StatusUnauthorized, wAuthErr.Code)
}

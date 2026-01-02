package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Reusing Mocks from services package (but they are in 'services' package, distinct).
// Need to define temporary mocks or import exposed mocks if possible.
// Since mocks are often test-only, I'll define simple mocks here for the service dependencies
// OR better, mock the SERVICE itself if I can interface it.
// TeacherService is a struct, not an interface. I have to mock its dependencies (Repos).
// I will copy the minimal mocks needed for Repo interfaces.

type HMocks struct {
	Sched   *HMockSchedulerRepo
	LMS     *HMockLMSRepo
	Grading *HMockGradingRepo
}

type HMockSchedulerRepo struct { mock.Mock }
// Terms
func (m *HMockSchedulerRepo) CreateTerm(ctx context.Context, t *models.AcademicTerm) error { return nil }
func (m *HMockSchedulerRepo) GetTerm(ctx context.Context, id string) (*models.AcademicTerm, error) { return nil, nil }
func (m *HMockSchedulerRepo) ListTerms(ctx context.Context, tenantID string) ([]models.AcademicTerm, error) { return nil, nil }
func (m *HMockSchedulerRepo) UpdateTerm(ctx context.Context, t *models.AcademicTerm) error { return nil }
func (m *HMockSchedulerRepo) DeleteTerm(ctx context.Context, id string) error { return nil }
// Offerings
func (m *HMockSchedulerRepo) CreateOffering(ctx context.Context, o *models.CourseOffering) error { return nil }
func (m *HMockSchedulerRepo) GetOffering(ctx context.Context, id string) (*models.CourseOffering, error) { return nil, nil }
func (m *HMockSchedulerRepo) ListOfferings(ctx context.Context, tID, termID string) ([]models.CourseOffering, error) { return nil, nil }
func (m *HMockSchedulerRepo) ListOfferingsByInstructor(ctx context.Context, iID, termID string) ([]models.CourseOffering, error) { 
    args := m.Called(ctx, iID, termID)
    if args.Get(0) == nil { return nil, args.Error(1) }
    return args.Get(0).([]models.CourseOffering), args.Error(1)
}
func (m *HMockSchedulerRepo) UpdateOffering(ctx context.Context, o *models.CourseOffering) error { return nil }
// Staff
func (m *HMockSchedulerRepo) AddStaff(ctx context.Context, s *models.CourseStaff) error { return nil }
func (m *HMockSchedulerRepo) ListStaff(ctx context.Context, oID string) ([]models.CourseStaff, error) { return nil, nil }
func (m *HMockSchedulerRepo) RemoveStaff(ctx context.Context, id string) error { return nil }
// Sessions
func (m *HMockSchedulerRepo) CreateSession(ctx context.Context, s *models.ClassSession) error { return nil }
func (m *HMockSchedulerRepo) UpdateSession(ctx context.Context, s *models.ClassSession) error { return nil }
func (m *HMockSchedulerRepo) DeleteSession(ctx context.Context, id string) error { return nil }
func (m *HMockSchedulerRepo) ListSessions(ctx context.Context, oID string, s, e time.Time) ([]models.ClassSession, error) { return nil, nil }
func (m *HMockSchedulerRepo) ListSessionsForTerm(ctx context.Context, tID string) ([]models.ClassSession, error) { return nil, nil }
func (m *HMockSchedulerRepo) ListSessionsByRoom(ctx context.Context, rID string, s, e time.Time) ([]models.ClassSession, error) { return nil, nil }
func (m *HMockSchedulerRepo) ListSessionsByInstructor(ctx context.Context, iID string, s, e time.Time) ([]models.ClassSession, error) {
    args := m.Called(ctx, iID, s, e)
    if args.Get(0) == nil { return nil, args.Error(1) }
    return args.Get(0).([]models.ClassSession), args.Error(1)
}

type HMockLMSRepo struct{ mock.Mock }
func (m *HMockLMSRepo) EnrollStudent(ctx context.Context, e *models.CourseEnrollment) error { return nil }
func (m *HMockLMSRepo) GetCourseRoster(ctx context.Context, offeringID string) ([]models.CourseEnrollment, error) {
	args := m.Called(ctx, offeringID)
	return args.Get(0).([]models.CourseEnrollment), args.Error(1)
}
func (m *HMockLMSRepo) GetStudentEnrollments(ctx context.Context, studentID string) ([]models.CourseEnrollment, error) { return nil, nil }
func (m *HMockLMSRepo) UpdateEnrollmentStatus(ctx context.Context, id, s string) error { return nil }
func (m *HMockLMSRepo) CreateSubmission(ctx context.Context, s *models.ActivitySubmission) error { return nil }
func (m *HMockLMSRepo) GetSubmission(ctx context.Context, a, s string) (*models.ActivitySubmission, error) { return nil, nil }
func (m *HMockLMSRepo) ListSubmissions(ctx context.Context, o string) ([]models.ActivitySubmission, error) { return nil, nil }
func (m *HMockLMSRepo) MarkAttendance(ctx context.Context, a *models.ClassAttendance) error { return nil }
func (m *HMockLMSRepo) GetSessionAttendance(ctx context.Context, s string) ([]models.ClassAttendance, error) { return nil, nil }

// Need full interface compliance
// Copying from lms_repository.go interface
// (Already implemented above)

type HMockGradingRepo struct{ mock.Mock }
func (m *HMockGradingRepo) CreateSchema(ctx context.Context, s *models.GradingSchema) error { return nil }
func (m *HMockGradingRepo) GetSchema(ctx context.Context, id string) (*models.GradingSchema, error) { return nil, nil }
func (m *HMockGradingRepo) ListSchemas(ctx context.Context, t string) ([]models.GradingSchema, error) { return nil, nil }
func (m *HMockGradingRepo) GetDefaultSchema(ctx context.Context, t string) (*models.GradingSchema, error) { return nil, nil }
func (m *HMockGradingRepo) UpdateSchema(ctx context.Context, s *models.GradingSchema) error { return nil }
func (m *HMockGradingRepo) DeleteSchema(ctx context.Context, id string) error { return nil }
func (m *HMockGradingRepo) CreateEntry(ctx context.Context, e *models.GradebookEntry) error { return nil }
func (m *HMockGradingRepo) GetEntry(ctx context.Context, id string) (*models.GradebookEntry, error) { return nil, nil }
func (m *HMockGradingRepo) GetEntryByActivity(ctx context.Context, o, a, s string) (*models.GradebookEntry, error) { return nil, nil }
func (m *HMockGradingRepo) ListEntries(ctx context.Context, o string) ([]models.GradebookEntry, error) { 
	args := m.Called(ctx, o)
	return args.Get(0).([]models.GradebookEntry), args.Error(1)
}
func (m *HMockGradingRepo) ListStudentEntries(ctx context.Context, s string) ([]models.GradebookEntry, error) { return nil, nil }

func setupTeacherHandler() (*TeacherHandler, *HMocks) {
	sched := new(HMockSchedulerRepo)
	lms := new(HMockLMSRepo)
	grad := new(HMockGradingRepo)
	svc := services.NewTeacherService(sched, lms, grad)
	return NewTeacherHandler(svc), &HMocks{Sched: sched, LMS: lms, Grading: grad}
}

func TestTeacherHandler_GetDashboardStats(t *testing.T) {
	h, mocks := setupTeacherHandler()
	
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/teacher/dashboard", nil)
	c.Set("claims", jwt.MapClaims{"sub": "inst-1"})

	// Dependencies
	mocks.Sched.On("ListSessionsByInstructor", mock.Anything, "inst-1", mock.Anything, mock.Anything).Return([]models.ClassSession{}, nil)
	mocks.Sched.On("ListOfferingsByInstructor", mock.Anything, "inst-1", "").Return([]models.CourseOffering{}, nil)

	h.GetDashboardStats(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTeacherHandler_GetCourseRoster(t *testing.T) {
	h, mocks := setupTeacherHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/teacher/courses/123/roster", nil)
	c.Params = gin.Params{{Key: "id", Value: "123"}}

	mocks.LMS.On("GetCourseRoster", mock.Anything, "123").Return([]models.CourseEnrollment{}, nil)

	h.GetCourseRoster(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTeacherHandler_GetGradebook(t *testing.T) {
	h, mocks := setupTeacherHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/teacher/courses/123/gradebook", nil)
	c.Params = gin.Params{{Key: "id", Value: "123"}}

	mocks.Grading.On("ListEntries", mock.Anything, "123").Return([]models.GradebookEntry{}, nil)

	h.GetGradebook(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

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

type HMockSchedulerRepo struct{ mock.Mock }

// Terms
func (m *HMockSchedulerRepo) CreateTerm(ctx context.Context, t *models.AcademicTerm) error {
	return m.Called(ctx, t).Error(0)
}
func (m *HMockSchedulerRepo) GetTerm(ctx context.Context, id string) (*models.AcademicTerm, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AcademicTerm), args.Error(1)
}
func (m *HMockSchedulerRepo) ListTerms(ctx context.Context, tenantID string) ([]models.AcademicTerm, error) {
	args := m.Called(ctx, tenantID)
	return args.Get(0).([]models.AcademicTerm), args.Error(1)
}
func (m *HMockSchedulerRepo) UpdateTerm(ctx context.Context, t *models.AcademicTerm) error {
	return m.Called(ctx, t).Error(0)
}
func (m *HMockSchedulerRepo) DeleteTerm(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// Offerings
func (m *HMockSchedulerRepo) CreateOffering(ctx context.Context, o *models.CourseOffering) error {
	return m.Called(ctx, o).Error(0)
}
func (m *HMockSchedulerRepo) GetOffering(ctx context.Context, id string) (*models.CourseOffering, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CourseOffering), args.Error(1)
}
func (m *HMockSchedulerRepo) ListOfferings(ctx context.Context, tID, termID string) ([]models.CourseOffering, error) {
	args := m.Called(ctx, tID, termID)
	return args.Get(0).([]models.CourseOffering), args.Error(1)
}
func (m *HMockSchedulerRepo) ListOfferingsByInstructor(ctx context.Context, iID, termID string) ([]models.CourseOffering, error) {
	args := m.Called(ctx, iID, termID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.CourseOffering), args.Error(1)
}
func (m *HMockSchedulerRepo) UpdateOffering(ctx context.Context, o *models.CourseOffering) error {
	return m.Called(ctx, o).Error(0)
}

// Staff
func (m *HMockSchedulerRepo) AddStaff(ctx context.Context, s *models.CourseStaff) error {
	return m.Called(ctx, s).Error(0)
}
func (m *HMockSchedulerRepo) ListStaff(ctx context.Context, oID string) ([]models.CourseStaff, error) {
	args := m.Called(ctx, oID)
	return args.Get(0).([]models.CourseStaff), args.Error(1)
}
func (m *HMockSchedulerRepo) RemoveStaff(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// Sessions
func (m *HMockSchedulerRepo) CreateSession(ctx context.Context, s *models.ClassSession) error {
	return m.Called(ctx, s).Error(0)
}
func (m *HMockSchedulerRepo) UpdateSession(ctx context.Context, s *models.ClassSession) error {
	return m.Called(ctx, s).Error(0)
}
func (m *HMockSchedulerRepo) DeleteSession(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *HMockSchedulerRepo) ListSessions(ctx context.Context, oID string, s, e time.Time) ([]models.ClassSession, error) {
	args := m.Called(ctx, oID, s, e)
	return args.Get(0).([]models.ClassSession), args.Error(1)
}
func (m *HMockSchedulerRepo) ListSessionsForTerm(ctx context.Context, tID string) ([]models.ClassSession, error) {
	args := m.Called(ctx, tID)
	return args.Get(0).([]models.ClassSession), args.Error(1)
}
func (m *HMockSchedulerRepo) ListSessionsByRoom(ctx context.Context, rID string, s, e time.Time) ([]models.ClassSession, error) {
	args := m.Called(ctx, rID, s, e)
	return args.Get(0).([]models.ClassSession), args.Error(1)
}
func (m *HMockSchedulerRepo) ListSessionsByInstructor(ctx context.Context, iID string, s, e time.Time) ([]models.ClassSession, error) {
	args := m.Called(ctx, iID, s, e)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ClassSession), args.Error(1)
}
func (m *HMockSchedulerRepo) AddCohortToOffering(ctx context.Context, offeringID, cohortID string) error {
	return m.Called(ctx, offeringID, cohortID).Error(0)
}
func (m *HMockSchedulerRepo) GetOfferingCohorts(ctx context.Context, offeringID string) ([]string, error) {
	args := m.Called(ctx, offeringID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}
func (m *HMockSchedulerRepo) ListSessionsForCohorts(ctx context.Context, cohortIDs []string, s, e time.Time) ([]models.ClassSession, error) {
	args := m.Called(ctx, cohortIDs, s, e)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ClassSession), args.Error(1)
}

type HMockLMSRepo struct{ mock.Mock }

func (m *HMockLMSRepo) EnrollStudent(ctx context.Context, e *models.CourseEnrollment) error {
	return nil
}
func (m *HMockLMSRepo) GetCourseRoster(ctx context.Context, offeringID string) ([]models.CourseEnrollment, error) {
	args := m.Called(ctx, offeringID)
	return args.Get(0).([]models.CourseEnrollment), args.Error(1)
}
func (m *HMockLMSRepo) GetStudentEnrollments(ctx context.Context, studentID string) ([]models.CourseEnrollment, error) {
	return nil, nil
}
func (m *HMockLMSRepo) UpdateEnrollmentStatus(ctx context.Context, id, s string) error { return nil }
func (m *HMockLMSRepo) CreateSubmission(ctx context.Context, s *models.ActivitySubmission) error {
	return nil
}
func (m *HMockLMSRepo) GetSubmission(ctx context.Context, id string) (*models.ActivitySubmission, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ActivitySubmission), args.Error(1)
}
func (m *HMockLMSRepo) GetSubmissionByStudent(ctx context.Context, activityID, studentID string) (*models.ActivitySubmission, error) {
	args := m.Called(ctx, activityID, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ActivitySubmission), args.Error(1)
}
func (m *HMockLMSRepo) ListSubmissions(ctx context.Context, offeringID string) ([]models.ActivitySubmission, error) {
	args := m.Called(ctx, offeringID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ActivitySubmission), args.Error(1)
}
func (m *HMockLMSRepo) MarkAttendance(ctx context.Context, att *models.ClassAttendance) error {
	args := m.Called(ctx, att)
	return args.Error(0)
}
func (m *HMockLMSRepo) CreateAnnotation(ctx context.Context, ann models.SubmissionAnnotation) (*models.SubmissionAnnotation, error) {
	args := m.Called(ctx, ann)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SubmissionAnnotation), args.Error(1)
}
func (m *HMockLMSRepo) ListAnnotations(ctx context.Context, submissionID string) ([]models.SubmissionAnnotation, error) {
	args := m.Called(ctx, submissionID)
	return args.Get(0).([]models.SubmissionAnnotation), args.Error(1)
}
func (m *HMockLMSRepo) DeleteAnnotation(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *HMockLMSRepo) GetSessionAttendance(ctx context.Context, s string) ([]models.ClassAttendance, error) {
	return nil, nil
}

// Need full interface compliance
// Copying from lms_repository.go interface
// (Already implemented above)

type HMockGradingRepo struct{ mock.Mock }

func (m *HMockGradingRepo) CreateSchema(ctx context.Context, s *models.GradingSchema) error {
	return nil
}
func (m *HMockGradingRepo) GetSchema(ctx context.Context, id string) (*models.GradingSchema, error) {
	return nil, nil
}
func (m *HMockGradingRepo) ListSchemas(ctx context.Context, t string) ([]models.GradingSchema, error) {
	return nil, nil
}
func (m *HMockGradingRepo) GetDefaultSchema(ctx context.Context, t string) (*models.GradingSchema, error) {
	return nil, nil
}
func (m *HMockGradingRepo) UpdateSchema(ctx context.Context, s *models.GradingSchema) error {
	return nil
}
func (m *HMockGradingRepo) DeleteSchema(ctx context.Context, id string) error { return nil }
func (m *HMockGradingRepo) CreateEntry(ctx context.Context, e *models.GradebookEntry) error {
	return nil
}
func (m *HMockGradingRepo) GetEntry(ctx context.Context, id string) (*models.GradebookEntry, error) {
	return nil, nil
}
func (m *HMockGradingRepo) GetEntryByActivity(ctx context.Context, o, a, s string) (*models.GradebookEntry, error) {
	return nil, nil
}
func (m *HMockGradingRepo) ListEntries(ctx context.Context, o string) ([]models.GradebookEntry, error) {
	args := m.Called(ctx, o)
	return args.Get(0).([]models.GradebookEntry), args.Error(1)
}
func (m *HMockGradingRepo) ListStudentEntries(ctx context.Context, s string) ([]models.GradebookEntry, error) {
	return nil, nil
}

func setupTeacherHandler() (*TeacherHandler, *HMocks) {
	sched := new(HMockSchedulerRepo)
	lms := new(HMockLMSRepo)
	grad := new(HMockGradingRepo)
	svc := services.NewTeacherService(sched, lms, grad, nil)
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

func TestTeacherHandler_GetCourseStudents(t *testing.T) {
	h, mocks := setupTeacherHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/teacher/courses/123/students", nil)
	c.Params = gin.Params{{Key: "id", Value: "123"}}

	mocks.LMS.On("GetCourseRoster", mock.Anything, "123").Return([]models.CourseEnrollment{
		{StudentID: "stu-1", StudentName: "Jane Doe", StudentEmail: "jane@example.com"},
	}, nil)
	mocks.LMS.On("ListSubmissions", mock.Anything, "123").Return([]models.ActivitySubmission{
		{ID: "sub-1", StudentID: "stu-1", ActivityID: "act-1", Status: "SUBMITTED", SubmittedAt: time.Now()},
	}, nil)
	mocks.Grading.On("ListEntries", mock.Anything, "123").Return([]models.GradebookEntry{
		{StudentID: "stu-1", ActivityID: "act-1", Score: 85, MaxScore: 100, Grade: "B", GradedAt: time.Now()},
	}, nil)

	h.GetCourseStudents(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "stu-1")
}

func TestTeacherHandler_GetAtRiskStudents(t *testing.T) {
	h, mocks := setupTeacherHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/teacher/courses/123/at-risk", nil)
	c.Params = gin.Params{{Key: "id", Value: "123"}}

	// High risk due to inactivity (>=14 days).
	mocks.LMS.On("GetCourseRoster", mock.Anything, "123").Return([]models.CourseEnrollment{
		{StudentID: "stu-1", StudentName: "Jane Doe", StudentEmail: "jane@example.com"},
	}, nil)
	mocks.LMS.On("ListSubmissions", mock.Anything, "123").Return([]models.ActivitySubmission{
		{ID: "sub-1", StudentID: "stu-1", ActivityID: "act-1", Status: "SUBMITTED", SubmittedAt: time.Now().Add(-20 * 24 * time.Hour)},
	}, nil)
	mocks.Grading.On("ListEntries", mock.Anything, "123").Return([]models.GradebookEntry{}, nil)

	h.GetAtRiskStudents(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "\"risk_level\"")
}

func TestTeacherHandler_GetStudentActivityLog(t *testing.T) {
	h, mocks := setupTeacherHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/teacher/students/stu-1/activity?course_offering_id=123&limit=10", nil)
	c.Params = gin.Params{{Key: "id", Value: "stu-1"}}

	mocks.LMS.On("ListSubmissions", mock.Anything, "123").Return([]models.ActivitySubmission{
		{ID: "sub-1", StudentID: "stu-1", ActivityID: "act-1", ActivityTitle: "Assignment 1", Status: "SUBMITTED", SubmittedAt: time.Now()},
	}, nil)
	mocks.Grading.On("ListEntries", mock.Anything, "123").Return([]models.GradebookEntry{
		{StudentID: "stu-1", ActivityID: "act-1", Score: 90, MaxScore: 100, Grade: "A", GradedAt: time.Now()},
	}, nil)

	h.GetStudentActivityLog(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "\"kind\"")
}

func TestTeacherHandler_GetMySchedule(t *testing.T) {
	h, _ := setupTeacherHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/teacher/schedule", nil)
	c.Set("claims", jwt.MapClaims{"sub": "inst-1"})

	h.GetMySchedule(c)
	// Currently returns 501
	assert.Equal(t, http.StatusNotImplemented, w.Code)
}

func TestTeacherHandler_GetMyCourses(t *testing.T) {
	h, mocks := setupTeacherHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/teacher/courses", nil)
	c.Set("claims", jwt.MapClaims{"sub": "inst-1"})

	mocks.Sched.On("ListOfferingsByInstructor", mock.Anything, "inst-1", "").Return([]models.CourseOffering{{ID: "c1"}}, nil)

	h.GetMyCourses(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTeacherHandler_GetSubmissions(t *testing.T) {
	h, mocks := setupTeacherHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/teacher/submissions", nil)
	c.Set("claims", jwt.MapClaims{"sub": "inst-1"})

	// Service first gets courses for instructor
	mocks.Sched.On("ListOfferingsByInstructor", mock.Anything, "inst-1", "").Return([]models.CourseOffering{{ID: "c1"}}, nil)
	// Then gets submissions for each course
	mocks.LMS.On("ListSubmissions", mock.Anything, "c1").Return([]models.ActivitySubmission{{ID: "s1"}}, nil)

	h.GetSubmissions(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/dto"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	pb "github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type hStudentUserRepo struct {
	byID map[string]*models.User
}

func (s *hStudentUserRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	if s.byID == nil {
		return nil, nil
	}
	return s.byID[id], nil
}

type hStudentJourneyRepo struct {
	state map[string]string
}

func (s *hStudentJourneyRepo) GetJourneyState(ctx context.Context, userID, tenantID string) (map[string]string, error) {
	return s.state, nil
}

type hStudentLMSRepo struct {
	enrollments []models.CourseEnrollment
}

func (s *hStudentLMSRepo) GetStudentEnrollments(ctx context.Context, studentID string) ([]models.CourseEnrollment, error) {
	return s.enrollments, nil
}

func (s *hStudentLMSRepo) CreateSubmission(ctx context.Context, sub *models.ActivitySubmission) error {
	return nil
}

func (s *hStudentLMSRepo) GetSubmissionByStudent(ctx context.Context, activityID, studentID string) (*models.ActivitySubmission, error) {
	return nil, nil
}

type hStudentSchedulerRepo struct {
	offering *models.CourseOffering
	staff    []models.CourseStaff
	sessions []models.ClassSession
}

func (s *hStudentSchedulerRepo) GetOffering(ctx context.Context, id string) (*models.CourseOffering, error) {
	if s.offering != nil && s.offering.ID == id {
		return s.offering, nil
	}
	return nil, nil
}

func (s *hStudentSchedulerRepo) ListStaff(ctx context.Context, offeringID string) ([]models.CourseStaff, error) {
	return s.staff, nil
}

func (s *hStudentSchedulerRepo) ListSessions(ctx context.Context, offeringID string, startDate, endDate time.Time) ([]models.ClassSession, error) {
	return s.sessions, nil
}

type hStudentCurriculumRepo struct {
	byID map[string]*models.Course
}

func (s *hStudentCurriculumRepo) GetCourse(ctx context.Context, id string) (*models.Course, error) {
	if s.byID == nil {
		return nil, nil
	}
	return s.byID[id], nil
}

func (s *hStudentCurriculumRepo) ListCourses(ctx context.Context, tenantID string, programID *string) ([]models.Course, error) {
	if s.byID == nil {
		return nil, nil
	}
	var list []models.Course
	for _, c := range s.byID {
		list = append(list, *c)
	}
	return list, nil
}

type hStudentGradingRepo struct {
	entries []models.GradebookEntry
}

func (s *hStudentGradingRepo) ListStudentEntries(ctx context.Context, studentID string) ([]models.GradebookEntry, error) {
	return s.entries, nil
}

func TestStudentHandler_GetDashboard_UnauthorizedWithoutContext(t *testing.T) {
	h := NewStudentHandler(nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/student/dashboard", nil)

	h.GetDashboard(c)
	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestStudentHandler_GetDashboard_Success(t *testing.T) {
	pbm := &pb.Manager{
		DefaultLocale: "en",
		Nodes: map[string]pb.Node{
			"n1": {ID: "n1", Title: map[string]string{"en": "Profile"}},
			"n2": {ID: "n2", Title: map[string]string{"en": "Proposal"}, Prerequisites: []string{"n1"}},
		},
		NodeWorlds: map[string]string{"n1": "W1", "n2": "W2"},
	}

	userRepo := &hStudentUserRepo{
		byID: map[string]*models.User{
			"stud-1": {ID: "stud-1", Program: "PhD in Medicine"},
		},
	}
	journeyRepo := &hStudentJourneyRepo{state: map[string]string{"n1": "done", "n2": "todo"}}
	gradingRepo := &hStudentGradingRepo{entries: []models.GradebookEntry{}}

	svc := services.NewStudentService(
		userRepo,
		journeyRepo,
		&hStudentLMSRepo{},
		&hStudentSchedulerRepo{},
		&hStudentCurriculumRepo{},
		gradingRepo,
		nil,
		nil,
		nil,
		pbm,
	)
	h := NewStudentHandler(svc)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/student/dashboard", nil)
	c.Set("tenant_id", "tenant-1")
	c.Set("userID", "stud-1")

	h.GetDashboard(c)
	require.Equal(t, http.StatusOK, w.Code)

	var payload dto.StudentDashboard
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &payload))
	require.Equal(t, "PhD in Medicine", payload.Program.Title)
	require.Equal(t, 2, payload.Program.TotalNodes)
	require.Len(t, payload.UpcomingDeadlines, 1)
	require.Equal(t, "n2", payload.UpcomingDeadlines[0].ID)
}

func TestStudentHandler_ListCourses_Success(t *testing.T) {
	nextDate := time.Now().Add(72 * time.Hour)
	roomID := "room-1"

	userRepo := &hStudentUserRepo{
		byID: map[string]*models.User{
			"inst-1": {ID: "inst-1", FirstName: "Aida", LastName: "Baken", Email: "aida@example.com"},
		},
	}
	lmsRepo := &hStudentLMSRepo{
		enrollments: []models.CourseEnrollment{
			{ID: "enr-1", CourseOfferingID: "off-1", Status: "ENROLLED"},
		},
	}
	schedulerRepo := &hStudentSchedulerRepo{
		offering: &models.CourseOffering{ID: "off-1", CourseID: "course-1", TermID: "term-1", Section: "A", DeliveryFormat: "IN_PERSON"},
		staff:    []models.CourseStaff{{UserID: "inst-1", Role: "INSTRUCTOR", IsPrimary: true}},
		sessions: []models.ClassSession{{ID: "sess-1", CourseOfferingID: "off-1", Date: nextDate, StartTime: "10:00", EndTime: "11:30", RoomID: &roomID, Type: "LECTURE"}},
	}
	currRepo := &hStudentCurriculumRepo{
		byID: map[string]*models.Course{
			"course-1": {ID: "course-1", Code: "MED101", Title: `{"en":"Intro to Medicine"}`},
		},
	}

	svc := services.NewStudentService(
		userRepo,
		&hStudentJourneyRepo{},
		lmsRepo,
		schedulerRepo,
		currRepo,
		&hStudentGradingRepo{},
		nil,
		nil,
		nil,
		nil,
	)
	h := NewStudentHandler(svc)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/student/courses", nil)
	c.Set("tenant_id", "tenant-1")
	c.Set("userID", "stud-1")

	h.ListCourses(c)
	require.Equal(t, http.StatusOK, w.Code)

	var payload []dto.StudentCourse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &payload))
	require.Len(t, payload, 1)
	require.Equal(t, "MED101", payload[0].Code)
	require.Equal(t, "Intro to Medicine", payload[0].Title)
	require.NotNil(t, payload[0].NextSession)
	require.Equal(t, nextDate.Format("2006-01-02"), payload[0].NextSession.Date)
}

func TestStudentHandler_ListAssignments_Success(t *testing.T) {
	pbm := &pb.Manager{
		DefaultLocale: "en",
		Nodes: map[string]pb.Node{
			"n1": {ID: "n1", Title: map[string]string{"en": "Profile"}},
			"n2": {ID: "n2", Title: map[string]string{"en": "Proposal"}, Prerequisites: []string{"n1"}},
		},
		NodeWorlds: map[string]string{"n1": "W1", "n2": "W2"},
	}

	userRepo := &hStudentUserRepo{byID: map[string]*models.User{"stud-1": {ID: "stud-1"}}}
	journeyRepo := &hStudentJourneyRepo{state: map[string]string{"n1": "done", "n2": "submitted"}}

	svc := services.NewStudentService(
		userRepo,
		journeyRepo,
		&hStudentLMSRepo{},
		&hStudentSchedulerRepo{},
		&hStudentCurriculumRepo{},
		&hStudentGradingRepo{},
		nil,
		nil,
		nil,
		pbm,
	)
	h := NewStudentHandler(svc)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/student/assignments", nil)
	c.Set("tenant_id", "tenant-1")
	c.Set("userID", "stud-1")

	h.ListAssignments(c)
	require.Equal(t, http.StatusOK, w.Code)

	var payload []dto.StudentAssignment
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &payload))
	require.Len(t, payload, 1)
	require.Equal(t, "n2", payload[0].ID)
	require.Equal(t, "urgent", payload[0].Severity)
}

func TestStudentHandler_ListGrades_Success(t *testing.T) {
	gradedAt := time.Date(2025, 12, 31, 10, 30, 0, 0, time.UTC)

	gradingRepo := &hStudentGradingRepo{
		entries: []models.GradebookEntry{
			{
				ID:               "g1",
				CourseOfferingID: "off-1",
				ActivityID:       "act-1",
				StudentID:        "stud-1",
				Score:            95,
				MaxScore:         100,
				Grade:            "A",
				Feedback:         "Excellent",
				GradedByID:       "inst-1",
				GradedAt:         gradedAt,
			},
		},
	}
	schedulerRepo := &hStudentSchedulerRepo{
		offering: &models.CourseOffering{ID: "off-1", CourseID: "course-1"},
	}
	currRepo := &hStudentCurriculumRepo{
		byID: map[string]*models.Course{
			"course-1": {ID: "course-1", Code: "MED101", Title: `{"en":"Intro to Medicine"}`},
		},
	}

	svc := services.NewStudentService(
		&hStudentUserRepo{},
		&hStudentJourneyRepo{},
		&hStudentLMSRepo{},
		schedulerRepo,
		currRepo,
		gradingRepo,
		nil,
		nil,
		nil,
		nil,
	)
	h := NewStudentHandler(svc)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/student/grades", nil)
	c.Set("tenant_id", "tenant-1")
	c.Set("userID", "stud-1")

	h.ListGrades(c)
	require.Equal(t, http.StatusOK, w.Code)

	var payload []dto.StudentGradeEntry
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &payload))
	require.Len(t, payload, 1)
	require.NotNil(t, payload[0].CourseCode)
	require.Equal(t, "MED101", *payload[0].CourseCode)
	require.Equal(t, gradedAt.Format(time.RFC3339), payload[0].GradedAt)
}

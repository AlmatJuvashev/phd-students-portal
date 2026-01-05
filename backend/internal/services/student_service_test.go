package services

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	pb "github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubStudentUserRepo struct {
	byID map[string]*models.User
}

func (s *stubStudentUserRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	if s.byID == nil {
		return nil, nil
	}
	return s.byID[id], nil
}

type stubStudentJourneyRepo struct {
	stateByUser map[string]map[string]string
	calls       int
}

func (s *stubStudentJourneyRepo) GetJourneyState(ctx context.Context, userID, tenantID string) (map[string]string, error) {
	s.calls++
	if s.stateByUser == nil {
		return map[string]string{}, nil
	}
	if state, ok := s.stateByUser[userID]; ok {
		return state, nil
	}
	return map[string]string{}, nil
}

type stubStudentLMSRepo struct {
	enrollmentsByStudent map[string][]models.CourseEnrollment
}

func (s *stubStudentLMSRepo) GetStudentEnrollments(ctx context.Context, studentID string) ([]models.CourseEnrollment, error) {
	if s.enrollmentsByStudent == nil {
		return nil, nil
	}
	return s.enrollmentsByStudent[studentID], nil
}

func (s *stubStudentLMSRepo) CreateSubmission(ctx context.Context, sub *models.ActivitySubmission) error {
	return nil
}

func (s *stubStudentLMSRepo) GetSubmissionByStudent(ctx context.Context, activityID, studentID string) (*models.ActivitySubmission, error) {
	return nil, nil
}

type stubStudentSchedulerRepo struct {
	offeringsByID map[string]*models.CourseOffering
	staffByOffer  map[string][]models.CourseStaff
	sessionsByOff map[string][]models.ClassSession
}

func (s *stubStudentSchedulerRepo) GetOffering(ctx context.Context, id string) (*models.CourseOffering, error) {
	if s.offeringsByID == nil {
		return nil, nil
	}
	return s.offeringsByID[id], nil
}

func (s *stubStudentSchedulerRepo) ListStaff(ctx context.Context, offeringID string) ([]models.CourseStaff, error) {
	if s.staffByOffer == nil {
		return nil, nil
	}
	return s.staffByOffer[offeringID], nil
}

func (s *stubStudentSchedulerRepo) ListSessions(ctx context.Context, offeringID string, startDate, endDate time.Time) ([]models.ClassSession, error) {
	if s.sessionsByOff == nil {
		return nil, nil
	}
	return s.sessionsByOff[offeringID], nil
}

type stubStudentCurriculumRepo struct {
	coursesByID map[string]*models.Course
}

func (s *stubStudentCurriculumRepo) GetCourse(ctx context.Context, id string) (*models.Course, error) {
	if s.coursesByID == nil {
		return nil, nil
	}
	return s.coursesByID[id], nil
}

type stubStudentGradingRepo struct {
	entriesByStudent map[string][]models.GradebookEntry
}

func (s *stubStudentGradingRepo) ListStudentEntries(ctx context.Context, studentID string) ([]models.GradebookEntry, error) {
	if s.entriesByStudent == nil {
		return nil, nil
	}
	return s.entriesByStudent[studentID], nil
}

func TestStudentService_GetDashboard_ComputesProgressAndUpcoming(t *testing.T) {
	pbm := &pb.Manager{
		DefaultLocale: "en",
		Nodes: map[string]pb.Node{
			"n1": {ID: "n1", Title: map[string]string{"en": "Profile"}},
			"n2": {ID: "n2", Title: map[string]string{"en": "Research Proposal"}, Prerequisites: []string{"n1"}},
			"n3": {ID: "n3", Title: map[string]string{"en": "Ethics Approval"}, Prerequisites: []string{"n1"}},
		},
		NodeWorlds: map[string]string{"n1": "W1", "n2": "W2", "n3": "W1"},
	}

	userRepo := &stubStudentUserRepo{
		byID: map[string]*models.User{
			"stud-1": {ID: "stud-1", Program: "PhD in Medicine"},
		},
	}
	journeyRepo := &stubStudentJourneyRepo{
		stateByUser: map[string]map[string]string{
			"stud-1": {"n1": "done", "n2": "submitted", "n3": "todo"},
		},
	}
	gradingRepo := &stubStudentGradingRepo{
		entriesByStudent: map[string][]models.GradebookEntry{"stud-1": {}},
	}

	svc := NewStudentService(
		userRepo,
		journeyRepo,
		&stubStudentLMSRepo{},
		&stubStudentSchedulerRepo{},
		&stubStudentCurriculumRepo{},
		gradingRepo,
		nil,
		nil,
		pbm,
	)

	res, err := svc.GetDashboard(context.Background(), "tenant-1", "stud-1")
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "PhD in Medicine", res.Program.Title)
	assert.Equal(t, 3, res.Program.TotalNodes)
	assert.Equal(t, 1, res.Program.CompletedNodes)
	assert.InDelta(t, 33.333, res.Program.ProgressPercent, 0.01)

	require.Len(t, res.UpcomingDeadlines, 2)
	assert.Equal(t, "n2", res.UpcomingDeadlines[0].ID)
	assert.Equal(t, "urgent", res.UpcomingDeadlines[0].Severity)
	assert.NotNil(t, res.UpcomingDeadlines[0].Link)
	assert.Equal(t, "/journey", *res.UpcomingDeadlines[0].Link)
}

func TestStudentService_ListCourses_MapsInstructorAndNextSession(t *testing.T) {
	now := time.Now()
	roomID := "room-1"
	meetingURL := "https://meet.example.com/room"
	nextDate := now.Add(48 * time.Hour)

	userRepo := &stubStudentUserRepo{
		byID: map[string]*models.User{
			"inst-1": {ID: "inst-1", FirstName: "Aida", LastName: "Baken", Email: "aida@example.com"},
		},
	}

	lmsRepo := &stubStudentLMSRepo{
		enrollmentsByStudent: map[string][]models.CourseEnrollment{
			"stud-1": {
				{ID: "enr-1", CourseOfferingID: "off-1", Status: "ENROLLED"},
			},
		},
	}

	schedulerRepo := &stubStudentSchedulerRepo{
		offeringsByID: map[string]*models.CourseOffering{
			"off-1": {ID: "off-1", CourseID: "course-1", TermID: "term-1", Section: "A", DeliveryFormat: "IN_PERSON"},
		},
		staffByOffer: map[string][]models.CourseStaff{
			"off-1": {
				{UserID: "inst-1", Role: "INSTRUCTOR", IsPrimary: true},
			},
		},
		sessionsByOff: map[string][]models.ClassSession{
			"off-1": {
				{ID: "sess-1", CourseOfferingID: "off-1", Date: nextDate, StartTime: "10:00", EndTime: "11:30", RoomID: &roomID, MeetingURL: &meetingURL, Type: "LECTURE"},
			},
		},
	}

	currRepo := &stubStudentCurriculumRepo{
		coursesByID: map[string]*models.Course{
			"course-1": {ID: "course-1", Code: "MED101", Title: `{"en":"Intro to Medicine"}`},
		},
	}

	svc := NewStudentService(
		userRepo,
		&stubStudentJourneyRepo{},
		lmsRepo,
		schedulerRepo,
		currRepo,
		&stubStudentGradingRepo{},
		nil,
		nil,
		nil,
	)

	list, err := svc.ListCourses(context.Background(), "tenant-1", "stud-1")
	require.NoError(t, err)
	require.Len(t, list, 1)

	c := list[0]
	assert.Equal(t, "enr-1", c.EnrollmentID)
	assert.Equal(t, "off-1", c.CourseOfferingID)
	assert.Equal(t, "course-1", c.CourseID)
	assert.Equal(t, "MED101", c.Code)
	assert.Equal(t, "Intro to Medicine", c.Title)
	require.NotNil(t, c.InstructorName)
	assert.Equal(t, "Aida Baken", *c.InstructorName)

	require.NotNil(t, c.NextSession)
	assert.Equal(t, "sess-1", c.NextSession.ID)
	assert.Equal(t, nextDate.Format("2006-01-02"), c.NextSession.Date)
	assert.Equal(t, "10:00", c.NextSession.StartTime)
	assert.Equal(t, "11:30", c.NextSession.EndTime)
	assert.Equal(t, "LECTURE", c.NextSession.Type)
	assert.Equal(t, &meetingURL, c.NextSession.MeetingURL)
}

func TestStudentService_ListGrades_EnrichesCourseInfo(t *testing.T) {
	gradedAt := time.Now()

	gradingRepo := &stubStudentGradingRepo{
		entriesByStudent: map[string][]models.GradebookEntry{
			"stud-1": {
				{
					ID:               "g1",
					CourseOfferingID: "off-1",
					ActivityID:       "act-1",
					StudentID:        "stud-1",
					Score:            90,
					MaxScore:         100,
					Grade:            "A",
					Feedback:         "Great work",
					GradedByID:       "inst-1",
					GradedAt:         gradedAt,
				},
			},
		},
	}

	schedulerRepo := &stubStudentSchedulerRepo{
		offeringsByID: map[string]*models.CourseOffering{
			"off-1": {ID: "off-1", CourseID: "course-1"},
		},
	}

	currRepo := &stubStudentCurriculumRepo{
		coursesByID: map[string]*models.Course{
			"course-1": {ID: "course-1", Code: "MED101", Title: `{"en":"Intro to Medicine"}`},
		},
	}

	svc := NewStudentService(
		&stubStudentUserRepo{},
		&stubStudentJourneyRepo{},
		&stubStudentLMSRepo{},
		schedulerRepo,
		currRepo,
		gradingRepo,
		nil,
		nil,
		nil,
	)

	list, err := svc.ListGrades(context.Background(), "tenant-1", "stud-1")
	require.NoError(t, err)
	require.Len(t, list, 1)

	g := list[0]
	assert.Equal(t, "g1", g.ID)
	assert.Equal(t, "off-1", g.CourseOfferingID)
	assert.Equal(t, "course-1", g.CourseID)
	require.NotNil(t, g.CourseCode)
	require.NotNil(t, g.CourseTitle)
	assert.Equal(t, "MED101", *g.CourseCode)
	assert.Equal(t, "Intro to Medicine", *g.CourseTitle)
	assert.Equal(t, "A", g.Grade)
	assert.Equal(t, gradedAt.Format(time.RFC3339), g.GradedAt)
}

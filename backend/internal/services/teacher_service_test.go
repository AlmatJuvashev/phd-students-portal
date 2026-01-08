package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks

type MockLMSRepo struct {
	mock.Mock
}

func (m *MockLMSRepo) EnrollStudent(ctx context.Context, e *models.CourseEnrollment) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockLMSRepo) GetCourseRoster(ctx context.Context, offeringID string) ([]models.CourseEnrollment, error) {
	args := m.Called(ctx, offeringID)
	return args.Get(0).([]models.CourseEnrollment), args.Error(1)
}
func (m *MockLMSRepo) GetStudentEnrollments(ctx context.Context, studentID string) ([]models.CourseEnrollment, error) {
	args := m.Called(ctx, studentID)
	return args.Get(0).([]models.CourseEnrollment), args.Error(1)
}
func (m *MockLMSRepo) UpdateEnrollmentStatus(ctx context.Context, id, status string) error {
	return m.Called(ctx, id, status).Error(0)
}
func (m *MockLMSRepo) CreateSubmission(ctx context.Context, s *models.ActivitySubmission) error {
	return m.Called(ctx, s).Error(0)
}
func (m *MockLMSRepo) GetSubmission(ctx context.Context, id string) (*models.ActivitySubmission, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ActivitySubmission), args.Error(1)
}
func (m *MockLMSRepo) GetSubmissionByStudent(ctx context.Context, activityID, studentID string) (*models.ActivitySubmission, error) {
	args := m.Called(ctx, activityID, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ActivitySubmission), args.Error(1)
}
func (m *MockLMSRepo) ListSubmissions(ctx context.Context, offeringID string) ([]models.ActivitySubmission, error) {
	args := m.Called(ctx, offeringID)
	return args.Get(0).([]models.ActivitySubmission), args.Error(1)
}
func (m *MockLMSRepo) MarkAttendance(ctx context.Context, att *models.ClassAttendance) error {
	return m.Called(ctx, att).Error(0)
}
func (m *MockLMSRepo) CreateAnnotation(ctx context.Context, ann models.SubmissionAnnotation) (*models.SubmissionAnnotation, error) {
	args := m.Called(ctx, ann)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SubmissionAnnotation), args.Error(1)
}
func (m *MockLMSRepo) ListAnnotations(ctx context.Context, submissionID string) ([]models.SubmissionAnnotation, error) {
	args := m.Called(ctx, submissionID)
	return args.Get(0).([]models.SubmissionAnnotation), args.Error(1)
}
func (m *MockLMSRepo) DeleteAnnotation(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockLMSRepo) GetSessionAttendance(ctx context.Context, sessionID string) ([]models.ClassAttendance, error) {
	args := m.Called(ctx, sessionID)
	return args.Get(0).([]models.ClassAttendance), args.Error(1)
}

// MockGradingRepo is already defined in grading_service_test.go
// Tests

func TestTeacherService_GetDashboardStats(t *testing.T) {
	mockSched := new(MockSchedulerRepo)
	mockLMS := new(MockLMSRepo)
	mockGrading := new(MockGradingRepo)
	svc := NewTeacherService(mockSched, mockLMS, mockGrading, nil)

	instructorID := "inst-123"

	// Mock ListSessionsByInstructor (Today)
	todaySessions := []models.ClassSession{
		{ID: "sess-1", StartTime: "09:00", Date: time.Now()}, // Passed? or future? Logic checks time.
	}
	// Logic matches startOfDay to endOfDay.
	mockSched.On("ListSessionsByInstructor", mock.Anything, instructorID, mock.Anything, mock.Anything).Return(todaySessions, nil)

	// Mock ListOfferings (Active Courses)
	offerings := []models.CourseOffering{
		{ID: "off-1", CourseID: "course-A"},
		{ID: "off-2", CourseID: "course-B"},
	}
	mockSched.On("ListOfferingsByInstructor", mock.Anything, instructorID, "").Return(offerings, nil)

	// Mock Pending Grading (ListSubmissions)
	// off-1 has 1 pending
	mockLMS.On("ListSubmissions", mock.Anything, "off-1").Return([]models.ActivitySubmission{
		{Status: "SUBMITTED"},
		{Status: "GRADED"},
	}, nil)
	// off-2 has 0 pending
	mockLMS.On("ListSubmissions", mock.Anything, "off-2").Return([]models.ActivitySubmission{}, nil)

	stats, err := svc.GetDashboardStats(context.Background(), instructorID)
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, 2, stats.ActiveCourses)
	assert.Equal(t, 1, stats.TodayClassesCount)
	assert.Equal(t, 1, stats.PendingGrading) // 1 pending from off-1
}

func TestTeacherService_GetMyCourses(t *testing.T) {
	mockSched := new(MockSchedulerRepo)
	svc := NewTeacherService(mockSched, new(MockLMSRepo), new(MockGradingRepo), nil)

	instructorID := "inst-123"
	expected := []models.CourseOffering{{ID: "off-1"}}
	mockSched.On("ListOfferingsByInstructor", mock.Anything, instructorID, "").Return(expected, nil)

	result, err := svc.GetMyCourses(context.Background(), instructorID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestTeacherService_GetCourseRoster(t *testing.T) {
	mockLMS := new(MockLMSRepo)
	svc := NewTeacherService(new(MockSchedulerRepo), mockLMS, new(MockGradingRepo), nil)

	offeringID := "off-1"
	expected := []models.CourseEnrollment{{Status: "ENROLLED"}}
	mockLMS.On("GetCourseRoster", mock.Anything, offeringID).Return(expected, nil)

	result, err := svc.GetCourseRoster(context.Background(), offeringID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestTeacherService_GetSubmissions(t *testing.T) {
	mockSched := new(MockSchedulerRepo)
	mockLMS := new(MockLMSRepo)
	svc := NewTeacherService(mockSched, mockLMS, new(MockGradingRepo), nil)

	instructorID := "inst-123"
	offerings := []models.CourseOffering{{ID: "off-1"}}
	mockSched.On("ListOfferingsByInstructor", mock.Anything, instructorID, "").Return(offerings, nil)

	expected := []models.ActivitySubmission{{ID: "sub-1"}}
	mockLMS.On("ListSubmissions", mock.Anything, "off-1").Return(expected, nil)

	result, err := svc.GetSubmissions(context.Background(), instructorID)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "sub-1", result[0].ID)
	mockLMS.AssertExpectations(t)
}

func TestTeacherService_Annotations(t *testing.T) {
	mockScheduler := new(MockSchedulerRepo)
	mockLMS := new(MockLMSRepo)
	mockGrading := new(MockGradingRepo)
	svc := NewTeacherService(mockScheduler, mockLMS, mockGrading, nil)

	ctx := context.Background()
	content := "Good job"
	ann := models.SubmissionAnnotation{
		SubmissionID: "sub-1",
		Content:      &content,
		Color:        "#FF0000",
	}

	mockLMS.On("CreateAnnotation", ctx, ann).Return(&ann, nil)
	mockLMS.On("ListAnnotations", ctx, "sub-1").Return([]models.SubmissionAnnotation{ann}, nil)
	mockLMS.On("DeleteAnnotation", ctx, "ann-1").Return(nil)

	// Create
	created, err := svc.AddAnnotation(ctx, ann)
	assert.NoError(t, err)
	assert.Equal(t, "sub-1", created.SubmissionID)

	// List
	list, err := svc.GetAnnotationsForSubmission(ctx, "sub-1")
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	// Delete
	err = svc.RemoveAnnotation(ctx, "ann-1")
	assert.NoError(t, err)

	mockLMS.AssertExpectations(t)
}

// MockContentRepo
type MockContentRepo struct {
	mock.Mock
}
func (m *MockContentRepo) GetModule(ctx context.Context, id string) (*models.CourseModule, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.CourseModule), args.Error(1)
}
func (m *MockContentRepo) GetLesson(ctx context.Context, id string) (*models.CourseLesson, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.CourseLesson), args.Error(1)
}
func (m *MockContentRepo) GetActivity(ctx context.Context, id string) (*models.CourseActivity, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.CourseActivity), args.Error(1)
}
func (m *MockContentRepo) ListModules(ctx context.Context, courseID string) ([]models.CourseModule, error) {
	args := m.Called(ctx, courseID)
	return args.Get(0).([]models.CourseModule), args.Error(1)
}
func (m *MockContentRepo) ListLessons(ctx context.Context, moduleID string) ([]models.CourseLesson, error) {
	args := m.Called(ctx, moduleID)
	return args.Get(0).([]models.CourseLesson), args.Error(1)
}
func (m *MockContentRepo) ListActivities(ctx context.Context, lessonID string) ([]models.CourseActivity, error) {
	args := m.Called(ctx, lessonID)
	return args.Get(0).([]models.CourseActivity), args.Error(1)
}

// Missing methods required by interface
func (m *MockContentRepo) CreateModule(ctx context.Context, mod *models.CourseModule) error {
	args := m.Called(ctx, mod)
	return args.Error(0)
}
func (m *MockContentRepo) UpdateModule(ctx context.Context, mod *models.CourseModule) error {
	args := m.Called(ctx, mod)
	return args.Error(0)
}
func (m *MockContentRepo) DeleteModule(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockContentRepo) CreateLesson(ctx context.Context, l *models.CourseLesson) error {
	args := m.Called(ctx, l)
	return args.Error(0)
}
func (m *MockContentRepo) UpdateLesson(ctx context.Context, l *models.CourseLesson) error {
	args := m.Called(ctx, l)
	return args.Error(0)
}
func (m *MockContentRepo) DeleteLesson(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockContentRepo) CreateActivity(ctx context.Context, a *models.CourseActivity) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}
func (m *MockContentRepo) UpdateActivity(ctx context.Context, a *models.CourseActivity) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}
func (m *MockContentRepo) DeleteActivity(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}


func TestTeacherService_GetCourseStudents_RiskLogic(t *testing.T) {
	mockSched := new(MockSchedulerRepo)
	mockLMS := new(MockLMSRepo)
	mockGrading := new(MockGradingRepo)
	mockContent := new(MockContentRepo)
	
	svc := NewTeacherService(mockSched, mockLMS, mockGrading, mockContent)

	offeringID := "off-1"
	
	// 1. Roster
	mockLMS.On("GetCourseRoster", mock.Anything, offeringID).Return([]models.CourseEnrollment{
		{StudentID: "s1", StudentName: "Student One"},
		{StudentID: "s2", StudentName: "Student Two"}, // Inactive
	}, nil)

	// 2. Content (Total Items calculation)
	mockSched.On("GetOffering", mock.Anything, offeringID).Return(&models.CourseOffering{CourseID: "c1"}, nil)
	mockContent.On("ListModules", mock.Anything, "c1").Return([]models.CourseModule{{ID: "m1"}}, nil)
	mockContent.On("ListLessons", mock.Anything, "m1").Return([]models.CourseLesson{{ID: "l1"}}, nil)
	mockContent.On("ListActivities", mock.Anything, "l1").Return([]models.CourseActivity{
		{ID: "a1", Type: "assignment"},
		{ID: "a2", Type: "assignment"},
	}, nil)
	
	// 3. Submissions
	lastWeek := time.Now().Add(-7 * 24 * time.Hour)
	longAgo := time.Now().Add(-35 * 24 * time.Hour)

	mockLMS.On("ListSubmissions", mock.Anything, offeringID).Return([]models.ActivitySubmission{
		{StudentID: "s1", ActivityID: "a1", Status: "SUBMITTED", SubmittedAt: lastWeek}, // 1/2 done
		{StudentID: "s2", ActivityID: "a1", Status: "SUBMITTED", SubmittedAt: longAgo},  // Inactive
	}, nil)

	// 4. Grades
	mockGrading.On("ListEntries", mock.Anything, offeringID).Return([]models.GradebookEntry{
		{StudentID: "s1", Score: 90, MaxScore: 100},
		{StudentID: "s2", Score: 40, MaxScore: 100},
	}, nil)

	profiles, err := svc.GetCourseStudents(context.Background(), offeringID)
	assert.NoError(t, err)
	assert.Len(t, profiles, 2)

	// Check s1 (Active, Good Grade)
	// We iterate roster, so order should match roster if not map iteration
	// But it uses studentSubmissionCounts iteration? No, iterates roster.
	p1 := profiles[0]
	assert.Equal(t, "s1", p1.StudentID)
	assert.Equal(t, 50.0, p1.OverallProgress)
	assert.Equal(t, 90.0, p1.AverageGrade)
	assert.Equal(t, "medium", p1.RiskLevel) // Inactive 7 days -> medium

	// Check s2 (Inactive, Bad Grade)
	p2 := profiles[1]
	assert.Equal(t, "s2", p2.StudentID)
	assert.Equal(t, 40.0, p2.AverageGrade)
	assert.Equal(t, "critical", p2.RiskLevel) // Inactive 30+ days -> critical
	assert.Contains(t, p2.RiskFactors, "35 days inactive")
}

func TestTeacherService_GetStudentActivity(t *testing.T) {
	mockLMS := new(MockLMSRepo)
	mockGrading := new(MockGradingRepo)
	svc := NewTeacherService(nil, mockLMS, mockGrading, nil)

	offID := "off-1"
	studID := "s1"
	now := time.Now()

	mockLMS.On("ListSubmissions", mock.Anything, offID).Return([]models.ActivitySubmission{
		{ID: "sub1", StudentID: studID, ActivityTitle: "Lab 1", SubmittedAt: now},
	}, nil)

	mockGrading.On("ListEntries", mock.Anything, offID).Return([]models.GradebookEntry{
		{StudentID: studID, ActivityID: "act1", Score: 95, MaxScore: 100, GradedAt: now.Add(time.Hour)},
	}, nil)

	events, err := svc.GetStudentActivity(context.Background(), studID, offID, 10)
	assert.NoError(t, err)
	assert.Len(t, events, 2)
	assert.Equal(t, "grade", events[0].Kind) // newer first
	assert.Equal(t, "submission", events[1].Kind)
}

func TestTeacherService_GetGradebook(t *testing.T) {
	mockGrading := new(MockGradingRepo)
	svc := NewTeacherService(nil, nil, mockGrading, nil)

	offeringID := "off-1"
	expected := []models.GradebookEntry{
		{StudentID: "s1", ActivityID: "a1", Score: 90},
	}
	mockGrading.On("ListEntries", mock.Anything, offeringID).Return(expected, nil)

	result, err := svc.GetGradebook(context.Background(), offeringID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestTeacherService_GetCourseAtRisk(t *testing.T) {
	mockSched := new(MockSchedulerRepo)
	mockLMS := new(MockLMSRepo)
	mockGrading := new(MockGradingRepo)
	mockContent := new(MockContentRepo) // Need content repo to avoid nil check in countOfferingWorkItems

	svc := NewTeacherService(mockSched, mockLMS, mockGrading, mockContent)

	offeringID := "off-1"

	// 1. Roster
	mockLMS.On("GetCourseRoster", mock.Anything, offeringID).Return([]models.CourseEnrollment{
		{StudentID: "safe", StudentName: "Safe Student"},
		{StudentID: "risk", StudentName: "Risky Student"},
	}, nil)

	// 2. Content (Total Items = 10)
	mockSched.On("GetOffering", mock.Anything, offeringID).Return(&models.CourseOffering{CourseID: "c1"}, nil)
	mockContent.On("ListModules", mock.Anything, "c1").Return([]models.CourseModule{{ID: "m1"}}, nil)
	mockContent.On("ListLessons", mock.Anything, "m1").Return([]models.CourseLesson{{ID: "l1"}}, nil)
	// 10 assignments
	var acts []models.CourseActivity
	for i := 0; i < 10; i++ {
		acts = append(acts, models.CourseActivity{ID: fmt.Sprintf("a%d", i), Type: "assignment"})
	}
	mockContent.On("ListActivities", mock.Anything, "l1").Return(acts, nil)

	// 3. Submissions
	// Safe student: Submitted all 10
	// Risky student: Submitted 0
	var subs []models.ActivitySubmission
	for i := 0; i < 10; i++ {
		subs = append(subs, models.ActivitySubmission{StudentID: "safe", ActivityID: fmt.Sprintf("a%d", i), Status: "SUBMITTED", SubmittedAt: time.Now()})
	}
	mockLMS.On("ListSubmissions", mock.Anything, offeringID).Return(subs, nil)

	// 4. Grades
	mockGrading.On("ListEntries", mock.Anything, offeringID).Return([]models.GradebookEntry{}, nil)

	// Run
	profiles, err := svc.GetCourseAtRisk(context.Background(), offeringID)
	assert.NoError(t, err)
	
	// Expect only "risk" student
	assert.Len(t, profiles, 1)
	assert.Equal(t, "risk", profiles[0].StudentID)
	assert.Contains(t, []string{"high", "critical"}, profiles[0].RiskLevel)
}

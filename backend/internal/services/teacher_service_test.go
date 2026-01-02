package services

import (
	"context"
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
	svc := NewTeacherService(mockSched, mockLMS, mockGrading)

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
	svc := NewTeacherService(mockSched, new(MockLMSRepo), new(MockGradingRepo))

	instructorID := "inst-123"
	expected := []models.CourseOffering{{ID: "off-1"}}
	mockSched.On("ListOfferingsByInstructor", mock.Anything, instructorID, "").Return(expected, nil)

	result, err := svc.GetMyCourses(context.Background(), instructorID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestTeacherService_GetCourseRoster(t *testing.T) {
	mockLMS := new(MockLMSRepo)
	svc := NewTeacherService(new(MockSchedulerRepo), mockLMS, new(MockGradingRepo))

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
	svc := NewTeacherService(mockSched, mockLMS, new(MockGradingRepo))

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
	svc := NewTeacherService(mockScheduler, mockLMS, mockGrading)

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


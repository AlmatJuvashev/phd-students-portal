package services

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type RiskMockAnalyticsRepo struct {
	mock.Mock
}

func (m *RiskMockAnalyticsRepo) GetMonitorMetrics(ctx context.Context, filter models.FilterParams) (models.MonitorMetrics, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(models.MonitorMetrics), args.Error(1)
}
func (m *RiskMockAnalyticsRepo) GetStudentsByStage(ctx context.Context) ([]models.StudentStageStats, error) {
	return nil, nil
}
func (m *RiskMockAnalyticsRepo) GetAdvisorLoad(ctx context.Context) ([]models.AdvisorLoadStats, error) {
	return nil, nil
}
func (m *RiskMockAnalyticsRepo) GetOverdueTasks(ctx context.Context) ([]models.OverdueTaskStats, error) {
	return nil, nil
}

// Implement new methods
func (m *RiskMockAnalyticsRepo) SaveRiskSnapshot(ctx context.Context, s *models.RiskSnapshot) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}
func (m *RiskMockAnalyticsRepo) GetStudentRiskHistory(ctx context.Context, studentID string) ([]models.RiskSnapshot, error) {
	return nil, nil
}
func (m *RiskMockAnalyticsRepo) GetHighRiskStudents(ctx context.Context, threshold float64) ([]models.RiskSnapshot, error) {
	return nil, nil
}
func (m *RiskMockAnalyticsRepo) GetTotalStudents(ctx context.Context, filter models.FilterParams) (int, error) {
    return 0, nil
}
func (m *RiskMockAnalyticsRepo) GetNodeCompletionCount(ctx context.Context, nodeID string, filter models.FilterParams) (int, error) {
    return 0, nil
}
func (m *RiskMockAnalyticsRepo) GetDurationForNodes(ctx context.Context, nodeIDs []string, filter models.FilterParams) ([]float64, error) {
    return nil, nil
}
func (m *RiskMockAnalyticsRepo) GetBottleneck(ctx context.Context, filter models.FilterParams) (string, int, error) {
    return "", 0, nil
}
func (m *RiskMockAnalyticsRepo) GetProfileFlagCount(ctx context.Context, flagColumn string, threshold float64, filter models.FilterParams) (int, error) {
    return 0, nil
}


type RiskMockAttendanceRepo struct {
	mock.Mock
}

func (m *RiskMockAttendanceRepo) RecordAttendance(ctx context.Context, sessionID string, record models.ClassAttendance) error {
	return nil
}

func (m *RiskMockAttendanceRepo) BatchUpsertAttendance(ctx context.Context, sessionID string, records []models.ClassAttendance, recordedBy string) error {
	return nil
}
func (m *RiskMockAttendanceRepo) GetSessionAttendance(ctx context.Context, sessionID string) ([]models.ClassAttendance, error) {
	return nil, nil
}
func (m *RiskMockAttendanceRepo) GetStudentAttendance(ctx context.Context, studentID string) ([]models.ClassAttendance, error) {
	args := m.Called(ctx, studentID)
	// Return nil or generic
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ClassAttendance), args.Error(1)
}

type RiskMockLMSRepo struct {
	mock.Mock
}

func (m *RiskMockLMSRepo) EnrollStudent(ctx context.Context, enrollment *models.CourseEnrollment) error { return nil }
func (m *RiskMockLMSRepo) GetCourseRoster(ctx context.Context, offeringID string) ([]models.CourseEnrollment, error) { return nil, nil }
func (m *RiskMockLMSRepo) GetStudentEnrollments(ctx context.Context, studentID string) ([]models.CourseEnrollment, error) { return nil, nil }
func (m *RiskMockLMSRepo) UpdateEnrollmentStatus(ctx context.Context, id, status string) error { return nil }
func (m *RiskMockLMSRepo) CreateSubmission(ctx context.Context, sub *models.ActivitySubmission) error { return nil }
func (m *RiskMockLMSRepo) GetSubmission(ctx context.Context, id string) (*models.ActivitySubmission, error) { return nil, nil }
func (m *RiskMockLMSRepo) GetSubmissionByStudent(ctx context.Context, activityID, studentID string) (*models.ActivitySubmission, error) {
    args := m.Called(ctx, activityID, studentID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.ActivitySubmission), args.Error(1)
}
func (m *RiskMockLMSRepo) ListSubmissions(ctx context.Context, activityID string) ([]models.ActivitySubmission, error) { return nil, nil }
func (m *RiskMockLMSRepo) CreateAnnotation(ctx context.Context, a models.SubmissionAnnotation) (*models.SubmissionAnnotation, error) { return &a, nil }
func (m *RiskMockLMSRepo) ListAnnotations(ctx context.Context, submissionID string) ([]models.SubmissionAnnotation, error) { return nil, nil }
func (m *RiskMockLMSRepo) DeleteAnnotation(ctx context.Context, id string) error { return nil }
func (m *RiskMockLMSRepo) MarkAttendance(ctx context.Context, att *models.ClassAttendance) error { return nil }
func (m *RiskMockLMSRepo) GetSessionAttendance(ctx context.Context, sessionID string) ([]models.ClassAttendance, error) { return nil, nil }

func TestAnalyticsService_CalculateStudentRisk(t *testing.T) {
	mockRepo := new(RiskMockAnalyticsRepo)
	mockAtt := new(RiskMockAttendanceRepo)
	mockLMS := new(RiskMockLMSRepo)

	svc := NewAnalyticsService(mockRepo, mockLMS, mockAtt, nil)
	ctx := context.Background()
	studentID := "student-1"

	t.Run("High Attendance Risk", func(t *testing.T) {
		// Mock Attendance: 4 classes, 2 ABSENT -> 50% attendance
		mockAtt.On("GetStudentAttendance", ctx, studentID).Return([]models.ClassAttendance{
			{Status: models.AttendancePresent},
			{Status: models.AttendancePresent},
			{Status: models.AttendanceAbsent},
			{Status: models.AttendanceAbsent},
		}, nil)

		// Mock LMS for Grades (placeholder logic uses GetSubmissionByStudent)
		// But in current "CalculateStudentRisk" we ignore submissions result. 
		// Just need mock to not panic if called.
		// However, I changed code to call GetSubmissionByStudent.
		// Wait, "GetSubmissionByStudent" args? (ctx, activityID, studentID).
		// In previous logic I put "" as activityID.
		mockLMS.On("GetSubmissionByStudent", ctx, "", studentID).Return(&models.ActivitySubmission{}, nil)
		mockLMS.On("GetProfileFlagCount", ctx, studentID).Return(map[string]int{}, nil)

		risk, err := svc.CalculateStudentRisk(ctx, studentID)
		assert.NoError(t, err)
		assert.NotNil(t, risk)
		assert.Equal(t, studentID, risk.StudentID)
		assert.GreaterOrEqual(t, risk.RiskScore, 25.0)

		mockAtt.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Low Risk - Good Attendance", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockAtt.ExpectedCalls = nil
		mockLMS.ExpectedCalls = nil

		mockAtt.On("GetStudentAttendance", ctx, studentID).Return([]models.ClassAttendance{
			{Status: models.AttendancePresent},
			{Status: models.AttendancePresent},
			{Status: models.AttendancePresent}, // 100%
		}, nil)
		mockLMS.On("GetSubmissionByStudent", ctx, "", studentID).Return(&models.ActivitySubmission{}, nil)
		
		mockRepo.On("SaveRiskSnapshot", ctx, mock.MatchedBy(func(s *models.RiskSnapshot) bool {
			return s.RiskScore < 10.0 // Should be 0 from attendance
		})).Return(nil)

		risk, err := svc.CalculateStudentRisk(ctx, studentID)
		assert.NoError(t, err)
		assert.Less(t, risk.RiskScore, 10.0) // Just strictly less than low threshold
	})
}

func TestAnalyticsService_RunBatchRiskAnalysis(t *testing.T) {
	// 1. Mocks
	mockRepo := new(RiskMockAnalyticsRepo)
	mockLMS := new(RiskMockLMSRepo)
	mockAtt := new(RiskMockAttendanceRepo)
	mockUser := new(MockUserRepository) // Using existing Mock from common_mocks_test.go

	// 2. Setup
	svc := NewAnalyticsService(mockRepo, mockLMS, mockAtt, mockUser)
	ctx := context.Background()

	// 3. User Mock
	stubUser := models.User{ID: "u1", Role: "student", IsActive: true}
	
	active := true
	filter := repository.UserFilter{Role: "student", Active: &active}
	// Note: Pagination Check. First call returns 1 user.
	// Next loop iteration calculates offset += 100.
	// Implementation calls List again?
	// The implementation loop:
	// users, total, err := s.userRepo.List(...)
	// if len(users) == 0 -> break.
	// process users.
	// offset += limit.
	// if offset >= total -> break.
	
	// So if List returns total=1, limit=100. offset starts 0.
	// loop 1: List(offset=0) -> returns [u1], total=1. Process u1. offset becomes 100.
	// check 100 >= 1? Yes. Break.
	// So List Is called ONCE.
	
	mockUser.On("List", ctx, filter, repository.Pagination{Limit: 100, Offset: 0}).
		Return([]models.User{stubUser}, 1, nil)

	// 4. Calc Risk Mocks
	// Attendance
	rawAtt := []models.ClassAttendance{
		{Status: models.AttendancePresent}, {Status: models.AttendanceAbsent},
	}
	mockAtt.On("GetStudentAttendance", ctx, "u1").Return(rawAtt, nil)

	// LMS
	mockLMS.On("GetProfileFlagCount", ctx, "u1").Return(map[string]int{"rp_required": 0}, nil)
	mockLMS.On("GetSubmissionByStudent", ctx, "", "u1").Return(&models.ActivitySubmission{}, nil)
	// (Note: CalculateStudentRisk logic might call GetSubmissionByStudent)

	// 5. Save Mock
	mockRepo.On("SaveRiskSnapshot", ctx, mock.MatchedBy(func(s *models.RiskSnapshot) bool {
		return s.StudentID == "u1" && s.RiskScore == 25.0
	})).Return(nil)

	// 6. Execute
	count, err := svc.RunBatchRiskAnalysis(ctx)

	// 7. Verify
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
	
	mockUser.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockAtt.AssertExpectations(t)
}

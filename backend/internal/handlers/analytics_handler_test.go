package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAnalyticsRepo struct {
	mock.Mock
}

func (m *mockAnalyticsRepo) GetStudentsByStage(ctx context.Context) ([]models.StudentStageStats, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.StudentStageStats), args.Error(1)
}
func (m *mockAnalyticsRepo) GetAdvisorLoad(ctx context.Context) ([]models.AdvisorLoadStats, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.AdvisorLoadStats), args.Error(1)
}
func (m *mockAnalyticsRepo) GetOverdueTasks(ctx context.Context) ([]models.OverdueTaskStats, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.OverdueTaskStats), args.Error(1)
}
func (m *mockAnalyticsRepo) GetTotalStudents(ctx context.Context, filter models.FilterParams) (int, error) {
	args := m.Called(ctx, filter)
	return args.Int(0), args.Error(1)
}
func (m *mockAnalyticsRepo) GetNodeCompletionCount(ctx context.Context, nodeID string, filter models.FilterParams) (int, error) {
	args := m.Called(ctx, nodeID, filter)
	return args.Int(0), args.Error(1)
}
func (m *mockAnalyticsRepo) GetDurationForNodes(ctx context.Context, nodeIDs []string, filter models.FilterParams) ([]float64, error) {
	args := m.Called(ctx, nodeIDs, filter)
	return args.Get(0).([]float64), args.Error(1)
}
func (m *mockAnalyticsRepo) GetBottleneck(ctx context.Context, filter models.FilterParams) (string, int, error) {
	args := m.Called(ctx, filter)
	return args.String(0), args.Int(1), args.Error(2)
}
func (m *mockAnalyticsRepo) GetProfileFlagCount(ctx context.Context, key string, minVal float64, filter models.FilterParams) (int, error) {
	args := m.Called(ctx, key, minVal, filter)
	return args.Int(0), args.Error(1)
}
func (m *mockAnalyticsRepo) SaveRiskSnapshot(ctx context.Context, s *models.RiskSnapshot) error {
	return m.Called(ctx, s).Error(0)
}
func (m *mockAnalyticsRepo) GetStudentRiskHistory(ctx context.Context, studentID string) ([]models.RiskSnapshot, error) {
	args := m.Called(ctx, studentID)
	return args.Get(0).([]models.RiskSnapshot), args.Error(1)
}
func (m *mockAnalyticsRepo) GetHighRiskStudents(ctx context.Context, threshold float64) ([]models.RiskSnapshot, error) {
	args := m.Called(ctx, threshold)
	return args.Get(0).([]models.RiskSnapshot), args.Error(1)
}
func (m *mockAnalyticsRepo) GetSessionAttendance(ctx context.Context, sessionID string) ([]models.ClassAttendance, error) {
	args := m.Called(ctx, sessionID)
	return args.Get(0).([]models.ClassAttendance), args.Error(1)
}

func TestAnalyticsHandler_GetMonitorMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mockAnalyticsRepo)
	svc := services.NewAnalyticsService(mockRepo, nil, nil, nil)
	h := NewAnalyticsHandler(svc)

	filter := models.FilterParams{TenantID: "t1"}

	mockRepo.On("GetTotalStudents", mock.Anything, filter).Return(10, nil)
	mockRepo.On("GetNodeCompletionCount", mock.Anything, mock.Anything, filter).Return(5, nil)
	mockRepo.On("GetDurationForNodes", mock.Anything, mock.Anything, filter).Return([]float64{1.0}, nil)
	mockRepo.On("GetBottleneck", mock.Anything, filter).Return("node1", 1, nil)
	mockRepo.On("GetProfileFlagCount", mock.Anything, mock.Anything, mock.Anything, filter).Return(2, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/monitor", nil)
	c.Set("tenant_id", "t1")

	h.GetMonitorMetrics(c)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(10), resp["total_students_count"])
}

func TestAnalyticsHandler_GetHighRiskStudents(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mockAnalyticsRepo)
	svc := services.NewAnalyticsService(mockRepo, nil, nil, nil)
	h := NewAnalyticsHandler(svc)

	mockRepo.On("GetHighRiskStudents", mock.Anything, 50.0).Return([]models.RiskSnapshot{}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/high-risk?threshold=50.0", nil)

	h.GetHighRiskStudents(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAnalyticsHandler_GetStageStats(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mockAnalyticsRepo)
	svc := services.NewAnalyticsService(mockRepo, nil, nil, nil)
	h := NewAnalyticsHandler(svc)

	mockRepo.On("GetStudentsByStage", mock.Anything).Return([]models.StudentStageStats{
		{Stage: "Research", Count: 5},
	}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/stages", nil)

	h.GetStageStats(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Research")
}

func TestAnalyticsHandler_GetOverdueStats(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mockAnalyticsRepo)
	svc := services.NewAnalyticsService(mockRepo, nil, nil, nil)
	h := NewAnalyticsHandler(svc)

	mockRepo.On("GetOverdueTasks", mock.Anything).Return([]models.OverdueTaskStats{
		{NodeID: "Submission", Count: 3},
	}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/overdue", nil)

	h.GetOverdueStats(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Submission")
}

// Minimal mocks for dependencies
type mockUserRepo struct { mock.Mock }
func (m *mockUserRepo) Create(ctx context.Context, u *models.User) (string, error) { return "", m.Called(ctx, u).Error(0) }
func (m *mockUserRepo) GetByID(ctx context.Context, id string) (*models.User, error) { args := m.Called(ctx, id); return args.Get(0).(*models.User), args.Error(1) }
func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) { args := m.Called(ctx, email); return args.Get(0).(*models.User), args.Error(1) }
func (m *mockUserRepo) Update(ctx context.Context, u *models.User) error { return m.Called(ctx, u).Error(0) }
func (m *mockUserRepo) Delete(ctx context.Context, id string) error { return m.Called(ctx, id).Error(0) }
func (m *mockUserRepo) List(ctx context.Context, filter repository.UserFilter, pagination repository.Pagination) ([]models.User, int, error) {
	args := m.Called(ctx, filter, pagination)
	return args.Get(0).([]models.User), args.Int(1), args.Error(2)
}
func (m *mockUserRepo) AssignRole(ctx context.Context, userID, roleID string) error { return m.Called(ctx, userID, roleID).Error(0) }
func (m *mockUserRepo) RemoveRole(ctx context.Context, userID, roleID string) error { return m.Called(ctx, userID, roleID).Error(0) }
func (m *mockUserRepo) GetRoles(ctx context.Context, userID string) ([]models.Role, error) { args := m.Called(ctx, userID); return args.Get(0).([]models.Role), args.Error(1) }
func (m *mockUserRepo) UpdateLastLogin(ctx context.Context, userID string) error { return m.Called(ctx, userID).Error(0) }
func (m *mockUserRepo) CheckRateLimit(ctx context.Context, userID, action string, window time.Duration) (int, error) { return 0, nil }
func (m *mockUserRepo) RecordRateLimit(ctx context.Context, userID, action string) error { return nil }
func (m *mockUserRepo) CreatePasswordResetToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error { return nil }
func (m *mockUserRepo) GetPasswordResetToken(ctx context.Context, tokenHash string) (string, time.Time, error) { return "", time.Time{}, nil }
func (m *mockUserRepo) DeletePasswordResetToken(ctx context.Context, tokenHash string) error { return nil }
func (m *mockUserRepo) GetTenantRoles(ctx context.Context, userID, tenantID string) ([]string, error) { return nil, nil }
func (m *mockUserRepo) LinkAdvisor(ctx context.Context, studentID, advisorID, tenantID string) error { return nil }
func (m *mockUserRepo) ReplaceAdvisors(ctx context.Context, studentID string, advisorIDs []string, tenantID string) error { return nil }
func (m *mockUserRepo) CreateEmailVerificationToken(ctx context.Context, userID, newEmail, token string, expiresAt time.Time) error { return nil }
func (m *mockUserRepo) GetEmailVerificationToken(ctx context.Context, token string) (string, string, string, error) { return "", "", "", nil }
func (m *mockUserRepo) DeleteEmailVerificationToken(ctx context.Context, token string) error { return nil }
func (m *mockUserRepo) GetPendingEmailVerification(ctx context.Context, userID string) (string, error) { return "", nil }
func (m *mockUserRepo) LogProfileAudit(ctx context.Context, userID, field, oldValue, newValue, changedBy string) error { return nil }
func (m *mockUserRepo) SyncProfileSubmissions(ctx context.Context, userID string, formData map[string]string, tenantID string) error { return nil }
func (m *mockUserRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) { args := m.Called(ctx, username); return args.Get(0).(*models.User), args.Error(1) }
func (m *mockUserRepo) UpdatePassword(ctx context.Context, id string, hash string) error { return nil }
func (m *mockUserRepo) UpdateAvatar(ctx context.Context, id string, avatarURL string) error { return nil }
func (m *mockUserRepo) SetActive(ctx context.Context, id string, active bool) error { return nil }
func (m *mockUserRepo) Exists(ctx context.Context, username string) (bool, error) { return true, nil }
func (m *mockUserRepo) EmailExists(ctx context.Context, email string, excludeID string) (bool, error) { return true, nil }
func (m *mockUserRepo) GetUserRoles(ctx context.Context, userID string) ([]string, error) { return nil, nil }

type mockAttRepo struct { mock.Mock }
func (m *mockAttRepo) CreateAttendance(ctx context.Context, a *models.ClassAttendance) error { return m.Called(ctx, a).Error(0) }
func (m *mockAttRepo) GetStudentAttendance(ctx context.Context, sID string) ([]models.ClassAttendance, error) { args := m.Called(ctx, sID); return args.Get(0).([]models.ClassAttendance), args.Error(1) }
func (m *mockAttRepo) GetSessionAttendance(ctx context.Context, sID string) ([]models.ClassAttendance, error) { args := m.Called(ctx, sID); return args.Get(0).([]models.ClassAttendance), args.Error(1) }
func (m *mockAttRepo) UpdateAttendance(ctx context.Context, a *models.ClassAttendance) error { return m.Called(ctx, a).Error(0) }
func (m *mockAttRepo) BatchUpsertAttendance(ctx context.Context, sID string, recs []models.ClassAttendance, by string) error { return nil }
func (m *mockAttRepo) RecordAttendance(ctx context.Context, sID string, rec models.ClassAttendance) error { return nil }

type mockLMSRepo struct { mock.Mock }
// Implement minimal methods for Analysis
func (m *mockLMSRepo) GetSubmissionByStudent(ctx context.Context, aID, sID string) (*models.ActivitySubmission, error) { 
	args := m.Called(ctx, aID, sID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.ActivitySubmission), args.Error(1) 
}
// Stubs for interface compliance (huge interface, adding minimal panic preventers or assumes usage restricted)
// If NewAnalyticsService only uses GetSubmissionByStudent, this might be enough if I cast it properly?
// But NewAnalyticsService expects repository.LMSRepository interface.
// Implementation must satisfy ALL methods.
// This is painful manually.
// I will try to use a "embedded" repo struct if possible or just implement all stubs.
// Or I can just pass `nil` if I mock the *CalculateStudentRisk* method? No, unlikely.
// I'll implement stubs.
func (m *mockLMSRepo) EnrollStudent(ctx context.Context, e *models.CourseEnrollment) error { return nil }
func (m *mockLMSRepo) GetCourseRoster(ctx context.Context, offeringID string) ([]models.CourseEnrollment, error) { return nil, nil }
func (m *mockLMSRepo) GetStudentEnrollments(ctx context.Context, sID string) ([]models.CourseEnrollment, error) { return nil, nil }
func (m *mockLMSRepo) UpdateEnrollmentStatus(ctx context.Context, id, s string) error { return nil }
func (m *mockLMSRepo) CreateSubmission(ctx context.Context, s *models.ActivitySubmission) error { return nil }
func (m *mockLMSRepo) GetSubmission(ctx context.Context, id string) (*models.ActivitySubmission, error) { return nil, nil }
func (m *mockLMSRepo) ListSubmissions(ctx context.Context, oID string) ([]models.ActivitySubmission, error) { return nil, nil }
func (m *mockLMSRepo) MarkAttendance(ctx context.Context, att *models.ClassAttendance) error { return nil }
func (m *mockLMSRepo) CreateAnnotation(ctx context.Context, ann models.SubmissionAnnotation) (*models.SubmissionAnnotation, error) { return nil, nil }
func (m *mockLMSRepo) ListAnnotations(ctx context.Context, reqID string) ([]models.SubmissionAnnotation, error) { return nil, nil }
func (m *mockLMSRepo) DeleteAnnotation(ctx context.Context, id string) error { return nil }
func (m *mockLMSRepo) GetSessionAttendance(ctx context.Context, sID string) ([]models.ClassAttendance, error) { return nil, nil }

func TestAnalyticsHandler_HandleBatchRiskAnalysis(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mockAnalyticsRepo)
	mockUser := new(mockUserRepo)
	mockAtt := new(mockAttRepo)
	mockLMS := new(mockLMSRepo)
	
	// Service needs UserRepo to list students, and AttRepo to calc risk (default implementation)
	svc := services.NewAnalyticsService(mockRepo, mockLMS, mockAtt, mockUser)
	h := NewAnalyticsHandler(svc)

	// 1. List Users returns 1 student
	mockUser.On("List", mock.Anything, mock.Anything, mock.Anything).Return([]models.User{{ID: "u1"}}, 1, nil).Once()
	// Loop calls List again for next page -> returns empty
	mockUser.On("List", mock.Anything, mock.Anything, mock.Anything).Return([]models.User{}, 1, nil).Once()

	// 2. Calculate Risk calls AttRepo
	mockAtt.On("GetStudentAttendance", mock.Anything, "u1").Return([]models.ClassAttendance{}, nil)
	// And LMS Repo (GetSubmissionByStudent)
	mockLMS.On("GetSubmissionByStudent", mock.Anything, "", "u1").Return(nil, nil)
	
	// 3. Save Snapshot
	// (And creates snapshot)

	// 3. Save Snapshot
	mockRepo.On("SaveRiskSnapshot", mock.Anything, mock.Anything).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/batch-risk", nil)

	h.HandleBatchRiskAnalysis(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Batch analysis completed")
}

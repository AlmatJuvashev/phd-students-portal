package services

import (
	"context"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository implements repository.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}
func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) { ret := m.Called(ctx, email); return ret.Get(0).(*models.User), ret.Error(1) }
func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) { ret := m.Called(ctx, username); return ret.Get(0).(*models.User), ret.Error(1) }
func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error { return m.Called(ctx, user).Error(0) }
func (m *MockUserRepository) UpdatePassword(ctx context.Context, id, hash string) error { return m.Called(ctx, id, hash).Error(0) }
func (m *MockUserRepository) UpdateAvatar(ctx context.Context, id, avatarURL string) error { return m.Called(ctx, id, avatarURL).Error(0) }
func (m *MockUserRepository) SetActive(ctx context.Context, id string, active bool) error { return m.Called(ctx, id, active).Error(0) }
func (m *MockUserRepository) Exists(ctx context.Context, username string) (bool, error) { ret := m.Called(ctx, username); return ret.Bool(0), ret.Error(1) }
func (m *MockUserRepository) EmailExists(ctx context.Context, email, excludeUserID string) (bool, error) { ret := m.Called(ctx, email, excludeUserID); return ret.Bool(0), ret.Error(1) }
func (m *MockUserRepository) List(ctx context.Context, filter repository.UserFilter, pagination repository.Pagination) ([]models.User, int, error) { 
	ret := m.Called(ctx, filter, pagination)
	return ret.Get(0).([]models.User), ret.Int(1), ret.Error(2)
}
func (m *MockUserRepository) CreatePasswordResetToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error { return m.Called(ctx, userID, tokenHash, expiresAt).Error(0) }
func (m *MockUserRepository) GetPasswordResetToken(ctx context.Context, tokenHash string) (string, time.Time, error) { 
	ret := m.Called(ctx, tokenHash)
	return ret.String(0), ret.Get(1).(time.Time), ret.Error(2)
}
func (m *MockUserRepository) DeletePasswordResetToken(ctx context.Context, tokenHash string) error { return m.Called(ctx, tokenHash).Error(0) }
func (m *MockUserRepository) GetTenantRoles(ctx context.Context, userID, tenantID string) ([]string, error) { 
	ret := m.Called(ctx, userID, tenantID)
	if ret.Get(0) == nil { return nil, ret.Error(1) }
	return ret.Get(0).([]string), ret.Error(1)
}
func (m *MockUserRepository) LinkAdvisor(ctx context.Context, studentID, advisorID, tenantID string) error { return m.Called(ctx, studentID, advisorID, tenantID).Error(0) }
func (m *MockUserRepository) ReplaceAdvisors(ctx context.Context, studentID string, advisorIDs []string, tenantID string) error { return m.Called(ctx, studentID, advisorIDs, tenantID).Error(0) }
func (m *MockUserRepository) CheckRateLimit(ctx context.Context, userID, action string, window time.Duration) (int, error) { 
	ret := m.Called(ctx, userID, action, window)
	return ret.Int(0), ret.Error(1)
}
func (m *MockUserRepository) RecordRateLimit(ctx context.Context, userID, action string) error { return m.Called(ctx, userID, action).Error(0) }
func (m *MockUserRepository) CreateEmailVerificationToken(ctx context.Context, userID, newEmail, token string, expiresAt time.Time) error { return m.Called(ctx, userID, newEmail, token, expiresAt).Error(0) }
func (m *MockUserRepository) GetEmailVerificationToken(ctx context.Context, token string) (string, string, string, error) { 
	ret := m.Called(ctx, token)
	return ret.String(0), ret.String(1), ret.String(2), ret.Error(3)
}
func (m *MockUserRepository) DeleteEmailVerificationToken(ctx context.Context, token string) error { return m.Called(ctx, token).Error(0) }
func (m *MockUserRepository) GetPendingEmailVerification(ctx context.Context, userID string) (string, error) { 
	ret := m.Called(ctx, userID)
	return ret.String(0), ret.Error(1)
}
func (m *MockUserRepository) LogProfileAudit(ctx context.Context, userID, field, oldValue, newValue, changedBy string) error { return m.Called(ctx, userID, field, oldValue, newValue, changedBy).Error(0) }
func (m *MockUserRepository) SyncProfileSubmissions(ctx context.Context, userID string, formData map[string]string, tenantID string) error { return m.Called(ctx, userID, formData, tenantID).Error(0) }

// MockMailer implements mailer.Mailer
type MockMailer struct {
	mock.Mock
}

func (m *MockMailer) SendNotificationEmail(to, subject, body string) error {
	return m.Called(to, subject, body).Error(0)
}

func (m *MockMailer) SendStateChangeNotification(to, studentName, nodeID, oldState, newState, frontendURL string) error {
	return nil
}

// MockAnalyticsRepository implements repository.AnalyticsRepository
type MockAnalyticsRepository struct {
	mock.Mock
}

func (m *MockAnalyticsRepository) GetStudentsByStage(ctx context.Context) ([]models.StudentStageStats, error) {
	return nil, nil // Not needed for current tests
}
func (m *MockAnalyticsRepository) GetAdvisorLoad(ctx context.Context) ([]models.AdvisorLoadStats, error) {
	return nil, nil
}
func (m *MockAnalyticsRepository) GetOverdueTasks(ctx context.Context) ([]models.OverdueTaskStats, error) {
	return nil, nil
}
func (m *MockAnalyticsRepository) GetTotalStudents(ctx context.Context, filter models.FilterParams) (int, error) {
	args := m.Called(ctx, filter)
	return args.Int(0), args.Error(1)
}
func (m *MockAnalyticsRepository) GetNodeCompletionCount(ctx context.Context, nodeID string, filter models.FilterParams) (int, error) {
	args := m.Called(ctx, nodeID, filter)
	return args.Int(0), args.Error(1)
}
func (m *MockAnalyticsRepository) GetDurationForNodes(ctx context.Context, nodeIDs []string, filter models.FilterParams) ([]float64, error) {
	args := m.Called(ctx, nodeIDs, filter)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]float64), args.Error(1)
}
func (m *MockAnalyticsRepository) GetBottleneck(ctx context.Context, filter models.FilterParams) (string, int, error) {
	args := m.Called(ctx, filter)
	return args.String(0), args.Int(1), args.Error(2)
}
func (m *MockAnalyticsRepository) GetProfileFlagCount(ctx context.Context, key string, minVal float64, filter models.FilterParams) (int, error) {
	args := m.Called(ctx, key, minVal, filter)
	return args.Int(0), args.Error(1)
}

// Risk Analytics
func (m *MockAnalyticsRepository) SaveRiskSnapshot(ctx context.Context, s *models.RiskSnapshot) error {
	return m.Called(ctx, s).Error(0)
}
func (m *MockAnalyticsRepository) GetStudentRiskHistory(ctx context.Context, studentID string) ([]models.RiskSnapshot, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.RiskSnapshot), args.Error(1)
}
func (m *MockAnalyticsRepository) GetHighRiskStudents(ctx context.Context, threshold float64) ([]models.RiskSnapshot, error) {
	args := m.Called(ctx, threshold)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.RiskSnapshot), args.Error(1)
}

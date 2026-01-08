package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/auth"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTenantRepo
type MockTenantRepo struct {
	mock.Mock
}
func (m *MockTenantRepo) AddUserToTenant(ctx context.Context, userID, tenantID, role string, isPrimary bool) error {
	return m.Called(ctx, userID, tenantID, role, isPrimary).Error(0)
}
func (m *MockTenantRepo) GetByID(ctx context.Context, id string) (*models.Tenant, error) { return nil, nil }
func (m *MockTenantRepo) GetBySlug(ctx context.Context, slug string) (*models.Tenant, error) { return nil, nil }
func (m *MockTenantRepo) ListForUser(ctx context.Context, userID string) ([]models.TenantMembershipView, error) { return nil, nil }
func (m *MockTenantRepo) GetPrimaryTenant(ctx context.Context, userID string) (*models.Tenant, error) { return nil, nil }
func (m *MockTenantRepo) Create(ctx context.Context, t *models.Tenant) (string, error) { return "", nil }
func (m *MockTenantRepo) Update(ctx context.Context, id string, updates map[string]interface{}) (*models.Tenant, error) { return nil, nil }
func (m *MockTenantRepo) Delete(ctx context.Context, id string) error { return nil }
func (m *MockTenantRepo) UpdateServices(ctx context.Context, id string, services []string) (string, error) { return "", nil }
func (m *MockTenantRepo) UpdateLogo(ctx context.Context, id, url string) error { return nil }
func (m *MockTenantRepo) Exists(ctx context.Context, id string) (bool, error) { return false, nil }
func (m *MockTenantRepo) ListAllWithStats(ctx context.Context) ([]models.TenantStatsView, error) { return nil, nil }
func (m *MockTenantRepo) GetWithStats(ctx context.Context, id string) (*models.TenantStatsView, error) { return nil, nil }
func (m *MockTenantRepo) GetUserMembership(ctx context.Context, userID, tenantID string) (*models.TenantMembershipView, error) { return nil, nil }
func (m *MockTenantRepo) GetRole(ctx context.Context, userID, tenantID string) (string, error) { return "", nil }
func (m *MockTenantRepo) RemoveUser(ctx context.Context, userID, tenantID string) error { return nil }

// Remove MockEmailSender (Already defined in auth_service_test.go)

// MockStorageClient

// MockStorageClient
type MockStorageClient struct { mock.Mock }
func (m *MockStorageClient) PresignPut(ctx context.Context, key, contentType string, expiry time.Duration) (string, error) {
	args := m.Called(ctx, key, contentType, expiry)
	return args.String(0), args.Error(1)
}
func (m *MockStorageClient) Upload(ctx context.Context, key string, data []byte, contentType string) error { return nil }
func (m *MockStorageClient) Delete(ctx context.Context, key string) error { return nil }

// MockUserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, filter repository.UserFilter, pagination repository.Pagination) ([]models.User, int, error) {
	args := m.Called(ctx, filter, pagination)
	return args.Get(0).([]models.User), args.Int(1), args.Error(2)
}

func (m *MockUserRepository) UpdatePassword(ctx context.Context, id string, hash string) error {
	args := m.Called(ctx, id, hash)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateAvatar(ctx context.Context, id string, avatarURL string) error {
	args := m.Called(ctx, id, avatarURL)
	return args.Error(0)
}

func (m *MockUserRepository) SetActive(ctx context.Context, id string, active bool) error {
	args := m.Called(ctx, id, active)
	return args.Error(0)
}

func (m *MockUserRepository) Exists(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) EmailExists(ctx context.Context, email string, excludeUserID string) (bool, error) {
	args := m.Called(ctx, email, excludeUserID)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) LinkAdvisor(ctx context.Context, studentID, advisorID, tenantID string) error {
	args := m.Called(ctx, studentID, advisorID, tenantID)
	return args.Error(0)
}

func (m *MockUserRepository) ReplaceAdvisors(ctx context.Context, studentID string, advisorIDs []string, tenantID string) error {
	args := m.Called(ctx, studentID, advisorIDs, tenantID)
	return args.Error(0)
}

func (m *MockUserRepository) CreatePasswordResetToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	args := m.Called(ctx, userID, tokenHash, expiresAt)
	return args.Error(0)
}

func (m *MockUserRepository) GetPasswordResetToken(ctx context.Context, tokenHash string) (string, time.Time, error) {
	args := m.Called(ctx, tokenHash)
	return args.String(0), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockUserRepository) DeletePasswordResetToken(ctx context.Context, tokenHash string) error {
	args := m.Called(ctx, tokenHash)
	return args.Error(0)
}

func (m *MockUserRepository) GetTenantRoles(ctx context.Context, userID, tenantID string) ([]string, error) {
	args := m.Called(ctx, userID, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockUserRepository) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockUserRepository) CheckRateLimit(ctx context.Context, userID, action string, window time.Duration) (int, error) {
	args := m.Called(ctx, userID, action, window)
	return args.Int(0), args.Error(1)
}

func (m *MockUserRepository) RecordRateLimit(ctx context.Context, userID, action string) error {
	args := m.Called(ctx, userID, action)
	return args.Error(0)
}

func (m *MockUserRepository) CreateEmailVerificationToken(ctx context.Context, userID, newEmail, token string, expiresAt time.Time) error {
	args := m.Called(ctx, userID, newEmail, token, expiresAt)
	return args.Error(0)
}

func (m *MockUserRepository) GetEmailVerificationToken(ctx context.Context, token string) (string, string, string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.String(1), args.String(2), args.Error(3)
}

func (m *MockUserRepository) DeleteEmailVerificationToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockUserRepository) GetPendingEmailVerification(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *MockUserRepository) LogProfileAudit(ctx context.Context, userID, field, oldValue, newValue, changedBy string) error {
	args := m.Called(ctx, userID, field, oldValue, newValue, changedBy)
	return args.Error(0)
}

func (m *MockUserRepository) SyncProfileSubmissions(ctx context.Context, userID string, formData map[string]string, tenantID string) error {
	args := m.Called(ctx, userID, formData, tenantID)
	return args.Error(0)
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := services.NewUserService(mockRepo, nil, nil, config.AppConfig{}, nil, nil)

	ctx := context.Background()
	req := services.CreateUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Role:      "student",
	}

	// Expect Exists check for username uniqueness
	mockRepo.On("Exists", ctx, mock.AnythingOfType("string")).Return(false, nil).Once()

	// Expect Create call
	mockRepo.On("Create", ctx, mock.MatchedBy(func(u *models.User) bool {
		return u.FirstName == "John" && 
		       u.LastName == "Doe" && 
			   u.Email == "john@example.com" &&
			   len(u.PasswordHash) > 0 // Password should be hashed
	})).Return("new-uuid", nil)

	user, tempPass, err := svc.CreateUser(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "new-uuid", user.ID)
	assert.NotEmpty(t, user.Username)
	assert.NotEmpty(t, tempPass)

	mockRepo.AssertExpectations(t)
}

func TestUserService_ListUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := services.NewUserService(mockRepo, nil, nil, config.AppConfig{}, nil, nil)

	ctx := context.Background()
	filter := repository.UserFilter{Role: "student"}
	pagination := repository.Pagination{Limit: 10, Offset: 0}

	expectedUsers := []models.User{
		{ID: "1", FirstName: "Alice"},
		{ID: "2", FirstName: "Bob"},
	}

	mockRepo.On("List", ctx, filter, pagination).Return(expectedUsers, 2, nil)

	users, total, err := svc.ListUsers(ctx, filter, pagination)
	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, users, 2)
	assert.Equal(t, "Alice", users[0].FirstName)
}

func TestUserService_CreateUser_Full(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockTenant := new(MockTenantRepo)
	svc := services.NewUserService(mockRepo, mockTenant, nil, config.AppConfig{}, nil, nil)

	ctx := context.Background()
	req := services.CreateUserRequest{
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane@example.com",
		Role:      "student",
		TenantID:  "tenant-1",
		AdvisorIDs: []string{"adv-1"},
	}

	mockRepo.On("Exists", ctx, mock.Anything).Return(false, nil).Once()
	mockRepo.On("Create", ctx, mock.Anything).Return("user-1", nil)
	// Tenant
	mockTenant.On("AddUserToTenant", ctx, "user-1", "tenant-1", "student", true).Return(nil)
	// Advisor
	mockRepo.On("LinkAdvisor", ctx, "user-1", "adv-1", "tenant-1").Return(nil)

	user, _, err := svc.CreateUser(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, "user-1", user.ID)
	mockRepo.AssertExpectations(t)
	mockTenant.AssertExpectations(t)
}

func TestUserService_AdminUpdateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := services.NewUserService(mockRepo, nil, nil, config.AppConfig{}, nil, nil)

	ctx := context.Background()
	targetID := "user-2"
	
	// 1. Success
	targetUser := &models.User{ID: targetID, Role: "student", Email: "old@test.com"}
	mockRepo.On("GetByID", ctx, targetID).Return(targetUser, nil).Once()
	mockRepo.On("Update", ctx, mock.MatchedBy(func(u *models.User) bool {
		return u.Email == "new@test.com" && u.FirstName == "NewName"
	})).Return(nil).Once()
	mockRepo.On("ReplaceAdvisors", ctx, targetID, []string{"adv-2"}, "tenant-1").Return(nil).Once()

	req := services.AdminUpdateUserRequest{
		TargetUserID: targetID,
		FirstName:    "NewName",
		Email:        "new@test.com", 
		Role:         "student",
		TenantID:     "tenant-1",
		AdvisorIDs:   []string{"adv-2"},
	}

	err := svc.AdminUpdateUser(ctx, req, "admin")
	assert.NoError(t, err)

	// 2. Fail Superadmin Edit
	super := &models.User{ID: "super", Role: "superadmin"}
	mockRepo.On("GetByID", ctx, "super").Return(super, nil).Once()
	
	err = svc.AdminUpdateUser(ctx, services.AdminUpdateUserRequest{TargetUserID: "super"}, "admin")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot edit superadmin")
}

func TestUserService_ChangePassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := services.NewUserService(mockRepo, nil, nil, config.AppConfig{}, nil, nil)
	ctx := context.Background()

	// Mock GetByID
	hash, _ := auth.HashPassword("oldpass")
	user := &models.User{ID: "u1", PasswordHash: hash}
	mockRepo.On("GetByID", ctx, "u1").Return(user, nil)

	// Mock UpdatePassword
	mockRepo.On("UpdatePassword", ctx, "u1", mock.AnythingOfType("string")).Return(nil)

	// Success
	err := svc.ChangePassword(ctx, "u1", "oldpass", "newpass")
	assert.NoError(t, err)

	// Fail: Wrong old pass
	err = svc.ChangePassword(ctx, "u1", "wrong", "newpass")
	assert.Error(t, err)
}

func TestUserService_UpdateProfile(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockEmail := new(MockEmailSender)
	svc := services.NewUserService(mockRepo, nil, nil, config.AppConfig{}, mockEmail, nil)
	ctx := context.Background()

	// Mock GetByID
	hash, _ := auth.HashPassword("validpass")
	user := &models.User{ID: "u1", PasswordHash: hash, Email: "old@test.com", FirstName: "Old"}
	mockRepo.On("GetByID", ctx, "u1").Return(user, nil)

	// Rate Limit
	mockRepo.On("CheckRateLimit", ctx, "u1", "profile_update", mock.Anything).Return(1, nil)
	mockRepo.On("RecordRateLimit", ctx, "u1", "profile_update").Return(nil)

	// Update (Non-sensitive)
	req := services.UpdateProfileRequest{
		UserID:          "u1",
		CurrentPassword: "validpass",
		Bio:             "New Bio",
	}
	mockRepo.On("Update", ctx, mock.MatchedBy(func(u *models.User) bool {
		return u.Bio == "New Bio"
	})).Return(nil).Once()

	resp, err := svc.UpdateProfile(ctx, req)
	assert.NoError(t, err)
	assert.Contains(t, resp["message"], "profile updated successfully")

	// Update (Email Change)
	reqEmail := services.UpdateProfileRequest{
		UserID:          "u1",
		CurrentPassword: "validpass",
		Email:           "new@test.com",
	}
	// Re-mock Get because separate call or just rely on previous if mocked Once? 
	// Better to use different test or re-setup. Assuming test runs sequentially.
	// We mocked GetByID without .Once() initially, so it persists? No, I mocked it without .Once() above.
	// But `Update` was .Once().
	
	mockRepo.On("Update", ctx, mock.Anything).Return(nil).Once()
	mockRepo.On("EmailExists", ctx, "new@test.com", "u1").Return(false, nil)
	mockRepo.On("CreateEmailVerificationToken", ctx, "u1", "new@test.com", mock.Anything, mock.Anything).Return(nil)
	mockEmail.On("SendEmailVerification", "new@test.com", mock.Anything, mock.Anything).Return(nil)
	mockEmail.On("SendEmailChangeNotification", "old@test.com", mock.Anything).Return(nil)
	mockRepo.On("LogProfileAudit", ctx, "u1", "email", "old@test.com", mock.MatchedBy(func(s string) bool {
		return s == "new@test.com (pending)"
	}), "u1").Return(nil)

	resp, err = svc.UpdateProfile(ctx, reqEmail)
	assert.NoError(t, err)
	assert.Equal(t, "verification_email_sent", resp["message"])
}

func TestUserService_UpdateAvatar(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := services.NewUserService(mockRepo, nil, nil, config.AppConfig{}, nil, nil)
	ctx := context.Background()

	mockRepo.On("UpdateAvatar", ctx, "u1", "http://avatar.jpg").Return(nil)

	err := svc.UpdateAvatar(ctx, "u1", "http://avatar.jpg")
	assert.NoError(t, err)
}

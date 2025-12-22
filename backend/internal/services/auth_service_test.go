package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/auth"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEmailSender
type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) SendEmailVerification(to, token, userName string) error {
	args := m.Called(to, token, userName)
	return args.Error(0)
}

func (m *MockEmailSender) SendEmailChangeNotification(to, userName string) error {
	args := m.Called(to, userName)
	return args.Error(0)
}

func (m *MockEmailSender) SendAddedToRoomNotification(to, userName, roomName string) error {
	args := m.Called(to, userName, roomName)
	return args.Error(0)
}

func (m *MockEmailSender) SendPasswordResetEmail(to, token, userName string) error {
	args := m.Called(to, token, userName)
	return args.Error(0)
}

func TestAuthService_Login(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockEmail := new(MockEmailSender)
	cfg := config.AppConfig{JWTSecret: "testsec", JWTExpDays: 1}
	svc := services.NewAuthService(mockRepo, mockEmail, cfg)

	ctx := context.Background()
	password := "securepass"
	hash, _ := auth.HashPassword(password)

	t.Run("Success - Student Role", func(t *testing.T) {
		user := &models.User{
			ID:           "u1",
			Username:     "student1",
			PasswordHash: hash,
			Role:         models.RoleStudent,
			IsActive:     true,
		}

		mockRepo.On("GetByUsername", ctx, "student1").Return(user, nil).Once()

		resp, err := svc.Login(ctx, "student1", password, "")
		
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "student", resp.Role)
		assert.NotEmpty(t, resp.Token)
	})

	t.Run("Success - Tenant Role", func(t *testing.T) {
		user := &models.User{
			ID:           "u2",
			Username:     "advisor1",
			PasswordHash: hash,
			Role:         models.RoleStudent, // Default role
			IsActive:     true,
		}
		tenantID := "t1"

		mockRepo.On("GetByUsername", ctx, "advisor1").Return(user, nil).Once()
		mockRepo.On("GetTenantRole", ctx, user.ID, tenantID).Return("advisor", nil).Once()

		resp, err := svc.Login(ctx, "advisor1", password, tenantID)
		
		assert.NoError(t, err)
		assert.Equal(t, "advisor", resp.Role)
	})

	t.Run("Fail - Wrong Password", func(t *testing.T) {
		user := &models.User{
			ID:           "u3",
			Username:     "user3",
			PasswordHash: hash,
			IsActive:     true,
		}

		mockRepo.On("GetByUsername", ctx, "user3").Return(user, nil).Once()

		_, err := svc.Login(ctx, "user3", "wrong", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "credentials")
	})

	t.Run("Fail - Inactive", func(t *testing.T) {
		user := &models.User{
			ID:           "u4",
			Username:     "user4",
			PasswordHash: hash,
			IsActive:     false,
		}

		mockRepo.On("GetByUsername", ctx, "user4").Return(user, nil).Once()

		_, err := svc.Login(ctx, "user4", password, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "inactive")
	})
}

func TestAuthService_RequestPasswordReset(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockEmail := new(MockEmailSender)
	cfg := config.AppConfig{}
	svc := services.NewAuthService(mockRepo, mockEmail, cfg)

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		user := &models.User{
			ID:        "u1",
			Email:     "test@example.com",
			Username:  "testuser",
			FirstName: "TestName",
		}

		mockRepo.On("GetByEmail", ctx, "test@example.com").Return(user, nil).Once()
		mockRepo.On("CreatePasswordResetToken", ctx, user.ID, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(nil).Once()
		mockEmail.On("SendPasswordResetEmail", user.Email, mock.AnythingOfType("string"), user.FirstName).Return(nil).Once()

		err := svc.RequestPasswordReset(ctx, "test@example.com")
		assert.NoError(t, err)
	})
}

func TestAuthService_ResetPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockEmail := new(MockEmailSender)
	cfg := config.AppConfig{}
	svc := services.NewAuthService(mockRepo, mockEmail, cfg)

	ctx := context.Background()
	token := "valid-token"
	newPass := "newpass"
	// Token hash logic: sha256 of "valid-token"
	// We need to match precise hash in mock expectation if service hashes it.
	// Service does: hash := sha256...
	
	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetPasswordResetToken", ctx, mock.AnythingOfType("string")).Return("u1", time.Now().Add(1*time.Hour), nil).Once()
		mockRepo.On("UpdatePassword", ctx, "u1", mock.AnythingOfType("string")).Return(nil).Once()
		mockRepo.On("DeletePasswordResetToken", ctx, mock.AnythingOfType("string")).Return(nil).Once()

		err := svc.ResetPassword(ctx, token, newPass)
		assert.NoError(t, err)
	})
}

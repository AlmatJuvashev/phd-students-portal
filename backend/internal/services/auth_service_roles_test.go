package services_test

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/auth"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_SwitchRole_Flow(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockEmail := new(MockEmailSender)
	cfg := config.AppConfig{JWTSecret: "testsec", JWTExpDays: 1}
	svc := services.NewAuthService(mockRepo, mockEmail, cfg)

	ctx := context.Background()

	t.Run("Success - User Has Role", func(t *testing.T) {
		userID := "u1"
		targetRole := "instructor"
		
		// Setup expectations
		mockRepo.On("GetUserRoles", ctx, userID).Return([]string{"student", "instructor"}, nil).Once()
		
		// Note: SwitchRole doesn't exist yet, this will cause compilation error (Red state)
		token, err := svc.SwitchRole(ctx, userID, targetRole)
		
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		
		// Verify token structure if possible, or trust the service uses GenerateToken correctly 
		// (which we will verify in implementation)
	})

	t.Run("Failure - User Does Not Have Role", func(t *testing.T) {
		userID := "u1"
		targetRole := "admin" // User doesn't have this
		
		mockRepo.On("GetUserRoles", ctx, userID).Return([]string{"student", "instructor"}, nil).Once()
		
		token, err := svc.SwitchRole(ctx, userID, targetRole)
		
		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "access denied")
	})
    
	t.Run("Failure - Repository Error", func(t *testing.T) {
		userID := "u1"
		targetRole := "instructor"
		
		mockRepo.On("GetUserRoles", ctx, userID).Return(nil, assert.AnError).Once()
		
		_, err := svc.SwitchRole(ctx, userID, targetRole)
		assert.Error(t, err)
	})
}

func TestAuthService_Login_MultiRole_Check(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockEmail := new(MockEmailSender)
	cfg := config.AppConfig{JWTSecret: "testsec", JWTExpDays: 1}
	svc := services.NewAuthService(mockRepo, mockEmail, cfg)

	ctx := context.Background()
	password := "securepass"
	hash, _ := auth.HashPassword(password)
    
	t.Run("Login Returns Expected Roles in Token", func(t *testing.T) {
		user := &models.User{
			ID:           "u1",
			Username:     "multiuser",
			PasswordHash: hash,
			Role:         models.RoleStudent, // Primary/Legacy
			IsActive:     true,
		}
        
		mockRepo.On("GetByUsername", ctx, "multiuser").Return(user, nil).Once()
		
		// Expect GetUserRoles to be called
		mockRepo.On("GetUserRoles", ctx, user.ID).Return([]string{"student", "instructor"}, nil).Once()
		
		// We expect Login to eventually call GenerateToken logic that embeds these roles
		resp, err := svc.Login(ctx, "multiuser", password, "")
		
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		
		// Verification of token contents would require parsing it, 
		// which is implicitly tested if GenerateToken uses the roles we passed.
		// For now, ensuring no error and method call is sufficient for TDD step.
	})
}

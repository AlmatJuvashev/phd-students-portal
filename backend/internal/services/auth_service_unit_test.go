package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/auth"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_Login_Unit(t *testing.T) {
	mockRepo := NewHandwrittenMockUserRepository()
	cfg := config.AppConfig{JWTSecret: "secret", JWTExpDays: 1}
	svc := services.NewAuthService(mockRepo, nil, cfg)
	ctx := context.Background()

	pw, _ := auth.HashPassword("pass")

	t.Run("Success Normal User", func(t *testing.T) {
		mockRepo.GetByUsernameFunc = func(ctx context.Context, u string) (*models.User, error) {
			return &models.User{ID: "u1", Username: u, PasswordHash: pw, IsActive: true, Role: "student"}, nil
		}
		mockRepo.GetTenantRoleFunc = func(ctx context.Context, uid, tid string) (string, error) {
			return "student", nil
		}
		res, err := svc.Login(ctx, "john", "pass", "t1")
		assert.NoError(t, err)
		assert.Equal(t, "student", res.Role)
	})

	t.Run("Success Superadmin", func(t *testing.T) {
		mockRepo.GetByUsernameFunc = func(ctx context.Context, u string) (*models.User, error) {
			return &models.User{ID: "adm1", Username: u, PasswordHash: pw, IsActive: true, Role: "superadmin"}, nil
		}
		res, err := svc.Login(ctx, "admin", "pass", "t1")
		assert.NoError(t, err)
		assert.Equal(t, "superadmin", res.Role)
		assert.True(t, res.IsSuperadmin)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockRepo.GetByUsernameFunc = func(ctx context.Context, u string) (*models.User, error) {
			return nil, repository.ErrNotFound
		}
		_, err := svc.Login(ctx, "ghost", "pass", "t1")
		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
	})

	t.Run("Tenant Denied", func(t *testing.T) {
		mockRepo.GetByUsernameFunc = func(ctx context.Context, u string) (*models.User, error) {
			return &models.User{ID: "u1", Username: u, PasswordHash: pw, IsActive: true, Role: "student"}, nil
		}
		mockRepo.GetTenantRoleFunc = func(ctx context.Context, uid, tid string) (string, error) {
			return "", errors.New("not member")
		}
		_, err := svc.Login(ctx, "john", "pass", "t1")
		assert.Error(t, err)
		assert.Equal(t, "access denied to this portal", err.Error())
	})

	t.Run("User Inactive", func(t *testing.T) {
		mockRepo.GetByUsernameFunc = func(ctx context.Context, u string) (*models.User, error) {
			return &models.User{ID: "u1", IsActive: false, PasswordHash: pw}, nil
		}
		_, err := svc.Login(ctx, "john", "pass", "t1")
		assert.Error(t, err)
		assert.Equal(t, "account inactive", err.Error())
	})

	t.Run("No Tenant Success", func(t *testing.T) {
		mockRepo.GetByUsernameFunc = func(ctx context.Context, u string) (*models.User, error) {
			return &models.User{ID: "u1", IsActive: true, Role: "student", PasswordHash: pw}, nil
		}
		res, err := svc.Login(ctx, "john", "pass", "")
		assert.NoError(t, err)
		assert.Equal(t, "student", res.Role)
	})
}

func TestAuthService_PasswordReset_Unit(t *testing.T) {
	mockRepo := NewHandwrittenMockUserRepository()
	mockEmail := NewManualEmailSender()
	
	user := &models.User{
		ID:        "u1",
		Email:     "test@test.com",
		FirstName: "John",
	}
	
	mockRepo.GetByEmailFunc = func(ctx context.Context, email string) (*models.User, error) {
		return user, nil
	}
	
	var capturedHash string
	mockRepo.CreatePasswordResetTokenFunc = func(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
		capturedHash = tokenHash
		return nil
	}
	
	svc := services.NewAuthService(mockRepo, mockEmail, config.AppConfig{})
	
	// Request reset
	err := svc.RequestPasswordReset(context.Background(), "test@test.com")
	assert.NoError(t, err)
	assert.NotEmpty(t, capturedHash)
	
	// Reset password
	mockRepo.GetPasswordResetTokenFunc = func(ctx context.Context, tokenHash string) (string, time.Time, error) {
		if tokenHash == capturedHash {
			return "u1", time.Now().Add(time.Hour), nil
		}
		return "", time.Time{}, repository.ErrNotFound
	}
	
	mockRepo.UpdatePasswordFunc = func(ctx context.Context, userID, hash string) error {
		return nil
	}
	mockRepo.DeletePasswordResetTokenFunc = func(ctx context.Context, tokenHash string) error {
		return nil
	}
	
	// Use some token... wait, svc.RequestPasswordReset sends it via email.
	// I'll just assume I have it.
	err = svc.ResetPassword(context.Background(), "dummy-token", "new-password")
	// Actually I need the REAL token to match the hash I generated in the test.
	// But the service generates the token INTERNALLY and only sends it via email.
	// To test this properly without exposing the token, I have to either:
	// 1. Mock the email sender to capture the token.
	// 2. OR modify the service to return the token (but it shouldn't).
	
	// Let's capture the token from email service!
	var capturedToken string
	mockEmail.SendPasswordResetEmailFunc = func(to, token, userName string) error {
		capturedToken = token
		return nil
	}
	
	_ = svc.RequestPasswordReset(context.Background(), "test@test.com")
	
	err = svc.ResetPassword(context.Background(), capturedToken, "new-password")
	assert.NoError(t, err)

	t.Run("ResetPassword Expired", func(t *testing.T) {
		mockRepo.GetPasswordResetTokenFunc = func(ctx context.Context, h string) (string, time.Time, error) {
			return "u1", time.Now().Add(-time.Hour), nil
		}
		err := svc.ResetPassword(context.Background(), "token123", "newpass")
		assert.Error(t, err)
		assert.Equal(t, "token expired", err.Error())
	})

	t.Run("ResetPassword Invalid", func(t *testing.T) {
		mockRepo.GetPasswordResetTokenFunc = func(ctx context.Context, h string) (string, time.Time, error) {
			return "", time.Time{}, errors.New("not found")
		}
		err := svc.ResetPassword(context.Background(), "invalid", "newpass")
		assert.Error(t, err)
		assert.Equal(t, "invalid or expired token", err.Error())
	})

	t.Run("RequestPasswordReset Repo Error", func(t *testing.T) {
		mockRepo.GetByEmailFunc = func(ctx context.Context, e string) (*models.User, error) {
			return &models.User{ID: "u1"}, nil
		}
		mockRepo.CreatePasswordResetTokenFunc = func(ctx context.Context, u, h string, e time.Time) error {
			return assert.AnError
		}
		err := svc.RequestPasswordReset(context.Background(), "test@test.com")
		assert.Error(t, err)
	})

	t.Run("ResetPassword Repo Errors", func(t *testing.T) {
		mockRepo.GetPasswordResetTokenFunc = func(ctx context.Context, h string) (string, time.Time, error) {
			return "u1", time.Now().Add(time.Hour), nil
		}
		mockRepo.UpdatePasswordFunc = func(ctx context.Context, u, h string) error {
			return assert.AnError
		}
		err := svc.ResetPassword(context.Background(), "token", "pass")
		assert.Error(t, err)

		mockRepo.UpdatePasswordFunc = func(ctx context.Context, u, h string) error { return nil }
		mockRepo.DeletePasswordResetTokenFunc = func(ctx context.Context, h string) error {
			return assert.AnError
		}
		err = svc.ResetPassword(context.Background(), "token", "pass")
		assert.Error(t, err)
	})
}

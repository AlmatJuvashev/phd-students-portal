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
)

func TestUserService_CreateUser_Unit(t *testing.T) {
	mockRepo := NewHandwrittenMockUserRepository()
	mockRepo.ExistsFunc = func(ctx context.Context, username string) (bool, error) {
		return false, nil
	}
	mockRepo.CreateFunc = func(ctx context.Context, user *models.User) (string, error) {
		return "u1", nil
	}

	svc := services.NewUserService(mockRepo, nil, nil, config.AppConfig{}, nil, nil)
	
	req := services.CreateUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@doe.com",
		Role:      "student",
	}
	
	user, tempPass, err := svc.CreateUser(context.Background(), req)
	
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, tempPass)
	assert.Equal(t, "u1", user.ID)
	assert.Contains(t, user.Username, "jd")
}

func TestUserService_UpdateProfile_Unit(t *testing.T) {
	mockRepo := NewHandwrittenMockUserRepository()
	mockEmail := NewManualEmailSender()
	svc := services.NewUserService(mockRepo, nil, nil, config.AppConfig{}, mockEmail, nil)
	ctx := context.Background()

	pw, _ := auth.HashPassword("old")
	mockRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.User, error) {
		return &models.User{ID: id, PasswordHash: pw, Email: "old@ex.com"}, nil
	}
	mockRepo.CheckRateLimitFunc = func(ctx context.Context, id, act string, w time.Duration) (int, error) {
		return 0, nil
	}
	mockRepo.UpdateFunc = func(ctx context.Context, u *models.User) error {
		return nil
	}
	mockRepo.CreateEmailVerificationTokenFunc = func(ctx context.Context, userID, newEmail, token string, expiresAt time.Time) error {
		return nil
	}
	mockRepo.EmailExistsFunc = func(ctx context.Context, email, exclude string) (bool, error) {
		return false, nil
	}
	mockRepo.LogProfileAuditFunc = func(ctx context.Context, u, f, o, n, c string) error {
		return nil
	}

	t.Run("Success Basic", func(t *testing.T) {
		req := services.UpdateProfileRequest{UserID: "u1", CurrentPassword: "old", Bio: "NewBio"}
		res, err := svc.UpdateProfile(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, "profile updated successfully", res["message"])
	})

	t.Run("Email Change", func(t *testing.T) {
		req := services.UpdateProfileRequest{UserID: "u1", CurrentPassword: "old", Email: "new@ex.com"}
		res, err := svc.UpdateProfile(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, "verification_email_sent", res["message"])
	})

	t.Run("Rate Limit", func(t *testing.T) {
		mockRepo.CheckRateLimitFunc = func(ctx context.Context, id, act string, w time.Duration) (int, error) {
			return 500, nil
		}
		req := services.UpdateProfileRequest{UserID: "u1", CurrentPassword: "old"}
		_, err := svc.UpdateProfile(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rate limit exceeded")
	})

	t.Run("Incorrect Password", func(t *testing.T) {
		req := services.UpdateProfileRequest{UserID: "u1", CurrentPassword: "wrong"}
		_, err := svc.UpdateProfile(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "incorrect password")
	})

	t.Run("Email Taken", func(t *testing.T) {
		mockRepo.CheckRateLimitFunc = func(ctx context.Context, id, act string, w time.Duration) (int, error) { return 0, nil }
		mockRepo.EmailExistsFunc = func(ctx context.Context, e, exc string) (bool, error) { return true, nil }
		req := services.UpdateProfileRequest{UserID: "u1", CurrentPassword: "old", Email: "taken@test.com"}
		_, err := svc.UpdateProfile(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already in use")
	})
}

func TestUserService_UsernameGen_Unit(t *testing.T) {
	mockRepo := NewHandwrittenMockUserRepository()
	svc := services.NewUserService(mockRepo, nil, nil, config.AppConfig{}, nil, nil)
	
	t.Run("Collision Case", func(t *testing.T) {
		calls := 0
		mockRepo.ExistsFunc = func(ctx context.Context, u string) (bool, error) {
			calls++
			if calls < 3 { return true, nil }
			return false, nil
		}
		mockRepo.CreateFunc = func(ctx context.Context, u *models.User) (string, error) { return "u1", nil }
		
		req := services.CreateUserRequest{FirstName: "John", LastName: "Doe", Email: "j@d.com", Role: "student"}
		u, _, err := svc.CreateUser(context.Background(), req)
		assert.NoError(t, err)
		assert.Contains(t, u.Username, "jd")
		assert.GreaterOrEqual(t, calls, 3)
	})

	t.Run("Foreign Characters", func(t *testing.T) {
		mockRepo.ExistsFunc = func(ctx context.Context, u string) (bool, error) { return false, nil }
		mockRepo.CreateFunc = func(ctx context.Context, u *models.User) (string, error) { return "u1", nil }
		
		req := services.CreateUserRequest{FirstName: "Алма", LastName: "Жу", Email: "a@j.com", Role: "student"}
		u, _, err := svc.CreateUser(context.Background(), req)
		assert.NoError(t, err)
		assert.NotEmpty(t, u.Username)
	})
}

func TestUserService_ResetPassword_Unit(t *testing.T) {
	mockRepo := NewHandwrittenMockUserRepository()
	svc := services.NewUserService(mockRepo, nil, nil, config.AppConfig{}, nil, nil)
	ctx := context.Background()

	t.Run("ResetPasswordForUser Success", func(t *testing.T) {
		mockRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.User, error) {
			return &models.User{ID: id}, nil
		}
		mockRepo.UpdatePasswordFunc = func(ctx context.Context, id, hash string) error { return nil }
		u, pass, err := svc.ResetPasswordForUser(ctx, "u1")
		assert.NoError(t, err)
		assert.NotNil(t, u)
		assert.NotEmpty(t, pass)
	})

	t.Run("ResetPasswordForUser Repo Error", func(t *testing.T) {
		mockRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.User, error) {
			return nil, assert.AnError
		}
		_, _, err := svc.ResetPasswordForUser(ctx, "u1")
		assert.Error(t, err)
	})
}

func TestUserService_BasicMethods(t *testing.T) {
	mockRepo := NewHandwrittenMockUserRepository()
	svc := services.NewUserService(mockRepo, nil, nil, config.AppConfig{}, nil, nil)
	ctx := context.Background()

	t.Run("ChangePassword", func(t *testing.T) {
		pw, _ := auth.HashPassword("current")
		mockRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.User, error) {
			return &models.User{ID: id, PasswordHash: pw}, nil
		}
		mockRepo.UpdatePasswordFunc = func(ctx context.Context, id, hash string) error {
			return nil
		}
		err := svc.ChangePassword(ctx, "u1", "current", "new")
		assert.NoError(t, err)

		err = svc.ChangePassword(ctx, "u1", "wrong", "new")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "incorrect password")
	})

	t.Run("SyncProfileSubmissions", func(t *testing.T) {
		mockRepo.SyncProfileSubmissionsFunc = func(ctx context.Context, u string, d map[string]string, t string) error {
			return nil
		}
		err := svc.SyncProfileSubmissions(ctx, "u1", map[string]string{"k": "v"}, "t1")
		assert.NoError(t, err)
	})

	t.Run("GetPendingEmailVerification", func(t *testing.T) {
		mockRepo.GetPendingEmailVerificationFunc = func(ctx context.Context, u string) (string, error) {
			return "pending@ex.com", nil
		}
		email, err := svc.GetPendingEmailVerification(ctx, "u1")
		assert.NoError(t, err)
		assert.Equal(t, "pending@ex.com", email)
	})

	_, _ = svc.GetByID(ctx, "u1")
	
	t.Run("UpdateUser Success", func(t *testing.T) {
		mockRepo.UpdateFunc = func(ctx context.Context, u *models.User) error { return nil }
		err := svc.UpdateUser(ctx, &models.User{ID: "u1"})
		assert.NoError(t, err)
	})

	t.Run("UpdateUser Error", func(t *testing.T) {
		mockRepo.UpdateFunc = func(ctx context.Context, u *models.User) error { return assert.AnError }
		err := svc.UpdateUser(ctx, &models.User{ID: "u1"})
		assert.Error(t, err)
	})

	t.Run("SetActive Success", func(t *testing.T) {
		mockRepo.SetActiveFunc = func(ctx context.Context, id string, a bool) error { return nil }
		err := svc.SetActive(ctx, "u1", true)
		assert.NoError(t, err)
	})

	t.Run("SetActive Error", func(t *testing.T) {
		mockRepo.SetActiveFunc = func(ctx context.Context, id string, a bool) error { return assert.AnError }
		err := svc.SetActive(ctx, "u1", true)
		assert.Error(t, err)
	})

	t.Run("ChangePassword Failures", func(t *testing.T) {
		mockRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.User, error) {
			h, _ := auth.HashPassword("old")
			return &models.User{ID: id, PasswordHash: h}, nil
		}
		
		// Wrong old pass
		err := svc.ChangePassword(ctx, "u1", "wrong", "new")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "incorrect password")

		// Repo error
		mockRepo.UpdatePasswordFunc = func(ctx context.Context, id, hash string) error { return assert.AnError }
		err = svc.ChangePassword(ctx, "u1", "old", "new")
		assert.Error(t, err)
	})

	t.Run("ForceUpdatePassword Error", func(t *testing.T) {
		mockRepo.UpdatePasswordFunc = func(ctx context.Context, id, hash string) error { return assert.AnError }
		err := svc.ForceUpdatePassword(ctx, "u1", "new")
		assert.Error(t, err)
	})

	t.Run("ResetPassword Error", func(t *testing.T) {
		mockRepo.UpdatePasswordFunc = func(ctx context.Context, id, hash string) error { return assert.AnError }
		mockRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.User, error) { return &models.User{ID: id}, nil }
		_, err := svc.ResetPassword(ctx, "u1")
		assert.Error(t, err)
	})

	t.Run("ForceUpdatePassword Success", func(t *testing.T) {
		mockRepo.UpdatePasswordFunc = func(ctx context.Context, id, hash string) error {
			if id != "u1" { return assert.AnError }
			return nil
		}
		err := svc.ForceUpdatePassword(ctx, "u1", "newpass")
		assert.NoError(t, err)
	})

	t.Run("ResetPassword Success", func(t *testing.T) {
		mockRepo.UpdatePasswordFunc = func(ctx context.Context, id, hash string) error { return nil }
		pass, err := svc.ResetPassword(ctx, "u1")
		assert.NoError(t, err)
		assert.NotEmpty(t, pass)
	})
}

func TestUserService_AdminUpdateUser_Unit(t *testing.T) {
	mockRepo := NewHandwrittenMockUserRepository()
	user := &models.User{ID: "u1", Role: "student", IsActive: true}
	mockRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.User, error) {
		return user, nil
	}
	mockRepo.UpdateFunc = func(ctx context.Context, u *models.User) error { return nil }
	mockRepo.SetActiveFunc = func(ctx context.Context, id string, a bool) error { return nil }
	mockRepo.LogProfileAuditFunc = func(ctx context.Context, u, f, o, n, c string) error { return nil }
	
	svc := services.NewUserService(mockRepo, nil, nil, config.AppConfig{}, nil, nil)
	
	t.Run("Update Role and Info", func(t *testing.T) {
		req := services.AdminUpdateUserRequest{
			TargetUserID: "u1",
			FirstName:    "New",
			LastName:     "Name",
			Role:         "advisor",
		}
		err := svc.AdminUpdateUser(context.Background(), req, "admin")
		assert.NoError(t, err)
		assert.Equal(t, models.Role("advisor"), user.Role)
	})

	t.Run("Repo Error", func(t *testing.T) {
		mockRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.User, error) {
			return nil, assert.AnError
		}
		err := svc.AdminUpdateUser(context.Background(), services.AdminUpdateUserRequest{TargetUserID: "u1"}, "admin")
		assert.Error(t, err)
	})
}

func TestUserService_Avatar_Unit(t *testing.T) {
	mockRepo := NewHandwrittenMockUserRepository()
	mockStorage := &services.MockStorageClient{}
	svc := services.NewUserService(mockRepo, nil, nil, config.AppConfig{S3Endpoint: "http://s3", S3Bucket: "bucket"}, nil, mockStorage)
	ctx := context.Background()

	t.Run("Presign Success", func(t *testing.T) {
		mockStorage.PresignPutFn = func(ctx context.Context, k, c string, e time.Duration) (string, error) {
			return "presigned-url", nil
		}
		url, _, public, err := svc.PresignAvatarUpload(ctx, "u1", "me.png", "image/png", 1024)
		assert.NoError(t, err)
		assert.Equal(t, "presigned-url", url)
		assert.Contains(t, public, "http://s3/bucket/avatars/u1")
	})

	t.Run("UpdateAvatar Success", func(t *testing.T) {
		mockRepo.UpdateAvatarFunc = func(ctx context.Context, id, url string) error {
			return nil
		}
		err := svc.UpdateAvatar(ctx, "u1", "http://public/avatar.png")
		assert.NoError(t, err)
	})
}

func TestUserService_VerifyEmail_Unit(t *testing.T) {
	mockRepo := NewHandwrittenMockUserRepository()
	ctx := context.Background()
	svc := services.NewUserService(mockRepo, nil, nil, config.AppConfig{}, nil, nil)

	t.Run("Success", func(t *testing.T) {
		mockRepo.GetEmailVerificationTokenFunc = func(ctx context.Context, token string) (string, string, string, error) {
			return "u1", "new@test.com", "v", nil
		}
		mockRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.User, error) {
			return &models.User{ID: id, Email: "old@test.com"}, nil
		}
		mockRepo.UpdateFunc = func(ctx context.Context, u *models.User) error { return nil }

		newEmail, err := svc.VerifyEmailChange(ctx, "token123")
		assert.NoError(t, err)
		assert.Equal(t, "new@test.com", newEmail)
	})

	t.Run("Token Fail", func(t *testing.T) {
		mockRepo.GetEmailVerificationTokenFunc = func(ctx context.Context, token string) (string, string, string, error) {
			return "", "", "", assert.AnError
		}
		_, err := svc.VerifyEmailChange(ctx, "token123")
		assert.Error(t, err)
	})

	t.Run("Get User Fail", func(t *testing.T) {
		mockRepo.GetEmailVerificationTokenFunc = func(ctx context.Context, token string) (string, string, string, error) {
			return "u1", "n@e.com", "v", nil
		}
		mockRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.User, error) {
			return nil, assert.AnError
		}
		_, err := svc.VerifyEmailChange(ctx, "token123")
		assert.Error(t, err)
	})
}

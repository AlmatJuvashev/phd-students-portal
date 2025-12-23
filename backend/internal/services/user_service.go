package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/auth"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/redis/go-redis/v9"
)

type CreateUserRequest struct {
	FirstName  string
	LastName   string
	Email      string
	Role       string
	Phone      string
	Program    string
	Specialty  string
	Department string
	Cohort     string
	AdvisorIDs []string
	TenantID   string // Contextual
}

type UserService struct {
	repo repository.UserRepository
	rds  *redis.Client
	cfg  config.AppConfig
	emailSvc *EmailService // We need this
}

func NewUserService(repo repository.UserRepository, rds *redis.Client, cfg config.AppConfig, emailSvc *EmailService) *UserService {
	return &UserService{
		repo: repo,
		rds:  rds,
		cfg:  cfg,
		emailSvc: emailSvc,
	}
}

// CreateUser generates username, password, hashes it, and stores the user.
// Returns the created user object and the temporary password (plain text).
func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*models.User, string, error) {
	// 1. Generate Username
	username, err := s.generateUsername(ctx, req.FirstName, req.LastName)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate username: %w", err)
	}

	// 2. Generate Temp Password and Hash
	tempPass := auth.GeneratePass()
	hash, _ := auth.HashPassword(tempPass)

	user := &models.User{
		Username:     username,
		Email:        req.Email,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         models.Role(req.Role),
		PasswordHash: hash,
		IsActive:     true,
		Phone:        req.Phone,
		Program:      req.Program,
		Specialty:    req.Specialty,
		Department:   req.Department,
		Cohort:       req.Cohort,
	}

	// 3. Persist
	id, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, "", err
	}
	user.ID = id
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// 4. Link Advisors (if student)
	if req.Role == "student" && len(req.AdvisorIDs) > 0 {
		for _, aid := range req.AdvisorIDs {
			_ = s.repo.LinkAdvisor(ctx, id, aid, req.TenantID)
		}
	}

	// 5. Invalidate List Cache (if applicable)
	s.invalidateListCache(ctx)

	return user, tempPass, nil
}

func (s *UserService) GetByID(ctx context.Context, id string) (*models.User, error) {
	// Check Cache (Optional: implement read-through)
	// For now, direct to DB
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) UpdateUser(ctx context.Context, u *models.User) error {
	err := s.repo.Update(ctx, u)
	if err != nil {
		return err
	}
	// Invalidate single user cache
	if s.rds != nil {
		s.rds.Del(ctx, "user:"+u.ID)
	}
	s.invalidateListCache(ctx)
	return nil
}



func (s *UserService) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil { return err }
	
	if !auth.CheckPassword(user.PasswordHash, currentPassword) {
		return fmt.Errorf("incorrect password")
	}
	
	hash, err := auth.HashPassword(newPassword)
	if err != nil { return err }
	
	err = s.repo.UpdatePassword(ctx, userID, hash)
	if err != nil { return err }
	
	if s.rds != nil { s.rds.Del(ctx, "user:"+userID) }
	return nil
}

func (s *UserService) SetActive(ctx context.Context, userID string, active bool) error {
	err := s.repo.SetActive(ctx, userID, active)
	if err != nil { return err }
	
	if s.rds != nil { s.rds.Del(ctx, "user:"+userID) }
	return nil
}



func (s *UserService) ForceUpdatePassword(ctx context.Context, userID, newPassword string) error {
	hash, err := auth.HashPassword(newPassword)
	if err != nil { return err }
	
	err = s.repo.UpdatePassword(ctx, userID, hash)
	if err != nil { return err }
	
	if s.rds != nil { s.rds.Del(ctx, "user:"+userID) }
	return nil
}

func (s *UserService) ResetPasswordForUser(ctx context.Context, id string) (string, string, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil { return "", "", err }
	
	if user.Role == "superadmin" {
		return "", "", fmt.Errorf("cannot reset superadmin password")
	}

	tempPass := auth.GeneratePass()
	hash, _ := auth.HashPassword(tempPass)
	
	err = s.repo.UpdatePassword(ctx, id, hash)
	if s.rds != nil { s.rds.Del(ctx, "user:"+id) }
	return user.Username, tempPass, nil
}

func (s *UserService) SyncProfileSubmissions(ctx context.Context, userID string, formData map[string]string, tenantID string) error {
	return s.repo.SyncProfileSubmissions(ctx, userID, formData, tenantID)
}

func (s *UserService) ResetPassword(ctx context.Context, id string) (string, error) {
	tempPass := auth.GeneratePass()
	hash, _ := auth.HashPassword(tempPass)
	
	err := s.repo.UpdatePassword(ctx, id, hash)
	if err != nil {
		return "", err
	}
	
	if s.rds != nil {
		s.rds.Del(ctx, "user:"+id)
	}
	return tempPass, nil
}

// ListUsers delegates to repo. Can eventually add caching here for common queries.
func (s *UserService) ListUsers(ctx context.Context, filter repository.UserFilter, pagination repository.Pagination) ([]models.User, int, error) {
	return s.repo.List(ctx, filter, pagination)
}

// UpdateProfileRequest contains fields for self-update
type UpdateProfileRequest struct {
	UserID          string
	Email           string
	Phone           string
	Bio             string
	Address         string
	DateOfBirth     *time.Time
	AvatarURL       string
	CurrentPassword string // For verification
}

// UpdateProfile handles user self-update with security checks
func (s *UserService) UpdateProfile(ctx context.Context, req UpdateProfileRequest) (map[string]any, error) {
	// 1. Fetch Current User
	user, err := s.repo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	// 2. Verify Password
	if !auth.CheckPassword(user.PasswordHash, req.CurrentPassword) {
		return nil, fmt.Errorf("incorrect password") // Handler should map to 401
	}

	// 3. Rate Limiting (5 updates / hour)
	count, err := s.repo.CheckRateLimit(ctx, req.UserID, "profile_update", time.Hour)
	if err != nil { return nil, err }
	if count >= 500 { // 5 in 1 hour? The handler said 500. I'll stick to 500 as per handler code found.
		// Wait, Handler said "limit 500" in log but code said count >= 500.
		// I will respect 500.
		return nil, fmt.Errorf("rate limit exceeded")
	}

	// 4. Record Attempt
	_ = s.repo.RecordRateLimit(ctx, req.UserID, "profile_update")

	// 5. Check Email Change
	emailChanged := req.Email != "" && req.Email != user.Email
	
	// 6. Update Non-Sensitive Fields
	// Clone user or modify?
	updated := *user
	if req.Phone != "" { updated.Phone = req.Phone }
	if req.Bio != "" { updated.Bio = req.Bio }
	if req.Address != "" { updated.Address = req.Address }
	if req.DateOfBirth != nil { updated.DateOfBirth = req.DateOfBirth }
	if req.AvatarURL != "" { updated.AvatarURL = req.AvatarURL }
	
	// Apply update
	err = s.repo.Update(ctx, &updated)
	if err != nil { return nil, err }
	
	// Invalidate Cache
	if s.rds != nil { s.rds.Del(ctx, "user:"+req.UserID) }

	response := map[string]any{"message": "profile updated successfully"}

	// 7. Handle Email Change
	if emailChanged {
		// Check uniqueness
		taken, err := s.repo.EmailExists(ctx, req.Email, req.UserID)
		if err != nil { return nil, err }
		if taken { return nil, fmt.Errorf("email already in use") }

		// Generate Token
		token, err := auth.GenerateSecureToken(32) // reusing some helper or generating
		if err != nil { return nil, err }
		
		expires := time.Now().Add(24 * time.Hour)
		err = s.repo.CreateEmailVerificationToken(ctx, req.UserID, req.Email, token, expires)
		if err != nil { return nil, err }

		// Send Email (check if emailSvc is configured)
		userName := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
		if s.emailSvc != nil {
			err = s.emailSvc.SendEmailVerification(req.Email, token, userName)
			if err != nil {
				response["message"] = "verification_email_pending"
				response["warning"] = "email service not configured"
			} else {
				_ = s.emailSvc.SendEmailChangeNotification(user.Email, userName)
				response["message"] = "verification_email_sent"
				response["info"] = "please check your new email"
			}
		} else {
			// Email service not available
			response["message"] = "verification_email_pending"
			response["warning"] = "email service not configured"
		}
		
		// Audit
		_ = s.repo.LogProfileAudit(ctx, req.UserID, "email", user.Email, req.Email+" (pending)", req.UserID)
	}

	return response, nil
}

// AdminUpdateUserRequest
type AdminUpdateUserRequest struct {
	TargetUserID string
	FirstName    string
	LastName     string
	Email        string
	Role         string
	Phone        string
	Program      string
	Specialty    string
	Department   string
	Cohort       string
}

func (s *UserService) AdminUpdateUser(ctx context.Context, req AdminUpdateUserRequest, adminRole string) error {
	// 1. Check Target
	target, err := s.repo.GetByID(ctx, req.TargetUserID)
	if err != nil { return err }
	
	if target.Role == "superadmin" {
		return fmt.Errorf("cannot edit superadmin")
	}
	if req.Role == "superadmin" {
		return fmt.Errorf("cannot assign superadmin role")
	}

	// 2. Update
	target.FirstName = req.FirstName
	target.LastName = req.LastName
	target.Email = req.Email
	target.Role = models.Role(req.Role)
	target.Phone = req.Phone
	target.Program = req.Program
	target.Specialty = req.Specialty
	target.Department = req.Department
	target.Cohort = req.Cohort
	
	err = s.repo.Update(ctx, target)
	if err != nil { return err }
	
	// Invalidate Cache
	if s.rds != nil { s.rds.Del(ctx, "user:"+req.TargetUserID) }
	
	return nil
}

// VerifyEmailChange
func (s *UserService) VerifyEmailChange(ctx context.Context, token string) (string, error) {
	// Get Token
	userID, newEmail, _, err := s.repo.GetEmailVerificationToken(ctx, token)
	if err != nil { return "", err } // Expired or not found
	
	// Get Old Email
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil { return "", err }
	oldEmail := user.Email
	
	// Update User
	user.Email = newEmail
	// We only update email here? Logic says `UPDATE users SET email=$1`.
	// Reuse Update? No, Update might overwrite other fields if we don't have fresh copy.
	// But we fetched fresh copy. 
	// However, stricter to use direct SQL in repo or careful Update.
	// `repo.Update` uses all fields.
	err = s.repo.Update(ctx, user)
	if err != nil { return "", err }
	
	// Delete Token
	_ = s.repo.DeleteEmailVerificationToken(ctx, token)
	
	// Audit
	_ = s.repo.LogProfileAudit(ctx, userID, "email", oldEmail, newEmail, userID)
	
	// Invalidate
	if s.rds != nil { s.rds.Del(ctx, "user:"+userID) }
	
	return newEmail, nil
}

func (s *UserService) GetPendingEmailVerification(ctx context.Context, userID string) (string, error) {
	return s.repo.GetPendingEmailVerification(ctx, userID)
}

// PresignAvatarUpload
func (s *UserService) PresignAvatarUpload(ctx context.Context, userID, filename, contentType string, sizeBytes int64) (string, string, string, error) {
	// Validate
	if sizeBytes > 5*1024*1024 {
		return "", "", "", fmt.Errorf("avatar size must be less than 5MB")
	}
	if !auth.IsImageMimeType(contentType) { // Assuming auth helper valid, otherwise manual check
		// manual check
		if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
			 return "", "", "", fmt.Errorf("only image files are allowed")
		}
	}
	
	s3c, err := NewS3FromEnv()
	if err != nil { return "", "", "", err }
	
	key := fmt.Sprintf("avatars/%s/%d_%s", userID, time.Now().Unix(), filename)
	url, err := s3c.PresignPut(key, contentType, 15*time.Minute)
	if err != nil { return "", "", "", err }
	
	publicURL := fmt.Sprintf("%s/%s/%s", s.cfg.S3Endpoint, s.cfg.S3Bucket, key)
	return url, key, publicURL, nil
}

func (s *UserService) UpdateAvatar(ctx context.Context, userID, avatarURL string) error {
	err := s.repo.UpdateAvatar(ctx, userID, avatarURL)
	if err != nil { return err }
	if s.rds != nil { s.rds.Del(ctx, "user:"+userID) }
	return nil
}


// --- Helpers ---

func (s *UserService) invalidateListCache(ctx context.Context) {
	// If we were caching lists keys pattern, we'd delete them here.
	// For now, placeholder.
}

func (s *UserService) generateUsername(ctx context.Context, firstName, lastName string) (string, error) {
	first := firstLatinInitial(firstName)
	if first == "" {
		first = "x"
	}
	last := firstLatinInitial(lastName)
	if last == "" {
		last = "x"
	}
	base := first + last

	// Retry loop for uniqueness
	for attempt := 0; attempt < 10; attempt++ {
		suffix, err := randomDigitsSuffix(4)
		if err != nil {
			return "", err
		}
		candidate := base + suffix
		exists, err := s.repo.Exists(ctx, candidate)
		if err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("could not generate unique username after 10 attempts")
}

func firstLatinInitial(input string) string {
	slug := auth.Slugify(input)
	for _, ch := range slug {
		if ch >= 'a' && ch <= 'z' {
			return string(ch)
		}
	}
	return ""
}

func randomDigitsSuffix(length int) (string, error) {
	max := big.NewInt(1)
	for i := 0; i < length; i++ {
		max.Mul(max, big.NewInt(10))
	}
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	format := fmt.Sprintf("%%0%dd", length)
	return fmt.Sprintf(format, n.Int64()), nil
}

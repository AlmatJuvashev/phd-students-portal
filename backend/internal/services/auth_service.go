package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/auth"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type AuthService struct {
	repo  repository.UserRepository
	email EmailSender
	cfg   config.AppConfig
}

func NewAuthService(repo repository.UserRepository, email EmailSender, cfg config.AppConfig) *AuthService {
	return &AuthService{
		repo:  repo,
		email: email,
		cfg:   cfg,
	}
}

type LoginResponse struct {
	Token          string
	Role           string
	IsSuperadmin   bool
	UserID         string
	ActiveRole     string
	AvailableRoles []string
}

func (s *AuthService) Login(ctx context.Context, username, password string, tenantID string) (*LoginResponse, error) {
	log.Printf("[AuthService.Login] Attempting login for username=%s, tenantID=%s", username, tenantID)
	
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		log.Printf("[AuthService.Login] User not found: username=%s, error=%v", username, err)
		return nil, errors.New("invalid credentials")
	}
	log.Printf("[AuthService.Login] Found user: id=%s, role=%s, isActive=%v", user.ID, user.Role, user.IsActive)

	if !user.IsActive {
		log.Printf("[AuthService.Login] User inactive: id=%s", user.ID)
		return nil, errors.New("account inactive")
	}

	if !auth.CheckPassword(user.PasswordHash, password) {
		log.Printf("[AuthService.Login] Password mismatch for user=%s, hash_prefix=%s", username, user.PasswordHash[:30])
		return nil, errors.New("invalid credentials")
	}
	log.Printf("[AuthService.Login] Password verified for user=%s", username)

	// Verify Tenant Access
	var roles []string
	if user.Role == "superadmin" {
		roles = []string{"superadmin"}
		log.Printf("[AuthService.Login] User is superadmin, bypassing tenant check")
	} else if tenantID != "" {
		// Check membership
		log.Printf("[AuthService.Login] Checking tenant membership: userID=%s, tenantID=%s", user.ID, tenantID)
		tenantRoles, err := s.repo.GetTenantRoles(ctx, user.ID, tenantID)
		if err != nil {
			log.Printf("[AuthService.Login] Tenant access denied: userID=%s, tenantID=%s, error=%v", user.ID, tenantID, err)
			return nil, errors.New("access denied to this portal")
		}
		roles = tenantRoles
		log.Printf("[AuthService.Login] Tenant access granted: roles=%v", roles)
	} else {
		// No tenant context (e.g. platform admin login? or just resolving user role)
		// Fetch all roles from user_roles table
		log.Printf("[AuthService.Login] No tenant context, fetching global roles for userID=%s", user.ID)
		r, err := s.repo.GetUserRoles(ctx, user.ID)
		if err != nil && err.Error() != "record not found" { 
			// If error, fall back to user.Role?
			// For TDD phase 1, we expect GetUserRoles to work or return empty.
			// Existing user.Role is backup.
			roles = []string{string(user.Role)}
		} else {
			roles = r
			if len(roles) == 0 {
				roles = []string{string(user.Role)} // Fallback if no roles in table
			}
		}
		log.Printf("[AuthService.Login] User roles=%v", roles)
	}

	// Determine Active Role
	activeRole := ""
	if len(roles) > 0 {
		activeRole = roles[0]
	}

	token, err := s.GenerateToken(user.ID, roles, activeRole, tenantID, user.Role == "superadmin")
	if err != nil {
		log.Printf("[AuthService.Login] Token generation failed: %v", err)
		return nil, err
	}
	log.Printf("[AuthService.Login] Login successful: userID=%s, roles=%v, tenantID=%s", user.ID, roles, tenantID)
	
	// Assuming LoginResponse.Role is deprecated or we send primary role
	primaryRole := ""
	if len(roles) > 0 {
		primaryRole = roles[0]
	}

	return &LoginResponse{
		Token:          token,
		UserID:         user.ID,
		Role:           primaryRole, // Active Role (legacy support)
		IsSuperadmin:   user.Role == "superadmin",
		ActiveRole:     primaryRole,
		AvailableRoles: roles,
	}, nil
}


func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil // Return nil to avoid email enumeration
	}

	tokenBytes := make([]byte, 32)
	_, _ = rand.Read(tokenBytes)
	token := hex.EncodeToString(tokenBytes)
	
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])
	expiresAt := time.Now().Add(1 * time.Hour)

	if err := s.repo.CreatePasswordResetToken(ctx, user.ID, tokenHash, expiresAt); err != nil {
		return err
	}

	return s.email.SendPasswordResetEmail(email, token, user.FirstName)
}

func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	userID, expiresAt, err := s.repo.GetPasswordResetToken(ctx, tokenHash)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	if time.Now().After(expiresAt) {
		return errors.New("token expired")
	}

	newHash, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := s.repo.UpdatePassword(ctx, userID, newHash); err != nil {
		return err
	}

	return s.repo.DeletePasswordResetToken(ctx, tokenHash)
}

func (s *AuthService) SwitchRole(ctx context.Context, userID, targetRole string) (string, error) {
	// 1. Get all available roles for user (global + tenant specific logic if needed)
	// For Phase 1, we assume GetUserRoles returns all valid roles
	roles, err := s.repo.GetUserRoles(ctx, userID)
	if err != nil {
		return "", err
	}

	// 2. Verify target role is allowed
	allowed := false
	for _, r := range roles {
		if r == targetRole {
			allowed = true
			break
		}
	}
	if !allowed {
		return "", errors.New("access denied")
	}

	// 3. Generate new token with active role
	// Note: We might need tenantID here if switching within a tenant context. 
	// For now, passing empty tenantID implies global switch or purely role-based.
	// TODO: Add tenantID support to SwitchRole if needed.
	return s.GenerateToken(userID, roles, targetRole, "", false)
}

func (s *AuthService) GenerateToken(userID string, roles []string, activeRole string, tenantID string, isSuperadmin bool) (string, error) {
	return auth.GenerateJWTWithTenant(userID, roles, activeRole, tenantID, isSuperadmin, []byte(s.cfg.JWTSecret), s.cfg.JWTExpDays)
}

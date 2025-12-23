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
	Token        string
	Role         string
	IsSuperadmin bool
	UserID       string
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
	var role string
	if user.Role == "superadmin" {
		role = "superadmin"
		log.Printf("[AuthService.Login] User is superadmin, bypassing tenant check")
		// Superadmin has access to everything effectively, but token needs a role.
		// If tenantID is provided, we should check if they are "operating as" superadmin?
		// Existing logic: "Verify user has access to this tenant (unless superadmin)"
		// If superadmin, role stays superadmin.
	} else if tenantID != "" {
		// Check membership
		log.Printf("[AuthService.Login] Checking tenant membership: userID=%s, tenantID=%s", user.ID, tenantID)
		tenantRole, err := s.repo.GetTenantRole(ctx, user.ID, tenantID)
		if err != nil {
			log.Printf("[AuthService.Login] Tenant access denied: userID=%s, tenantID=%s, error=%v", user.ID, tenantID, err)
			return nil, errors.New("access denied to this portal")
		}
		role = tenantRole
		log.Printf("[AuthService.Login] Tenant access granted: role=%s", role)
	} else {
		// No tenant context (e.g. platform admin login? or just resolving user role)
		role = string(user.Role)
		log.Printf("[AuthService.Login] No tenant context, using user role=%s", role)
	}

	token, err := s.GenerateToken(user.ID, role, tenantID, user.Role == "superadmin")
	if err != nil {
		log.Printf("[AuthService.Login] Token generation failed: %v", err)
		return nil, err
	}
	log.Printf("[AuthService.Login] Login successful: userID=%s, role=%s, tenantID=%s", user.ID, role, tenantID)
	
	return &LoginResponse{
		Token:        token,
		UserID:       user.ID,
		Role:         role,
		IsSuperadmin: user.Role == "superadmin", 
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

func (s *AuthService) GenerateToken(userID, role, tenantID string, isSuperadmin bool) (string, error) {
	return auth.GenerateJWTWithTenant(userID, role, tenantID, isSuperadmin, []byte(s.cfg.JWTSecret), s.cfg.JWTExpDays)
}

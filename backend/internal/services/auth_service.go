package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
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
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("account inactive")
	}

	if !auth.CheckPassword(user.PasswordHash, password) {
		return nil, errors.New("invalid credentials")
	}

	// Verify Tenant Access
	var role string
	if user.Role == "superadmin" {
		role = "superadmin"
		// Superadmin has access to everything effectively, but token needs a role.
		// If tenantID is provided, we should check if they are "operating as" superadmin?
		// Existing logic: "Verify user has access to this tenant (unless superadmin)"
		// If superadmin, role stays superadmin.
	} else if tenantID != "" {
		// Check membership
		tenantRole, err := s.repo.GetTenantRole(ctx, user.ID, tenantID)
		if err != nil {
			return nil, errors.New("access denied to this portal")
		}
		role = tenantRole
	} else {
		// No tenant context (e.g. platform admin login? or just resolving user role)
		role = string(user.Role)
	}

	token, err := s.GenerateToken(user.ID, role, tenantID, user.Role == "superadmin")
	if err != nil {
		return nil, err
	}
	
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

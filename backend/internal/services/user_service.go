package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/auth"
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
}

func NewUserService(repo repository.UserRepository, rds *redis.Client) *UserService {
	return &UserService{
		repo: repo,
		rds:  rds,
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
			_ = s.repo.LinkAdvisor(ctx, id, aid)
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

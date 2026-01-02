package services

import (
	"context"
	"log"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/google/uuid"
)

type AuthzService struct {
	repo repository.RBACRepository
}

func NewAuthzService(repo repository.RBACRepository) *AuthzService {
	return &AuthzService{repo: repo}
}

// HasPermission checks if user has specific permission in a context, considering inheritance.
// Inheritance: Global > Tenant > Context (e.g. Course)
func (s *AuthzService) HasPermission(ctx context.Context, userID uuid.UUID, permSlug string, contextType string, contextID uuid.UUID) (bool, error) {
	// 1. Check Global Context (e.g. Superadmin)
	if allowed, err := s.checkContext(ctx, userID, permSlug, models.ContextGlobal, uuid.Nil); err != nil {
		return false, err
	} else if allowed {
		return true, nil
	}

	// 2. Check Tenant Context (if applicable fallback)
	// For simplicity, assuming caller passes specific context. 
	// Real implementation might look up TenantID from CourseID. 
	// For this MVP, we just check the requested context and Global.

	// 3. Check Requested Context
	return s.checkContext(ctx, userID, permSlug, contextType, contextID)
}

func (s *AuthzService) checkContext(ctx context.Context, userID uuid.UUID, permSlug string, cType string, cID uuid.UUID) (bool, error) {
	roles, err := s.repo.GetUserRolesInContext(ctx, userID, cType, cID)
	if err != nil {
		return false, err
	}

	for _, role := range roles {
		perms, err := s.repo.GetRolePermissions(ctx, role.ID)
		if err != nil {
			log.Printf("Error fetching permissions for role %s: %v", role.Name, err)
			continue
		}
		for _, p := range perms {
			if p == permSlug || p == "*" { // Support wildcard
				return true, nil
			}
		}
	}
	return false, nil
}

package services

import (
	"context"
	"errors"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type ResourceService struct {
	repo repository.ResourceRepository
}

func NewResourceService(repo repository.ResourceRepository) *ResourceService {
	return &ResourceService{repo: repo}
}

// Buildings

func (s *ResourceService) CreateBuilding(ctx context.Context, b *models.Building) error {
	if b.TenantID == "" {
		return errors.New("tenant_id is required")
	}
	if b.Name == "" {
		return errors.New("name is required")
	}
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()
	return s.repo.CreateBuilding(ctx, b)
}

func (s *ResourceService) GetBuilding(ctx context.Context, id string) (*models.Building, error) {
	return s.repo.GetBuilding(ctx, id)
}

func (s *ResourceService) ListBuildings(ctx context.Context, tenantID string) ([]models.Building, error) {
	return s.repo.ListBuildings(ctx, tenantID)
}

func (s *ResourceService) UpdateBuilding(ctx context.Context, b *models.Building) error {
	b.UpdatedAt = time.Now()
	return s.repo.UpdateBuilding(ctx, b)
}

func (s *ResourceService) DeleteBuilding(ctx context.Context, id string) error {
	return s.repo.DeleteBuilding(ctx, id)
}

// Rooms

func (s *ResourceService) CreateRoom(ctx context.Context, r *models.Room) error {
	if r.BuildingID == "" {
		return errors.New("building_id is required")
	}
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.Capacity < 0 {
		return errors.New("capacity cannot be negative")
	}
	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()
	return s.repo.CreateRoom(ctx, r)
}

func (s *ResourceService) GetRoom(ctx context.Context, id string) (*models.Room, error) {
	return s.repo.GetRoom(ctx, id)
}

func (s *ResourceService) ListRooms(ctx context.Context, buildingID string) ([]models.Room, error) {
	return s.repo.ListRooms(ctx, buildingID)
}

func (s *ResourceService) UpdateRoom(ctx context.Context, r *models.Room) error {
	if r.Capacity < 0 {
		return errors.New("capacity cannot be negative")
	}
	r.UpdatedAt = time.Now()
	return s.repo.UpdateRoom(ctx, r)
}

func (s *ResourceService) DeleteRoom(ctx context.Context, id string) error {
	return s.repo.DeleteRoom(ctx, id)
}

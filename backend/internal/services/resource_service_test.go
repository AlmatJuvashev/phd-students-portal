package services

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockResourceRepo for service testing
type MockResourceRepo struct {
	mock.Mock
}

func (m *MockResourceRepo) CreateBuilding(ctx context.Context, b *models.Building) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}
func (m *MockResourceRepo) GetBuilding(ctx context.Context, id string) (*models.Building, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Building), args.Error(1)
}
func (m *MockResourceRepo) ListBuildings(ctx context.Context, tenantID string) ([]models.Building, error) {
	args := m.Called(ctx, tenantID)
	return args.Get(0).([]models.Building), args.Error(1)
}
func (m *MockResourceRepo) UpdateBuilding(ctx context.Context, b *models.Building) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}
func (m *MockResourceRepo) DeleteBuilding(ctx context.Context, id string, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}
// Room mocks
func (m *MockResourceRepo) CreateRoom(ctx context.Context, r *models.Room) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}
func (m *MockResourceRepo) GetRoom(ctx context.Context, id string) (*models.Room, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Room), args.Error(1)
}
func (m *MockResourceRepo) ListRooms(ctx context.Context, tenantID string, buildingID string) ([]models.Room, error) {
	args := m.Called(ctx, tenantID, buildingID)
	return args.Get(0).([]models.Room), args.Error(1)
}
func (m *MockResourceRepo) UpdateRoom(ctx context.Context, r *models.Room) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}
func (m *MockResourceRepo) DeleteRoom(ctx context.Context, id string, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}
func (m *MockResourceRepo) SetAvailability(ctx context.Context, avail *models.InstructorAvailability) error {
	args := m.Called(ctx, avail)
	return args.Error(0)
}
func (m *MockResourceRepo) GetAvailability(ctx context.Context, instructorID string) ([]models.InstructorAvailability, error) {
	args := m.Called(ctx, instructorID)
	return args.Get(0).([]models.InstructorAvailability), args.Error(1)
}
func (m *MockResourceRepo) SetRoomAttribute(ctx context.Context, attr *models.RoomAttribute) error {
	args := m.Called(ctx, attr)
	return args.Error(0)
}
func (m *MockResourceRepo) GetRoomAttributes(ctx context.Context, roomID string) ([]models.RoomAttribute, error) {
	args := m.Called(ctx, roomID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.RoomAttribute), args.Error(1)
}

// Ensure mock implements interface
var _ repository.ResourceRepository = (*MockResourceRepo)(nil)

func TestResourceService_Buildings(t *testing.T) {
	mockRepo := new(MockResourceRepo)
	svc := NewResourceService(mockRepo)
	ctx := context.Background()

	// 1. Success
	b := &models.Building{
		TenantID: "t1",
		Name:     "Main Hall",
		Address:  "123 St",
	}
	mockRepo.On("CreateBuilding", ctx, b).Return(nil)
	err := svc.CreateBuilding(ctx, b)
	assert.NoError(t, err)
	assert.NotZero(t, b.CreatedAt)

	// 2. Validation
	bInvalid := &models.Building{TenantID: "t1"} // Missing Name
	err = svc.CreateBuilding(ctx, bInvalid)
	assert.Error(t, err)
	assert.Equal(t, "name is required", err.Error())
}

func TestResourceService_Rooms(t *testing.T) {
	mockRepo := new(MockResourceRepo)
	svc := NewResourceService(mockRepo)
	ctx := context.Background()

	// 1. Success
	r := &models.Room{
		BuildingID: "b1",
		Name:       "101",
		Capacity:   30,
		Type:       "lecture",
	}
	mockRepo.On("CreateRoom", ctx, r).Return(nil)
	err := svc.CreateRoom(ctx, r)
	assert.NoError(t, err)

	// 2. Validation
	rInvalid := &models.Room{BuildingID: "b1", Name: "102", Capacity: -1} // Negative Capacity
	err = svc.CreateRoom(ctx, rInvalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "capacity cannot be negative")
}

package services

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestResourceService_Buildings(t *testing.T) {
	repo := new(MockResourceRepository)
	svc := NewResourceService(repo)
	ctx := context.Background()

	t.Run("CreateBuilding_Success", func(t *testing.T) {
		b := &models.Building{TenantID: "t1", Name: "Engineering Hall"}
		repo.On("CreateBuilding", ctx, mock.MatchedBy(func(arg *models.Building) bool {
			return arg.TenantID == "t1" && arg.Name == "Engineering Hall" && !arg.CreatedAt.IsZero()
		})).Return(nil)

		err := svc.CreateBuilding(ctx, b)
		assert.NoError(t, err)
	})

	t.Run("CreateBuilding_MissingTenant", func(t *testing.T) {
		b := &models.Building{Name: "No Tenant"}
		err := svc.CreateBuilding(ctx, b)
		assert.Error(t, err)
		assert.Equal(t, "tenant_id is required", err.Error())
	})

	t.Run("ListBuildings", func(t *testing.T) {
		repo.On("ListBuildings", ctx, "t1").Return([]models.Building{{ID: "b1", Name: "Gym"}}, nil)
		list, err := svc.ListBuildings(ctx, "t1")
		assert.NoError(t, err)
		assert.Len(t, list, 1)
		assert.Equal(t, "Gym", list[0].Name)
	})

	t.Run("UpdateBuilding", func(t *testing.T) {
		b := &models.Building{ID: "b1", Name: "Gym v2"}
		repo.On("UpdateBuilding", ctx, mock.MatchedBy(func(arg *models.Building) bool {
			return arg.ID == "b1" && arg.Name == "Gym v2" && !arg.UpdatedAt.IsZero()
		})).Return(nil)
		
		err := svc.UpdateBuilding(ctx, b)
		assert.NoError(t, err)
	})

	t.Run("DeleteBuilding", func(t *testing.T) {
		repo.On("DeleteBuilding", ctx, "b1", "u1").Return(nil)
		err := svc.DeleteBuilding(ctx, "b1", "u1")
		assert.NoError(t, err)
	})
	
	t.Run("GetBuilding", func(t *testing.T) {
		repo.On("GetBuilding", ctx, "b1").Return(&models.Building{ID: "b1"}, nil)
		res, err := svc.GetBuilding(ctx, "b1")
		assert.NoError(t, err)
		assert.Equal(t, "b1", res.ID)
	})
}

func TestResourceService_Rooms(t *testing.T) {
	repo := new(MockResourceRepository)
	svc := NewResourceService(repo)
	ctx := context.Background()

	t.Run("CreateRoom_Success", func(t *testing.T) {
		r := &models.Room{BuildingID: "b1", Name: "101", Capacity: 30}
		repo.On("CreateRoom", ctx, mock.MatchedBy(func(arg *models.Room) bool {
			return arg.BuildingID == "b1" && arg.Name == "101" && !arg.CreatedAt.IsZero()
		})).Return(nil)

		err := svc.CreateRoom(ctx, r)
		assert.NoError(t, err)
	})

	t.Run("CreateRoom_NegativeCapacity", func(t *testing.T) {
		r := &models.Room{BuildingID: "b1", Name: "101", Capacity: -5}
		err := svc.CreateRoom(ctx, r)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "capacity cannot be negative")
	})

	t.Run("ListRooms", func(t *testing.T) {
		repo.On("ListRooms", ctx, "t1", "b1").Return([]models.Room{{ID: "r1"}}, nil)
		list, err := svc.ListRooms(ctx, "t1", "b1")
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})

	t.Run("ListRooms_MissingTenant", func(t *testing.T) {
		_, err := svc.ListRooms(ctx, "", "b1")
		assert.Error(t, err)
		assert.Equal(t, "tenant_id is required", err.Error())
	})

	t.Run("UpdateRoom", func(t *testing.T) {
		r := &models.Room{ID: "r1", Capacity: 50}
		repo.On("UpdateRoom", ctx, mock.MatchedBy(func(arg *models.Room) bool {
			return arg.ID == "r1" && arg.Capacity == 50 && !arg.UpdatedAt.IsZero()
		})).Return(nil)

		err := svc.UpdateRoom(ctx, r)
		assert.NoError(t, err)
	})

	t.Run("DeleteRoom", func(t *testing.T) {
		repo.On("DeleteRoom", ctx, "r1", "u1").Return(nil)
		err := svc.DeleteRoom(ctx, "r1", "u1")
		assert.NoError(t, err)
	})

	t.Run("GetRoom", func(t *testing.T) {
		repo.On("GetRoom", ctx, "r1").Return(&models.Room{ID: "r1"}, nil)
		res, err := svc.GetRoom(ctx, "r1")
		assert.NoError(t, err)
		assert.Equal(t, "r1", res.ID)
	})
}

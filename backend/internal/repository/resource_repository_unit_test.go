package repository

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestResourceRepository_Buildings(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLResourceRepository(sqlxDB)
	ctx := context.Background()

	b := &models.Building{
		TenantID:    "t1",
		Name:        "Main Hall",
		Address:     "123 St",
		Description: "{}",
		IsActive:    true,
	}

	// Create
	mock.ExpectQuery(`INSERT INTO buildings`).
		WithArgs(b.TenantID, b.Name, b.Address, b.Description, b.IsActive).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow("b1", time.Now(), time.Now()))

	err = repo.CreateBuilding(ctx, b)
	assert.NoError(t, err)
	assert.Equal(t, "b1", b.ID)

	// List
	mock.ExpectQuery(`SELECT \* FROM buildings WHERE tenant_id=\$1`).
		WithArgs("t1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow("b1", "Main Hall"))
	
	list, err := repo.ListBuildings(ctx, "t1")
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	// Update
	b.Name = "Updated Hall"
	mock.ExpectExec(`UPDATE buildings`).
		WithArgs(b.Name, b.Address, b.Description, b.IsActive, b.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.UpdateBuilding(ctx, b)
	assert.NoError(t, err)

	// Delete
	mock.ExpectExec(`DELETE FROM buildings`).
		WithArgs("b1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.DeleteBuilding(ctx, "b1")
	assert.NoError(t, err)
	
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestResourceRepository_Rooms(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLResourceRepository(sqlxDB)
	ctx := context.Background()

	r := &models.Room{
		BuildingID: "b1",
		Name:       "101",
		Capacity:   50,
		Type:       "lecture",
		Features:   "[]",
		IsActive:   true,
	}

	// Create
	mock.ExpectQuery(`INSERT INTO rooms`).
		WithArgs(r.BuildingID, r.Name, r.Capacity, r.Type, r.Features, r.IsActive).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow("r1", time.Now(), time.Now()))

	err = repo.CreateRoom(ctx, r)
	assert.NoError(t, err)
	assert.Equal(t, "r1", r.ID)

	// List
	mock.ExpectQuery(`SELECT \* FROM rooms WHERE building_id=\$1`).
		WithArgs("b1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow("r1", "101"))
	
	list, err := repo.ListRooms(ctx, "b1")
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	// Update
	r.Capacity = 60
	mock.ExpectExec(`UPDATE rooms`).
		WithArgs(r.Name, r.Capacity, r.Type, r.Features, r.IsActive, r.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.UpdateRoom(ctx, r)
	assert.NoError(t, err)

	// Delete
	mock.ExpectExec(`DELETE FROM rooms`).
		WithArgs("r1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.DeleteRoom(ctx, "r1")
	assert.NoError(t, err)
	
	assert.NoError(t, mock.ExpectationsWereMet())
}

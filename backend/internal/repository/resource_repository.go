package repository

import (
	"context"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type ResourceRepository interface {
	// Buildings
	CreateBuilding(ctx context.Context, b *models.Building) error
	GetBuilding(ctx context.Context, id string) (*models.Building, error)
	ListBuildings(ctx context.Context, tenantID string) ([]models.Building, error)
	UpdateBuilding(ctx context.Context, b *models.Building) error
	DeleteBuilding(ctx context.Context, id string) error

	// Rooms
	CreateRoom(ctx context.Context, r *models.Room) error
	GetRoom(ctx context.Context, id string) (*models.Room, error)
	ListRooms(ctx context.Context, buildingID string) ([]models.Room, error)
	UpdateRoom(ctx context.Context, r *models.Room) error
	DeleteRoom(ctx context.Context, id string) error

	// Availability
	SetAvailability(ctx context.Context, avail *models.InstructorAvailability) error
	GetAvailability(ctx context.Context, instructorID string) ([]models.InstructorAvailability, error)
}

type SQLResourceRepository struct {
	db *sqlx.DB
}

func NewSQLResourceRepository(db *sqlx.DB) *SQLResourceRepository {
	return &SQLResourceRepository{db: db}
}

// --- Buildings ---

func (r *SQLResourceRepository) CreateBuilding(ctx context.Context, b *models.Building) error {
	return r.db.QueryRowxContext(ctx, `
		INSERT INTO buildings (tenant_id, name, address, description, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`,
		b.TenantID, b.Name, b.Address, b.Description, b.IsActive,
	).Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt)
}

func (r *SQLResourceRepository) GetBuilding(ctx context.Context, id string) (*models.Building, error) {
	var b models.Building
	err := sqlx.GetContext(ctx, r.db, &b, `SELECT * FROM buildings WHERE id=$1`, id)
	return &b, err
}

func (r *SQLResourceRepository) ListBuildings(ctx context.Context, tenantID string) ([]models.Building, error) {
	var list []models.Building
	err := sqlx.SelectContext(ctx, r.db, &list, `SELECT * FROM buildings WHERE tenant_id=$1 ORDER BY name ASC`, tenantID)
	return list, err
}

func (r *SQLResourceRepository) UpdateBuilding(ctx context.Context, b *models.Building) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE buildings SET name=$1, address=$2, description=$3, is_active=$4, updated_at=now()
		WHERE id=$5`,
		b.Name, b.Address, b.Description, b.IsActive, b.ID)
	return err
}

func (r *SQLResourceRepository) DeleteBuilding(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM buildings WHERE id=$1`, id)
	return err
}

// --- Rooms ---

func (r *SQLResourceRepository) CreateRoom(ctx context.Context, rm *models.Room) error {
	return r.db.QueryRowxContext(ctx, `
		INSERT INTO rooms (building_id, name, capacity, floor, department_id, type, features, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`,
		rm.BuildingID, rm.Name, rm.Capacity, rm.Floor, rm.DepartmentID, rm.Type, rm.Features, rm.IsActive,
	).Scan(&rm.ID, &rm.CreatedAt, &rm.UpdatedAt)
}

func (r *SQLResourceRepository) GetRoom(ctx context.Context, id string) (*models.Room, error) {
	var rm models.Room
	err := sqlx.GetContext(ctx, r.db, &rm, `SELECT * FROM rooms WHERE id=$1`, id)
	return &rm, err
}

func (r *SQLResourceRepository) ListRooms(ctx context.Context, buildingID string) ([]models.Room, error) {
	var list []models.Room
	err := sqlx.SelectContext(ctx, r.db, &list, `SELECT * FROM rooms WHERE building_id=$1 ORDER BY floor, name ASC`, buildingID)
	return list, err
}

func (r *SQLResourceRepository) UpdateRoom(ctx context.Context, rm *models.Room) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE rooms SET name=$1, capacity=$2, floor=$3, department_id=$4, type=$5, features=$6, is_active=$7, updated_at=now()
		WHERE id=$8`,
		rm.Name, rm.Capacity, rm.Floor, rm.DepartmentID, rm.Type, rm.Features, rm.IsActive, rm.ID)
	return err
}

func (r *SQLResourceRepository) DeleteRoom(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM rooms WHERE id=$1`, id)
	return err
}

// --- Availability ---

func (r *SQLResourceRepository) SetAvailability(ctx context.Context, avail *models.InstructorAvailability) error {
	// Upsert logic or simple insert? For simplicity, we'll just insert/update on conflict if checking ID, 
	// but here we might want to just create new entries.
	// Let's assume frontend manages IDs or we wipe and replace.
	// For now: Simple Create. logic for overlap check should be in service or FE.
	return r.db.QueryRowxContext(ctx, `
		INSERT INTO instructor_availability (instructor_id, day_of_week, start_time, end_time, is_unavailable)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`,
		avail.InstructorID, avail.DayOfWeek, avail.StartTime, avail.EndTime, avail.IsUnavailable,
	).Scan(&avail.ID, &avail.CreatedAt)
}

func (r *SQLResourceRepository) GetAvailability(ctx context.Context, instructorID string) ([]models.InstructorAvailability, error) {
	var list []models.InstructorAvailability
	err := sqlx.SelectContext(ctx, r.db, &list, `
		SELECT * FROM instructor_availability 
		WHERE instructor_id=$1 
		ORDER BY day_of_week, start_time`, instructorID)
	return list, err
}

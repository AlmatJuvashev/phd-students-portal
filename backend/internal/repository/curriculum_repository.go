package repository

import (
	"context"
	"database/sql"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type CurriculumRepository interface {
	// Programs
	CreateProgram(ctx context.Context, p *models.Program) error
	GetProgram(ctx context.Context, id string) (*models.Program, error)
	ListPrograms(ctx context.Context, tenantID string) ([]models.Program, error)
	UpdateProgram(ctx context.Context, p *models.Program) error
	DeleteProgram(ctx context.Context, id string) error
	
	// Courses
	CreateCourse(ctx context.Context, c *models.Course) error
	GetCourse(ctx context.Context, id string) (*models.Course, error)
	ListCourses(ctx context.Context, tenantID string, programID *string) ([]models.Course, error)
	UpdateCourse(ctx context.Context, c *models.Course) error
	DeleteCourse(ctx context.Context, id string) error
	
	// Journey Maps (Playbooks)
	CreateJourneyMap(ctx context.Context, jm *models.JourneyMap) error
	GetJourneyMapByProgram(ctx context.Context, programID string) (*models.JourneyMap, error)
	
	// Node Definitions
	CreateNodeDefinition(ctx context.Context, nd *models.JourneyNodeDefinition) error
	GetNodeDefinitions(ctx context.Context, journeyMapID string) ([]models.JourneyNodeDefinition, error)
	DeleteNodeDefinition(ctx context.Context, id string) error
	
	// Cohorts
	CreateCohort(ctx context.Context, c *models.Cohort) error
	ListCohorts(ctx context.Context, programID string) ([]models.Cohort, error)
}

type SQLCurriculumRepository struct {
	db *sqlx.DB
}

func NewSQLCurriculumRepository(db *sqlx.DB) *SQLCurriculumRepository {
	return &SQLCurriculumRepository{db: db}
}

// --- Programs ---

func (r *SQLCurriculumRepository) CreateProgram(ctx context.Context, p *models.Program) error {
	return r.db.QueryRowxContext(ctx, `
		INSERT INTO programs (tenant_id, code, title, description, credits, duration_months, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`,
		p.TenantID, p.Code, p.Title, p.Description, p.Credits, p.DurationMonths, p.IsActive,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *SQLCurriculumRepository) GetProgram(ctx context.Context, id string) (*models.Program, error) {
	var p models.Program
	err := sqlx.GetContext(ctx, r.db, &p, `SELECT * FROM programs WHERE id=$1`, id)
	return &p, err
}

func (r *SQLCurriculumRepository) ListPrograms(ctx context.Context, tenantID string) ([]models.Program, error) {
	var programs []models.Program
	err := sqlx.SelectContext(ctx, r.db, &programs, `
		SELECT * FROM programs WHERE tenant_id=$1 ORDER BY created_at DESC`, tenantID)
	return programs, err
}

func (r *SQLCurriculumRepository) UpdateProgram(ctx context.Context, p *models.Program) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE programs SET code=$1, title=$2, description=$3, credits=$4, duration_months=$5, is_active=$6, updated_at=now()
		WHERE id=$7`,
		p.Code, p.Title, p.Description, p.Credits, p.DurationMonths, p.IsActive, p.ID)
	return err
}

func (r *SQLCurriculumRepository) DeleteProgram(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM programs WHERE id=$1`, id)
	return err
}

// --- Courses ---

func (r *SQLCurriculumRepository) CreateCourse(ctx context.Context, c *models.Course) error {
	return r.db.QueryRowxContext(ctx, `
		INSERT INTO courses (tenant_id, program_id, code, title, description, credits, workload_hours, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`,
		c.TenantID, c.ProgramID, c.Code, c.Title, c.Description, c.Credits, c.WorkloadHours, c.IsActive,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *SQLCurriculumRepository) GetCourse(ctx context.Context, id string) (*models.Course, error) {
	var c models.Course
	err := sqlx.GetContext(ctx, r.db, &c, `SELECT * FROM courses WHERE id=$1`, id)
	return &c, err
}

func (r *SQLCurriculumRepository) ListCourses(ctx context.Context, tenantID string, programID *string) ([]models.Course, error) {
	var courses []models.Course
	query := `SELECT * FROM courses WHERE tenant_id=$1`
	args := []interface{}{tenantID}
	
	if programID != nil {
		query += ` AND program_id=$2`
		args = append(args, *programID)
	}
	query += ` ORDER BY code ASC`
	
	err := sqlx.SelectContext(ctx, r.db, &courses, query, args...)
	return courses, err
}

func (r *SQLCurriculumRepository) UpdateCourse(ctx context.Context, c *models.Course) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE courses SET program_id=$1, code=$2, title=$3, description=$4, credits=$5, workload_hours=$6, is_active=$7, updated_at=now()
		WHERE id=$8`,
		c.ProgramID, c.Code, c.Title, c.Description, c.Credits, c.WorkloadHours, c.IsActive, c.ID)
	return err
}

func (r *SQLCurriculumRepository) DeleteCourse(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM courses WHERE id=$1`, id)
	return err
}

// --- Journey Maps ---

func (r *SQLCurriculumRepository) CreateJourneyMap(ctx context.Context, jm *models.JourneyMap) error {
	return r.db.QueryRowxContext(ctx, `
		INSERT INTO journey_maps (program_id, title, version, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`,
		jm.ProgramID, jm.Title, jm.Version, jm.IsActive,
	).Scan(&jm.ID, &jm.CreatedAt)
}

func (r *SQLCurriculumRepository) GetJourneyMapByProgram(ctx context.Context, programID string) (*models.JourneyMap, error) {
	var jm models.JourneyMap
	err := sqlx.GetContext(ctx, r.db, &jm, `SELECT * FROM journey_maps WHERE program_id=$1 LIMIT 1`, programID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &jm, err
}

// --- Node Definitions ---

func (r *SQLCurriculumRepository) CreateNodeDefinition(ctx context.Context, nd *models.JourneyNodeDefinition) error {
	return r.db.QueryRowxContext(ctx, `
		INSERT INTO journey_node_definitions (journey_map_id, parent_node_id, slug, type, title, description, module_key, coordinates, config, prerequisites)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at`,
		nd.JourneyMapID, nd.ParentNodeID, nd.Slug, nd.Type, nd.Title, nd.Description, nd.ModuleKey, nd.Coordinates, nd.Config, pq.Array(nd.Prerequisites),
	).Scan(&nd.ID, &nd.CreatedAt)
}

func (r *SQLCurriculumRepository) GetNodeDefinitions(ctx context.Context, journeyMapID string) ([]models.JourneyNodeDefinition, error) {
	var nodes []models.JourneyNodeDefinition
	err := sqlx.SelectContext(ctx, r.db, &nodes, `SELECT * FROM journey_node_definitions WHERE journey_map_id=$1 ORDER BY module_key, slug`, journeyMapID)
	return nodes, err
}

func (r *SQLCurriculumRepository) DeleteNodeDefinition(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM journey_node_definitions WHERE id=$1`, id)
	return err
}

// --- Cohorts ---

func (r *SQLCurriculumRepository) CreateCohort(ctx context.Context, c *models.Cohort) error {
	return r.db.QueryRowxContext(ctx, `
		INSERT INTO cohorts (program_id, name, start_date, end_date, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`,
		c.ProgramID, c.Name, c.StartDate, c.EndDate, c.IsActive,
	).Scan(&c.ID, &c.CreatedAt)
}

func (r *SQLCurriculumRepository) ListCohorts(ctx context.Context, programID string) ([]models.Cohort, error) {
	var cohorts []models.Cohort
	err := sqlx.SelectContext(ctx, r.db, &cohorts, `SELECT * FROM cohorts WHERE program_id=$1 ORDER BY start_date DESC`, programID)
	return cohorts, err
}

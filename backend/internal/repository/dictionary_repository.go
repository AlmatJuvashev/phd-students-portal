package repository

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type DictionaryRepository interface {
	// Programs
	ListPrograms(ctx context.Context, tenantID string, activeOnly bool) ([]models.Program, error)
	CreateProgram(ctx context.Context, tenantID, name, code string) (string, error)
	UpdateProgram(ctx context.Context, tenantID, id string, name, code string, isActive *bool) error
	DeleteProgram(ctx context.Context, tenantID, id string) error

	// Specialties
	ListSpecialties(ctx context.Context, tenantID string, activeOnly bool, programID string) ([]models.Specialty, error)
	CreateSpecialty(ctx context.Context, tenantID, name, code string, programIDs []string) (string, error)
	UpdateSpecialty(ctx context.Context, tenantID, id, name, code string, isActive *bool, programIDs []string) error
	DeleteSpecialty(ctx context.Context, tenantID, id string) error

	// Cohorts
	ListCohorts(ctx context.Context, tenantID string, activeOnly bool) ([]models.Cohort, error)
	CreateCohort(ctx context.Context, tenantID, name, startDate, endDate string) (string, error)
	UpdateCohort(ctx context.Context, tenantID, id, name, startDate, endDate string, isActive *bool) error
	DeleteCohort(ctx context.Context, tenantID, id string) error

	// Departments
	ListDepartments(ctx context.Context, tenantID string, activeOnly bool) ([]models.Department, error)
	CreateDepartment(ctx context.Context, tenantID, name, code string) (string, error)
	UpdateDepartment(ctx context.Context, tenantID, id, name, code string, isActive *bool) error
	DeleteDepartment(ctx context.Context, tenantID, id string) error
}

type SQLDictionaryRepository struct {
	db *sqlx.DB
}

func NewSQLDictionaryRepository(db *sqlx.DB) *SQLDictionaryRepository {
	return &SQLDictionaryRepository{db: db}
}

// Helper for nullable strings
func dictNullable(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Helper to convert int to string
func itoa(i int) string {
	return strconv.Itoa(i)
}

// --- Programs ---

func (r *SQLDictionaryRepository) ListPrograms(ctx context.Context, tenantID string, activeOnly bool) ([]models.Program, error) {
	query := `SELECT id, name, COALESCE(code, '') as code, is_active, 
              to_char(created_at,'YYYY-MM-DD"T"HH24:MI:SSZ') as created_at,
              to_char(updated_at,'YYYY-MM-DD"T"HH24:MI:SSZ') as updated_at 
              FROM programs WHERE tenant_id = $1`
	args := []interface{}{tenantID}
	
	if activeOnly {
		query += ` AND is_active = true`
	}
	query += ` ORDER BY name`

	var programs []models.Program
	if err := r.db.SelectContext(ctx, &programs, query, args...); err != nil {
		return nil, err
	}
	if programs == nil {
		return []models.Program{}, nil
	}
	return programs, nil
}

func (r *SQLDictionaryRepository) CreateProgram(ctx context.Context, tenantID, name, code string) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, `INSERT INTO programs (name, code, tenant_id) VALUES ($1, $2, $3) RETURNING id`, 
		name, dictNullable(code), tenantID).Scan(&id)
	return id, err
}

func (r *SQLDictionaryRepository) UpdateProgram(ctx context.Context, tenantID, id string, name, code string, isActive *bool) error {
	setParts := []string{"updated_at = now()"}
	args := []interface{}{}
	argId := 1

	if name != "" {
		setParts = append(setParts, "name = $"+itoa(argId))
		args = append(args, name)
		argId++
	}
	if code != "" {
		setParts = append(setParts, "code = $"+itoa(argId))
		args = append(args, code)
		argId++
	}
	if isActive != nil {
		setParts = append(setParts, "is_active = $"+itoa(argId))
		args = append(args, *isActive)
		argId++
	}

	args = append(args, id, tenantID)
	query := "UPDATE programs SET " + strings.Join(setParts, ", ") + " WHERE id = $" + itoa(argId) + " AND tenant_id = $" + itoa(argId+1)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *SQLDictionaryRepository) DeleteProgram(ctx context.Context, tenantID, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE programs SET is_active = false, updated_at = now() WHERE id = $1 AND tenant_id = $2`, id, tenantID)
	return err
}

// --- Specialties ---

func (r *SQLDictionaryRepository) ListSpecialties(ctx context.Context, tenantID string, activeOnly bool, programID string) ([]models.Specialty, error) {
	query := `SELECT id, name, COALESCE(code, '') as code, is_active, 
              to_char(created_at,'YYYY-MM-DD"T"HH24:MI:SSZ') as created_at,
              to_char(updated_at,'YYYY-MM-DD"T"HH24:MI:SSZ') as updated_at 
              FROM specialties WHERE tenant_id = $1`
	
	args := []interface{}{tenantID}
	nextArg := 2

	if activeOnly {
		query += ` AND is_active = true`
	}
	if programID != "" {
		query += ` AND EXISTS (SELECT 1 FROM specialty_programs WHERE specialty_id = specialties.id AND program_id = $` + itoa(nextArg) + `)`
		args = append(args, programID)
	}
	query += ` ORDER BY name`

	type dbSpecialty struct {
		ID        string `db:"id"`
		Name      string `db:"name"`
		Code      string `db:"code"`
		IsActive  bool   `db:"is_active"`
		CreatedAt string `db:"created_at"`
		UpdatedAt string `db:"updated_at"`
	}

	var dbSpecialties []dbSpecialty
	if err := r.db.SelectContext(ctx, &dbSpecialties, query, args...); err != nil {
		return nil, err
	}

	specialties := make([]models.Specialty, len(dbSpecialties))
	for i, s := range dbSpecialties {
		var programIDs []string
		err := r.db.SelectContext(ctx, &programIDs, `SELECT program_id FROM specialty_programs WHERE specialty_id = $1`, s.ID)
		if err != nil {
			programIDs = []string{}
		}

		specialties[i] = models.Specialty{
			ID:         s.ID,
			Name:       s.Name,
			Code:       s.Code,
			ProgramIDs: programIDs,
			IsActive:   s.IsActive,
			CreatedAt:  s.CreatedAt,
			UpdatedAt:  s.UpdatedAt,
		}
	}

	return specialties, nil
}

func (r *SQLDictionaryRepository) CreateSpecialty(ctx context.Context, tenantID, name, code string, programIDs []string) (string, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var id string
	err = tx.QueryRowContext(ctx, `INSERT INTO specialties (name, code, tenant_id) VALUES ($1, $2, $3) RETURNING id`, 
		name, dictNullable(code), tenantID).Scan(&id)
	if err != nil {
		return "", err
	}

	for _, pid := range programIDs {
		if pid != "" && pid != "no_program" {
			_, err := tx.ExecContext(ctx, `INSERT INTO specialty_programs (specialty_id, program_id) VALUES ($1, $2)`, id, pid)
			if err != nil {
				// We can log or ignore, legacy logic ignored errors
				continue 
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return id, nil
}

func (r *SQLDictionaryRepository) UpdateSpecialty(ctx context.Context, tenantID, id, name, code string, isActive *bool, programIDs []string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	setParts := []string{"updated_at = now()"}
	args := []interface{}{}
	argId := 1

	if name != "" {
		setParts = append(setParts, "name = $"+itoa(argId))
		args = append(args, name)
		argId++
	}
	if code != "" {
		setParts = append(setParts, "code = $"+itoa(argId))
		args = append(args, code)
		argId++
	}
	if isActive != nil {
		setParts = append(setParts, "is_active = $"+itoa(argId))
		args = append(args, *isActive)
		argId++
	}

	if len(setParts) > 1 {
		args = append(args, id, tenantID)
		query := "UPDATE specialties SET " + strings.Join(setParts, ", ") + " WHERE id = $" + itoa(argId) + " AND tenant_id = $" + itoa(argId+1)
		_, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}
	}

	if programIDs != nil {
		_, _ = tx.ExecContext(ctx, `DELETE FROM specialty_programs WHERE specialty_id = $1`, id)
		for _, pid := range programIDs {
			if pid != "" && pid != "no_program" {
				_, err := tx.ExecContext(ctx, `INSERT INTO specialty_programs (specialty_id, program_id) VALUES ($1, $2)`, id, pid)
				if err != nil {
					continue
				}
			}
		}
	}

	return tx.Commit()
}

func (r *SQLDictionaryRepository) DeleteSpecialty(ctx context.Context, tenantID, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE specialties SET is_active = false, updated_at = now() WHERE id = $1 AND tenant_id = $2`, id, tenantID)
	return err
}

// --- Cohorts ---

func (r *SQLDictionaryRepository) ListCohorts(ctx context.Context, tenantID string, activeOnly bool) ([]models.Cohort, error) {
	query := `SELECT id, name, start_date, end_date, is_active, created_at, updated_at 
              FROM cohorts WHERE tenant_id = $1`
	args := []interface{}{tenantID}

	if activeOnly {
		query += ` AND is_active = true`
	}
	query += ` ORDER BY name DESC`

	type dbCohort struct {
		ID        string     `db:"id"`
		Name      string     `db:"name"`
		StartDate *time.Time `db:"start_date"`
		EndDate   *time.Time `db:"end_date"`
		IsActive  bool       `db:"is_active"`
		CreatedAt time.Time  `db:"created_at"`
		UpdatedAt time.Time  `db:"updated_at"`
	}

	var dbCohorts []dbCohort
	if err := r.db.SelectContext(ctx, &dbCohorts, query, args...); err != nil {
		return nil, err
	}

	cohorts := make([]models.Cohort, len(dbCohorts))
	for i, dc := range dbCohorts {
		c := models.Cohort{
			ID:        dc.ID,
			Name:      dc.Name,
			IsActive:  dc.IsActive,
			CreatedAt: dc.CreatedAt,
			UpdatedAt: dc.UpdatedAt,
		}
		if dc.StartDate != nil {
			c.StartDate = *dc.StartDate
		}
		if dc.EndDate != nil {
			c.EndDate = *dc.EndDate
		}
		cohorts[i] = c
	}
	
	if cohorts == nil {
		return []models.Cohort{}, nil
	}
	return cohorts, nil
}

func (r *SQLDictionaryRepository) CreateCohort(ctx context.Context, tenantID, name, startDate, endDate string) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, `INSERT INTO cohorts (name, start_date, end_date, tenant_id) VALUES ($1, $2, $3, $4) RETURNING id`, 
		name, dictNullable(startDate), dictNullable(endDate), tenantID).Scan(&id)
	return id, err
}

func (r *SQLDictionaryRepository) UpdateCohort(ctx context.Context, tenantID, id, name, startDate, endDate string, isActive *bool) error {
	setParts := []string{"updated_at = now()"}
	args := []interface{}{}
	argId := 1

	if name != "" {
		setParts = append(setParts, "name = $"+itoa(argId))
		args = append(args, name)
		argId++
	}
	if startDate != "" {
		setParts = append(setParts, "start_date = $"+itoa(argId))
		args = append(args, dictNullable(startDate))
		argId++
	}
	if endDate != "" {
		setParts = append(setParts, "end_date = $"+itoa(argId))
		args = append(args, dictNullable(endDate))
		argId++
	}
	if isActive != nil {
		setParts = append(setParts, "is_active = $"+itoa(argId))
		args = append(args, *isActive)
		argId++
	}

	args = append(args, id, tenantID)
	query := "UPDATE cohorts SET " + strings.Join(setParts, ", ") + " WHERE id = $" + itoa(argId) + " AND tenant_id = $" + itoa(argId+1)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *SQLDictionaryRepository) DeleteCohort(ctx context.Context, tenantID, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE cohorts SET is_active = false, updated_at = now() WHERE id = $1 AND tenant_id = $2`, id, tenantID)
	return err
}

// --- Departments ---

func (r *SQLDictionaryRepository) ListDepartments(ctx context.Context, tenantID string, activeOnly bool) ([]models.Department, error) {
	query := `SELECT id, name, COALESCE(code, '') as code, is_active, 
              to_char(created_at,'YYYY-MM-DD"T"HH24:MI:SSZ') as created_at,
              to_char(updated_at,'YYYY-MM-DD"T"HH24:MI:SSZ') as updated_at 
              FROM departments WHERE tenant_id = $1`
	args := []interface{}{tenantID}

	if activeOnly {
		query += ` AND is_active = true`
	}
	query += ` ORDER BY name ASC`

	var departments []models.Department
	if err := r.db.SelectContext(ctx, &departments, query, args...); err != nil {
		return nil, err
	}
	if departments == nil {
		return []models.Department{}, nil
	}
	return departments, nil
}

func (r *SQLDictionaryRepository) CreateDepartment(ctx context.Context, tenantID, name, code string) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, `INSERT INTO departments (name, code, tenant_id) VALUES ($1, $2, $3) RETURNING id`, 
		name, dictNullable(code), tenantID).Scan(&id)
	return id, err
}

func (r *SQLDictionaryRepository) UpdateDepartment(ctx context.Context, tenantID, id, name, code string, isActive *bool) error {
	setParts := []string{"updated_at = now()"}
	args := []interface{}{}
	argId := 1

	if name != "" {
		setParts = append(setParts, "name = $"+itoa(argId))
		args = append(args, name)
		argId++
	}
	if code != "" {
		setParts = append(setParts, "code = $"+itoa(argId))
		args = append(args, code)
		argId++
	}
	if isActive != nil {
		setParts = append(setParts, "is_active = $"+itoa(argId))
		args = append(args, *isActive)
		argId++
	}

	args = append(args, id, tenantID)
	query := "UPDATE departments SET " + strings.Join(setParts, ", ") + " WHERE id = $" + itoa(argId) + " AND tenant_id = $" + itoa(argId+1)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *SQLDictionaryRepository) DeleteDepartment(ctx context.Context, tenantID, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE departments SET is_active = false, updated_at = now() WHERE id = $1 AND tenant_id = $2`, id, tenantID)
	return err
}

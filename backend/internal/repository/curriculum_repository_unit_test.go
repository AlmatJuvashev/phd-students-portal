package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCurriculumRepository_CreateProgram(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLCurriculumRepository(sqlxDB)

	ctx := context.Background()
	p := &models.Program{
		TenantID:    "tenant-1",
		Code:        "PHD-CS",
		Title:       `{"en": "PhD Computer Science"}`,
		Description: `{"en": "Advanced research"}`,
		Credits:     240,
		DurationMonths: 36,
		IsActive:    true,
	}

	mock.ExpectQuery(`INSERT INTO programs`).
		WithArgs(p.TenantID, p.Code, p.Title, p.Description, p.Credits, p.DurationMonths, p.IsActive).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow("prog-1", time.Now(), time.Now()))

	err = repo.CreateProgram(ctx, p)
	assert.NoError(t, err)
	assert.Equal(t, "prog-1", p.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCurriculumRepository_GetUpdateDeleteProgram(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLCurriculumRepository(sqlxDB)
	ctx := context.Background()

	// GetProgram
	mock.ExpectQuery(`SELECT \* FROM programs WHERE id=\$1`).
		WithArgs("prog-1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "code"}).AddRow("prog-1", "P1"))

	p, err := repo.GetProgram(ctx, "prog-1")
	assert.NoError(t, err)
	assert.Equal(t, "prog-1", p.ID)

	// UpdateProgram
	p.Title = `{"en": "Updated"}`
	mock.ExpectExec(`UPDATE programs SET code=\$1, title=\$2`).
		WithArgs(p.Code, p.Title, p.Description, p.Credits, p.DurationMonths, p.IsActive, p.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateProgram(ctx, p)
	assert.NoError(t, err)

	// DeleteProgram
	mock.ExpectExec(`DELETE FROM programs WHERE id=\$1`).
		WithArgs("prog-1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.DeleteProgram(ctx, "prog-1")
	assert.NoError(t, err)
	
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCurriculumRepository_CreateCourse(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLCurriculumRepository(sqlxDB)

	ctx := context.Background()
	progID := "prog-1"
	c := &models.Course{
		TenantID:    "tenant-1",
		ProgramID:   &progID,
		Code:        "CS101",
		Title:       `{"en": "Intro to AI"}`,
		Description: `{"en": "Basics"}`,
		Credits:     5,
		WorkloadHours: 150,
		IsActive:    true,
	}

	mock.ExpectQuery(`INSERT INTO courses`).
		WithArgs(c.TenantID, c.ProgramID, c.DepartmentID, c.Code, c.Title, c.Description, c.Credits, c.WorkloadHours, c.IsActive).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow("course-1", time.Now(), time.Now()))

	err = repo.CreateCourse(ctx, c)
	assert.NoError(t, err)
	assert.Equal(t, "course-1", c.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCurriculumRepository_GetUpdateDeleteCourse(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLCurriculumRepository(sqlxDB)
	ctx := context.Background()

	// GetCourse
	mock.ExpectQuery(`SELECT \* FROM courses WHERE id=\$1`).
		WithArgs("c1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "code"}).AddRow("c1", "C1"))

	c, err := repo.GetCourse(ctx, "c1")
	assert.NoError(t, err)
	assert.Equal(t, "c1", c.ID)

	// UpdateCourse
	c.Code = "C2"
	mock.ExpectExec(`UPDATE courses SET program_id=\$1`).
		WithArgs(c.ProgramID, c.DepartmentID, c.Code, c.Title, c.Description, c.Credits, c.WorkloadHours, c.IsActive, c.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateCourse(ctx, c)
	assert.NoError(t, err)

	// DeleteCourse
	mock.ExpectExec(`DELETE FROM courses WHERE id=\$1`).
		WithArgs("c1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.DeleteCourse(ctx, "c1")
	assert.NoError(t, err)
	
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCurriculumRepository_CreateJourneyMapAndNodes(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLCurriculumRepository(sqlxDB)

	ctx := context.Background()
	jm := &models.JourneyMap{
		ProgramID: "prog-1",
		Title:     `{"en": "Standard Path"}`,
		Version:   "1.0.0",
		IsActive:  true,
	}

	// 1. Create Journey Map
	mock.ExpectQuery(`INSERT INTO journey_maps`).
		WithArgs(jm.ProgramID, jm.Title, jm.Version, jm.IsActive).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
			AddRow("map-1", time.Now()))

	err = repo.CreateJourneyMap(ctx, jm)
	assert.NoError(t, err)
	assert.Equal(t, "map-1", jm.ID)

	// 2. Create Node Definition
	nd := &models.JourneyNodeDefinition{
		JourneyMapID: "map-1",
		Slug:         "node-1",
		Type:         "task",
		Title:        `{"en": "Submit Proposal"}`,
		Description:  `{"en": "Description"}`,
		ModuleKey:    "I",
		Coordinates:  `{"x": 100, "y": 100}`,
		Config:       `{}`,
		Prerequisites: []string{"node-0"},
	}

	mock.ExpectQuery(`INSERT INTO journey_node_definitions`).
		WithArgs(nd.JourneyMapID, nd.ParentNodeID, nd.Slug, nd.Type, nd.Title, nd.Description, nd.ModuleKey, nd.Coordinates, nd.Config, pq.Array(nd.Prerequisites)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
			AddRow("node-def-1", time.Now()))

	err = repo.CreateNodeDefinition(ctx, nd)
	assert.NoError(t, err)
	assert.Equal(t, "node-def-1", nd.ID)

	// 3. GetJourneyMapByProgram (Found)
	mock.ExpectQuery(`SELECT \* FROM journey_maps WHERE program_id=\$1 LIMIT 1`).
		WithArgs("prog-1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "program_id"}).AddRow("map-1", "prog-1"))
	foundMap, err := repo.GetJourneyMapByProgram(ctx, "prog-1")
	assert.NoError(t, err)
	assert.NotNil(t, foundMap)

	// 4. GetJourneyMapByProgram (Not Found)
	mock.ExpectQuery(`SELECT \* FROM journey_maps WHERE program_id=\$1 LIMIT 1`).
		WithArgs("prog-unknown").
		WillReturnError(sql.ErrNoRows)
	foundMap, err = repo.GetJourneyMapByProgram(ctx, "prog-unknown")
	assert.NoError(t, err) // Should handle NoRows gracefully
	assert.Nil(t, foundMap)
	
	// 5. GetNodeDefinitions
	mock.ExpectQuery(`SELECT \* FROM journey_node_definitions WHERE journey_map_id=\$1`).
		WithArgs("map-1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "slug"}).AddRow("n1", "slug1"))
	nodes, err := repo.GetNodeDefinitions(ctx, "map-1")
	assert.NoError(t, err)
	assert.Len(t, nodes, 1)

	// 6. DeleteNodeDefinition
	mock.ExpectExec(`DELETE FROM journey_node_definitions WHERE id=\$1`).
		WithArgs("n1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.DeleteNodeDefinition(ctx, "n1")
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCurriculumRepository_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLCurriculumRepository(sqlxDB)

	ctx := context.Background()

	// List Programs
	mock.ExpectQuery(`SELECT \* FROM programs WHERE tenant_id=\$1`).
		WithArgs("tenant-1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "code", "title"}).
			AddRow("prog-1", "P1", `{"en": "Title"}`))

	progs, err := repo.ListPrograms(ctx, "tenant-1")
	assert.NoError(t, err)
	assert.Len(t, progs, 1)

	// List Courses
	mock.ExpectQuery(`SELECT \* FROM courses WHERE tenant_id=\$1`).
		WithArgs("tenant-1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "code"}).
			AddRow("c1", "C1"))

	courses, err := repo.ListCourses(ctx, "tenant-1", nil)
	assert.NoError(t, err)
	assert.Len(t, courses, 1)
	
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCurriculumRepository_Cohorts(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLCurriculumRepository(sqlxDB)
	ctx := context.Background()

	c := &models.Cohort{
		ProgramID: "prog-1",
		Name: "Winter 2024",
		IsActive: true,
	}

	// Create
	mock.ExpectQuery(`INSERT INTO cohorts`).
		WithArgs(c.ProgramID, c.Name, c.StartDate, c.EndDate, c.IsActive).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow("coh-1", time.Now()))
	
	err = repo.CreateCohort(ctx, c)
	assert.NoError(t, err)

	// List
	mock.ExpectQuery(`SELECT \* FROM cohorts WHERE program_id=\$1`).
		WithArgs("prog-1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow("coh-1", "Winter 2024"))
	
	list, err := repo.ListCohorts(ctx, "prog-1")
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	assert.NoError(t, mock.ExpectationsWereMet())
}

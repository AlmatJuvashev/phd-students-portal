package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSQLDictionaryRepository_Programs(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLDictionaryRepository(sqlxDB)

	tenantID := "t1"

	t.Run("ListPrograms Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "tenant_id", "code", "name", "title", "description", "credits", "duration_months", "is_active", "created_at", "updated_at"}).
			AddRow("p1", tenantID, "P01", "Program 1", "{}", "{}", 0, 36, true, time.Now(), time.Now())

		mock.ExpectQuery("SELECT (.+) FROM programs WHERE tenant_id = \\$1").
			WithArgs(tenantID).
			WillReturnRows(rows)

		programs, err := repo.ListPrograms(context.Background(), tenantID, false)
		assert.NoError(t, err)
		assert.Len(t, programs, 1)
		assert.Equal(t, "Program 1", programs[0].Name)
	})

	t.Run("CreateProgram Success", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO programs").
			WithArgs("New Program", "NP01", tenantID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("p2"))

		id, err := repo.CreateProgram(context.Background(), tenantID, "New Program", "NP01")
		assert.NoError(t, err)
		assert.Equal(t, "p2", id)
	})

	t.Run("UpdateProgram Success", func(t *testing.T) {
		// update_at = now(), name = $1. id = $2, tenant_id = $3
		mock.ExpectExec("UPDATE programs SET (.+) WHERE id = \\$2 AND tenant_id = \\$3").
			WithArgs("Updated Name", "p1", tenantID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateProgram(context.Background(), tenantID, "p1", "Updated Name", "", nil)
		assert.NoError(t, err)
	})

	t.Run("DeleteProgram Success", func(t *testing.T) {
		mock.ExpectExec("UPDATE programs SET is_active = false").
			WithArgs("p1", tenantID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeleteProgram(context.Background(), tenantID, "p1")
		assert.NoError(t, err)
	})
}

func TestSQLDictionaryRepository_Specialties(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLDictionaryRepository(sqlxDB)

	tenantID := "t1"

	t.Run("ListSpecialties Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "code", "is_active", "created_at", "updated_at"}).
			AddRow("s1", "Specialty 1", "S01", true, "2023-01-01T00:00:00Z", "2023-01-01T00:00:00Z")

		mock.ExpectQuery("SELECT (.+) FROM specialties WHERE tenant_id = \\$1").
			WithArgs(tenantID).
			WillReturnRows(rows)
		
		mock.ExpectQuery("SELECT program_id FROM specialty_programs WHERE specialty_id = \\$1").
			WithArgs("s1").
			WillReturnRows(sqlmock.NewRows([]string{"program_id"}).AddRow("p1"))

		specialties, err := repo.ListSpecialties(context.Background(), tenantID, false, "")
		assert.NoError(t, err)
		assert.Len(t, specialties, 1)
		assert.Equal(t, "Specialty 1", specialties[0].Name)
		assert.Equal(t, []string{"p1"}, specialties[0].ProgramIDs)
	})

	t.Run("CreateSpecialty Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO specialties").
			WithArgs("New Spec", "NS01", tenantID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("s2"))
		
		mock.ExpectExec("INSERT INTO specialty_programs").
			WithArgs("s2", "p1").
			WillReturnResult(sqlmock.NewResult(0, 1))
		
		mock.ExpectCommit()

		id, err := repo.CreateSpecialty(context.Background(), tenantID, "New Spec", "NS01", []string{"p1"})
		assert.NoError(t, err)
		assert.Equal(t, "s2", id)
	})

	t.Run("UpdateSpecialty Success", func(t *testing.T) {
		mock.ExpectBegin()
		// Update basic fields
		mock.ExpectExec("UPDATE specialties SET (.+) WHERE id = \\$2 AND tenant_id = \\$3").
			WithArgs("Updated Spec", "s1", tenantID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		
		// Update programs relations
		mock.ExpectExec("DELETE FROM specialty_programs WHERE specialty_id = \\$1").
			WithArgs("s1").
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec("INSERT INTO specialty_programs").
			WithArgs("s1", "p2").
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		err := repo.UpdateSpecialty(context.Background(), tenantID, "s1", "Updated Spec", "", nil, []string{"p2"})
		assert.NoError(t, err)
	})

	t.Run("DeleteSpecialty Success", func(t *testing.T) {
		mock.ExpectExec("UPDATE specialties SET is_active = false").
			WithArgs("s1", tenantID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeleteSpecialty(context.Background(), tenantID, "s1")
		assert.NoError(t, err)
	})
}

func TestSQLDictionaryRepository_Cohorts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLDictionaryRepository(sqlxDB)

	tenantID := "t1"

	t.Run("ListCohorts Success", func(t *testing.T) {
		startTime, _ := time.Parse("2006-01-02", "2023-01-01")
		endTime, _ := time.Parse("2006-01-02", "2023-12-31")
		rows := sqlmock.NewRows([]string{"id", "name", "start_date", "end_date", "is_active", "created_at", "updated_at"}).
			AddRow("c1", "Cohort 2023", startTime, endTime, true, startTime, startTime)

		mock.ExpectQuery("SELECT (.+) FROM cohorts WHERE tenant_id = \\$1").
			WithArgs(tenantID).
			WillReturnRows(rows)

		cohorts, err := repo.ListCohorts(context.Background(), tenantID, false)
		assert.NoError(t, err)
		assert.Len(t, cohorts, 1)
	})

	t.Run("CreateCohort Success", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO cohorts").
			WithArgs("Cohort 2024", "2024-01-01", "2024-12-31", tenantID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("c2"))

		id, err := repo.CreateCohort(context.Background(), tenantID, "Cohort 2024", "2024-01-01", "2024-12-31")
		assert.NoError(t, err)
		assert.Equal(t, "c2", id)
	})

	t.Run("UpdateCohort Success", func(t *testing.T) {
		mock.ExpectExec("UPDATE cohorts SET (.+) WHERE id = \\$2 AND tenant_id = \\$3").
			WithArgs("Updated Cohort", "c1", tenantID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateCohort(context.Background(), tenantID, "c1", "Updated Cohort", "", "", nil)
		assert.NoError(t, err)
	})

	t.Run("DeleteCohort Success", func(t *testing.T) {
		mock.ExpectExec("UPDATE cohorts SET is_active = false").
			WithArgs("c1", tenantID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeleteCohort(context.Background(), tenantID, "c1")
		assert.NoError(t, err)
	})
}

func TestSQLDictionaryRepository_Departments(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLDictionaryRepository(sqlxDB)

	tenantID := "t1"

	t.Run("ListDepartments Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "code", "is_active", "created_at", "updated_at"}).
			AddRow("d1", "Department 1", "D01", true, "2023-01-01T00:00:00Z", "2023-01-01T00:00:00Z")

		mock.ExpectQuery("SELECT (.+) FROM departments WHERE tenant_id = \\$1").
			WithArgs(tenantID).
			WillReturnRows(rows)

		departments, err := repo.ListDepartments(context.Background(), tenantID, false)
		assert.NoError(t, err)
		assert.Len(t, departments, 1)
	})

	t.Run("CreateDepartment Success", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO departments").
			WithArgs("New Dept", "ND01", tenantID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("d2"))

		id, err := repo.CreateDepartment(context.Background(), tenantID, "New Dept", "ND01")
		assert.NoError(t, err)
		assert.Equal(t, "d2", id)
	})

	t.Run("UpdateDepartment Success", func(t *testing.T) {
		mock.ExpectExec("UPDATE departments SET (.+) WHERE id = \\$2 AND tenant_id = \\$3").
			WithArgs("Updated Dept", "d1", tenantID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateDepartment(context.Background(), tenantID, "d1", "Updated Dept", "", nil)
		assert.NoError(t, err)
	})

	t.Run("DeleteDepartment Success", func(t *testing.T) {
		mock.ExpectExec("UPDATE departments SET is_active = false").
			WithArgs("d1", tenantID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeleteDepartment(context.Background(), tenantID, "d1")
		assert.NoError(t, err)
	})
}

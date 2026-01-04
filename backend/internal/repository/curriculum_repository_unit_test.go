package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func newTestCurriculumRepo(t *testing.T) (*SQLCurriculumRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLCurriculumRepository(sqlxDB)
	return repo, mock, func() { db.Close() }
}

func TestSQLCurriculumRepository_Programs(t *testing.T) {
	repo, mock, teardown := newTestCurriculumRepo(t)
	defer teardown()
	ctx := context.Background()

	t.Run("CreateProgram", func(t *testing.T) {
		prog := &models.Program{
			TenantID:       "t1",
			Code:           "PHD-CS",
			Name:           "PHD-CS",
			Title:          "PhD Computer Science",
			Description:    "Desc",
			Credits:        180,
			DurationMonths: 48,
			IsActive:       true,
		}

		mock.ExpectQuery("INSERT INTO programs").
			WithArgs("t1", "PHD-CS", "PHD-CS", "PhD Computer Science", "Desc", 180, 48, true).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow("p1", time.Now(), time.Now()))

		err := repo.CreateProgram(ctx, prog)
		assert.NoError(t, err)
		assert.Equal(t, "p1", prog.ID)
	})

	t.Run("GetProgram", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "code", "title"}).
			AddRow("p1", "PHD-CS", "PhD Computer Science")
		mock.ExpectQuery("SELECT \\* FROM programs WHERE id=\\$1").
			WithArgs("p1").
			WillReturnRows(rows)

		p, err := repo.GetProgram(ctx, "p1")
		assert.NoError(t, err)
		assert.Equal(t, "PHD-CS", p.Code)
	})

	t.Run("ListPrograms", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "code"}).
			AddRow("p1", "PHD-CS").
			AddRow("p2", "PHD-MATH")
		mock.ExpectQuery("SELECT \\* FROM programs WHERE tenant_id=\\$1 ORDER BY created_at DESC").
			WithArgs("t1").
			WillReturnRows(rows)

		list, err := repo.ListPrograms(ctx, "t1")
		assert.NoError(t, err)
		assert.Len(t, list, 2)
	})

	t.Run("UpdateProgram", func(t *testing.T) {
		prog := &models.Program{
			ID:             "p1",
			Code:           "PHD-CS-V2",
			Title:          "Updated Title",
			Description:    "New Desc",
			Credits:        240,
			DurationMonths: 60,
			IsActive:       false,
		}
		mock.ExpectExec("UPDATE programs SET").
			WithArgs("PHD-CS-V2", "Updated Title", "New Desc", 240, 60, false, "p1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateProgram(ctx, prog)
		assert.NoError(t, err)
	})

	t.Run("DeleteProgram", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM programs WHERE id=\\$1").
			WithArgs("p1").
			WillReturnResult(sqlmock.NewResult(1, 1))
		err := repo.DeleteProgram(ctx, "p1")
		assert.NoError(t, err)
	})
}

func TestSQLCurriculumRepository_Courses(t *testing.T) {
	repo, mock, teardown := newTestCurriculumRepo(t)
	defer teardown()
	ctx := context.Background()

	t.Run("CreateCourse", func(t *testing.T) {
		course := &models.Course{
			TenantID:      "t1",
			ProgramID:     nil,
			DepartmentID:  nil,
			Code:          "CS101",
			Title:         "Intro",
			Description:   "Desc",
			Credits:       6,
			WorkloadHours: 120,
			IsActive:      true,
		}

		mock.ExpectQuery("INSERT INTO courses").
			WithArgs("t1", nil, nil, "CS101", "Intro", "Desc", 6, 120, true).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow("c1", time.Now(), time.Now()))

		err := repo.CreateCourse(ctx, course)
		assert.NoError(t, err)
		assert.Equal(t, "c1", course.ID)
	})

	t.Run("GetCourse", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "code"}).AddRow("c1", "CS101")
		mock.ExpectQuery("SELECT \\* FROM courses WHERE id=\\$1").
			WithArgs("c1").
			WillReturnRows(rows)
		c, err := repo.GetCourse(ctx, "c1")
		assert.NoError(t, err)
		assert.Equal(t, "CS101", c.Code)
	})

	t.Run("ListCourses", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "code"}).AddRow("c1", "CS101")
		mock.ExpectQuery("SELECT \\* FROM courses WHERE tenant_id=\\$1 ORDER BY code ASC").
			WithArgs("t1").
			WillReturnRows(rows)

		list, err := repo.ListCourses(ctx, "t1", nil)
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})

	t.Run("UpdateCourse", func(t *testing.T) {
		course := &models.Course{
			ID:          "c1",
			Code:        "CS101-V2",
			Title:       "Updated",
			Description: "Desc",
		}
		mock.ExpectExec("UPDATE courses SET").
			WithArgs(nil, nil, "CS101-V2", "Updated", "Desc", 0, 0, false, "c1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateCourse(ctx, course)
		assert.NoError(t, err)
	})

	t.Run("DeleteCourse", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM courses WHERE id=\\$1").
			WithArgs("c1").
			WillReturnResult(sqlmock.NewResult(1, 1))
		err := repo.DeleteCourse(ctx, "c1")
		assert.NoError(t, err)
	})
}

func TestSQLCurriculumRepository_JourneyMaps(t *testing.T) {
	repo, mock, teardown := newTestCurriculumRepo(t)
	defer teardown()
	ctx := context.Background()

	t.Run("CreateJourneyMap", func(t *testing.T) {
		jm := &models.JourneyMap{
			ProgramID: "p1",
			Title:     "Map 1",
			Version:   "1.0",
			Config:    "{}",
			IsActive:  true,
		}
		mock.ExpectQuery("INSERT INTO program_versions").
			WithArgs("p1", "Map 1", "1.0", "{}", true).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow("jm1", time.Now(), time.Now()))

		err := repo.CreateJourneyMap(ctx, jm)
		assert.NoError(t, err)
		assert.Equal(t, "jm1", jm.ID)
	})

	t.Run("GetJourneyMapByProgram", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM program_versions").
			WithArgs("p1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow("jm1", "Map 1"))

		jm, err := repo.GetJourneyMapByProgram(ctx, "p1")
		assert.NoError(t, err)
		assert.Equal(t, "Map 1", jm.Title)
	})

	t.Run("GetJourneyMapByProgram_NotFound", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM program_versions").
			WithArgs("p1").
			WillReturnError(sql.ErrNoRows)

		jm, err := repo.GetJourneyMapByProgram(ctx, "p1")
		assert.NoError(t, err)
		assert.Nil(t, jm)
	})
}

func TestSQLCurriculumRepository_NodeDefinitions(t *testing.T) {
	repo, mock, teardown := newTestCurriculumRepo(t)
	defer teardown()
	ctx := context.Background()

	t.Run("CreateNodeDefinition", func(t *testing.T) {
		nd := &models.JourneyNodeDefinition{
			JourneyMapID:  "jm1",
			Slug:          "node-1",
			Type:          "task",
			Title:         "Node 1",
			Description:   "Desc",
			ModuleKey:     "I",
			Coordinates:   `{"x":0,"y":0}`,
			Config:        "{}",
			Prerequisites: []string{"start"},
		}
		mock.ExpectQuery("INSERT INTO program_version_node_definitions").
			WithArgs("jm1", nil, "node-1", "task", "Node 1", "Desc", "I", `{"x":0,"y":0}`, "{}", pq.Array([]string{"start"})).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow("nd1", time.Now(), time.Now()))

		err := repo.CreateNodeDefinition(ctx, nd)
		assert.NoError(t, err)
		assert.Equal(t, "nd1", nd.ID)
	})

	t.Run("GetNodeDefinitions", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM program_version_node_definitions WHERE program_version_id=\\$1").
			WithArgs("jm1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "slug"}).AddRow("nd1", "node-1"))

		nodes, err := repo.GetNodeDefinitions(ctx, "jm1")
		assert.NoError(t, err)
		assert.Len(t, nodes, 1)
	})
	
	t.Run("GetNodeDefinition", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM program_version_node_definitions WHERE id=\\$1").
			WithArgs("nd1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "slug"}).AddRow("nd1", "node-1"))
		
		node, err := repo.GetNodeDefinition(ctx, "nd1")
		assert.NoError(t, err)
		assert.Equal(t, "node-1", node.Slug)
	})
	
	t.Run("UpdateNodeDefinition", func(t *testing.T) {
		nd := &models.JourneyNodeDefinition{
			ID: "nd1",
			Title: "Updated",
			Prerequisites: []string{"a"},
		}
		mock.ExpectExec("UPDATE program_version_node_definitions").
			WithArgs("Updated", "", "", "", pq.Array([]string{"a"}), "", "", "nd1").
			WillReturnResult(sqlmock.NewResult(1, 1))
			
		err := repo.UpdateNodeDefinition(ctx, nd)
		assert.NoError(t, err)
	})
	
	t.Run("DeleteNodeDefinition", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM program_version_node_definitions WHERE id=\\$1").
			WithArgs("nd1").
			WillReturnResult(sqlmock.NewResult(1, 1))
		err := repo.DeleteNodeDefinition(ctx, "nd1")
		assert.NoError(t, err)
	})
}

func TestSQLCurriculumRepository_Cohorts(t *testing.T) {
	repo, mock, teardown := newTestCurriculumRepo(t)
	defer teardown()
	ctx := context.Background()

	t.Run("CreateCohort", func(t *testing.T) {
		cohort := &models.Cohort{
			ProgramID: "p1",
			Name:      "Winter 2024",
			IsActive:  true,
		}
		mock.ExpectQuery("INSERT INTO cohorts").
			WithArgs("p1", "Winter 2024", sqlmock.AnyArg(), sqlmock.AnyArg(), true).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow("c1", time.Now()))

		err := repo.CreateCohort(ctx, cohort)
		assert.NoError(t, err)
		assert.Equal(t, "c1", cohort.ID)
	})

	t.Run("ListCohorts", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM cohorts WHERE program_id=\\$1").
			WithArgs("p1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow("c1", "Winter 2024"))

		list, err := repo.ListCohorts(ctx, "p1")
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})
}

func TestSQLCurriculumRepository_CourseRequirements(t *testing.T) {
	repo, mock, teardown := newTestCurriculumRepo(t)
	defer teardown()
	ctx := context.Background()

	t.Run("SetCourseRequirement_Success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO course_requirements").
			WithArgs("course-1", "Lab", "true").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.SetCourseRequirement(ctx, &models.CourseRequirement{CourseID: "course-1", Key: "Lab", Value: "true"})
		assert.NoError(t, err)
	})

	t.Run("SetCourseRequirement_Error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO course_requirements").
			WithArgs("course-1", "Lab", "true").
			WillReturnError(fmt.Errorf("db error"))

		err := repo.SetCourseRequirement(ctx, &models.CourseRequirement{CourseID: "course-1", Key: "Lab", Value: "true"})
		assert.Error(t, err)
	})

	t.Run("GetCourseRequirements_Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"course_id", "key", "value"}).
			AddRow("course-1", "Lab", "true").
			AddRow("course-1", "OS", "Linux")

		mock.ExpectQuery(`SELECT \* FROM course_requirements WHERE course_id=\$1`).
			WithArgs("course-1").
			WillReturnRows(rows)

		reqs, err := repo.GetCourseRequirements(ctx, "course-1")
		assert.NoError(t, err)
		assert.Len(t, reqs, 2)
		assert.Equal(t, "Lab", reqs[0].Key)
		assert.Equal(t, "true", reqs[0].Value)
	})
}

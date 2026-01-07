package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func newMockRubricRepo(t *testing.T) (*SQLRubricRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLRubricRepository(sqlxDB)

	return repo, mock, func() {
		db.Close()
	}
}

func TestSQLRubricRepository_CreateRubric(t *testing.T) {
	repo, mock, teardown := newMockRubricRepo(t)
	defer teardown()

	ctx := context.Background()

	rubric := &models.Rubric{
		CourseOfferingID: "offering-1",
		Title:            "Final Essay Rubric",
		Description:      "Rubric for final assignment",
		IsGlobal:         false,
		Criteria: []models.RubricCriterion{
			{
				Title:       "Grammar",
				Description: "Grammar usage",
				Weight:      30,
				Levels: []models.RubricLevel{
					{Title: "Excellent", Points: 10, Position: 0},
				},
				Position: 0,
			},
		},
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		
		// 1. Insert Rubric
		mock.ExpectQuery(`INSERT INTO rubrics`).
			WithArgs(rubric.CourseOfferingID, rubric.Title, rubric.Description, rubric.IsGlobal, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("rubric-1"))

		// 2. Insert Criterion
		mock.ExpectQuery(`INSERT INTO rubric_criteria`).
			WithArgs("rubric-1", "Grammar", "Grammar usage", 30.0, 0, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("crit-1"))

		// 3. Insert Level
		mock.ExpectExec(`INSERT INTO rubric_levels`).
			WithArgs("crit-1", "Excellent", "", 10.0, 0, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := repo.CreateRubric(ctx, rubric)
		assert.NoError(t, err)
		assert.Equal(t, "rubric-1", rubric.ID)
		assert.Equal(t, "crit-1", rubric.Criteria[0].ID)
	})
	
	t.Run("Rollback on Error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO rubrics`).WillReturnError(fmt.Errorf("db error"))
		mock.ExpectRollback()

		err := repo.CreateRubric(ctx, rubric)
		assert.Error(t, err)
	})
}

func TestSQLRubricRepository_GetRubric(t *testing.T) {
	repo, mock, teardown := newMockRubricRepo(t)
	defer teardown()
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// 1. Fetches Rubric
		rubricRow := sqlmock.NewRows([]string{"id", "title"}).AddRow("r1", "Essay Rubric")
		mock.ExpectQuery(`SELECT \* FROM rubrics WHERE id=\$1`).
			WithArgs("r1").
			WillReturnRows(rubricRow)
		
		// 2. Fetches Criteria
		critRow := sqlmock.NewRows([]string{"id", "rubric_id", "title"}).AddRow("c1", "r1", "Grammar")
		mock.ExpectQuery(`SELECT \* FROM rubric_criteria WHERE rubric_id=\$1 ORDER BY position ASC`).
			WithArgs("r1").
			WillReturnRows(critRow)

		// 3. Fetches Levels for Criteria
		levelRow := sqlmock.NewRows([]string{"id", "criterion_id", "title"}).AddRow("l1", "c1", "Good")
		// sqlx.In logic: translates `criterion_id IN (?)` to `criterion_id IN (?)` for sqlmock driver
		mock.ExpectQuery(`SELECT \* FROM rubric_levels WHERE criterion_id IN \(\?\) ORDER BY position ASC`).
			WithArgs("c1").
			WillReturnRows(levelRow)

		r, err := repo.GetRubric(ctx, "r1")
		assert.NoError(t, err)
		assert.Equal(t, "Essay Rubric", r.Title)
		assert.Len(t, r.Criteria, 1)
		assert.Len(t, r.Criteria[0].Levels, 1)
	})
}

func TestSQLRubricRepository_ListRubrics(t *testing.T) {
	repo, mock, teardown := newMockRubricRepo(t)
	defer teardown()
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title"}).AddRow("r1", "R1").AddRow("r2", "R2")
		mock.ExpectQuery(`SELECT \* FROM rubrics WHERE course_offering_id=\$1`).
			WithArgs("off-1").
			WillReturnRows(rows)

		list, err := repo.ListRubrics(ctx, "off-1")
		assert.NoError(t, err)
		assert.Len(t, list, 2)
	})
}

func TestSQLRubricRepository_SubmitGrade(t *testing.T) {
	repo, mock, teardown := newMockRubricRepo(t)
	defer teardown()
	ctx := context.Background()

	grade := &models.RubricGrade{
		SubmissionID: "sub-1",
		RubricID:     "rubric-1",
		GraderID:     strPtr("grader-1"),
		TotalScore:   95.0,
		Comments:     strPtr("Good job"),
		Items: []models.RubricGradeItem{
			{CriterionID: "c1", LevelID: strPtr("l1"), PointsAwarded: 10},
		},
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		
		// 1. Insert Grade Header
		mock.ExpectQuery(`INSERT INTO rubric_grades`).
			WithArgs(grade.SubmissionID, grade.RubricID, grade.GraderID, grade.TotalScore, grade.Comments, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("grade-1"))

		// 2. Insert Item
		mock.ExpectExec(`INSERT INTO rubric_grade_items`).
			WithArgs("grade-1", "c1", "l1", 10.0, nil, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// 3. Update Submission Status
		mock.ExpectExec(`UPDATE activity_submissions SET status='GRADED'`).
			WithArgs(grade.TotalScore, grade.SubmissionID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := repo.SubmitGrade(ctx, grade)
		assert.NoError(t, err)
		assert.Equal(t, "grade-1", grade.ID)
	})
}

func TestSQLRubricRepository_GetGrade(t *testing.T) {
	repo, mock, teardown := newMockRubricRepo(t)
	defer teardown()
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// 1. Get Grade
		rows := sqlmock.NewRows([]string{"id", "total_score"}).AddRow("grade-1", 90.0)
		mock.ExpectQuery(`SELECT \* FROM rubric_grades WHERE submission_id=\$1`).
			WithArgs("sub-1").
			WillReturnRows(rows)

		// 2. Get Items
		itemRows := sqlmock.NewRows([]string{"id", "points_awarded"}).AddRow("item-1", 10.0)
		mock.ExpectQuery(`SELECT \* FROM rubric_grade_items WHERE rubric_grade_id=\$1`).
			WithArgs("grade-1").
			WillReturnRows(itemRows)

		g, err := repo.GetGrade(ctx, "sub-1")
		assert.NoError(t, err)
		assert.Equal(t, 90.0, g.TotalScore)
		assert.Len(t, g.Items, 1)
	})
}



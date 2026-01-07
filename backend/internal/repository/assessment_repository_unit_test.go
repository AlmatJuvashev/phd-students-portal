package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/assert"
)

func newMockAssessmentRepo(t *testing.T) (AssessmentRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSQLAssessmentRepository(sqlxDB)

	return repo, mock, func() {
		db.Close()
	}
}

func TestSQLAssessmentRepository_QuestionBanks(t *testing.T) {
	repo, mock, teardown := newMockAssessmentRepo(t)
	defer teardown()

	ctx := context.Background()
	bank := models.QuestionBank{
		TenantID:       "t1",
		Title:         "Bank 1",
		Description:   toPtr("Desc"),
		Subject:       toPtr("Sub"),
		BloomsTaxonomy: toPtr(models.BloomsApplication),
		IsPublic:      true,
		CreatedBy:     "u1",
	}

	t.Run("CreateQuestionBank", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "tenant_id", "title"}).
			AddRow("b1", "t1", "Bank 1")
		
		mock.ExpectQuery("INSERT INTO question_banks").
			WithArgs(bank.TenantID, bank.Title, bank.Description, bank.Subject, bank.BloomsTaxonomy, bank.IsPublic, bank.CreatedBy).
			WillReturnRows(rows)

		res, err := repo.CreateQuestionBank(ctx, bank)
		assert.NoError(t, err)
		assert.Equal(t, "b1", res.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetQuestionBank", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM question_banks WHERE id = \\$1").
			WithArgs("b1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow("b1", "Bank 1"))

		res, err := repo.GetQuestionBank(ctx, "b1")
		assert.NoError(t, err)
		assert.Equal(t, "Bank 1", res.Title)
	})

	t.Run("ListQuestionBanks", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM question_banks WHERE tenant_id = \\$1").
			WithArgs("t1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow("b1", "Bank 1").AddRow("b2", "Bank 2"))

		res, err := repo.ListQuestionBanks(ctx, "t1")
		assert.NoError(t, err)
		assert.Len(t, res, 2)
	})

	t.Run("UpdateQuestionBank", func(t *testing.T) {
		mock.ExpectExec("UPDATE question_banks").
			WithArgs(bank.Title, bank.Description, bank.Subject, bank.BloomsTaxonomy, bank.IsPublic, "b1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		bank.ID = "b1"
		err := repo.UpdateQuestionBank(ctx, bank)
		assert.NoError(t, err)
	})

	t.Run("DeleteQuestionBank", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM question_banks WHERE id=\\$1").
			WithArgs("b1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.DeleteQuestionBank(ctx, "b1")
		assert.NoError(t, err)
	})
}

func TestSQLAssessmentRepository_Questions(t *testing.T) {
	repo, mock, teardown := newMockAssessmentRepo(t)
	defer teardown()

	ctx := context.Background()
	q := models.Question{
		BankID:          "b1",
		Type:            models.QuestionTypeMCQ,
		Stem:            "Question ?",
		PointsDefault:   10,
		DifficultyLevel: toPtr(models.DifficultyMedium),
		Options: []models.QuestionOption{
			{Text: "Opt 1", IsCorrect: true},
			{Text: "Opt 2", IsCorrect: false},
		},
	}

	t.Run("CreateQuestion", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO questions").
			WithArgs(q.BankID, q.Type, q.Stem, q.MediaURL, q.PointsDefault, q.DifficultyLevel, q.LearningOutcomeID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("q1"))
		
		mock.ExpectExec("INSERT INTO question_options").
			WithArgs("q1", q.Options[0].Text, q.Options[0].IsCorrect, 0, q.Options[0].Feedback).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO question_options").
			WithArgs("q1", q.Options[1].Text, q.Options[1].IsCorrect, 1, q.Options[1].Feedback).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		res, err := repo.CreateQuestion(ctx, q)
		assert.NoError(t, err)
		assert.Equal(t, "q1", res.ID)
	})

	t.Run("GetQuestion", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM questions WHERE id = \\$1").
			WithArgs("q1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "stem"}).AddRow("q1", "Question ?"))
		
		mock.ExpectQuery("SELECT \\* FROM question_options WHERE question_id = \\$1").
			WithArgs("q1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "text"}).AddRow("o1", "Opt 1").AddRow("o2", "Opt 2"))

		res, err := repo.GetQuestion(ctx, "q1")
		assert.NoError(t, err)
		assert.Len(t, res.Options, 2)
	})

	t.Run("UpdateQuestion", func(t *testing.T) {
		q.ID = "q1"
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE questions SET").
			WithArgs(q.Type, q.Stem, q.MediaURL, q.PointsDefault, q.DifficultyLevel, q.LearningOutcomeID, "q1").
			WillReturnResult(sqlmock.NewResult(1, 1))
		
		mock.ExpectExec("DELETE FROM question_options WHERE question_id=\\$1").
			WithArgs("q1").
			WillReturnResult(sqlmock.NewResult(1, 1))
		
		mock.ExpectExec("INSERT INTO question_options").
			WithArgs("q1", q.Options[0].Text, q.Options[0].IsCorrect, 0, q.Options[0].Feedback).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO question_options").
			WithArgs("q1", q.Options[1].Text, q.Options[1].IsCorrect, 1, q.Options[1].Feedback).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.UpdateQuestion(ctx, q)
		assert.NoError(t, err)
	})

	t.Run("DeleteQuestion", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM questions WHERE id=\\$1").
			WithArgs("q1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.DeleteQuestion(ctx, "q1")
		assert.NoError(t, err)
	})
}

func TestSQLAssessmentRepository_Assessments(t *testing.T) {
	repo, mock, teardown := newMockAssessmentRepo(t)
	defer teardown()

	ctx := context.Background()
	a := models.Assessment{
		TenantID:          "t1",
		CourseOfferingID:  "co1",
		Title:             "Quiz 1",
		TimeLimitMinutes: toPtr(30),
		PassingScore:      60,
		SecuritySettings: types.JSONText("{}"),
	}

	t.Run("CreateAssessment", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO assessments").
			WithArgs(a.TenantID, a.CourseOfferingID, a.Title, a.Description, a.TimeLimitMinutes, a.AvailableFrom, a.AvailableUntil, a.ShuffleQuestions, a.GradingPolicy, a.SecuritySettings, a.PassingScore, a.CreatedBy).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("a1"))

		res, err := repo.CreateAssessment(ctx, a)
		assert.NoError(t, err)
		assert.Equal(t, "a1", res.ID)
	})

	t.Run("GetAssessment", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM assessments WHERE id = \\$1").
			WithArgs("a1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow("a1", "Quiz 1"))

		res, err := repo.GetAssessment(ctx, "a1")
		assert.NoError(t, err)
		assert.Equal(t, "Quiz 1", res.Title)
	})

	t.Run("ListAssessments", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM assessments WHERE tenant_id=\\$1 AND course_offering_id=\\$2").
			WithArgs("t1", "co1").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("a1"))

		res, err := repo.ListAssessments(ctx, "t1", "co1")
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("UpdateAssessment", func(t *testing.T) {
		a.ID = "a1"
		mock.ExpectExec("UPDATE assessments SET").
			WithArgs(a.CourseOfferingID, a.Title, a.Description, a.TimeLimitMinutes, a.AvailableFrom, a.AvailableUntil, a.ShuffleQuestions, a.GradingPolicy, a.SecuritySettings, a.PassingScore, "a1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateAssessment(ctx, a)
		assert.NoError(t, err)
	})

	t.Run("DeleteAssessment", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM assessments WHERE id=\\$1").
			WithArgs("a1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.DeleteAssessment(ctx, "a1")
		assert.NoError(t, err)
	})
}

func TestSQLAssessmentRepository_Attempts(t *testing.T) {
	repo, mock, teardown := newMockAssessmentRepo(t)
	defer teardown()

	ctx := context.Background()

	t.Run("CreateAttempt", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO assessment_attempts").
			WithArgs("a1", "s1", models.AttemptStatusInProgress).
			WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow("at1", models.AttemptStatusInProgress))

		res, err := repo.CreateAttempt(ctx, models.AssessmentAttempt{AssessmentID: "a1", StudentID: "s1"})
		assert.NoError(t, err)
		assert.Equal(t, "at1", res.ID)
	})

	t.Run("SaveItemResponse", func(t *testing.T) {
		resp := models.ItemResponse{
			AttemptID:  "at1",
			QuestionID: "q1",
			Score:      5,
			IsCorrect:  true,
		}
		mock.ExpectExec("INSERT INTO item_responses").
			WithArgs(resp.AttemptID, resp.QuestionID, resp.SelectedOptionID, resp.TextResponse, resp.Score, resp.IsCorrect).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.SaveItemResponse(ctx, resp)
		assert.NoError(t, err)
	})

	t.Run("CompleteAttempt", func(t *testing.T) {
		mock.ExpectExec("UPDATE assessment_attempts SET finished_at=NOW\\(\\), score=\\$1, status=\\$2 WHERE id=\\$3").
			WithArgs(85.5, models.AttemptStatusSubmitted, "at1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.CompleteAttempt(ctx, "at1", 85.5)
		assert.NoError(t, err)
	})

	t.Run("LogProctoringEvent", func(t *testing.T) {
		log := models.ProctoringLog{
			AttemptID: "at1",
			EventType: models.ProctoringEventTabSwitch,
			Metadata:  types.JSONText(`{"count": 1}`),
		}
		mock.ExpectExec("INSERT INTO proctoring_logs").
			WithArgs(log.AttemptID, log.EventType, log.Metadata).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.LogProctoringEvent(ctx, log)
		assert.NoError(t, err)
	})
}

func TestSQLAssessmentRepository_Complex(t *testing.T) {
	repo, mock, teardown := newMockAssessmentRepo(t)
	defer teardown()

	ctx := context.Background()

	t.Run("GetAssessmentQuestions", func(t *testing.T) {
		// 1. Fetch Questions
		mock.ExpectQuery("SELECT q.* FROM questions q").
			WithArgs("a1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "stem"}).
				AddRow("q1", "Q1").
				AddRow("q2", "Q2"))
		
		// 2. Fetch Options
		mock.ExpectQuery("SELECT \\* FROM question_options WHERE question_id IN \\(\\$1, \\$2\\)").
			WithArgs("q1", "q2").
			WillReturnRows(sqlmock.NewRows([]string{"id", "question_id", "text"}).
				AddRow("o1", "q1", "Opt 1").
				AddRow("o2", "q2", "Opt 2"))

		res, err := repo.GetAssessmentQuestions(ctx, "a1")
		assert.NoError(t, err)
		assert.Len(t, res, 2)
		assert.Len(t, res[0].Options, 1)
		assert.Len(t, res[1].Options, 1)
	})

	t.Run("GetAssessmentQuestions - No Rows", func(t *testing.T) {
		mock.ExpectQuery("SELECT q.* FROM questions q").
			WithArgs("a1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "stem"}))

		res, err := repo.GetAssessmentQuestions(ctx, "a1")
		assert.NoError(t, err)
		assert.Nil(t, res)
	})

	t.Run("GetAssessmentQuestions - DB Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT q.* FROM questions q").
			WithArgs("a1").
			WillReturnError(sql.ErrConnDone)

		_, err := repo.GetAssessmentQuestions(ctx, "a1")
		assert.Error(t, err)
	})
}

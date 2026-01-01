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

func TestCourseContentRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLCourseContentRepository(sqlxDB)
	ctx := context.Background()

	// --- Module ---
	m := &models.CourseModule{CourseID: "c1", Title: "M1", Order: 1, IsActive: true}
	mock.ExpectQuery(`INSERT INTO course_modules`).
		WithArgs(m.CourseID, m.Title, m.Order, m.IsActive).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow("m1", time.Now(), time.Now()))
	
	err = repo.CreateModule(ctx, m)
	assert.NoError(t, err)
	assert.Equal(t, "m1", m.ID)

	mock.ExpectQuery(`SELECT \* FROM course_modules WHERE course_id=\$1`).
		WithArgs("c1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow("m1", "M1"))
	listM, err := repo.ListModules(ctx, "c1")
	assert.NoError(t, err)
	assert.Len(t, listM, 1)

	// --- Lesson ---
	l := &models.CourseLesson{ModuleID: "m1", Title: "L1", Order: 1, IsActive: true}
	mock.ExpectQuery(`INSERT INTO course_lessons`).
		WithArgs(l.ModuleID, l.Title, l.Order, l.IsActive).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow("l1", time.Now(), time.Now()))
	
	err = repo.CreateLesson(ctx, l)
	assert.NoError(t, err)

	// --- Activity ---
	a := &models.CourseActivity{
		LessonID: "l1", 
		Type: "video", 
		Title: "Intro Video", 
		Order: 1, 
		Points: 10, 
		IsOptional: false, 
		Content: `{"videoUrls": ["http://a.com"]}`,
	}
	mock.ExpectQuery(`INSERT INTO course_activities`).
		WithArgs(a.LessonID, a.Type, a.Title, a.Order, a.Points, a.IsOptional, a.Content).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow("a1", time.Now(), time.Now()))
	
	err = repo.CreateActivity(ctx, a)
	assert.NoError(t, err)

	// --- Update & Get & Delete Module ---
	m.Title = "M1 Updated"
	mock.ExpectExec(`UPDATE course_modules`).
		WithArgs(m.Title, m.Order, m.IsActive, m.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.UpdateModule(ctx, m)
	assert.NoError(t, err)

	mock.ExpectQuery(`SELECT \* FROM course_modules WHERE id=\$1`).
		WithArgs("m1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow("m1", "M1 Updated"))
	gotM, err := repo.GetModule(ctx, "m1")
	assert.NoError(t, err)
	assert.Equal(t, "M1 Updated", gotM.Title)

	mock.ExpectExec(`DELETE FROM course_modules`).
		WithArgs("m1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.DeleteModule(ctx, "m1")
	assert.NoError(t, err)

	// --- Update & Get & Delete Lesson ---
	l.Title = "L1 Updated"
	mock.ExpectExec(`UPDATE course_lessons`).
		WithArgs(l.Title, l.Order, l.IsActive, l.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.UpdateLesson(ctx, l)
	assert.NoError(t, err)

	mock.ExpectQuery(`SELECT \* FROM course_lessons WHERE id=\$1`).
		WithArgs("l1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow("l1", "L1 Updated"))
	gotL, err := repo.GetLesson(ctx, "l1")
	assert.NoError(t, err)
	assert.Equal(t, "L1 Updated", gotL.Title)

	mock.ExpectExec(`DELETE FROM course_lessons`).
		WithArgs("l1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.DeleteLesson(ctx, "l1")
	assert.NoError(t, err)

	// --- Update & Get & Delete Activity ---
	a.Title = "A1 Updated"
	mock.ExpectExec(`UPDATE course_activities`).
		WithArgs(a.Type, a.Title, a.Order, a.Points, a.IsOptional, a.Content, a.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.UpdateActivity(ctx, a)
	assert.NoError(t, err)

	mock.ExpectQuery(`SELECT \* FROM course_activities WHERE id=\$1`).
		WithArgs("a1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow("a1", "A1 Updated"))
	gotA, err := repo.GetActivity(ctx, "a1")
	assert.NoError(t, err)
	assert.Equal(t, "A1 Updated", gotA.Title)

	mock.ExpectExec(`DELETE FROM course_activities`).
		WithArgs("a1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.DeleteActivity(ctx, "a1")
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

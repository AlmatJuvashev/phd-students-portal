package services_test

import (
	"errors"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestNotifyAdvisorsOnSubmission_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	t.Run("DB Error on Student Name", func(t *testing.T) {
		mock.ExpectQuery("SELECT COALESCE").WithArgs("s1").WillReturnError(errors.New("db error"))
		mock.ExpectQuery("SELECT advisor_id").WithArgs("s1").WillReturnRows(sqlmock.NewRows([]string{"advisor_id"}).AddRow("a1"))
		mock.ExpectExec("INSERT INTO admin_notifications").WillReturnResult(sqlmock.NewResult(1, 1))

		err := services.NotifyAdvisorsOnSubmission(sqlxDB, "s1", "n1", "ni1", "")
		assert.NoError(t, err) // It logs but doesn't fail on name error
	})

	t.Run("DB Error on Advisors List", func(t *testing.T) {
		mock.ExpectQuery("SELECT COALESCE").WithArgs("s1").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Student"))
		mock.ExpectQuery("SELECT advisor_id").WithArgs("s1").WillReturnError(errors.New("db error"))

		err := services.NotifyAdvisorsOnSubmission(sqlxDB, "s1", "n1", "ni1", "")
		assert.Error(t, err)
	})

	t.Run("DB Error on Notification Insert", func(t *testing.T) {
		mock.ExpectQuery("SELECT COALESCE").WithArgs("s1").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Student"))
		mock.ExpectQuery("SELECT advisor_id").WithArgs("s1").WillReturnRows(sqlmock.NewRows([]string{"advisor_id"}).AddRow("a1"))
		mock.ExpectExec("INSERT INTO admin_notifications").WillReturnError(errors.New("insert error"))

		err := services.NotifyAdvisorsOnSubmission(sqlxDB, "s1", "n1", "ni1", "")
		assert.Error(t, err)
	})
}

func TestGetAdvisorsForStudent_Unit(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mock.ExpectQuery("SELECT advisor_id").WithArgs("s1").WillReturnRows(sqlmock.NewRows([]string{"advisor_id"}).AddRow("a1").AddRow("a2"))
	
	ids, err := services.GetAdvisorsForStudent(sqlxDB, "s1")
	assert.NoError(t, err)
	assert.Equal(t, []string{"a1", "a2"}, ids)
}

func TestHasAdvisors_Unit(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mock.ExpectQuery("SELECT COUNT").WithArgs("s1").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	
	has, err := services.HasAdvisors(sqlxDB, "s1")
	assert.NoError(t, err)
	assert.True(t, has)
}

package repository

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSQLUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLUserRepository(sqlxDB)

	mock.ExpectQuery("INSERT INTO users").
		WithArgs("jdoe1234", "john@example.com", "John", "Doe", models.RoleStudent, "hash", true, nil, nil, nil, nil, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("uid-123"))

	id, err := repo.Create(context.Background(), &models.User{
		Username: "jdoe1234",
		Email: "john@example.com",
		FirstName: "John",
		LastName: "Doe",
		Role: models.RoleStudent,
		PasswordHash: "hash",
	})

	assert.NoError(t, err)
	assert.Equal(t, "uid-123", id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSQLUserRepository_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLUserRepository(sqlxDB)

	// Mock Count Query
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM users`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))

	// Mock List Query with JOIN
	// We use QueryMatcherRegexp under the hood, so we need to account for wildcards.
	// The query has newlines, so we need to be careful.
	mock.ExpectQuery(`SELECT u.id, u.username.*FROM users u.*LEFT JOIN LATERAL`).
		WithArgs(10, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "first_name", "last_name", "role", "is_active", "created_at", "phone", "program", "specialty", "department", "cohort"}).
			AddRow("1", "user1", "u1@e.com", "User", "One", "student", true, nil, "", "CS", "AI", "Eng", "2024"))

	users, total, err := repo.List(context.Background(), UserFilter{}, Pagination{Limit: 10, Offset: 0})

	assert.NoError(t, err)
	assert.Equal(t, 10, total)
	if assert.NotEmpty(t, users) {
		assert.Equal(t, "CS", users[0].Program)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

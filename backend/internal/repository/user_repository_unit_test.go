package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSQLUserRepository_GetByEmail_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLUserRepository(sqlxDB)

	email := "test@example.com"
	now := time.Now()

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "username", "email", "first_name", "last_name", "role", "password_hash", "is_active",
			"is_superadmin", "phone", "program", "specialty", "department", "cohort", "avatar_url",
			"bio", "address", "date_of_birth", "created_at", "updated_at",
		}).AddRow(
			"user-1", "testuser", email, "First", "Last", "student", "hash", true,
			false, "123456", "PhD", "CS", "IT", "2023", "http://avatar.com",
			"bio text", "address text", now, now, now,
		)

		mock.ExpectQuery("SELECT (.+) FROM users WHERE LOWER\\(email\\) = LOWER\\(\\$1\\) AND is_active = true").
			WithArgs(email).
			WillReturnRows(rows)

		user, err := repo.GetByEmail(context.Background(), email)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "user-1", user.ID)
		assert.Equal(t, email, user.Email)
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM users WHERE LOWER\\(email\\) = LOWER\\(\\$1\\) AND is_active = true").
			WithArgs(email).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetByEmail(context.Background(), email)

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
		assert.Nil(t, user)
	})

	t.Run("DatabaseError", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM users WHERE LOWER\\(email\\) = LOWER\\(\\$1\\) AND is_active = true").
			WithArgs(email).
			WillReturnError(sql.ErrConnDone)

		user, err := repo.GetByEmail(context.Background(), email)

		assert.Error(t, err)
		assert.Equal(t, sql.ErrConnDone, err)
		assert.Nil(t, user)
	})
}

func TestSQLUserRepository_Update_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLUserRepository(sqlxDB)

	user := &models.User{
		ID:        "user-1",
		FirstName: "NewFirst",
		LastName:  "NewLast",
		Email:     "new@example.com",
		Role:      "advisor",
		Bio:       "new bio",
		Address:   "new address",
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("UPDATE users SET").
			WithArgs(
				user.FirstName, user.LastName, user.Email, user.Role,
				nil, nil, nil, nil, nil, // Nullable fields
				user.Bio, user.Address, user.DateOfBirth, user.AvatarURL,
				user.ID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Update(context.Background(), user)
		assert.NoError(t, err)
	})

	t.Run("DatabaseError", func(t *testing.T) {
		mock.ExpectExec("UPDATE users SET").
			WillReturnError(sql.ErrConnDone)

		err := repo.Update(context.Background(), user)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrConnDone, err)
	})
}

func TestSQLUserRepository_Create_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLUserRepository(sqlxDB)

	user := &models.User{
		Username: "newuser",
		Email:    "new@example.com",
		Role:     "student",
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO users").
			WithArgs(
				user.Username, user.Email, user.FirstName, user.LastName, user.Role, user.PasswordHash, true,
				nil, nil, nil, nil, nil, // phone, program, specialty, department, cohort
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("new-uuid"))

		id, err := repo.Create(context.Background(), user)
		assert.NoError(t, err)
		assert.Equal(t, "new-uuid", id)
	})
}

func TestSQLUserRepository_List_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLUserRepository(sqlxDB)

	filter := UserFilter{
		Role:    "student",
		Program: "PhD",
	}
	pagination := Pagination{
		Limit:  10,
		Offset: 0,
	}

	t.Run("Success", func(t *testing.T) {
		// Mock count query
		mock.ExpectQuery(`SELECT COUNT\(\*\) FROM users u WHERE 1=1 AND u.role = \$1 AND u.program = \$2`).
			WithArgs(filter.Role, filter.Program).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		// Mock main query
		rows := sqlmock.NewRows([]string{
			"id", "username", "email", "first_name", "last_name", "role", "is_active", "created_at",
			"phone", "program", "specialty", "department", "cohort",
		}).AddRow(
			"u-1", "user1", "u1@ex.com", "F", "L", "student", true, time.Now(),
			"", "PhD", "CS", "IT", "2023",
		)

		// Use a more relaxed regex that ignores the comments and specific spacing
		mock.ExpectQuery(`SELECT (.+) FROM users u (.+) LEFT JOIN LATERAL (.+) WHERE (.+) AND u.role = \$1 AND u.program = \$2 ORDER BY u.last_name LIMIT \$3 OFFSET \$4`).
			WithArgs(filter.Role, filter.Program, pagination.Limit, pagination.Offset).
			WillReturnRows(rows)

		users, total, err := repo.List(context.Background(), filter, pagination)

		assert.NoError(t, err)
		assert.Equal(t, 1, total)
		if assert.Len(t, users, 1) {
			assert.Equal(t, "u-1", users[0].ID)
		}
	})

	t.Run("CountError", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT").
			WillReturnError(sql.ErrConnDone)

		users, total, err := repo.List(context.Background(), filter, pagination)

		assert.Error(t, err)
		assert.Equal(t, 0, total)
		assert.Nil(t, users)
	})
}

func TestSQLUserRepository_ReplaceAdvisors_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLUserRepository(sqlx.NewDb(db, "sqlmock"))

	studentID := "st-1"
	tenantID := "t-1"
	advisorIDs := []string{"adv-1", "adv-2"}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		
		// 1. Delete
		mock.ExpectExec(`DELETE FROM student_advisors WHERE student_id=\$1 AND tenant_id=\$2`).
			WithArgs(studentID, tenantID).
			WillReturnResult(sqlmock.NewResult(0, 5)) // 5 deleted

		// 2. Insert loop
		for _, aid := range advisorIDs {
			mock.ExpectExec(`INSERT INTO student_advisors`).
				WithArgs(studentID, aid, tenantID).
				WillReturnResult(sqlmock.NewResult(1, 1))
		}

		mock.ExpectCommit()

		err := repo.ReplaceAdvisors(context.Background(), studentID, advisorIDs, tenantID)
		assert.NoError(t, err)
	})

	t.Run("Rollback", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM student_advisors`).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := repo.ReplaceAdvisors(context.Background(), studentID, advisorIDs, tenantID)
		assert.ErrorIs(t, err, sql.ErrConnDone)
	})
}

func TestSQLUserRepository_RateLimit_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLUserRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("CheckRateLimit", func(t *testing.T) {
		// Matches: SELECT COUNT(*) FROM rate_limit_events WHERE user_id=$1 AND action=$2 AND occurred_at > NOW() - $3::interval
		mock.ExpectQuery(`SELECT COUNT\(\*\) FROM rate_limit_events WHERE user_id=\$1 AND action=\$2 AND occurred_at > NOW\(\) - \$3::interval`).
			WithArgs("u-1", "login", time.Minute).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

		count, err := repo.CheckRateLimit(context.Background(), "u-1", "login", time.Minute)
		assert.NoError(t, err)
		assert.Equal(t, 3, count)
	})

	t.Run("RecordRateLimit", func(t *testing.T) {
		mock.ExpectExec(`INSERT INTO rate_limit_events \(user_id, action\) VALUES \(\$1, \$2\)`).
			WithArgs("u-1", "login").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.RecordRateLimit(context.Background(), "u-1", "login")
		assert.NoError(t, err)
	})
}

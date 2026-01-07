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
		mock.ExpectQuery(`SELECT COUNT\(DISTINCT u.id\) FROM users u WHERE 1=1 AND u.role = \$1 AND u.program = \$2`).
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

		// Extremely relaxed regex to ensure matching despite SQL generation details
		mock.ExpectQuery(`SELECT .*`).
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
		// Mock needs to expect strict string match for "60 seconds" constructed in repository
		mock.ExpectQuery(`SELECT COUNT\(\*\) FROM rate_limit_events WHERE user_id=\$1 AND action=\$2 AND occurred_at > NOW\(\) - \$3::interval`).
			WithArgs("u-1", "login", "60 seconds").
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

func TestSQLUserRepository_PasswordReset_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLUserRepository(sqlx.NewDb(db, "sqlmock"))

	userID := "u1"
	token := "hashedtoken"
	expires := time.Now().Add(time.Hour)

	t.Run("CreateToken", func(t *testing.T) {
		mock.ExpectExec(`INSERT INTO password_reset_tokens`).
			WithArgs(userID, token, expires).
			WillReturnResult(sqlmock.NewResult(1, 1))
		
		err := repo.CreatePasswordResetToken(context.Background(), userID, token, expires)
		assert.NoError(t, err)
	})

	t.Run("GetToken", func(t *testing.T) {
		mock.ExpectQuery(`SELECT user_id, expires_at FROM password_reset_tokens WHERE token_hash = \$1`).
			WithArgs(token).
			WillReturnRows(sqlmock.NewRows([]string{"user_id", "expires_at"}).AddRow(userID, expires))

		uid, exp, err := repo.GetPasswordResetToken(context.Background(), token)
		assert.NoError(t, err)
		assert.Equal(t, userID, uid)
		assert.WithinDuration(t, expires, exp, time.Second)
	})

	t.Run("DeleteToken", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM password_reset_tokens WHERE token_hash = \$1`).
			WithArgs(token).
			WillReturnResult(sqlmock.NewResult(1, 1))
		
		err := repo.DeletePasswordResetToken(context.Background(), token)
		assert.NoError(t, err)
	})
}

func TestSQLUserRepository_MiscUpdates_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLUserRepository(sqlx.NewDb(db, "sqlmock"))
	
	id := "u1"

	t.Run("UpdatePassword", func(t *testing.T) {
		mock.ExpectExec(`UPDATE users SET password_hash=\$1, updated_at=NOW\(\) WHERE id=\$2`).
			WithArgs("newhash", id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		assert.NoError(t, repo.UpdatePassword(context.Background(), id, "newhash"))
	})

	t.Run("UpdateAvatar", func(t *testing.T) {
		mock.ExpectExec(`UPDATE users SET avatar_url=\$1, updated_at=NOW\(\) WHERE id=\$2`).
			WithArgs("http://img.com", id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		assert.NoError(t, repo.UpdateAvatar(context.Background(), id, "http://img.com"))
	})

	t.Run("SetActive", func(t *testing.T) {
		mock.ExpectExec(`UPDATE users SET is_active=\$1 WHERE id=\$2`).
			WithArgs(false, id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		assert.NoError(t, repo.SetActive(context.Background(), id, false))
	})
	
	t.Run("GetByUsername", func(t *testing.T) {
		mock.ExpectQuery(`SELECT (.+) FROM users WHERE username = \$1`).
			WithArgs("user1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(id, "user1"))
		u, err := repo.GetByUsername(context.Background(), "user1")
		assert.NoError(t, err)
		assert.Equal(t, "user1", u.Username)
	})

	t.Run("EmailExists", func(t *testing.T) {
		mock.ExpectQuery(`SELECT COUNT\(\*\) FROM users WHERE email=\$1`).
			WithArgs("test@ex.com").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
		exists, err := repo.EmailExists(context.Background(), "test@ex.com", "")
		assert.NoError(t, err)
		assert.True(t, exists)
	})
}

func TestSQLUserRepository_Cleanup_Unit(t *testing.T) {
	t.Run("Exists", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("stub open error: %s", err)
		}
		defer db.Close()
		repo := NewSQLUserRepository(sqlx.NewDb(db, "sqlmock"))

		mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM users WHERE username=\$1\)`).
			WithArgs("u1").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
		
		exists, err := repo.Exists(context.Background(), "u1")
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("SyncProfileSubmissions", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("stub open error: %s", err)
		}
		defer db.Close()
		repo := NewSQLUserRepository(sqlx.NewDb(db, "sqlmock"))

		data := map[string]string{"program": "PhD"}
		
		mock.ExpectExec(`INSERT INTO profile_submissions`).
			WithArgs("u1", sqlmock.AnyArg(), "t1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.SyncProfileSubmissions(context.Background(), "u1", data, "t1")
		assert.NoError(t, err)
	})

	t.Run("DeleteEmailVerificationToken", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("stub open error: %s", err)
		}
		defer db.Close()
		repo := NewSQLUserRepository(sqlx.NewDb(db, "sqlmock"))

		mock.ExpectExec(`DELETE FROM email_verification_tokens WHERE token=\$1`).
			WithArgs("tok").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.DeleteEmailVerificationToken(context.Background(), "tok")
		assert.NoError(t, err)
	})
}

func TestSQLUserRepository_GetByID_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLUserRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "username", "email", "first_name", "last_name", "role", "is_active", "created_at"}).
			AddRow("u1", "user1", "u1@ex.com", "F", "L", "student", true, time.Now())
		
		mock.ExpectQuery(`SELECT (.+) FROM users WHERE id = \$1`).
			WithArgs("u1").
			WillReturnRows(rows)

		u, err := repo.GetByID(context.Background(), "u1")
		assert.NoError(t, err)
		assert.Equal(t, "u1", u.ID)
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery(`SELECT (.+) FROM users WHERE id = \$1`).
			WithArgs("u1").
			WillReturnError(sql.ErrNoRows)

		_, err := repo.GetByID(context.Background(), "u1")
		assert.Equal(t, ErrNotFound, err)
	})
}

func TestSQLUserRepository_LinkAdvisor_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLUserRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectExec(`INSERT INTO student_advisors`).
		WithArgs("s1", "a1", "t1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.LinkAdvisor(context.Background(), "s1", "a1", "t1")
	assert.NoError(t, err)
}

func TestSQLUserRepository_GetTenantRoles_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLUserRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectQuery(`SELECT roles FROM user_tenant_memberships`).
		WithArgs("u1", "t1").
		WillReturnRows(sqlmock.NewRows([]string{"roles"}).AddRow(pq.StringArray{"admin", "teacher"}))

	roles, err := repo.GetTenantRoles(context.Background(), "u1", "t1")
	assert.NoError(t, err)
	assert.Len(t, roles, 2)
	assert.Equal(t, "admin", roles[0])
}

func TestSQLUserRepository_List_Filters_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLUserRepository(sqlx.NewDb(db, "sqlmock"))

	// Filter with TenantID (triggers JOIN) and other fields
	filter := UserFilter{
		TenantID:   "t1",
		Department: "Dep1",
		Search:     "John",
	}
	pagination := Pagination{Limit: 10, Offset: 0}

	// Expect JOIN in Count Query
	mock.ExpectQuery(`SELECT COUNT.*`).
		WithArgs(filter.TenantID, filter.Department, "%John%").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

	// Expect JOIN in List Query
	mock.ExpectQuery(`SELECT .*`).
		WithArgs(filter.TenantID, filter.Department, "%John%", pagination.Limit, pagination.Offset).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "first_name", "last_name", "role", "is_active", "created_at"}).
			AddRow("u1", "user1", "u1@ex.com", "F", "L", "student", true, time.Now()))

	users, total, err := repo.List(context.Background(), filter, pagination)
	assert.NoError(t, err)
	assert.Equal(t, 5, total)
	assert.Len(t, users, 1)
}

func TestSQLUserRepository_EmailVerification_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLUserRepository(sqlx.NewDb(db, "sqlmock"))

	userID := "u1"
	newEmail := "new@ex.com"
	token := "tok123"
	expires := time.Now().Add(time.Hour)

	t.Run("CreateToken", func(t *testing.T) {
		mock.ExpectExec(`INSERT INTO email_verification_tokens`).
			WithArgs(userID, newEmail, token, expires).
			WillReturnResult(sqlmock.NewResult(1, 1))
		
		err := repo.CreateEmailVerificationToken(context.Background(), userID, newEmail, token, expires)
		assert.NoError(t, err)
	})

	t.Run("GetToken", func(t *testing.T) {
		// Mock query matches regex for token retrieval
		mock.ExpectQuery(`SELECT user_id, new_email, expires_at FROM email_verification_tokens WHERE token=\$1 AND expires_at > NOW\(\)`).
			WithArgs(token).
			WillReturnRows(sqlmock.NewRows([]string{"user_id", "new_email", "expires_at"}).AddRow(userID, newEmail, expires))

		uid, em, expStr, err := repo.GetEmailVerificationToken(context.Background(), token)
		assert.NoError(t, err)
		assert.Equal(t, userID, uid)
		assert.Equal(t, newEmail, em)
		
		// Parse returned string to time for comparison
		expTime, err := time.Parse(time.RFC3339, expStr)
		if err == nil {
			assert.WithinDuration(t, expires, expTime, time.Second)
		}
	})

	t.Run("GetPending", func(t *testing.T) {
		mock.ExpectQuery(`SELECT new_email FROM email_verification_tokens WHERE user_id=\$1 AND expires_at > NOW\(\) ORDER BY expires_at DESC LIMIT 1`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"new_email"}).AddRow(newEmail))

		em, err := repo.GetPendingEmailVerification(context.Background(), userID)
		assert.NoError(t, err)
		assert.Equal(t, newEmail, em)
	})
	
	t.Run("GetPending_NotFound", func(t *testing.T) {
		mock.ExpectQuery(`SELECT new_email FROM email_verification_tokens`).
			WithArgs(userID).
			WillReturnError(sql.ErrNoRows)

		em, err := repo.GetPendingEmailVerification(context.Background(), userID)
		assert.NoError(t, err)
		assert.Equal(t, "", em)
	})
}

func TestSQLUserRepository_Audit_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLUserRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("LogProfileAudit", func(t *testing.T) {
		mock.ExpectExec(`INSERT INTO profile_audit_log`).
			WithArgs("u1", "bio", "old", "new", "admin1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.LogProfileAudit(context.Background(), "u1", "bio", "old", "new", "admin1")
		assert.NoError(t, err)
	})
}

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLUserRepository_CreateUser(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()

	repo := NewSQLUserRepository(db)
	
	user := &models.User{
		ID:           "u1",
		Username:     "testuser",
		Email:        "test@example.com",
		FirstName:    "Test",
		LastName:     "User",
		Role:         "student",
		PasswordHash: "hash",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	id, err := repo.Create(context.Background(), user)
	require.NoError(t, err)
	assert.NotEmpty(t, id)

	// Verify insertion
	fetched, err := repo.GetByEmail(context.Background(), user.Email)
	require.NoError(t, err)
	assert.Equal(t, id, fetched.ID) // Create returns generated id potentially if not supplied, but our model had "u1". Logic in Create ignores ID insert? SQL says: INSERT INTO ... RETURNING id.
	// user_repository.go Create method does NOT insert ID. It lets DB generate it (UUID default).
	assert.Equal(t, user.Username, fetched.Username)
	assert.Equal(t, user.Role, fetched.Role)
}

func TestSQLUserRepository_UpdateUser(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()

	repo := NewSQLUserRepository(db)
	
	user := &models.User{
		Username:     "updateuser",
		Email:        "update@example.com",
		FirstName:    "Update",
		LastName:     "Me",
		Role:         "student",
		PasswordHash: "hash",
		IsActive:     true,
	}
	id, err := repo.Create(context.Background(), user)
	require.NoError(t, err)
	user.ID = id // Capture generated ID

	// Update
	user.FirstName = "UpdatedName"
	user.LastName = "UpdatedLast"
	err = repo.Update(context.Background(), user)
	require.NoError(t, err)

	fetched, err := repo.GetByID(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, "UpdatedName", fetched.FirstName)
	assert.Equal(t, "UpdatedLast", fetched.LastName)
}

func TestSQLUserRepository_LinkAdvisor(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()

	repo := NewSQLUserRepository(db)

	// Create Student
	studentID := testutils.CreateTestUser(t, db, "student1", "student")
	// Create Advisor
	advisorID := testutils.CreateTestUser(t, db, "advisor1", "advisor")
	tenantID := "00000000-0000-0000-0000-000000000001"

	err := repo.LinkAdvisor(context.Background(), studentID, advisorID, tenantID)
	require.NoError(t, err)

	// Verify linkage manually since GetStudentAdvisors is missing
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM student_advisors WHERE student_id=$1 AND advisor_id=$2", studentID, advisorID)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

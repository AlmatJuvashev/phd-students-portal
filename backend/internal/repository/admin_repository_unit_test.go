package repository

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSQLAdminRepository_ListStudentsForMonitor_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLAdminRepository(sqlxDB)

	tenantID := "tenant-1"
	filter := models.FilterParams{
		TenantID: tenantID,
		Program:  "PhD",
	}

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "name", "email", "phone", "program", "department", "cohort", "current_node_id",
		}).AddRow(
			"student-1", "John Doe", "john@example.com", "123", "PhD", "CS", "2023", "node-1",
		)

		// Regex to match the complex base query with subqueries
		mock.ExpectQuery(`SELECT u.id, (.+) FROM users u JOIN user_tenant_memberships utm ON utm.user_id = u.id WHERE u.is_active=true AND u.role='student' AND utm.tenant_id=\$1 AND (.+) ORDER BY u.last_name, u.first_name`).
			WithArgs(tenantID, filter.Program).
			WillReturnRows(rows)

		students, err := repo.ListStudentsForMonitor(context.Background(), filter)

		assert.NoError(t, err)
		assert.Len(t, students, 1)
		assert.Equal(t, "student-1", students[0].ID)
		assert.Equal(t, "John Doe", students[0].Name)
	})

	t.Run("WithAdvisorFilter", func(t *testing.T) {
		advisorID := "advisor-1"
		filterWithAdvisor := models.FilterParams{
			TenantID:  tenantID,
			AdvisorID: advisorID,
		}

		rows := sqlmock.NewRows([]string{
			"id", "name", "email", "phone", "program", "department", "cohort", "current_node_id",
		}).AddRow(
			"student-1", "John Doe", "john@example.com", "123", "PhD", "CS", "2023", "node-1",
		)

		mock.ExpectQuery(`SELECT (.+) JOIN student_advisors sa ON sa.student_id=u.id WHERE (.+) sa.advisor_id=\$2`).
			WithArgs(tenantID, advisorID).
			WillReturnRows(rows)

		students, err := repo.ListStudentsForMonitor(context.Background(), filterWithAdvisor)

		assert.NoError(t, err)
		assert.Len(t, students, 1)
	})
}

func TestSQLAdminRepository_GetAttachmentCounts_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLAdminRepository(sqlxDB)

	instanceID := "instance-1"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT
		COALESCE\(SUM\(CASE WHEN a.status='submitted' THEN 1 ELSE 0 END\),0\) AS submitted,
		COALESCE\(SUM\(CASE WHEN a.status IN \('approved', 'approved_with_comments'\) THEN 1 ELSE 0 END\),0\) AS approved,
		COALESCE\(SUM\(CASE WHEN a.status='rejected' THEN 1 ELSE 0 END\),0\) AS rejected
		FROM node_instance_slot_attachments a
		JOIN node_instance_slots s ON s.id=a.slot_id
		WHERE s.node_instance_id=\$1 AND a.is_active`).
			WithArgs(instanceID).
			WillReturnRows(sqlmock.NewRows([]string{"submitted", "approved", "rejected"}).AddRow(5, 3, 1))

		submitted, approved, rejected, err := repo.GetAttachmentCounts(context.Background(), instanceID)

		assert.NoError(t, err)
		assert.Equal(t, 5, submitted)
		assert.Equal(t, 3, approved)
		assert.Equal(t, 1, rejected)
	})
}

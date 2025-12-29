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

func TestSQLAdminRepository_GetStudentDetails_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLAdminRepository(sqlx.NewDb(db, "sqlmock"))

	studentID := "st-1"
	tenantID := "t-1"

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "email", "phone", "first_name", "last_name", "program", "department", "cohort",
		}).AddRow(
			studentID, "jane@example.com", "555-0100", "Jane", "Doe", "PhD", "Biology", "2024",
		)

		// Regex matching the actual query in repository
		mock.ExpectQuery(`SELECT u.id, COALESCE\(u.email,('')\) AS email, COALESCE\(ps.form_data->>'phone',('')\) AS phone, u.first_name, u.last_name, COALESCE\(ps.form_data->>'program',('')\) AS program, COALESCE\(ps.form_data->>'department',('')\) AS department, COALESCE\(ps.form_data->>'cohort',('')\) AS cohort FROM users u LEFT JOIN profile_submissions ps ON ps.user_id = u.id JOIN user_tenant_memberships utm ON utm.user_id = u.id WHERE u.id=\$1 AND u.role='student' AND utm.tenant_id=\$2`).
			WithArgs(studentID, tenantID).
			WillReturnRows(rows)

		details, err := repo.GetStudentDetails(context.Background(), studentID, tenantID)
		assert.NoError(t, err)
		assert.Equal(t, "Jane", details.FirstName)
		assert.Equal(t, "jane@example.com", details.Email)
	})
}

func TestSQLAdminRepository_Attachments_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLAdminRepository(sqlx.NewDb(db, "sqlmock"))

	attID := "att-1"
	status := "approved"
	note := "Good job"
	actorID := "admin-1"

	t.Run("UpdateStatus", func(t *testing.T) {
		mock.ExpectExec(`UPDATE node_instance_slot_attachments SET status=\$1, review_note=\$2, approved_by=\$3, approved_at=now\(\) WHERE id=\$4`).
			WithArgs(status, note, actorID, attID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateAttachmentStatus(context.Background(), attID, status, note, actorID)
		assert.NoError(t, err)
	})
}

func TestSQLAdminRepository_GetStudentNodeInstances_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLAdminRepository(sqlx.NewDb(db, "sqlmock"))

	studentID := "st-1"
	t.Run("Success", func(t *testing.T) {
		// sqlmock Rows must match columns selected in the query
		// Query: DISTINCT ON ... id, tenant_id, user_id, playbook_version_id, node_id, state, opened_at, submitted_at, updated_at
		rows := sqlmock.NewRows([]string{"id", "tenant_id", "user_id", "playbook_version_id", "node_id", "state", "opened_at", "submitted_at", "updated_at"}).
			AddRow("inst-1", "t-1", studentID, "v-1", "node-1", "in_progress", time.Now(), nil, time.Now())

		mock.ExpectQuery(`SELECT DISTINCT ON \(node_id\) id, tenant_id, user_id, playbook_version_id, node_id, state, opened_at, submitted_at, updated_at FROM node_instances WHERE user_id=\$1 ORDER BY node_id, updated_at DESC`).
			WithArgs(studentID).
			WillReturnRows(rows)

		insts, err := repo.GetStudentNodeInstances(context.Background(), studentID)
		assert.NoError(t, err)
		if assert.Len(t, insts, 1) {
			assert.Equal(t, "node-1", insts[0].NodeID)
		}
	})
}

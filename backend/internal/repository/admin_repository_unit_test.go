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

func TestSQLAdminRepository_GetAntiplagCount_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLAdminRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT COUNT\(\*\) FROM node_instances WHERE playbook_version_id=\? AND node_id='S1_antiplag' AND state='done' AND user_id IN \(\?, \?\)`).
			WithArgs("v1", "s1", "s2").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

		count, err := repo.GetAntiplagCount(context.Background(), []string{"s1", "s2"}, "v1")
		assert.NoError(t, err)
		assert.Equal(t, 5, count)
	})
}

func TestSQLAdminRepository_GetW2Durations_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLAdminRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("Success", func(t *testing.T) {
		start := time.Now().Add(-48 * time.Hour)
		end := time.Now()
		
		rows := sqlmock.NewRows([]string{"user_id", "start", "end"}).
			AddRow("u1", start, end) // 2 days

		mock.ExpectQuery(`SELECT user_id, MIN\(updated_at\) as start, MAX\(updated_at\) as end FROM node_instances WHERE playbook_version_id=\? AND node_id IN \(\?\) AND user_id IN \(\?\) GROUP BY user_id`).
			WithArgs("v1", "n1", "u1").
			WillReturnRows(rows)

		durations, err := repo.GetW2Durations(context.Background(), []string{"u1"}, "v1", []string{"n1"})
		assert.NoError(t, err)
		assert.Len(t, durations, 1)
		assert.InDelta(t, 2.0, durations[0], 0.1)
	})
}

func TestSQLAdminRepository_GetBottleneck_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLAdminRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("Success", func(t *testing.T) {
		since := time.Now().Add(-24 * time.Hour)
		rows := sqlmock.NewRows([]string{"node_id", "cnt"}).AddRow("node-1", 10)

		mock.ExpectQuery(`SELECT node_id, COUNT\(\*\) as cnt 
		FROM node_instances 
		WHERE playbook_version_id=\? AND user_id IN \(\?\) AND state IN \('waiting','needs_fixes'\) AND updated_at >= \? 
		GROUP BY node_id ORDER BY cnt DESC LIMIT 1`).
			WithArgs("v1", "u1", since).
			WillReturnRows(rows)

		nodeID, count, err := repo.GetBottleneck(context.Background(), []string{"u1"}, "v1", since)
		assert.NoError(t, err)
		assert.Equal(t, "node-1", nodeID)
		assert.Equal(t, 10, count)
	})
}

func TestSQLAdminRepository_GetNodeFiles_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLAdminRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("Success", func(t *testing.T) {
		// 1. Get Instance ID Query
		mock.ExpectQuery(`SELECT id FROM node_instances WHERE user_id=\$1 AND node_id=\$2 ORDER BY updated_at DESC LIMIT 1`).
			WithArgs("u1", "n1").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("inst-1"))

		// 2. Main Files Query matches db tags in models.NodeFile
		cols := []string{
			"slot_key", "attachment_id", "filename", "size_bytes", "status", "review_note", "is_active",
			"attached_at", "approved_at", "approved_by", "version_id", "mime_type",
			"uploaded_by", "reviewed_doc_id", "reviewed_at", "reviewed_mime_type", "reviewed_by_name",
		}
		
		rows := sqlmock.NewRows(cols).
			AddRow(
				"k1", "att-1", "f1.pdf", 100, "submitted", "", true,
				"2023-01-01T00:00:00Z", nil, nil, "v1", "application/pdf",
				"John Doe", "rv1", "2023-01-02T00:00:00Z", "application/pdf", "Reviewer Name",
			)

		// Use the exact query structure with broad regex to avoid whitespace issues
		queryRegex := `SELECT s\.slot_key, a\.id as attachment_id, a\.filename, a\.size_bytes, a\.status, a\.review_note, a\.is_active,.*` +
			`to_char\(a\.attached_at, 'YYYY-MM-DD"T"HH24:MI:SSZ'\) as attached_at,.*` +
			`to_char\(a\.approved_at, 'YYYY-MM-DD"T"HH24:MI:SSZ'\) as approved_at,.*` +
			`a\.approved_by, dv\.id AS version_id, dv\.mime_type,.*` +
			`COALESCE\(u\.first_name\|\|' '\|\|u\.last_name,''\) AS uploaded_by,.*` +
			`a\.reviewed_document_version_id as reviewed_doc_id,.*` +
			`to_char\(a\.reviewed_at, 'YYYY-MM-DD"T"HH24:MI:SSZ'\) as reviewed_at,.*` +
			`rdv\.mime_type AS reviewed_mime_type,.*` +
			`COALESCE\(ru\.first_name\|\|' '\|\|ru\.last_name,''\) AS reviewed_by_name.*` +
			`FROM node_instance_slots s.*` +
			`JOIN node_instance_slot_attachments a ON a\.slot_id=s\.id.*` +
			`JOIN document_versions dv ON dv\.id=a\.document_version_id.*` +
			`LEFT JOIN users u ON u\.id=a\.attached_by.*` +
			`LEFT JOIN document_versions rdv ON rdv\.id=a\.reviewed_document_version_id.*` +
			`LEFT JOIN users ru ON ru\.id=a\.reviewed_by.*` +
			`WHERE s\.node_instance_id=\$1.*` +
			`ORDER BY a\.attached_at ASC`

		mock.ExpectQuery(queryRegex).
			WithArgs("inst-1").
			WillReturnRows(rows)

		files, err := repo.GetNodeFiles(context.Background(), "u1", "n1")
		assert.NoError(t, err)
		if assert.Len(t, files, 1) {
			f := files[0]
			assert.Equal(t, "f1.pdf", f.Filename)
			assert.Equal(t, "John Doe", f.UploadedBy)
			
			// Verify post-processing
			assert.NotNil(t, f.ReviewedDocument)
			assert.Equal(t, "rv1", f.ReviewedDocument.VersionID)
			assert.Contains(t, f.ReviewedDocument.DownloadURL, "rv1")
			assert.Equal(t, "Reviewer Name", f.ReviewedDocument.ReviewedBy)
		}
	})
}

func TestSQLAdminRepository_GetStudentJourneyNodes_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLAdminRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("Success", func(t *testing.T) {
		// 1. Nodes Query
		mock.ExpectQuery(`SELECT DISTINCT ON \(node_id\) id, node_id, state, to_char\(updated_at, 'YYYY-MM-DD"T"HH24:MI:SSZ'\) as updated_at FROM node_instances WHERE user_id=\$1 ORDER BY node_id, updated_at DESC`).
			WithArgs("u1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "node_id", "state", "updated_at"}).
				AddRow("inst-1", "n1", "done", "2023-01-01T00:00:00Z"))

		// 2. Attachments Query
		// Matches: SELECT s.node_instance_id, a.filename, a.size_bytes, to_char(a.attached_at, ...), dv.id as version_id
		// FROM node_instance_slot_attachments a JOIN ... WHERE ... IN (?) ...
		mock.ExpectQuery(`SELECT s.node_instance_id, a.filename, a.size_bytes, (.+) FROM node_instance_slot_attachments a (.+) WHERE s.node_instance_id IN \(\?\) AND a.is_active=true`).
			WithArgs("inst-1").
			WillReturnRows(sqlmock.NewRows([]string{"node_instance_id", "filename", "size_bytes", "attached_at", "version_id"}).
				AddRow("inst-1", "file1.pdf", 1024, "2023-01-01T00:00:00Z", "v1"))

		nodes, err := repo.GetStudentJourneyNodes(context.Background(), "u1")
		assert.NoError(t, err)
		if assert.Len(t, nodes, 1) {
			assert.Equal(t, "n1", nodes[0].NodeID)
			assert.Equal(t, 1, nodes[0].Attachments)
			assert.Len(t, nodes[0].Files, 1)
			assert.Equal(t, "file1.pdf", nodes[0].Files[0].Filename)
		}
	})
}

func TestSQLAdminRepository_GetAdvisorsForStudents_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLAdminRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("Success", func(t *testing.T) {
		studentIDs := []string{"s1", "s2"}
		rows := sqlmock.NewRows([]string{"student_id", "id", "name", "email"}).
			AddRow("s1", "a1", "Advisor 1", "a1@ex.com").
			AddRow("s1", "a2", "Advisor 2", "a2@ex.com").
			AddRow("s2", "a1", "Advisor 1", "a1@ex.com")

		mock.ExpectQuery(`SELECT (.+) FROM student_advisors sa JOIN users u ON u.id=sa.advisor_id WHERE sa.student_id IN \(\?, \?\)`).
			WithArgs("s1", "s2").
			WillReturnRows(rows)

		advisors, err := repo.GetAdvisorsForStudents(context.Background(), studentIDs)
		assert.NoError(t, err)
		assert.Len(t, advisors, 2)
		assert.Len(t, advisors["s1"], 2)
		assert.Len(t, advisors["s2"], 1)
	})
}

func TestSQLAdminRepository_GetDoneCountsForStudents_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLAdminRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("Success", func(t *testing.T) {
		studentIDs := []string{"s1", "s2"}
		rows := sqlmock.NewRows([]string{"user_id", "count"}).
			AddRow("s1", 5).
			AddRow("s2", 3)

		mock.ExpectQuery(`SELECT user_id, COUNT\(\*\) FROM node_instances WHERE state='done' AND user_id IN \(\?, \?\) GROUP BY user_id`).
			WithArgs("s1", "s2").
			WillReturnRows(rows)

		counts, err := repo.GetDoneCountsForStudents(context.Background(), studentIDs)
		assert.NoError(t, err)
		assert.Equal(t, 5, counts["s1"])
		assert.Equal(t, 3, counts["s2"])
	})
}

func TestSQLAdminRepository_GetLastUpdatesForStudents_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLAdminRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("Success", func(t *testing.T) {
		studentIDs := []string{"s1"}
		now := time.Now()

		// 1. Node Instances
		mock.ExpectQuery(`SELECT user_id, MAX\(updated_at\) FROM node_instances WHERE user_id IN \(\?\) GROUP BY user_id`).
			WithArgs("s1").
			WillReturnRows(sqlmock.NewRows([]string{"user_id", "max"}).AddRow("s1", now.Add(-time.Hour)))

		// 2. Revisions
		mock.ExpectQuery(`SELECT ni.user_id, MAX\(r.created_at\) FROM node_instance_form_revisions r JOIN node_instances ni ON ni.id=r.node_instance_id WHERE ni.user_id IN \(\?\) GROUP BY ni.user_id`).
			WithArgs("s1").
			WillReturnRows(sqlmock.NewRows([]string{"user_id", "max"}).AddRow("s1", now))

		// 3. Attachments
		mock.ExpectQuery(`SELECT ni.user_id, MAX\(a.attached_at\) FROM node_instance_slot_attachments a JOIN node_instance_slots s ON s.id=a.slot_id JOIN node_instances ni ON ni.id=s.node_instance_id WHERE ni.user_id IN \(\?\) GROUP BY ni.user_id`).
			WithArgs("s1").
			WillReturnRows(sqlmock.NewRows([]string{"user_id", "max"}).AddRow("s1", now.Add(-2*time.Hour)))

		// 4. Events
		mock.ExpectQuery(`SELECT ni.user_id, MAX\(e.created_at\) FROM node_events e JOIN node_instances ni ON ni.id=e.node_instance_id WHERE ni.user_id IN \(\?\) GROUP BY ni.user_id`).
			WithArgs("s1").
			WillReturnRows(sqlmock.NewRows([]string{"user_id", "max"}).AddRow("s1", now.Add(-3*time.Hour)))

		updates, err := repo.GetLastUpdatesForStudents(context.Background(), studentIDs)
		assert.NoError(t, err)
		assert.True(t, updates["s1"].Equal(now))
	})
}

func TestSQLAdminRepository_GetRPRequiredForStudents_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLAdminRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("Success", func(t *testing.T) {
		studentIDs := []string{"s1", "s2"}
		rows := sqlmock.NewRows([]string{"user_id", "form_data"}).
			AddRow("s1", []byte(`{"years_since_graduation": 4}`)).
			AddRow("s2", []byte(`{"years_since_graduation": 1}`))

		mock.ExpectQuery(`SELECT user_id, form_data FROM profile_submissions WHERE user_id IN \(\?, \?\)`).
			WithArgs("s1", "s2").
			WillReturnRows(rows)

		required, err := repo.GetRPRequiredForStudents(context.Background(), studentIDs)
		assert.NoError(t, err)
		assert.True(t, required["s1"])
		assert.False(t, required["s2"])
	})
}

func TestSQLAdminRepository_GetAdminUnreadCount_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLAdminRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT COUNT\(\*\) FROM admin_notifications WHERE is_read = false`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

		count, err := repo.GetAdminUnreadCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 3, count)
	})
}

func TestSQLAdminRepository_NodeState_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLAdminRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("UpdateNodeInstanceState", func(t *testing.T) {
		mock.ExpectExec("UPDATE node_instances SET state=\\$1, submitted_at=COALESCE\\(submitted_at, now\\(\\)\\), updated_at=now\\(\\) WHERE id=\\$2").
			WithArgs("submitted", "inst-1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateNodeInstanceState(context.Background(), "inst-1", "submitted")
		assert.NoError(t, err)
	})

	t.Run("UpdateAllNodeInstances", func(t *testing.T) {
		mock.ExpectExec("UPDATE node_instances SET state=\\$1, updated_at=now\\(\\) WHERE user_id=\\$2 AND node_id=\\$3 AND id != \\$4").
			WithArgs("archived", "s1", "n1", "inst-1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateAllNodeInstances(context.Background(), "s1", "n1", "inst-1", "archived")
		assert.NoError(t, err)
	})

	t.Run("UpsertJourneyState", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO journey_states").
			WithArgs("t1", "s1", "n1", "done").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpsertJourneyState(context.Background(), "t1", "s1", "n1", "done")
		assert.NoError(t, err)
	})

	t.Run("LogNodeEvent", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO node_events").
			WithArgs("inst-1", "submit", sqlmock.AnyArg(), "actor-1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.LogNodeEvent(context.Background(), "inst-1", "submit", "actor-1", map[string]any{"foo": "bar"})
		assert.NoError(t, err)
	})
}

func TestSQLAdminRepository_CreateReminders_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLAdminRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("Success", func(t *testing.T) {
		studentIDs := []string{"s1", "s2"}
		due := "2023-12-31T23:59:59Z"

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO reminders").WithArgs("s1", "Title", "Msg", due, "admin-1").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO reminders").WithArgs("s2", "Title", "Msg", due, "admin-1").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.CreateReminders(context.Background(), studentIDs, "Title", "Msg", &due, "admin-1")
		assert.NoError(t, err)
	})
}

func TestSQLAdminRepository_AttachmentGaps_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLAdminRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("GetAttachmentMeta", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"instance_id", "slot_id", "slot_key", "node_id", "student_id", "state", "tenant_id", "filename", "status", "document_id"}).
			AddRow("inst-1", "slot-1", "key-1", "node-1", "s1", "in_progress", "t1", "file.pdf", "submitted", "doc-1")

		mock.ExpectQuery("SELECT (.+) FROM node_instance_slot_attachments a").
			WithArgs("att-1").
			WillReturnRows(rows)

		meta, err := repo.GetAttachmentMeta(context.Background(), "att-1")
		assert.NoError(t, err)
		assert.Equal(t, "inst-1", meta.InstanceID)
	})

	t.Run("UploadReviewedDocument", func(t *testing.T) {
		mock.ExpectExec("UPDATE node_instance_slot_attachments").
			WithArgs("v1", "admin-1", "att-1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UploadReviewedDocument(context.Background(), "att-1", "v1", "admin-1")
		assert.NoError(t, err)
	})

	t.Run("GetLatestAttachmentStatus", func(t *testing.T) {
		mock.ExpectQuery("SELECT a.status FROM node_instance_slot_attachments a").
			WithArgs("inst-1").
			WillReturnRows(sqlmock.NewRows([]string{"status"}).AddRow("approved"))

		status, err := repo.GetLatestAttachmentStatus(context.Background(), "inst-1")
		assert.NoError(t, err)
		assert.Equal(t, "approved", status)
	})
}

func TestSQLAdminRepository_Notifications_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLAdminRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("CreateNotification", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO notifications").
			WithArgs("r1", "Title", "Msg", "/link", "info", "t1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.CreateNotification(context.Background(), "r1", "Title", "Msg", "/link", "info", "t1")
		assert.NoError(t, err)
	})

	t.Run("MarkAllAdminNotificationsRead", func(t *testing.T) {
		mock.ExpectExec("UPDATE admin_notifications SET is_read = true").
			WillReturnResult(sqlmock.NewResult(1, 4))

		err := repo.MarkAllAdminNotificationsRead(context.Background())
		assert.NoError(t, err)
	})
}


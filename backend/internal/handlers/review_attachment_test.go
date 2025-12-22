package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	pb "github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// reviewTestFixtures holds all the IDs needed for attachment review tests
type reviewTestFixtures struct {
	TenantID     string
	StudentID    string
	AdvisorID    string
	AdminID      string
	PVVersionID  string
	NodeInstID   string
	SlotID       string
	DocID        string
	DocVersionID string
	AttachmentID string
}

// setupReviewFixtures creates complete test data for attachment review tests
// Includes: tenant, student, advisor, admin, playbook version, node instance, slot, document, attachment
func setupReviewFixtures(t *testing.T, db interface {
	Exec(string, ...interface{}) (interface{}, error)
	QueryRow(string, ...interface{}) interface{ Scan(...interface{}) error }
}) reviewTestFixtures {
	testDB, teardown := testutils.SetupTestDB()
	t.Cleanup(teardown)

	f := reviewTestFixtures{
		TenantID:    uuid.New().String(),
		StudentID:   uuid.New().String(),
		AdvisorID:   uuid.New().String(),
		AdminID:     uuid.New().String(),
		PVVersionID: uuid.New().String(),
		NodeInstID:  uuid.New().String(),
		SlotID:      uuid.New().String(),
		DocID:       uuid.New().String(),
	}

	tenantSlug := "revtest-" + f.TenantID[:8]

	// Create tenant
	_, err := testDB.Exec(`INSERT INTO tenants (id, slug, name, tenant_type, is_active) 
		VALUES ($1, $2, 'Test', 'university', true)
		ON CONFLICT (id) DO NOTHING`, f.TenantID, tenantSlug)
	require.NoError(t, err)

	// Create users  
	_, err = testDB.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) VALUES 
		($1, 'teststudent', 'student@t.com', 'Test', 'Student', 'student', 'hash', true),
		($2, 'testadvisor', 'advisor@t.com', 'Test', 'Advisor', 'advisor', 'hash', true),
		($3, 'testadmin', 'admin@t.com', 'Test', 'Admin', 'admin', 'hash', true)
		ON CONFLICT (id) DO NOTHING`, f.StudentID, f.AdvisorID, f.AdminID)
	require.NoError(t, err)

	// Assign advisor to student
	_, err = testDB.Exec(`INSERT INTO student_advisors (student_id, advisor_id, tenant_id) VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING`, f.StudentID, f.AdvisorID, f.TenantID)
	require.NoError(t, err)

	// Create playbook version
	_, err = testDB.Exec(`INSERT INTO playbook_versions (id, tenant_id, version, checksum, raw_json) 
		VALUES ($1, $2, 'v1', 'c', '{}')
		ON CONFLICT (id) DO NOTHING`, f.PVVersionID, f.TenantID)
	require.NoError(t, err)

	// Create node instance
	_, err = testDB.Exec(`INSERT INTO node_instances (id, tenant_id, node_id, user_id, state, playbook_version_id) 
		VALUES ($1, $2, 'confirm_task', $3, 'submitted', $4)
		ON CONFLICT (id) DO NOTHING`, f.NodeInstID, f.TenantID, f.StudentID, f.PVVersionID)
	require.NoError(t, err)

	// Create slot
	_, err = testDB.Exec(`INSERT INTO node_instance_slots (id, tenant_id, node_instance_id, slot_key) 
		VALUES ($1, $2, $3, 'doc')
		ON CONFLICT (id) DO NOTHING`, f.SlotID, f.TenantID, f.NodeInstID)
	require.NoError(t, err)

	// Create document
	_, err = testDB.Exec(`INSERT INTO documents (id, tenant_id, user_id, title, kind) 
		VALUES ($1, $2, $3, 'Test Doc', 'other')
		ON CONFLICT (id) DO NOTHING`, f.DocID, f.TenantID, f.StudentID)
	require.NoError(t, err)

	// Create document version
	err = testDB.QueryRow(`INSERT INTO document_versions (tenant_id, document_id, storage_path, mime_type, size_bytes, uploaded_by) 
		VALUES ($1, $2, 'path/to/doc', 'application/pdf', 1024, $3) RETURNING id`, f.TenantID, f.DocID, f.StudentID).Scan(&f.DocVersionID)
	require.NoError(t, err)

	// Create attachment
	err = testDB.QueryRow(`INSERT INTO node_instance_slot_attachments (slot_id, document_version_id, is_active, status, filename, size_bytes, attached_by) 
		VALUES ($1, $2, true, 'submitted', 'test.pdf', 1024, $3) RETURNING id`, f.SlotID, f.DocVersionID, f.StudentID).Scan(&f.AttachmentID)
	require.NoError(t, err)

	return f
}

func TestReviewAttachment_AdvisorApproves(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Setup fixtures
	f := reviewTestFixtures{
		TenantID:    uuid.New().String(),
		StudentID:   uuid.New().String(),
		AdvisorID:   uuid.New().String(),
		PVVersionID: uuid.New().String(),
		NodeInstID:  uuid.New().String(),
		SlotID:      uuid.New().String(),
		DocID:       uuid.New().String(),
	}

	tenantSlug := "adv-approve-" + f.TenantID[:8]
	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type, is_active) VALUES ($1, $2, 'T', 'university', true)`, f.TenantID, tenantSlug)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) VALUES 
		($1, 'stu', 's@t.com', 'S', 'T', 'student', 'h', true),
		($2, 'adv', 'a@t.com', 'A', 'D', 'advisor', 'h', true)`, f.StudentID, f.AdvisorID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO student_advisors (student_id, advisor_id, tenant_id) VALUES ($1, $2, $3)`, f.StudentID, f.AdvisorID, f.TenantID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO playbook_versions (id, tenant_id, version, checksum, raw_json) VALUES ($1, $2, 'v1', 'c', '{}')`, f.PVVersionID, f.TenantID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO node_instances (id, tenant_id, node_id, user_id, state, playbook_version_id) VALUES ($1, $2, 'confirm_task', $3, 'submitted', $4)`, f.NodeInstID, f.TenantID, f.StudentID, f.PVVersionID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO node_instance_slots (id, tenant_id, node_instance_id, slot_key) VALUES ($1, $2, $3, 'doc')`, f.SlotID, f.TenantID, f.NodeInstID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO documents (id, tenant_id, user_id, title, kind) VALUES ($1, $2, $3, 'Test Doc', 'other')`, f.DocID, f.TenantID, f.StudentID)
	require.NoError(t, err)

	err = db.QueryRow(`INSERT INTO document_versions (tenant_id, document_id, storage_path, mime_type, size_bytes, uploaded_by) VALUES ($1, $2, 'path', 'application/pdf', 100, $3) RETURNING id`, f.TenantID, f.DocID, f.StudentID).Scan(&f.DocVersionID)
	require.NoError(t, err)

	err = db.QueryRow(`INSERT INTO node_instance_slot_attachments (slot_id, document_version_id, is_active, status, filename, size_bytes, attached_by) VALUES ($1, $2, true, 'submitted', 'test.pdf', 100, $3) RETURNING id`, f.SlotID, f.DocVersionID, f.StudentID).Scan(&f.AttachmentID)
	require.NoError(t, err)

	// Setup handler
	pbm := &pb.Manager{VersionID: f.PVVersionID, Nodes: map[string]pb.Node{"confirm_task": {}}}
	repo := repository.NewSQLAdminRepository(db)
	svc := services.NewAdminService(repo, pbm, config.AppConfig{})
	jRepo := repository.NewSQLJourneyRepository(db)
	jSvc := services.NewJourneyService(jRepo, pbm, config.AppConfig{}, nil, nil, nil)
	h := handlers.NewAdminHandler(config.AppConfig{}, pbm, svc, jSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": f.AdvisorID, "role": "advisor"})
		c.Next()
	})
	r.PATCH("/admin/attachments/:attachmentId/review", h.ReviewAttachment)

	// Test: Advisor approves
	body := map[string]string{"status": "approved", "note": "Good work!"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("PATCH", "/admin/attachments/"+f.AttachmentID+"/review", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Logf("Response: %s", w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify node state is now "done"
	var state string
	err = db.Get(&state, `SELECT state FROM node_instances WHERE id=$1`, f.NodeInstID)
	require.NoError(t, err)
	assert.Equal(t, "done", state)
}

func TestReviewAttachment_AdvisorApprovesWithComments(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	f := reviewTestFixtures{
		TenantID:    uuid.New().String(),
		StudentID:   uuid.New().String(),
		AdvisorID:   uuid.New().String(),
		PVVersionID: uuid.New().String(),
		NodeInstID:  uuid.New().String(),
		SlotID:      uuid.New().String(),
		DocID:       uuid.New().String(),
	}

	tenantSlug := "awc-" + f.TenantID[:8]
	db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type, is_active) VALUES ($1, $2, 'T', 'university', true)`, f.TenantID, tenantSlug)
	db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) VALUES ($1, 'stu', 's@t.com', 'S', 'T', 'student', 'h', true), ($2, 'adv', 'a@t.com', 'A', 'D', 'advisor', 'h', true)`, f.StudentID, f.AdvisorID)
	db.Exec(`INSERT INTO student_advisors (student_id, advisor_id, tenant_id) VALUES ($1, $2, $3)`, f.StudentID, f.AdvisorID, f.TenantID)
	db.Exec(`INSERT INTO playbook_versions (id, tenant_id, version, checksum, raw_json) VALUES ($1, $2, 'v1', 'c', '{}')`, f.PVVersionID, f.TenantID)
	db.Exec(`INSERT INTO node_instances (id, tenant_id, node_id, user_id, state, playbook_version_id) VALUES ($1, $2, 'n', $3, 'submitted', $4)`, f.NodeInstID, f.TenantID, f.StudentID, f.PVVersionID)
	db.Exec(`INSERT INTO node_instance_slots (id, tenant_id, node_instance_id, slot_key) VALUES ($1, $2, $3, 'doc')`, f.SlotID, f.TenantID, f.NodeInstID)
	db.Exec(`INSERT INTO documents (id, tenant_id, user_id, title, kind) VALUES ($1, $2, $3, 'Test', 'other')`, f.DocID, f.TenantID, f.StudentID)
	db.QueryRow(`INSERT INTO document_versions (tenant_id, document_id, storage_path, mime_type, size_bytes, uploaded_by) VALUES ($1, $2, 'p', 'pdf', 100, $3) RETURNING id`, f.TenantID, f.DocID, f.StudentID).Scan(&f.DocVersionID)
	db.QueryRow(`INSERT INTO node_instance_slot_attachments (slot_id, document_version_id, is_active, status, filename, size_bytes, attached_by) VALUES ($1, $2, true, 'submitted', 't.pdf', 100, $3) RETURNING id`, f.SlotID, f.DocVersionID, f.StudentID).Scan(&f.AttachmentID)

	pbm := &pb.Manager{VersionID: f.PVVersionID, Nodes: map[string]pb.Node{"n": {}}}
	repo := repository.NewSQLAdminRepository(db)
	svc := services.NewAdminService(repo, pbm, config.AppConfig{})
	jRepo := repository.NewSQLJourneyRepository(db)
	jSvc := services.NewJourneyService(jRepo, pbm, config.AppConfig{}, nil, nil, nil)
	h := handlers.NewAdminHandler(config.AppConfig{}, pbm, svc, jSvc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": f.AdvisorID, "role": "advisor"})
		c.Next()
	})
	r.PATCH("/admin/attachments/:attachmentId/review", h.ReviewAttachment)

	body := map[string]string{"status": "approved_with_comments", "note": "Minor typos but OK"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("PATCH", "/admin/attachments/"+f.AttachmentID+"/review", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Logf("Response: %s", w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var state string
	db.Get(&state, `SELECT state FROM node_instances WHERE id=$1`, f.NodeInstID)
	assert.Equal(t, "done", state, "approved_with_comments should set node to done")
}

func TestReviewAttachment_AdvisorUnassignedStudent_Forbidden(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	f := reviewTestFixtures{
		TenantID:    uuid.New().String(),
		StudentID:   uuid.New().String(),
		AdvisorID:   uuid.New().String(),         // Unassigned advisor
		PVVersionID: uuid.New().String(),
		NodeInstID:  uuid.New().String(),
		SlotID:      uuid.New().String(),
		DocID:       uuid.New().String(),
	}
	otherAdvisorID := uuid.New().String()

	tenantSlug := "unassigned-" + f.TenantID[:8]
	db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type, is_active) VALUES ($1, $2, 'T', 'university', true)`, f.TenantID, tenantSlug)
	db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) VALUES ($1, 'stu', 's@t.com', 'S', 'T', 'student', 'h', true), ($2, 'adv1', 'a1@t.com', 'A1', 'D', 'advisor', 'h', true), ($3, 'adv2', 'a2@t.com', 'A2', 'D', 'advisor', 'h', true)`, f.StudentID, f.AdvisorID, otherAdvisorID)
	// Only assign otherAdvisorID, NOT f.AdvisorID
	db.Exec(`INSERT INTO student_advisors (student_id, advisor_id, tenant_id) VALUES ($1, $2, $3)`, f.StudentID, otherAdvisorID, f.TenantID)
	db.Exec(`INSERT INTO playbook_versions (id, tenant_id, version, checksum, raw_json) VALUES ($1, $2, 'v1', 'c', '{}')`, f.PVVersionID, f.TenantID)
	db.Exec(`INSERT INTO node_instances (id, tenant_id, node_id, user_id, state, playbook_version_id) VALUES ($1, $2, 'n', $3, 'submitted', $4)`, f.NodeInstID, f.TenantID, f.StudentID, f.PVVersionID)
	db.Exec(`INSERT INTO node_instance_slots (id, tenant_id, node_instance_id, slot_key) VALUES ($1, $2, $3, 'doc')`, f.SlotID, f.TenantID, f.NodeInstID)
	db.Exec(`INSERT INTO documents (id, tenant_id, user_id, title, kind) VALUES ($1, $2, $3, 'Test', 'other')`, f.DocID, f.TenantID, f.StudentID)
	db.QueryRow(`INSERT INTO document_versions (tenant_id, document_id, storage_path, mime_type, size_bytes, uploaded_by) VALUES ($1, $2, 'p', 'pdf', 100, $3) RETURNING id`, f.TenantID, f.DocID, f.StudentID).Scan(&f.DocVersionID)
	db.QueryRow(`INSERT INTO node_instance_slot_attachments (slot_id, document_version_id, is_active, status, filename, size_bytes, attached_by) VALUES ($1, $2, true, 'submitted', 't.pdf', 100, $3) RETURNING id`, f.SlotID, f.DocVersionID, f.StudentID).Scan(&f.AttachmentID)

	pbm := &pb.Manager{VersionID: f.PVVersionID, Nodes: map[string]pb.Node{"n": {}}}
	repo := repository.NewSQLAdminRepository(db)
	svc := services.NewAdminService(repo, pbm, config.AppConfig{})
	jRepo := repository.NewSQLJourneyRepository(db)
	jSvc := services.NewJourneyService(jRepo, pbm, config.AppConfig{}, nil, nil, nil)
	h := handlers.NewAdminHandler(config.AppConfig{}, pbm, svc, jSvc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		// Use f.AdvisorID which is NOT assigned to student
		c.Set("claims", jwt.MapClaims{"sub": f.AdvisorID, "role": "advisor"})
		c.Next()
	})
	r.PATCH("/admin/attachments/:attachmentId/review", h.ReviewAttachment)

	body := map[string]string{"status": "approved"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("PATCH", "/admin/attachments/"+f.AttachmentID+"/review", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code, "Unassigned advisor should get 403")
}

func TestReviewAttachment_AdminApprovesAnyStudent(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	f := reviewTestFixtures{
		TenantID:    uuid.New().String(),
		StudentID:   uuid.New().String(),
		AdminID:     uuid.New().String(),
		PVVersionID: uuid.New().String(),
		NodeInstID:  uuid.New().String(),
		SlotID:      uuid.New().String(),
		DocID:       uuid.New().String(),
	}

	tenantSlug := "admin-approve-" + f.TenantID[:8]
	db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type, is_active) VALUES ($1, $2, 'T', 'university', true)`, f.TenantID, tenantSlug)
	db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) VALUES ($1, 'stu', 's@t.com', 'S', 'T', 'student', 'h', true), ($2, 'adm', 'adm@t.com', 'Admin', 'User', 'admin', 'h', true)`, f.StudentID, f.AdminID)
	// No advisor assigned - admin should still be able to review
	db.Exec(`INSERT INTO playbook_versions (id, tenant_id, version, checksum, raw_json) VALUES ($1, $2, 'v1', 'c', '{}')`, f.PVVersionID, f.TenantID)
	db.Exec(`INSERT INTO node_instances (id, tenant_id, node_id, user_id, state, playbook_version_id) VALUES ($1, $2, 'n', $3, 'submitted', $4)`, f.NodeInstID, f.TenantID, f.StudentID, f.PVVersionID)
	db.Exec(`INSERT INTO node_instance_slots (id, tenant_id, node_instance_id, slot_key) VALUES ($1, $2, $3, 'doc')`, f.SlotID, f.TenantID, f.NodeInstID)
	db.Exec(`INSERT INTO documents (id, tenant_id, user_id, title, kind) VALUES ($1, $2, $3, 'Test', 'other')`, f.DocID, f.TenantID, f.StudentID)
	db.QueryRow(`INSERT INTO document_versions (tenant_id, document_id, storage_path, mime_type, size_bytes, uploaded_by) VALUES ($1, $2, 'p', 'pdf', 100, $3) RETURNING id`, f.TenantID, f.DocID, f.StudentID).Scan(&f.DocVersionID)
	db.QueryRow(`INSERT INTO node_instance_slot_attachments (slot_id, document_version_id, is_active, status, filename, size_bytes, attached_by) VALUES ($1, $2, true, 'submitted', 't.pdf', 100, $3) RETURNING id`, f.SlotID, f.DocVersionID, f.StudentID).Scan(&f.AttachmentID)

	pbm := &pb.Manager{VersionID: f.PVVersionID, Nodes: map[string]pb.Node{"n": {}}}
	repo := repository.NewSQLAdminRepository(db)
	svc := services.NewAdminService(repo, pbm, config.AppConfig{})
	jRepo := repository.NewSQLJourneyRepository(db)
	jSvc := services.NewJourneyService(jRepo, pbm, config.AppConfig{}, nil, nil, nil)
	h := handlers.NewAdminHandler(config.AppConfig{}, pbm, svc, jSvc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": f.AdminID, "role": "admin"})
		c.Next()
	})
	r.PATCH("/admin/attachments/:attachmentId/review", h.ReviewAttachment)

	body := map[string]string{"status": "approved", "note": "Admin approved"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("PATCH", "/admin/attachments/"+f.AttachmentID+"/review", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Logf("Response: %s", w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code, "Admin should be able to approve any student")

	var state string
	db.Get(&state, `SELECT state FROM node_instances WHERE id=$1`, f.NodeInstID)
	assert.Equal(t, "done", state)
}

func TestNotifyAdvisorsOnSubmission(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	tenantID := uuid.New().String()
	tenantSlug := "notify-" + tenantID[:8]
	db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type, is_active) VALUES ($1, $2, 'T', 'university', true)`, tenantID, tenantSlug)

	studentID := uuid.New().String()
	advisor1ID := uuid.New().String()
	advisor2ID := uuid.New().String()
	db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) VALUES 
		($1, 's', 's@t.com', 'Test', 'Student', 'student', 'h', true),
		($2, 'a1', 'a1@t.com', 'A1', 'D', 'advisor', 'h', true),
		($3, 'a2', 'a2@t.com', 'A2', 'D', 'advisor', 'h', true)`, studentID, advisor1ID, advisor2ID)

	// Assign BOTH advisors to the student
	db.Exec(`INSERT INTO student_advisors (student_id, advisor_id, tenant_id) VALUES ($1, $2, $3), ($1, $4, $3)`, studentID, advisor1ID, tenantID, advisor2ID)

	pvID := uuid.New().String()
	db.Exec(`INSERT INTO playbook_versions (id, tenant_id, version, checksum, raw_json) VALUES ($1, $2, 'v1', 'c', '{}')`, pvID, tenantID)

	niID := uuid.New().String()
	db.Exec(`INSERT INTO node_instances (id, tenant_id, node_id, user_id, state, playbook_version_id) VALUES ($1, $2, 'confirm_task', $3, 'submitted', $4)`, niID, tenantID, studentID, pvID)

	// Call notify function
	err := services.NotifyAdvisorsOnSubmission(db, studentID, "confirm_task", niID, "")
	require.NoError(t, err)

	// Verify notification was created
	var count int
	err = db.Get(&count, `SELECT COUNT(*) FROM admin_notifications WHERE student_id=$1 AND event_type='document_submitted'`, studentID)
	require.NoError(t, err)
	assert.Equal(t, 1, count, "Expected 1 shared notification to be created")

	// Verify the message contains student name
	var message string
	db.Get(&message, `SELECT message FROM admin_notifications WHERE student_id=$1 AND event_type='document_submitted'`, studentID)
	assert.Contains(t, message, "Test Student")
}

func TestGetAdvisorsForStudent(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	studentID := uuid.New().String()
	advisor1ID := uuid.New().String()
	advisor2ID := uuid.New().String()

	db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) VALUES 
		($1, 's', 's@t.com', 'S', 'T', 'student', 'h', true),
		($2, 'a1', 'a1@t.com', 'A1', 'D', 'advisor', 'h', true),
		($3, 'a2', 'a2@t.com', 'A2', 'D', 'advisor', 'h', true)`, studentID, advisor1ID, advisor2ID)

	tenantID := uuid.New().String()
	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type, is_active) VALUES ($1, $2, 'T', 'university', true)`, tenantID, "slug-"+tenantID)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO student_advisors (student_id, advisor_id, tenant_id) VALUES ($1, $2, $3), ($1, $4, $3)`, studentID, advisor1ID, tenantID, advisor2ID)
	require.NoError(t, err)

	advisors, err := services.GetAdvisorsForStudent(db, studentID)
	require.NoError(t, err)
	assert.Len(t, advisors, 2)
	assert.Contains(t, advisors, advisor1ID)
	assert.Contains(t, advisors, advisor2ID)
}

func TestHasAdvisors(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	studentWithAdvisor := uuid.New().String()
	studentWithoutAdvisor := uuid.New().String()
	advisorID := uuid.New().String()
	tenantID := uuid.New().String()

	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type, is_active) VALUES ($1, $2, 'T', 'university', true)`, tenantID, "slug-"+tenantID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) VALUES 
		($1, 's1', 's1@t.com', 'S', '1', 'student', 'h', true),
		($2, 's2', 's2@t.com', 'S', '2', 'student', 'h', true),
		($3, 'a', 'a@t.com', 'A', '1', 'advisor', 'h', true)`, studentWithAdvisor, studentWithoutAdvisor, advisorID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO student_advisors (student_id, advisor_id, tenant_id) VALUES ($1, $2, $3)`, studentWithAdvisor, advisorID, tenantID)
	require.NoError(t, err)
	// Deleted duplicate setup code

	// Student with advisor
	has, err := services.HasAdvisors(db, studentWithAdvisor)
	require.NoError(t, err)
	assert.True(t, has)

	// Student without advisor
	has, err = services.HasAdvisors(db, studentWithoutAdvisor)
	require.NoError(t, err)
	assert.False(t, has)
}

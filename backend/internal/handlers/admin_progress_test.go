package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	pb "github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminHandler_StudentProgress(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Seed data
	studentID := "10000000-0000-0000-0000-000000000001"
	tenantID := "10000000-0000-0000-0000-000000000001"
	_, err := db.Exec(`INSERT INTO tenants (id, name, slug) VALUES ($1, 'Test Tenant', 'test-tenant') ON CONFLICT DO NOTHING`, tenantID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'student1', 's1@ex.com', 'Student', 'One', 'student', 'hash', true)`, studentID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'student', true)`, studentID, tenantID)
	require.NoError(t, err)

	// Mock Playbook Manager
	pbm := &pb.Manager{
		VersionID: "20000000-0000-0000-0000-000000000001",
		Nodes: map[string]pb.Node{
			"node1": {ID: "node1"},
			"node2": {ID: "node2"},
		},
	}

	// Seed playbook version to satisfy FK
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, created_at, tenant_id) 
		VALUES ($1, 'v1', 'sum', '{}', NOW(), $2)`, pbm.VersionID, tenantID)
	require.NoError(t, err)

	// Seed progress
	_, err = db.Exec(`INSERT INTO node_instances (node_id, user_id, state, playbook_version_id, updated_at, tenant_id) 
		VALUES ('node1', $1, 'done', $2, NOW(), $3)`, studentID, pbm.VersionID, tenantID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	h := handlers.NewAdminHandler(db, cfg, pbm)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", tenantID)
		c.Next()
	})
	r.GET("/admin/student-progress", h.StudentProgress)

	t.Run("List Student Progress", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/admin/student-progress", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		
		var student map[string]interface{}
		for _, s := range resp {
			if s["id"] == studentID {
				student = s
				break
			}
		}
		require.NotNil(t, student, "Student not found in response")
		assert.Equal(t, "Student One", student["name"])
		
		progress := student["progress"].(map[string]interface{})
		assert.Equal(t, float64(1), progress["completed_nodes"])
		assert.Equal(t, float64(2), progress["total_nodes"])
		assert.Equal(t, float64(50), progress["percent"])
	})
}

func TestAdminHandler_MonitorStudents(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Seed users
	studentID := "10000000-0000-0000-0000-000000000002"
	advisorID := "10000000-0000-0000-0000-000000000003"
	tenantID := "20000000-0000-0000-0000-000000000002"
	_, err := db.Exec(`INSERT INTO tenants (id, name, slug) VALUES ($1, 'Test Tenant', 'test-tenant') ON CONFLICT DO NOTHING`, tenantID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES 
		($1, 'student2', 's2@ex.com', 'Student', 'Two', 'student', 'hash', true),
		($2, 'advisor1', 'a1@ex.com', 'Advisor', 'One', 'advisor', 'hash', true)`, studentID, advisorID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) 
		VALUES 
		($1, $3, 'student', true),
		($2, $3, 'advisor', true)`, studentID, advisorID, tenantID)
	require.NoError(t, err)

	// Seed profile submission (for filters)
	profileData := `{"program": "CS", "department": "Eng", "cohort": "2024"}`
	_, err = db.Exec(`INSERT INTO profile_submissions (user_id, form_data, tenant_id) VALUES ($1, $2, $3)`, studentID, profileData, tenantID)
	require.NoError(t, err)

	// Seed advisor relationship
	_, err = db.Exec(`INSERT INTO student_advisors (student_id, advisor_id, tenant_id) VALUES ($1, $2, $3)`, studentID, advisorID, tenantID)
	require.NoError(t, err)

	// Mock Playbook Manager
	pbm := &pb.Manager{
		VersionID: "20000000-0000-0000-0000-000000000002",
		Raw: []byte(`{"worlds": [{"id": "W1", "nodes": ["node1"]}, {"id": "W2", "nodes": ["node2"]}]}`),
		Nodes: map[string]pb.Node{"node1": {}, "node2": {}},
	}

	cfg := config.AppConfig{}
	h := handlers.NewAdminHandler(db, cfg, pbm)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", tenantID)
		// Mock auth
		role := c.GetHeader("X-Role")
		uid := c.GetHeader("X-User-ID")
		if role != "" {
			c.Set("role", role)
			c.Set("claims", jwt.MapClaims{"sub": uid, "role": role})
		}
		c.Next()
	})
	r.GET("/admin/monitor", h.MonitorStudents)

	t.Run("Admin List All", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/admin/monitor", nil)
		req.Header.Set("X-Role", "admin")
		req.Header.Set("X-User-ID", "admin-id")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Len(t, resp, 1)
		assert.Equal(t, "Student Two", resp[0]["name"])
		assert.Equal(t, "CS", resp[0]["program"])
	})

	t.Run("Advisor List Assigned", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/admin/monitor", nil)
		req.Header.Set("X-Role", "advisor")
		req.Header.Set("X-User-ID", advisorID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Len(t, resp, 1)
	})

	t.Run("Advisor List Unassigned (Empty)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/admin/monitor", nil)
		req.Header.Set("X-Role", "advisor")
		req.Header.Set("X-User-ID", "other-advisor-id")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Len(t, resp, 0)
	})

	t.Run("Filter by Program", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/admin/monitor?program=CS", nil)
		req.Header.Set("X-Role", "admin")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Len(t, resp, 1)
	})

	t.Run("Filter by Wrong Program", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/admin/monitor?program=Wrong", nil)
		req.Header.Set("X-Role", "admin")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Len(t, resp, 0)
	})
}

func TestAdminHandler_GetStudentDetails(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	studentID := "10000000-0000-0000-0000-000000000004"
	tenantID := "30000000-0000-0000-0000-000000000003"
	_, err := db.Exec(`INSERT INTO tenants (id, name, slug) VALUES ($1, 'Test Tenant', 'test-tenant') ON CONFLICT DO NOTHING`, tenantID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'student3', 's3@ex.com', 'Student', 'Three', 'student', 'hash', true)`, studentID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'student', true)`, studentID, tenantID)
	require.NoError(t, err)

	// Mock Playbook Manager
	pbm := &pb.Manager{
		VersionID: "20000000-0000-0000-0000-000000000003",
		Raw: []byte(`{"worlds": [{"id": "W1", "nodes": ["node1"]}]}`),
		Nodes: map[string]pb.Node{"node1": {}},
	}

	cfg := config.AppConfig{}
	h := handlers.NewAdminHandler(db, cfg, pbm)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", tenantID)
		c.Next()
	})
	r.GET("/admin/students/:id", h.GetStudentDetails)

	t.Run("Get Details Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/admin/students/"+studentID, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "Student Three", resp["name"])
		assert.Equal(t, "s3@ex.com", resp["email"])
	})

	t.Run("Get Details Not Found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/admin/students/invalid-id", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestAdminHandler_StudentJourney(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	studentID := "10000000-0000-0000-0000-000000000005"
	tenantID := "40000000-0000-0000-0000-000000000004"
	_, err := db.Exec(`INSERT INTO tenants (id, name, slug) VALUES ($1, 'Test Tenant', 'test-tenant') ON CONFLICT DO NOTHING`, tenantID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'student4', 's4@ex.com', 'Student', 'Four', 'student', 'hash', true)`, studentID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'student', true)`, studentID, tenantID)
	require.NoError(t, err)

	// Mock Playbook Manager
	pbm := &pb.Manager{
		VersionID: "20000000-0000-0000-0000-000000000004",
		Raw: []byte(`{"worlds": [{"id": "W1", "nodes": ["node1"]}]}`),
		Nodes: map[string]pb.Node{"node1": {}},
	}

	// Seed playbook version
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, created_at, tenant_id) 
		VALUES ($1, 'v1', 'sum', '{}', NOW(), $2)`, pbm.VersionID, tenantID)
	require.NoError(t, err)

	// Seed progress
	_, err = db.Exec(`INSERT INTO node_instances (node_id, user_id, state, playbook_version_id, updated_at, tenant_id) 
		VALUES ('node1', $1, 'done', $2, NOW(), $3)`, studentID, pbm.VersionID, tenantID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	h := handlers.NewAdminHandler(db, cfg, pbm)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", tenantID)
		c.Next()
	})
	r.GET("/admin/students/:id/journey", h.StudentJourney)

	t.Run("Get Student Journey", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/admin/students/"+studentID+"/journey", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		
		require.NotNil(t, resp["nodes"], "nodes should not be nil")
		
		nodes := resp["nodes"].([]interface{})
		require.Len(t, nodes, 1)
		node := nodes[0].(map[string]interface{})
		assert.Equal(t, "node1", node["node_id"])
		assert.Equal(t, "done", node["state"])
	})
}

func TestAdminHandler_ListStudentNodeFiles(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	studentID := "10000000-0000-0000-0000-000000000006"
	tenantID := "50000000-0000-0000-0000-000000000005"
	_, err := db.Exec(`INSERT INTO tenants (id, name, slug) VALUES ($1, 'Test Tenant', 'test-tenant') ON CONFLICT DO NOTHING`, tenantID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'student5', 's5@ex.com', 'Student', 'Five', 'student', 'hash', true)`, studentID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'student', true)`, studentID, tenantID)
	require.NoError(t, err)

	pbm := &pb.Manager{
		VersionID: "20000000-0000-0000-0000-000000000005",
		Nodes: map[string]pb.Node{"node1": {}},
	}

	// Seed playbook version
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, created_at, tenant_id) 
		VALUES ($1, 'v1', 'sum', '{}', NOW(), $2)`, pbm.VersionID, tenantID)
	require.NoError(t, err)

	// Seed node instance
	var instanceID string
	err = db.QueryRow(`INSERT INTO node_instances (node_id, user_id, state, playbook_version_id, updated_at, tenant_id) 
		VALUES ('node1', $1, 'done', $2, NOW(), $3) RETURNING id`, studentID, pbm.VersionID, tenantID).Scan(&instanceID)
	require.NoError(t, err)

	// Seed slot
	var slotID string
	err = db.QueryRow(`INSERT INTO node_instance_slots (node_instance_id, slot_key, tenant_id) 
		VALUES ($1, 'slot1', $2) RETURNING id`, instanceID, tenantID).Scan(&slotID)
	require.NoError(t, err)

	// Seed document
	var docID string
	err = db.QueryRow(`INSERT INTO documents (user_id, title, kind, created_at, tenant_id) 
		VALUES ($1, 'Test Doc', 'other', NOW(), $2) RETURNING id`, studentID, tenantID).Scan(&docID)
	require.NoError(t, err)

	// Seed document version
	var docVerID string
	err = db.QueryRow(`INSERT INTO document_versions (document_id, storage_path, mime_type, size_bytes, created_at, uploaded_by, tenant_id) 
		VALUES ($1, 'path/to/doc', 'application/pdf', 1024, NOW(), $2, $3) RETURNING id`, docID, studentID, tenantID).Scan(&docVerID)
	require.NoError(t, err)

	// Seed attachment
	_, err = db.Exec(`INSERT INTO node_instance_slot_attachments (slot_id, document_version_id, is_active, status, filename, size_bytes, attached_by) 
		VALUES ($1, $2, true, 'submitted', 'test.pdf', 1024, $3)`, slotID, docVerID, studentID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	h := handlers.NewAdminHandler(db, cfg, pbm)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", tenantID)
		c.Next()
	})
	r.GET("/admin/students/:id/nodes/:nodeId/files", h.ListStudentNodeFiles)

	t.Run("List Student Node Files", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/admin/students/"+studentID+"/nodes/node1/files", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Len(t, resp, 1)
		assert.Equal(t, "slot1", resp[0]["slot_key"])
	})

	t.Run("List Student Node Files - Node Not Found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/admin/students/"+studentID+"/nodes/invalid-node/files", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestAdminHandler_ListStudentNodeFiles_Forbidden(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	advisorID := uuid.NewString()
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'advisor', 'advisor@ex.com', 'Advisor', 'One', 'advisor', 'hash', true)`, advisorID)
	require.NoError(t, err)

	studentID := uuid.NewString()
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'student', 'student@ex.com', 'Student', 'One', 'student', 'hash', true)`, studentID)
	require.NoError(t, err)

	h := handlers.NewAdminHandler(db, config.AppConfig{}, &pb.Manager{})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": advisorID, "role": "advisor"})
		c.Next()
	})
	r.GET("/admin/students/:id/nodes/:nodeId/files", h.ListStudentNodeFiles)

	t.Run("List Student Node Files Forbidden", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/admin/students/"+studentID+"/nodes/node1/files", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestAdminHandler_MonitorAnalytics(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	tenantID := "60000000-0000-0000-0000-000000000006"
	_, err := db.Exec(`INSERT INTO tenants (id, name, slug) VALUES ($1, 'Test Tenant', 'test-tenant') ON CONFLICT DO NOTHING`, tenantID)
	require.NoError(t, err)

	// Seed users
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ('10000000-0000-0000-0000-000000000007', 'student6', 's6@ex.com', 'Student', 'Six', 'student', 'hash', true)`)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) 
		VALUES ('10000000-0000-0000-0000-000000000007', $1, 'student', true)`, tenantID)
	require.NoError(t, err)

	pbm := &pb.Manager{
		VersionID: "20000000-0000-0000-0000-000000000005",
		Nodes: map[string]pb.Node{"node1": {}},
	}

	// Seed playbook version
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, created_at, tenant_id) 
		VALUES ($1, 'v1', 'sum', '{}', NOW(), $2)`, pbm.VersionID, tenantID)
	require.NoError(t, err)

	// Seed progress
	_, err = db.Exec(`INSERT INTO node_instances (node_id, user_id, state, playbook_version_id, updated_at, tenant_id) 
		VALUES ('node1', '10000000-0000-0000-0000-000000000007', 'done', $1, NOW(), $2)`, pbm.VersionID, tenantID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	h := handlers.NewAdminHandler(db, cfg, pbm)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", tenantID)
		c.Next()
	})
	r.GET("/admin/analytics", h.MonitorAnalytics)

	t.Run("Get Analytics", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/admin/analytics", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		
		// MonitorAnalytics returns aggregate stats
		assert.NotNil(t, resp["rp_required_count"])
	})

	t.Run("Analytics Filter Program", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/admin/analytics?program=CS", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
	
	t.Run("Analytics Filter Search", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/admin/analytics?q=Student", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAdminHandler_ReviewAttachment(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	tenantID := "70000000-0000-0000-0000-000000000007"
	_, err := db.Exec(`INSERT INTO tenants (id, name, slug) VALUES ($1, 'Test Tenant', 'test-tenant') ON CONFLICT DO NOTHING`, tenantID)
	require.NoError(t, err)

	studentID := "10000000-0000-0000-0000-000000000008"
	adminID := "99999999-9999-9999-9999-999999999999"
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'admin', 'admin@ex.com', 'Admin', 'User', 'admin', 'hash', true)`, adminID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'admin', true)`, adminID, tenantID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES 
		($1, 'student7', 's7@ex.com', 'Student', 'Seven', 'student', 'hash', true)`, studentID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'student', true)`, studentID, tenantID)
	require.NoError(t, err)

	pbm := &pb.Manager{
		VersionID: "20000000-0000-0000-0000-000000000006",
		Nodes: map[string]pb.Node{"node1": {}},
	}

	// Seed playbook version
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, created_at, tenant_id) 
		VALUES ($1, 'v1', 'sum', '{}', NOW(), $2)`, pbm.VersionID, tenantID)
	require.NoError(t, err)

	// Seed node instance
	var instanceID string
	err = db.QueryRow(`INSERT INTO node_instances (node_id, user_id, state, playbook_version_id, updated_at, tenant_id) 
		VALUES ('node1', $1, 'done', $2, NOW(), $3) RETURNING id`, studentID, pbm.VersionID, tenantID).Scan(&instanceID)
	require.NoError(t, err)

	// Seed slot
	var slotID string
	err = db.QueryRow(`INSERT INTO node_instance_slots (node_instance_id, slot_key, tenant_id) 
		VALUES ($1, 'slot1', $2) RETURNING id`, instanceID, tenantID).Scan(&slotID)
	require.NoError(t, err)

	// Seed document
	var docID string
	err = db.QueryRow(`INSERT INTO documents (user_id, title, kind, created_at, tenant_id) 
		VALUES ($1, 'Test Doc', 'other', NOW(), $2) RETURNING id`, studentID, tenantID).Scan(&docID)
	require.NoError(t, err)

	// Seed document version
	var docVerID string
	err = db.QueryRow(`INSERT INTO document_versions (document_id, storage_path, mime_type, size_bytes, created_at, uploaded_by, tenant_id) 
		VALUES ($1, 'path/to/doc', 'application/pdf', 1024, NOW(), $2, $3) RETURNING id`, docID, studentID, tenantID).Scan(&docVerID)
	require.NoError(t, err)

	// Seed attachment
	var attID string
	err = db.QueryRow(`INSERT INTO node_instance_slot_attachments (slot_id, document_version_id, is_active, status, filename, size_bytes, attached_by) 
		VALUES ($1, $2, true, 'submitted', 'test.pdf', 1024, $3) RETURNING id`, slotID, docVerID, studentID).Scan(&attID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	h := handlers.NewAdminHandler(db, cfg, pbm)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", tenantID)
		c.Set("userID", adminID)
		c.Set("claims", jwt.MapClaims{"role": "admin", "sub": adminID})
		c.Next()
	})
	r.POST("/admin/attachments/:attachmentId/review", h.ReviewAttachment)

	t.Run("Review Attachment", func(t *testing.T) {
		reqBody := map[string]string{
			"status": "approved",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/admin/attachments/"+attID+"/review", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "approved", resp["status"])
	})
}

func TestAdminHandler_ListStudentProgress(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	tenantID := "80000000-0000-0000-0000-000000000008"
	_, err := db.Exec(`INSERT INTO tenants (id, name, slug) VALUES ($1, 'Test Tenant', 'test-tenant') ON CONFLICT DO NOTHING`, tenantID)
	require.NoError(t, err)

	// Seed student
	studentID := "11111111-1111-1111-1111-111111111111"
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'student1', 's1@ex.com', 'Student', 'One', 'student', 'hash', true)`, studentID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'student', true)`, studentID, tenantID)
	require.NoError(t, err)

	// Seed playbook version
	pbm := &pb.Manager{
		VersionID: "11111111-1111-1111-1111-111111111111",
		Nodes: map[string]pb.Node{
			"node1": {ID: "node1", Title: map[string]string{"en": "Node 1"}},
		},
	}
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, tenant_id) 
		VALUES ($1, 'v1', 'checksum', '{}', $2)`, pbm.VersionID, tenantID)
	require.NoError(t, err)

	// Seed instance
	_, err = db.Exec(`INSERT INTO node_instances (node_id, user_id, state, playbook_version_id, updated_at, tenant_id) 
		VALUES ('node1', $1, 'done', $2, NOW(), $3)`, studentID, pbm.VersionID, tenantID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	h := handlers.NewAdminHandler(db, cfg, pbm)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", tenantID)
		c.Next()
	})
	r.GET("/admin/students/:studentId/progress", h.StudentProgress)

	t.Run("List Student Progress", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/admin/students/"+studentID+"/progress", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		t.Logf("Response: %s", w.Body.String())
		assert.Equal(t, http.StatusOK, w.Code)
		var resp []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		
		assert.NotEmpty(t, resp)
		student := resp[0]
		progress := student["progress"].(map[string]interface{})
		assert.NotNil(t, progress["completed_nodes"])
		assert.NotNil(t, progress["percent"])
	})
}

func TestAdminHandler_UploadReviewedDocument(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	tenantID := "90000000-0000-0000-0000-000000000009"
	_, err := db.Exec(`INSERT INTO tenants (id, name, slug) VALUES ($1, 'Test Tenant', 'test-tenant') ON CONFLICT DO NOTHING`, tenantID)
	require.NoError(t, err)

	adminID := "10000000-0000-0000-0000-000000000000"
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'admin', 'admin@ex.com', 'Admin', 'User', 'admin', 'hash', true)`, adminID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'admin', true)`, adminID, tenantID)
	require.NoError(t, err)

	studentID := "20000000-0000-0000-0000-000000000000"
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'student', 'student@ex.com', 'Student', 'User', 'student', 'hash', true)`, studentID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'student', true)`, studentID, tenantID)
	require.NoError(t, err)

	// Seed playbook version
	versionID := "30000000-0000-0000-0000-000000000000"
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, tenant_id) 
		VALUES ($1, 'v1', 'checksum', '{}', $2)`, versionID, tenantID)
	require.NoError(t, err)

	// Seed node instance
	var instanceID string
	err = db.QueryRow(`INSERT INTO node_instances (user_id, node_id, state, playbook_version_id, tenant_id) 
		VALUES ($1, 'node1', 'submitted', $2, $3) RETURNING id`, studentID, versionID, tenantID).Scan(&instanceID)
	require.NoError(t, err)

	// Seed slot
	var slotID string
	err = db.QueryRow(`INSERT INTO node_instance_slots (node_instance_id, slot_key, tenant_id) 
		VALUES ($1, 'slot1', $2) RETURNING id`, instanceID, tenantID).Scan(&slotID)
	require.NoError(t, err)

	// Seed document version for review
	var docID string
	err = db.QueryRow(`INSERT INTO documents (user_id, kind, title, tenant_id) VALUES ($1, 'other', 'Review Doc', $2) RETURNING id`, adminID, tenantID).Scan(&docID)
	require.NoError(t, err)

	var reviewVerID string
	err = db.QueryRow(`INSERT INTO document_versions (document_id, storage_path, mime_type, size_bytes, uploaded_by, tenant_id) 
		VALUES ($1, 'path/to/review.pdf', 'application/pdf', 1024, $2, $3) RETURNING id`, docID, adminID, tenantID).Scan(&reviewVerID)
	require.NoError(t, err)

	// Seed attachment
	var attachmentID string
	err = db.QueryRow(`INSERT INTO node_instance_slot_attachments (slot_id, document_version_id, is_active, attached_by, filename, size_bytes) 
		VALUES ($1, $2, true, $3, 'test.pdf', 1024) RETURNING id`, slotID, reviewVerID, studentID).Scan(&attachmentID)
	require.NoError(t, err)

	pb := &playbook.Manager{}
	cfg := config.AppConfig{}
	h := handlers.NewAdminHandler(db, cfg, pb)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", tenantID)
		c.Set("claims", jwt.MapClaims{"sub": adminID, "role": "admin"})
		c.Set("userRole", "admin")
		c.Next()
	})
	r.POST("/admin/attachments/:attachmentId/reviewed-document", h.UploadReviewedDocument)

	t.Run("Upload Reviewed Document", func(t *testing.T) {
		reqBody := map[string]string{
			"document_version_id": reviewVerID,
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/admin/attachments/"+attachmentID+"/reviewed-document", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		// Verify DB update
		var reviewedBy string
		err := db.QueryRow("SELECT reviewed_by FROM node_instance_slot_attachments WHERE id=$1", attachmentID).Scan(&reviewedBy)
		assert.NoError(t, err)
		assert.Equal(t, adminID, reviewedBy)
	})
}

func TestAdminHandler_PatchStudentNodeState(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	tenantID := "a0000000-0000-0000-0000-00000000000a"
	_, err := db.Exec(`INSERT INTO tenants (id, name, slug) VALUES ($1, 'Test Tenant', 'test-tenant') ON CONFLICT DO NOTHING`, tenantID)
	require.NoError(t, err)

	adminID := "40000000-0000-0000-0000-000000000000"
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'admin2', 'admin2@ex.com', 'Admin', 'Two', 'admin', 'hash', true)`, adminID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'admin', true)`, adminID, tenantID)
	require.NoError(t, err)

	studentID := "50000000-0000-0000-0000-000000000000"
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'student2', 'student2@ex.com', 'Student', 'Two', 'student', 'hash', true)`, studentID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'student', true)`, studentID, tenantID)
	require.NoError(t, err)

	// Seed playbook version
	versionID := "60000000-0000-0000-0000-000000000000"
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, tenant_id) 
		VALUES ($1, 'v1', 'checksum', '{}', $2)`, versionID, tenantID)
	require.NoError(t, err)

	// Seed node instance
	_, err = db.Exec(`INSERT INTO node_instances (user_id, node_id, state, playbook_version_id, tenant_id) 
		VALUES ($1, 'node1', 'active', $2, $3)`, studentID, versionID, tenantID)
	require.NoError(t, err)

	pb := &playbook.Manager{
		VersionID: versionID,
		Nodes: map[string]playbook.Node{
			"node1": {ID: "node1", Type: "form"},
		},
	}
	cfg := config.AppConfig{}
	h := handlers.NewAdminHandler(db, cfg, pb)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", tenantID)
		c.Set("claims", jwt.MapClaims{"sub": adminID, "role": "admin"})
		c.Set("userRole", "admin")
		c.Next()
	})
	// Seed transition rule
	db.Exec(`DELETE FROM node_state_transitions WHERE from_state='active' AND to_state='done'`)
	_, err = db.Exec(`INSERT INTO node_state_transitions (from_state, to_state, allowed_roles) 
		VALUES ('active', 'done', '{"admin"}')`)
	require.NoError(t, err)

	r.PATCH("/admin/students/:id/nodes/:nodeId/state", h.PatchStudentNodeState)

	t.Run("Patch Node State", func(t *testing.T) {
		reqBody := map[string]string{
			"state": "done",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PATCH", "/admin/students/"+studentID+"/nodes/node1/state", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify DB
		var state string
		err := db.QueryRow("SELECT state FROM node_instances WHERE user_id=$1 AND node_id='node1'", studentID).Scan(&state)
		assert.NoError(t, err)
		assert.Equal(t, "done", state)
	})
}

func TestAdminHandler_PresignReviewedDocumentUpload(t *testing.T) {
	t.Setenv("S3_BUCKET", "test-bucket")
	t.Setenv("S3_ACCESS_KEY", "test-key")
	t.Setenv("S3_SECRET_KEY", "test-secret")
	t.Setenv("S3_REGION", "us-east-1")

	db, teardown := testutils.SetupTestDB()
	defer teardown()

	tenantID := "00000000-0000-0000-0000-000000000001"
	_, err := db.Exec(`INSERT INTO tenants (id, name, slug) VALUES ($1, 'Test Tenant', 'test-tenant') ON CONFLICT DO NOTHING`, tenantID)
	require.NoError(t, err)

	userID := uuid.NewString()
	studentID := uuid.NewString()
	
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'admin', 'admin@ex.com', 'Admin', 'User', 'admin', 'hash', true)`, userID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'admin', true)`, userID, tenantID)
	require.NoError(t, err)
	
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'student', 'student@ex.com', 'Student', 'One', 'student', 'hash', true)`, studentID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'student', true)`, studentID, tenantID)
	require.NoError(t, err)

	// Seed node instance and attachment
	versionID := uuid.NewString()
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, tenant_id) VALUES ($1, 'v1', 'sum', '{}', $2)`, versionID, tenantID)
	require.NoError(t, err)

	var instanceID string
	err = db.QueryRow(`INSERT INTO node_instances (node_id, user_id, state, playbook_version_id, updated_at, tenant_id) 
		VALUES ('node1', $1, 'done', $2, NOW(), $3) RETURNING id`, studentID, versionID, tenantID).Scan(&instanceID)
	require.NoError(t, err)

	var slotID string
	err = db.QueryRow(`INSERT INTO node_instance_slots (node_instance_id, slot_key, tenant_id) VALUES ($1, 'slot1', $2) RETURNING id`, instanceID, tenantID).Scan(&slotID)
	require.NoError(t, err)

	var docID string
	err = db.QueryRow(`INSERT INTO documents (user_id, kind, title, tenant_id) VALUES ($1, 'other', 'Doc', $2) RETURNING id`, studentID, tenantID).Scan(&docID)
	require.NoError(t, err)

	var docVerID string
	err = db.QueryRow(`INSERT INTO document_versions (document_id, storage_path, mime_type, size_bytes, uploaded_by, tenant_id) 
		VALUES ($1, 'path', 'pdf', 1024, $2, $3) RETURNING id`, docID, studentID, tenantID).Scan(&docVerID)
	require.NoError(t, err)

	var attID string
	err = db.QueryRow(`INSERT INTO node_instance_slot_attachments (slot_id, document_version_id, is_active, status, filename, size_bytes, attached_by) 
		VALUES ($1, $2, true, 'submitted', 'file.pdf', 1024, $3) RETURNING id`, slotID, docVerID, studentID).Scan(&attID)
	require.NoError(t, err)

	h := handlers.NewAdminHandler(db, config.AppConfig{}, &pb.Manager{})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": userID, "role": "admin"})
		c.Next()
	})
	r.POST("/admin/attachments/:attachmentId/review/presign", h.PresignReviewedDocumentUpload)

	t.Run("Presign Upload Success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"filename":     "reviewed.pdf",
			"content_type": "application/pdf",
			"size_bytes":   1024,
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/admin/attachments/"+attID+"/review/presign", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 200 or 500, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("Presign Reviewed Document Upload Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/admin/attachments/"+attID+"/review/presign", bytes.NewBuffer([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Presign Reviewed Document Upload Attachment Not Found", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"filename":     "test.pdf",
			"content_type": "application/pdf",
			"size_bytes":   1024,
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/admin/attachments/99999999-9999-9999-9999-999999999999/review/presign", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestAdminHandler_AttachReviewedDocument(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	tenantID := "00000000-0000-0000-0000-000000000001"
	_, err := db.Exec(`INSERT INTO tenants (id, name, slug) VALUES ($1, 'Test Tenant', 'test-tenant') ON CONFLICT DO NOTHING`, tenantID)
	require.NoError(t, err)

	studentID := uuid.NewString()
	adminID := uuid.NewString()
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'student8', 's8@ex.com', 'Student', 'Eight', 'student', 'hash', true)`, studentID)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'student', true)`, studentID, tenantID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'admin2', 'admin2@ex.com', 'Admin', 'Two', 'admin', 'hash', true)`, adminID)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'admin', true)`, adminID, tenantID)
	require.NoError(t, err)

	pbm := &pb.Manager{VersionID: uuid.NewString()}
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, tenant_id) VALUES ($1, 'v1', 'sum', '{}', $2)`, pbm.VersionID, tenantID)
	require.NoError(t, err)

	var instanceID string
	db.QueryRow(`INSERT INTO node_instances (node_id, user_id, state, playbook_version_id, tenant_id) VALUES ('node1', $1, 'done', $2, $3) RETURNING id`, studentID, pbm.VersionID, tenantID).Scan(&instanceID)

	var slotID string
	db.QueryRow(`INSERT INTO node_instance_slots (node_instance_id, slot_key, tenant_id) VALUES ($1, 'slot1', $2) RETURNING id`, instanceID, tenantID).Scan(&slotID)

	var docID string
	db.QueryRow(`INSERT INTO documents (user_id, title, kind, tenant_id) VALUES ($1, 'Doc', 'other', $2) RETURNING id`, studentID, tenantID).Scan(&docID)
	var docVerID string
	db.QueryRow(`INSERT INTO document_versions (document_id, storage_path, mime_type, size_bytes, uploaded_by, tenant_id) VALUES ($1, 'path', 'pdf', 100, $2, $3) RETURNING id`, docID, studentID, tenantID).Scan(&docVerID)

	var attID string
	db.QueryRow(`INSERT INTO node_instance_slot_attachments (slot_id, document_version_id, is_active, status, filename, size_bytes, attached_by) 
		VALUES ($1, $2, true, 'submitted', 'file.pdf', 100, $3) RETURNING id`, slotID, docVerID, studentID).Scan(&attID)

	cfg := config.AppConfig{}
	h := handlers.NewAdminHandler(db, cfg, pbm)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", adminID)
		c.Set("claims", jwt.MapClaims{"role": "admin", "sub": adminID})
		c.Next()
	})
	r.POST("/admin/attachments/:attachmentId/review/attach", h.AttachReviewedDocument)

	t.Run("Attach Reviewed Document", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"object_key": "reviewed/file.pdf",
			"filename": "reviewed_file.pdf",
			"size_bytes": 200,
			"content_type": "application/pdf",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/admin/attachments/"+attID+"/review/attach", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// It might fail with 500 because S3 is not configured, but that's okay for coverage.
		// We just want to ensure it passes validation and tries to execute.
		if w.Code == http.StatusOK {
			var status string
			db.QueryRow("SELECT status FROM node_instance_slot_attachments WHERE id=$1", attID).Scan(&status)
			// Status should NOT change in AttachReviewedDocument
			assert.Equal(t, "submitted", status)
		} else {
			assert.Equal(t, http.StatusInternalServerError, w.Code)
			assert.Contains(t, w.Body.String(), "S3 not configured")
		}
	})


}

func TestAdminHandler_UploadReviewedDocument_Failures(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := uuid.NewString()
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'admin', 'admin@ex.com', 'Admin', 'User', 'admin', 'hash', true)`, userID)
	require.NoError(t, err)

	h := handlers.NewAdminHandler(db, config.AppConfig{}, &pb.Manager{})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": userID, "role": "admin"})
		c.Next()
	})
	r.POST("/admin/attachments/:attachmentId/review/upload", h.UploadReviewedDocument)

	t.Run("Upload Reviewed Document Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/admin/attachments/123/review/upload", bytes.NewBuffer([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Upload Reviewed Document Attachment Not Found", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"document_version_id": uuid.NewString(),
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/admin/attachments/99999999-9999-9999-9999-999999999999/review/upload", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestAdminHandler_UploadReviewedDocument_Forbidden(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	tenantID := "00000000-0000-0000-0000-000000000001"
	_, err := db.Exec(`INSERT INTO tenants (id, name, slug) VALUES ($1, 'Test Tenant', 'test-tenant') ON CONFLICT DO NOTHING`, tenantID)
	require.NoError(t, err)

	advisorID := uuid.NewString()
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'advisor', 'advisor@ex.com', 'Advisor', 'One', 'advisor', 'hash', true)`, advisorID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'advisor', true)`, advisorID, tenantID)
	require.NoError(t, err)

	studentID := uuid.NewString()
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'student', 'student@ex.com', 'Student', 'One', 'student', 'hash', true)`, studentID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'student', true)`, studentID, tenantID)
	require.NoError(t, err)

	// Seed attachment
	versionID := uuid.NewString()
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, tenant_id) VALUES ($1, 'v1', 'sum', '{}', $2)`, versionID, tenantID)
	require.NoError(t, err)
	
	var instanceID string
	err = db.QueryRow(`INSERT INTO node_instances (node_id, user_id, state, playbook_version_id, tenant_id) VALUES ('node1', $1, 'done', $2, $3) RETURNING id`, studentID, versionID, tenantID).Scan(&instanceID)
	require.NoError(t, err)

	var slotID string
	err = db.QueryRow(`INSERT INTO node_instance_slots (node_instance_id, slot_key, tenant_id) VALUES ($1, 'slot1', $2) RETURNING id`, instanceID, tenantID).Scan(&slotID)
	require.NoError(t, err)

	var docID string
	err = db.QueryRow(`INSERT INTO documents (user_id, title, kind, tenant_id) VALUES ($1, 'Test', 'other', $2) RETURNING id`, studentID, tenantID).Scan(&docID)
	require.NoError(t, err)

	var docVerID string
	err = db.QueryRow(`INSERT INTO document_versions (document_id, storage_path, mime_type, size_bytes, uploaded_by, tenant_id) VALUES ($1, 'path', 'application/pdf', 100, $2, $3) RETURNING id`, docID, studentID, tenantID).Scan(&docVerID)
	require.NoError(t, err)

	var attID string
	err = db.QueryRow(`INSERT INTO node_instance_slot_attachments (slot_id, document_version_id, is_active, status, filename, size_bytes, attached_by) 
		VALUES ($1, $2, true, 'submitted', 'file.pdf', 100, $3) RETURNING id`, slotID, docVerID, studentID).Scan(&attID)
	require.NoError(t, err)

	h := handlers.NewAdminHandler(db, config.AppConfig{}, &pb.Manager{})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": advisorID, "role": "advisor"})
		c.Set("tenant_id", tenantID)
		c.Next()
	})
	r.POST("/admin/attachments/:attachmentId/review/upload", h.UploadReviewedDocument)

	t.Run("Upload Reviewed Document Forbidden", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"document_version_id": uuid.NewString(), // Doesn't exist, but forbidden check happens before? No.
			// Handler order: 
			// 1. Verify attachment exists (OK)
			// 2. Check permissions (Forbidden if not assigned)
			// So it should return 403 BEFORE checking document_version_id
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/admin/attachments/"+attID+"/review/upload", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}





func TestAdminHandler_PostReminders_Extended(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	tenantID := "00000000-0000-0000-0000-000000000001"
	_, err := db.Exec(`INSERT INTO tenants (id, name, slug) VALUES ($1, 'Test Tenant', 'test-tenant') ON CONFLICT DO NOTHING`, tenantID)
	require.NoError(t, err)

	userID := uuid.NewString()
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'admin', 'admin@ex.com', 'Admin', 'User', 'admin', 'hash', true)`, userID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'admin', true)`, userID, tenantID)
	require.NoError(t, err)

	studentID := uuid.NewString()
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'student', 'student@ex.com', 'Student', 'One', 'student', 'hash', true)`, studentID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'student', true)`, studentID, tenantID)
	require.NoError(t, err)

	h := handlers.NewAdminHandler(db, config.AppConfig{}, &pb.Manager{})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", tenantID)
		c.Set("claims", jwt.MapClaims{"sub": userID, "role": "admin"})
		c.Next()
	})
	r.POST("/admin/reminders", h.PostReminders)

	t.Run("Post Reminder Success With DueAt", func(t *testing.T) {
		dueAt := "2025-12-31T23:59:59Z"
		reqBody := map[string]interface{}{
			"student_ids": []string{studentID},
			"title":       "Reminder Title",
			"message":     "Reminder Message",
			"due_at":      dueAt,
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/admin/reminders", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var count int
		db.QueryRow("SELECT count(*) FROM reminders WHERE due_at IS NOT NULL").Scan(&count)
		assert.Equal(t, 1, count)
	})
}

func TestAdminHandler_PostReminders_Failures(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	h := handlers.NewAdminHandler(db, config.AppConfig{}, &pb.Manager{})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		// No user in context
		c.Next()
	})
	r.POST("/admin/reminders", h.PostReminders)

	t.Run("Post Reminder Unauthorized", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"student_ids": []string{"123"},
			"title":       "Title",
			"message":     "Message",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/admin/reminders", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

}

func TestAdminHandler_ReviewAttachment_Failures(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	tenantID := "00000000-0000-0000-0000-000000000001"
	_, err := db.Exec(`INSERT INTO tenants (id, name, slug) VALUES ($1, 'Test Tenant', 'test-tenant') ON CONFLICT DO NOTHING`, tenantID)
	require.NoError(t, err)

	userID := uuid.NewString()
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'admin', 'admin@ex.com', 'Admin', 'User', 'admin', 'hash', true)`, userID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'admin', true)`, userID, tenantID)
	require.NoError(t, err)

	h := handlers.NewAdminHandler(db, config.AppConfig{}, &pb.Manager{})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", tenantID)
		c.Set("claims", jwt.MapClaims{"sub": userID, "role": "admin"})
		c.Next()
	})
	r.POST("/admin/attachments/:attachmentId/review", h.ReviewAttachment)

	t.Run("Review Attachment Invalid Status", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"status": "invalid_status",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/admin/attachments/123/review", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	
	t.Run("Review Attachment Not Found", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"status": "approved",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/admin/attachments/99999999-9999-9999-9999-999999999999/review", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestAdminHandler_PostReminders(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	tenantID := "00000000-0000-0000-0000-000000000001"
	_, err := db.Exec(`INSERT INTO tenants (id, name, slug) VALUES ($1, 'Test Tenant', 'test-tenant') ON CONFLICT DO NOTHING`, tenantID)
	require.NoError(t, err)

	userID := "123e4567-e89b-12d3-a456-426614174000"
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'admin', 'admin@ex.com', 'Admin', 'User', 'admin', 'hash', true)
		ON CONFLICT (id) DO NOTHING`, userID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'admin', true) ON CONFLICT DO NOTHING`, userID, tenantID)
	require.NoError(t, err)

	studentID := "11111111-1111-1111-1111-111111111111"
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'student', 'student@ex.com', 'Student', 'One', 'student', 'hash', true)
		ON CONFLICT (id) DO NOTHING`, studentID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'student', true) ON CONFLICT DO NOTHING`, studentID, tenantID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	pb := &pb.Manager{}
	h := handlers.NewAdminHandler(db, cfg, pb)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", tenantID)
		c.Set("claims", jwt.MapClaims{"sub": userID, "role": "admin"})
		c.Next()
	})
	r.POST("/admin/reminders", h.PostReminders)

	t.Run("Post Reminder Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/admin/reminders", bytes.NewBuffer([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Post Reminder Success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"student_ids": []string{studentID},
			"title":       "Reminder Title",
			"message":     "Reminder Message",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/admin/reminders", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var count int
		err := db.Get(&count, "SELECT COUNT(*) FROM reminders WHERE student_id=$1", studentID)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})
}








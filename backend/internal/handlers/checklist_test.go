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
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChecklistHandler_ListModules(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Seed modules
	_, err := db.Exec(`INSERT INTO checklist_modules (id, code, title, sort_order) VALUES ('11111111-1111-1111-1111-111111111111', 'I', 'Module 1', 1)`)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	repo := repository.NewSQLChecklistRepository(db)
	svc := services.NewChecklistService(repo)
	h := handlers.NewChecklistHandler(svc, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/checklist/modules", h.ListModules)

	req, _ := http.NewRequest("GET", "/checklist/modules", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp, 1)
	assert.Equal(t, "I", resp[0]["code"])
}

func TestChecklistHandler_ListStepsByModule(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Seed modules and steps
	_, err := db.Exec(`INSERT INTO checklist_modules (id, code, title, sort_order) VALUES ('11111111-1111-1111-1111-111111111111', 'I', 'Module 1', 1)`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO checklist_steps (id, module_id, code, title, sort_order, requires_upload) 
		VALUES ('22222222-2222-2222-2222-222222222222', '11111111-1111-1111-1111-111111111111', '1.1', 'Step 1', 1, false)`)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	repo := repository.NewSQLChecklistRepository(db)
	svc := services.NewChecklistService(repo)
	h := handlers.NewChecklistHandler(svc, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/checklist/steps", h.ListStepsByModule)

	req, _ := http.NewRequest("GET", "/checklist/steps?module=I", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp, 1)
	assert.Equal(t, "1.1", resp[0]["code"])
}

func TestChecklistHandler_ListStudentSteps(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "10000000-0000-0000-0000-000000000001"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'student1', 's1@ex.com', 'Student', 'One', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO checklist_modules (id, code, title, sort_order) VALUES ('11111111-1111-1111-1111-111111111111', 'I', 'Module 1', 1)`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO checklist_steps (id, module_id, code, title, sort_order, requires_upload) 
		VALUES ('22222222-2222-2222-2222-222222222222', '11111111-1111-1111-1111-111111111111', '1.1', 'Step 1', 1, false)`)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO student_steps (user_id, step_id, status) VALUES ($1, '22222222-2222-2222-2222-222222222222', 'done')`, userID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	repo := repository.NewSQLChecklistRepository(db)
	svc := services.NewChecklistService(repo)
	h := handlers.NewChecklistHandler(svc, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/checklist/students/:id/steps", h.ListStudentSteps)

	req, _ := http.NewRequest("GET", "/checklist/students/"+userID+"/steps", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp, 1)
	assert.Equal(t, "22222222-2222-2222-2222-222222222222", resp[0]["step_id"])
	assert.Equal(t, "done", resp[0]["status"])
}

func TestChecklistHandler_UpdateStudentStep(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "10000000-0000-0000-0000-000000000001"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'student1', 's1@ex.com', 'Student', 'One', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO checklist_modules (id, code, title, sort_order) VALUES ('11111111-1111-1111-1111-111111111111', 'I', 'Module 1', 1)`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO checklist_steps (id, module_id, code, title, sort_order, requires_upload) 
		VALUES ('22222222-2222-2222-2222-222222222222', '11111111-1111-1111-1111-111111111111', '1.1', 'Step 1', 1, false)`)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	repo := repository.NewSQLChecklistRepository(db)
	svc := services.NewChecklistService(repo)
	h := handlers.NewChecklistHandler(svc, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.PUT("/checklist/students/:id/steps/:stepId", h.UpdateStudentStep)

	reqBody := map[string]interface{}{"status": "submitted", "data": map[string]string{"foo": "bar"}}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/checklist/students/"+userID+"/steps/22222222-2222-2222-2222-222222222222", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var status string
	db.QueryRow("SELECT status FROM student_steps WHERE user_id=$1 AND step_id='22222222-2222-2222-2222-222222222222'", userID).Scan(&status)
	assert.Equal(t, "submitted", status)
}

func TestChecklistHandler_AdvisorInbox(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "10000000-0000-0000-0000-000000000001"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'student1', 's1@ex.com', 'Student', 'One', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO checklist_modules (id, code, title, sort_order) VALUES ('11111111-1111-1111-1111-111111111111', 'I', 'Module 1', 1)`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO checklist_steps (id, module_id, code, title, sort_order, requires_upload) 
		VALUES ('22222222-2222-2222-2222-222222222222', '11111111-1111-1111-1111-111111111111', '1.1', 'Step 1', 1, false)`)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO student_steps (user_id, step_id, status) VALUES ($1, '22222222-2222-2222-2222-222222222222', 'submitted')`, userID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	repo := repository.NewSQLChecklistRepository(db)
	svc := services.NewChecklistService(repo)
	h := handlers.NewChecklistHandler(svc, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/checklist/advisor/inbox", h.AdvisorInbox)

	req, _ := http.NewRequest("GET", "/checklist/advisor/inbox", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp, 1)
	assert.Equal(t, userID, resp[0]["student_id"])
}

func TestChecklistHandler_ApproveStep(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "10000000-0000-0000-0000-000000000001"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'student1', 's1@ex.com', 'Student', 'One', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO checklist_modules (id, code, title, sort_order) VALUES ('11111111-1111-1111-1111-111111111111', 'I', 'Module 1', 1)`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO checklist_steps (id, module_id, code, title, sort_order, requires_upload) 
		VALUES ('22222222-2222-2222-2222-222222222222', '11111111-1111-1111-1111-111111111111', '1.1', 'Step 1', 1, false)`)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO student_steps (user_id, step_id, status) VALUES ($1, '22222222-2222-2222-2222-222222222222', 'submitted')`, userID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	repo := repository.NewSQLChecklistRepository(db)
	svc := services.NewChecklistService(repo)
	h := handlers.NewChecklistHandler(svc, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	// Create tenant and document for comment attachment
	tenantID := "00000000-0000-0000-0000-000000000001"
	_, err = db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type, is_active) 
		VALUES ($1, 'test', 'Test Tenant', 'university', true) ON CONFLICT DO NOTHING`, tenantID)
	require.NoError(t, err)
	
	_, err = db.Exec(`INSERT INTO documents (id, user_id, title, kind, tenant_id, created_at, updated_at) 
		VALUES ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', $1, 'Thesis', 'dissertation', $2, now(), now())`, userID, tenantID)
	require.NoError(t, err)

	// Add middleware to set tenant_id in context
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", tenantID)
		c.Next()
	})

	r.POST("/checklist/students/:id/steps/:stepId/approve", h.ApproveStep)

	req, _ := http.NewRequest("POST", "/checklist/students/"+userID+"/steps/22222222-2222-2222-2222-222222222222/approve", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var status string
	db.QueryRow("SELECT status FROM student_steps WHERE user_id=$1 AND step_id='22222222-2222-2222-2222-222222222222'", userID).Scan(&status)
	assert.Equal(t, "done", status)
}

func TestChecklistHandler_ReturnStep(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "10000000-0000-0000-0000-000000000001"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'student1', 's1@ex.com', 'Student', 'One', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO checklist_modules (id, code, title, sort_order) VALUES ('11111111-1111-1111-1111-111111111111', 'I', 'Module 1', 1)`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO checklist_steps (id, module_id, code, title, sort_order, requires_upload) 
		VALUES ('22222222-2222-2222-2222-222222222222', '11111111-1111-1111-1111-111111111111', '1.1', 'Step 1', 1, false)`)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO student_steps (user_id, step_id, status) VALUES ($1, '22222222-2222-2222-2222-222222222222', 'submitted')`, userID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	repo := repository.NewSQLChecklistRepository(db)
	svc := services.NewChecklistService(repo)
	h := handlers.NewChecklistHandler(svc, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	// Create tenant and document for comment attachment
	tenantID := "00000000-0000-0000-0000-000000000001"
	_, err = db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type, is_active) 
		VALUES ($1, 'test', 'Test Tenant', 'university', true) ON CONFLICT DO NOTHING`, tenantID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO documents (id, user_id, title, kind, tenant_id, created_at, updated_at) 
		VALUES ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', $1, 'Thesis', 'dissertation', $2, now(), now())`, userID, tenantID)
	require.NoError(t, err)

	// Add middleware to set tenant_id in context
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", tenantID)
		c.Next()
	})

	r.POST("/checklist/students/:id/steps/:stepId/return", h.ReturnStep)

	reqBody := map[string]string{"comment": "Please fix this"}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/checklist/students/"+userID+"/steps/22222222-2222-2222-2222-222222222222/return", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var status string
	db.QueryRow("SELECT status FROM student_steps WHERE user_id=$1 AND step_id='22222222-2222-2222-2222-222222222222'", userID).Scan(&status)
	assert.Equal(t, "needs_changes", status)
}

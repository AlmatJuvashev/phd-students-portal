package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProgramBuilderHandler_CreateNode(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Seed Tenant first (FK)
	progID := "11111111-1111-1111-1111-111111111111"
	tenantID := "22222222-2222-2222-2222-222222222222"
	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type, is_active) VALUES 
		($1, 'test-builder', 'Test Builder Tenant', 'university', true)`, tenantID)
	require.NoError(t, err)

	// Seed Program
	_, err = db.Exec(`INSERT INTO programs (id, tenant_id, code, name, title, is_active) VALUES 
		($1, $2, 'P1', 'Legacy Name', '{"en": "Program 1"}', true)`, progID, tenantID)
	require.NoError(t, err)

	repo := repository.NewSQLCurriculumRepository(db)
	svc := services.NewProgramBuilderService(repo)
	h := handlers.NewProgramBuilderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/programs/:id/builder/nodes", h.CreateNode)
	r.GET("/programs/:id/builder/nodes", h.GetNodes)

	// Case 1: Create a valid formEntry node
	// Note: title, description, coordinates, config are JSONB fields stored as strings.
	// API clients must send these as JSON strings (pre-serialized).
	reqBody := map[string]interface{}{
		"node": map[string]interface{}{
			"type":        "formEntry",
			"slug":        "form-1",
			"title":       `{"en": "Form 1"}`,
			"description": `{"en": "Test description"}`,
			"module_key":  "M1",
			"coordinates": `{"x": 0, "y": 0}`,
		},
		"config": map[string]interface{}{
			"fields": []map[string]interface{}{
				{"key": "f1", "type": "text", "label": map[string]string{"en": "Field 1"}},
			},
		},
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/programs/"+progID+"/builder/nodes", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Debug: Print response body on failure
	if w.Code != http.StatusCreated {
		t.Logf("CreateNode response: %s", w.Body.String())
	}
	assert.Equal(t, http.StatusCreated, w.Code)
	
	// Verify it was created
	reqGet, _ := http.NewRequest("GET", "/programs/"+progID+"/builder/nodes", nil)
	wGet := httptest.NewRecorder()
	r.ServeHTTP(wGet, reqGet)
	assert.Equal(t, http.StatusOK, wGet.Code)
	assert.Contains(t, wGet.Body.String(), "form-1")
}

func TestProgramBuilderHandler_UpdateJourneyMap(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	progID := "11111111-1111-1111-1111-111111111111"
	tenantID := "22222222-2222-2222-2222-222222222222"
	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type, is_active) VALUES ($1, 't1', 'T1', 'univ', true)`, tenantID)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO programs (id, tenant_id, code, name, title, is_active) VALUES ($1, $2, 'P1', 'N1', '{}', true)`, progID, tenantID)
	require.NoError(t, err)
	// Create journey map entry (ensure draft)
	// Table is program_versions. columns: id, program_id, version, is_active, config. title defaults?
	_, err = db.Exec(`INSERT INTO program_versions (id, program_id, version, is_active, config, title) VALUES ($1, $1, 1, false, '{}', '{}')`, progID)
	require.NoError(t, err)

	repo := repository.NewSQLCurriculumRepository(db)
	svc := services.NewProgramBuilderService(repo)
	h := handlers.NewProgramBuilderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.PUT("/programs/:id/builder/map", h.UpdateJourneyMap)

	// Success Case
	body, _ := json.Marshal(map[string]interface{}{
		"map": map[string]interface{}{
			"version": "2",
			"is_active": true,
		},
	})
	req, _ := http.NewRequest("PUT", "/programs/"+progID+"/builder/map", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Logf("UpdateJourneyMap failed: %s", w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	// Error: No fields
	bodyEmpty, _ := json.Marshal(map[string]interface{}{"map": map[string]interface{}{}})
	reqEmpty, _ := http.NewRequest("PUT", "/programs/"+progID+"/builder/map", bytes.NewBuffer(bodyEmpty))
	wEmpty := httptest.NewRecorder()
	r.ServeHTTP(wEmpty, reqEmpty)
	assert.Equal(t, http.StatusBadRequest, wEmpty.Code)
}

func TestProgramBuilderHandler_UpdateNode(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	progID := "11111111-1111-1111-1111-111111111111"
	tenantID := "22222222-2222-2222-2222-222222222222"
	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type, is_active) VALUES ($1, 't1', 'T1', 'univ', true)`, tenantID)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO programs (id, tenant_id, code, name, title, is_active) VALUES ($1, $2, 'P1', 'N1', '{}', true)`, progID, tenantID)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO program_versions (id, program_id, version, is_active, config, title) VALUES ($1, $1, 1, false, '{}', '{}')`, progID)
	require.NoError(t, err)
	
	nodeID := "33333333-3333-3333-3333-333333333333"
	// program_version_node_definitions: id, program_version_id, slug, type, title, description, coordinates, config, prerequisites, module_key
	_, err = db.Exec(`INSERT INTO program_version_node_definitions (id, program_version_id, slug, type, title, description, coordinates, config, prerequisites, module_key) 
		VALUES ($1, $2, 'n1', 'step', '{}', '{}', '{}', '{}', '{}', 'm1')`, nodeID, progID)
	require.NoError(t, err)

	repo := repository.NewSQLCurriculumRepository(db)
	svc := services.NewProgramBuilderService(repo)
	h := handlers.NewProgramBuilderHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.PUT("/programs/:id/builder/nodes/:nodeId", h.UpdateNode)

	// Update slug
	body, _ := json.Marshal(map[string]interface{}{
		"node": map[string]interface{}{
			"slug": "n1-updated",
		},
	})
	req, _ := http.NewRequest("PUT", "/programs/"+progID+"/builder/nodes/"+nodeID, bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Verify update
	var newSlug string
	err = db.QueryRow("SELECT slug FROM program_version_node_definitions WHERE id=$1", nodeID).Scan(&newSlug)
	require.NoError(t, err)
	assert.Equal(t, "n1-updated", newSlug)

	// Error: Not found
	reqNF, _ := http.NewRequest("PUT", "/programs/"+progID+"/builder/nodes/non-existent-id", bytes.NewBuffer(body))
	wNF := httptest.NewRecorder()
	r.ServeHTTP(wNF, reqNF)
	assert.Equal(t, http.StatusInternalServerError, wNF.Code) // Current implementation returns 500 on sql.ErrNoRows
}

func TestProgramBuilderHandler_CreateNode_Errors(t *testing.T) {
	// Re-using SetupTestDB pattern but simpler just for validation check if possible?
	// The handler checks validation BEFORE DB access in some cases, but EnsureDraftMap is called early.
	// So we need DB.
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLCurriculumRepository(db)
	svc := services.NewProgramBuilderService(repo)
	h := handlers.NewProgramBuilderHandler(svc)
	
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/programs/:id/builder/nodes", h.CreateNode)

	// Missing Slug
	body, _ := json.Marshal(map[string]interface{}{
		"node": map[string]interface{}{ "type": "step" }, // slug missing
	})
	req, _ := http.NewRequest("POST", "/programs/p1/builder/nodes", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Missing Type
	body2, _ := json.Marshal(map[string]interface{}{
		"node": map[string]interface{}{ "slug": "s1" }, // type missing
	})
	req2, _ := http.NewRequest("POST", "/programs/p1/builder/nodes", bytes.NewBuffer(body2))
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusBadRequest, w2.Code)
}

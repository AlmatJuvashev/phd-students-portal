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

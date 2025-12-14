package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDataIsolation_BetweenTenants(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// 1. Setup Tenants
	tenantA := "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	tenantB := "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
	_, err := db.Exec(`INSERT INTO tenants (id, name, slug) VALUES ($1, 'Tenant A', 'tenant-a'), ($2, 'Tenant B', 'tenant-b') ON CONFLICT DO NOTHING`, tenantA, tenantB)
	require.NoError(t, err)

	// 2. Setup Users for Tenant A
	adminA := uuid.NewString()
	studentA := uuid.NewString()
	_, err = db.Exec(`INSERT INTO users (id, username, email, role, first_name, last_name, password_hash, is_active) 
		VALUES ($1, 'adminA', 'adminA@ex.com', 'admin', 'Admin', 'A', 'hash', true), ($2, 'studentA', 'studentA@ex.com', 'student', 'Student', 'A', 'hash', true)`, adminA, studentA)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'admin', true), ($3, $2, 'student', true)`, adminA, tenantA, studentA)
	require.NoError(t, err)

	// 3. Setup Users for Tenant B
	adminB := uuid.NewString()
	studentB := uuid.NewString()
	_, err = db.Exec(`INSERT INTO users (id, username, email, role, first_name, last_name, password_hash, is_active) 
		VALUES ($1, 'adminB', 'adminB@ex.com', 'admin', 'Admin', 'B', 'hash', true), ($2, 'studentB', 'studentB@ex.com', 'student', 'Student', 'B', 'hash', true)`, adminB, studentB)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'admin', true), ($3, $2, 'student', true)`, adminB, tenantB, studentB)
	require.NoError(t, err)

	// 4. Create Data in Tenant A
	// Document
	// 5. Create Data in Tenant B
	// Document
	docB := uuid.NewString()
	db.Exec(`INSERT INTO documents (id, user_id, title, kind, tenant_id) VALUES ($1, $2, 'Doc B', 'other', $3)`, docB, studentB, tenantB)

	// Setup Handler
	// We need a handler that accesses these resources. AdminHandler's GetStudentDetails or MonitorStudents or Documents check.
	// Using AdminHandler for testing isolation on list endpoints.
	h := handlers.NewAdminHandler(db, config.AppConfig{}, &playbook.Manager{})

	gin.SetMode(gin.TestMode)

	// Test Case 1: Admin A listing students - Should see Student A, NOT Student B
	t.Run("Admin A sees only Student A", func(t *testing.T) {
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("claims", jwt.MapClaims{"sub": adminA, "role": "admin"})
			c.Set("tenant_id", tenantA)
			c.Next()
		})
		r.GET("/admin/students", h.MonitorStudents)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/admin/students?role=student", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var students []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &students)
		
		// Verify
		foundStudentA := false
		foundStudentB := false
		for _, s := range students {
			id := s["id"].(string)
			if id == studentA { foundStudentA = true }
			if id == studentB { foundStudentB = true }
		}
		assert.True(t, foundStudentA, "Should find student A")
		assert.False(t, foundStudentB, "Should NOT find student B")
	})

	// Test Case 2: Admin B list students - Should see Student B, NOT Student A
	t.Run("Admin B sees only Student B", func(t *testing.T) {
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("claims", jwt.MapClaims{"sub": adminB, "role": "admin"})
			c.Set("tenant_id", tenantB)
			c.Next()
		})
		r.GET("/admin/students", h.MonitorStudents)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/admin/students?role=student", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var students []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &students)
		
		// Verify
		foundStudentA := false
		foundStudentB := false
		for _, s := range students {
			id := s["id"].(string)
			if id == studentA { foundStudentA = true }
			if id == studentB { foundStudentB = true }
		}
		assert.False(t, foundStudentA, "Should NOT find student A")
		assert.True(t, foundStudentB, "Should find student B")
	})

	// Test Case 3: Admin A trying to access Student B details -> 404 (Not Found in tenant)
	t.Run("Admin A accesses Student B details -> 404", func(t *testing.T) {
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("claims", jwt.MapClaims{"sub": adminA, "role": "admin"})
			c.Set("tenant_id", tenantA)
			c.Next()
		})
		r.GET("/admin/students/:id", h.GetStudentDetails)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/admin/students/"+studentB, nil)
		r.ServeHTTP(w, req)

		// Should receive 404 because student B is not in tenant A
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

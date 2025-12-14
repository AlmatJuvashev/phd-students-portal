package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDictionaryHandler_Programs(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	h := handlers.NewDictionaryHandler(db)

	gin.SetMode(gin.TestMode)
	tenantID := "00000000-0000-0000-0000-000000000001"
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", tenantID)
		c.Next()
	})
	r.POST("/dictionaries/programs", h.CreateProgram)
	r.GET("/dictionaries/programs", h.ListPrograms)
	r.PATCH("/dictionaries/programs/:id", h.UpdateProgram)
	r.DELETE("/dictionaries/programs/:id", h.DeleteProgram)

	var programID string

	t.Run("Create Program", func(t *testing.T) {
		reqBody := map[string]string{
			"name": "Computer Science",
			"code": "CS",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/dictionaries/programs", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]string
		json.Unmarshal(w.Body.Bytes(), &resp)
		programID = resp["id"]
	})

	t.Run("Create Program Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/dictionaries/programs", bytes.NewBuffer([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("List Programs", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/dictionaries/programs", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Len(t, resp, 1)
		assert.Equal(t, "Computer Science", resp[0]["name"])
	})

	t.Run("Update Program", func(t *testing.T) {
		reqBody := map[string]string{
			"name": "CS Updated",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PATCH", "/dictionaries/programs/"+programID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var name string
		db.QueryRow("SELECT name FROM programs WHERE id=$1", programID).Scan(&name)
		assert.Equal(t, "CS Updated", name)
	})

	t.Run("Delete Program", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/dictionaries/programs/"+programID, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var isActive bool
		db.QueryRow("SELECT is_active FROM programs WHERE id=$1", programID).Scan(&isActive)
		assert.False(t, isActive)
	})
}

func TestDictionaryHandler_Departments(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	h := handlers.NewDictionaryHandler(db)

	gin.SetMode(gin.TestMode)
	tenantID := "00000000-0000-0000-0000-000000000001"
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", tenantID)
		c.Next()
	})
	r.POST("/dictionaries/departments", h.CreateDepartment)
	r.GET("/dictionaries/departments", h.ListDepartments)
	r.PATCH("/dictionaries/departments/:id", h.UpdateDepartment)
	r.DELETE("/dictionaries/departments/:id", h.DeleteDepartment)

	var deptID string

	t.Run("Create Department", func(t *testing.T) {
		reqBody := map[string]string{
			"name": "Engineering",
			"code": "ENG",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/dictionaries/departments", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]string
		json.Unmarshal(w.Body.Bytes(), &resp)
		deptID = resp["id"]
	})

	t.Run("List Departments", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/dictionaries/departments", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Len(t, resp, 1)
		assert.Equal(t, "Engineering", resp[0]["name"])
	})

	t.Run("Update Department", func(t *testing.T) {
		reqBody := map[string]string{
			"name": "Engineering Updated",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PATCH", "/dictionaries/departments/"+deptID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var name string
		db.QueryRow("SELECT name FROM departments WHERE id=$1", deptID).Scan(&name)
		assert.Equal(t, "Engineering Updated", name)
	})

	t.Run("Delete Department", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/dictionaries/departments/"+deptID, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var isActive bool
		db.QueryRow("SELECT is_active FROM departments WHERE id=$1", deptID).Scan(&isActive)
		assert.False(t, isActive)
	})
}

func TestDictionaryHandler_Cohorts(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	h := handlers.NewDictionaryHandler(db)

	gin.SetMode(gin.TestMode)
	tenantID := "00000000-0000-0000-0000-000000000001"
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", tenantID)
		c.Next()
	})
	r.POST("/dictionaries/cohorts", h.CreateCohort)
	r.GET("/dictionaries/cohorts", h.ListCohorts)
	r.PATCH("/dictionaries/cohorts/:id", h.UpdateCohort)
	r.DELETE("/dictionaries/cohorts/:id", h.DeleteCohort)

	var cohortID string

	t.Run("Create Cohort", func(t *testing.T) {
		reqBody := map[string]string{
			"name": "2024-2025",
			"start_date": "2024-09-01",
			"end_date": "2025-06-30",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/dictionaries/cohorts", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]string
		json.Unmarshal(w.Body.Bytes(), &resp)
		cohortID = resp["id"]
	})

	t.Run("List Cohorts", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/dictionaries/cohorts", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Len(t, resp, 1)
		assert.Equal(t, "2024-2025", resp[0]["name"])
	})

	t.Run("Update Cohort", func(t *testing.T) {
		reqBody := map[string]string{
			"name": "2024-2025 Updated",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PATCH", "/dictionaries/cohorts/"+cohortID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		// Verify
		var name string
		db.QueryRow("SELECT name FROM cohorts WHERE id=$1", cohortID).Scan(&name)
		assert.Equal(t, "2024-2025 Updated", name)
	})

	t.Run("Delete Cohort", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/dictionaries/cohorts/"+cohortID, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		// Verify inactive
		var isActive bool
		db.QueryRow("SELECT is_active FROM cohorts WHERE id=$1", cohortID).Scan(&isActive)
		assert.False(t, isActive)
	})
}

func TestDictionaryHandler_Specialties(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	h := handlers.NewDictionaryHandler(db)

	gin.SetMode(gin.TestMode)
	tenantID := "00000000-0000-0000-0000-000000000001"
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", tenantID)
		c.Next()
	})
	r.GET("/dictionaries/specialties", h.ListSpecialties)

	// Seed specialty
	// Seed specialty
	_, err := db.Exec(`INSERT INTO specialties (id, name, code, tenant_id) VALUES ('10000000-0000-0000-0000-000000000001', 'Computer Science', 'CS101', $1)`, tenantID)
	require.NoError(t, err)

	t.Run("List Specialties", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/dictionaries/specialties", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Len(t, resp, 1)
		assert.Equal(t, "Computer Science", resp[0]["name"])
	})
}

func TestDictionaryHandler_SpecialtyCRUD(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	h := handlers.NewDictionaryHandler(db)

	gin.SetMode(gin.TestMode)
	tenantID := "00000000-0000-0000-0000-000000000001"
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", tenantID)
		c.Next()
	})
	r.POST("/dictionaries/specialties", h.CreateSpecialty)
	r.PUT("/dictionaries/specialties/:id", h.UpdateSpecialty)
	r.DELETE("/dictionaries/specialties/:id", h.DeleteSpecialty)

	var specID string

	t.Run("Create Specialty", func(t *testing.T) {
		reqBody := map[string]string{"name": "New Spec", "code": "NS001"}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/dictionaries/specialties", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		specID = resp["id"].(string)
	})

	t.Run("Create Specialty Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/dictionaries/specialties", bytes.NewBuffer([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Update Specialty", func(t *testing.T) {
		reqBody := map[string]string{"name": "Updated Spec", "code": "NS001"}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PUT", "/dictionaries/specialties/"+specID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete Specialty", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/dictionaries/specialties/"+specID, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

package handlers_test

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSuperadminTenantsHandler_MockFailures(t *testing.T) {
	gin.SetMode(gin.TestMode)

	setup := func(t *testing.T) (*gin.Engine, *handlers.SuperadminTenantsHandler, sqlmock.Sqlmock) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("unexpected error '%s' when opening a stub database connection", err)
		}
		
		sqlxDB := sqlx.NewDb(db, "sqlmock")
		
		tenantRepo := repository.NewSQLTenantRepository(sqlxDB)
		adminRepo := repository.NewSQLSuperAdminRepository(sqlxDB)
		tenantSvc := services.NewTenantService(tenantRepo)
		adminSvc := services.NewSuperAdminService(adminRepo)

		h := handlers.NewSuperadminTenantsHandler(tenantSvc, adminSvc, config.AppConfig{})
		r := gin.New()
		
		// Auth middleware mock if needed, but here we can just pass context via recorder or simple middleware
		r.Use(func(c *gin.Context) {
			c.Set("userID", "admin-id")
			c.Next()
		})
		
		return r, h, mock
	}

	t.Run("ListTenants DB Error", func(t *testing.T) {
		r, h, mock := setup(t)

		// ListAllWithStats calls List
		mock.ExpectQuery("SELECT .* FROM tenants").
			WillReturnError(errors.New("db error"))

		r.GET("/tenants", h.ListTenants)
		req, _ := http.NewRequest("GET", "/tenants", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "failed to fetch tenants")
	})

	t.Run("GetTenant DB Error", func(t *testing.T) {
		r, h, mock := setup(t)

		// GetWithStats calls GetByID
		mock.ExpectQuery("SELECT .* FROM tenants WHERE id").
			WithArgs("t1").
			WillReturnError(errors.New("db error"))

		r.GET("/tenants/:id", h.GetTenant)
		req, _ := http.NewRequest("GET", "/tenants/t1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Handler returns 404 on error? Let's check logic:
		// if err != nil { c.JSON(http.StatusNotFound... "tenant not found") }
		// Wait, implementation: 
		// if err != nil { c.JSON(http.StatusNotFound, ... "tenant not found") }
		// So it treats DB error as Not Found?
		assert.Equal(t, http.StatusNotFound, w.Code) 
		assert.Contains(t, w.Body.String(), "tenant not found")
	})

	t.Run("UpdateTenant DB Error", func(t *testing.T) {
		r, h, mock := setup(t)

		mock.ExpectExec("UPDATE tenants SET").
			WillReturnError(errors.New("db error"))

		r.PUT("/tenants/:id", h.UpdateTenant)
		body := `{"name":"New Name"}`
		req, _ := http.NewRequest("PUT", "/tenants/t1", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "failed to update tenant")
	})

	t.Run("UploadLogo DB Error", func(t *testing.T) {
		r, h, mock := setup(t)

		mock.ExpectExec("UPDATE tenants SET logo_url").
			WillReturnError(errors.New("db error"))

		r.POST("/tenants/:id/logo", h.UploadLogo)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		header := make(map[string][]string)
		header["Content-Disposition"] = []string{`form-data; name="logo"; filename="logo.png"`}
		header["Content-Type"] = []string{"image/png"}
		part, _ := writer.CreatePart(header)
		part.Write([]byte("fake"))
		writer.Close()

		req, _ := http.NewRequest("POST", "/tenants/t1/logo", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "failed to update logo")
	})
}

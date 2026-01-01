package handlers_test

import (
	"errors"
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

func TestChecklistHandler_ServiceErrors(t *testing.T) {
	// Setup sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	// Setup stack
	repo := repository.NewSQLChecklistRepository(sqlxDB)
	svc := services.NewChecklistService(repo)
	h := handlers.NewChecklistHandler(svc, config.AppConfig{})

	gin.SetMode(gin.TestMode)

	t.Run("ListModules DB Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, code, title, sort_order FROM checklist_modules").
			WillReturnError(errors.New("db error"))

		r := gin.New()
		r.GET("/checklist/modules", h.ListModules)

		req, _ := http.NewRequest("GET", "/checklist/modules", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "failed to list modules")
	})

	t.Run("ListStepsByModule DB Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, code, title, requires_upload, sort_order FROM checklist_steps").
			WithArgs("mod1").
			WillReturnError(errors.New("db error"))

		r := gin.New()
		r.GET("/checklist/steps", h.ListStepsByModule)

		req, _ := http.NewRequest("GET", "/checklist/steps?module=mod1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "failed to list steps")
	})

	t.Run("ListStudentSteps DB Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT step_id, status FROM student_steps").
			WithArgs("u1").
			WillReturnError(errors.New("db error"))

		r := gin.New()
		r.GET("/checklist/students/:id/steps", h.ListStudentSteps)

		req, _ := http.NewRequest("GET", "/checklist/students/u1/steps", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "failed to list student steps")
	})

	t.Run("UpdateStudentStep DB Error", func(t *testing.T) {
		// Expect Exec for Update
		mock.ExpectExec("INSERT INTO student_steps").
			WithArgs("u1", "s1", "submitted", sqlmock.AnyArg()).
			WillReturnError(errors.New("db error"))

		r := gin.New()
		r.PUT("/checklist/students/:id/steps/:stepId", h.UpdateStudentStep)

		body := `{"status":"submitted", "data":{}}`
		req, _ := http.NewRequest("PUT", "/checklist/students/u1/steps/s1", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("AdvisorInbox DB Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT ss.user_id").
			WillReturnError(errors.New("db error"))

		r := gin.New()
		r.GET("/checklist/advisor/inbox", h.AdvisorInbox)

		req, _ := http.NewRequest("GET", "/checklist/advisor/inbox", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

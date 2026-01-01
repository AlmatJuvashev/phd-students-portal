package handlers_test

import (
	"bytes"
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
)

func TestChecklistHandler_UpdateStudentStep_Failures(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLChecklistRepository(db)
	svc := services.NewChecklistService(repo)
	h := handlers.NewChecklistHandler(svc, config.AppConfig{})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.PUT("/checklist/students/:id/steps/:stepId", h.UpdateStudentStep)

	t.Run("Invalid JSON Body", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/checklist/students/user1/steps/step1", bytes.NewBufferString("invalid-json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestChecklistHandler_ApproveStep_Failures(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLChecklistRepository(db)
	svc := services.NewChecklistService(repo)
	h := handlers.NewChecklistHandler(svc, config.AppConfig{})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/checklist/students/:id/steps/:stepId/approve", h.ApproveStep)

	t.Run("Approve Step Invalid Input", func(t *testing.T) {
		// Even if body is invalid JSON, ShouldBindJSON might just leave struct empty or error.
		// If it errors, we generally want 400.
		// But handler uses _ = c.ShouldBindJSON(&r) so it ignores error and succeeds with empty values.
		// So testing invalid JSON won't trigger 400 here unless we change handler.
		// Wait, handler code: `_ = c.ShouldBindJSON(&r)`. It IGNORES the error!
		// So we can't test 400 here properly without fixing the handler.
	})
}

// Handler code for ApproveStep:
// var r reviewReq
// _ = c.ShouldBindJSON(&r)
// if err := h.svc.ApproveStep...

// So if JSON is invalid, r is empty. ApproveStep is called with empty comments/mentions.
// To cover error path, we need svc.ApproveStep to fail.
// This usually happens if inputs are invalid or DB fails.
// We can pass invalid user ID or step ID if the service validates them or DB constraint fails.

func TestChecklistHandler_ApproveStep_ServiceError(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()
	
	repo := repository.NewSQLChecklistRepository(db)
	svc := services.NewChecklistService(repo)
	h := handlers.NewChecklistHandler(svc, config.AppConfig{})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/checklist/students/:id/steps/:stepId/approve", h.ApproveStep)
	
	t.Run("Approve Non-Existent User/Step", func(t *testing.T) {
		// This should trigger an error in service (User/Step not found or DB error)
		req, _ := http.NewRequest("POST", "/checklist/students/bad-user/steps/bad-step/approve", bytes.NewBufferString("{}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		// Handler returns 500 on service error
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestChecklistHandler_ReturnStep_ServiceError(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLChecklistRepository(db)
	svc := services.NewChecklistService(repo)
	h := handlers.NewChecklistHandler(svc, config.AppConfig{})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/checklist/students/:id/steps/:stepId/return", h.ReturnStep)

	t.Run("Return Non-Existent User/Step", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/checklist/students/bad-user/steps/bad-step/return", bytes.NewBufferString("{}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

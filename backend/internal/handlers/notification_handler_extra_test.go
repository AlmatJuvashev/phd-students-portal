package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNotificationHandler_Failures(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLNotificationRepository(db)
	svc := services.NewNotificationService(repo)
	h := handlers.NewNotificationHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", "user1")
		c.Next()
	})
	r.PUT("/notifications/:id/read", h.MarkAsRead)

	t.Run("Mark As Read Invalid UUID", func(t *testing.T) {
		// Passing an invalid UUID should trigger DB error in repository -> service error -> 500
		req, _ := http.NewRequest("PUT", "/notifications/invalid-uuid/read", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Expect 500 because handler blindly returns 500 on error
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})
}

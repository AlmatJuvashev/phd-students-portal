package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMeHandler_Me(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "40000000-aaaa-4000-4000-400000000000"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'meuser', 'me@ex.com', 'Me', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	// Pass nil for redis client
	h := handlers.NewMeHandler(db, cfg, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", map[string]interface{}{"sub": userID})
		c.Next()
	})
	r.GET("/auth/me", h.Me)

	t.Run("Get Me Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/auth/me", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "meuser", resp["username"])
		assert.Equal(t, "Me", resp["first_name"])
	})

	t.Run("Get Me Not Found", func(t *testing.T) {
		// Use a different user ID that doesn't exist
		r2 := gin.New()
		r2.Use(func(c *gin.Context) {
			c.Set("claims", map[string]interface{}{"sub": "99999999-aaaa-9999-9999-999999999999"})
			c.Next()
		})
		r2.GET("/auth/me", h.Me)

		req, _ := http.NewRequest("GET", "/auth/me", nil)
		w := httptest.NewRecorder()
		r2.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

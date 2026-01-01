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
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestChatHandler_ServiceErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Helper to setup isolated environment
	setup := func(t *testing.T) (*gin.Engine, *handlers.ChatHandler, sqlmock.Sqlmock) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		// Note: we don't defer db.Close() here easily because it closes check. 
		// Ideally we rely on garbage collection or cleaner structure, but for small unit tests it's okay-ish 
		// or we return a cleanup func.
		
		sqlxDB := sqlx.NewDb(db, "sqlmock")
		repo := repository.NewSQLChatRepository(sqlxDB)
		svc := services.NewChatService(repo, nil, config.AppConfig{})
		h := handlers.NewChatHandler(svc, config.AppConfig{})

		r := gin.New()
		
		return r, h, mock
	}

	t.Run("CreateRoom DB Error", func(t *testing.T) {
		r, h, mock := setup(t)
		
		mock.ExpectQuery("INSERT INTO chat_rooms").
			WillReturnError(errors.New("db error"))

		r.Use(func(c *gin.Context) {
			c.Set("claims", jwt.MapClaims{"sub": "u1"})
			c.Set("tenant_id", "t1")
			c.Next()
		})
		r.POST("/chat/rooms", h.CreateRoom)

		body := `{"name":"Room1", "type":"cohort"}`
		req, _ := http.NewRequest("POST", "/chat/rooms", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "failed to create room")
	})

	t.Run("ListRooms DB Error", func(t *testing.T) {
		r, h, mock := setup(t)

		mock.ExpectQuery("SELECT .* FROM chat_rooms").
			WillReturnError(errors.New("db error"))

		r.Use(func(c *gin.Context) {
			c.Set("claims", jwt.MapClaims{"sub": "u1"})
			c.Set("tenant_id", "t1")
			c.Next()
		})
		r.GET("/chat/rooms", h.ListRooms)

		req, _ := http.NewRequest("GET", "/chat/rooms", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "failed to list rooms")
	})

	t.Run("CreateMessage Membership Check Error", func(t *testing.T) {
		r, h, mock := setup(t)

		// IsMember query fails
		mock.ExpectQuery("SELECT COUNT").
			WillReturnError(errors.New("db error"))

		r.Use(func(c *gin.Context) {
			c.Set("claims", jwt.MapClaims{"sub": "u1"})
			c.Next()
		})
		r.POST("/chat/rooms/:roomId/messages", h.CreateMessage)

		req, _ := http.NewRequest("POST", "/chat/rooms/r1/messages", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "membership check failed")
	})

	t.Run("AddMember DB Error", func(t *testing.T) {
		r, h, mock := setup(t)

		// Service calls repo.AddMember directly which does INSERT...SELECT
		mock.ExpectExec("INSERT INTO chat_room_members").WillReturnError(errors.New("db error"))

		r.POST("/chat/rooms/:roomId/members", h.AddMember)

		body := `{"user_id":"u2","role_in_room":"member"}`
		req, _ := http.NewRequest("POST", "/chat/rooms/r1/members", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "failed to add member")
	})

	t.Run("AddRoomMembersBatch DB Error", func(t *testing.T) {
		r, h, mock := setup(t)

		// Service likely iterates and calls AddMember (INSERT...)
		// We expect at least one failure
		mock.ExpectExec("INSERT INTO chat_room_members").
			WillReturnError(errors.New("db error"))

		r.POST("/chat/rooms/:roomId/members/batch", h.AddRoomMembersBatch)

		body := `{"user_ids":["u2","u3"]}`
		req, _ := http.NewRequest("POST", "/chat/rooms/r1/members/batch", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		// Since DB failed, added_count might be 0
		assert.Contains(t, w.Body.String(), "added_count")
	})
}

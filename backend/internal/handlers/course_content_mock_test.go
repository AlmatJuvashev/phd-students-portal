package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCourseContentHandler_Module(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewSQLCourseContentRepository(sqlxDB)
	svc := services.NewCourseContentService(repo)
	handler := NewCourseContentHandler(svc)

	// Create Module
	mock.ExpectQuery(`INSERT INTO course_modules`).
		WithArgs("c1", "M1", 1, true).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow("m1", time.Now(), time.Now()))

	m := models.CourseModule{CourseID: "c1", Title: "M1", Order: 1, IsActive: true}
	body, _ := json.Marshal(m)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/content/modules", bytes.NewBuffer(body))
	handler.CreateModule(c)
	assert.Equal(t, http.StatusCreated, w.Code)

	// List Modules
	mock.ExpectQuery(`SELECT \* FROM course_modules WHERE course_id=\$1`).
		WithArgs("c1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow("m1", "M1"))
	
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request, _ = http.NewRequest("GET", "/api/content/modules?course_id=c1", nil)
	handler.ListModules(c2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestCourseContentHandler_Activity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewSQLCourseContentRepository(sqlxDB)
	svc := services.NewCourseContentService(repo)
	handler := NewCourseContentHandler(svc)

	// Create Activity
	mock.ExpectQuery(`INSERT INTO course_activities`).
		WithArgs("l1", "text", "A1", 1, 0, false, "{}").
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow("a1", time.Now(), time.Now()))

	a := models.CourseActivity{LessonID: "l1", Type: "text", Title: "A1", Order: 1}
	body, _ := json.Marshal(a)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/content/activities", bytes.NewBuffer(body))
	handler.CreateActivity(c)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCourseContentHandler_FullCRUD(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewSQLCourseContentRepository(sqlxDB)
	svc := services.NewCourseContentService(repo)
	handler := NewCourseContentHandler(svc)

	// We'll test one Update flow (e.g., Module) to verify binder & service call
	mock.ExpectExec(`UPDATE course_modules`).
		WithArgs("M1 Updated", 2, true, "m1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	m := models.CourseModule{ID: "m1", Title: "M1 Updated", Order: 2, IsActive: true}
	body, _ := json.Marshal(m)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "m1"}}
	c.Request, _ = http.NewRequest("PUT", "/api/content/modules/m1", bytes.NewBuffer(body))
	handler.UpdateModule(c)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test Delete Flow
	mock.ExpectExec(`DELETE FROM course_modules`).
		WithArgs("m1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Params = gin.Params{{Key: "id", Value: "m1"}}
	c2.Request, _ = http.NewRequest("DELETE", "/api/content/modules/m1", nil)
	handler.DeleteModule(c2)
	assert.Equal(t, http.StatusOK, w2.Code)

	// Test Lesson Create
	mock.ExpectQuery(`INSERT INTO course_lessons`).
		WithArgs("m1", "L1", 1, true).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow("l1", time.Now(), time.Now()))
	
	l := models.CourseLesson{ModuleID: "m1", Title: "L1", Order: 1, IsActive: true}
	bodyL, _ := json.Marshal(l)
	w3 := httptest.NewRecorder()
	c3, _ := gin.CreateTestContext(w3)
	c3.Request, _ = http.NewRequest("POST", "/api/content/lessons", bytes.NewBuffer(bodyL))
	handler.CreateLesson(c3)
	assert.Equal(t, http.StatusCreated, w3.Code)
}

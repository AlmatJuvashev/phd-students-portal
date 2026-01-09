package handlers

import (
	"bytes"
	"database/sql"
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



func TestCurriculumHandler_CreateProgram_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewSQLCurriculumRepository(sqlxDB)
	svc := services.NewCurriculumService(repo)
	handler := NewCurriculumHandler(svc)

	// Mock DB Expectation for Repo call
	mock.ExpectQuery(`INSERT INTO programs`).
		WithArgs("tenant-1", "P1", "P1", `{"en":"Title"}`, "{\"en\":\"Desc\"}", 120, 36, true).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow("prog-1", time.Now(), time.Now()))

	// Request
	p := models.Program{
		Code: "P1",
		Title: `{"en":"Title"}`,
		Description: strPtr(`{"en":"Desc"}`),
		Credits: 120,
		DurationMonths: 36,
		IsActive: true,
	}
	body, _ := json.Marshal(p)
	req, _ := http.NewRequest("POST", "/api/curriculum/programs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	// Mock Middleware context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("tenant_id", "tenant-1") // Simulating middleware

	handler.CreateProgram(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCurriculumHandler_ListPrograms_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewSQLCurriculumRepository(sqlxDB)
	svc := services.NewCurriculumService(repo)
	handler := NewCurriculumHandler(svc)

	// Mock DB
	mock.ExpectQuery(`SELECT \* FROM programs WHERE tenant_id=\$1`).
		WithArgs("tenant-1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "code", "title", "tenant_id"}).
			AddRow("prog-1", "P1", `{"en":"Title"}`, "tenant-1"))

	// Request
	req, _ := http.NewRequest("GET", "/api/curriculum/programs", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("tenant_id", "tenant-1")

	handler.ListPrograms(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCurriculumHandler_GetProgram_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewSQLCurriculumRepository(sqlxDB)
	svc := services.NewCurriculumService(repo)
	handler := NewCurriculumHandler(svc)

	mock.ExpectQuery(`SELECT \* FROM programs WHERE id=\$1`).
		WithArgs("prog-1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "code"}).AddRow("prog-1", "P1"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/curriculum/programs/prog-1", nil)
	c.Params = gin.Params{{Key: "id", Value: "prog-1"}}
	
	handler.GetProgram(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCurriculumHandler_GetProgram_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewSQLCurriculumRepository(sqlxDB)
	svc := services.NewCurriculumService(repo)
	handler := NewCurriculumHandler(svc)

	mock.ExpectQuery(`SELECT \* FROM programs WHERE id=\$1`).
		WithArgs("prog-missing").
		WillReturnError(sql.ErrNoRows) // sqlx.Get returns error on no rows

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/curriculum/programs/prog-missing", nil)
	c.Params = gin.Params{{Key: "id", Value: "prog-missing"}}
	
	handler.GetProgram(c)
	// sqlx.Get returns error, so service returns error. 
	// To distinguish 404 vs 500, repo/service needs fine-tuning, but for now 500 is expected on error.
	// Actually, repo.GetProgram returns error. Service returns error. Handler checks err != nil -> 500.
	assert.Equal(t, http.StatusInternalServerError, w.Code) 
}

func TestCurriculumHandler_UpdateProgram_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewSQLCurriculumRepository(sqlxDB)
	svc := services.NewCurriculumService(repo)
	handler := NewCurriculumHandler(svc)

	p := models.Program{Code: "P2"}
	mock.ExpectExec(`UPDATE programs SET code=\$1`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	body, _ := json.Marshal(p)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/api/curriculum/programs/prog-1", bytes.NewBuffer(body))
	c.Params = gin.Params{{Key: "id", Value: "prog-1"}}
	
	handler.UpdateProgram(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCurriculumHandler_DeleteProgram_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewSQLCurriculumRepository(sqlxDB)
	svc := services.NewCurriculumService(repo)
	handler := NewCurriculumHandler(svc)

	mock.ExpectExec(`DELETE FROM programs WHERE id=\$1`).
		WithArgs("prog-1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/api/curriculum/programs/prog-1", nil)
	c.Params = gin.Params{{Key: "id", Value: "prog-1"}}
	
	handler.DeleteProgram(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCurriculumHandler_CreateCourse_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewSQLCurriculumRepository(sqlxDB)
	svc := services.NewCurriculumService(repo)
	handler := NewCurriculumHandler(svc)

	mock.ExpectQuery(`INSERT INTO courses`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow("c1", time.Now(), time.Now()))

	cObj := models.Course{Code: "C101", Title: "Intro"}
	body, _ := json.Marshal(cObj)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/curriculum/courses", bytes.NewBuffer(body))
	c.Set("tenant_id", "t1")
	
	handler.CreateCourse(c)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCurriculumHandler_ListCourses_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewSQLCurriculumRepository(sqlxDB)
	svc := services.NewCurriculumService(repo)
	handler := NewCurriculumHandler(svc)

	mock.ExpectQuery(`SELECT \* FROM courses WHERE tenant_id=\$1`).
		WithArgs("t1").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("c1"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/curriculum/courses", nil)
	c.Set("tenant_id", "t1")
	
	handler.ListCourses(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

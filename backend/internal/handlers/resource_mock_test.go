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

func TestResourceHandler_CreateBuilding_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewSQLResourceRepository(sqlxDB)
	svc := services.NewResourceService(repo)
	handler := NewResourceHandler(svc)

	mock.ExpectQuery(`INSERT INTO buildings`).
		WithArgs("t1", "Main", "123 St", "{}", true, "u1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow("b1", time.Now(), time.Now()))

	b := models.Building{Name: "Main", Address: "123 St", Description: "{}", IsActive: true}
	body, _ := json.Marshal(b)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/resources/buildings", bytes.NewBuffer(body))
	c.Set("tenant_id", "t1")
	c.Set("userID", "u1")

	handler.CreateBuilding(c)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestResourceHandler_ListBuildings_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewSQLResourceRepository(sqlxDB)
	svc := services.NewResourceService(repo)
	handler := NewResourceHandler(svc)

	mock.ExpectQuery(`SELECT \* FROM buildings WHERE tenant_id=\$1 AND deleted_at IS NULL`).
		WithArgs("t1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow("b1", "Main"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/resources/buildings", nil)
	c.Set("tenant_id", "t1")

	handler.ListBuildings(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestResourceHandler_CreateRoom_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewSQLResourceRepository(sqlxDB)
	svc := services.NewResourceService(repo)
	handler := NewResourceHandler(svc)

	buildingID := "00000000-0000-0000-0000-000000000001"

	// 9 args: building_id, name, capacity, floor, department_id, type, features, is_active, created_by
	mock.ExpectQuery(`INSERT INTO rooms`).
		WithArgs(buildingID, "101", 30, 1, nil, "lab", "[]", true, "u1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow("r1", time.Now(), time.Now()))

	r := models.Room{BuildingID: buildingID, Name: "101", Capacity: 30, Floor: 1, Type: "lab", Features: "[]", IsActive: true}
	body, _ := json.Marshal(r)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/resources/rooms", bytes.NewBuffer(body))
	c.Set("userID", "u1")

	handler.CreateRoom(c)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestResourceHandler_ListRooms_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewSQLResourceRepository(sqlxDB)
	svc := services.NewResourceService(repo)
	handler := NewResourceHandler(svc)

	buildingID := "00000000-0000-0000-0000-000000000001"

	mock.ExpectQuery(`SELECT[\s\S]*FROM rooms r[\s\S]*JOIN buildings b[\s\S]*WHERE b\.tenant_id = \$1[\s\S]*r\.building_id = \$2`).
		WithArgs("t1", buildingID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow("r1", "101"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/resources/rooms?building_id="+buildingID, nil)
	c.Set("tenant_id", "t1")

	handler.ListRooms(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestResourceHandler_GetUpdateDeleteBuilding(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewSQLResourceRepository(sqlxDB)
	svc := services.NewResourceService(repo)
	handler := NewResourceHandler(svc)

	// Get
	mock.ExpectQuery(`SELECT \* FROM buildings WHERE id=\$1 AND deleted_at IS NULL`).
		WithArgs("b1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow("b1", "Main"))
	w1 := httptest.NewRecorder()
	c1, _ := gin.CreateTestContext(w1)
	c1.Request, _ = http.NewRequest("GET", "/api/resources/buildings/b1", nil)
	c1.Params = gin.Params{{Key: "id", Value: "b1"}}
	handler.GetBuilding(c1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Update
	mock.ExpectExec(`UPDATE buildings`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	b := models.Building{Name: "Updated"}
	body, _ := json.Marshal(b)
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request, _ = http.NewRequest("PUT", "/api/resources/buildings/b1", bytes.NewBuffer(body))
	c2.Params = gin.Params{{Key: "id", Value: "b1"}}
	c2.Set("userID", "u1")
	handler.UpdateBuilding(c2)
	assert.Equal(t, http.StatusOK, w2.Code)

	// Delete
	mock.ExpectExec(`UPDATE buildings SET deleted_at=now\(\), updated_by=\$2 WHERE id=\$1`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	w3 := httptest.NewRecorder()
	c3, _ := gin.CreateTestContext(w3)
	c3.Request, _ = http.NewRequest("DELETE", "/api/resources/buildings/b1", nil)
	c3.Params = gin.Params{{Key: "id", Value: "b1"}}
	c3.Set("userID", "u1")
	handler.DeleteBuilding(c3)
	assert.Equal(t, http.StatusOK, w3.Code)
}

func TestResourceHandler_UpdateDeleteRoom(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewSQLResourceRepository(sqlxDB)
	svc := services.NewResourceService(repo)
	handler := NewResourceHandler(svc)

	buildingID := "00000000-0000-0000-0000-000000000001"

	// Update
	mock.ExpectExec(`UPDATE rooms`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	r := models.Room{BuildingID: buildingID, Name: "102", Capacity: 30, Floor: 1, Type: "classroom", Features: "[]", IsActive: true}
	body, _ := json.Marshal(r)
	w1 := httptest.NewRecorder()
	c1, _ := gin.CreateTestContext(w1)
	c1.Request, _ = http.NewRequest("PUT", "/api/resources/rooms/r1", bytes.NewBuffer(body))
	c1.Params = gin.Params{{Key: "id", Value: "r1"}}
	c1.Set("userID", "u1")
	handler.UpdateRoom(c1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Delete
	mock.ExpectExec(`UPDATE rooms SET deleted_at=now\(\), updated_by=\$2 WHERE id=\$1`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request, _ = http.NewRequest("DELETE", "/api/resources/rooms/r1", nil)
	c2.Params = gin.Params{{Key: "id", Value: "r1"}}
	c2.Set("userID", "u1")
	handler.DeleteRoom(c2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestResourceHandler_ErrorPaths(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewSQLResourceRepository(sqlxDB)
	svc := services.NewResourceService(repo)
	handler := NewResourceHandler(svc)

	t.Run("CreateBuilding_BindError", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/resources/buildings", bytes.NewBufferString("invalid"))
		handler.CreateBuilding(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("CreateBuilding_ServiceError", func(t *testing.T) {
		mock.ExpectQuery(`INSERT INTO buildings`).WillReturnError(assert.AnError) // Repo error
		
		b := models.Building{Name: "B1"}
		body, _ := json.Marshal(b)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/resources/buildings", bytes.NewBuffer(body))
		c.Set("tenant_id", "t1")
		c.Set("userID", "u1")
		
		handler.CreateBuilding(c)
		assert.Equal(t, http.StatusBadRequest, w.Code) // Handler returns 400 on service error unfortunately? Let's check handler. Yes, line 56.
	})

	t.Run("ListBuildings_ServiceError", func(t *testing.T) {
		mock.ExpectQuery(`SELECT`).WillReturnError(assert.AnError)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/resources/buildings", nil)
		c.Set("tenant_id", "t1")
		
		handler.ListBuildings(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("GetBuilding_ServiceError", func(t *testing.T) {
		mock.ExpectQuery(`SELECT`).WillReturnError(assert.AnError)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/resources/buildings/b1", nil)
		c.Params = gin.Params{{Key: "id", Value: "b1"}}
		
		handler.GetBuilding(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("UpdateBuilding_ServiceError", func(t *testing.T) {
		mock.ExpectExec(`UPDATE`).WillReturnError(assert.AnError)
		b := models.Building{Name: "Upd"}
		body, _ := json.Marshal(b)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/resources/buildings/b1", bytes.NewBuffer(body))
		c.Params = gin.Params{{Key: "id", Value: "b1"}}
		c.Set("userID", "u1")

		handler.UpdateBuilding(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
	
	t.Run("DeleteBuilding_ServiceError", func(t *testing.T) {
		mock.ExpectExec(`UPDATE`).WillReturnError(assert.AnError)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("DELETE", "/api/resources/buildings/b1", nil)
		c.Params = gin.Params{{Key: "id", Value: "b1"}}
		c.Set("userID", "u1")

		handler.DeleteBuilding(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

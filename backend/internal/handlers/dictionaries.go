package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// DictionaryHandler handles CRUD for Programs and Specialties
type DictionaryHandler struct {
	db *sqlx.DB
}

func NewDictionaryHandler(db *sqlx.DB) *DictionaryHandler {
	return &DictionaryHandler{db: db}
}

// --- Programs ---

type Program struct {
	ID        string `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	Code      string `db:"code" json:"code"`
	IsActive  bool   `db:"is_active" json:"is_active"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

type createProgramReq struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code"`
}

type updateProgramReq struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	IsActive *bool  `json:"is_active"`
}

func (h *DictionaryHandler) ListPrograms(c *gin.Context) {
	activeOnly := c.Query("active") == "true"
	query := `SELECT id, name, COALESCE(code, '') as code, is_active, 
              to_char(created_at,'YYYY-MM-DD"T"HH24:MI:SSZ') as created_at,
              to_char(updated_at,'YYYY-MM-DD"T"HH24:MI:SSZ') as updated_at 
              FROM programs`
	if activeOnly {
		query += ` WHERE is_active = true`
	}
	query += ` ORDER BY name`

	var programs []Program
	if err := h.db.Select(&programs, query); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if programs == nil {
		programs = []Program{}
	}
	c.JSON(http.StatusOK, programs)
}

func (h *DictionaryHandler) CreateProgram(c *gin.Context) {
	var req createProgramReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id string
	err := h.db.QueryRow(`INSERT INTO programs (name, code) VALUES ($1, $2) RETURNING id`, req.Name, nullable(req.Code)).Scan(&id)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			c.JSON(http.StatusConflict, gin.H{"error": "Program with this name already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *DictionaryHandler) UpdateProgram(c *gin.Context) {
	id := c.Param("id")
	var req updateProgramReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build dynamic update query
	setParts := []string{"updated_at = now()"}
	args := []interface{}{}
	argId := 1

	if req.Name != "" {
		setParts = append(setParts, "name = $"+itoa(argId))
		args = append(args, req.Name)
		argId++
	}
	if req.Code != "" {
		setParts = append(setParts, "code = $"+itoa(argId))
		args = append(args, req.Code)
		argId++
	}
	if req.IsActive != nil {
		setParts = append(setParts, "is_active = $"+itoa(argId))
		args = append(args, *req.IsActive)
		argId++
	}

	args = append(args, id)
	query := "UPDATE programs SET " + strings.Join(setParts, ", ") + " WHERE id = $" + itoa(argId)

	_, err := h.db.Exec(query, args...)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			c.JSON(http.StatusConflict, gin.H{"error": "Program with this name already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *DictionaryHandler) DeleteProgram(c *gin.Context) {
	id := c.Param("id")
	// Soft delete
	_, err := h.db.Exec(`UPDATE programs SET is_active = false, updated_at = now() WHERE id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// --- Specialties ---

type Specialty struct {
	ID         string   `db:"id" json:"id"`
	Name       string   `db:"name" json:"name"`
	Code       string   `db:"code" json:"code"`
	ProgramIDs []string `json:"program_ids"` // Multiple programs
	IsActive   bool     `db:"is_active" json:"is_active"`
	CreatedAt  string   `db:"created_at" json:"created_at"`
	UpdatedAt  string   `db:"updated_at" json:"updated_at"`
}

type createSpecialtyReq struct {
	Name       string   `json:"name" binding:"required"`
	Code       string   `json:"code"`
	ProgramIDs []string `json:"program_ids"` // Multiple programs
}

type updateSpecialtyReq struct {
	Name       string   `json:"name"`
	Code       string   `json:"code"`
	ProgramIDs []string `json:"program_ids"` // Multiple programs
	IsActive   *bool    `json:"is_active"`
}

func (h *DictionaryHandler) ListSpecialties(c *gin.Context) {
	activeOnly := c.Query("active") == "true"
	programID := c.Query("program_id")

	// Get base specialty data
	query := `SELECT id, name, COALESCE(code, '') as code, is_active, 
              to_char(created_at,'YYYY-MM-DD"T"HH24:MI:SSZ') as created_at,
              to_char(updated_at,'YYYY-MM-DD"T"HH24:MI:SSZ') as updated_at 
              FROM specialties WHERE 1=1`
	
	if activeOnly {
		query += ` AND is_active = true`
	}
	if programID != "" {
		query += ` AND EXISTS (SELECT 1 FROM specialty_programs WHERE specialty_id = specialties.id AND program_id = $1)`
	}
	query += ` ORDER BY name`

	type dbSpecialty struct {
		ID        string `db:"id"`
		Name      string `db:"name"`
		Code      string `db:"code"`
		IsActive  bool   `db:"is_active"`
		CreatedAt string `db:"created_at"`
		UpdatedAt string `db:"updated_at"`
	}

	var dbSpecialties []dbSpecialty
	args := []interface{}{}
	if programID != "" {
		args = append(args, programID)
	}

	if err := h.db.Select(&dbSpecialties, query, args...); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// For each specialty, fetch its program IDs
	specialties := make([]Specialty, len(dbSpecialties))
	for i, s := range dbSpecialties {
		var programIDs []string
		err := h.db.Select(&programIDs, `SELECT program_id FROM specialty_programs WHERE specialty_id = $1`, s.ID)
		if err != nil {
			programIDs = []string{}
		}

		specialties[i] = Specialty{
			ID:         s.ID,
			Name:       s.Name,
			Code:       s.Code,
			ProgramIDs: programIDs,
			IsActive:   s.IsActive,
			CreatedAt:  s.CreatedAt,
			UpdatedAt:  s.UpdatedAt,
		}
	}


	c.JSON(http.StatusOK, specialties)
}

func (h *DictionaryHandler) CreateSpecialty(c *gin.Context) {
	var req createSpecialtyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create specialty
	var id string
	err := h.db.QueryRow(`INSERT INTO specialties (name, code) VALUES ($1, $2) RETURNING id`, 
		req.Name, nullable(req.Code)).Scan(&id)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			c.JSON(http.StatusConflict, gin.H{"error": "Specialty with this name already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Link to programs if provided
	for _, programID := range req.ProgramIDs {
		if programID != "" && programID != "no_program" {
			_, err := h.db.Exec(`INSERT INTO specialty_programs (specialty_id, program_id) VALUES ($1, $2)`, id, programID)
			if err != nil {
				// Log error but don't fail the creation
				continue
			}
		}
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *DictionaryHandler) UpdateSpecialty(c *gin.Context) {
	id := c.Param("id")
	var req updateSpecialtyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	setParts := []string{"updated_at = now()"}
	args := []interface{}{}
	argId := 1

	if req.Name != "" {
		setParts = append(setParts, "name = $"+itoa(argId))
		args = append(args, req.Name)
		argId++
	}
	if req.Code != "" {
		setParts = append(setParts, "code = $"+itoa(argId))
		args = append(args, req.Code)
		argId++
	}
	if req.IsActive != nil {
		setParts = append(setParts, "is_active = $"+itoa(argId))
		args = append(args, *req.IsActive)
		argId++
	}

	// Update specialty basic data
	if len(setParts) > 1 { // More than just updated_at
		args = append(args, id)
		query := "UPDATE specialties SET " + strings.Join(setParts, ", ") + " WHERE id = $" + itoa(argId)
		_, err := h.db.Exec(query, args...)
		if err != nil {
			if strings.Contains(err.Error(), "unique constraint") {
				c.JSON(http.StatusConflict, gin.H{"error": "Specialty with this name already exists"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Update program relationships if provided
	if req.ProgramIDs != nil {
		// Delete existing relationships
		_, _ = h.db.Exec(`DELETE FROM specialty_programs WHERE specialty_id = $1`, id)
		// Add new relationships
		for _, programID := range req.ProgramIDs {
			if programID != "" && programID != "no_program" {
				_, err := h.db.Exec(`INSERT INTO specialty_programs (specialty_id, program_id) VALUES ($1, $2)`, id, programID)
				if err != nil {
					continue
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *DictionaryHandler) DeleteSpecialty(c *gin.Context) {
	id := c.Param("id")
	_, err := h.db.Exec(`UPDATE specialties SET is_active = false, updated_at = now() WHERE id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// --- Cohorts ---

type Cohort struct {
	ID        string `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	StartDate string `db:"start_date" json:"start_date"`
	EndDate   string `db:"end_date" json:"end_date"`
	IsActive  bool   `db:"is_active" json:"is_active"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

type createCohortReq struct {
	Name      string `json:"name" binding:"required"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type updateCohortReq struct {
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	IsActive  *bool  `json:"is_active"`
}

func (h *DictionaryHandler) ListCohorts(c *gin.Context) {
	activeOnly := c.Query("active") == "true"
	query := `SELECT id, name, 
              COALESCE(to_char(start_date, 'YYYY-MM-DD'), '') as start_date,
              COALESCE(to_char(end_date, 'YYYY-MM-DD'), '') as end_date,
              is_active, 
              to_char(created_at,'YYYY-MM-DD"T"HH24:MI:SSZ') as created_at,
              to_char(updated_at,'YYYY-MM-DD"T"HH24:MI:SSZ') as updated_at 
              FROM cohorts`
	if activeOnly {
		query += ` WHERE is_active = true`
	}
	query += ` ORDER BY name DESC` // Usually want newest cohorts first

	var cohorts []Cohort
	if err := h.db.Select(&cohorts, query); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if cohorts == nil {
		cohorts = []Cohort{}
	}
	c.JSON(http.StatusOK, cohorts)
}

func (h *DictionaryHandler) CreateCohort(c *gin.Context) {
	var req createCohortReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id string
	err := h.db.QueryRow(`INSERT INTO cohorts (name, start_date, end_date) VALUES ($1, $2, $3) RETURNING id`, 
		req.Name, nullable(req.StartDate), nullable(req.EndDate)).Scan(&id)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			c.JSON(http.StatusConflict, gin.H{"error": "Cohort with this name already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *DictionaryHandler) UpdateCohort(c *gin.Context) {
	id := c.Param("id")
	var req updateCohortReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	setParts := []string{"updated_at = now()"}
	args := []interface{}{}
	argId := 1

	if req.Name != "" {
		setParts = append(setParts, "name = $"+itoa(argId))
		args = append(args, req.Name)
		argId++
	}
	if req.StartDate != "" {
		setParts = append(setParts, "start_date = $"+itoa(argId))
		args = append(args, nullable(req.StartDate))
		argId++
	}
	if req.EndDate != "" {
		setParts = append(setParts, "end_date = $"+itoa(argId))
		args = append(args, nullable(req.EndDate))
		argId++
	}
	if req.IsActive != nil {
		setParts = append(setParts, "is_active = $"+itoa(argId))
		args = append(args, *req.IsActive)
		argId++
	}

	args = append(args, id)
	query := "UPDATE cohorts SET " + strings.Join(setParts, ", ") + " WHERE id = $" + itoa(argId)

	_, err := h.db.Exec(query, args...)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			c.JSON(http.StatusConflict, gin.H{"error": "Cohort with this name already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *DictionaryHandler) DeleteCohort(c *gin.Context) {
	id := c.Param("id")
	_, err := h.db.Exec(`UPDATE cohorts SET is_active = false, updated_at = now() WHERE id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// --- Departments ---

type Department struct {
	ID        string `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	Code      string `db:"code" json:"code"`
	IsActive  bool   `db:"is_active" json:"is_active"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

type createDepartmentReq struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code"`
}

type updateDepartmentReq struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	IsActive *bool  `json:"is_active"`
}

func (h *DictionaryHandler) ListDepartments(c *gin.Context) {
	activeOnly := c.Query("active") == "true"
	query := `SELECT id, name, COALESCE(code, '') as code, is_active, 
              to_char(created_at,'YYYY-MM-DD"T"HH24:MI:SSZ') as created_at,
              to_char(updated_at,'YYYY-MM-DD"T"HH24:MI:SSZ') as updated_at 
              FROM departments`
	if activeOnly {
		query += ` WHERE is_active = true`
	}
	query += ` ORDER BY name ASC`

	var departments []Department
	if err := h.db.Select(&departments, query); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if departments == nil {
		departments = []Department{}
	}
	c.JSON(http.StatusOK, departments)
}

func (h *DictionaryHandler) CreateDepartment(c *gin.Context) {
	var req createDepartmentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id string
	err := h.db.QueryRow(`INSERT INTO departments (name, code) VALUES ($1, $2) RETURNING id`, 
		req.Name, nullable(req.Code)).Scan(&id)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			c.JSON(http.StatusConflict, gin.H{"error": "Department with this name already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *DictionaryHandler) UpdateDepartment(c *gin.Context) {
	id := c.Param("id")
	var req updateDepartmentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	setParts := []string{"updated_at = now()"}
	args := []interface{}{}
	argId := 1

	if req.Name != "" {
		setParts = append(setParts, "name = $"+itoa(argId))
		args = append(args, req.Name)
		argId++
	}
	if req.Code != "" {
		setParts = append(setParts, "code = $"+itoa(argId))
		args = append(args, nullable(req.Code))
		argId++
	}
	if req.IsActive != nil {
		setParts = append(setParts, "is_active = $"+itoa(argId))
		args = append(args, *req.IsActive)
		argId++
	}

	args = append(args, id)
	query := "UPDATE departments SET " + strings.Join(setParts, ", ") + " WHERE id = $" + itoa(argId)

	_, err := h.db.Exec(query, args...)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			c.JSON(http.StatusConflict, gin.H{"error": "Department with this name already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *DictionaryHandler) DeleteDepartment(c *gin.Context) {
	id := c.Param("id")
	_, err := h.db.Exec(`UPDATE departments SET is_active = false, updated_at = now() WHERE id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// Helper to convert int to string
func itoa(i int) string {
	// Simple implementation using strconv
	return strconv.Itoa(i)
}

package handlers

import (
	"net/http"
	"strings"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/middleware"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

// DictionaryHandler handles CRUD for Programs and Specialties
type DictionaryHandler struct {
	svc *services.DictionaryService
}

func NewDictionaryHandler(svc *services.DictionaryService) *DictionaryHandler {
	return &DictionaryHandler{svc: svc}
}

// --- Programs ---

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
	tenantID := middleware.GetTenantID(c)
	activeOnly := c.Query("active") == "true"
	
	programs, err := h.svc.ListPrograms(c.Request.Context(), tenantID, activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, programs)
}

func (h *DictionaryHandler) CreateProgram(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	var req createProgramReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.svc.CreateProgram(c.Request.Context(), tenantID, req.Name, req.Code)
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
	tenantID := middleware.GetTenantID(c)
	var req updateProgramReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.UpdateProgram(c.Request.Context(), tenantID, id, req.Name, req.Code, req.IsActive)
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
	tenantID := middleware.GetTenantID(c)
	
	err := h.svc.DeleteProgram(c.Request.Context(), tenantID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// --- Specialties ---

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
	tenantID := middleware.GetTenantID(c)

	specialties, err := h.svc.ListSpecialties(c.Request.Context(), tenantID, activeOnly, programID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, specialties)
}

func (h *DictionaryHandler) CreateSpecialty(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	var req createSpecialtyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.svc.CreateSpecialty(c.Request.Context(), tenantID, req.Name, req.Code, req.ProgramIDs)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			c.JSON(http.StatusConflict, gin.H{"error": "Specialty with this name already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *DictionaryHandler) UpdateSpecialty(c *gin.Context) {
	id := c.Param("id")
	tenantID := middleware.GetTenantID(c)
	var req updateSpecialtyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.UpdateSpecialty(c.Request.Context(), tenantID, id, req.Name, req.Code, req.IsActive, req.ProgramIDs)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			c.JSON(http.StatusConflict, gin.H{"error": "Specialty with this name already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *DictionaryHandler) DeleteSpecialty(c *gin.Context) {
	id := c.Param("id")
	tenantID := middleware.GetTenantID(c)
	err := h.svc.DeleteSpecialty(c.Request.Context(), tenantID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// --- Cohorts ---

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
	tenantID := middleware.GetTenantID(c)
	activeOnly := c.Query("active") == "true"
	
	cohorts, err := h.svc.ListCohorts(c.Request.Context(), tenantID, activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cohorts)
}

func (h *DictionaryHandler) CreateCohort(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	var req createCohortReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.svc.CreateCohort(c.Request.Context(), tenantID, req.Name, req.StartDate, req.EndDate)
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
	tenantID := middleware.GetTenantID(c)
	var req updateCohortReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.UpdateCohort(c.Request.Context(), tenantID, id, req.Name, req.StartDate, req.EndDate, req.IsActive)
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
	tenantID := middleware.GetTenantID(c)
	err := h.svc.DeleteCohort(c.Request.Context(), tenantID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// --- Departments ---

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
	tenantID := middleware.GetTenantID(c)
	
	departments, err := h.svc.ListDepartments(c.Request.Context(), tenantID, activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, departments)
}

func (h *DictionaryHandler) CreateDepartment(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	var req createDepartmentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.svc.CreateDepartment(c.Request.Context(), tenantID, req.Name, req.Code)
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
	tenantID := middleware.GetTenantID(c)
	var req updateDepartmentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.UpdateDepartment(c.Request.Context(), tenantID, id, req.Name, req.Code, req.IsActive)
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
	tenantID := middleware.GetTenantID(c)
	err := h.svc.DeleteDepartment(c.Request.Context(), tenantID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

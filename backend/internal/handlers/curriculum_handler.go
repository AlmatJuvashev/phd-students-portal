package handlers

import (
	"net/http"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/middleware"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type CurriculumHandler struct {
	svc *services.CurriculumService
}

func NewCurriculumHandler(svc *services.CurriculumService) *CurriculumHandler {
	return &CurriculumHandler{svc: svc}
}

// Programs

func (h *CurriculumHandler) ListPrograms(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	list, err := h.svc.ListPrograms(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *CurriculumHandler) GetProgram(c *gin.Context) {
	id := c.Param("id")
	p, err := h.svc.GetProgram(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if p == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "program not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *CurriculumHandler) CreateProgram(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	var p models.Program
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.TenantID = tenantID
	
	if err := h.svc.CreateProgram(c.Request.Context(), &p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *CurriculumHandler) UpdateProgram(c *gin.Context) {
	id := c.Param("id")
	var p models.Program
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.ID = id
	// Maintain tenant_id from existing or context? 
	// Ideally we check ownership. For now assume middleware validated access.
	
	if err := h.svc.UpdateProgram(c.Request.Context(), &p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *CurriculumHandler) DeleteProgram(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteProgram(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// Courses

func (h *CurriculumHandler) ListCourses(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	progID := c.Query("program_id")
	var pID *string
	if progID != "" {
		pID = &progID
	}
	
	list, err := h.svc.ListCourses(c.Request.Context(), tenantID, pID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *CurriculumHandler) CreateCourse(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	var obj models.Course
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	obj.TenantID = tenantID
	
	if err := h.svc.CreateCourse(c.Request.Context(), &obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, obj)
}

// ... Additional helper methods for other entities can be added here

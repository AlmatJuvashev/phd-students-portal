package handlers

import (
	"net/http"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type CourseContentHandler struct {
	svc *services.CourseContentService
}

func NewCourseContentHandler(svc *services.CourseContentService) *CourseContentHandler {
	return &CourseContentHandler{svc: svc}
}

// Modules

func (h *CourseContentHandler) ListModules(c *gin.Context) {
	courseID := c.Query("course_id")
	if courseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "course_id required"})
		return
	}
	list, err := h.svc.ListModules(c.Request.Context(), courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *CourseContentHandler) CreateModule(c *gin.Context) {
	var m models.CourseModule
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.CreateModule(c.Request.Context(), &m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, m)
}

func (h *CourseContentHandler) UpdateModule(c *gin.Context) {
	id := c.Param("id")
	var m models.CourseModule
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m.ID = id
	if err := h.svc.UpdateModule(c.Request.Context(), &m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, m)
}

func (h *CourseContentHandler) DeleteModule(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteModule(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// Lessons

func (h *CourseContentHandler) ListLessons(c *gin.Context) {
	moduleID := c.Query("module_id")
	if moduleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "module_id required"})
		return
	}
	list, err := h.svc.ListLessons(c.Request.Context(), moduleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *CourseContentHandler) CreateLesson(c *gin.Context) {
	var l models.CourseLesson
	if err := c.ShouldBindJSON(&l); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.CreateLesson(c.Request.Context(), &l); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, l)
}

func (h *CourseContentHandler) UpdateLesson(c *gin.Context) {
	id := c.Param("id")
	var l models.CourseLesson
	if err := c.ShouldBindJSON(&l); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	l.ID = id
	if err := h.svc.UpdateLesson(c.Request.Context(), &l); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, l)
}

func (h *CourseContentHandler) DeleteLesson(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteLesson(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// Activities

func (h *CourseContentHandler) ListActivities(c *gin.Context) {
	lessonID := c.Query("lesson_id")
	if lessonID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lesson_id required"})
		return
	}
	list, err := h.svc.ListActivities(c.Request.Context(), lessonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *CourseContentHandler) CreateActivity(c *gin.Context) {
	var a models.CourseActivity
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.CreateActivity(c.Request.Context(), &a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, a)
}

func (h *CourseContentHandler) UpdateActivity(c *gin.Context) {
	id := c.Param("id")
	var a models.CourseActivity
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	a.ID = id
	if err := h.svc.UpdateActivity(c.Request.Context(), &a); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, a)
}

func (h *CourseContentHandler) DeleteActivity(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteActivity(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

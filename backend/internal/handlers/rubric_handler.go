package handlers

import (
	"net/http"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type RubricHandler struct {
	svc *services.RubricService
}

func NewRubricHandler(svc *services.RubricService) *RubricHandler {
	return &RubricHandler{svc: svc}
}

// CreateRubric POST /courses/:id/rubrics
func (h *RubricHandler) CreateRubric(c *gin.Context) {
	courseID := c.Param("id")
	var r models.Rubric
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r.CourseOfferingID = courseID
	
	created, err := h.svc.CreateRubric(c.Request.Context(), &r)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// ListRubrics GET /courses/:id/rubrics
func (h *RubricHandler) ListRubrics(c *gin.Context) {
	courseID := c.Param("id")
	list, err := h.svc.ListRubrics(c.Request.Context(), courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// GetRubric GET /rubrics/:id
func (h *RubricHandler) GetRubric(c *gin.Context) {
	id := c.Param("id")
	r, err := h.svc.GetRubric(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, r)
}

// SubmitGrade POST /submissions/:id/rubric_grade
func (h *RubricHandler) SubmitGrade(c *gin.Context) {
	submissionID := c.Param("id")
	var input services.GradeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.SubmissionID = submissionID
	input.GraderID = userIDFromClaims(c)

	grade, err := h.svc.SubmitGrade(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, grade)
}

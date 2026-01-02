package handlers

import (
	"context"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BulkEnroller interface {
	ImportStudents(ctx context.Context, r io.Reader, tenantID string) (int, []error)
}

type BulkHandler struct {
	service BulkEnroller
}

func NewBulkHandler(service BulkEnroller) *BulkHandler {
	return &BulkHandler{service: service}
}

// BulkEnrollStudents godoc
// @Summary Bulk enroll students
// @Description Import students from CSV file
// @Tags admin
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV File (first_name, last_name, email, role)"
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/bulk/enroll [post]
func (h *BulkHandler) BulkEnrollStudents(c *gin.Context) {
	// Get Tenant ID from context (set by middleware)
	tenantID := c.GetString("tenantID")
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant context missing"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload required"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer file.Close()

	count, errs := h.service.ImportStudents(c.Request.Context(), file, tenantID)
	
	response := gin.H{
		"success_count": count,
	}
	
	if len(errs) > 0 {
		var errMsgs []string
		for _, e := range errs {
			errMsgs = append(errMsgs, e.Error())
		}
		response["errors"] = errMsgs
		// Partial success is still 200 OK often, or 207 Multi-Status.
		// For simplicity, 200 with error list.
	}

	c.JSON(http.StatusOK, response)
}

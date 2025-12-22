package handlers

import (
	"net/http"
	"strings"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	svc *services.SearchService
	cfg config.AppConfig
}

func NewSearchHandler(svc *services.SearchService, cfg config.AppConfig) *SearchHandler {
	return &SearchHandler{svc: svc, cfg: cfg}
}

func (h *SearchHandler) GlobalSearch(c *gin.Context) {
	query := strings.TrimSpace(c.Query("q"))
	if len(query) < 2 {
		c.JSON(http.StatusOK, []models.SearchResult{})
		return
	}

	// RBAC
	role := roleFromContext(c)
	userID := userIDFromClaims(c)

	results, err := h.svc.GlobalSearch(c.Request.Context(), query, role, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

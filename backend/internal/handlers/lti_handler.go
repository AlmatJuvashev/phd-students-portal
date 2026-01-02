package handlers

import (
	"net/http"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type LTIHandler struct {
	svc *services.LTIService
	cfg config.AppConfig
}

func NewLTIHandler(svc *services.LTIService, cfg config.AppConfig) *LTIHandler {
	return &LTIHandler{svc: svc, cfg: cfg}
}

// RegisterTool - Admin only
func (h *LTIHandler) RegisterTool(c *gin.Context) {
	var req models.CreateToolParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Tenant check is usually middleware, or explicit here
	// For MVP allow passing tenant_id
	
	tool, err := h.svc.RegisterTool(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register tool"})
		return
	}
	c.JSON(http.StatusCreated, tool)
}

func (h *LTIHandler) ListTools(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id required"})
		return
	}
	tools, err := h.svc.ListTools(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tools)
}

// LoginInit starts the launch flow
func (h *LTIHandler) LoginInit(c *gin.Context) {
	toolID := c.Query("tool_id")
	targetLinkURI := c.Query("target_link_uri")
	userID := c.GetString("userID") // From middleware
	
	if toolID == "" || targetLinkURI == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tool_id and target_link_uri required"})
		return
	}

	redirectURL, err := h.svc.GenerateLoginInit(c.Request.Context(), toolID, userID, targetLinkURI)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.Redirect(http.StatusFound, redirectURL)
}

func (h *LTIHandler) GetJWKS(c *gin.Context) {
	jwks, err := h.svc.GetJWKS(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, jwks)
}

// Launch endpoint would handle the POST back from Tool
func (h *LTIHandler) Launch(c *gin.Context) {
	// Phase 19.3
	c.JSON(http.StatusNotImplemented, gin.H{"error": "launch validation not ready"})
}

package handlers

import (
	"net/http"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type ProgramBuilderHandler struct {
	builderService *services.ProgramBuilderService
}

func NewProgramBuilderHandler(builderService *services.ProgramBuilderService) *ProgramBuilderHandler {
	return &ProgramBuilderHandler{builderService: builderService}
}

// GetNodes returns all nodes for a program's journey map.
// GET /api/admin/programs/:id/builder/nodes
func (h *ProgramBuilderHandler) GetNodes(c *gin.Context) {
	programID := c.Param("id")
	
	// Ensure map exists first
	if _, err := h.builderService.EnsureDraftMap(c.Request.Context(), programID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to ensure draft map: " + err.Error()})
		return
	}

	nodes, err := h.builderService.GetNodes(c.Request.Context(), programID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nodes)
}

// CreateNode adds a new node to the program journey.
// POST /api/admin/programs/:id/builder/nodes
func (h *ProgramBuilderHandler) CreateNode(c *gin.Context) {
	programID := c.Param("id")

	var req struct {
		Node   models.JourneyNodeDefinition `json:"node"`
		Config models.ProgramNodeConfig     `json:"config"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jm, err := h.builderService.EnsureDraftMap(c.Request.Context(), programID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	created, err := h.builderService.CreateNode(c.Request.Context(), jm.ID, req.Node, req.Config)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// UpdateNode updates an existing node config.
// PUT /api/admin/programs/:id/builder/nodes/:nodeId
func (h *ProgramBuilderHandler) UpdateNode(c *gin.Context) {
	nodeID := c.Param("nodeId")

	var config models.ProgramNodeConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.builderService.UpdateNodeConfig(c.Request.Context(), nodeID, config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

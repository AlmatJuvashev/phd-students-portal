package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type ProgramBuilderHandler struct {
	builderService *services.ProgramBuilderService
}

func NewProgramBuilderHandler(builderService *services.ProgramBuilderService) *ProgramBuilderHandler {
	return &ProgramBuilderHandler{builderService: builderService}
}

// GetJourneyMap returns the full map structure.
// GET /api/programs/:id/builder/map
func (h *ProgramBuilderHandler) GetJourneyMap(c *gin.Context) {
	programID := c.Param("id")
	
	fullMap, err := h.builderService.GetJourneyMap(c.Request.Context(), programID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fullMap)
}

// UpdateJourneyMap updates the program version config (phases/layout) and/or metadata.
// PUT /api/programs/:id/builder/map
func (h *ProgramBuilderHandler) UpdateJourneyMap(c *gin.Context) {
	programID := c.Param("id")

	raw, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	var payload services.UpdateMapInput
	var wrapped struct {
		Map       services.UpdateMapInput `json:"map"`
		JourneyMap services.UpdateMapInput `json:"journey_map"`
	}
	if err := json.Unmarshal(raw, &wrapped); err == nil && (wrapped.Map != (services.UpdateMapInput{}) || wrapped.JourneyMap != (services.UpdateMapInput{})) {
		if wrapped.Map != (services.UpdateMapInput{}) {
			payload = wrapped.Map
		} else {
			payload = wrapped.JourneyMap
		}
	} else if err := json.Unmarshal(raw, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if payload.Title == nil && payload.Version == nil && payload.Config == nil && payload.Phases == nil && payload.IsActive == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	updated, err := h.builderService.UpdateJourneyMap(c.Request.Context(), programID, payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
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

	raw, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	type nodePayload struct {
		ParentNodeID  *string          `json:"parent_node_id,omitempty"`
		Slug          string           `json:"slug"`
		Type          string           `json:"type"`
		Title         json.RawMessage  `json:"title"`
		Description   json.RawMessage  `json:"description"`
		ModuleKey     string           `json:"module_key"`
		Coordinates   json.RawMessage  `json:"coordinates"`
		Config        json.RawMessage  `json:"config"`
		Prerequisites []string         `json:"prerequisites"`
	}

	var payload nodePayload
	var wrapped struct {
		Node nodePayload `json:"node"`
	}
	if err := json.Unmarshal(raw, &wrapped); err == nil && (wrapped.Node.Slug != "" || wrapped.Node.Type != "") {
		payload = wrapped.Node
	} else if err := json.Unmarshal(raw, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if payload.Slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug is required"})
		return
	}
	if payload.Type == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type is required"})
		return
	}

	jm, err := h.builderService.EnsureDraftMap(c.Request.Context(), programID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	node := &models.JourneyNodeDefinition{
		ParentNodeID:  payload.ParentNodeID,
		Slug:          payload.Slug,
		Type:          payload.Type,
		Title:         string(payload.Title),
		Description:   string(payload.Description),
		ModuleKey:     payload.ModuleKey,
		Coordinates:   string(payload.Coordinates),
		Config:        string(payload.Config),
		Prerequisites: pq.StringArray(payload.Prerequisites),
	}

	if err := h.builderService.CreateNode(c.Request.Context(), jm.ID, node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.builderService.GetBuilderNode(c.Request.Context(), node.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// UpdateNode updates an existing node config.
// PUT /api/admin/programs/:id/builder/nodes/:nodeId
func (h *ProgramBuilderHandler) UpdateNode(c *gin.Context) {
	nodeID := c.Param("nodeId")

	raw, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	type updatePayload struct {
		ParentNodeID  *string          `json:"parent_node_id,omitempty"`
		Slug          *string          `json:"slug,omitempty"`
		Type          *string          `json:"type,omitempty"`
		Title         *json.RawMessage `json:"title,omitempty"`
		Description   *json.RawMessage `json:"description,omitempty"`
		ModuleKey     *string          `json:"module_key,omitempty"`
		Coordinates   *json.RawMessage `json:"coordinates,omitempty"`
		Config        *json.RawMessage `json:"config,omitempty"`
		Prerequisites *[]string        `json:"prerequisites,omitempty"`
	}

	var payload updatePayload
	var wrapped struct {
		Node updatePayload `json:"node"`
	}
	if err := json.Unmarshal(raw, &wrapped); err == nil && wrapped.Node != (updatePayload{}) {
		payload = wrapped.Node
	} else if err := json.Unmarshal(raw, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	node, err := h.builderService.GetNode(c.Request.Context(), nodeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if node == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "node not found"})
		return
	}

	if payload.ParentNodeID != nil {
		node.ParentNodeID = payload.ParentNodeID
	}
	if payload.Slug != nil {
		node.Slug = *payload.Slug
	}
	if payload.Type != nil {
		node.Type = *payload.Type
	}
	if payload.Title != nil {
		node.Title = string(*payload.Title)
	}
	if payload.Description != nil {
		node.Description = string(*payload.Description)
	}
	if payload.ModuleKey != nil {
		node.ModuleKey = *payload.ModuleKey
	}
	if payload.Coordinates != nil {
		node.Coordinates = string(*payload.Coordinates)
	}
	if payload.Config != nil {
		node.Config = string(*payload.Config)
	}
	if payload.Prerequisites != nil {
		node.Prerequisites = pq.StringArray(*payload.Prerequisites)
	}

	if err := h.builderService.UpdateNode(c.Request.Context(), node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.builderService.GetBuilderNode(c.Request.Context(), nodeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

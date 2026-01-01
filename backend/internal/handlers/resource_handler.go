package handlers

import (
	"net/http"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/middleware"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type ResourceHandler struct {
	svc *services.ResourceService
}

func NewResourceHandler(svc *services.ResourceService) *ResourceHandler {
	return &ResourceHandler{svc: svc}
}

// Buildings

func (h *ResourceHandler) ListBuildings(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	list, err := h.svc.ListBuildings(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *ResourceHandler) GetBuilding(c *gin.Context) {
	id := c.Param("id")
	b, err := h.svc.GetBuilding(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Simplified, could check for NoRows
		return
	}
	c.JSON(http.StatusOK, b)
}

func (h *ResourceHandler) CreateBuilding(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	var b models.Building
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	b.TenantID = tenantID
	
	if err := h.svc.CreateBuilding(c.Request.Context(), &b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, b)
}

func (h *ResourceHandler) UpdateBuilding(c *gin.Context) {
	id := c.Param("id")
	var b models.Building
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	b.ID = id
	if err := h.svc.UpdateBuilding(c.Request.Context(), &b); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, b)
}

func (h *ResourceHandler) DeleteBuilding(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteBuilding(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// Rooms

func (h *ResourceHandler) ListRooms(c *gin.Context) {
	buildingID := c.Query("building_id")
	if buildingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "building_id required"})
		return
	}
	list, err := h.svc.ListRooms(c.Request.Context(), buildingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *ResourceHandler) CreateRoom(c *gin.Context) {
	var r models.Room
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// BuildingID must be in body
	if err := h.svc.CreateRoom(c.Request.Context(), &r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, r)
}

func (h *ResourceHandler) UpdateRoom(c *gin.Context) {
	id := c.Param("id")
	var r models.Room
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r.ID = id
	if err := h.svc.UpdateRoom(c.Request.Context(), &r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, r)
}

func (h *ResourceHandler) DeleteRoom(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteRoom(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

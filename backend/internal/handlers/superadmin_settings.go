package handlers

import (
	"net/http"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

// SuperadminSettingsHandler handles global settings operations
type SuperadminSettingsHandler struct {
	adminSvc *services.SuperAdminService
	cfg      config.AppConfig
}

// NewSuperadminSettingsHandler creates a new handler
func NewSuperadminSettingsHandler(adminSvc *services.SuperAdminService, cfg config.AppConfig) *SuperadminSettingsHandler {
	return &SuperadminSettingsHandler{adminSvc: adminSvc, cfg: cfg}
}

// ListSettings returns all global settings
func (h *SuperadminSettingsHandler) ListSettings(c *gin.Context) {
	category := c.Query("category")
	settings, err := h.adminSvc.ListSettings(c.Request.Context(), category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch settings"})
		return
	}
	c.JSON(http.StatusOK, settings)
}

// GetSetting returns a single setting
func (h *SuperadminSettingsHandler) GetSetting(c *gin.Context) {
	key := c.Param("key")
	setting, err := h.adminSvc.GetSetting(c.Request.Context(), key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch setting"})
		return
	}
	if setting == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "setting not found"})
		return
	}
	c.JSON(http.StatusOK, setting)
}

// UpdateSettingRequest
type UpdateSettingRequest struct {
	Value       interface{} `json:"value" binding:"required"`
	Description *string     `json:"description"`
	Category    *string     `json:"category"`
}

// UpdateSetting updates/creates a global setting
func (h *SuperadminSettingsHandler) UpdateSetting(c *gin.Context) {
	key := c.Param("key")

	var req UpdateSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := models.UpdateSettingParams{
		Value:       req.Value,
		Description: req.Description,
		Category:    req.Category,
		UpdatedBy:   c.GetString("userID"),
	}

	setting, err := h.adminSvc.UpdateSetting(c.Request.Context(), key, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update setting"})
		return
	}

	_ = h.adminSvc.LogActivity(c.Request.Context(), models.ActivityLogParams{
		UserID:      strPtr(c.GetString("userID")),
		Action:      "update",
		EntityType:  "setting",
		EntityID:    key,
		Description: "Updated setting: " + key,
		IPAddress:   c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
	})

	c.JSON(http.StatusOK, setting)
}

// DeleteSetting deletes a setting
func (h *SuperadminSettingsHandler) DeleteSetting(c *gin.Context) {
	key := c.Param("key")
	err := h.adminSvc.DeleteSetting(c.Request.Context(), key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete setting"})
		return
	}

	_ = h.adminSvc.LogActivity(c.Request.Context(), models.ActivityLogParams{
		UserID:      strPtr(c.GetString("userID")),
		Action:      "delete",
		EntityType:  "setting",
		EntityID:    key,
		Description: "Deleted setting: " + key,
		IPAddress:   c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
	})

	c.JSON(http.StatusOK, gin.H{"message": "setting deleted"})
}

// GetCategories returns unique setting categories
func (h *SuperadminSettingsHandler) GetCategories(c *gin.Context) {
	cats, err := h.adminSvc.GetCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch categories"})
		return
	}
	c.JSON(http.StatusOK, cats)
}

// BulkUpdateRequest
type BulkUpdateRequest struct {
	Settings map[string]interface{} `json:"settings" binding:"required"`
}

// BulkUpdate updates multiple settings
func (h *SuperadminSettingsHandler) BulkUpdate(c *gin.Context) {
	var req BulkUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated := 0
	userID := c.GetString("userID")
	// Service doesn't have BulkUpdate, so loop here or add to service.
	// Will loop here for simplicity as I didn't add BulkUpdate to Repository.
	
	for key, value := range req.Settings {
		params := models.UpdateSettingParams{
			Value:     value,
			UpdatedBy: userID,
		}
		// Try update each, ignore errors? Or stop?
		// Handler logic was: "if err == nil { updated++ }"
		_, err := h.adminSvc.UpdateSetting(c.Request.Context(), key, params)
		if err == nil {
			updated++
		}
	}

	_ = h.adminSvc.LogActivity(c.Request.Context(), models.ActivityLogParams{
		UserID:      strPtr(userID),
		Action:      "bulk_update",
		EntityType:  "settings",
		Description: "Bulk updated settings",
		IPAddress:   c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
	})

	c.JSON(http.StatusOK, gin.H{"updated": updated})
}

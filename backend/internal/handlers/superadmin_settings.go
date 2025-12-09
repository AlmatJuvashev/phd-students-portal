package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// SuperadminSettingsHandler handles global settings operations for superadmins
type SuperadminSettingsHandler struct {
	db  *sqlx.DB
	cfg config.AppConfig
}

// NewSuperadminSettingsHandler creates a new superadmin settings handler
func NewSuperadminSettingsHandler(db *sqlx.DB, cfg config.AppConfig) *SuperadminSettingsHandler {
	return &SuperadminSettingsHandler{db: db, cfg: cfg}
}

// SettingResponse is the API response for a setting
type SettingResponse struct {
	Key         string          `json:"key" db:"key"`
	Value       json.RawMessage `json:"value" db:"value"`
	Description *string         `json:"description" db:"description"`
	Category    string          `json:"category" db:"category"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
	UpdatedBy   *string         `json:"updated_by" db:"updated_by"`
}

// ListSettings returns all global settings
func (h *SuperadminSettingsHandler) ListSettings(c *gin.Context) {
	// Optional category filter
	category := c.Query("category")

	query := `
		SELECT key, value, description, COALESCE(category, 'general') as category, updated_at, updated_by
		FROM global_settings
	`
	var args []interface{}
	if category != "" {
		query += " WHERE category = $1"
		args = append(args, category)
	}
	query += " ORDER BY category, key"

	var settings []SettingResponse
	var err error
	if len(args) > 0 {
		err = h.db.Select(&settings, query, args...)
	} else {
		err = h.db.Select(&settings, query)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch settings"})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// GetSetting returns a single setting by key
func (h *SuperadminSettingsHandler) GetSetting(c *gin.Context) {
	key := c.Param("key")

	query := `
		SELECT key, value, description, COALESCE(category, 'general') as category, updated_at, updated_by
		FROM global_settings
		WHERE key = $1
	`

	var setting SettingResponse
	err := h.db.Get(&setting, query, key)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "setting not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch setting"})
		return
	}

	c.JSON(http.StatusOK, setting)
}

// UpdateSettingRequest is the request body for updating a setting
type UpdateSettingRequest struct {
	Value       interface{} `json:"value" binding:"required"`
	Description *string     `json:"description"`
	Category    *string     `json:"category"`
}

// UpdateSetting updates a global setting
func (h *SuperadminSettingsHandler) UpdateSetting(c *gin.Context) {
	key := c.Param("key")

	var req UpdateSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert value to JSON
	valueJSON, err := json.Marshal(req.Value)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid value format"})
		return
	}

	userID := c.GetString("userID")

	query := `
		INSERT INTO global_settings (key, value, description, category, updated_at, updated_by)
		VALUES ($1, $2, $3, COALESCE($4, 'general'), now(), $5)
		ON CONFLICT (key) DO UPDATE SET
			value = $2,
			description = COALESCE($3, global_settings.description),
			category = COALESCE($4, global_settings.category),
			updated_at = now(),
			updated_by = $5
		RETURNING key, value, description, category, updated_at, updated_by
	`

	var setting SettingResponse
	err = h.db.QueryRowx(query, key, valueJSON, req.Description, req.Category, userID).StructScan(&setting)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update setting"})
		return
	}

	// Log activity
	logActivity(h.db, c, "update", "setting", key, "Updated setting: "+key, nil)

	c.JSON(http.StatusOK, setting)
}

// DeleteSetting deletes a global setting
func (h *SuperadminSettingsHandler) DeleteSetting(c *gin.Context) {
	key := c.Param("key")

	result, err := h.db.Exec(`DELETE FROM global_settings WHERE key = $1`, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete setting"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "setting not found"})
		return
	}

	// Log activity
	logActivity(h.db, c, "delete", "setting", key, "Deleted setting: "+key, nil)

	c.JSON(http.StatusOK, gin.H{"message": "setting deleted"})
}

// GetCategories returns distinct setting categories
func (h *SuperadminSettingsHandler) GetCategories(c *gin.Context) {
	var categories []string
	h.db.Select(&categories, `SELECT DISTINCT COALESCE(category, 'general') as category FROM global_settings ORDER BY category`)
	c.JSON(http.StatusOK, categories)
}

// BulkUpdateRequest is the request body for bulk updating settings
type BulkUpdateRequest struct {
	Settings map[string]interface{} `json:"settings" binding:"required"`
}

// BulkUpdate updates multiple settings at once
func (h *SuperadminSettingsHandler) BulkUpdate(c *gin.Context) {
	var req BulkUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("userID")
	updated := 0

	for key, value := range req.Settings {
		valueJSON, err := json.Marshal(value)
		if err != nil {
			continue
		}

		query := `
			INSERT INTO global_settings (key, value, updated_at, updated_by)
			VALUES ($1, $2, now(), $3)
			ON CONFLICT (key) DO UPDATE SET
				value = $2,
				updated_at = now(),
				updated_by = $3
		`
		_, err = h.db.Exec(query, key, valueJSON, userID)
		if err == nil {
			updated++
		}
	}

	// Log activity
	logActivity(h.db, c, "bulk_update", "settings", "", "Bulk updated settings", nil)

	c.JSON(http.StatusOK, gin.H{"updated": updated})
}

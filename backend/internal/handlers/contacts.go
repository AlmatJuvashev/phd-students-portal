package handlers

import (
	"net/http"
	"strings"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type ContactsHandler struct {
	svc *services.ContactService
}

func NewContactsHandler(svc *services.ContactService) *ContactsHandler {
	return &ContactsHandler{svc: svc}
}

type contactPayload struct {
	Name      map[string]string `json:"name"`
	Title     map[string]string `json:"title"`
	Email     *string           `json:"email"`
	Phone     *string           `json:"phone"`
	SortOrder *int              `json:"sort_order"`
	IsActive  *bool             `json:"is_active"`
}

func (h *ContactsHandler) PublicList(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant context required"})
		return
	}

	contacts, err := h.svc.ListPublic(c.Request.Context(), tenantID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, contacts)
}

func (h *ContactsHandler) AdminList(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant context required"})
		return
	}
	includeInactive := c.Query("all") == "true"
	
	contacts, err := h.svc.ListAdmin(c.Request.Context(), tenantID.(string), includeInactive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, contacts)
}

func (h *ContactsHandler) Create(c *gin.Context) {
	var req contactPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(req.Name) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	tenantIDStr := c.GetString("tenant_id")
	if tenantIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant context required"})
		return
	}

	sortOrder := 0
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	contact := models.Contact{
		Name:      models.LocalizedMap(req.Name),
		Title:     models.LocalizedMap(req.Title),
		Email:     req.Email,
		Phone:     req.Phone,
		SortOrder: sortOrder,
		IsActive:  isActive,
	}

	id, err := h.svc.Create(c.Request.Context(), tenantIDStr, contact)
	if err != nil {
		if strings.Contains(err.Error(), "contacts_name_key") {
			c.JSON(http.StatusConflict, gin.H{"error": "Contact already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *ContactsHandler) Update(c *gin.Context) {
	id := c.Param("id")
	tenantIDStr := c.GetString("tenant_id")
	
	var req contactPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]interface{})

	if len(req.Name) > 0 {
		updates["name"] = models.LocalizedMap(req.Name)
	}
	if req.Title != nil {
		updates["title"] = models.LocalizedMap(req.Title)
	}
	if req.Email != nil {
		updates["email"] = req.Email
	}
	if req.Phone != nil {
		updates["phone"] = req.Phone
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	if err := h.svc.Update(c.Request.Context(), tenantIDStr, id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *ContactsHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	tenantIDStr := c.GetString("tenant_id")
	
	if err := h.svc.Delete(c.Request.Context(), tenantIDStr, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}


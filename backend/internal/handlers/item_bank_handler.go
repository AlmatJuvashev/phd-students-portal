package handlers

import (
	"net/http"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/middleware"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type ItemBankHandler struct {
	svc *services.ItemBankService
}

func NewItemBankHandler(svc *services.ItemBankService) *ItemBankHandler {
	return &ItemBankHandler{svc: svc}
}

type updateBankRequest struct {
	Title          *string                 `json:"title"`
	Description    **string                `json:"description"`
	Subject        **string                `json:"subject"`
	BloomsTaxonomy **models.BloomsTaxonomy `json:"blooms_taxonomy"`
	IsPublic       *bool                   `json:"is_public"`
}

// CreateBank - POST /api/item-banks/banks
func (h *ItemBankHandler) CreateBank(c *gin.Context) {
	var b models.QuestionBank
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	b.TenantID = middleware.GetTenantID(c)
	b.CreatedBy = middleware.GetUserID(c)

	if err := h.svc.CreateBank(c.Request.Context(), &b); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, b)
}

// ListBanks - GET /api/item-banks/banks
func (h *ItemBankHandler) ListBanks(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	list, err := h.svc.ListBanks(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// UpdateBank - PUT /api/item-banks/banks/:bankId
func (h *ItemBankHandler) UpdateBank(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	bankID := c.Param("bankId")

	existing, err := h.svc.GetBank(c.Request.Context(), bankID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existing == nil || existing.TenantID != tenantID {
		c.JSON(http.StatusNotFound, gin.H{"error": "bank not found"})
		return
	}

	var req updateBankRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Subject != nil {
		existing.Subject = *req.Subject
	}
	if req.BloomsTaxonomy != nil {
		existing.BloomsTaxonomy = *req.BloomsTaxonomy
	}
	if req.IsPublic != nil {
		existing.IsPublic = *req.IsPublic
	}

	if err := h.svc.UpdateBank(c.Request.Context(), existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, existing)
}

// DeleteBank - DELETE /api/item-banks/banks/:bankId
func (h *ItemBankHandler) DeleteBank(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	bankID := c.Param("bankId")

	existing, err := h.svc.GetBank(c.Request.Context(), bankID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existing == nil || existing.TenantID != tenantID {
		c.JSON(http.StatusNotFound, gin.H{"error": "bank not found"})
		return
	}

	if err := h.svc.DeleteBank(c.Request.Context(), bankID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// CreateItem - POST /api/item-banks/banks/:bankId/items
func (h *ItemBankHandler) CreateItem(c *gin.Context) {
	bankID := c.Param("bankId")
	var item models.Question
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item.BankID = bankID

	if err := h.svc.CreateItem(c.Request.Context(), &item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

type updateItemRequest struct {
	Type              *models.QuestionType     `json:"type"`
	Stem              *string                  `json:"stem"`
	MediaURL          **string                 `json:"media_url"`
	PointsDefault     *float64                 `json:"points_default"`
	DifficultyLevel   **models.DifficultyLevel `json:"difficulty_level"`
	LearningOutcomeID **string                 `json:"learning_outcome_id"`
	Options           *[]models.QuestionOption `json:"options"`
}

// ListItems - GET /api/item-banks/banks/:bankId/items
func (h *ItemBankHandler) ListItems(c *gin.Context) {
	bankID := c.Param("bankId")
	list, err := h.svc.ListItems(c.Request.Context(), bankID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// UpdateItem - PUT /api/item-banks/banks/:bankId/items/:itemId
func (h *ItemBankHandler) UpdateItem(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	bankID := c.Param("bankId")
	itemID := c.Param("itemId")

	bank, err := h.svc.GetBank(c.Request.Context(), bankID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if bank == nil || bank.TenantID != tenantID {
		c.JSON(http.StatusNotFound, gin.H{"error": "bank not found"})
		return
	}

	existing, err := h.svc.GetItem(c.Request.Context(), itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existing == nil || existing.BankID != bankID {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	var req updateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Type != nil {
		existing.Type = *req.Type
	}
	if req.Stem != nil {
		existing.Stem = *req.Stem
	}
	if req.MediaURL != nil {
		existing.MediaURL = *req.MediaURL
	}
	if req.PointsDefault != nil {
		existing.PointsDefault = *req.PointsDefault
	}
	if req.DifficultyLevel != nil {
		existing.DifficultyLevel = *req.DifficultyLevel
	}
	if req.LearningOutcomeID != nil {
		existing.LearningOutcomeID = *req.LearningOutcomeID
	}
	if req.Options != nil {
		existing.Options = *req.Options
	}

	if err := h.svc.UpdateItem(c.Request.Context(), existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existing)
}

// DeleteItem - DELETE /api/item-banks/banks/:bankId/items/:itemId
func (h *ItemBankHandler) DeleteItem(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	bankID := c.Param("bankId")
	itemID := c.Param("itemId")

	bank, err := h.svc.GetBank(c.Request.Context(), bankID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if bank == nil || bank.TenantID != tenantID {
		c.JSON(http.StatusNotFound, gin.H{"error": "bank not found"})
		return
	}

	existing, err := h.svc.GetItem(c.Request.Context(), itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existing == nil || existing.BankID != bankID {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	if err := h.svc.DeleteItem(c.Request.Context(), itemID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

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

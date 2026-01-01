package handlers

import (
	"net/http"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/middleware"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type GovernanceHandler struct {
	svc *services.GovernanceService
}

func NewGovernanceHandler(svc *services.GovernanceService) *GovernanceHandler {
	return &GovernanceHandler{svc: svc}
}

// SubmitProposal - POST /api/governance/proposals
func (h *GovernanceHandler) SubmitProposal(c *gin.Context) {
	var p models.Proposal
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	p.TenantID = middleware.GetTenantID(c)
	p.RequesterID = middleware.GetUserID(c) // Uses the helper we defined earlier

	if err := h.svc.SubmitProposal(c.Request.Context(), &p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

// ListProposals - GET /api/governance/proposals?status=pending
func (h *GovernanceHandler) ListProposals(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	status := c.Query("status")
	
	list, err := h.svc.ListProposals(c.Request.Context(), tenantID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

type ReviewRequest struct {
	Status  string `json:"status" binding:"required,oneof=approved rejected"`
	Comment string `json:"comment"`
}

// ReviewProposal - POST /api/governance/proposals/:id/review
func (h *GovernanceHandler) ReviewProposal(c *gin.Context) {
	proposalID := c.Param("id")
	var req ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	reviewerID := middleware.GetUserID(c)
	
	if err := h.svc.ReviewProposal(c.Request.Context(), proposalID, reviewerID, req.Status, req.Comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "review recorded", "status": req.Status})
}

// GetProposal - GET /api/governance/proposals/:id
func (h *GovernanceHandler) GetProposal(c *gin.Context) {
	id := c.Param("id")
	p, err := h.svc.GetProposal(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Fetch reviews as well? Or separate endpoint? 
	// For detail view, usually good to include or fetch separately.
	// We'll keep it simple for now.
	c.JSON(http.StatusOK, p)
}

// ListReviews - GET /api/governance/proposals/:id/reviews
func (h *GovernanceHandler) ListReviews(c *gin.Context) {
	id := c.Param("id")
	list, err := h.svc.GetReviews(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

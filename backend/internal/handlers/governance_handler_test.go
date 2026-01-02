package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// HMockGovernanceRepo implements repository.GovernanceRepository
type HMockGovernanceRepo struct {
	mock.Mock
}

func (m *HMockGovernanceRepo) CreateProposal(ctx context.Context, p *models.Proposal) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}
func (m *HMockGovernanceRepo) GetProposal(ctx context.Context, id string) (*models.Proposal, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Proposal), args.Error(1)
}
func (m *HMockGovernanceRepo) ListProposals(ctx context.Context, tenantID string, statusFilter string) ([]models.Proposal, error) {
	args := m.Called(ctx, tenantID, statusFilter)
	return args.Get(0).([]models.Proposal), args.Error(1)
}
func (m *HMockGovernanceRepo) UpdateProposalStatus(ctx context.Context, id string, status string, currentStep int) error {
	args := m.Called(ctx, id, status, currentStep)
	return args.Error(0)
}
func (m *HMockGovernanceRepo) CreateReview(ctx context.Context, r *models.ProposalReview) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}
func (m *HMockGovernanceRepo) ListReviews(ctx context.Context, proposalID string) ([]models.ProposalReview, error) {
	args := m.Called(ctx, proposalID)
	return args.Get(0).([]models.ProposalReview), args.Error(1)
}

func setupGovernanceHandler() (*GovernanceHandler, *HMockGovernanceRepo) {
	repo := new(HMockGovernanceRepo)
	svc := services.NewGovernanceService(repo)
	return NewGovernanceHandler(svc), repo
}

// --- Tests ---

func TestGovernanceHandler_SubmitProposal(t *testing.T) {
	h, repo := setupGovernanceHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	body := `{"title":"Change Req", "type":"curriculum_change"}`
	c.Request, _ = http.NewRequest("POST", "/proposals", strings.NewReader(body))
	// Mock middleware values
	c.Set("tenant_id", "t1")
	c.Set("userID", "u1")
	c.Set("claims", jwt.MapClaims{"sub": "u1"})

	repo.On("CreateProposal", mock.Anything, mock.MatchedBy(func(p *models.Proposal) bool {
		return p.Title == "Change Req" && p.TenantID == "t1" && p.RequesterID == "u1"
	})).Return(nil)

	h.SubmitProposal(c)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGovernanceHandler_ListProposals(t *testing.T) {
	h, repo := setupGovernanceHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/proposals?status=pending", nil)
	c.Set("tenant_id", "t1")
	c.Set("userID", "u1")

	repo.On("ListProposals", mock.Anything, "t1", "pending").Return([]models.Proposal{}, nil)

	h.ListProposals(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGovernanceHandler_GetProposal(t *testing.T) {
	h, repo := setupGovernanceHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/proposals/p1", nil)
	c.Params = gin.Params{{Key: "id", Value: "p1"}}
	c.Set("userID", "u1")

	repo.On("GetProposal", mock.Anything, "p1").Return(&models.Proposal{ID: "p1", Title: "P1"}, nil)

	h.GetProposal(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGovernanceHandler_ReviewProposal(t *testing.T) {
	h, repo := setupGovernanceHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	body := `{"status":"approved", "comment":"LGTM"}`
	c.Request, _ = http.NewRequest("POST", "/proposals/p1/review", strings.NewReader(body))
	c.Params = gin.Params{{Key: "id", Value: "p1"}}
	c.Set("userID", "reviewer1")
	c.Set("claims", jwt.MapClaims{"sub": "reviewer1"})

	// Service Logic: Get, Check Status, Create Review, Update Status
	repo.On("GetProposal", mock.Anything, "p1").Return(&models.Proposal{ID: "p1", Status: "pending", CurrentStep: 1}, nil)
	
	repo.On("CreateReview", mock.Anything, mock.MatchedBy(func(r *models.ProposalReview) bool {
		return r.ProposalID == "p1" && r.ReviewerID == "reviewer1" && r.Status == "approved"
	})).Return(nil)

	// Approved increments step
	repo.On("UpdateProposalStatus", mock.Anything, "p1", "approved", 2).Return(nil)

	h.ReviewProposal(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGovernanceHandler_ListReviews(t *testing.T) {
	h, repo := setupGovernanceHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/proposals/p1/reviews", nil)
	c.Params = gin.Params{{Key: "id", Value: "p1"}}

	repo.On("ListReviews", mock.Anything, "p1").Return([]models.ProposalReview{}, nil)

	h.ListReviews(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

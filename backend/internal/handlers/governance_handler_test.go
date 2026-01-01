package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGovernanceRepo struct {
	mock.Mock
}
func (m *MockGovernanceRepo) CreateProposal(ctx context.Context, p *models.Proposal) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}
func (m *MockGovernanceRepo) GetProposal(ctx context.Context, id string) (*models.Proposal, error) { 
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Proposal), args.Error(1)
}
func (m *MockGovernanceRepo) ListProposals(ctx context.Context, t, s string) ([]models.Proposal, error) { return nil, nil }
func (m *MockGovernanceRepo) UpdateProposalStatus(ctx context.Context, id, s string, step int) error { return nil }
func (m *MockGovernanceRepo) CreateReview(ctx context.Context, r *models.ProposalReview) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}
func (m *MockGovernanceRepo) ListReviews(ctx context.Context, id string) ([]models.ProposalReview, error) { return nil, nil }

func TestGovernanceHandler_SubmitProposal(t *testing.T) {
	mockRepo := new(MockGovernanceRepo)
	svc := services.NewGovernanceService(mockRepo)
	h := handlers.NewGovernanceHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", "t-1")
		c.Set("userID", "u-1")
	})
	r.POST("/governance/proposals", h.SubmitProposal)

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("CreateProposal", mock.Anything, mock.MatchedBy(func(p *models.Proposal) bool {
			return p.Title == "Change" && p.RequesterID == "u-1"
		})).Return(nil)

		reqBody := map[string]interface{}{"title": "Change"}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/governance/proposals", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

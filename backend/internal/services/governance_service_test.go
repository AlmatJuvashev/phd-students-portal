package services

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGovernanceRepo
type MockGovernanceRepo struct {
	mock.Mock
}
func (m *MockGovernanceRepo) CreateProposal(ctx context.Context, p *models.Proposal) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}
func (m *MockGovernanceRepo) GetProposal(ctx context.Context, id string) (*models.Proposal, error) {
	args := m.Called(ctx, id)
	// Properly handle nil/type assertion
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Proposal), args.Error(1)
}
func (m *MockGovernanceRepo) ListProposals(ctx context.Context, t, s string) ([]models.Proposal, error) {
	args := m.Called(ctx, t, s)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.Proposal), args.Error(1)
}
func (m *MockGovernanceRepo) UpdateProposalStatus(ctx context.Context, id, s string, step int) error {
	args := m.Called(ctx, id, s, step)
	return args.Error(0)
}
func (m *MockGovernanceRepo) CreateReview(ctx context.Context, r *models.ProposalReview) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}
func (m *MockGovernanceRepo) ListReviews(ctx context.Context, id string) ([]models.ProposalReview, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.ProposalReview), args.Error(1)
}

func TestGovernanceService_ReviewProposal(t *testing.T) {
	mockRepo := new(MockGovernanceRepo)
	svc := NewGovernanceService(mockRepo)
	ctx := context.Background()
	
	proposalID := "prop-1"
	
	// 1. Approve Logic
	t.Run("Approve Pending Proposal", func(t *testing.T) {
		// Existing proposal state
		existing := &models.Proposal{ID: proposalID, Status: "pending", CurrentStep: 1}
		
		mockRepo.On("GetProposal", ctx, proposalID).Return(existing, nil).Once()
		mockRepo.On("CreateReview", ctx, mock.Anything).Return(nil).Once()
		// Should transition to approved, step increments
		mockRepo.On("UpdateProposalStatus", ctx, proposalID, "approved", 2).Return(nil).Once()

		err := svc.ReviewProposal(ctx, proposalID, "approver-1", "approved", "LGTM")
		assert.NoError(t, err)
	})

	// 2. Reject Logic
	t.Run("Reject Pending Proposal", func(t *testing.T) {
		existing := &models.Proposal{ID: proposalID, Status: "pending", CurrentStep: 1}

		mockRepo.On("GetProposal", ctx, proposalID).Return(existing, nil).Once()
		mockRepo.On("CreateReview", ctx, mock.Anything).Return(nil).Once()
		// Should transition to rejected, step stays or doesn't matter much (logic says step keeps same)
		mockRepo.On("UpdateProposalStatus", ctx, proposalID, "rejected", 1).Return(nil).Once()

		err := svc.ReviewProposal(ctx, proposalID, "approver-1", "rejected", "No")
		assert.NoError(t, err)
	})
}

func TestGovernanceService_SubmitProposal(t *testing.T) {
	mockRepo := new(MockGovernanceRepo)
	svc := NewGovernanceService(mockRepo)
	ctx := context.Background()

	t.Run("Submit Valid Proposal", func(t *testing.T) {
		proposal := &models.Proposal{
			Title: "New Policy",
		}
		
		mockRepo.On("CreateProposal", ctx, mock.MatchedBy(func(p *models.Proposal) bool {
			return p.Status == "pending" && p.CurrentStep == 1 && string(p.Data) == "{}"
		})).Return(nil)

		err := svc.SubmitProposal(ctx, proposal)
		assert.NoError(t, err)
	})

	t.Run("Submit Empty Title", func(t *testing.T) {
		proposal := &models.Proposal{
			Title: "",
		}
		err := svc.SubmitProposal(ctx, proposal)
		assert.Error(t, err)
		assert.Equal(t, "title is required", err.Error())
	})
}

func TestGovernanceService_ReadOps(t *testing.T) {
	mockRepo := new(MockGovernanceRepo)
	svc := NewGovernanceService(mockRepo)
	ctx := context.Background()

	t.Run("ListProposals", func(t *testing.T) {
		mockRepo.On("ListProposals", ctx, "t1", "pending").Return([]models.Proposal{{ID: "p1"}}, nil)
		res, err := svc.ListProposals(ctx, "t1", "pending")
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("GetProposal", func(t *testing.T) {
		mockRepo.On("GetProposal", ctx, "p1").Return(&models.Proposal{ID: "p1"}, nil)
		res, err := svc.GetProposal(ctx, "p1")
		assert.NoError(t, err)
		assert.Equal(t, "p1", res.ID)
	})

	t.Run("GetReviews", func(t *testing.T) {
		mockRepo.On("ListReviews", ctx, "p1").Return([]models.ProposalReview{{ID: "r1"}}, nil)
		res, err := svc.GetReviews(ctx, "p1")
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})
}

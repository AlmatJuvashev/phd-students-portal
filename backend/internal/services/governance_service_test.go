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
func (m *MockGovernanceRepo) ListProposals(ctx context.Context, t, s string) ([]models.Proposal, error) { return nil, nil }
func (m *MockGovernanceRepo) UpdateProposalStatus(ctx context.Context, id, s string, step int) error {
	args := m.Called(ctx, id, s, step)
	return args.Error(0)
}
func (m *MockGovernanceRepo) CreateReview(ctx context.Context, r *models.ProposalReview) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}
func (m *MockGovernanceRepo) ListReviews(ctx context.Context, id string) ([]models.ProposalReview, error) { return nil, nil }

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

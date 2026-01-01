package services

import (
	"context"
	"errors"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type GovernanceService struct {
	repo repository.GovernanceRepository
}

func NewGovernanceService(repo repository.GovernanceRepository) *GovernanceService {
	return &GovernanceService{repo: repo}
}

// SubmitProposal creates a new change request.
func (s *GovernanceService) SubmitProposal(ctx context.Context, p *models.Proposal) error {
	if p.Title == "" {
		return errors.New("title is required")
	}
	p.Status = "pending"
	p.CurrentStep = 1
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	
	// Default Data to empty JSON object if null
	if p.Data == nil {
		p.Data = []byte("{}")
	}

	return s.repo.CreateProposal(ctx, p)
}

// ReviewProposal records a reviewer's decision.
func (s *GovernanceService) ReviewProposal(ctx context.Context, proposalID, reviewerID, status, comment string) error {
	// Validate current status
	p, err := s.repo.GetProposal(ctx, proposalID)
	if err != nil {
		return err
	}
	if p.Status != "pending" {
		return errors.New("proposal is not in pending state")
	}

	// Record Review
	review := &models.ProposalReview{
		ProposalID: proposalID,
		ReviewerID: reviewerID,
		Status:     status,
		Comment:    comment,
	}
	if err := s.repo.CreateReview(ctx, review); err != nil {
		return err
	}

	// Update Proposal Status
	// Logic: If Rejected -> Main Status Rejected. 
	// If Approved -> For now, immediately Approved. In future, check steps.
	newStatus := "pending"
	newStep := p.CurrentStep

	if status == "rejected" {
		newStatus = "rejected"
	} else if status == "approved" {
		newStatus = "approved"
		newStep++ 
	} else {
		return errors.New("invalid status, must be approved or rejected")
	}

	return s.repo.UpdateProposalStatus(ctx, proposalID, newStatus, newStep)
}

func (s *GovernanceService) ListProposals(ctx context.Context, tenantID, status string) ([]models.Proposal, error) {
	return s.repo.ListProposals(ctx, tenantID, status)
}

func (s *GovernanceService) GetProposal(ctx context.Context, id string) (*models.Proposal, error) {
	return s.repo.GetProposal(ctx, id)
}

func (s *GovernanceService) GetReviews(ctx context.Context, proposalID string) ([]models.ProposalReview, error) {
	return s.repo.ListReviews(ctx, proposalID)
}

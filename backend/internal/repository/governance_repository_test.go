package repository

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/assert"
)

func TestSQLGovernanceRepository_CreateProposal(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLGovernanceRepository(sqlxDB)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		p := &models.Proposal{
			TenantID:    "t-1",
			RequesterID: "u-1",
			Type:        "curriculum",
			Title:       "New Course",
			Status:      "pending",
			Data:        types.JSONText("{}"),
			CurrentStep: 1,
		}

		rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow("prop-1", time.Now(), time.Now())

		mock.ExpectQuery("INSERT INTO proposals").
			WithArgs(p.TenantID, p.RequesterID, p.Type, p.TargetID, p.Title, p.Description, p.Status, p.Data, p.CurrentStep).
			WillReturnRows(rows)

		err := repo.CreateProposal(ctx, p)
		assert.NoError(t, err)
		assert.Equal(t, "prop-1", p.ID)
	})
}

func TestSQLGovernanceRepository_CreateReview(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLGovernanceRepository(sqlxDB)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		rev := &models.ProposalReview{
			ProposalID: "prop-1",
			ReviewerID: "rev-1",
			Status:     "approved",
			Comment:    "LGTM",
		}

		rows := sqlmock.NewRows([]string{"id", "created_at"}).
			AddRow("rev-id-1", time.Now())

		mock.ExpectQuery("INSERT INTO proposal_reviews").
			WithArgs(rev.ProposalID, rev.ReviewerID, rev.Status, rev.Comment).
			WillReturnRows(rows)

		err := repo.CreateReview(ctx, rev)
		assert.NoError(t, err)
		assert.Equal(t, "rev-id-1", rev.ID)
	})
}

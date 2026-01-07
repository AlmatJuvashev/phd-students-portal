package services

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuditService_CreateLearningOutcome(t *testing.T) {
	ctx := context.Background()
	repo := new(MockAuditRepository)
	curRepo := new(MockCurriculumRepository)
	svc := NewAuditService(repo, curRepo)

	outcome := &models.LearningOutcome{TenantID: "t1", ID: "o1", Code: "PLO1"}
	repo.On("CreateLearningOutcome", ctx, outcome).Return(nil)
	repo.On("LogCurriculumChange", ctx, mock.MatchedBy(func(log *models.CurriculumChangeLog) bool {
		return log.EntityType == "outcome" && log.Action == "created" && log.EntityID == "o1"
	})).Return(nil)

	err := svc.CreateLearningOutcome(ctx, outcome, "u1")
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestAuditService_UpdateLearningOutcome(t *testing.T) {
	ctx := context.Background()
	repo := new(MockAuditRepository)
	curRepo := new(MockCurriculumRepository)
	svc := NewAuditService(repo, curRepo)

	outcome := &models.LearningOutcome{TenantID: "t1", ID: "o1", Code: "PLO1_v2"}
	repo.On("GetLearningOutcome", ctx, "o1").Return(&models.LearningOutcome{ID: "o1", Code: "PLO1"}, nil)
	repo.On("UpdateLearningOutcome", ctx, outcome).Return(nil)
	repo.On("LogCurriculumChange", ctx, mock.MatchedBy(func(log *models.CurriculumChangeLog) bool {
		return log.EntityType == "outcome" && log.Action == "updated" && log.EntityID == "o1"
	})).Return(nil)

	err := svc.UpdateLearningOutcome(ctx, outcome, "u1")
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestAuditService_GenerateProgramSummary(t *testing.T) {
	ctx := context.Background()
	repo := new(MockAuditRepository)
	curRepo := new(MockCurriculumRepository)
	svc := NewAuditService(repo, curRepo)

	programID := "p1"
	tenantID := "t1"

	curRepo.On("GetProgram", ctx, programID).Return(&models.Program{ID: programID, Credits: 120}, nil)
	curRepo.On("ListCourses", ctx, tenantID, &programID).Return([]models.Course{{ID: "c1", Credits: 5}, {ID: "c2", Credits: 3}}, nil)
	repo.On("ListLearningOutcomes", ctx, tenantID, &programID, (*string)(nil)).Return([]models.LearningOutcome{{ID: "lo1"}}, nil)

	report, err := svc.GenerateProgramSummary(ctx, tenantID, programID)
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, 8, report.TotalCredits)
	assert.Equal(t, 2, report.TotalCourses)
	assert.Equal(t, 1, report.TotalOutcomes)
}

func TestAuditService_Others(t *testing.T) {
	ctx := context.Background()
	repo := new(MockAuditRepository)
	curRepo := new(MockCurriculumRepository)
	svc := NewAuditService(repo, curRepo)

	t.Run("ListLearningOutcomes", func(t *testing.T) {
		repo.On("ListLearningOutcomes", ctx, "t1", mock.Anything, mock.Anything).Return([]models.LearningOutcome{{ID: "lo1"}}, nil)
		res, err := svc.ListLearningOutcomes(ctx, "t1", nil, nil)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("DeleteLearningOutcome", func(t *testing.T) {
		repo.On("GetLearningOutcome", ctx, "o1").Return(&models.LearningOutcome{ID: "o1"}, nil)
		repo.On("DeleteLearningOutcome", ctx, "o1").Return(nil)
		repo.On("LogCurriculumChange", ctx, mock.Anything).Return(nil)
		err := svc.DeleteLearningOutcome(ctx, "t1", "o1", "u1")
		assert.NoError(t, err)
	})

	t.Run("LinkOutcomeToAssessment", func(t *testing.T) {
		repo.On("LinkOutcomeToAssessment", ctx, "o1", "nd1", 1.0).Return(nil)
		err := svc.LinkOutcomeToAssessment(ctx, "o1", "nd1", 1.0)
		assert.NoError(t, err)
	})

	t.Run("ListCurriculumChanges", func(t *testing.T) {
		repo.On("ListCurriculumChanges", ctx, mock.Anything).Return([]models.CurriculumChangeLog{{ID: "l1"}}, nil)
		res, err := svc.ListCurriculumChanges(ctx, models.AuditReportFilter{})
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})
}

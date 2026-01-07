package services

import (
	"context"
	"errors"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAnalyticsService_GetMonitorMetrics(t *testing.T) {
	mockRepo := new(MockAnalyticsRepository)
	service := NewAnalyticsService(mockRepo, nil, nil, nil)

	ctx := context.Background()
	filter := models.FilterParams{TenantID: "tenant-1"}

	t.Run("Success", func(t *testing.T) {
		// Mock behavior
		mockRepo.On("GetTotalStudents", ctx, filter).Return(100, nil)
		mockRepo.On("GetNodeCompletionCount", ctx, "S1_antiplag", filter).Return(50, nil)
		mockRepo.On("GetDurationForNodes", ctx, mock.Anything, filter).Return([]float64{10.0, 20.0, 30.0}, nil)
		mockRepo.On("GetBottleneck", ctx, filter).Return("S3_thesis", 15, nil)
		mockRepo.On("GetProfileFlagCount", ctx, "years_since_graduation", 3.0, filter).Return(5, nil)

		metrics, err := service.GetMonitorMetrics(ctx, filter)

		assert.NoError(t, err)
		assert.Equal(t, 50.0, metrics.ComplianceRate)    // 50/100 * 100
		assert.Equal(t, 20.0, metrics.StageMedianDays) // Median of 10,20,30 is 20
		assert.Equal(t, "S3_thesis", metrics.BottleneckNodeID)
		assert.Equal(t, 15, metrics.BottleneckCount)
		assert.Equal(t, 5, metrics.ProfileFlagCount)
		
		mockRepo.AssertExpectations(t)
	})

	t.Run("EmptyTotalStudents", func(t *testing.T) {
		// If total students is 0, should return empty metrics without calling other methods
		mockRepo.ExpectedCalls = nil // Clear previous
		mockRepo.On("GetTotalStudents", ctx, filter).Return(0, nil)

		metrics, err := service.GetMonitorMetrics(ctx, filter)

		assert.NoError(t, err)
		assert.Equal(t, 0.0, metrics.ComplianceRate)
		assert.Equal(t, 0.0, metrics.StageMedianDays)
		
		mockRepo.AssertNotCalled(t, "GetNodeCompletionCount")
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetTotalStudents", ctx, filter).Return(0, errors.New("db error"))

		_, err := service.GetMonitorMetrics(ctx, filter)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
	
	t.Run("MedianEvenCount", func(t *testing.T) {
		// Test median logic for even number of elements
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetTotalStudents", ctx, filter).Return(10, nil)
		mockRepo.On("GetNodeCompletionCount", ctx, mock.Anything, filter).Return(5, nil)
		// 10, 20, 30, 40 -> Median is (20+30)/2 = 25
		mockRepo.On("GetDurationForNodes", ctx, mock.Anything, filter).Return([]float64{10.0, 20.0, 30.0, 40.0}, nil)
		mockRepo.On("GetBottleneck", ctx, filter).Return("", 0, nil)
		mockRepo.On("GetProfileFlagCount", ctx, mock.Anything, mock.Anything, filter).Return(0, nil)

		metrics, err := service.GetMonitorMetrics(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, 25.0, metrics.StageMedianDays)
	})
}

func TestAnalyticsService_HighRisk(t *testing.T) {
	mockRepo := new(MockAnalyticsRepository)
	service := NewAnalyticsService(mockRepo, nil, nil, nil)
	ctx := context.Background()

	t.Run("GetHighRiskStudents", func(t *testing.T) {
		expected := []models.RiskSnapshot{{StudentID: "s1", RiskScore: 0.85}}
		mockRepo.On("GetHighRiskStudents", ctx, 0.7).Return(expected, nil)

		res, err := service.GetHighRiskStudents(ctx, 0.7)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}

func TestAnalyticsService_PassThroughs(t *testing.T) {
	mockRepo := new(MockAnalyticsRepository)
	service := NewAnalyticsService(mockRepo, nil, nil, nil)
	ctx := context.Background()

	mockRepo.On("GetStudentsByStage", ctx).Return([]models.StudentStageStats{}, nil)
	mockRepo.On("GetAdvisorLoad", ctx).Return([]models.AdvisorLoadStats{}, nil)
	mockRepo.On("GetOverdueTasks", ctx).Return([]models.OverdueTaskStats{}, nil)

	_, err := service.GetStudentsByStage(ctx)
	assert.NoError(t, err)

	_, err = service.GetAdvisorLoad(ctx)
	assert.NoError(t, err)

	_, err = service.GetOverdueTasks(ctx)
	assert.NoError(t, err)
}

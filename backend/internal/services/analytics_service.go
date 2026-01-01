package services

import (
	"context"
	"sort"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type AnalyticsConfig struct {
	ComplianceNodeID string   // e.g. "S1_antiplag"
	StageNodeIDs     []string // e.g. Nodes belonging to W2 for median calculation
	ProfileFlagKey   string   // e.g. "years_since_graduation" or "is_grant"
	ProfileFlagMin   float64  // e.g. 3.0
}

func DefaultAnalyticsConfig() *AnalyticsConfig {
	return &AnalyticsConfig{
		ComplianceNodeID: "S1_antiplag",
		// Placeholder for W2 nodes. In a real dynamic system, these would be fetched from a Config Service or DB.
		StageNodeIDs:     []string{"S2_topic_approval", "S2_advisor_approval", "S2_plan_approval", "S2_research_methodology"},
		ProfileFlagKey:   "years_since_graduation",
		ProfileFlagMin:   3.0,
	}
}

type AnalyticsService struct {
	repo   repository.AnalyticsRepository
	config *AnalyticsConfig
}

func NewAnalyticsService(repo repository.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{
		repo:   repo,
		config: DefaultAnalyticsConfig(),
	}
}

// WithConfig allows overriding default config (e.g. per tenant/program in future)
func (s *AnalyticsService) WithConfig(cfg *AnalyticsConfig) *AnalyticsService {
	s.config = cfg
	return s
}

func (s *AnalyticsService) GetMonitorMetrics(ctx context.Context, filter models.FilterParams) (*models.MonitorMetrics, error) {
	metrics := &models.MonitorMetrics{}

	// 1. Compliance Rate (Antiplag)
	// Get total students matching filter
	// Optimization: If total is 0, return 0.
	total, err := s.repo.GetTotalStudents(ctx, filter)
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return metrics, nil // Empty metrics
	}

	// Get count of students who completed the compliance node
	compCount, err := s.repo.GetNodeCompletionCount(ctx, s.config.ComplianceNodeID, filter)
	if err != nil {
		return nil, err
	}
	metrics.ComplianceRate = (float64(compCount) / float64(total)) * 100.0

	// 2. Stage Median Days (W2)
	durations, err := s.repo.GetDurationForNodes(ctx, s.config.StageNodeIDs, filter)
	if err != nil {
		return nil, err
	}
	metrics.StageMedianDays = calculateMedian(durations)

	// 3. Bottleneck
	bNode, bCount, err := s.repo.GetBottleneck(ctx, filter)
	if err != nil {
		return nil, err
	}
	metrics.BottleneckNodeID = bNode
	metrics.BottleneckCount = bCount

	// 4. Profile Flag (RP Required)
	// Only if ProfileFlagKey is set
	if s.config.ProfileFlagKey != "" {
		rpCount, err := s.repo.GetProfileFlagCount(ctx, s.config.ProfileFlagKey, s.config.ProfileFlagMin, filter)
		if err != nil {
			return nil, err
		}
		metrics.ProfileFlagCount = rpCount
	}

	return metrics, nil
}

func calculateMedian(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sort.Float64s(values)
	n := len(values)
	if n%2 == 1 {
		return values[n/2]
	}
	return (values[n/2-1] + values[n/2]) / 2.0
}

// --- Legacy Delegates ---

func (s *AnalyticsService) GetStudentsByStage(ctx context.Context) ([]models.StudentStageStats, error) {
	return s.repo.GetStudentsByStage(ctx)
}

func (s *AnalyticsService) GetAdvisorLoad(ctx context.Context) ([]models.AdvisorLoadStats, error) {
	return s.repo.GetAdvisorLoad(ctx)
}

func (s *AnalyticsService) GetOverdueTasks(ctx context.Context) ([]models.OverdueTaskStats, error) {
	return s.repo.GetOverdueTasks(ctx)
}

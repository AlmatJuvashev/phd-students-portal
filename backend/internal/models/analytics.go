package models

type StudentStageStats struct {
	Stage string `db:"stage" json:"stage"`
	Count int    `db:"count" json:"count"`
}

type AdvisorLoadStats struct {
	AdvisorName string `db:"advisor_name" json:"advisor_name"`
	StudentCount int   `db:"student_count" json:"student_count"`
}

type OverdueTaskStats struct {
	NodeID string `db:"node_id" json:"node_id"`
	Count  int    `db:"count" json:"count"`
}

// MonitorMetrics aggregates high-level dashboard stats
type MonitorMetrics struct {
	TotalStudentsCount int     `json:"total_students_count"` // Total students in filtered population
	ComplianceRate     float64 `json:"compliance_rate"`      // Generic "Antiplag" etc.
	StageMedianDays    float64 `json:"stage_median_days"`    // Generic "W2" etc.
	BottleneckNodeID   string  `json:"bottleneck_node_id"`
	BottleneckCount    int     `json:"bottleneck_count"`
	ProfileFlagCount   int     `json:"profile_flag_count"` // Generic "RP Required"
}

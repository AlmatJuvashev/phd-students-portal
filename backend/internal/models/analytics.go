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

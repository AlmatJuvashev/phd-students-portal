package models

// Program moved to curriculum.go
// Cohort moved to curriculum.go


type Specialty struct {
	ID         string   `db:"id" json:"id"`
	Name       string   `db:"name" json:"name"`
	Code       string   `db:"code" json:"code"`
	ProgramIDs []string `json:"program_ids"` // Multiple programs
	IsActive   bool     `db:"is_active" json:"is_active"`
	CreatedAt  string   `db:"created_at" json:"created_at"`
	UpdatedAt  string   `db:"updated_at" json:"updated_at"`
}

type Department struct {
	ID        string `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	Code      string `db:"code" json:"code"`
	IsActive  bool   `db:"is_active" json:"is_active"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

package models

import (
	"time"

	"github.com/jmoiron/sqlx/types"
)

// --- Enums ---

type QuestionType string

const (
	QuestionTypeMCQ       QuestionType = "MCQ"
	QuestionTypeMRQ       QuestionType = "MRQ"
	QuestionTypeTrueFalse QuestionType = "TRUE_FALSE"
	QuestionTypeText      QuestionType = "TEXT"
	QuestionTypeLikert    QuestionType = "LIKERT"
)

type DifficultyLevel string

const (
	DifficultyEasy   DifficultyLevel = "EASY"
	DifficultyMedium DifficultyLevel = "MEDIUM"
	DifficultyHard   DifficultyLevel = "HARD"
)

type BloomsTaxonomy string

const (
	BloomsKnowledge     BloomsTaxonomy = "KNOWLEDGE"
	BloomsComprehension BloomsTaxonomy = "COMPREHENSION"
	BloomsApplication   BloomsTaxonomy = "APPLICATION"
	BloomsAnalysis      BloomsTaxonomy = "ANALYSIS"
	BloomsSynthesis     BloomsTaxonomy = "SYNTHESIS"
	BloomsEvaluation    BloomsTaxonomy = "EVALUATION"
)

type GradingPolicy string

const (
	GradingPolicyAutomatic    GradingPolicy = "AUTOMATIC"
	GradingPolicyManualReview GradingPolicy = "MANUAL_REVIEW"
)

type AttemptStatus string

const (
	AttemptStatusInProgress AttemptStatus = "IN_PROGRESS"
	AttemptStatusSubmitted  AttemptStatus = "SUBMITTED"
	AttemptStatusGraded     AttemptStatus = "GRADED"
)

// --- Structs ---

type QuestionBank struct {
	ID             string         `db:"id" json:"id"`
	TenantID       string         `db:"tenant_id" json:"tenant_id"`
	Title          string         `db:"title" json:"title"`
	Description    *string        `db:"description" json:"description,omitempty"`
	Subject        *string        `db:"subject" json:"subject,omitempty"`
	BloomsTaxonomy *BloomsTaxonomy `db:"blooms_taxonomy" json:"blooms_taxonomy,omitempty"`
	IsPublic       bool           `db:"is_public" json:"is_public"`
	CreatedBy      string         `db:"created_by" json:"created_by"`
	CreatedAt      time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at" json:"updated_at"`
}

type Question struct {
	ID                string           `db:"id" json:"id"`
	BankID            string           `db:"bank_id" json:"bank_id"`
	Type              QuestionType     `db:"type" json:"type"`
	Stem              string           `db:"stem" json:"stem"`
	MediaURL          *string          `db:"media_url" json:"media_url,omitempty"`
	PointsDefault     float64          `db:"points_default" json:"points_default"`
	DifficultyLevel   *DifficultyLevel `db:"difficulty_level" json:"difficulty_level,omitempty"`
	LearningOutcomeID *string          `db:"learning_outcome_id" json:"learning_outcome_id,omitempty"`
	CreatedAt         time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time        `db:"updated_at" json:"updated_at"`

	// Joins
	Options []QuestionOption `db:"-" json:"options,omitempty"`
}

type QuestionOption struct {
	ID         string  `db:"id" json:"id"`
	QuestionID string  `db:"question_id" json:"question_id"`
	Text       string  `db:"text" json:"text"`
	IsCorrect  bool    `db:"is_correct" json:"is_correct"` // Hidden from student during exam
	SortOrder  int     `db:"sort_order" json:"sort_order"`
	Feedback   *string `db:"feedback" json:"feedback,omitempty"`
}

type Assessment struct {
	ID               string        `db:"id" json:"id"`
	TenantID         string        `db:"tenant_id" json:"tenant_id"`
	CourseOfferingID string        `db:"course_offering_id" json:"course_offering_id"`
	Title            string        `db:"title" json:"title"`
	Description      *string       `db:"description" json:"description,omitempty"`
	TimeLimitMinutes *int          `db:"time_limit_minutes" json:"time_limit_minutes,omitempty"`
	AvailableFrom    *time.Time    `db:"available_from" json:"available_from,omitempty"`
	AvailableUntil   *time.Time    `db:"available_until" json:"available_until,omitempty"`
	ShuffleQuestions bool          `db:"shuffle_questions" json:"shuffle_questions"`
	GradingPolicy    GradingPolicy    `db:"grading_policy" json:"grading_policy"`
	SecuritySettings types.JSONText   `db:"security_settings" json:"security_settings"` // Stores SecuritySettings struct
	PassingScore     float64          `db:"passing_score" json:"passing_score"`
	CreatedBy        string        `db:"created_by" json:"created_by"`
	CreatedAt        time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time     `db:"updated_at" json:"updated_at"`

	// Joins
	Sections []AssessmentSection `db:"-" json:"sections,omitempty"`
}

type AssessmentSection struct {
	ID           string           `db:"id" json:"id"`
	AssessmentID string           `db:"assessment_id" json:"assessment_id"`
	Title        *string          `db:"title" json:"title,omitempty"`
	Instructions *string          `db:"instructions" json:"instructions,omitempty"`
	SortOrder    int              `db:"sort_order" json:"sort_order"`
	
	Items        []AssessmentItem `db:"-" json:"items,omitempty"`
}

type AssessmentItem struct {
	ID              string  `db:"id" json:"id"`
	AssessmentID    string  `db:"assessment_id" json:"assessment_id"`
	SectionID       *string `db:"section_id" json:"section_id,omitempty"`
	QuestionID      string  `db:"question_id" json:"question_id"`
	PointsOverride  *float64 `db:"points_override" json:"points_override,omitempty"`
	SortOrder       int     `db:"sort_order" json:"sort_order"`

	// Joins
	Question *Question `db:"-" json:"question,omitempty"`
}

type AssessmentAttempt struct {
	ID           string        `db:"id" json:"id"`
	AssessmentID string        `db:"assessment_id" json:"assessment_id"`
	StudentID    string        `db:"student_id" json:"student_id"`
	StartedAt    time.Time     `db:"started_at" json:"started_at"`
	FinishedAt   *time.Time    `db:"finished_at" json:"finished_at,omitempty"`
	Score        float64       `db:"score" json:"score"`
	Status       AttemptStatus `db:"status" json:"status"`
}

type ItemResponse struct {
	ID               string    `db:"id" json:"id"`
	AttemptID        string    `db:"attempt_id" json:"attempt_id"`
	QuestionID       string    `db:"question_id" json:"question_id"`
	SelectedOptionID *string   `db:"selected_option_id" json:"selected_option_id,omitempty"`
	TextResponse     *string   `db:"text_response" json:"text_response,omitempty"`
	Score            float64   `db:"score" json:"score"`
	IsCorrect        bool      `db:"is_correct" json:"is_correct"`
	GradedAt         *time.Time `db:"graded_at" json:"graded_at,omitempty"`
}

// --- Proctoring Models ---

type ProctoringEventType string

const (
	ProctoringEventTabSwitch      ProctoringEventType = "TAB_SWITCH"
	ProctoringEventWindowBlur     ProctoringEventType = "WINDOW_BLUR"
	ProctoringEventFullscreenExit ProctoringEventType = "FULLSCREEN_EXIT"
	ProctoringEventMouseLeave     ProctoringEventType = "MOUSE_LEAVE"
	ProctoringEventDeviceChange   ProctoringEventType = "DEVICE_CHANGE"
)

type SecuritySettings struct {
	FullScreenMode    bool `json:"full_screen_mode"`
	TrackTabSwitches  bool `json:"track_tab_switches"`
	MaxViolations     int  `json:"max_violations"`      // 0 = unlimited
	AutoSubmitOnLimit bool `json:"auto_submit_on_limit"`
	RecordWebcam      bool `json:"record_webcam"`       // Future placeholder
}

type ProctoringLog struct {
	ID         string              `db:"id" json:"id"`
	AttemptID  string              `db:"attempt_id" json:"attempt_id"`
	EventType  ProctoringEventType `db:"event_type" json:"event_type"`
	OccurredAt time.Time           `db:"occurred_at" json:"occurred_at"`
	Metadata   types.JSONText      `db:"metadata" json:"metadata,omitempty"`
}

// Ensure Enums implement Value/Scan if needed (simplified here as strings usually work with sqlx/pq if cast)

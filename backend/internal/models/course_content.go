package models

import "time"

// CourseModule represents a high-level section or chapter in a course
type CourseModule struct {
	ID        string    `db:"id" json:"id"`
	CourseID  string    `db:"course_id" json:"course_id"`
	Title     string    `db:"title" json:"title"`
	Order     int       `db:"sort_order" json:"order"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	
	// Relationships
	Lessons []CourseLesson `json:"lessons,omitempty"`
}

// CourseLesson represents a specific topic within a module
type CourseLesson struct {
	ID        string    `db:"id" json:"id"`
	ModuleID  string    `db:"module_id" json:"module_id"`
	Title     string    `db:"title" json:"title"`
	Order     int       `db:"sort_order" json:"order"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	// Relationships
	Activities []CourseActivity `json:"activities,omitempty"`
}

// CourseActivity represents a specific content item (Text, Video, Quiz, etc.)
type CourseActivity struct {
	ID          string    `db:"id" json:"id"`
	LessonID    string    `db:"lesson_id" json:"lesson_id"`
	Type        string    `db:"type" json:"type"` // text, video, quiz, survey, assignment
	Title       string    `db:"title" json:"title"`
	Order       int       `db:"sort_order" json:"order"`
	Points      int       `db:"points" json:"points"`
	IsOptional  bool      `db:"is_optional" json:"is_optional"`
	Content     string    `db:"content" json:"content"` // JSONB: stores videoUrls, quizConfig, text content etc.
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// --- Quiz & Survey Schemas (Stored in Content JSONB) ---

type QuizConfig struct {
	TimeLimitMinutes int            `json:"timeLimit"`
	PassingScore     int            `json:"passingScore"`
	ShuffleQuestions bool           `json:"shuffleQuestions"`
	ShowResults      bool           `json:"showResults"`
	Questions        []QuizQuestion `json:"questions"`
}

type QuizQuestion struct {
	ID                string       `json:"id"`
	Type              string       `json:"type"` // multiple_choice, multi_select, etc.
	Text              string       `json:"text"`
	Points            int          `json:"points"`
	Hint              string       `json:"hint,omitempty"`
	FeedbackCorrect   string       `json:"feedbackCorrect,omitempty"`
	FeedbackIncorrect string       `json:"feedbackIncorrect,omitempty"`
	Options           []QuizOption `json:"options,omitempty"`
	MatrixRows        []string     `json:"matrixRows,omitempty"`
	MatrixCols        []string     `json:"matrixCols,omitempty"`
}

type QuizOption struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	IsCorrect bool   `json:"isCorrect"`
}

type SurveyConfig struct {
	Anonymous       bool             `json:"anonymous"`
	ShowProgressBar bool             `json:"showProgressBar"`
	Questions       []SurveyQuestion `json:"questions"`
}

type SurveyQuestion struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"` // rating_stars, scale_10, etc.
	Text       string         `json:"text"`
	Required   bool           `json:"required"`
	Options    []SurveyOption `json:"options,omitempty"` // No IsCorrect
	MatrixRows []string       `json:"matrixRows,omitempty"`
	MatrixCols []string       `json:"matrixCols,omitempty"`
}

type SurveyOption struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

package repository

import (
	"context"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type AssessmentRepository interface {
	// Question Banks
	CreateQuestionBank(ctx context.Context, bank models.QuestionBank) (*models.QuestionBank, error)
	GetQuestionBank(ctx context.Context, id string) (*models.QuestionBank, error)
	ListQuestionBanks(ctx context.Context, tenantID string) ([]models.QuestionBank, error)
	UpdateQuestionBank(ctx context.Context, bank models.QuestionBank) error
	DeleteQuestionBank(ctx context.Context, id string) error

	// Questions
	CreateQuestion(ctx context.Context, q models.Question) (*models.Question, error)
	GetQuestion(ctx context.Context, id string) (*models.Question, error)
	ListQuestionsByBank(ctx context.Context, bankID string) ([]models.Question, error)
	UpdateQuestion(ctx context.Context, q models.Question) error
	DeleteQuestion(ctx context.Context, id string) error

	// Assessments
	CreateAssessment(ctx context.Context, a models.Assessment) (*models.Assessment, error)
	GetAssessment(ctx context.Context, id string) (*models.Assessment, error)

	// Attempts & Security
	CreateAttempt(ctx context.Context, attempt models.AssessmentAttempt) (*models.AssessmentAttempt, error)
	SaveItemResponse(ctx context.Context, response models.ItemResponse) error
	CompleteAttempt(ctx context.Context, attemptID string, score float64) error
	GetAttempt(ctx context.Context, id string) (*models.AssessmentAttempt, error)
	ListResponses(ctx context.Context, attemptID string) ([]models.ItemResponse, error)
	LogProctoringEvent(ctx context.Context, log models.ProctoringLog) error
	CountProctoringEvents(ctx context.Context, attemptID string) (int, error)

	// Helper to fetch full test content for grading/taking
	GetAssessmentQuestions(ctx context.Context, assessmentID string) ([]models.Question, error)
}

type SQLAssessmentRepository struct {
	db *sqlx.DB
}

func NewSQLAssessmentRepository(db *sqlx.DB) AssessmentRepository {
	return &SQLAssessmentRepository{db: db}
}

// ... (Question Banks code remains same)

// --- Helper Implementation ---

func (r *SQLAssessmentRepository) GetAssessmentQuestions(ctx context.Context, assessmentID string) ([]models.Question, error) {
	// 1. Fetch Questions linked to Assessment
	query := `
		SELECT q.* 
		FROM questions q
		JOIN assessment_items ai ON ai.question_id = q.id
		WHERE ai.assessment_id = $1
		ORDER BY ai.sort_order
	`
	var questions []models.Question
	err := r.db.SelectContext(ctx, &questions, query, assessmentID)
	if err != nil {
		return nil, err
	}
	if len(questions) == 0 {
		return nil, nil
	}

	// 2. Fetch Options for these questions
	questionIDs := make([]string, len(questions))
	for i, q := range questions {
		questionIDs[i] = q.ID
	}

	optsQuery, args, err := sqlx.In(`SELECT * FROM question_options WHERE question_id IN (?) ORDER BY sort_order`, questionIDs)
	if err != nil {
		return nil, err
	}
	optsQuery = r.db.Rebind(optsQuery)

	var options []models.QuestionOption
	err = r.db.SelectContext(ctx, &options, optsQuery, args...)
	if err != nil {
		return nil, err
	}

	// 3. Map Options to Questions
	optMap := make(map[string][]models.QuestionOption)
	for _, opt := range options {
		optMap[opt.QuestionID] = append(optMap[opt.QuestionID], opt)
	}

	for i := range questions {
		questions[i].Options = optMap[questions[i].ID]
	}

	return questions, nil
}

func (r *SQLAssessmentRepository) ListResponses(ctx context.Context, attemptID string) ([]models.ItemResponse, error) {
	var responses []models.ItemResponse
	err := r.db.SelectContext(ctx, &responses, `SELECT * FROM item_responses WHERE attempt_id = $1`, attemptID)
	return responses, err
}

// --- Question Banks ---

func (r *SQLAssessmentRepository) CreateQuestionBank(ctx context.Context, bank models.QuestionBank) (*models.QuestionBank, error) {
	err := r.db.QueryRowxContext(ctx, `
		INSERT INTO question_banks (tenant_id, title, description, subject, blooms_taxonomy, is_public, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING *
	`, bank.TenantID, bank.Title, bank.Description, bank.Subject, bank.BloomsTaxonomy, bank.IsPublic, bank.CreatedBy).StructScan(&bank)
	if err != nil {
		return nil, err
	}
	return &bank, nil
}

func (r *SQLAssessmentRepository) GetQuestionBank(ctx context.Context, id string) (*models.QuestionBank, error) {
	var bank models.QuestionBank
	err := r.db.GetContext(ctx, &bank, `SELECT * FROM question_banks WHERE id = $1`, id)
	return &bank, err
}

func (r *SQLAssessmentRepository) ListQuestionBanks(ctx context.Context, tenantID string) ([]models.QuestionBank, error) {
	var banks []models.QuestionBank
	err := r.db.SelectContext(ctx, &banks, `SELECT * FROM question_banks WHERE tenant_id = $1`, tenantID)
	return banks, err
}

func (r *SQLAssessmentRepository) UpdateQuestionBank(ctx context.Context, bank models.QuestionBank) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE question_banks
		SET title=$1, description=$2, subject=$3, blooms_taxonomy=$4, is_public=$5, updated_at=NOW()
		WHERE id=$6
	`, bank.Title, bank.Description, bank.Subject, bank.BloomsTaxonomy, bank.IsPublic, bank.ID)
	return err
}

func (r *SQLAssessmentRepository) DeleteQuestionBank(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM question_banks WHERE id=$1`, id)
	return err
}

// --- Questions ---

func (r *SQLAssessmentRepository) CreateQuestion(ctx context.Context, q models.Question) (*models.Question, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert Question
	err = tx.QueryRowxContext(ctx, `
		INSERT INTO questions (bank_id, type, stem, media_url, points_default, difficulty_level, learning_outcome_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING *
	`, q.BankID, q.Type, q.Stem, q.MediaURL, q.PointsDefault, q.DifficultyLevel, q.LearningOutcomeID).StructScan(&q)
	if err != nil {
		return nil, err
	}

	// Insert Options
	if len(q.Options) > 0 {
		for i, opt := range q.Options {
			_, err = tx.ExecContext(ctx, `
				INSERT INTO question_options (question_id, text, is_correct, sort_order, feedback)
				VALUES ($1, $2, $3, $4, $5)
			`, q.ID, opt.Text, opt.IsCorrect, i, opt.Feedback) // Use generic index as sort order if not specified
			if err != nil {
				return nil, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &q, nil // Note: Options struct not re-populated with IDs here for simplicity, typically handled by Get
}

func (r *SQLAssessmentRepository) GetQuestion(ctx context.Context, id string) (*models.Question, error) {
	var q models.Question
	// Fetch Question
	err := r.db.GetContext(ctx, &q, `SELECT * FROM questions WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}

	// Fetch Options
	err = r.db.SelectContext(ctx, &q.Options, `SELECT * FROM question_options WHERE question_id = $1 ORDER BY sort_order`, id)
	if err != nil {
		return nil, err
	}

	return &q, nil
}

func (r *SQLAssessmentRepository) ListQuestionsByBank(ctx context.Context, bankID string) ([]models.Question, error) {
	var questions []models.Question
	err := r.db.SelectContext(ctx, &questions, `SELECT * FROM questions WHERE bank_id = $1 ORDER BY created_at DESC`, bankID)
	// Not eager loading options for list view to save performance
	return questions, err
}

func (r *SQLAssessmentRepository) UpdateQuestion(ctx context.Context, q models.Question) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update Question Data
	_, err = tx.ExecContext(ctx, `
		UPDATE questions 
		SET type=$1, stem=$2, media_url=$3, points_default=$4, difficulty_level=$5, learning_outcome_id=$6, updated_at=NOW()
		WHERE id=$7
	`, q.Type, q.Stem, q.MediaURL, q.PointsDefault, q.DifficultyLevel, q.LearningOutcomeID, q.ID)
	if err != nil {
		return err
	}

	// Update Options: Simplest strategy is Delete All & Re-Insert
	// For production we might diff, but this ensures consistency easily
	_, err = tx.ExecContext(ctx, `DELETE FROM question_options WHERE question_id=$1`, q.ID)
	if err != nil {
		return err
	}

	if len(q.Options) > 0 {
		for i, opt := range q.Options {
			_, err = tx.ExecContext(ctx, `
				INSERT INTO question_options (question_id, text, is_correct, sort_order, feedback)
				VALUES ($1, $2, $3, $4, $5)
			`, q.ID, opt.Text, opt.IsCorrect, i, opt.Feedback)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *SQLAssessmentRepository) DeleteQuestion(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM questions WHERE id=$1`, id)
	return err
}

// --- Assessments ---

func (r *SQLAssessmentRepository) CreateAssessment(ctx context.Context, a models.Assessment) (*models.Assessment, error) {
	err := r.db.QueryRowxContext(ctx, `
		INSERT INTO assessments (tenant_id, course_offering_id, title, description, time_limit_minutes, available_from, available_until, shuffle_questions, grading_policy, passing_score, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING *
	`, a.TenantID, a.CourseOfferingID, a.Title, a.Description, a.TimeLimitMinutes, a.AvailableFrom, a.AvailableUntil, a.ShuffleQuestions, a.GradingPolicy, a.PassingScore, a.CreatedBy).StructScan(&a)
	return &a, err
}

func (r *SQLAssessmentRepository) GetAssessment(ctx context.Context, id string) (*models.Assessment, error) {
	var a models.Assessment
	err := r.db.GetContext(ctx, &a, `SELECT * FROM assessments WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}

	// Fetch Items/Sections would go here
	// skipping for brevity in this initial implementation
	return &a, nil
}

// --- Attempts ---

func (r *SQLAssessmentRepository) CreateAttempt(ctx context.Context, attempt models.AssessmentAttempt) (*models.AssessmentAttempt, error) {
	err := r.db.QueryRowxContext(ctx, `
		INSERT INTO assessment_attempts (assessment_id, student_id, started_at, status, score)
		VALUES ($1, $2, NOW(), $3, 0)
		RETURNING *
	`, attempt.AssessmentID, attempt.StudentID, models.AttemptStatusInProgress).StructScan(&attempt)
	return &attempt, err
}

func (r *SQLAssessmentRepository) SaveItemResponse(ctx context.Context, response models.ItemResponse) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO item_responses (attempt_id, question_id, selected_option_id, text_response, score, is_correct)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (attempt_id, question_id) 
		DO UPDATE SET selected_option_id = EXCLUDED.selected_option_id, text_response = EXCLUDED.text_response, score = EXCLUDED.score, is_correct = EXCLUDED.is_correct
	`, response.AttemptID, response.QuestionID, response.SelectedOptionID, response.TextResponse, response.Score, response.IsCorrect)
	return err
}

func (r *SQLAssessmentRepository) CompleteAttempt(ctx context.Context, attemptID string, score float64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE assessment_attempts 
		SET finished_at=NOW(), score=$1, status=$2 
		WHERE id=$3
	`, score, models.AttemptStatusSubmitted, attemptID)
	return err
}

func (r *SQLAssessmentRepository) LogProctoringEvent(ctx context.Context, log models.ProctoringLog) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO proctoring_logs (attempt_id, event_type, metadata)
		VALUES ($1, $2, $3)
	`, log.AttemptID, log.EventType, log.Metadata)
	return err
}

func (r *SQLAssessmentRepository) CountProctoringEvents(ctx context.Context, attemptID string) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM proctoring_logs WHERE attempt_id=$1`, attemptID)
	return count, err
}

func (r *SQLAssessmentRepository) GetAttempt(ctx context.Context, id string) (*models.AssessmentAttempt, error) {
	var attempt models.AssessmentAttempt
	err := r.db.GetContext(ctx, &attempt, `SELECT * FROM assessment_attempts WHERE id = $1`, id)
	return &attempt, err
}

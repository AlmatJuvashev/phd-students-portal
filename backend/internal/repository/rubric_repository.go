package repository

import (
	"context"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type RubricRepository interface {
	CreateRubric(ctx context.Context, r *models.Rubric) error
	GetRubric(ctx context.Context, id string) (*models.Rubric, error)
	ListRubrics(ctx context.Context, courseID string) ([]models.Rubric, error)
	
	SubmitGrade(ctx context.Context, g *models.RubricGrade) error
	GetGrade(ctx context.Context, submissionID string) (*models.RubricGrade, error)
}

type SQLRubricRepository struct {
	db *sqlx.DB
}

func NewSQLRubricRepository(db *sqlx.DB) *SQLRubricRepository {
	return &SQLRubricRepository{db: db}
}

// CreateRubric Deep Insert
func (r *SQLRubricRepository) CreateRubric(ctx context.Context, d *models.Rubric) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Insert Rubric
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()
	err = tx.QueryRowxContext(ctx, `
		INSERT INTO rubrics (course_offering_id, title, description, is_global, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		d.CourseOfferingID, d.Title, d.Description, d.IsGlobal, d.CreatedAt, d.UpdatedAt,
	).Scan(&d.ID)
	if err != nil {
		return err
	}

	// 2. Insert Criteria & Levels
	for i, c := range d.Criteria {
		c.RubricID = d.ID
		c.CreatedAt = time.Now()
		c.UpdatedAt = time.Now()
		c.Position = i
		
		err = tx.QueryRowxContext(ctx, `
			INSERT INTO rubric_criteria (rubric_id, title, description, weight, position, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			c.RubricID, c.Title, c.Description, c.Weight, c.Position, c.CreatedAt, c.UpdatedAt,
		).Scan(&c.ID)
		if err != nil {
			return err
		}

		for j, l := range c.Levels {
			l.CriterionID = c.ID
			l.CreatedAt = time.Now()
			l.UpdatedAt = time.Now()
			l.Position = j
			
			_, err = tx.ExecContext(ctx, `
				INSERT INTO rubric_levels (criterion_id, title, description, points, position, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7)`,
				l.CriterionID, l.Title, l.Description, l.Points, l.Position, l.CreatedAt, l.UpdatedAt,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *SQLRubricRepository) GetRubric(ctx context.Context, id string) (*models.Rubric, error) {
	var rub models.Rubric
	err := r.db.GetContext(ctx, &rub, "SELECT * FROM rubrics WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	
	// Fetch Criteria
	var criteria []models.RubricCriterion
	err = r.db.SelectContext(ctx, &criteria, "SELECT * FROM rubric_criteria WHERE rubric_id=$1 ORDER BY position ASC", id)
	if err != nil {
		return &rub, nil // Return rubric even if no criteria?
	}
	
	// Fetch Levels for each criterion? N+1 query.
	// Optimize: Fetch all levels for these criteria
	if len(criteria) > 0 {
		var criIDs []string
		for _, c := range criteria {
			criIDs = append(criIDs, c.ID)
		}
		query, args, _ := sqlx.In("SELECT * FROM rubric_levels WHERE criterion_id IN (?) ORDER BY position ASC", criIDs)
		query = r.db.Rebind(query)
		
		var allLevels []models.RubricLevel
		if err = r.db.SelectContext(ctx, &allLevels, query, args...); err == nil {
			// Map levels to criteria
			levelMap := make(map[string][]models.RubricLevel)
			for _, l := range allLevels {
				levelMap[l.CriterionID] = append(levelMap[l.CriterionID], l)
			}
			for i := range criteria {
				criteria[i].Levels = levelMap[criteria[i].ID]
			}
		}
	}
	
	rub.Criteria = criteria
	return &rub, nil
}

func (r *SQLRubricRepository) ListRubrics(ctx context.Context, courseID string) ([]models.Rubric, error) {
	var list []models.Rubric
	err := r.db.SelectContext(ctx, &list, "SELECT * FROM rubrics WHERE course_offering_id=$1 ORDER BY created_at DESC", courseID)
	return list, err
}

func (r *SQLRubricRepository) SubmitGrade(ctx context.Context, g *models.RubricGrade) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	g.CreatedAt = time.Now()
	g.UpdatedAt = time.Now()

	// Insert Grade Header
	err = tx.QueryRowxContext(ctx, `
		INSERT INTO rubric_grades (submission_id, rubric_id, grader_id, total_score, comments, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		g.SubmissionID, g.RubricID, g.GraderID, g.TotalScore, g.Comments, g.CreatedAt, g.UpdatedAt,
	).Scan(&g.ID)
	if err != nil {
		return err
	}

	// Insert items
	for _, item := range g.Items {
		item.RubricGradeID = g.ID
		item.CreatedAt = time.Now()
		_, err = tx.ExecContext(ctx, `
			INSERT INTO rubric_grade_items (rubric_grade_id, criterion_id, level_id, points_awarded, comments, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			item.RubricGradeID, item.CriterionID, item.LevelID, item.PointsAwarded, item.Comments, item.CreatedAt,
		)
		if err != nil {
			return err
		}
	}
	
	// Update Submission Score & Status
	// Assuming max possible score is untracked, we just set the score.
	_, err = tx.ExecContext(ctx, `
		UPDATE activity_submissions SET status='GRADED', score=$1, graded_at=NOW(), updated_at=NOW()
		WHERE id=$2`, g.TotalScore, g.SubmissionID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *SQLRubricRepository) GetGrade(ctx context.Context, submissionID string) (*models.RubricGrade, error) {
	var g models.RubricGrade
	err := r.db.GetContext(ctx, &g, "SELECT * FROM rubric_grades WHERE submission_id=$1", submissionID)
	if err != nil {
		return nil, err
	}
	
	// Fetch items
	var items []models.RubricGradeItem
	_ = r.db.SelectContext(ctx, &items, "SELECT * FROM rubric_grade_items WHERE rubric_grade_id=$1", g.ID)
	g.Items = items
	return &g, nil
}

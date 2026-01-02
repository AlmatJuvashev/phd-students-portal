package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type AnalyticsRepository interface {
	GetStudentsByStage(ctx context.Context) ([]models.StudentStageStats, error)
	GetAdvisorLoad(ctx context.Context) ([]models.AdvisorLoadStats, error)
	GetOverdueTasks(ctx context.Context) ([]models.OverdueTaskStats, error)

	// Agnostic Metrics
	GetTotalStudents(ctx context.Context, filter models.FilterParams) (int, error)
	GetNodeCompletionCount(ctx context.Context, nodeID string, filter models.FilterParams) (int, error)
	GetDurationForNodes(ctx context.Context, nodeIDs []string, filter models.FilterParams) ([]float64, error)
	GetBottleneck(ctx context.Context, filter models.FilterParams) (string, int, error)
	GetProfileFlagCount(ctx context.Context, key string, min float64, filter models.FilterParams) (int, error)

	// Risk Analytics
	SaveRiskSnapshot(ctx context.Context, s *models.RiskSnapshot) error
	GetStudentRiskHistory(ctx context.Context, studentID string) ([]models.RiskSnapshot, error)
	GetHighRiskStudents(ctx context.Context, threshold float64) ([]models.RiskSnapshot, error)
}

type SQLAnalyticsRepository struct {
	db *sqlx.DB
}

func NewSQLAnalyticsRepository(db *sqlx.DB) *SQLAnalyticsRepository {
	return &SQLAnalyticsRepository{db: db}
}

// --- Specific/Legacy Methods ---

func (r *SQLAnalyticsRepository) GetStudentsByStage(ctx context.Context) ([]models.StudentStageStats, error) {
	query := `
		SELECT 
			COALESCE(js.state, 'Not Started') as stage, 
			COUNT(*) as count
		FROM users u
		LEFT JOIN journey_states js ON u.id = js.user_id
		WHERE u.role = 'student'
		GROUP BY js.state`
	
	var stats []models.StudentStageStats
	err := r.db.SelectContext(ctx, &stats, query)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (r *SQLAnalyticsRepository) GetAdvisorLoad(ctx context.Context) ([]models.AdvisorLoadStats, error) {
	query := `
		SELECT 
			COALESCE(u.first_name || ' ' || u.last_name, 'Unknown') as advisor_name,
			COUNT(sa.student_id) as student_count
		FROM users u
		JOIN student_advisors sa ON u.id = sa.advisor_id
		WHERE u.role = 'advisor'
		GROUP BY u.id, u.first_name, u.last_name
		ORDER BY student_count DESC`
	
	var stats []models.AdvisorLoadStats
	err := r.db.SelectContext(ctx, &stats, query)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (r *SQLAnalyticsRepository) GetOverdueTasks(ctx context.Context) ([]models.OverdueTaskStats, error) {
	query := `
		SELECT 
			nd.node_id,
			COUNT(*) as count
		FROM node_deadlines nd
		LEFT JOIN node_instances ni ON nd.user_id = ni.user_id AND nd.node_id = ni.node_id
		WHERE nd.due_at < NOW() 
		AND (ni.state IS NULL OR ni.state != 'done')
		GROUP BY nd.node_id
		ORDER BY count DESC`
	
	var stats []models.OverdueTaskStats
	err := r.db.SelectContext(ctx, &stats, query)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

// --- Agnostic Implementations ---

func (r *SQLAnalyticsRepository) GetTotalStudents(ctx context.Context, filter models.FilterParams) (int, error) {
	base := `SELECT COUNT(DISTINCT u.id) FROM users u JOIN user_tenant_memberships utm ON utm.user_id = u.id`
	where, args := r.buildFilter(filter)
	
	query := base + " WHERE " + strings.Join(where, " AND ")
	query = r.db.Rebind(query)
	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	return count, err
}

func (r *SQLAnalyticsRepository) GetNodeCompletionCount(ctx context.Context, nodeID string, filter models.FilterParams) (int, error) {
	// Count users matching filter who have instance(nodeID) in 'done'
	base := `SELECT COUNT(DISTINCT u.id) 
			 FROM users u 
			 JOIN user_tenant_memberships utm ON utm.user_id = u.id
			 JOIN node_instances ni ON ni.user_id = u.id`
	
	where, args := r.buildFilter(filter)
	where = append(where, "ni.node_id = ?", "ni.state = 'done'")
	args = append(args, nodeID)
	
	query := base + " WHERE " + strings.Join(where, " AND ")
	query = r.db.Rebind(query)
	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	return count, err
}

func (r *SQLAnalyticsRepository) GetStageDurationStats(ctx context.Context, stageID string, filter models.FilterParams) ([]float64, error) {
	// Get start/end times for node_instances in this stage (approximation via list of nodes?)
	// Or query journey_states history if available?
	// Plan: Use node_instances for generic "W2" approach (which is just a set of nodes or a 'stage' tag?)
	// Actually, the previous implementation assumed stages 'W2' etc are mapped or defined.
	// But `journey_states` table has `node_id`. 'W2' is usually a 'World' (Frontend concept).
	// If backend tracks it: We can check `journey_states` if it is stored there?
	// Existing data.ts shows "id": "W2".
	// The original sql implementation checked for specific nodes within W2.
	// For AGNOSTIC approach, if "StageID" is passed, we might need to know which nodes comprise it, OR assume 'StageID' is stored in journey_states?
	// The `journey_states` table keys are `node_id` (uuid per user+node).
    // Let's assume we pass a list of node IDs that represent the stage, OR simply query `node_instances` time range if `stageID` corresponds to a set of nodes.
    // BUT, the prompt said `GetStageDurationStats(stageID)`.
    // Let's implement looking up min/max updated_at for *all* node instances of a user if we can't map W2 directly.
    // Wait, the "Median Days in W2" usually means time from First Node Start to Last Node End in that world.
    
    // Simplification for Agnostic: We calculate time for a *single* node or aggregate?
    // Let's stick to the previous logic: `GetW2Durations` took `w2Nodes []string`.
    // So we should pass `nodeIDs []string` representing the stage.
    // Actually, let's keep it simple: `StageID` might not be enough without a map.
    // But the Repo method signature I proposed was `GetStageDurationStats(ctx, stageID string, ...)`
    // This implies the repo knows what 'stageID' means? No, that violates agnostic.
    // Let's change signature to `GetDurationForNodes(ctx, nodeIDs []string, filter)`.
    // But to respect the Interface I just defined in thought...
    // Let's assume the Service resolves Stage -> Nodes and passes Nodes.
    // I will change the impl here to `GetStageDurationStats` but it actually expects `nodeIDs` joined or I will change the interface to accept `[]string`.
    // Let's change the interface in this file to `GetDurationForNodes`.
    
    return nil, fmt.Errorf("use GetDurationForNodes instead")
}

// Corrected Method
func (r *SQLAnalyticsRepository) GetDurationForNodes(ctx context.Context, nodeIDs []string, filter models.FilterParams) ([]float64, error) {
	if len(nodeIDs) == 0 { return nil, nil }
	
	base := `SELECT u.id, MIN(ni.updated_at), MAX(ni.updated_at)
			 FROM users u
			 JOIN user_tenant_memberships utm ON utm.user_id = u.id
			 JOIN node_instances ni ON ni.user_id = u.id`
	
	where, args := r.buildFilter(filter)
	
	// Add Node Constraints
	query, a2, _ := sqlx.In("ni.node_id IN (?)", nodeIDs)
	where = append(where, query)
	args = append(args, a2...)
	
	fullQuery := base + " WHERE " + strings.Join(where, " AND ") + " GROUP BY u.id"
	fullQuery = r.db.Rebind(fullQuery)
	
	rows, err := r.db.QueryxContext(ctx, fullQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var durations []float64
	for rows.Next() {
		var uid string
		var start, end time.Time
		if err := rows.Scan(&uid, &start, &end); err == nil && !end.Before(start) {
			days := end.Sub(start).Hours() / 24.0
			durations = append(durations, days)
		}
	}
	return durations, nil
}


func (r *SQLAnalyticsRepository) GetBottleneck(ctx context.Context, filter models.FilterParams) (string, int, error) {
	// Find node with max pending/in_progress count
	base := `SELECT ni.node_id, COUNT(*) as cnt
			 FROM users u
			 JOIN user_tenant_memberships utm ON utm.user_id = u.id
			 JOIN node_instances ni ON ni.user_id = u.id`
			 
	where, args := r.buildFilter(filter)
	where = append(where, "ni.state IN ('waiting', 'needs_fixes', 'in_progress', 'pending')") // Added common states
	
	query := base + " WHERE " + strings.Join(where, " AND ") + " GROUP BY ni.node_id ORDER BY cnt DESC LIMIT 1"
	query = r.db.Rebind(query)
	
	var nodeID string
	var count int
	err := r.db.QueryRowxContext(ctx, query, args...).Scan(&nodeID, &count)
	if err != nil {
		// No rows is fine
		return "", 0, nil
	}
	return nodeID, count, nil
}

func (r *SQLAnalyticsRepository) GetProfileFlagCount(ctx context.Context, key string, minVal float64, filter models.FilterParams) (int, error) {
	// JSON path check
	base := `SELECT COUNT(DISTINCT u.id)
			 FROM users u
			 JOIN user_tenant_memberships utm ON utm.user_id = u.id
			 LEFT JOIN profile_submissions ps ON ps.user_id = u.id`
			 
	where, args := r.buildFilter(filter)
	// Add JSON check. Assumes Postgres.
	// We want latest submission? The join above might duplicate.
	// Better: `COALESCE(u.program, (SELECT form_data->>'program' ...))` logic from AdminRepo is for attributes.
	// Here we specifically check profile form data value.
	// "CAST(ps.form_data->>? AS NUMERIC) > ?"
	
	// Ensure we only check the LATEST submission if multiple exist (though profile_submissions usually log history)
	// For simplicity in this step, let's assume one submission or latest logic isn't strictly enforced by simple JOIN.
	// Correct way: JOIN LATERAL or Subquery.
	// Subquery approach:
	where = append(where, fmt.Sprintf("(SELECT CAST(form_data->>'%s' AS NUMERIC) FROM profile_submissions WHERE user_id=u.id ORDER BY submitted_at DESC LIMIT 1) > ?", key))
	args = append(args, minVal)

	query := base + " WHERE " + strings.Join(where, " AND ")
	query = r.db.Rebind(query)
	
	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	return count, err
}

// buildFilter reuse
func (r *SQLAnalyticsRepository) buildFilter(f models.FilterParams) ([]string, []interface{}) {
	var where []string
	var args []interface{}
	
	where = append(where, "u.role = 'student'", "u.is_active = true")
	
	if f.TenantID != "" {
		where = append(where, "utm.tenant_id = ?")
		args = append(args, f.TenantID)
	}
	if f.AdvisorID != "" {
		// Helper subquery exists check
		where = append(where, "EXISTS (SELECT 1 FROM student_advisors sa WHERE sa.student_id = u.id AND sa.advisor_id = ?)")
		args = append(args, f.AdvisorID)
	}
	if f.Program != "" {
		where = append(where, "COALESCE(u.program, '') = ?")
		args = append(args, f.Program)
	}
	if f.Department != "" {
		where = append(where, "COALESCE(u.department, '') = ?")
		args = append(args, f.Department)
	}
	if f.Cohort != "" {
		where = append(where, "COALESCE(u.cohort, '') = ?")
		args = append(args, f.Cohort)
	}
	
	return where, args
}

// --- Risk Analytics Implementation ---

func (r *SQLAnalyticsRepository) SaveRiskSnapshot(ctx context.Context, s *models.RiskSnapshot) error {
	// Insert snapshot
	s.CreatedAt = time.Now()
	// Marshal RiskFactors for DB if needed, but model handles it?
	// RiskSnapshot struct has RiskFactors (JSONB) and RawFactors?
	// Let's assume passed struct has RawFactors populated if using that, or we assume RiskFactors is []byte.
	// The model definition: RiskFactors types.JSONText `db:"risk_factors"`
	
	query := `
		INSERT INTO student_risk_snapshots (student_id, risk_score, risk_factors, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id`
	
	err := r.db.QueryRowxContext(ctx, query, s.StudentID, s.RiskScore, s.RiskFactors, s.CreatedAt).Scan(&s.ID)
	return err
}

func (r *SQLAnalyticsRepository) GetStudentRiskHistory(ctx context.Context, studentID string) ([]models.RiskSnapshot, error) {
	var snapshots []models.RiskSnapshot
	query := `SELECT * FROM student_risk_snapshots WHERE student_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &snapshots, query, studentID)
	return snapshots, err
}

func (r *SQLAnalyticsRepository) GetHighRiskStudents(ctx context.Context, threshold float64) ([]models.RiskSnapshot, error) {
	// Get LATEST snapshot for each student that is above threshold
	// Use DISTINCT ON (student_id) logic
	query := `
		SELECT DISTINCT ON (student_id) *
		FROM student_risk_snapshots
		WHERE risk_score >= $1
		ORDER BY student_id, created_at DESC`
		
	var snapshots []models.RiskSnapshot
	err := r.db.SelectContext(ctx, &snapshots, query, threshold)
	
	// Also populate StudentName for convenience?
	// We might need a JOIN query to do that efficiently.
	if err == nil && len(snapshots) > 0 {
		// Populate names... or update query
		// Let's update query to join users
		queryWithJoin := `
			SELECT DISTINCT ON (srs.student_id) srs.*, u.first_name || ' ' || u.last_name as student_name
			FROM student_risk_snapshots srs
			JOIN users u ON srs.student_id = u.id
			WHERE srs.risk_score >= $1
			ORDER BY srs.student_id, srs.created_at DESC`
		err = r.db.SelectContext(ctx, &snapshots, queryWithJoin, threshold)
	}
	
	return snapshots, err
}

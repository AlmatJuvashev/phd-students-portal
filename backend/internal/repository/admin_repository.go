package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AdminRepository interface {
	ListStudentProgress(ctx context.Context, tenantID, playbookVersionID string) ([]models.StudentProgressSummary, error)
	ListStudentsForMonitor(ctx context.Context, filter models.FilterParams) ([]models.StudentMonitorRow, error)
	GetAnalytics(ctx context.Context, filter models.FilterParams, playbookVersionID string) (*models.AdminAnalytics, error)
	GetStudentDetails(ctx context.Context, studentID, tenantID string) (*models.StudentDetails, error)
	
	// Helper batch loaders
	GetAdvisorsForStudents(ctx context.Context, studentIDs []string) (map[string][]models.AdvisorSummary, error)
	GetDoneCountsForStudents(ctx context.Context, studentIDs []string) (map[string]int, error)
	GetLastUpdatesForStudents(ctx context.Context, studentIDs []string) (map[string]time.Time, error)
	GetRPRequiredForStudents(ctx context.Context, studentIDs []string) (map[string]bool, error)
	
	// Single student graph
	GetStudentNodeInstances(ctx context.Context, studentID string) ([]models.NodeInstance, error) 
	
	// Analytics Aggregates
	GetAntiplagCount(ctx context.Context, studentIDs []string, playbookVersionID string) (int, error)
	GetW2Durations(ctx context.Context, studentIDs []string, playbookVersionID string, w2Nodes []string) ([]float64, error)
	GetBottleneck(ctx context.Context, studentIDs []string, playbookVersionID string, since time.Time) (string, int, error)
	
	// Refactoring additions
	CheckAdvisorAccess(ctx context.Context, studentID, advisorID string) (bool, error)
	GetStudentJourneyNodes(ctx context.Context, studentID string) ([]models.StudentJourneyNode, error)
	GetNodeFiles(ctx context.Context, studentID, nodeID string) ([]models.NodeFile, error)
	
	// Attachment Review
	GetAttachmentMeta(ctx context.Context, attachmentID string) (*models.AttachmentMeta, error)
	GetLatestAttachmentStatus(ctx context.Context, instanceID string) (string, error)
	GetAttachmentCounts(ctx context.Context, instanceID string) (submitted, approved, rejected int, err error)
	UpdateAttachmentStatus(ctx context.Context, attachmentID, status, note, actorID string) error
	UploadReviewedDocument(ctx context.Context, attachmentID, versionID, actorID string) error
	LogNodeEvent(ctx context.Context, instanceID, eventType, actorID string, payload map[string]any) error
	
	// Node State Management
	UpdateNodeInstanceState(ctx context.Context, instanceID, state string) error
	UpdateAllNodeInstances(ctx context.Context, studentID, nodeID, instanceID, state string) error
	UpsertJourneyState(ctx context.Context, tenantID, studentID, nodeID, state string) error
	
	// Reminders
	CreateReminders(ctx context.Context, studentIDs []string, title, message string, dueAt *string, createdBy string) error
	
	// Notifications
	CreateNotification(ctx context.Context, recipientID, title, message, link, nType, tenantID string) error
	
	// Reviewed Docs
	CreateReviewedDocumentVersion(ctx context.Context, docID, storagePath, objKey, bucket, mimeType string, sizeBytes int64, actorID, etag, tenantID string) (string, error)

	// Admin Notifications
	ListAdminNotifications(ctx context.Context, unreadOnly bool) ([]models.AdminNotification, error)
	GetAdminUnreadCount(ctx context.Context) (int, error)
	MarkAdminNotificationRead(ctx context.Context, id string) error
	MarkAllAdminNotificationsRead(ctx context.Context) error
}

type SQLAdminRepository struct {
	db *sqlx.DB
}

func NewSQLAdminRepository(db *sqlx.DB) *SQLAdminRepository {
	return &SQLAdminRepository{db: db}
}

func (r *SQLAdminRepository) ListStudentProgress(ctx context.Context, tenantID, playbookVersionID string) ([]models.StudentProgressSummary, error) {
	// Optimized query combining counts and latest node
	query := `
		SELECT 
			u.id, 
			(u.first_name || ' ' || u.last_name) as name, 
			u.email, 
			u.role,
			COUNT(ni.id) FILTER (WHERE ni.state='done' AND ni.playbook_version_id=$2) as completed_nodes,
			(SELECT node_id FROM node_instances WHERE user_id=u.id AND playbook_version_id=$2 ORDER BY updated_at DESC LIMIT 1) as current_node_id,
			(SELECT to_char(MAX(updated_at), 'YYYY-MM-DD"T"HH24:MI:SSZ') FROM node_instances WHERE user_id=u.id AND playbook_version_id=$2) as last_submission_at
		FROM users u
		JOIN user_tenant_memberships utm ON utm.user_id = u.id
		LEFT JOIN node_instances ni ON ni.user_id = u.id AND ni.playbook_version_id=$2
		WHERE utm.tenant_id = $1 AND u.role='student' AND u.is_active=true
		GROUP BY u.id
		ORDER BY u.last_name
	`
	var rows []models.StudentProgressSummary
	err := r.db.SelectContext(ctx, &rows, query, tenantID, playbookVersionID)
	return rows, err
}

func (r *SQLAdminRepository) ListStudentsForMonitor(ctx context.Context, filter models.FilterParams) ([]models.StudentMonitorRow, error) {
	// Base query construction similar to handler but cleaner
	base := `SELECT u.id, (u.first_name||' '||u.last_name) AS name, COALESCE(u.email,'') AS email,
			COALESCE(u.phone, (SELECT form_data->>'phone' FROM profile_submissions ps WHERE ps.user_id=u.id ORDER BY ps.submitted_at DESC LIMIT 1), '') AS phone,
			COALESCE(u.program, (SELECT form_data->>'program' FROM profile_submissions ps WHERE ps.user_id=u.id ORDER BY ps.submitted_at DESC LIMIT 1), '') AS program,
			COALESCE(u.department, (SELECT form_data->>'department' FROM profile_submissions ps WHERE ps.user_id=u.id ORDER BY ps.submitted_at DESC LIMIT 1), '') AS department,
			COALESCE(u.cohort, (SELECT form_data->>'cohort' FROM profile_submissions ps WHERE ps.user_id=u.id ORDER BY ps.submitted_at DESC LIMIT 1), '') AS cohort,
			(SELECT node_id FROM node_instances WHERE user_id=u.id ORDER BY updated_at DESC LIMIT 1) as current_node_id
			FROM users u
			JOIN user_tenant_memberships utm ON utm.user_id = u.id`

	whereConditions := []string{"u.is_active=true", "u.role='student'", "utm.tenant_id=$1"}
	args := []interface{}{filter.TenantID}

	if filter.AdvisorID != "" {
		// Only join if filtering by advisor
		base += " JOIN student_advisors sa ON sa.student_id=u.id"
		whereConditions = append(whereConditions, fmt.Sprintf("sa.advisor_id=$%d", len(args)+1))
		args = append(args, filter.AdvisorID)
	}

	// Helper to add condition
	addCond := func(clause string, val interface{}) {
		whereConditions = append(whereConditions, fmt.Sprintf(clause, len(args)+1))
		args = append(args, val)
	}

	if filter.Program != "" {
		addCond("COALESCE(u.program, (SELECT form_data->>'program' FROM profile_submissions WHERE user_id=u.id ORDER BY submitted_at DESC LIMIT 1))=$%d", filter.Program)
	}
	if filter.Department != "" {
		addCond("COALESCE(u.department, (SELECT form_data->>'department' FROM profile_submissions WHERE user_id=u.id ORDER BY submitted_at DESC LIMIT 1))=$%d", filter.Department)
	}
	if filter.Cohort != "" {
		addCond("COALESCE(u.cohort, (SELECT form_data->>'cohort' FROM profile_submissions WHERE user_id=u.id ORDER BY submitted_at DESC LIMIT 1))=$%d", filter.Cohort)
	}
	if filter.Query != "" {
		idx := len(args) + 1
		qClause := fmt.Sprintf("((u.first_name ILIKE '%%' || $%d || '%%') OR (u.last_name ILIKE '%%' || $%d || '%%') OR (u.email ILIKE '%%' || $%d || '%%'))", idx, idx, idx)
		// Simpler query for optimization, original handler had robust checks. Keeping it consistent.
		whereConditions = append(whereConditions, qClause)
		args = append(args, filter.Query)
	}

	fullQuery := base + " WHERE " + strings.Join(whereConditions, " AND ") + " ORDER BY u.last_name, u.first_name"
	if filter.Limit > 0 {
		fullQuery += fmt.Sprintf(" LIMIT %d", filter.Limit)
	}

	var rows []models.StudentMonitorRow
	err := r.db.SelectContext(ctx, &rows, fullQuery, args...)
	return rows, err
}

func (r *SQLAdminRepository) GetAdvisorsForStudents(ctx context.Context, studentIDs []string) (map[string][]models.AdvisorSummary, error) {
	if len(studentIDs) == 0 {
		return nil, nil
	}
	query, args, err := sqlx.In(`
		SELECT sa.student_id, u.id, (u.first_name||' '||u.last_name) AS name, COALESCE(u.email,'') as email
		FROM student_advisors sa 
		JOIN users u ON u.id=sa.advisor_id 
		WHERE sa.student_id IN (?)`, studentIDs)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[string][]models.AdvisorSummary)
	for rows.Next() {
		var sid string
		var adv models.AdvisorSummary
		if err := rows.Scan(&sid, &adv.ID, &adv.Name, &adv.Email); err == nil {
			out[sid] = append(out[sid], adv)
		}
	}
	return out, nil
}

func (r *SQLAdminRepository) GetDoneCountsForStudents(ctx context.Context, studentIDs []string) (map[string]int, error) {
	if len(studentIDs) == 0 {
		return nil, nil
	}
	// Across all versions, as per handler logic
	query, args, err := sqlx.In(`SELECT user_id, COUNT(*) FROM node_instances WHERE state='done' AND user_id IN (?) GROUP BY user_id`, studentIDs)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	out := make(map[string]int)
	for rows.Next() {
		var uid string
		var cnt int
		if err := rows.Scan(&uid, &cnt); err == nil {
			out[uid] = cnt
		}
	}
	return out, nil
}

func (r *SQLAdminRepository) GetLastUpdatesForStudents(ctx context.Context, studentIDs []string) (map[string]time.Time, error) {
	// Batch optimize the GREATEST query
	// Complex: Postgres doesn't easily support lateral join on list of IDs passed as IN clause unless we use UNNEST.
	// But we can do 4 separate aggregate queries grouped by user_id, then merge in Go. This is O(4) queries instead of O(N).
	if len(studentIDs) == 0 {
		return nil, nil
	}

	result := make(map[string]time.Time)
	
	// Helper to merge max time
	merge := func(uid string, t time.Time) {
		if cur, ok := result[uid]; !ok || t.After(cur) {
			result[uid] = t
		}
	}

	// 1. Node Instances
	q1, a1, _ := sqlx.In("SELECT user_id, MAX(updated_at) FROM node_instances WHERE user_id IN (?) GROUP BY user_id", studentIDs)
	if rows, err := r.db.QueryContext(ctx, r.db.Rebind(q1), a1...); err == nil {
		for rows.Next() {
			var uid string
			var t time.Time
			rows.Scan(&uid, &t)
			merge(uid, t)
		}
		rows.Close()
	}

	// 2. Revisions
	q2, a2, _ := sqlx.In("SELECT ni.user_id, MAX(r.created_at) FROM node_instance_form_revisions r JOIN node_instances ni ON ni.id=r.node_instance_id WHERE ni.user_id IN (?) GROUP BY ni.user_id", studentIDs)
	if rows, err := r.db.QueryContext(ctx, r.db.Rebind(q2), a2...); err == nil {
		for rows.Next() {
			var uid string
			var t time.Time
			rows.Scan(&uid, &t)
			merge(uid, t)
		}
		rows.Close()
	}

	// 3. Attachments
	q3, a3, _ := sqlx.In("SELECT ni.user_id, MAX(a.attached_at) FROM node_instance_slot_attachments a JOIN node_instance_slots s ON s.id=a.slot_id JOIN node_instances ni ON ni.id=s.node_instance_id WHERE ni.user_id IN (?) GROUP BY ni.user_id", studentIDs)
	if rows, err := r.db.QueryContext(ctx, r.db.Rebind(q3), a3...); err == nil {
		for rows.Next() {
			var uid string
			var t time.Time
			rows.Scan(&uid, &t)
			merge(uid, t)
		}
		rows.Close()
	}

	// 4. Events
	q4, a4, _ := sqlx.In("SELECT ni.user_id, MAX(e.created_at) FROM node_events e JOIN node_instances ni ON ni.id=e.node_instance_id WHERE ni.user_id IN (?) GROUP BY ni.user_id", studentIDs)
	if rows, err := r.db.QueryContext(ctx, r.db.Rebind(q4), a4...); err == nil {
		for rows.Next() {
			var uid string
			var t time.Time
			rows.Scan(&uid, &t)
			merge(uid, t)
		}
		rows.Close()
	}

	return result, nil
}

func (r *SQLAdminRepository) GetRPRequiredForStudents(ctx context.Context, studentIDs []string) (map[string]bool, error) {
	if len(studentIDs) == 0 {
		return nil, nil
	}
	query, args, _ := sqlx.In("SELECT user_id, form_data FROM profile_submissions WHERE user_id IN (?)", studentIDs)
	query = r.db.Rebind(query)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	out := make(map[string]bool)
	for rows.Next() {
		var uid string
		var raw []byte
		if err := rows.Scan(&uid, &raw); err == nil {
			// Check JSON: years_since_graduation > 3
			// Quick string check or unmarshal? Unmarshal safer.
			// Optimization: query it inside SQL directly using ->> if predictable.
			// Handler logic: y > 3. 
			// Let's rely on Service to parse, or do it here. 
			// Doing it here matches "Batch Loader" pattern.
			var m map[string]interface{}
			if json.Unmarshal(raw, &m) == nil {
				if y, ok := m["years_since_graduation"].(float64); ok && y > 3 {
					out[uid] = true
				}
			}
		}
	}
	return out, nil
}

// TODO: Implement GetAnalytics and GetStudentDetails moving logic from handler
func (r *SQLAdminRepository) GetStudentDetails(ctx context.Context, studentID, tenantID string) (*models.StudentDetails, error) {
	// ... implementation mirroring handler query ...
	var user models.StudentDetails
	query := `SELECT u.id, COALESCE(u.email,'') AS email, 
			 COALESCE(ps.form_data->>'phone','') AS phone,
			 u.first_name, u.last_name,
			 COALESCE(ps.form_data->>'program','') AS program,
			 COALESCE(ps.form_data->>'department','') AS department,
			 COALESCE(ps.form_data->>'cohort','') AS cohort
	  FROM users u
	  LEFT JOIN profile_submissions ps ON ps.user_id = u.id
	  JOIN user_tenant_memberships utm ON utm.user_id = u.id
	  WHERE u.id=$1 AND u.role='student' AND utm.tenant_id=$2`
	
	err := r.db.GetContext(ctx, &user, query, studentID, tenantID)
	if err != nil {
		return nil, err
	}
	user.Name = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	return &user, nil
}

func (r *SQLAdminRepository) GetStudentNodeInstances(ctx context.Context, studentID string) ([]models.NodeInstance, error) {
	var instances []models.NodeInstance
	// Return latest per node_id
	query := `
		SELECT DISTINCT ON (node_id) id, tenant_id, user_id, playbook_version_id, node_id, state, opened_at, submitted_at, updated_at
		FROM node_instances 
		WHERE user_id=$1 
		ORDER BY node_id, updated_at DESC
	`
	fmt.Printf("[GetStudentNodeInstances] Query for studentID=%s\n", studentID)
	err := r.db.SelectContext(ctx, &instances, query, studentID)
	fmt.Printf("[GetStudentNodeInstances] Returned %d instances, err=%v\n", len(instances), err)
	return instances, err
}


func (r *SQLAdminRepository) GetAnalytics(ctx context.Context, filter models.FilterParams, playbookVersionID string) (*models.AdminAnalytics, error) {
	// Not implemented directly, used by service composition
	return nil, nil // Service orchestrates this
}

// Aggregation helpers for Analytics Service
func (r *SQLAdminRepository) GetAntiplagCount(ctx context.Context, studentIDs []string, playbookVersionID string) (int, error) {
	if len(studentIDs) == 0 {
		return 0, nil
	}
	query, args, err := sqlx.In("SELECT COUNT(*) FROM node_instances WHERE playbook_version_id=? AND node_id='S1_antiplag' AND state='done' AND user_id IN (?)", playbookVersionID, studentIDs)
	if err != nil {
		return 0, err
	}
	query = r.db.Rebind(query)
	var count int
	err = r.db.GetContext(ctx, &count, query, args...)
	return count, err
}

func (r *SQLAdminRepository) GetW2Durations(ctx context.Context, studentIDs []string, playbookVersionID string, w2Nodes []string) ([]float64, error) {
	if len(studentIDs) == 0 || len(w2Nodes) == 0 {
		return nil, nil
	}
	// We need min/max per student. 
	// Doing it in one query is tricky with IN (?) for generic SQL. 
	// But we can GROUP BY user_id.
	// SELECT user_id, MIN(updated_at), MAX(updated_at) ... GROUP BY user_id
	
	q, args, err := sqlx.In("SELECT user_id, MIN(updated_at) as start, MAX(updated_at) as end FROM node_instances WHERE playbook_version_id=? AND node_id IN (?) AND user_id IN (?) GROUP BY user_id", playbookVersionID, w2Nodes, studentIDs)
	if err != nil {
		return nil, err
	}
	q = r.db.Rebind(q)
	
	rows, err := r.db.QueryxContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var durations []float64
	for rows.Next() {
		var uid string
		var start, end time.Time
		if err := rows.Scan(&uid, &start, &end); err == nil {
			if !end.Before(start) {
				days := end.Sub(start).Hours() / 24.0
				durations = append(durations, days)
			}
		}
	}
	return durations, nil
}

func (r *SQLAdminRepository) GetBottleneck(ctx context.Context, studentIDs []string, playbookVersionID string, since time.Time) (string, int, error) {
	if len(studentIDs) == 0 {
		return "", 0, nil
	}
	// "waiting" or "needs_fixes"
	q, args, err := sqlx.In(`SELECT node_id, COUNT(*) as cnt 
		FROM node_instances 
		WHERE playbook_version_id=? AND user_id IN (?) AND state IN ('waiting','needs_fixes') AND updated_at >= ? 
		GROUP BY node_id ORDER BY cnt DESC LIMIT 1`, playbookVersionID, studentIDs, since)
	if err != nil {
		return "", 0, err
	}
	q = r.db.Rebind(q)
	
	var nodeID string
	var count int
	err = r.db.QueryRowContext(ctx, q, args...).Scan(&nodeID, &count)
	if err == sql.ErrNoRows {
		return "", 0, nil
	}
	return nodeID, count, err
}

func (r *SQLAdminRepository) CheckAdvisorAccess(ctx context.Context, studentID, advisorID string) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM student_advisors WHERE student_id=$1 AND advisor_id=$2)`, studentID, advisorID)
	return exists, err
}

func (r *SQLAdminRepository) GetStudentJourneyNodes(ctx context.Context, studentID string) ([]models.StudentJourneyNode, error) {
	// 1. Get Distinct Instances
	query := `
		SELECT DISTINCT ON (node_id) id, node_id, state, to_char(updated_at, 'YYYY-MM-DD"T"HH24:MI:SSZ') as updated_at
		FROM node_instances 
		WHERE user_id=$1 
		ORDER BY node_id, updated_at DESC
	`
	var nodes []models.StudentJourneyNode
	err := r.db.SelectContext(ctx, &nodes, query, studentID)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	// 2. Batch load attachments
	// Map instanceID -> []Files
	instanceIDs := make([]string, len(nodes))
	for i, n := range nodes {
		instanceIDs[i] = n.ID
	}
	
	q2, args, err := sqlx.In(`
		SELECT s.node_instance_id, a.filename, a.size_bytes, to_char(a.attached_at, 'YYYY-MM-DD"T"HH24:MI:SSZ') as attached_at, dv.id as version_id
		FROM node_instance_slot_attachments a 
		JOIN node_instance_slots s ON s.id=a.slot_id 
		JOIN document_versions dv ON dv.id=a.document_version_id
		WHERE s.node_instance_id IN (?) AND a.is_active=true
		ORDER BY a.attached_at DESC
	`, instanceIDs)
	if err != nil {
		return nil, err // Should not happen with valid input
	}
	q2 = r.db.Rebind(q2)
	
	rows, err := r.db.QueryxContext(ctx, q2, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	filesMap := make(map[string][]models.NodeSimplifiedFile)
	for rows.Next() {
		var iid string
		var f models.NodeSimplifiedFile
		if err := rows.Scan(&iid, &f.Filename, &f.SizeBytes, &f.AttachedAt, &f.VersionID); err == nil {
			f.DownloadURL = fmt.Sprintf("/api/documents/versions/%s/download", f.VersionID)
			filesMap[iid] = append(filesMap[iid], f)
		}
	}

	// 3. Assign back
	for i := range nodes {
		nodes[i].Files = filesMap[nodes[i].ID]
		if nodes[i].Files == nil {
			nodes[i].Files = []models.NodeSimplifiedFile{}
		}
		nodes[i].Attachments = len(nodes[i].Files)
	}

	return nodes, nil
}

func (r *SQLAdminRepository) GetNodeFiles(ctx context.Context, studentID, nodeID string) ([]models.NodeFile, error) {
	// Find instance
	var instanceID string
	err := r.db.GetContext(ctx, &instanceID, `SELECT id FROM node_instances WHERE user_id=$1 AND node_id=$2 ORDER BY updated_at DESC LIMIT 1`, studentID, nodeID)
	if err != nil {
		return nil, err
	}

	query := `SELECT s.slot_key, a.id as attachment_id, a.filename, a.size_bytes, a.status, a.review_note, a.is_active,
		to_char(a.attached_at, 'YYYY-MM-DD"T"HH24:MI:SSZ') as attached_at, 
		to_char(a.approved_at, 'YYYY-MM-DD"T"HH24:MI:SSZ') as approved_at, 
		a.approved_by, dv.id AS version_id, dv.mime_type,
		COALESCE(u.first_name||' '||u.last_name,'') AS uploaded_by,
		a.reviewed_document_version_id as reviewed_doc_id,
		to_char(a.reviewed_at, 'YYYY-MM-DD"T"HH24:MI:SSZ') as reviewed_at,
		rdv.mime_type AS reviewed_mime_type,
		COALESCE(ru.first_name||' '||ru.last_name,'') AS reviewed_by_name
		FROM node_instance_slots s
		JOIN node_instance_slot_attachments a ON a.slot_id=s.id
		JOIN document_versions dv ON dv.id=a.document_version_id
		LEFT JOIN users u ON u.id=a.attached_by
		LEFT JOIN document_versions rdv ON rdv.id=a.reviewed_document_version_id
		LEFT JOIN users ru ON ru.id=a.reviewed_by
		WHERE s.node_instance_id=$1
		ORDER BY a.attached_at ASC`
		
	// Need to handle nulls. SQLX struct scan is good but types must match perfectly or use standard Scan.
	// models.NodeFile uses pointers for nullables.
	// But `reviewed_mime_type` in query above calculates `rdv.id`. Wait, `reviewed_mime_type` should be `rdv.mime_type`.
	
	rows, err := r.db.QueryxContext(ctx, query, instanceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.NodeFile
	for rows.Next() {
		var nf models.NodeFile
		// We can scan into struct if db tags match. 
		// "reviewed_doc_id" matches struct tag.
		// "reviewed_mime_type" in struct is *string.
		// Let's rely on struct scan if tags are correct.
		// Query column names: `attachment_id`, `uploaded_by`, etc.
		// I'll adjust the query to match struct tags 100%.
		err := rows.StructScan(&nf)
		if err != nil {
			// Handle mismatch or sql nulls if struct fields aren't pointers?
			// NodeFile has *string for nullable fields.
			continue
		}
		
		// Post process
		if nf.VersionID != "" {
			nf.DownloadURL = fmt.Sprintf("/api/documents/versions/%s/download", nf.VersionID)
		}

		if nf.ReviewedDocID != nil && *nf.ReviewedDocID != "" {
			nf.ReviewedDocument = &models.ReviewedDocInfo{
				VersionID: *nf.ReviewedDocID,
				DownloadURL: fmt.Sprintf("/api/documents/versions/%s/download", *nf.ReviewedDocID),
			}
			if nf.ReviewedMimeType != nil { nf.ReviewedDocument.MimeType = *nf.ReviewedMimeType }
			if nf.ReviewedByName != nil { nf.ReviewedDocument.ReviewedBy = *nf.ReviewedByName }
			if nf.ReviewedAt != nil { nf.ReviewedDocument.ReviewedAt = *nf.ReviewedAt }
		}
		
		out = append(out, nf)
	}
	return out, nil
}

func (r *SQLAdminRepository) GetAttachmentMeta(ctx context.Context, attachmentID string) (*models.AttachmentMeta, error) {
	var meta models.AttachmentMeta
	query := `SELECT ni.id AS instance_id, s.id AS slot_id, s.slot_key, ni.node_id, ni.user_id as student_id, ni.state, ni.tenant_id, a.filename, a.status, dv.document_id
		FROM node_instance_slot_attachments a
		JOIN node_instance_slots s ON s.id=a.slot_id
		JOIN node_instances ni ON ni.id=s.node_instance_id
		JOIN document_versions dv ON dv.id=a.document_version_id
		WHERE a.id=$1 AND a.is_active=true`
	err := r.db.GetContext(ctx, &meta, query, attachmentID)
	return &meta, err
}

func (r *SQLAdminRepository) UpdateAttachmentStatus(ctx context.Context, attachmentID, status, note, actorID string) error {
	var err error
	if status == "approved" || status == "approved_with_comments" || status == "rejected" {
		_, err = r.db.ExecContext(ctx, `UPDATE node_instance_slot_attachments SET 
			status=$1, review_note=$2, approved_by=$3, approved_at=now() WHERE id=$4`, status, note, actorID, attachmentID)
	} else {
		_, err = r.db.ExecContext(ctx, `UPDATE node_instance_slot_attachments SET 
			status=$1, review_note=$2, approved_by=NULL, approved_at=NULL WHERE id=$3`, status, note, attachmentID)
	}
	return err
}

func (r *SQLAdminRepository) UploadReviewedDocument(ctx context.Context, attachmentID, versionID, actorID string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE node_instance_slot_attachments SET 
		reviewed_document_version_id=$1, reviewed_by=$2, reviewed_at=now() WHERE id=$3`, versionID, actorID, attachmentID)
	return err
}

func (r *SQLAdminRepository) GetLatestAttachmentStatus(ctx context.Context, instanceID string) (string, error) {
	var status string
	query := `SELECT a.status
		FROM node_instance_slot_attachments a
		JOIN node_instance_slots s ON s.id=a.slot_id
		WHERE s.node_instance_id=$1 AND a.is_active
		ORDER BY a.attached_at DESC
		LIMIT 1`
	err := r.db.GetContext(ctx, &status, query, instanceID)
	if err == sql.ErrNoRows {
		return "", nil // Not found is empty string
	}
	return status, err
}

func (r *SQLAdminRepository) GetAttachmentCounts(ctx context.Context, instanceID string) (submitted, approved, rejected int, err error) {
	query := `SELECT
		COALESCE(SUM(CASE WHEN a.status='submitted' THEN 1 ELSE 0 END),0) AS submitted,
		COALESCE(SUM(CASE WHEN a.status IN ('approved', 'approved_with_comments') THEN 1 ELSE 0 END),0) AS approved,
		COALESCE(SUM(CASE WHEN a.status='rejected' THEN 1 ELSE 0 END),0) AS rejected
		FROM node_instance_slot_attachments a
		JOIN node_instance_slots s ON s.id=a.slot_id
		WHERE s.node_instance_id=$1 AND a.is_active`
	err = r.db.QueryRowxContext(ctx, query, instanceID).Scan(&submitted, &approved, &rejected)
	return
}

func (r *SQLAdminRepository) UpdateNodeInstanceState(ctx context.Context, instanceID, state string) error {
	query := "UPDATE node_instances SET state=$1, updated_at=now() WHERE id=$2"
	if state == "submitted" || state == "under_review" {
		query = "UPDATE node_instances SET state=$1, submitted_at=COALESCE(submitted_at, now()), updated_at=now() WHERE id=$2"
	}
	_, err := r.db.ExecContext(ctx, query, state, instanceID)
	return err
}

func (r *SQLAdminRepository) UpdateAllNodeInstances(ctx context.Context, studentID, nodeID, instanceID, state string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE node_instances SET state=$1, updated_at=now() 
			WHERE user_id=$2 AND node_id=$3 AND id != $4`,
			state, studentID, nodeID, instanceID)
	return err
}

func (r *SQLAdminRepository) UpsertJourneyState(ctx context.Context, tenantID, studentID, nodeID, state string) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO journey_states (tenant_id, user_id, node_id, state)
			VALUES ($1,$2,$3,$4)
			ON CONFLICT (user_id,node_id) DO UPDATE SET state=$4, updated_at=now()`, tenantID, studentID, nodeID, state)
	return err
}

func (r *SQLAdminRepository) LogNodeEvent(ctx context.Context, instanceID, eventType, actorID string, payload map[string]any) error {
	pBytes, _ := json.Marshal(payload)
	_, err := r.db.ExecContext(ctx, `INSERT INTO node_events (id, node_instance_id, type, payload, created_by) 
		VALUES (gen_random_uuid(), $1, $2, $3, $4)`, instanceID, eventType, pBytes, actorID)
	return err
}

func (r *SQLAdminRepository) CreateReminders(ctx context.Context, studentIDs []string, title, message string, dueAt *string, createdBy string) error {
	if len(studentIDs) == 0 { return nil }
	// Batch insert?
	// Using loop inside tx is safer for simple implementation, or batched query building.
	// Since expected count is low (batch select), loop is fine.
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil { return err }
	defer tx.Rollback()
	
	for _, sid := range studentIDs {
		if dueAt != nil {
			_, err = tx.ExecContext(ctx, `INSERT INTO reminders (student_id,title,message,due_at,created_by) VALUES ($1,$2,$3,$4,$5)`, sid, title, message, *dueAt, createdBy)
		} else {
			_, err = tx.ExecContext(ctx, `INSERT INTO reminders (student_id,title,message,created_by) VALUES ($1,$2,$3,$4)`, sid, title, message, createdBy)
		}
		if err != nil { return err }
	}
	return tx.Commit()
}

func (r *SQLAdminRepository) CreateNotification(ctx context.Context, recipientID, title, message, link, nType, tenantID string) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO notifications (id, recipient_id, title, message, link, type, tenant_id) 
			VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6)`, 
			recipientID, title, message, link, nType, tenantID)
	return err
}

func (r *SQLAdminRepository) CreateReviewedDocumentVersion(ctx context.Context, docID, storagePath, objKey, bucket, mimeType string, sizeBytes int64, actorID, etag, tenantID string) (string, error) {
	// Generate ID in Go to avoid "postgres inconsistent types deduced" error with RETURNING clause
	id := uuid.New().String()
	
	// Use standard Exec, no RETURNING needed
	query := `INSERT INTO document_versions (
		id, document_id, storage_path, object_key, bucket, mime_type, size_bytes, uploaded_by, etag, tenant_id
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.db.ExecContext(ctx, query, 
		id, 
		docID, 
		storagePath, 
		objKey, 
		bucket, 
		mimeType, 
		sizeBytes, 
		actorID, 
		nullableString(etag), 
		tenantID,
	)
	
	if err != nil {
		return "", err
	}
	return id, nil
}

func nullableString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

func (r *SQLAdminRepository) ListAdminNotifications(ctx context.Context, unreadOnly bool) ([]models.AdminNotification, error) {
	query := `
		SELECT 
			n.id,
			n.student_id,
			COALESCE(u.first_name || ' ' || u.last_name, 'Unknown') as student_name,
			COALESCE(u.email, '') as student_email,
			n.node_id,
			COALESCE(n.node_instance_id::text, '') as node_instance_id,
			n.event_type,
			n.is_read,
			n.message,
			n.metadata::text as metadata,
			to_char(n.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') as created_at
		FROM admin_notifications n
		JOIN users u ON u.id = n.student_id
	`

	if unreadOnly {
		query += " WHERE n.is_read = false"
	}

	query += " ORDER BY n.created_at DESC LIMIT 100"

	var notifications []models.AdminNotification
	err := r.db.SelectContext(ctx, &notifications, query)
	if err != nil {
		return nil, err
	}
	if notifications == nil {
		notifications = []models.AdminNotification{}
	}
	return notifications, nil
}

// Fixed signature to match Interface
func (r *SQLAdminRepository) GetAdminUnreadCount(ctx context.Context) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM admin_notifications WHERE is_read = false")
	return count, err
}

func (r *SQLAdminRepository) MarkAdminNotificationRead(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE admin_notifications SET is_read = true WHERE id = $1", id)
	return err
}

func (r *SQLAdminRepository) MarkAllAdminNotificationsRead(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, "UPDATE admin_notifications SET is_read = true WHERE is_read = false")
	return err
}


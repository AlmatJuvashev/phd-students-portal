package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type JourneyRepository interface {
	// State
	GetJourneyState(ctx context.Context, userID, tenantID string) (map[string]string, error)
	UpsertJourneyState(ctx context.Context, userID, nodeID, state, tenantID string) error
	ResetJourney(ctx context.Context, userID, tenantID string) error
	
	// Scoreboard
	GetDoneNodes(ctx context.Context, tenantID string) ([]models.JourneyState, error)
	GetUsersByIDs(ctx context.Context, ids []string) ([]models.User, error)
	
	// Node Instances
	GetNodeInstance(ctx context.Context, userID, nodeID string) (*models.NodeInstance, error)
	GetNodeInstanceByID(ctx context.Context, instanceID string) (*models.NodeInstance, error)
	CreateNodeInstance(ctx context.Context, tenantID, userID, versionID, nodeID, state string, locale *string) (string, error)
	UpdateNodeInstanceState(ctx context.Context, instanceID, oldState, newState string) error
	GetAllowedTransitionRoles(ctx context.Context, fromState, toState string) ([]string, error)
	
	// Submissions
	GetNodeInstanceSlots(ctx context.Context, instanceID string) ([]models.NodeInstanceSlot, error)
	GetNodeInstanceAttachments(ctx context.Context, instanceID string) ([]models.NodeInstanceSlotAttachment, error)
	GetFullSubmissionSlots(ctx context.Context, instanceID string) ([]models.SubmissionSlotDTO, error)
	GetNodeOutcomes(ctx context.Context, instanceID string) ([]models.NodeOutcome, error)
	UpsertSubmission(ctx context.Context, instanceID string, currentRev int, locale *string) error
	GetFormRevision(ctx context.Context, instanceID string, rev int) ([]byte, error)
	InsertFormRevision(ctx context.Context, instanceID string, rev int, data []byte, editedBy string) error
	InsertOutcome(ctx context.Context, instanceID, value, decidedBy, note string) error
	
	// Events
	LogNodeEvent(ctx context.Context, instanceID, eventType, actorID string, payload map[string]any) error
	
	// Slots/Attachments
	CreateSlot(ctx context.Context, instanceID, slotKey, tenantID string, required bool, multiplicity string, mime []string) (string, error)
	GetSlot(ctx context.Context, instanceID, slotKey string) (*models.NodeInstanceSlot, error)
	CreateAttachment(ctx context.Context, slotID, docVerID, status, filename, attachedBy string, sizeBytes int64) (string, error)
	
	// Transaction
	WithTx(ctx context.Context, fn func(repo JourneyRepository) error) error
}

type SQLJourneyRepository struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func NewSQLJourneyRepository(db *sqlx.DB) *SQLJourneyRepository {
	return &SQLJourneyRepository{db: db}
}

// Helper to create a repo from transaction
func (r *SQLJourneyRepository) withTx(tx *sqlx.Tx) *SQLJourneyRepository {
	return &SQLJourneyRepository{db: r.db, tx: tx}
}

// Queries use r.q which selects db or tx
func (r *SQLJourneyRepository) q() sqlx.ExtContext {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

// WithTx executes function in transaction
func (r *SQLJourneyRepository) WithTx(ctx context.Context, fn func(JourneyRepository) error) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	repoTx := r.withTx(tx)
	if err := fn(repoTx); err != nil {
		return err
	}
	
	return tx.Commit()
}

// GetJourneyState returns map of nodeID -> state
func (r *SQLJourneyRepository) GetJourneyState(ctx context.Context, userID, tenantID string) (map[string]string, error) {
	rows, err := r.q().QueryxContext(ctx, `SELECT node_id, state FROM journey_states WHERE user_id=$1 AND tenant_id=$2`, userID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	m := make(map[string]string)
	for rows.Next() {
		var nid, st string
		if err := rows.Scan(&nid, &st); err == nil {
			m[nid] = st
		}
	}
	return m, nil
}

// UpsertJourneyState inserts or updates state
func (r *SQLJourneyRepository) UpsertJourneyState(ctx context.Context, userID, nodeID, state, tenantID string) error {
	_, err := r.q().ExecContext(ctx, `INSERT INTO journey_states (user_id, node_id, state, tenant_id, updated_at) 
		VALUES ($1, $2, $3, $4, now())
		ON CONFLICT (user_id, node_id) DO UPDATE SET state=$3, updated_at=now()`, userID, nodeID, state, tenantID)
	return err
}

// ResetJourney deletes journey data excluding profile
func (r *SQLJourneyRepository) ResetJourney(ctx context.Context, userID, tenantID string) error {
	_, err := r.q().ExecContext(ctx, `DELETE FROM node_instances WHERE user_id=$1 AND tenant_id=$2 AND node_id <> 'S1_profile'`, userID, tenantID)
	if err != nil {
		return err
	}
	_, err = r.q().ExecContext(ctx, `DELETE FROM journey_states WHERE user_id=$1 AND tenant_id=$2 AND node_id <> 'S1_profile'`, userID, tenantID)
	return err
}

// GetDoneNodes fetches all 'done' states for tenant
func (r *SQLJourneyRepository) GetDoneNodes(ctx context.Context, tenantID string) ([]models.JourneyState, error) {
    var nodes []models.JourneyState
    err := sqlx.SelectContext(ctx, r.q(), &nodes, `SELECT user_id, node_id FROM journey_states WHERE state='done' AND tenant_id=$1`, tenantID)
    return nodes, err
}

// GetUsersByIDs fetches user details
func (r *SQLJourneyRepository) GetUsersByIDs(ctx context.Context, ids []string) ([]models.User, error) {
    if len(ids) == 0 {
        return nil, nil
    }
    query, args, err := sqlx.In(`SELECT id, email, first_name, last_name, avatar_url FROM users WHERE id IN (?)`, ids)
    if err != nil {
        return nil, err
    }
    query = r.db.Rebind(query)
    var users []models.User
    err = sqlx.SelectContext(ctx, r.q(), &users, query, args...)
    return users, err
}

// GetNodeInstance returns single instance by user/node
func (r *SQLJourneyRepository) GetNodeInstance(ctx context.Context, userID, nodeID string) (*models.NodeInstance, error) {
	var inst models.NodeInstance
	err := sqlx.GetContext(ctx, r.q(), &inst, `SELECT * FROM node_instances WHERE user_id=$1 AND node_id=$2 ORDER BY updated_at DESC LIMIT 1`, userID, nodeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found is not error
		}
		return nil, err
	}
	return &inst, nil
}

func (r *SQLJourneyRepository) GetNodeInstanceByID(ctx context.Context, instanceID string) (*models.NodeInstance, error) {
	var inst models.NodeInstance
	err := sqlx.GetContext(ctx, r.q(), &inst, `SELECT * FROM node_instances WHERE id=$1`, instanceID)
	return &inst, err
}

func (r *SQLJourneyRepository) CreateNodeInstance(ctx context.Context, tenantID, userID, versionID, nodeID, state string, locale *string) (string, error) {
	var id string
	err := r.q().QueryRowxContext(ctx, `INSERT INTO node_instances (tenant_id, user_id, playbook_version_id, node_id, state, locale, opened_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, now(), now()) RETURNING id`, tenantID, userID, versionID, nodeID, state, locale).Scan(&id)
	return id, err
}

func (r *SQLJourneyRepository) UpdateNodeInstanceState(ctx context.Context, instanceID, oldState, newState string) error {
	res, err := r.q().ExecContext(ctx, `UPDATE node_instances SET state=$1, updated_at=now() WHERE id=$2 AND state=$3`, newState, instanceID, oldState)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows // Or custom error "state changed or not found"
	}
	return nil
}

func (r *SQLJourneyRepository) GetAllowedTransitionRoles(ctx context.Context, fromState, toState string) ([]string, error) {
	var roles pq.StringArray
	err := r.q().QueryRowxContext(ctx, `SELECT allowed_roles FROM node_state_transitions WHERE from_state=$1 AND to_state=$2`, fromState, toState).Scan(&roles)
	if err != nil {
		return nil, err
	}
	return []string(roles), nil
}

func (r *SQLJourneyRepository) GetNodeInstanceSlots(ctx context.Context, instanceID string) ([]models.NodeInstanceSlot, error) {
	var slots []models.NodeInstanceSlot
	err := sqlx.SelectContext(ctx, r.q(), &slots, `SELECT * FROM node_instance_slots WHERE node_instance_id=$1`, instanceID)
	return slots, err
}

func (r *SQLJourneyRepository) GetNodeInstanceAttachments(ctx context.Context, instanceID string) ([]models.NodeInstanceSlotAttachment, error) {
	var atts []models.NodeInstanceSlotAttachment
	query := `SELECT a.* 
		FROM node_instance_slot_attachments a
		JOIN node_instance_slots s ON s.id = a.slot_id
		WHERE s.node_instance_id = $1 AND a.is_active = true`
	err := sqlx.SelectContext(ctx, r.q(), &atts, query, instanceID)
	return atts, err
}

func (r *SQLJourneyRepository) GetFullSubmissionSlots(ctx context.Context, instanceID string) ([]models.SubmissionSlotDTO, error) {
	rows, err := r.q().QueryxContext(ctx, `SELECT s.id, s.slot_key, s.required, s.multiplicity, s.mime_whitelist,
		a.id AS attachment_id, a.document_version_id, a.filename, a.size_bytes, a.attached_at, a.is_active,
		a.status, a.review_note, a.approved_at, a.approved_by,
		a.reviewed_document_version_id, a.reviewed_at, a.reviewed_by,
		rdv.size_bytes AS reviewed_size_bytes, rdv.mime_type AS reviewed_mime_type,
		COALESCE(ru.first_name || ' ' || ru.last_name, '') AS reviewed_by_name
		FROM node_instance_slots s
		LEFT JOIN node_instance_slot_attachments a ON a.slot_id=s.id
		LEFT JOIN document_versions rdv ON rdv.id=a.reviewed_document_version_id
		LEFT JOIN users ru ON ru.id=a.reviewed_by
		WHERE s.node_instance_id=$1
		ORDER BY s.slot_key, a.attached_at DESC`, instanceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Temporary struct for scanning join result
	type fullRow struct {
		SlotID                    string         `db:"id"`
		SlotKey                   string         `db:"slot_key"`
		Required                  bool           `db:"required"`
		Multiplicity              string         `db:"multiplicity"`
		Mime                      pq.StringArray `db:"mime_whitelist"`
		
		AttachmentID              sql.NullString `db:"attachment_id"`
		VersionID                 sql.NullString `db:"document_version_id"`
		Filename                  sql.NullString `db:"filename"`
		SizeBytes                 sql.NullInt64  `db:"size_bytes"`
		AttachedAt                sql.NullTime   `db:"attached_at"`
		IsActive                  sql.NullBool   `db:"is_active"`
		Status                    sql.NullString `db:"status"`
		ReviewNote                sql.NullString `db:"review_note"`
		ApprovedAt                sql.NullTime   `db:"approved_at"`
		ApprovedBy                sql.NullString `db:"approved_by"`
		
		ReviewedDocumentVersionID sql.NullString `db:"reviewed_document_version_id"`
		ReviewedAt                sql.NullTime   `db:"reviewed_at"`
		ReviewedBy                sql.NullString `db:"reviewed_by"`
		ReviewedSizeBytes         sql.NullInt64  `db:"reviewed_size_bytes"`
		ReviewedMimeType          sql.NullString `db:"reviewed_mime_type"`
		ReviewedByName            sql.NullString `db:"reviewed_by_name"`
	}

	slotMap := make(map[string]*models.SubmissionSlotDTO)
	// Maintain order
	var orderedKeys []string

	for rows.Next() {
		var r fullRow
		if err := rows.StructScan(&r); err != nil {
			return nil, err
		}
		
		key := r.SlotKey
		slot, exists := slotMap[key]
		if !exists {
			slot = &models.SubmissionSlotDTO{
				ID:           r.SlotID,
				SlotKey:      r.SlotKey,
				Required:     r.Required,
				Multiplicity: r.Multiplicity,
				Mime:         []string(r.Mime),
				Attachments:  []models.SubmissionAttachmentDTO{},
			}
			slotMap[key] = slot
			orderedKeys = append(orderedKeys, key)
		}
		
		if r.AttachmentID.Valid && r.VersionID.Valid {
			att := models.SubmissionAttachmentDTO{
				AttachmentID:      r.AttachmentID.String,
				DocumentVersionID: r.VersionID.String,
				Filename:          r.Filename.String,
				SizeBytes:         r.SizeBytes.Int64,
				IsActive:          r.IsActive.Bool,
			}
			if r.AttachedAt.Valid { t := r.AttachedAt.Time; att.AttachedAt = &t }
			if r.Status.Valid { s := r.Status.String; att.Status = &s }
			if r.ReviewNote.Valid { s := r.ReviewNote.String; att.ReviewNote = &s }
			if r.ApprovedAt.Valid { t := r.ApprovedAt.Time; att.ApprovedAt = &t }
			if r.ApprovedBy.Valid { s := r.ApprovedBy.String; att.ApprovedBy = &s }
			// Reviewed doc
			if r.ReviewedDocumentVersionID.Valid {
				s := r.ReviewedDocumentVersionID.String
				att.ReviewedDocumentVersionID = &s
			}
			if r.ReviewedAt.Valid { t := r.ReviewedAt.Time; att.ReviewedAt = &t }
			if r.ReviewedBy.Valid { s := r.ReviewedBy.String; att.ReviewedBy = &s }
			if r.ReviewedSizeBytes.Valid { i := r.ReviewedSizeBytes.Int64; att.ReviewedSizeBytes = &i }
			if r.ReviewedMimeType.Valid { s := r.ReviewedMimeType.String; att.ReviewedMimeType = &s }
			if r.ReviewedByName.Valid { s := r.ReviewedByName.String; att.ReviewedByName = &s }
			
			slot.Attachments = append(slot.Attachments, att)
		}
	}
	
	// Convert map to slice
	res := make([]models.SubmissionSlotDTO, 0, len(orderedKeys))
	for _, k := range orderedKeys {
		res = append(res, *slotMap[k])
	}
	return res, nil
}

func (r *SQLJourneyRepository) GetNodeOutcomes(ctx context.Context, instanceID string) ([]models.NodeOutcome, error) {
	var outs []models.NodeOutcome
	err := sqlx.SelectContext(ctx, r.q(), &outs, `SELECT outcome_value, decided_by, note, created_at FROM node_outcomes WHERE node_instance_id=$1 ORDER BY created_at DESC`, instanceID)
	return outs, err
}

func (r *SQLJourneyRepository) GetFormRevision(ctx context.Context, instanceID string, rev int) ([]byte, error) {
	var data []byte
	err := r.q().QueryRowxContext(ctx, `SELECT form_data FROM node_instance_form_revisions WHERE node_instance_id=$1 AND rev=$2`, instanceID, rev).Scan(&data)
	return data, err
}

func (r *SQLJourneyRepository) InsertFormRevision(ctx context.Context, instanceID string, rev int, data []byte, editedBy string) error {
	_, err := r.q().ExecContext(ctx, `INSERT INTO node_instance_form_revisions (node_instance_id, rev, form_data, edited_by, created_at) VALUES ($1,$2,$3,$4,now())`, instanceID, rev, data, editedBy)
	return err
}

func (r *SQLJourneyRepository) InsertOutcome(ctx context.Context, instanceID, value, decidedBy, note string) error {
	_, err := r.q().ExecContext(ctx, `INSERT INTO node_outcomes (node_instance_id, outcome_value, decided_by, note, created_at) VALUES ($1,$2,$3,$4,now())`, instanceID, value, decidedBy, note)
	return err
}

func (r *SQLJourneyRepository) UpsertSubmission(ctx context.Context, instanceID string, currentRev int, locale *string) error {
	_, err := r.q().ExecContext(ctx, `UPDATE node_instances SET current_rev=$1, locale=$2, updated_at=now() WHERE id=$3`, currentRev, locale, instanceID)
	return err
}

func (r *SQLJourneyRepository) LogNodeEvent(ctx context.Context, instanceID, eventType, actorID string, payload map[string]any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = r.q().ExecContext(ctx, `INSERT INTO node_events (node_instance_id, event_type, payload, actor_id) VALUES ($1, $2, $3, $4)`, instanceID, eventType, data, actorID)
	return err
}

func (r *SQLJourneyRepository) CreateSlot(ctx context.Context, instanceID, slotKey, tenantID string, required bool, multiplicity string, mime []string) (string, error) {
	var id string
	mimeArr := pq.Array(mime)
	if len(mime) == 0 {
		mimeArr = pq.Array([]string{})
	}
	err := r.q().QueryRowxContext(ctx, `INSERT INTO node_instance_slots (node_instance_id, slot_key, tenant_id, required, multiplicity, mime_whitelist) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, instanceID, slotKey, tenantID, required, multiplicity, mimeArr).Scan(&id)
	return id, err
}

func (r *SQLJourneyRepository) GetSlot(ctx context.Context, instanceID, slotKey string) (*models.NodeInstanceSlot, error) {
	var slot models.NodeInstanceSlot
	err := sqlx.GetContext(ctx, r.q(), &slot, `SELECT * FROM node_instance_slots WHERE node_instance_id=$1 AND slot_key=$2`, instanceID, slotKey)
	if err != nil {
		return nil, err
	}
	return &slot, nil
}

func (r *SQLJourneyRepository) CreateAttachment(ctx context.Context, slotID, docVerID, status, filename, attachedBy string, sizeBytes int64) (string, error) {
	var id string
	err := r.q().QueryRowxContext(ctx, `INSERT INTO node_instance_slot_attachments (slot_id, document_version_id, is_active, status, filename, attached_by, size_bytes, attached_at)
		VALUES ($1, $2, true, $3, $4, $5, $6, now()) RETURNING id`, slotID, docVerID, status, filename, attachedBy, sizeBytes).Scan(&id)
	return id, err
}

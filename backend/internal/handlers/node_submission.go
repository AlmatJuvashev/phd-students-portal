package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type NodeSubmissionHandler struct {
	db  *sqlx.DB
	cfg config.AppConfig
	pb  *playbook.Manager
}

func NewNodeSubmissionHandler(db *sqlx.DB, cfg config.AppConfig, pb *playbook.Manager) *NodeSubmissionHandler {
	return &NodeSubmissionHandler{db: db, cfg: cfg, pb: pb}
}

type nodeInstanceRecord struct {
	ID         string `db:"id"`
	NodeID     string `db:"node_id"`
	State      string `db:"state"`
	CurrentRev int    `db:"current_rev"`
	Locale     string `db:"locale"`
}

// GET /api/journey/nodes/:nodeId/submission
func (h *NodeSubmissionHandler) GetSubmission(c *gin.Context) {
	nodeID := c.Param("nodeId")
	uid := userIDFromClaims(c)
	if uid == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	locale := h.resolveLocale(c.Query("locale"))

	instance, err := h.ensureNodeInstance(c, uid, nodeID, locale)
	if err != nil {
		handleNodeErr(c, err)
		return
	}
	dto, err := h.buildSubmissionDTO(instance.ID)
	if err != nil {
		handleNodeErr(c, err)
		return
	}
	c.JSON(200, dto)
}

// GET /api/journey/profile
func (h *NodeSubmissionHandler) GetProfile(c *gin.Context) {
	uid := userIDFromClaims(c)
	if uid == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	var form json.RawMessage
	err := h.db.QueryRow(`SELECT form_data FROM profile_submissions WHERE user_id=$1`, uid).Scan(&form)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(404, gin.H{"error": "not_found"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Data(200, "application/json", form)
}

type submissionReq struct {
	FormData json.RawMessage `json:"form_data"`
	State    string          `json:"state"`
}

// PUT /api/journey/nodes/:nodeId/submission
func (h *NodeSubmissionHandler) PutSubmission(c *gin.Context) {
	nodeID := c.Param("nodeId")
	uid := userIDFromClaims(c)
	if uid == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	var req submissionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	locale := h.resolveLocale(c.Query("locale"))

	var err error
	err = h.withTx(func(tx *sqlx.Tx) error {
		inst, err := h.ensureNodeInstanceTx(tx, uid, nodeID, locale)
		if err != nil {
			return err
		}
		if len(req.FormData) != 0 {
			if nodeID == "S1_publications_list" {
				sanitized, _, err := normalizeApp7Payload(req.FormData)
				if err != nil {
					return err
				}
				req.FormData = sanitized
			}
			if err := h.appendFormRevision(tx, inst, uid, req.FormData); err != nil {
				return err
			}
			if nodeID == "S1_profile" {
				if err := h.upsertProfileSubmission(tx, uid, req.FormData); err != nil {
					return err
				}
			}
		}
		if req.State != "" && req.State != inst.State {
			role := roleFromContext(c)
			if err := h.transitionState(tx, inst, uid, role, req.State); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		handleNodeErr(c, err)
		return
	}
	// reload dto
	inst, _ := h.loadInstance(uid, nodeID)
	dto, err := h.buildSubmissionDTO(inst.ID)
	if err != nil {
		handleNodeErr(c, err)
		return
	}
	c.JSON(200, dto)
}

type nodeUploadPresignReq struct {
	SlotKey     string `json:"slot_key" binding:"required"`
	Filename    string `json:"filename" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
	SizeBytes   int64  `json:"size_bytes" binding:"required"`
}

// POST /api/journey/nodes/:nodeId/uploads/presign
func (h *NodeSubmissionHandler) PresignUpload(c *gin.Context) {
	nodeID := c.Param("nodeId")
	uid := userIDFromClaims(c)
	if uid == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	var req nodeUploadPresignReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	maxBytes := int64(h.cfg.FileUploadMaxMB) * 1024 * 1024
	if maxBytes > 0 && req.SizeBytes > maxBytes {
		c.JSON(400, gin.H{"error": fmt.Sprintf("file too large (max %d MB)", h.cfg.FileUploadMaxMB)})
		return
	}
	locale := h.resolveLocale(c.Query("locale"))
	var docID string
	var instanceID string
	var err error
	err = h.withTx(func(tx *sqlx.Tx) error {
		inst, err := h.ensureNodeInstanceTx(tx, uid, nodeID, locale)
		if err != nil {
			return err
		}
		slot, err := h.getSlot(tx, inst.ID, req.SlotKey)
		if err != nil {
			return err
		}
		if len(slot.MimeWhitelist) > 0 {
			valid := false
			for _, m := range slot.MimeWhitelist {
				if strings.EqualFold(m, req.ContentType) {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("invalid mime type")
			}
		}
		instanceID = inst.ID
		docID, err = h.ensureDocumentForSlot(tx, uid, nodeID, req.SlotKey)
		return err
	})
	if err != nil {
		handleNodeErr(c, err)
		return
	}
	s3c, err := services.NewS3FromEnv()
	if err != nil {
		handleNodeErr(c, err)
		return
	}
	if s3c == nil {
		c.JSON(400, gin.H{"error": "S3 not configured"})
		return
	}
	objectKey := storage.BuildNodeObjectKey(uid, nodeID, req.SlotKey, req.Filename)
	expires := time.Minute * 15
	url, err := s3c.PresignPut(objectKey, req.ContentType, expires)
	if err != nil {
		handleNodeErr(c, err)
		return
	}
	_ = instanceID // currently unused but reserved for future audit
	resp := gin.H{
		"upload_url":       url,
		"object_key":       objectKey,
		"document_id":      docID,
		"bucket":           s3c.Bucket(),
		"expires_in":       int(expires.Seconds()),
		"max_size_bytes":   maxBytes,
		"required_headers": map[string]string{"Content-Type": req.ContentType},
	}
	c.JSON(200, resp)
}

type nodeAttachReq struct {
	SlotKey     string `json:"slot_key" binding:"required"`
	Filename    string `json:"filename" binding:"required"`
	ObjectKey   string `json:"object_key" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
	SizeBytes   int64  `json:"size_bytes" binding:"required"`
	ETag        string `json:"etag"`
}

// POST /api/journey/nodes/:nodeId/uploads/attach
func (h *NodeSubmissionHandler) AttachUpload(c *gin.Context) {
	nodeID := c.Param("nodeId")
	uid := userIDFromClaims(c)
	if uid == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	role := roleFromContext(c)
	var req nodeAttachReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	s3c, err := services.NewS3FromEnv()
	if err != nil {
		handleNodeErr(c, err)
		return
	}
	if s3c == nil {
		c.JSON(400, gin.H{"error": "S3 not configured"})
		return
	}
	bucket := s3c.Bucket()
	locale := h.resolveLocale(c.Query("locale"))
	err = h.withTx(func(tx *sqlx.Tx) error {
		inst, err := h.ensureNodeInstanceTx(tx, uid, nodeID, locale)
		if err != nil {
			return err
		}
		slot, err := h.getSlot(tx, inst.ID, req.SlotKey)
		if err != nil {
			return err
		}
		docID, err := h.ensureDocumentForSlot(tx, uid, nodeID, req.SlotKey)
		if err != nil {
			return err
		}
		var versionID string
		err = tx.QueryRowx(`INSERT INTO document_versions (document_id, storage_path, object_key, bucket, mime_type, size_bytes, uploaded_by, etag)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
			docID, req.ObjectKey, req.ObjectKey, bucket, req.ContentType, req.SizeBytes, uid, nullableString(req.ETag)).Scan(&versionID)
		if err != nil {
			return err
		}
		_, err = tx.Exec(`UPDATE documents SET current_version_id=$1 WHERE id=$2`, versionID, docID)
		if err != nil {
			return err
		}
		if slot.Multiplicity == "single" {
			if _, err := tx.Exec(`UPDATE node_instance_slot_attachments SET is_active=false WHERE slot_id=$1 AND is_active=true`, slot.ID); err != nil {
				return err
			}
		}
		_, err = tx.Exec(`INSERT INTO node_instance_slot_attachments (slot_id, document_version_id, filename, size_bytes, attached_by, status)
			VALUES ($1,$2,$3,$4,$5,'submitted')`, slot.ID, versionID, req.Filename, req.SizeBytes, uid)
		if err != nil {
			return err
		}
		if err := h.insertEvent(tx, inst.ID, "file_attached", uid, map[string]any{"slot_key": req.SlotKey, "version_id": versionID}); err != nil {
			return err
		}
		if inst.State == "active" {
			_ = h.transitionState(tx, inst, uid, role, "submitted")
		}
		return nil
	})
	if err != nil {
		handleNodeErr(c, err)
		return
	}
	inst, _ := h.loadInstance(uid, nodeID)
	dto, err := h.buildSubmissionDTO(inst.ID)
	if err != nil {
		handleNodeErr(c, err)
		return
	}
	c.JSON(200, dto)
}

type stateReq struct {
	State string `json:"state" binding:"required"`
}

// PATCH /api/journey/nodes/:nodeId/state
func (h *NodeSubmissionHandler) PatchState(c *gin.Context) {
	nodeID := c.Param("nodeId")
	uid := userIDFromClaims(c)
	if uid == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	var req stateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	locale := h.resolveLocale(c.Query("locale"))
	err := h.withTx(func(tx *sqlx.Tx) error {
		inst, err := h.ensureNodeInstanceTx(tx, uid, nodeID, locale)
		if err != nil {
			return err
		}
		role := roleFromContext(c)
		if err := h.transitionState(tx, inst, uid, role, req.State); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		handleNodeErr(c, err)
		return
	}
	inst, _ := h.loadInstance(uid, nodeID)
	dto, err := h.buildSubmissionDTO(inst.ID)
	if err != nil {
		handleNodeErr(c, err)
		return
	}
	c.JSON(200, dto)
}

func (h *NodeSubmissionHandler) withTx(fn func(tx *sqlx.Tx) error) error {
	tx, err := h.db.Beginx()
	if err != nil {
		return err
	}
	err = fn(tx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (h *NodeSubmissionHandler) ensureNodeInstance(c *gin.Context, userID, nodeID, locale string) (*nodeInstanceRecord, error) {
	tx, err := h.db.Beginx()
	if err != nil {
		return nil, err
	}
	inst, err := h.ensureNodeInstanceTx(tx, userID, nodeID, locale)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return inst, nil
}

func (h *NodeSubmissionHandler) ensureNodeInstanceTx(tx *sqlx.Tx, userID, nodeID, locale string) (*nodeInstanceRecord, error) {
    inst, err := h.loadInstanceTx(tx, userID, nodeID)
    if err == nil {
        // Backfill missing upload slots if playbook now defines them
        if err := h.ensureSlotsForInstance(tx, inst.ID, nodeID); err != nil {
            return nil, err
        }
        return inst, nil
    }
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	nodeDef, ok := h.pb.NodeDefinition(nodeID)
	if !ok {
		return nil, fmt.Errorf("node not found in playbook")
	}
	var rec nodeInstanceRecord
	err = tx.QueryRowx(`INSERT INTO node_instances (user_id, playbook_version_id, node_id, state, locale)
        VALUES ($1,$2,$3,'active',$4)
        RETURNING id, node_id, state, current_rev, locale`, userID, h.pb.VersionID, nodeID, locale).Scan(&rec.ID, &rec.NodeID, &rec.State, &rec.CurrentRev, &rec.Locale)
	if err != nil {
		return nil, err
	}
    if nodeDef.Requirements != nil {
        for _, up := range nodeDef.Requirements.Uploads {
            req := up.Required
            mime := pq.Array(up.Mime)
            if len(up.Mime) == 0 {
                mime = pq.Array([]string{})
            }
            _, err := tx.Exec(`INSERT INTO node_instance_slots (node_instance_id, slot_key, required, multiplicity, mime_whitelist)
                VALUES ($1,$2,$3,'single',$4)`, rec.ID, up.Key, req, mime)
            if err != nil {
                return nil, err
            }
        }
    }
	if err := h.insertEvent(tx, rec.ID, "opened", userID, map[string]any{"locale": locale}); err != nil {
		return nil, err
	}
	_, _ = tx.Exec(`INSERT INTO journey_states (user_id, node_id, state)
        VALUES ($1,$2,'active')
        ON CONFLICT (user_id, node_id) DO UPDATE SET state='active', updated_at=now()`, userID, nodeID)
	return &rec, nil
}

// ensureSlotsForInstance inserts any missing slot rows for an existing node instance
// based on the current playbook definition (useful after adding requirements.uploads
// to a node in a newer playbook version).
func (h *NodeSubmissionHandler) ensureSlotsForInstance(tx *sqlx.Tx, instanceID, nodeID string) error {
    nodeDef, ok := h.pb.NodeDefinition(nodeID)
    if !ok || nodeDef.Requirements == nil || len(nodeDef.Requirements.Uploads) == 0 {
        return nil
    }
    // Load existing slot keys for this instance
    var existing []string
    if err := tx.Select(&existing, `SELECT slot_key FROM node_instance_slots WHERE node_instance_id=$1`, instanceID); err != nil {
        return err
    }
    present := map[string]struct{}{}
    for _, k := range existing {
        present[k] = struct{}{}
    }
    // Insert any missing slots
    for _, up := range nodeDef.Requirements.Uploads {
        if _, ok := present[up.Key]; ok {
            continue
        }
        mime := pq.Array(up.Mime)
        if len(up.Mime) == 0 {
            mime = pq.Array([]string{})
        }
        if _, err := tx.Exec(`INSERT INTO node_instance_slots (node_instance_id, slot_key, required, multiplicity, mime_whitelist)
            VALUES ($1,$2,$3,'single',$4)`, instanceID, up.Key, up.Required, mime); err != nil {
            return err
        }
    }
    return nil
}

func (h *NodeSubmissionHandler) loadInstanceTx(tx *sqlx.Tx, userID, nodeID string) (*nodeInstanceRecord, error) {
	var rec nodeInstanceRecord
	err := tx.QueryRowx(`SELECT id, node_id, state, current_rev, locale FROM node_instances WHERE user_id=$1 AND playbook_version_id=$2 AND node_id=$3`, userID, h.pb.VersionID, nodeID).Scan(&rec.ID, &rec.NodeID, &rec.State, &rec.CurrentRev, &rec.Locale)
	if err != nil {
		return nil, err
	}
	return &rec, nil
}

func (h *NodeSubmissionHandler) loadInstance(userID, nodeID string) (*nodeInstanceRecord, error) {
	var rec nodeInstanceRecord
	err := h.db.QueryRowx(`SELECT id, node_id, state, current_rev, locale FROM node_instances WHERE user_id=$1 AND playbook_version_id=$2 AND node_id=$3`, userID, h.pb.VersionID, nodeID).Scan(&rec.ID, &rec.NodeID, &rec.State, &rec.CurrentRev, &rec.Locale)
	if err != nil {
		return nil, err
	}
	return &rec, nil
}

func (h *NodeSubmissionHandler) appendFormRevision(tx *sqlx.Tx, inst *nodeInstanceRecord, userID string, data json.RawMessage) error {
	nextRev := inst.CurrentRev + 1
	if _, err := tx.Exec(`INSERT INTO node_instance_form_revisions (node_instance_id, rev, form_data, edited_by)
        VALUES ($1,$2,$3,$4)`, inst.ID, nextRev, data, userID); err != nil {
		return err
	}
	if _, err := tx.Exec(`UPDATE node_instances SET current_rev=$1, updated_at=now() WHERE id=$2`, nextRev, inst.ID); err != nil {
		return err
	}
	inst.CurrentRev = nextRev
	return h.insertEvent(tx, inst.ID, "draft_saved", userID, map[string]any{"rev": nextRev})
}

func (h *NodeSubmissionHandler) upsertProfileSubmission(tx *sqlx.Tx, userID string, data json.RawMessage) error {
	_, err := tx.Exec(`INSERT INTO profile_submissions (user_id, form_data)
        VALUES ($1, $2)
        ON CONFLICT (user_id)
        DO UPDATE SET form_data = EXCLUDED.form_data, updated_at = NOW()`, userID, data)
	return err
}

func (h *NodeSubmissionHandler) transitionState(tx *sqlx.Tx, inst *nodeInstanceRecord, userID, role, newState string) error {
	roles := []string{}

	err := tx.QueryRowx(`SELECT allowed_roles FROM node_state_transitions WHERE from_state=$1 AND to_state=$2`, inst.State, newState).Scan(pq.Array(&roles))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Iteration override: allow students to complete nodes directly to done
			if role == "student" && newState == "done" && (inst.State == "active" || inst.State == "submitted") {
				// treat as allowed below
				roles = []string{role}
			} else {
				return fmt.Errorf("transition not allowed")
			}
		} else {
			return err
		}
	}
	allowed := false
	for _, r := range roles {
		if r == role {
			allowed = true
			break
		}
	}
	// Additional safety: keep override even if roles were loaded but current role not listed
	if !allowed && role == "student" && newState == "done" && (inst.State == "active" || inst.State == "submitted") {
		allowed = true
	}
	if !allowed {
		return fmt.Errorf("role %s cannot transition from %s to %s", role, inst.State, newState)
	}
	previous := inst.State
	query := "UPDATE node_instances SET state=$1, updated_at=now() WHERE id=$2"
	if newState == "submitted" {
		query = "UPDATE node_instances SET state=$1, submitted_at=now(), updated_at=now() WHERE id=$2"
	}
	if _, err := tx.Exec(query, newState, inst.ID); err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO journey_states (user_id,node_id,state) VALUES ($1,$2,$3)
        ON CONFLICT (user_id,node_id) DO UPDATE SET state=$3, updated_at=now()`, userID, inst.NodeID, newState)
	if err != nil {
		return err
	}
	inst.State = newState
	payload := map[string]any{"from": previous, "to": newState}
	return h.insertEvent(tx, inst.ID, "state_changed", userID, payload)
}

func (h *NodeSubmissionHandler) getSlot(tx *sqlx.Tx, instanceID, slotKey string) (slotRecord, error) {
	var slot slotRecord
	err := tx.QueryRowx(`SELECT id, slot_key, required, multiplicity, mime_whitelist FROM node_instance_slots WHERE node_instance_id=$1 AND slot_key=$2`, instanceID, slotKey).Scan(&slot.ID, &slot.SlotKey, &slot.Required, &slot.Multiplicity, pq.Array(&slot.MimeWhitelist))
	if err != nil {
		return slotRecord{}, err
	}
	slot.NodeInstanceID = instanceID
	return slot, nil
}

type slotRecord struct {
	ID             string
	NodeInstanceID string
	SlotKey        string
	Required       bool
	Multiplicity   string
	MimeWhitelist  []string
}

func (h *NodeSubmissionHandler) ensureDocumentForSlot(tx *sqlx.Tx, userID, nodeID, slotKey string) (string, error) {
	title := fmt.Sprintf("node:%s:%s", nodeID, slotKey)
	var docID string
	err := tx.QueryRowx(`SELECT id FROM documents WHERE user_id=$1 AND title=$2`, userID, title).Scan(&docID)
	if err == nil {
		return docID, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}
	err = tx.QueryRowx(`INSERT INTO documents (user_id, kind, title) VALUES ($1,'node_slot',$2) RETURNING id`, userID, title).Scan(&docID)
	if err != nil {
		return "", err
	}
	return docID, nil
}

func (h *NodeSubmissionHandler) insertEvent(tx *sqlx.Tx, nodeInstanceID, eventType, actorID string, payload map[string]any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO node_events (node_instance_id, event_type, payload, actor_id) VALUES ($1,$2,$3,$4)`, nodeInstanceID, eventType, data, actorID)
	return err
}

func (h *NodeSubmissionHandler) buildSubmissionDTO(instanceID string) (gin.H, error) {
	var inst struct {
		ID         string `db:"id"`
		NodeID     string `db:"node_id"`
		State      string `db:"state"`
		Locale     string `db:"locale"`
		CurrentRev int    `db:"current_rev"`
	}
	err := h.db.QueryRowx(`SELECT id, node_id, state, locale, current_rev FROM node_instances WHERE id=$1`, instanceID).Scan(&inst.ID, &inst.NodeID, &inst.State, &inst.Locale, &inst.CurrentRev)
	if err != nil {
		return nil, err
	}
	dto := gin.H{
		"node_id":             inst.NodeID,
		"playbook_version_id": h.pb.VersionID,
		"state":               inst.State,
		"locale":              inst.Locale,
	}
	if inst.CurrentRev > 0 {
		var rev struct {
			Rev  int             `db:"rev"`
			Data json.RawMessage `db:"form_data"`
		}
		err := h.db.QueryRowx(`SELECT rev, form_data FROM node_instance_form_revisions WHERE node_instance_id=$1 AND rev=$2`, instanceID, inst.CurrentRev).Scan(&rev.Rev, &rev.Data)
		if err == nil {
			if inst.NodeID == "S1_publications_list" {
				if form, parseErr := buildApp7Form(rev.Data); parseErr == nil {
					summary := summarizeApp7(form.Sections)
					for key, val := range form.LegacyCounts {
						if val > 0 && summary[key] == 0 {
							summary[key] = val
						}
					}
					clientData := gin.H{
						"wos_scopus":  form.Sections.WosScopus,
						"kokson":      form.Sections.Kokson,
						"conferences": form.Sections.Conferences,
						"ip":          form.Sections.IP,
						"summary":     summary,
					}
					if len(form.LegacyCounts) > 0 {
						clientData["legacy_counts"] = form.LegacyCounts
					}
					dto["form"] = gin.H{"rev": rev.Rev, "data": clientData}
				} else {
					dto["form"] = gin.H{"rev": rev.Rev, "data": rev.Data}
				}
			} else {
				dto["form"] = gin.H{"rev": rev.Rev, "data": rev.Data}
			}
		}
	}
	slots, err := h.fetchSlots(instanceID)
	if err != nil {
		return nil, err
	}
	dto["slots"] = slots
	outcomes, err := h.fetchOutcomes(instanceID)
	if err != nil {
		return nil, err
	}
	if len(outcomes) > 0 {
		dto["outcomes"] = outcomes
	}
	return dto, nil
}

func (h *NodeSubmissionHandler) fetchSlots(instanceID string) ([]gin.H, error) {
	rows, err := h.db.Queryx(`SELECT s.id, s.slot_key, s.required, s.multiplicity, s.mime_whitelist,
		a.id AS attachment_id, a.document_version_id, a.filename, a.size_bytes, a.attached_at, a.is_active,
		a.status, a.review_note, a.approved_at, a.approved_by
		FROM node_instance_slots s
		LEFT JOIN node_instance_slot_attachments a ON a.slot_id=s.id
		WHERE s.node_instance_id=$1
		ORDER BY s.slot_key, a.attached_at DESC`, instanceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	type row struct {
		SlotID       string         `db:"id"`
		SlotKey      string         `db:"slot_key"`
		Required     bool           `db:"required"`
		Multiplicity string         `db:"multiplicity"`
		Mime         pq.StringArray `db:"mime_whitelist"`
		AttachmentID sql.NullString `db:"attachment_id"`
		VersionID    sql.NullString `db:"document_version_id"`
		Filename     sql.NullString `db:"filename"`
		SizeBytes    sql.NullInt64  `db:"size_bytes"`
		AttachedAt   sql.NullTime   `db:"attached_at"`
		IsActive     sql.NullBool   `db:"is_active"`
		Status       sql.NullString `db:"status"`
		ReviewNote   sql.NullString `db:"review_note"`
		ApprovedAt   sql.NullTime   `db:"approved_at"`
		ApprovedBy   sql.NullString `db:"approved_by"`
	}
	slots := []gin.H{}
	slotMap := map[string]gin.H{}
	for rows.Next() {
		var r row
		if err := rows.StructScan(&r); err != nil {
			return nil, err
		}
		key := r.SlotKey
		slot, exists := slotMap[key]
		if !exists {
			slot = gin.H{
				"key":          r.SlotKey,
				"required":     r.Required,
				"multiplicity": r.Multiplicity,
				"mime":         []string(r.Mime),
				"attachments":  []gin.H{},
			}
			slotMap[key] = slot
			slots = append(slots, slot)
		}
		if r.AttachmentID.Valid && r.VersionID.Valid {
			attachments := slot["attachments"].([]gin.H)
			att := gin.H{
				"version_id": r.VersionID.String,
				"filename":   r.Filename.String,
				"size_bytes": r.SizeBytes.Int64,
				"is_active":  r.IsActive.Bool,
			}
			if r.AttachedAt.Valid {
				att["attached_at"] = r.AttachedAt.Time.Format(time.RFC3339)
			}
			if r.Status.Valid {
				att["status"] = r.Status.String
			}
			if r.ReviewNote.Valid {
				att["review_note"] = r.ReviewNote.String
			}
			if r.ApprovedAt.Valid {
				att["approved_at"] = r.ApprovedAt.Time.Format(time.RFC3339)
			}
			if r.ApprovedBy.Valid {
				att["approved_by"] = r.ApprovedBy.String
			}
			att["download_url"] = fmt.Sprintf("/api/documents/versions/%s/download", r.VersionID.String)
			attachments = append(attachments, att)
			slot["attachments"] = attachments
		}
	}
	return slots, nil
}

func (h *NodeSubmissionHandler) fetchOutcomes(instanceID string) ([]gin.H, error) {
	rows, err := h.db.Queryx(`SELECT outcome_value, decided_by, note, created_at FROM node_outcomes WHERE node_instance_id=$1 ORDER BY created_at DESC`, instanceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var outs []gin.H
	for rows.Next() {
		var value, decidedBy, note string
		var created time.Time
		if err := rows.Scan(&value, &decidedBy, &note, &created); err != nil {
			return nil, err
		}
		outs = append(outs, gin.H{
			"value":      value,
			"decided_by": decidedBy,
			"note":       note,
			"created_at": created.Format(time.RFC3339),
		})
	}
	return outs, nil
}

func handleNodeErr(c *gin.Context, err error) {
	if err == nil {
		return
	}
	if strings.Contains(err.Error(), "not allowed") {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}
	var valErr *app7ValidationError
	if errors.As(err, &valErr) {
		c.JSON(400, gin.H{"error": valErr.Error()})
		return
	}
	switch {
	case errors.Is(err, sql.ErrNoRows):
		c.JSON(404, gin.H{"error": "not found"})
	default:
		c.JSON(500, gin.H{"error": err.Error()})
	}
}

func (h *NodeSubmissionHandler) resolveLocale(requested string) string {
	if requested != "" {
		return requested
	}
	if h.pb.DefaultLocale != "" {
		return h.pb.DefaultLocale
	}
	return "ru"
}

func roleFromContext(c *gin.Context) string {
	if val, ok := c.Get("claims"); ok {
		if claims, ok := val.(jwt.MapClaims); ok {
			if role, ok := claims["role"].(string); ok {
				return role
			}
		}
	}
	return ""
}

func nullableString(v string) interface{} {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	return v
}

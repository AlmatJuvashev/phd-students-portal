package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/mailer"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/google/uuid"
)

type JourneyService struct {
	repo    repository.JourneyRepository
	pb      *playbook.Manager
	cfg     config.AppConfig
	mailer  mailer.Mailer
	storage StorageClient
	docSvc  *DocumentService
}

func NewJourneyService(repo repository.JourneyRepository, pb *playbook.Manager, cfg config.AppConfig, mailer mailer.Mailer, storage StorageClient, docSvc *DocumentService) *JourneyService {
	return &JourneyService{
		repo:    repo,
		pb:      pb,
		cfg:     cfg,
		mailer:  mailer,
		storage: storage,
		docSvc:  docSvc,
	}
}

// GetState returns user's journey state map
func (s *JourneyService) GetState(ctx context.Context, userID, tenantID string) (map[string]string, error) {
	return s.repo.GetJourneyState(ctx, userID, tenantID)
}

// SetState upserts a state (Admin/Debug usage mainly)
func (s *JourneyService) SetState(ctx context.Context, userID, nodeID, state, tenantID string) error {
	// Validation
	allowed := map[string]bool{"locked": true, "active": true, "submitted": true, "waiting": true, "needs_fixes": true, "done": true}
	if !allowed[state] {
		return errors.New("invalid state")
	}
	return s.repo.UpsertJourneyState(ctx, userID, nodeID, state, tenantID)
}

// Reset clears user progress
func (s *JourneyService) Reset(ctx context.Context, userID, tenantID string) error {
	return s.repo.ResetJourney(ctx, userID, tenantID)
}

// Scoreboard Logic
// Scoreboard Logic
func (s *JourneyService) GetScoreboard(ctx context.Context, tenantID, currentUserID string) (*models.ScoreboardResponse, error) {
	// 1. Get all done nodes for tenant
	doneNodes, err := s.repo.GetDoneNodes(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// 2. Score Aggregation (Filter out W3)
	userScores := make(map[string]int)
	for _, dn := range doneNodes {
		// Only count known nodes
		if _, ok := s.pb.NodeDefinition(dn.NodeID); ok {
			// Check World ID
			worldID := s.pb.NodeWorldID(dn.NodeID)
			// Conditional Logic: Nodes from W3 are 0XP
			if worldID != "W3" {
				userScores[dn.UserID] += 100
			}
		}
	}

	// 3. Collect User IDs participating
	userIDs := make([]string, 0, len(userScores))
	for uid := range userScores {
		userIDs = append(userIDs, uid)
	}

	// 4. Fetch User Details for these IDs
	userInfoMap := make(map[string]models.User)
	if len(userIDs) > 0 {
		users, err := s.repo.GetUsersByIDs(ctx, userIDs)
		if err == nil {
			for _, usr := range users {
				userInfoMap[usr.ID] = usr
			}
		} else {
			log.Printf("[Scoreboard] DB Error fetching users: %v", err)
		}
	}

	// 5. Flatten to list for sorting
	var allEntries []models.ScoreboardEntry
	totalSum := 0
	for uid, score := range userScores {
		uInfo, found := userInfoMap[uid]
		name := "Unknown"
		avatar := ""
		if found {
			f := uInfo.FirstName
			l := uInfo.LastName
			name = strings.TrimSpace(f + " " + l)
			if name == "" {
				if uInfo.Email != "" {
					name = uInfo.Email
				} else {
					name = "Student"
				}
			}
			if uInfo.AvatarURL != "" {
				avatar = uInfo.AvatarURL
			}
		}

		allEntries = append(allEntries, models.ScoreboardEntry{
			UserID:     uid,
			Name:       name,
			Avatar:     avatar,
			TotalScore: score,
		})
		totalSum += score
	}

	// 6. Sort Descending
	sort.Slice(allEntries, func(i, j int) bool {
		return allEntries[i].TotalScore > allEntries[j].TotalScore
	})

	// 7. Assign Ranks
	for i := range allEntries {
		allEntries[i].Rank = i + 1
	}

	// 8. Construct Response
	var top5 []models.ScoreboardEntry
	if len(allEntries) > 5 {
		top5 = allEntries[:5]
	} else {
		top5 = allEntries
	}
	
	avg := 0
	if len(allEntries) > 0 {
		avg = totalSum / len(allEntries)
	}

	var me *models.ScoreboardEntry
	for _, e := range allEntries {
		if e.UserID == currentUserID {
			val := e
			me = &val
			break
		}
	}
	
	// If user has 0 score (no done nodes), they might not be in the list
	if me == nil {
		users, err := s.repo.GetUsersByIDs(ctx, []string{currentUserID})
		var self models.User
		if err == nil && len(users) > 0 {
			self = users[0]
		}
		
		f := self.FirstName
		l := self.LastName
		name := strings.TrimSpace(f + " " + l)
		if name == "" {
			if self.Email != "" {
				name = self.Email
			} else {
				name = "You"
			}
		}
		
		me = &models.ScoreboardEntry{
			UserID:     currentUserID,
			Name:       name,
			Avatar:     "", // can fetch if self populated correctly
			TotalScore: 0,
			Rank:       len(allEntries) + 1,
		}
		if self.AvatarURL != "" {
			me.Avatar = self.AvatarURL
		}
	}

	return &models.ScoreboardResponse{
		Top5:       top5,
		Average:    avg,
		Me:         me,
		TotalUsers: len(allEntries),
	}, nil
}

// ActivateNextNodes checks dependent nodes and activates them if all prerequisites are met
func (s *JourneyService) ActivateNextNodes(ctx context.Context, userID, completedNodeID, tenantID string) error {
	log.Printf("[ActivateNextNodes] Starting for user=%s node=%s", userID, completedNodeID)
	
	nodeDef, ok := s.pb.NodeDefinition(completedNodeID)
	if !ok {
		return nil // Should not happen with valid nodeID
	}
	
	if len(nodeDef.Next) == 0 {
		return nil
	}
	
	for _, nodeID := range nodeDef.Next {
		// 1. Check if we can activate this node (all prerequisites done)
		can, err := s.canActivate(ctx, userID, nodeID)
		if err != nil {
			log.Printf("[ActivateNextNodes] Error checking prerequisites for %s: %v", nodeID, err)
			continue
		}
		if !can {
			log.Printf("[ActivateNextNodes] Node %s prerequisites not yet met", nodeID)
			continue
		}

		// 2. Activate or Create
		inst, err := s.repo.GetNodeInstance(ctx, userID, nodeID)
		if inst != nil { // Exists
			if inst.State == "locked" {
				err = s.repo.UpdateNodeInstanceState(ctx, inst.ID, "locked", "active")
				if err == nil {
					log.Printf("Activated existing node %s", nodeID)
					_ = s.repo.UpsertJourneyState(ctx, userID, nodeID, "active", tenantID)
				}
			}
		} else if err == nil { 
			// Create
			id, err := s.repo.CreateNodeInstance(ctx, tenantID, userID, s.pb.VersionID, nodeID, "active", nil)
			if err != nil {
				log.Printf("[ActivateNextNodes] Error creating instance %s: %v", nodeID, err)
			} else {
				log.Printf("[ActivateNextNodes] Created new node instance %s for node %s", id, nodeID)
				_ = s.repo.UpsertJourneyState(ctx, userID, nodeID, "active", tenantID)
				
				// Log Event
				payload := map[string]any{"reason": "prerequisites_met", "source": completedNodeID}
				_ = s.repo.LogNodeEvent(ctx, id, "node_activated", userID, payload)
			}
		}
	}
	return nil
}

func (s *JourneyService) canActivate(ctx context.Context, userID, nodeID string) (bool, error) {
	nodeDef, ok := s.pb.NodeDefinition(nodeID)
	if !ok {
		return false, fmt.Errorf("node %s not found in playbook", nodeID)
	}

	if len(nodeDef.Prerequisites) == 0 {
		return true, nil
	}

	// Fetch all instances for user to check states
	// Optimization: Get states for specific nodes only? 
	// For now, journey state table is small and indexed by (user_id, node_id).
	// But we need to be sure they are 'done'.
	for _, preID := range nodeDef.Prerequisites {
		inst, err := s.repo.GetNodeInstance(ctx, userID, preID)
		if err != nil {
			return false, err
		}
		if inst == nil || inst.State != "done" {
			return false, nil
		}
	}

	return true, nil
}

func (s *JourneyService) verifyRequirements(ctx context.Context, inst *models.NodeInstance) error {
	// Use GetFullSubmissionSlots as it includes both SlotKey and Attachments
	slots, err := s.repo.GetFullSubmissionSlots(ctx, inst.ID)
	if err != nil {
		return err
	}

	for _, slot := range slots {
		if slot.Required {
			hasActive := false
			for _, a := range slot.Attachments {
				if a.IsActive {
					hasActive = true
					break
				}
			}
			if !hasActive {
				return fmt.Errorf("required file for slot '%s' is missing", slot.SlotKey)
			}
		}
	}

	return nil
}

// GetSubmission logic
func (s *JourneyService) GetSubmission(ctx context.Context, tenantID, userID, nodeID string, locale *string) (map[string]any, error) {
	// 1. Ensure Instance
	inst, err := s.EnsureNodeInstance(ctx, tenantID, userID, nodeID, locale)
	if err != nil {
		return nil, err
	}
	
	localeStr := ""
	if inst.Locale != nil {
		localeStr = *inst.Locale
	}

	dto := map[string]any{
		"node_id":             inst.NodeID,
		"playbook_version_id": inst.PlaybookVersionID,
		"state":               inst.State,
		"locale":              localeStr,
		"slots":               []any{}, // default
	}

	// Form Data & App7 Logic
	if inst.CurrentRev > 0 {
		revData, err := s.repo.GetFormRevision(ctx, inst.ID, inst.CurrentRev)
		if err == nil {
			// Default form DTO
			formDTO := map[string]any{
				"rev": inst.CurrentRev,
				"data": json.RawMessage(revData),
			}
			
			// Special handling for publications list
			if inst.NodeID == "S1_publications_list" {
				if form, parseErr := buildApp7Form(revData); parseErr == nil {
					summary := summarizeApp7(form.Sections)
					
					// Logic from original handler: check legacy counts?
					// Yes, merge if needed.
					for key, val := range form.LegacyCounts {
						if val > 0 && summary[key] == 0 {
							summary[key] = val
						}
					}
					
					clientData := map[string]any{
						"wos_scopus":  form.Sections.WosScopus,
						"kokson":      form.Sections.Kokson,
						"conferences": form.Sections.Conferences,
						"ip":          form.Sections.IP,
						"summary":     summary,
					}
					if len(form.LegacyCounts) > 0 {
						clientData["legacy_counts"] = form.LegacyCounts
					}
					formDTO["data"] = clientData
				}
			}
			dto["form"] = formDTO
		} else {
			log.Printf("[JourneyService] Error fetching revision: %v", err)
		}
	}

	// 2. Fetch Slots
	slots, err := s.repo.GetFullSubmissionSlots(ctx, inst.ID)
	if err != nil {
		return nil, err
	}
	dto["slots"] = slots

	// 3. Outcomes
	outs, err := s.repo.GetNodeOutcomes(ctx, inst.ID)
	if err != nil {
		return nil, err
	}
	if len(outs) > 0 {
		dto["outcomes"] = outs
	}
	
	return dto, nil
}

// EnsureNodeInstance finds or creates node instance
func (s *JourneyService) EnsureNodeInstance(ctx context.Context, tenantID, userID, nodeID string, locale *string) (*models.NodeInstance, error) {
	// 1. Try Load
	inst, err := s.repo.GetNodeInstance(ctx, userID, nodeID)
	if err != nil {
		return nil, err
	}
	
	if inst != nil {
		// Synch slots
		if err := s.ensureSlots(ctx, tenantID, inst.ID, nodeID); err != nil {
			return nil, err
		}
		return inst, nil
	}
	
	// 2. Create
	nodeDef, ok := s.pb.NodeDefinition(nodeID)
	if !ok {
		return nil, fmt.Errorf("node not found in playbook")
	}
	
	log.Printf("[JourneyService] Creating node instance: userID=%s nodeID=%s", userID, nodeID)
	id, err := s.repo.CreateNodeInstance(ctx, tenantID, userID, s.pb.VersionID, nodeID, "active", locale)
	if err != nil {
		return nil, err
	}
	
	// Create Slots
	if nodeDef.Requirements != nil {
		for _, up := range nodeDef.Requirements.Uploads {
			_, err := s.repo.CreateSlot(ctx, id, up.Key, tenantID, up.Required, "single", up.Mime)
			if err != nil {
				return nil, err
			}
		}
	}
	
	// Log Event
	_ = s.repo.LogNodeEvent(ctx, id, "opened", userID, map[string]any{"locale": locale})
	
	// Upsert Journey State
	_ = s.repo.UpsertJourneyState(ctx, userID, nodeID, "active", tenantID)
	
	// Return full object
	return s.repo.GetNodeInstanceByID(ctx, id)
}

func (s *JourneyService) ensureSlots(ctx context.Context, tenantID, instanceID, nodeID string) error {
	nodeDef, ok := s.pb.NodeDefinition(nodeID)
	if !ok || nodeDef.Requirements == nil || len(nodeDef.Requirements.Uploads) == 0 {
		return nil
	}
	
	existing, err := s.repo.GetNodeInstanceSlots(ctx, instanceID)
	if err != nil {
		return err
	}
	
	present := make(map[string]bool)
	for _, slot := range existing {
		present[slot.SlotKey] = true
	}
	
	for _, up := range nodeDef.Requirements.Uploads {
		if !present[up.Key] {
			_, err := s.repo.CreateSlot(ctx, instanceID, up.Key, tenantID, up.Required, "single", up.Mime)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// PutSubmission handles form data update and state transition
func (s *JourneyService) PutSubmission(ctx context.Context, tenantID, userID, role, nodeID string, locale *string, state string, formData []byte) error {
	// 1. Ensure Instance
	inst, err := s.EnsureNodeInstance(ctx, tenantID, userID, nodeID, locale)
	if err != nil {
		return err
	}

	// 2. Append Form Revision
	if len(formData) > 0 {
		inst.CurrentRev++
		err = s.repo.InsertFormRevision(ctx, inst.ID, inst.CurrentRev, formData, userID)
		if err != nil {
			return err
		}
		// Update instance rev/locale
		err = s.repo.UpsertSubmission(ctx, inst.ID, inst.CurrentRev, locale)
		if err != nil {
			return err
		}
		
		// Special handling for S1_profile node: sync data to users table
		if nodeID == "S1_profile" {
			err = s.syncProfileToUsers(ctx, tenantID, userID, formData)
			if err != nil {
				// Log error but don't fail the submission
				log.Printf("[PutSubmission] Failed to sync profile to users: %v", err)
			}
		}
	}

	// 3. Transition State if requested
	if state != "" && state != inst.State {
		err = s.transitionState(ctx, tenantID, inst, userID, role, state)
		if err != nil {
			return err
		}
	}
	
	// If state became done, activate next
	if state == "done" || (state == "" && inst.State == "done") {
		// Reload state in case transition happened
		refreshed, _ := s.repo.GetNodeInstanceByID(ctx, inst.ID)
		if refreshed != nil && refreshed.State == "done" {
			_ = s.ActivateNextNodes(ctx, userID, nodeID, tenantID)
		}
	}
	return nil
}

// PatchState handles state transition only
func (s *JourneyService) PatchState(ctx context.Context, tenantID, userID, role, nodeID, state string) error {
	inst, err := s.EnsureNodeInstance(ctx, tenantID, userID, nodeID, nil) // use existing locale
	if err != nil {
		return err
	}
	
	if state != "" && state != inst.State {
		err = s.transitionState(ctx, tenantID, inst, userID, role, state)
		if err != nil {
			return err
		}
		
		if state == "done" {
			_ = s.ActivateNextNodes(ctx, userID, nodeID, tenantID)
		}
	}
	return nil
}

// transitionState validation and execution
func (s *JourneyService) transitionState(ctx context.Context, tenantID string, inst *models.NodeInstance, userID, role, newState string) error {
	// Check allowed roles
	roles, err := s.repo.GetAllowedTransitionRoles(ctx, inst.State, newState)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	
	allowed := false
	// Logic from handler: "Iteration override: allow students to complete nodes directly to done"
	if role == "student" && newState == "done" && (inst.State == "active" || inst.State == "submitted") {
		allowed = true
	} else {
		for _, r := range roles {
			if r == role {
				allowed = true; break
			}
		}
	}
	
	if !allowed {
		return fmt.Errorf("role %s cannot transition from %s to %s", role, inst.State, newState)
	}
	
	// Requirement Verification for terminal states
	if newState == "submitted" || newState == "done" {
		if err := s.verifyRequirements(ctx, inst); err != nil {
			return fmt.Errorf("requirements not met: %w", err)
		}
	}

	oldState := inst.State
	err = s.repo.UpdateNodeInstanceState(ctx, inst.ID, oldState, newState)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("node state changed by another process (anticipated %s)", oldState)
		}
		return err
	}
	
	// Update Journey State
	_ = s.repo.UpsertJourneyState(ctx, userID, inst.NodeID, newState, tenantID)
	
	// Log Event
	payload := map[string]any{"from": oldState, "to": newState}
	_ = s.repo.LogNodeEvent(ctx, inst.ID, "state_changed", userID, payload)
	
	// Notify
	go s.sendStateChangeEmail(context.Background(), userID, inst.NodeID, oldState, newState)
	
	return nil
}

func (s *JourneyService) sendStateChangeEmail(ctx context.Context, userID, nodeID, fromState, toState string) {
	// Fetch user details? Repo method exists? 
	// GetUsersByIDs
	users, err := s.repo.GetUsersByIDs(context.Background(), []string{userID})
	if err != nil || len(users) == 0 {
		return
	}
	user := users[0]
	// Logic:
	studentName := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	
	if s.mailer != nil {
		subject := fmt.Sprintf("Node Status Change: %s (%s -> %s)", nodeID, fromState, toState)
		body := fmt.Sprintf("Student %s changed node %s status from %s to %s.", studentName, nodeID, fromState, toState)
		log.Printf("[JourneyService] Sending email: %s", subject)
		_ = s.mailer.SendNotificationEmail("admin@portal.kaznmu.kz", subject, body) 
	}
}

// PresignUpload generates S3 URL and returns both URL and object path
func (s *JourneyService) PresignUpload(ctx context.Context, userID, nodeID, slotKey, filename, contentType string, sizeBytes int64) (string, string, error) {
	// 1. Validate against playbook
	nodeDef, ok := s.pb.NodeDefinition(nodeID)
	if !ok {
		return "", "", errors.New("node not found in playbook")
	}

	var slotDef *playbook.UploadRequirement
	if nodeDef.Requirements != nil {
		for _, up := range nodeDef.Requirements.Uploads {
			if up.Key == slotKey {
				slotDef = &up
				break
			}
		}
	}

	if slotDef == nil {
		return "", "", errors.New("slot not found in node definition")
	}

	// MIME check
	if len(slotDef.Mime) > 0 {
		allowed := false
		for _, m := range slotDef.Mime {
			if m == contentType {
				allowed = true
				break
			}
		}
		if !allowed {
			return "", "", fmt.Errorf("mime type %s not allowed for this slot", contentType)
		}
	}

	// Size check
	maxBytes := int64(s.cfg.FileUploadMaxMB) * 1024 * 1024
	if sizeBytes > maxBytes {
		return "", "", fmt.Errorf("file size %d bytes is too large (max %dMB)", sizeBytes, s.cfg.FileUploadMaxMB)
	}

	// 2. Logic
	path := fmt.Sprintf("node_uploads/%s/%s/%s/%s", nodeID, slotKey, uuid.NewString(), filename)

	if s.storage == nil {
		return "", "", errors.New("storage client not available")
	}

	url, err := s.storage.PresignPut(ctx, path, contentType, 15*time.Minute)
	if err != nil {
		return "", "", err
	}

	return url, path, nil
}

// AttachUpload logic
func (s *JourneyService) AttachUpload(ctx context.Context, tenantID, userID, nodeID, slotKey, objectKey, filename string, sizeBytes int64) error {
	inst, err := s.EnsureNodeInstance(ctx, tenantID, userID, nodeID, nil)
	if err != nil {
		return err
	}

	slot, err := s.repo.GetSlot(ctx, inst.ID, slotKey)
	if err != nil {
		return err
	}

	if s.docSvc == nil {
		return errors.New("document service not available")
	}

	log.Printf("[JourneyService] AttachUpload: Attaching %s to slot %s", filename, slot.ID)

	// Multiplicity check: If single, deactivate previous attachments
	if slot.Multiplicity == "single" {
		err = s.repo.DeactivateSlotAttachments(ctx, slot.ID)
		if err != nil {
			return fmt.Errorf("failed to deactivate old attachments: %w", err)
		}
	}

	// 1. Create or Find Document
	// For simplicity, we create a new Document entry for each upload, 
	// or we could look up an existing one for this slot.
	// Standard approach: create a Document record if this is the first upload to the slot, 
	// but many slots allow multiple versions.
	
	docID, err := s.docSvc.CreateMetadata(ctx, CreateDocumentRequest{
		Title:    filename,
		Kind:     "node_slot",
		TenantID: tenantID,
		UserID:   userID,
	})
	if err != nil {
		return fmt.Errorf("failed to create document metadata: %w", err)
	}

	// 2. Create Document Version
	verID, err := s.docSvc.CreateVersion(ctx, docID, tenantID, userID, models.DocumentVersion{
		StoragePath: objectKey, // We use the S3 key here
		MimeType:    "application/octet-stream", // Should we pass this from handler?
		SizeBytes:   sizeBytes,
	})
	if err != nil {
		return fmt.Errorf("failed to create document version: %w", err)
	}

	// 3. Create Node Instance Slot Attachment
	_, err = s.repo.CreateAttachment(ctx, slot.ID, verID, "submitted", filename, userID, sizeBytes)
	if err != nil {
		return fmt.Errorf("failed to create slot attachment: %w", err)
	}

	return nil
}

// syncProfileToUsers syncs profile submission data to the users table
// This is called when the S1_profile node is submitted
func (s *JourneyService) syncProfileToUsers(ctx context.Context, tenantID, userID string, formData []byte) error {
	// Parse the form data
	var data map[string]interface{}
	if err := json.Unmarshal(formData, &data); err != nil {
		return fmt.Errorf("failed to parse form data: %w", err)
	}
	
	// Build update fields dynamically based on available fields
	fields := make(map[string]interface{})
	
	// Map profile fields to user table columns
	syncFields := []string{"program", "specialty", "department", "cohort"}
	for _, f := range syncFields {
		if val, ok := data[f]; ok {
			fields[f] = val
		}
	}
	
	// If no fields to update, return early
	if len(fields) == 0 {
		return nil
	}
	
	// Execute the update via repository
	err := s.repo.SyncProfileToUsers(ctx, userID, tenantID, fields)
	if err != nil {
		return fmt.Errorf("failed to update users table: %w", err)
	}
	
	log.Printf("[syncProfileToUsers] Successfully synced profile data for user %s", userID)
	return nil
}

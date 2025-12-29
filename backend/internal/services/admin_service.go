package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	pb "github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
)

type AdminService struct {
	repo    repository.AdminRepository
	pb      *pb.Manager
	cfg     config.AppConfig
	storage StorageClient
}

func NewAdminService(repo repository.AdminRepository, pbm *pb.Manager, cfg config.AppConfig, storage StorageClient) *AdminService {
	return &AdminService{
		repo:    repo,
		pb:      pbm,
		cfg:     cfg,
		storage: storage,
	}
}

func (s *AdminService) ListStudentProgress(ctx context.Context, tenantID string) ([]models.StudentProgressSummary, error) {
	summaries, err := s.repo.ListStudentProgress(ctx, tenantID, s.pb.VersionID)
	if err != nil {
		return nil, err
	}
	
	totalNodes := len(s.pb.Nodes)
	for i := range summaries {
		summaries[i].TotalNodes = totalNodes
		if totalNodes > 0 {
			summaries[i].Percent = float64(summaries[i].CompletedNodes) * 100.0 / float64(totalNodes)
		}
	}
	return summaries, nil
}

func (s *AdminService) MonitorStudents(ctx context.Context, filter models.FilterParams) ([]models.StudentMonitorRow, error) {
	// 1. Fetch base filtered list (paginated)
	rows, err := s.repo.ListStudentsForMonitor(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return []models.StudentMonitorRow{}, nil
	}

	ids := make([]string, len(rows))
	for i, r := range rows {
		ids[i] = r.ID
	}

	// 2. Batch fetch related data
	advisors, _ := s.repo.GetAdvisorsForStudents(ctx, ids)
	doneCounts, _ := s.repo.GetDoneCountsForStudents(ctx, ids)
	lastUpdates, _ := s.repo.GetLastUpdatesForStudents(ctx, ids)
	rpRequired, _ := s.repo.GetRPRequiredForStudents(ctx, ids)

	// DEBUG: Log playbook info
	fmt.Printf("[MonitorStudents] PlaybookVersionID=%s, TotalNodes=%d, StudentCount=%d\n", s.pb.VersionID, len(s.pb.Nodes), len(rows))
	fmt.Printf("[MonitorStudents] DoneCounts map: %+v\n", doneCounts)

	// 3. Merge and Compute
	totalNodes := len(s.pb.Nodes)
	_, worldNodes := s.getWorlds()
	w3Count := len(worldNodes["W3"])

	enriched := make([]models.StudentMonitorRow, 0, len(rows))
	
	for _, r := range rows {
		// Populate
		r.Advisors = advisors[r.ID]
		if r.Advisors == nil {
			r.Advisors = []models.AdvisorSummary{}
		}
		r.RPRequired = rpRequired[r.ID]
		
		// RP Logic
		if filter.RPRequired && !r.RPRequired {
			continue // Should have been filtered in repo if possible, but RP logic is complex JSON extraction
			// Ideally we move RP filtering to Repo SQL using JSON operators.
			// for now, strict filter here might reduce page size, which is a trade-off.
		}

		r.DoneCount = doneCounts[r.ID]
		if t, ok := lastUpdates[r.ID]; ok {
			r.LastUpdate = &t
		}

		// Calc Percent
		r.TotalNodes = totalNodes
		totalRequired := totalNodes
		if !r.RPRequired {
			totalRequired = totalNodes - w3Count
		}
		if totalRequired <= 0 { 
			totalRequired = totalNodes 
		}
		r.OverallProgressPct = 0.0
		if totalRequired > 0 {
			r.OverallProgressPct = float64(r.DoneCount) * 100.0 / float64(totalRequired)
		}

		// Stage logic - derive from CurrentNodeID
		nodeID := ""
		if r.CurrentNodeID != nil {
			nodeID = *r.CurrentNodeID
		}
		if r.CurrentNodeID != nil && *r.CurrentNodeID != "" {
			r.CurrentStage = s.pb.NodeWorldID(*r.CurrentNodeID)
		}
		if r.CurrentStage == "" {
			r.CurrentStage = "W1" // Default to first stage
		}
		// DEBUG: Log stage calculation
		fmt.Printf("[MonitorStudents] Student %s: CurrentNodeID=%s -> CurrentStage=%s, DoneCount=%d\n", r.ID, nodeID, r.CurrentStage, r.DoneCount)
		
		enriched = append(enriched, r)

	}

	return enriched, nil
}

// Helpers
func (s *AdminService) getWorlds() ([]string, map[string][]string) {
	// Simplified logic to extract world nodes. Using playbook manager if available.
	if s.pb == nil {
		return []string{}, map[string][]string{}
	}
	// Note: keys in NodeWorlds are NodeIDs, values are WorldIDs.
	// We want map[WorldID][]NodeID
	out := make(map[string][]string)
	// Hardcoded worlds order? Or pb.Raw parsing.
	// For simplicty:
	out["W1"] = s.pb.GetNodesByWorld("W1")
	out["W2"] = s.pb.GetNodesByWorld("W2")
	out["W3"] = s.pb.GetNodesByWorld("W3")
	return []string{"W1", "W2", "W3"}, out
}

func (s *AdminService) MonitorAnalytics(ctx context.Context, filter models.FilterParams) (*models.AdminAnalytics, error) {
	// 1. Get ALL matching student IDs (no limit)
	filter.Limit = 0
	rows, err := s.repo.ListStudentsForMonitor(ctx, filter)
	if err != nil {
		return nil, err
	}
	
	ids := make([]string, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.ID)
	}

	res := &models.AdminAnalytics{}
	if len(ids) == 0 {
		return res, nil
	}
	
	// 2. RP Required Count
	// Reuse batch loader.
	// Optimization: This loads filtering logic in code.
	// If RPRequired filter was ON, all rows are RPRequired.
	if filter.RPRequired {
		res.RPRequiredCount = len(ids)
	} else {
		rpMap, _ := s.repo.GetRPRequiredForStudents(ctx, ids)
		count := 0
		for _, v := range rpMap {
			if v {
				count++
			}
		}
		// Also filter IDs if rpOnly logic requires strict subset? 
		// Previous handler logic: If rpOnly=true, it filtered IDs before calc other stats.
		// My `repo.ListStudentsForMonitor` handles DB filters, but `RPRequired` is currently post-filter in handler/service `MonitorStudents`.
		// If `filter.RPRequired` is passed to repo but not restricted in SQL, `rows` contains everyone.
		// Let's check `ListStudentsForMonitor`: it does NOT filter by RPRequired in SQL.
		// So `rows` includes non-RP even if `filter.RPRequired=true`.
		// I must filter `ids` here.
		
		finalIDs := make([]string, 0, len(ids))
		for _, id := range ids {
			if rpMap[id] {
				finalIDs = append(finalIDs, id)
			} else if !filter.RPRequired {
				finalIDs = append(finalIDs, id)
			}
		}
		ids = finalIDs
		res.RPRequiredCount = count // Count of ACTUAL RP required in the original set
		// Wait, if `filter.RPRequired=true`, we should only consider RP IDs for other stats.
		if filter.RPRequired {
			res.RPRequiredCount = len(ids) // All remaining are RP
		}
	}
	
	if len(ids) == 0 {
		return res, nil
	}

	// 3. Antiplag
	if apCount, err := s.repo.GetAntiplagCount(ctx, ids, s.pb.VersionID); err == nil {
		res.AntiplagDonePercent = float64(apCount) * 100.0 / float64(len(ids))
	}
	
	// 4. Bottleneck
	startOfMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())
	if nid, cnt, err := s.repo.GetBottleneck(ctx, ids, s.pb.VersionID, startOfMonth); err == nil {
		res.BottleneckNodeID = nid
		res.BottleneckCount = cnt
	}

	// 5. W2 Median
	_, worlds := s.getWorlds()
	w2Nodes := worlds["W2"]
	if len(w2Nodes) > 0 {
		if durs, err := s.repo.GetW2Durations(ctx, ids, s.pb.VersionID, w2Nodes); err == nil && len(durs) > 0 {
			// Sort and median
			// Bubble sort for small N or standard sort
			for i := 1; i < len(durs); i++ {
				key := durs[i]
				j := i - 1
				for j >= 0 && durs[j] > key {
					durs[j+1] = durs[j]
					j--
				}
				durs[j+1] = key
			}
			mid := len(durs) / 2
			if len(durs)%2 == 1 {
				res.W2MedianDays = durs[mid]
			} else {
				res.W2MedianDays = (durs[mid-1] + durs[mid]) / 2
			}
		}
	}
	
	return res, nil
}

func (s *AdminService) GetStudentDetails(ctx context.Context, studentID, tenantID string) (*models.StudentDetails, error) {
	// 1. Fetch basic info
	details, err := s.repo.GetStudentDetails(ctx, studentID, tenantID)
	if err != nil {
		return nil, err
	}
	
	// 2. Advisors
	advMap, _ := s.repo.GetAdvisorsForStudents(ctx, []string{studentID})
	if list, ok := advMap[studentID]; ok {
		details.Advisors = list
	} else {
		details.Advisors = []models.AdvisorSummary{}
	}
	
	// 3. RP
	rpMap, _ := s.repo.GetRPRequiredForStudents(ctx, []string{studentID})
	details.RPRequired = rpMap[studentID]
	
	// 4. Progress
	// Reuse "GetStudentNodeInstances" or counts?
	// Handler used explicit queries for 'done' count, 'last node', 'stage done'.
	// We can compute ALL from `GetStudentNodeInstances` (all nodes for student).
	instances, _ := s.repo.GetStudentNodeInstances(ctx, studentID) // Returns map of latest instances
	
	// Computed logic
	totalNodes := len(s.pb.Nodes)
	details.TotalNodes = totalNodes
	_, worldNodes := s.getWorlds()
	w3Count := len(worldNodes["W3"])
	
	totalRequired := totalNodes
	if !details.RPRequired {
		totalRequired = totalNodes - w3Count
	}
	if totalRequired <= 0 { totalRequired = totalNodes }
	
	doneCount := 0
	var lastNodeID string
	var lastUpdate time.Time
	
	// DEBUG: Log instances info
	fmt.Printf("[GetStudentDetails] studentID=%s, PlaybookVersionID=%s, InstancesCount=%d\n", studentID, s.pb.VersionID, len(instances))
	for _, inst := range instances {
		fmt.Printf("[GetStudentDetails]   Instance: NodeID=%s, State=%s, Version=%s\n", inst.NodeID, inst.State, inst.PlaybookVersionID)
	}
	
	// Count done and find last update
	for _, inst := range instances {
		// Filter by version ??? Handler uses specific version stats?
		// Handler `GetStudentDetails` uses `h.pb.VersionID` for 'done' count.
		// `GetStudentNodeInstances` returns distinct on node_id (latest).
		// We should check if latest instance matches version?
		// Or just count state='done' and version match?
		// Handler line 585: `WHERE user_id=$1 AND playbook_version_id=$2 AND state='done'`
		// So stats specific to active version.
		if inst.PlaybookVersionID == s.pb.VersionID {
			if inst.State == "done" {
				doneCount++
			}
			// Determining last update
			if inst.UpdatedAt.After(lastUpdate) {
				lastUpdate = inst.UpdatedAt
				lastNodeID = inst.NodeID
			}
		}
	}
	
	// DEBUG: Log calculated values
	fmt.Printf("[GetStudentDetails] doneCount=%d, lastNodeID=%s, totalRequired=%d\n", doneCount, lastNodeID, totalRequired)
	
	details.OverallProgressPct = 0.0
	if totalRequired > 0 {
		details.OverallProgressPct = float64(doneCount) * 100.0 / float64(totalRequired)
	}
	
	// Stage logic
	stage := s.pb.NodeWorldID(lastNodeID)
	fmt.Printf("[GetStudentDetails] NodeWorldID(%s)=%s\n", lastNodeID, stage)
	if stage == "" { stage = "W1" }
	details.CurrentStage = stage
	details.StageTotal = len(worldNodes[stage])

	
	// Stage done (only for active version)
	stageDone := 0
	stageNodeSet := make(map[string]bool)
	for _, nid := range worldNodes[stage] {
		stageNodeSet[nid] = true
	}
	for _, inst := range instances {
		if inst.PlaybookVersionID == s.pb.VersionID && inst.State == "done" && stageNodeSet[inst.NodeID] {
			stageDone++
		}
	}
	details.StageDone = stageDone
	
	// Last Update Global
	// Handler uses `GREATEST` query again (line 612 is simpler in handler "MAX(updated_at)" on node_instances).
	// But `MonitorStudents` used GREATEST of 4 tables.
	// `GetStudentDetails` handler line 612 uses ONLY node_instances.
	// Let's stick to node_instance update for consistency with handler `GetStudentDetails`.
	// Wait, line 612 handler: `SELECT MAX(updated_at) ... node_instances`.
	
	// HOWEVER, getting true last activity is better.
	// Repository `GetLastUpdatesForStudents` uses the GREATEST logic.
	// Let's use that for accuracy.
	lastUpdatesMap, _ := s.repo.GetLastUpdatesForStudents(ctx, []string{studentID})
	if t, ok := lastUpdatesMap[studentID]; ok {
		ts := t.Format(time.RFC3339)
		details.LastUpdate = &ts
	}

	return details, nil
}

// GetStudentJourney returns the journey nodes state including attachments
func (s *AdminService) GetStudentJourney(ctx context.Context, studentID, role, callerID string) ([]models.StudentJourneyNode, error) {
	// RBAC: Check access if advisor
	if role == "advisor" {
		allowed, err := s.repo.CheckAdvisorAccess(ctx, studentID, callerID)
		if err != nil {
			return nil, err
		}
		if !allowed {
			return nil, errors.New("forbidden")
		}
	}
	
	return s.repo.GetStudentJourneyNodes(ctx, studentID)
}

// ListStudentNodeFiles returns files for a specific node
func (s *AdminService) ListStudentNodeFiles(ctx context.Context, studentID, nodeID, role, callerID string) ([]models.NodeFile, error) {
	// RBAC
	if role == "advisor" {
		allowed, err := s.repo.CheckAdvisorAccess(ctx, studentID, callerID)
		if err != nil {
			return nil, err
		}
		if !allowed {
			return nil, errors.New("forbidden")
		}
	}
	
	return s.repo.GetNodeFiles(ctx, studentID, nodeID)
}

// ReviewResult contains result of review
type ReviewResult struct {
	Status     string
	State      string
	ReviewNote *string
	ApprovedAt *string
	StudentID  string
	NodeID     string
}

// ReviewAttachment handles approval/rejection of docs
func (s *AdminService) ReviewAttachment(ctx context.Context, attachmentID, status, note, actorID, role, tenantID string) (*ReviewResult, error) {
	// Validate input
	// Get Meta
	meta, err := s.repo.GetAttachmentMeta(ctx, attachmentID)
	if err != nil {
		return nil, err // NotFound or DB error
	}
	
	// RBAC
	if role == "advisor" {
		allowed, err := s.repo.CheckAdvisorAccess(ctx, meta.StudentID, actorID)
		if err != nil {
			return nil, err
		}
		if !allowed {
			return nil, errors.New("forbidden")
		}
	}
	
	// Update Attachment
	err = s.repo.UpdateAttachmentStatus(ctx, attachmentID, status, note, actorID)
	if err != nil {
		return nil, err
	}
	
	// Log Event
	payload := map[string]any{"attachment_id": attachmentID, "status": status}
	if note != "" { payload["note"] = note }
	_ = s.repo.LogNodeEvent(ctx, meta.InstanceID, "attachment_reviewed", actorID, payload)
	
	// Node State Logic
	// Check latest file status (the one that matters for progress)
	latestStatus, err := s.repo.GetLatestAttachmentStatus(ctx, meta.InstanceID)
	if err != nil {
		return nil, err
	}
	
	// Get all counts if needed, but latest drives state mostly
	// Logic: If latest is approved -> done. If rejected -> needs_fixes.
	// Logic from Handler:
	// counts: submitted, approved, rejected
	submitted, approved, rejected, _ := s.repo.GetAttachmentCounts(ctx, meta.InstanceID)
	total := submitted + approved + rejected
	
	newState := meta.State
	if total > 0 {
		switch {
		case latestStatus == "approved":
			newState = "done"
		case latestStatus == "approved_with_comments":
			newState = "done"
		case latestStatus == "rejected":
			newState = "needs_fixes"
		case latestStatus == "submitted":
			newState = "under_review"
		default:
			newState = "submitted"
		}
	}
	
	if newState != meta.State {
		// Update Instance
		if err := s.repo.UpdateNodeInstanceState(ctx, meta.InstanceID, newState); err != nil {
			return nil, err
		}
		// Update All Instances (sync versions)
		if err := s.repo.UpdateAllNodeInstances(ctx, meta.StudentID, meta.NodeID, meta.InstanceID, newState); err != nil {
			// Warning
		}
		// Upsert Journey State
		_ = s.repo.UpsertJourneyState(ctx, meta.TenantID, meta.StudentID, meta.NodeID, newState)
		
		_ = s.repo.LogNodeEvent(ctx, meta.InstanceID, "state_changed", actorID, map[string]any{"from": meta.State, "to": newState})
		
		// Activate Next Nodes
		if newState == "done" {
			// TODO: Implement ActiveNextNodes logic in Service or Helper, or let Handler do it
		}
	}
	
	// Notifications
	title := "Document Reviewed: " + meta.Filename
	msg := "Your document has been reviewed."
	if status == "approved" || status == "approved_with_comments" {
		msg = "Your document has been approved."
	} else if status == "rejected" {
		msg = "Changes requested for your document."
		if note != "" {
			msg += " Note: " + note
		}
	}
	// Create notification
	if meta.TenantID != "" {
		_ = s.repo.CreateNotification(ctx, meta.StudentID, title, msg, "/journey", "document_review", meta.TenantID)
	}

	// Return result
	now := time.Now().Format(time.RFC3339)
	var rn *string
	if note != "" { rn = &note }
	return &ReviewResult{
		Status: status,
		State: newState,
		ReviewNote: rn,
		ApprovedAt: &now,
		StudentID: meta.StudentID,
		NodeID: meta.NodeID,
	}, nil
}

// UploadReviewedDocument
func (s *AdminService) UploadReviewedDocument(ctx context.Context, attachmentID, versionID, actorID, role string) (string, error) {
	meta, err := s.repo.GetAttachmentMeta(ctx, attachmentID)
	if err != nil { return "", err }
	
	if role == "advisor" {
		allowed, err := s.repo.CheckAdvisorAccess(ctx, meta.StudentID, actorID)
		if err != nil { return "", err }
		if !allowed { return "", errors.New("forbidden") }
	}
	
	err = s.repo.UploadReviewedDocument(ctx, attachmentID, versionID, actorID)
	if err != nil { return "", err }
	
	_ = s.repo.LogNodeEvent(ctx, meta.InstanceID, "reviewed_document_uploaded", actorID, map[string]any{
		"attachment_id": attachmentID, "reviewed_version_id": versionID,
	})
	
	return time.Now().Format(time.RFC3339), nil
}

func (s *AdminService) CreateReminders(ctx context.Context, studentIDs []string, title, message string, dueAt *string, callerID string) error {
	return s.repo.CreateReminders(ctx, studentIDs, title, message, dueAt, callerID)
}

func (s *AdminService) AttachReviewedDocument(ctx context.Context, attachmentID string, 
    storagePath, objKey, bucket, mimeType string, sizeBytes int64, etag string,
    actorID, role, tenantID string) (string, string, error) {
    
    meta, err := s.repo.GetAttachmentMeta(ctx, attachmentID)
    if err != nil { return "", "", err }
    
	if role == "advisor" {
		allowed, err := s.repo.CheckAdvisorAccess(ctx, meta.StudentID, actorID)
		if err != nil { return "", "", err }
		if !allowed { return "", "", errors.New("forbidden") }
	}
    
    versionID, err := s.repo.CreateReviewedDocumentVersion(ctx, meta.DocumentID, storagePath, objKey, bucket, mimeType, sizeBytes, actorID, etag, tenantID)
    if err != nil { return "", "", err }
    
    err = s.repo.UploadReviewedDocument(ctx, attachmentID, versionID, actorID)
    if err != nil { return "", "", err }
    
    payload := map[string]any{
		"attachment_id":       attachmentID,
		"reviewed_version_id": versionID,
		"filename":            meta.Filename,
	}
    _ = s.repo.LogNodeEvent(ctx, meta.InstanceID, "reviewed_document_uploaded", actorID, payload)
    
    return versionID, time.Now().Format(time.RFC3339), nil
}

// PresignReviewedDocumentUpload generates a presigned URL for reviewed document upload
func (s *AdminService) PresignReviewedDocumentUpload(ctx context.Context, attachmentID string, 
	filename, contentType string, sizeBytes int64, actorID, role string) (string, string, error) {
	
	// Verify attachment exists
	meta, err := s.repo.GetAttachmentMeta(ctx, attachmentID)
	if err != nil {
		return "", "", err
	}

	// Permission Check
	if role == "advisor" {
		allowed, err := s.repo.CheckAdvisorAccess(ctx, meta.StudentID, actorID)
		if err != nil { return "", "", err }
		if !allowed { return "", "", errors.New("forbidden") }
	}

	if s.storage == nil {
		return "", "", errors.New("storage client not available")
	}

	// Generate object key: reviewed_documents/{attachment_id}/{timestamp}-{filename}
	timestamp := time.Now().Format("20060102-150405")
	objectKey := fmt.Sprintf("reviewed_documents/%s/%s-%s", attachmentID, timestamp, filename)

	expires := GetPresignExpires()
	url, err := s.storage.PresignPut(ctx, objectKey, contentType, expires)
	if err != nil {
		return "", "", err
	}

	return url, objectKey, nil
}

// Admin Notifications

func (s *AdminService) ListNotifications(ctx context.Context, unreadOnly bool) ([]models.AdminNotification, error) {
	return s.repo.ListAdminNotifications(ctx, unreadOnly)
}

func (s *AdminService) GetUnreadNotificationCount(ctx context.Context) (int, error) {
	return s.repo.GetAdminUnreadCount(ctx)
}

func (s *AdminService) MarkNotificationAsRead(ctx context.Context, id string) error {
	return s.repo.MarkAdminNotificationRead(ctx, id)
}

func (s *AdminService) MarkAllNotificationsAsRead(ctx context.Context) error {
	return s.repo.MarkAllAdminNotificationsRead(ctx)
}

package services_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestJourneyService_PresignUpload_Hardening(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	tenantID := "00000000-0000-0000-0000-000000000001"
	userID := "11111111-1111-1111-1111-111111111111"
	versionID := "22222222-2222-2222-2222-222222222222"

	// Setup background data
	_, _ = db.Exec(`INSERT INTO tenants (id, name, slug) VALUES ($1, 'Test', 'test') ON CONFLICT DO NOTHING`, tenantID)
	_, _ = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash) VALUES ($1, 'testuser', 'test@test.com', 'Test', 'User', 'student', 'hash') ON CONFLICT DO NOTHING`, userID)
	
	rawJSON := `{
		"worlds": [{
			"id": "W1",
			"nodes": [{
				"id": "upload_node",
				"requirements": {
					"uploads": [
						{
							"key": "pdf_only",
							"mime": ["application/pdf"],
							"required": true
						},
						{
							"key": "images",
							"mime": ["image/jpeg", "image/png"]
						}
					]
				}
			}]
		}]
	}`
	_, _ = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, tenant_id) VALUES ($1, 'v1', 'sum1', $2, $3) ON CONFLICT DO NOTHING`, versionID, rawJSON, tenantID)

	pb := &playbook.Manager{
		VersionID: versionID,
		Nodes: map[string]playbook.Node{
			"upload_node": {
				ID: "upload_node",
				Requirements: &playbook.Requirements{
					Uploads: []playbook.UploadRequirement{
						{Key: "pdf_only", Mime: []string{"application/pdf"}},
						{Key: "images", Mime: []string{"image/jpeg", "image/png"}},
					},
				},
			},
		},
	}

	cfg := config.AppConfig{FileUploadMaxMB: 5}
	repo := repository.NewSQLJourneyRepository(db)
	
	// Mock storage client
	storage := &mockStorage{}
	svc := services.NewJourneyService(repo, pb, cfg, nil, storage, nil)

	tests := []struct {
		name        string
		slotKey     string
		contentType string
		sizeBytes   int64
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "Valid PDF",
			slotKey:     "pdf_only",
			contentType: "application/pdf",
			sizeBytes:   1024,
			wantErr:     false,
		},
		{
			name:        "Invalid MIME for PDF slot",
			slotKey:     "pdf_only",
			contentType: "image/png",
			sizeBytes:   1024,
			wantErr:     true,
			errMsg:      "mime type image/png not allowed",
		},
		{
			name:        "Valid JPEG",
			slotKey:     "images",
			contentType: "image/jpeg",
			sizeBytes:   1024,
			wantErr:     false,
		},
		{
			name:        "Valid PNG",
			slotKey:     "images",
			contentType: "image/png",
			sizeBytes:   1024,
			wantErr:     false,
		},
		{
			name:        "Invalid slot key",
			slotKey:     "nonexistent",
			contentType: "application/pdf",
			sizeBytes:   1024,
			wantErr:     true,
			errMsg:      "slot not found",
		},
		{
			name:        "File too large",
			slotKey:     "pdf_only",
			contentType: "application/pdf",
			sizeBytes:   6 * 1024 * 1024, // 6MB > 5MB
			wantErr:     true,
			errMsg:      "too large",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := svc.PresignUpload(context.Background(), userID, "upload_node", tt.slotKey, "test.file", tt.contentType, tt.sizeBytes)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Empty(t, url)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, url)
			}
		})
	}
}

func TestJourneyService_AttachUpload_Multiplicity(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	tenantID := "00000000-0000-0000-0000-000000000001"
	userID := "11111111-1111-1111-1111-111111111111"
	versionID := "22222222-2222-2222-2222-222222222222"

	// Setup background data
	_, _ = db.Exec(`INSERT INTO tenants (id, name, slug) VALUES ($1, 'Test', 'test') ON CONFLICT DO NOTHING`, tenantID)
	_, _ = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash) VALUES ($1, 'testuser', 'test@test.com', 'Test', 'User', 'student', 'hash') ON CONFLICT DO NOTHING`, userID)
	
	rawJSON := `{"worlds": [{"id": "W1", "nodes": [{"id": "node1", "requirements": {"uploads": [{"key": "single_slot"}]}}]}]}`
	_, _ = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, tenant_id) VALUES ($1, 'v1', 'sum1', $2, $3) ON CONFLICT DO NOTHING`, versionID, rawJSON, tenantID)

	pb := &playbook.Manager{
		VersionID: versionID,
		Nodes: map[string]playbook.Node{
			"node1": {
				ID: "node1",
				Requirements: &playbook.Requirements{
					Uploads: []playbook.UploadRequirement{{Key: "single_slot"}},
				},
			},
		},
	}

	cfg := config.AppConfig{}
	repo := repository.NewSQLJourneyRepository(db)
	docRepo := repository.NewSQLDocumentRepository(db)
	storage := &mockStorage{}
	docSvc := services.NewDocumentService(docRepo, cfg, storage)
	svc := services.NewJourneyService(repo, pb, cfg, nil, storage, docSvc)

	ctx := context.Background()
	slotKey := "single_slot"

	// 1. Attach first file
	err := svc.AttachUpload(ctx, tenantID, userID, "node1", slotKey, "path/file1.pdf", "file1.pdf", 1024)
	assert.NoError(t, err)

	// Verify one active attachment
	var activeCount int
	err = db.Get(&activeCount, `SELECT COUNT(*) FROM node_instance_slot_attachments a 
		JOIN node_instance_slots s ON a.slot_id = s.id
		JOIN node_instances i ON s.node_instance_id = i.id
		WHERE i.user_id = $1 AND s.slot_key = $2 AND a.is_active = true`, userID, slotKey)
	assert.NoError(t, err)
	assert.Equal(t, 1, activeCount)

	// 2. Attach second file (should deactivate first because default multiplicity is "single")
	err = svc.AttachUpload(ctx, tenantID, userID, "node1", slotKey, "path/file2.pdf", "file2.pdf", 2048)
	assert.NoError(t, err)

	// Verify still only one active attachment
	err = db.Get(&activeCount, `SELECT COUNT(*) FROM node_instance_slot_attachments a 
		JOIN node_instance_slots s ON a.slot_id = s.id
		JOIN node_instances i ON s.node_instance_id = i.id
		WHERE i.user_id = $1 AND s.slot_key = $2 AND a.is_active = true`, userID, slotKey)
	assert.NoError(t, err)
	assert.Equal(t, 1, activeCount)

	// Verify total attachments is 2
	var totalCount int
	err = db.Get(&totalCount, `SELECT COUNT(*) FROM node_instance_slot_attachments a 
		JOIN node_instance_slots s ON a.slot_id = s.id
		JOIN node_instances i ON s.node_instance_id = i.id
		WHERE i.user_id = $1 AND s.slot_key = $2`, userID, slotKey)
	assert.NoError(t, err)
	assert.Equal(t, 2, totalCount)
}

func TestJourneyService_PresignUpload_StorageFailure(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	tenantID := "00000000-0000-0000-0000-000000000001"
	userID := "11111111-1111-1111-1111-111111111111"
	versionID := "22222222-2222-2222-2222-222222222222"

	// Setup background data
	_, _ = db.Exec(`INSERT INTO tenants (id, name, slug) VALUES ($1, 'Test', 'test') ON CONFLICT DO NOTHING`, tenantID)
	_, _ = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash) VALUES ($1, 'testuser', 'test@test.com', 'Test', 'User', 'student', 'hash') ON CONFLICT DO NOTHING`, userID)
	
	rawJSON := `{"worlds": [{"id": "W1", "nodes": [{"id": "node1", "requirements": {"uploads": [{"key": "file"}]}}]}]}`
	_, _ = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, tenant_id) VALUES ($1, 'v1', 'sum1', $2, $3) ON CONFLICT DO NOTHING`, versionID, rawJSON, tenantID)

	pb := &playbook.Manager{
		VersionID: versionID,
		Nodes: map[string]playbook.Node{
			"node1": {
				ID: "node1",
				Requirements: &playbook.Requirements{
					Uploads: []playbook.UploadRequirement{{Key: "file"}},
				},
			},
		},
	}

	cfg := config.AppConfig{FileUploadMaxMB: 5}
	repo := repository.NewSQLJourneyRepository(db)
	storage := &errorStorage{}
	svc := services.NewJourneyService(repo, pb, cfg, nil, storage, nil)

	_, err := svc.PresignUpload(context.Background(), userID, "node1", "file", "test.pdf", "application/pdf", 1024)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "storage error")
}

func TestJourneyService_AttachUpload_RepoFailure(t *testing.T) {
	pb := &playbook.Manager{
		Nodes: map[string]playbook.Node{
			"node1": {ID: "node1"},
		},
	}
	mock := &MockJourneyRepository{
		GetNodeInstanceFunc: func(ctx context.Context, userID, nodeID string) (*models.NodeInstance, error) {
			// Return a mock instance so it proceeds to GetSlot
			return &models.NodeInstance{ID: "inst1"}, nil
		},
		GetSlotFunc: func(ctx context.Context, instanceID, slotKey string) (*models.NodeInstanceSlot, error) {
			return nil, fmt.Errorf("db error")
		},
	}
	svc := services.NewJourneyService(mock, pb, config.AppConfig{}, nil, nil, nil)

	err := svc.AttachUpload(context.Background(), "t1", "u1", "node1", "slot1", "up", "orig", 100)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}

type mockStorage struct{}

func (m *mockStorage) PresignPut(ctx context.Context, path, contentType string, expiry time.Duration) (string, error) {
	return fmt.Sprintf("https://mock-s3/%s", path), nil
}
func (m *mockStorage) PresignGet(ctx context.Context, path string, expiry time.Duration) (string, error) {
	return fmt.Sprintf("https://mock-s3-get/%s", path), nil
}
func (m *mockStorage) ObjectExists(ctx context.Context, key string) (bool, error) {
	return true, nil
}
func (m *mockStorage) Bucket() string {
	return "test-bucket"
}
func (m *mockStorage) Delete(ctx context.Context, path string) error { return nil }

type errorStorage struct{ mockStorage }
func (e *errorStorage) PresignPut(ctx context.Context, path, contentType string, expiry time.Duration) (string, error) {
	return "", fmt.Errorf("storage error")
}

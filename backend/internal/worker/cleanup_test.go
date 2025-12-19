package worker

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockS3Client implements S3ClientInterface for testing
type MockS3Client struct {
	ListObjectsFunc func(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	DeleteObjectFunc func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	DeletedKeys      []string
}

func (m *MockS3Client) ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	if m.ListObjectsFunc != nil {
		return m.ListObjectsFunc(ctx, params, optFns...)
	}
	return &s3.ListObjectsV2Output{}, nil
}

func (m *MockS3Client) DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	if params.Key != nil {
		m.DeletedKeys = append(m.DeletedKeys, *params.Key)
	}
	if m.DeleteObjectFunc != nil {
		return m.DeleteObjectFunc(ctx, params, optFns...)
	}
	return &s3.DeleteObjectOutput{}, nil
}

func TestCleanupWorker_RunCleanup(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// 1. Setup Data
	userID := "10000000-0000-0000-0000-000000000001"
	tenantID := "00000000-0000-0000-0000-000000000001"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'cleanupuser', 'cleanup@ex.com', 'Clean', 'Up', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	docID := "20000000-0000-0000-0000-000000000002"
	_, err = db.Exec(`INSERT INTO documents (id, tenant_id, user_id, kind, title) 
		VALUES ($1, $2, $3, 'other', 'Title')`, docID, tenantID, userID)
	require.NoError(t, err)

	// Referenced file (should NOT be deleted)
	referencedKey := "referenced-file.pdf"
	_, err = db.Exec(`INSERT INTO document_versions (tenant_id, document_id, storage_path, object_key, bucket, mime_type, size_bytes, uploaded_by)
		VALUES ($1, $2, 'local/path', $3, 'test-bucket', 'application/pdf', 100, $4)`, tenantID, docID, referencedKey, userID)
	require.NoError(t, err)

	// 2. Setup Mock S3
	oldTime := time.Now().Add(-48 * time.Hour)
	recentTime := time.Now().Add(-1 * time.Hour)

	mockS3 := &MockS3Client{
		ListObjectsFunc: func(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
			return &s3.ListObjectsV2Output{
				Contents: []types.Object{
					{Key: aws.String(referencedKey), LastModified: &oldTime}, // Should keep (referenced)
					{Key: aws.String("orphan-old.pdf"), LastModified: &oldTime}, // Should delete (orphan, old)
					{Key: aws.String("orphan-new.pdf"), LastModified: &recentTime}, // Should keep (orphan, new)
				},
				IsTruncated: aws.Bool(false),
			}, nil
		},
	}

	// 3. Run Worker
	worker := NewCleanupWorker(db, mockS3, "test-bucket")
	worker.runCleanup(context.Background())

	// 4. Verify Assertions
	assert.Len(t, mockS3.DeletedKeys, 1, "Should delete exactly one file")
	if len(mockS3.DeletedKeys) > 0 {
		assert.Equal(t, "orphan-old.pdf", mockS3.DeletedKeys[0])
	}
}

func TestCleanupWorker_RunCleanup_Pagination(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	oldTime := time.Now().Add(-48 * time.Hour)
	
	page1 := true
	mockS3 := &MockS3Client{
		ListObjectsFunc: func(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
			if page1 {
				page1 = false
				return &s3.ListObjectsV2Output{
					Contents: []types.Object{
						{Key: aws.String("orphan-1.pdf"), LastModified: &oldTime},
					},
					IsTruncated: aws.Bool(true),
					NextContinuationToken: aws.String("token"),
				}, nil
			}
			return &s3.ListObjectsV2Output{
				Contents: []types.Object{
					{Key: aws.String("orphan-2.pdf"), LastModified: &oldTime},
				},
				IsTruncated: aws.Bool(false),
			}, nil
		},
	}

	worker := NewCleanupWorker(db, mockS3, "test-bucket")
	worker.runCleanup(context.Background())

	assert.Len(t, mockS3.DeletedKeys, 2)
	assert.Contains(t, mockS3.DeletedKeys, "orphan-1.pdf")
	assert.Contains(t, mockS3.DeletedKeys, "orphan-2.pdf")
}

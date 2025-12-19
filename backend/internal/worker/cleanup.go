package worker

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jmoiron/sqlx"
)

// S3ClientInterface abstracts S3 operations for testing
type S3ClientInterface interface {
	DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

type CleanupWorker struct {
	db       *sqlx.DB
	s3Client S3ClientInterface
	bucket   string
}

func NewCleanupWorker(db *sqlx.DB, s3Client S3ClientInterface, bucket string) *CleanupWorker {
	return &CleanupWorker{
		db:       db,
		s3Client: s3Client,
		bucket:   bucket,
	}
}

// Start begins the background cleanup process
func (w *CleanupWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	log.Println("[CleanupWorker] Started - will run every 24 hours")

	// Run immediately on start
	w.runCleanup(ctx)

	for {
		select {
		case <-ticker.C:
			w.runCleanup(ctx)
		case <-ctx.Done():
			log.Println("[CleanupWorker] Stopped")
			return
		}
	}
}

func (w *CleanupWorker) runCleanup(ctx context.Context) {
	log.Println("[CleanupWorker] Starting orphan file cleanup...")

	if w.s3Client == nil {
		log.Println("[CleanupWorker] S3 client not configured, skipping cleanup")
		return
	}

	// Get all object keys from database
	var dbKeys []string
	err := w.db.Select(&dbKeys, `SELECT DISTINCT object_key FROM document_versions WHERE object_key IS NOT NULL AND object_key != ''`)
	if err != nil {
		log.Printf("[CleanupWorker] Error fetching DB keys: %v", err)
		return
	}

	// Create a map for quick lookup
	dbKeyMap := make(map[string]bool)
	for _, key := range dbKeys {
		dbKeyMap[key] = true
	}

	log.Printf("[CleanupWorker] Found %d files in database", len(dbKeys))

	// List all objects in S3 bucket
	// Note: Using ListObjectsV2 directly instead of Paginator for interface simplicity compatibility
	// In a full production mock, we might want to stick to Paginator but it requires more complex mocking.
	// For simplicity in this worker refactor, we'll iterate manually or assume a single page for the interface constraint,
	// OR better: keep Using paginator but we need to ensure our interface covers what paginator expects?
	// Paginator takes *s3.Client. It doesn't take an interface effortlessly without a wrapper.
	// To keep it testable, we'll manually loop with ContinuationToken which is what Paginator does.
	
	input := &s3.ListObjectsV2Input{
		Bucket: &w.bucket,
	}

	orphanCount := 0
	deletedCount := 0
	cutoffTime := time.Now().Add(-24 * time.Hour)

	for {
		output, err := w.s3Client.ListObjectsV2(ctx, input)
		if err != nil {
			log.Printf("[CleanupWorker] Error listing objects: %v", err)
			break
		}

		for _, object := range output.Contents {
			if object.Key == nil {
				continue
			}

			// Check if object exists in database
			if _, exists := dbKeyMap[*object.Key]; !exists {
				// Object is orphaned
				orphanCount++

				// Only delete if older than 24 hours
				if object.LastModified != nil && object.LastModified.Before(cutoffTime) {
					_, err := w.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
						Bucket: &w.bucket,
						Key:    object.Key,
					})
					if err != nil {
						log.Printf("[CleanupWorker] Error deleting orphan %s: %v", *object.Key, err)
					} else {
						deletedCount++
						log.Printf("[CleanupWorker] Deleted orphan file: %s (age: %v)", *object.Key, time.Since(*object.LastModified))
					}
				} else if object.LastModified != nil {
					log.Printf("[CleanupWorker] Found recent orphan (< 24h), keeping: %s", *object.Key)
				}
			}
		}

		if output.IsTruncated != nil && *output.IsTruncated {
			input.ContinuationToken = output.NextContinuationToken
		} else {
			break
		}
	}

	log.Printf("[CleanupWorker] Cleanup complete - Found %d orphans, deleted %d files", orphanCount, deletedCount)
}

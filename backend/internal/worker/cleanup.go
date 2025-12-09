package worker

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jmoiron/sqlx"
)

type CleanupWorker struct {
	db       *sqlx.DB
	s3Client *s3.Client
	bucket   string
}

func NewCleanupWorker(db *sqlx.DB, s3Client *s3.Client, bucket string) *CleanupWorker {
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
	paginator := s3.NewListObjectsV2Paginator(w.s3Client, &s3.ListObjectsV2Input{
		Bucket: &w.bucket,
	})

	orphanCount := 0
	deletedCount := 0
	cutoffTime := time.Now().Add(-24 * time.Hour)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			log.Printf("[CleanupWorker] Error listing objects: %v", err)
			break
		}

		for _, object := range page.Contents {
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
	}

	log.Printf("[CleanupWorker] Cleanup complete - Found %d orphans, deleted %d files", orphanCount, deletedCount)
}

package services_test

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestMockStorageClient(t *testing.T) {
	m := &services.MockStorageClient{}
	ctx := context.Background()

	url, err := m.PresignPut(ctx, "k1", "c1", 0)
	assert.NoError(t, err)
	assert.Contains(t, url, "k1")

	url, err = m.PresignGet(ctx, "k1", 0)
	assert.NoError(t, err)
	assert.Contains(t, url, "k1")

	exists, err := m.ObjectExists(ctx, "k1")
	assert.NoError(t, err)
	assert.True(t, exists)

	assert.Equal(t, "mock-bucket", m.Bucket())

	m.BucketName = "custom"
	assert.Equal(t, "custom", m.Bucket())

	m.ObjectExistsFn = func(ctx context.Context, key string) (bool, error) {
		return false, nil
	}
	exists, _ = m.ObjectExists(ctx, "k1")
	assert.False(t, exists)
}

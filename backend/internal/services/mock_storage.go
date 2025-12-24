package services

import (
	"context"
	"time"
)

type MockStorageClient struct {
	PresignPutFn   func(ctx context.Context, key, contentType string, expires time.Duration) (string, error)
	PresignGetFn   func(ctx context.Context, key string, expires time.Duration) (string, error)
	ObjectExistsFn func(ctx context.Context, key string) (bool, error)
	BucketName     string
}

func (m *MockStorageClient) PresignPut(ctx context.Context, key, contentType string, expires time.Duration) (string, error) {
	if m.PresignPutFn != nil {
		return m.PresignPutFn(ctx, key, contentType, expires)
	}
	return "http://mock-s3.com/" + key, nil
}

func (m *MockStorageClient) PresignGet(ctx context.Context, key string, expires time.Duration) (string, error) {
	if m.PresignGetFn != nil {
		return m.PresignGetFn(ctx, key, expires)
	}
	return "http://mock-s3.com/" + key, nil
}

func (m *MockStorageClient) ObjectExists(ctx context.Context, key string) (bool, error) {
	if m.ObjectExistsFn != nil {
		return m.ObjectExistsFn(ctx, key)
	}
	return true, nil
}

func (m *MockStorageClient) Bucket() string {
	if m.BucketName != "" {
		return m.BucketName
	}
	return "mock-bucket"
}

package services_test

import (
	"context"
	"os"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestS3Helpers(t *testing.T) {
	t.Run("ValidateContentType", func(t *testing.T) {
		assert.NoError(t, services.ValidateContentType("application/pdf"))
		assert.NoError(t, services.ValidateContentType("image/png"))
		assert.Error(t, services.ValidateContentType("application/evil"))
	})

	t.Run("ValidateFileSize", func(t *testing.T) {
		os.Setenv("S3_MAX_FILE_SIZE_MB", "10")
		defer os.Unsetenv("S3_MAX_FILE_SIZE_MB")

		assert.NoError(t, services.ValidateFileSize(5 * 1024 * 1024))
		assert.Error(t, services.ValidateFileSize(15 * 1024 * 1024))
		assert.Error(t, services.ValidateFileSize(0))
		assert.Error(t, services.ValidateFileSize(-1))
	})

	t.Run("GetPresignExpires", func(t *testing.T) {
		os.Setenv("S3_PRESIGN_EXPIRES_MINUTES", "30")
		defer os.Unsetenv("S3_PRESIGN_EXPIRES_MINUTES")
		
		exp := services.GetPresignExpires()
		assert.Equal(t, float64(30), exp.Minutes())
	})

	t.Run("S3Client_NilMethods", func(t *testing.T) {
		var client *services.S3Client
		ctx := context.Background()
		url, err := client.PresignPut(ctx, "", "", 0)
		assert.NoError(t, err)
		assert.Empty(t, url)

		url, err = client.PresignGet(ctx, "", 0)
		assert.NoError(t, err)
		assert.Empty(t, url)

		exists, err := client.ObjectExists(ctx, "")
		assert.NoError(t, err)
		assert.False(t, exists)

		assert.Empty(t, client.Bucket())
		assert.Nil(t, client.Client())
	})
}

func TestS3Client_EnvMissing(t *testing.T) {
	os.Unsetenv("S3_BUCKET")
	os.Unsetenv("S3_BUCKET_NAME")
	cl, err := services.NewS3FromEnv()
	assert.NoError(t, err)
	assert.Nil(t, cl)

	os.Setenv("S3_BUCKET", "test")
	os.Unsetenv("S3_ACCESS_KEY_ID")
	os.Unsetenv("S3_ACCESS_KEY")
	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("S3_SECRET_ACCESS_KEY")
	os.Unsetenv("S3_SECRET_KEY")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	
	cl, err = services.NewS3FromEnv()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be set")
}

func TestS3Client_PathStyle(t *testing.T) {
	os.Setenv("S3_BUCKET", "test")
	os.Setenv("S3_ACCESS_KEY", "ak")
	os.Setenv("S3_SECRET_KEY", "sk")
	
	t.Run("PathStyle True", func(t *testing.T) {
		os.Setenv("S3_USE_PATH_STYLE", "true")
		cl, err := services.NewS3FromEnv()
		assert.NoError(t, err)
		assert.NotNil(t, cl)
	})

	t.Run("Default PathStyle with Endpoint", func(t *testing.T) {
		os.Setenv("S3_ENDPOINT", "http://localhost:9000")
		os.Unsetenv("S3_USE_PATH_STYLE")
		cl, err := services.NewS3FromEnv()
		assert.NoError(t, err)
		assert.NotNil(t, cl)
	})
}

func TestS3Helpers_Int(t *testing.T) {
	t.Run("getEnvInt Defaults", func(t *testing.T) {
		os.Setenv("S3_MAX_FILE_SIZE_MB", "invalid")
		assert.NoError(t, services.ValidateFileSize(10 * 1024 * 1024)) // Should use default 100MB
		
		os.Setenv("S3_MAX_FILE_SIZE_MB", "5")
		assert.Error(t, services.ValidateFileSize(10 * 1024 * 1024))
	})
}

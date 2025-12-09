package services

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewS3FromEnv_NoBucket(t *testing.T) {
	// Clear all S3 env vars
	os.Unsetenv("S3_BUCKET")
	os.Unsetenv("S3_BUCKET_NAME")

	client, err := NewS3FromEnv()
	assert.NoError(t, err)
	assert.Nil(t, client)
}

func TestNewS3FromEnv_MissingCredentials(t *testing.T) {
	// Set bucket but no credentials
	os.Setenv("S3_BUCKET", "test-bucket")
	os.Unsetenv("S3_ACCESS_KEY_ID")
	os.Unsetenv("S3_ACCESS_KEY")
	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("S3_SECRET_ACCESS_KEY")
	os.Unsetenv("S3_SECRET_KEY")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	defer os.Unsetenv("S3_BUCKET")

	client, err := NewS3FromEnv()
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Contains(t, err.Error(), "S3_ACCESS_KEY and S3_SECRET_KEY must be set")
}

func TestNewS3FromEnv_WithConfig(t *testing.T) {
	os.Setenv("S3_BUCKET", "test-bucket")
	os.Setenv("S3_ENDPOINT", "http://localhost:9000")
	os.Setenv("S3_ACCESS_KEY_ID", "minioadmin")
	os.Setenv("S3_SECRET_ACCESS_KEY", "minioadmin")
	os.Setenv("S3_REGION", "us-east-1")
	defer func() {
		os.Unsetenv("S3_BUCKET")
		os.Unsetenv("S3_ENDPOINT")
		os.Unsetenv("S3_ACCESS_KEY_ID")
		os.Unsetenv("S3_SECRET_ACCESS_KEY")
		os.Unsetenv("S3_REGION")
	}()

	client, err := NewS3FromEnv()
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "test-bucket", client.Bucket())
}

func TestS3Client_Bucket_NilClient(t *testing.T) {
	var client *S3Client
	assert.Equal(t, "", client.Bucket())
}

func TestS3Client_Client_NilClient(t *testing.T) {
	var client *S3Client
	assert.Nil(t, client.Client())
}

func TestS3Client_PresignPut_NilClient(t *testing.T) {
	var client *S3Client
	url, err := client.PresignPut("key", "application/pdf", 15*time.Minute)
	assert.NoError(t, err)
	assert.Equal(t, "", url)
}

func TestS3Client_PresignGet_NilClient(t *testing.T) {
	var client *S3Client
	url, err := client.PresignGet("key", 15*time.Minute)
	assert.NoError(t, err)
	assert.Equal(t, "", url)
}

func TestS3Client_ObjectExists_NilClient(t *testing.T) {
	var client *S3Client
	exists, err := client.ObjectExists("key")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestValidateContentType_Allowed(t *testing.T) {
	allowedTypes := []string{
		"application/pdf",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"image/jpeg",
		"image/png",
		"text/plain",
	}

	for _, ct := range allowedTypes {
		err := ValidateContentType(ct)
		assert.NoError(t, err, "Expected %s to be allowed", ct)
	}
}

func TestValidateContentType_NotAllowed(t *testing.T) {
	err := ValidateContentType("application/x-executable")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported content type")

	err = ValidateContentType("video/mp4")
	assert.Error(t, err)
}

func TestValidateFileSize_Valid(t *testing.T) {
	// Default max is 100MB
	err := ValidateFileSize(1024 * 1024) // 1MB
	assert.NoError(t, err)

	err = ValidateFileSize(50 * 1024 * 1024) // 50MB
	assert.NoError(t, err)
}

func TestValidateFileSize_TooLarge(t *testing.T) {
	// Default max is 100MB
	err := ValidateFileSize(200 * 1024 * 1024) // 200MB
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds maximum")
}

func TestValidateFileSize_Invalid(t *testing.T) {
	err := ValidateFileSize(0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid file size")

	err = ValidateFileSize(-100)
	assert.Error(t, err)
}

func TestGetPresignExpires_Default(t *testing.T) {
	os.Unsetenv("S3_PRESIGN_EXPIRES_MINUTES")
	expires := GetPresignExpires()
	assert.Equal(t, 15*time.Minute, expires)
}

func TestGetPresignExpires_Custom(t *testing.T) {
	os.Setenv("S3_PRESIGN_EXPIRES_MINUTES", "30")
	defer os.Unsetenv("S3_PRESIGN_EXPIRES_MINUTES")

	expires := GetPresignExpires()
	assert.Equal(t, 30*time.Minute, expires)
}

func TestFirstNonEmpty(t *testing.T) {
	assert.Equal(t, "first", firstNonEmpty("first", "second"))
	assert.Equal(t, "second", firstNonEmpty("", "second", "third"))
	assert.Equal(t, "third", firstNonEmpty("", "", "third"))
	assert.Equal(t, "", firstNonEmpty("", "", ""))
	assert.Equal(t, "value", firstNonEmpty("  ", "value")) // Whitespace-only is not empty but trimmed
}

package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Config struct {
	Endpoint     string
	Region       string
	Bucket       string
	AccessKey    string
	SecretKey    string
	UsePathStyle bool
}

type StorageClient interface {
	PresignPut(ctx context.Context, key, contentType string, expires time.Duration) (string, error)
	PresignGet(ctx context.Context, key string, expires time.Duration) (string, error)
	ObjectExists(ctx context.Context, key string) (bool, error)
	Bucket() string
}

type S3Client struct {
	cfg    S3Config
	client *s3.Client
}

func NewS3FromEnv() (*S3Client, error) {
	bucket := firstNonEmpty(os.Getenv("S3_BUCKET"), os.Getenv("S3_BUCKET_NAME"))
	if bucket == "" {
		return nil, nil
	}
	region := firstNonEmpty(os.Getenv("S3_REGION"), os.Getenv("AWS_REGION"), "us-east-1")
	endpoint := os.Getenv("S3_ENDPOINT")
	access := firstNonEmpty(os.Getenv("S3_ACCESS_KEY_ID"), os.Getenv("S3_ACCESS_KEY"), os.Getenv("AWS_ACCESS_KEY_ID"))
	secret := firstNonEmpty(os.Getenv("S3_SECRET_ACCESS_KEY"), os.Getenv("S3_SECRET_KEY"), os.Getenv("AWS_SECRET_ACCESS_KEY"))
	
	// Require credentials if bucket is configured (security best practice)
	if access == "" || secret == "" {
		return nil, fmt.Errorf("S3_ACCESS_KEY and S3_SECRET_KEY must be set when S3_BUCKET is configured")
	}
	
	usePathStyleEnv := strings.ToLower(getEnv("S3_USE_PATH_STYLE", ""))
	usePathStyle := usePathStyleEnv == "true"
	if usePathStyleEnv == "" {
		usePathStyle = endpoint != ""
	}
	scfg := S3Config{
		Endpoint:     endpoint,
		Region:       region,
		Bucket:       bucket,
		AccessKey:    access,
		SecretKey:    secret,
		UsePathStyle: usePathStyle,
	}
	var cfg aws.Config
	var err error
	
	credProvider := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(access, secret, ""))
	
	loadOpts := []func(*config.LoadOptions) error{
		config.WithRegion(region),
		config.WithCredentialsProvider(credProvider),
	}
	cfg, err = config.LoadDefaultConfig(context.Background(), loadOpts...)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = scfg.UsePathStyle
		if scfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(scfg.Endpoint)
		}
	})
	return &S3Client{cfg: scfg, client: client}, nil
}

func (s *S3Client) PresignPut(ctx context.Context, objectKey, contentType string, expires time.Duration) (string, error) {
	if s == nil || s.client == nil {
		return "", nil
	}
	log.Printf("[S3] Presigning PUT for key=%s bucket=%s expires=%v", objectKey, s.cfg.Bucket, expires)
	ps := s3.NewPresignClient(s.client)
	req, err := ps.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      &s.cfg.Bucket,
		Key:         &objectKey,
		ContentType: &contentType,
	}, s3.WithPresignExpires(expires))
	if err != nil {
		log.Printf("[S3] PresignPut failed: %v", err)
		return "", err
	}
	return req.URL, nil
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func (s *S3Client) PresignGet(ctx context.Context, objectKey string, expires time.Duration) (string, error) {
	if s == nil || s.client == nil {
		return "", nil
	}
	log.Printf("[S3] Presigning GET for key=%s bucket=%s expires=%v", objectKey, s.cfg.Bucket, expires)
	ps := s3.NewPresignClient(s.client)
	req, err := ps.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.cfg.Bucket,
		Key:    &objectKey,
	}, s3.WithPresignExpires(expires))
	if err != nil {
		log.Printf("[S3] PresignGet failed: %v", err)
		return "", err
	}
	return req.URL, nil
}

func (s *S3Client) Bucket() string {
	if s == nil {
		return ""
	}
	return s.cfg.Bucket
}

// Client returns the underlying S3 client for advanced operations
func (s *S3Client) Client() *s3.Client {
	if s == nil {
		return nil
	}
	return s.client
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

// ObjectExists checks if an object exists in S3 bucket
func (s *S3Client) ObjectExists(ctx context.Context, objectKey string) (bool, error) {
	if s == nil || s.client == nil {
		return false, nil
	}
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &s.cfg.Bucket,
		Key:    &objectKey,
	})
	if err != nil {
		// Object doesn't exist or other error
		return false, err
	}
	return true, nil
}

// GetPresignExpires returns the presign URL expiration time from env or default
func GetPresignExpires() time.Duration {
	minutes := getEnvInt("S3_PRESIGN_EXPIRES_MINUTES", 15)
	return time.Duration(minutes) * time.Minute
}

// ValidateContentType checks if content type is allowed
func ValidateContentType(contentType string) error {
	allowedTypes := map[string]bool{
		"application/pdf":                                                   true,
		"application/msword":                                                true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"application/vnd.ms-excel":                                          true,
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
		"image/jpeg":                                                        true,
		"image/png":                                                         true,
		"image/gif":                                                         true,
		"text/plain":                                                        true,
		"application/zip":                                                   true,
	}
	if !allowedTypes[contentType] {
		return fmt.Errorf("unsupported content type: %s", contentType)
	}
	return nil
}

// ValidateFileSize checks if file size is within limits
func ValidateFileSize(sizeBytes int64) error {
	maxSize := int64(getEnvInt("S3_MAX_FILE_SIZE_MB", 100)) * 1024 * 1024 // Default 100MB
	if sizeBytes > maxSize {
		return fmt.Errorf("file size %d bytes exceeds maximum %d bytes", sizeBytes, maxSize)
	}
	if sizeBytes <= 0 {
		return fmt.Errorf("invalid file size: %d", sizeBytes)
	}
	return nil
}

func getEnvInt(key string, defaultVal int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return intVal
}

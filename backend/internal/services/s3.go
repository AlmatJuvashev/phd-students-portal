package services

import (
	"context"
	"os"
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

type S3Client struct {
	cfg    S3Config
	client *s3.Client
}

func NewS3FromEnv() (*S3Client, error) {
	endpoint := os.Getenv("S3_ENDPOINT")
	if endpoint == "" {
		return nil, nil
	} // feature disabled
	scfg := S3Config{
		Endpoint:     endpoint,
		Region:       getEnv("S3_REGION", "us-east-1"),
		Bucket:       getEnv("S3_BUCKET", "phd-portal"),
		AccessKey:    os.Getenv("S3_ACCESS_KEY"),
		SecretKey:    os.Getenv("S3_SECRET_KEY"),
		UsePathStyle: getEnv("S3_USE_PATH_STYLE", "true") == "true",
	}
	creds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(scfg.AccessKey, scfg.SecretKey, ""))
	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{URL: scfg.Endpoint, HostnameImmutable: true}, nil
	})
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(scfg.Region),
		config.WithCredentialsProvider(creds),
		config.WithEndpointResolverWithOptions(resolver),
	)
	if err != nil {
		return nil, err
	}
	return &S3Client{
		cfg:    scfg,
		client: s3.NewFromConfig(cfg, func(o *s3.Options) { o.UsePathStyle = scfg.UsePathStyle }),
	}, nil
}

func (s *S3Client) PresignPut(objectKey, contentType string, expires time.Duration) (string, error) {
	if s == nil || s.client == nil {
		return "", nil
	}
	ps := s3.NewPresignClient(s.client)
	req, err := ps.PresignPutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      &s.cfg.Bucket,
		Key:         &objectKey,
		ContentType: &contentType,
	}, s3.WithPresignExpires(expires))
	if err != nil {
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

func (s *S3Client) PresignGet(objectKey string, expires time.Duration) (string, error) {
	if s == nil || s.client == nil {
		return "", nil
	}
	ps := s3.NewPresignClient(s.client)
	req, err := ps.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: &s.cfg.Bucket,
		Key:    &objectKey,
	}, s3.WithPresignExpires(expires))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

package services

import (
	"context"
	"os"
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
	credProvider := aws.CredentialsProvider(nil)
	if access != "" && secret != "" {
		credProvider = aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(access, secret, ""))
	}
	if endpoint != "" {
		resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{URL: endpoint, HostnameImmutable: true}, nil
		})
		opts := []func(*config.LoadOptions) error{
			config.WithRegion(region),
			config.WithEndpointResolverWithOptions(resolver),
		}
		if credProvider != nil {
			opts = append(opts, config.WithCredentialsProvider(credProvider))
		}
		cfg, err = config.LoadDefaultConfig(context.Background(), opts...)
	} else {
		loadOpts := []func(*config.LoadOptions) error{
			config.WithRegion(region),
		}
		if credProvider != nil {
			loadOpts = append(loadOpts, config.WithCredentialsProvider(credProvider))
		}
		cfg, err = config.LoadDefaultConfig(context.Background(), loadOpts...)
	}
	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(cfg, func(o *s3.Options) { o.UsePathStyle = scfg.UsePathStyle })
	return &S3Client{cfg: scfg, client: client}, nil
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

func (s *S3Client) Bucket() string {
	if s == nil {
		return ""
	}
	return s.cfg.Bucket
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

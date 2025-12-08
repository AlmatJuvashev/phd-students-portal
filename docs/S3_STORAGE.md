# S3 File Storage Guide

## Overview

This application uses **S3-compatible object storage** for managing file uploads. This guide explains how the system works, how to configure it, and how to work with it as a developer.

---

## Table of Contents

1. [What is S3?](#what-is-s3)
2. [Architecture](#architecture)
3. [Configuration](#configuration)
4. [How File Upload Works](#how-file-upload-works)
5. [Backend Implementation](#backend-implementation)
6. [Frontend Implementation](#frontend-implementation)
7. [Local Development Setup](#local-development-setup)
8. [Production Deployment](#production-deployment)
9. [Troubleshooting](#troubleshooting)
10. [Security Best Practices](#security-best-practices)

---

## What is S3?

**Amazon S3 (Simple Storage Service)** is an object storage service. Unlike traditional file systems, S3 stores files as "objects" in "buckets" with unique keys (paths).

### Key Concepts:

- **Bucket**: A container for objects (like a root folder)
- **Object Key**: Unique identifier for a file (e.g., `nodes/user123/node456/slot1/20251115-abc123-document.pdf`)
- **Presigned URL**: Temporary URL that allows upload/download without exposing credentials
- **S3-Compatible Storage**: Services like MinIO, DigitalOcean Spaces, AWS S3 that follow the same API

### Why Use S3?

✅ **Scalable**: Handle millions of files without performance degradation  
✅ **Reliable**: Built-in redundancy and durability  
✅ **Secure**: Fine-grained access control and encryption  
✅ **Cost-effective**: Pay only for what you use  
✅ **Direct upload**: Files go straight from browser to S3 (no backend bottleneck)

---

## Architecture

### Upload Flow

```
┌─────────┐      1. Request presigned URL      ┌─────────┐
│ Browser ├────────────────────────────────────►│ Backend │
└─────────┘      (filename, content-type)      └────┬────┘
     │                                               │
     │                                               │ 2. Generate presigned URL
     │                                               │    (AWS SDK)
     │                                               │
     │           3. Return presigned URL        ┌────▼────┐
     │◄─────────────────────────────────────────┤   S3    │
     │           (expires in 15 min)            │  Client │
     │                                          └─────────┘
     │
     │           4. Direct PUT to S3
     ├──────────────────────────────────────────►┌─────────┐
     │           (file bytes)                    │   S3    │
     │                                           │ Bucket  │
     │           5. Success (200 OK + ETag)      └─────────┘
     │◄──────────────────────────────────────────┘
     │
     │           6. Confirm upload
     └──────────────────────────────────────────►┌─────────┐
                (object_key, etag, size)         │ Backend │
                                                  └─────────┘
```

### Download Flow

```
┌─────────┐      1. Request download URL       ┌─────────┐
│ Browser ├────────────────────────────────────►│ Backend │
└─────────┘      (document_id/node_id)         └────┬────┘
     │                                               │
     │                                               │ 2. Query DB for object_key
     │                                               │
     │           3. Generate presigned GET      ┌────▼────┐
     │◄─────────────────────────────────────────┤   S3    │
     │           (temporary download URL)       │  Client │
     │                                          └─────────┘
     │
     │           4. Direct GET from S3
     └──────────────────────────────────────────►┌─────────┐
                (download file)                  │   S3    │
                                                 │ Bucket  │
                                                 └─────────┘
```

---

## Configuration

### Environment Variables

Add these to your `backend/.env` file:

```bash
# S3-compatible storage (leave S3_BUCKET empty to disable)
S3_ENDPOINT=                        # For MinIO/custom: http://localhost:9000
                                    # For AWS S3: leave empty
S3_REGION=us-east-1                 # AWS region or MinIO region
S3_BUCKET=phd-portal                # Bucket name (required)
S3_ACCESS_KEY=YOUR_ACCESS_KEY       # Access key ID (required)
S3_SECRET_KEY=YOUR_SECRET_KEY       # Secret access key (required)
S3_USE_PATH_STYLE=true              # true for MinIO, false for AWS S3
S3_PRESIGN_EXPIRES_MINUTES=15       # How long presigned URLs are valid (default: 15)
S3_MAX_FILE_SIZE_MB=100             # Maximum file size in MB (default: 100)
```

### Alternative Variable Names

The system supports multiple environment variable names for compatibility:

| Primary         | Alternatives                                    |
| --------------- | ----------------------------------------------- |
| `S3_BUCKET`     | `S3_BUCKET_NAME`                                |
| `S3_ACCESS_KEY` | `S3_ACCESS_KEY_ID`, `AWS_ACCESS_KEY_ID`         |
| `S3_SECRET_KEY` | `S3_SECRET_ACCESS_KEY`, `AWS_SECRET_ACCESS_KEY` |
| `S3_REGION`     | `AWS_REGION`                                    |

---

## How File Upload Works

### Step-by-Step Process

#### 1. User Selects File (Frontend)

```tsx
// User clicks "Upload" button
const handleFileSelected = async (slotKey: string, files: FileList | null) => {
  const file = files?.[0];
  if (!file) return;

  // File metadata
  const contentType = file.type || "application/octet-stream";
  const sizeBytes = file.size;
  const filename = file.name;
```

#### 2. Request Presigned URL (Frontend → Backend)

```tsx
// Request presigned URL from backend
const presign = await presignNodeUpload(nodeId, {
  slot_key: slotKey,
  filename: filename,
  content_type: contentType,
  size_bytes: sizeBytes,
});

// Backend response:
// {
//   upload_url: "https://s3.amazonaws.com/...",
//   object_key: "nodes/user123/node456/...",
//   expires_in: 900,
//   bucket: "phd-portal"
// }
```

#### 3. Backend Generates Presigned URL

```go
// Validate file size and content type
if err := services.ValidateFileSize(req.SizeBytes); err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}
if err := services.ValidateContentType(req.ContentType); err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}

// Generate unique object key
objectKey := storage.BuildNodeObjectKey(uid, nodeID, req.SlotKey, req.Filename)
// Result: "nodes/d75a1afc.../S1_profile/main_doc/20251115-abc123-thesis.pdf"

// Create presigned PUT URL (valid for 15 minutes)
expires := services.GetPresignExpires()
url, err := s3c.PresignPut(objectKey, req.ContentType, expires)
```

#### 4. Direct Upload to S3 (Frontend → S3)

```tsx
// Upload file directly to S3 (NOT through backend)
const uploadResp = await fetch(presign.upload_url, {
  method: "PUT",
  headers: {
    "Content-Type": contentType,
  },
  body: file, // Raw file bytes
});

if (!uploadResp.ok) {
  throw new Error(`Upload failed (${uploadResp.status})`);
}

// Get ETag from response headers (used for validation)
const etag = uploadResp.headers.get("ETag")?.replace(/"/g, "");
```

#### 5. Confirm Upload (Frontend → Backend)

```tsx
// Tell backend that upload is complete
await attachNodeUpload(nodeId, {
  slot_key: slotKey,
  filename: filename,
  object_key: presign.object_key,
  content_type: contentType,
  size_bytes: sizeBytes,
  etag: etag, // For S3 validation
});

// Backend saves metadata to database
```

---

## Backend Implementation

### Core Service: `internal/services/s3.go`

#### Initialization

```go
func NewS3FromEnv() (*S3Client, error) {
    // Read configuration from environment
    bucket := os.Getenv("S3_BUCKET")
    if bucket == "" {
        return nil, nil // S3 disabled
    }

    // Require credentials for security
    access := os.Getenv("S3_ACCESS_KEY")
    secret := os.Getenv("S3_SECRET_KEY")
    if access == "" || secret == "" {
        return nil, fmt.Errorf("S3_ACCESS_KEY and S3_SECRET_KEY required")
    }

    // Create AWS SDK client
    credProvider := credentials.NewStaticCredentialsProvider(access, secret, "")
    cfg, err := config.LoadDefaultConfig(context.Background(),
        config.WithRegion(region),
        config.WithCredentialsProvider(credProvider),
    )

    client := s3.NewFromConfig(cfg, func(o *s3.Options) {
        o.UsePathStyle = usePathStyle // true for MinIO
    })

    return &S3Client{cfg: scfg, client: client}, nil
}
```

#### Presigned PUT (Upload)

```go
func (s *S3Client) PresignPut(objectKey, contentType string, expires time.Duration) (string, error) {
    log.Printf("[S3] Presigning PUT for key=%s bucket=%s expires=%v",
        objectKey, s.cfg.Bucket, expires)

    ps := s3.NewPresignClient(s.client)
    req, err := ps.PresignPutObject(context.Background(), &s3.PutObjectInput{
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
```

#### Presigned GET (Download)

```go
func (s *S3Client) PresignGet(objectKey string, expires time.Duration) (string, error) {
    log.Printf("[S3] Presigning GET for key=%s", objectKey)

    ps := s3.NewPresignClient(s.client)
    req, err := ps.PresignGetObject(context.Background(), &s3.GetObjectInput{
        Bucket: &s.cfg.Bucket,
        Key:    &objectKey,
    }, s3.WithPresignExpires(expires))

    return req.URL, nil
}
```

#### Validation Functions

```go
// Validate file size (default: 100MB max)
func ValidateFileSize(sizeBytes int64) error {
    maxSize := int64(getEnvInt("S3_MAX_FILE_SIZE_MB", 100)) * 1024 * 1024
    if sizeBytes > maxSize {
        return fmt.Errorf("file size %d bytes exceeds maximum %d bytes",
            sizeBytes, maxSize)
    }
    if sizeBytes <= 0 {
        return fmt.Errorf("invalid file size: %d", sizeBytes)
    }
    return nil
}

// Validate content type (whitelist)
func ValidateContentType(contentType string) error {
    allowedTypes := map[string]bool{
        "application/pdf":  true,
        "application/msword": true,
        "application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
        "image/jpeg": true,
        "image/png":  true,
        // ... more types
    }
    if !allowedTypes[contentType] {
        return fmt.Errorf("unsupported content type: %s", contentType)
    }
    return nil
}
```

### Object Key Generation: `internal/storage/object_key.go`

```go
func BuildNodeObjectKey(userID, nodeID, slotKey, filename string) string {
    // Sanitize inputs
    cleanName := sanitizeFilename(filename)
    segmentUser := sanitizeSegment(userID)
    segmentNode := sanitizeSegment(nodeID)
    segmentSlot := sanitizeSegment(slotKey)

    // Add timestamp and UUID for uniqueness
    stamp := time.Now().UTC().Format("20060102")
    uuid := uuid.NewString()

    // Result: nodes/d75a1afc/s1-profile/main-doc/20251115-abc123-thesis.pdf
    return fmt.Sprintf("nodes/%s/%s/%s/%s-%s-%s",
        segmentUser, segmentNode, segmentSlot, stamp, uuid, cleanName)
}
```

**Why this structure?**

- ✅ Organized by user/node/slot
- ✅ Timestamp for chronological ordering
- ✅ UUID prevents collisions
- ✅ Sanitized filenames (no special characters)

---

## Frontend Implementation

### API Client: `frontend/src/features/nodes/api.ts`

```typescript
// Request presigned URL
export async function presignNodeUpload(
  nodeId: string,
  payload: {
    slot_key: string;
    filename: string;
    content_type: string;
    size_bytes: number;
  }
) {
  return api(`/journey/nodes/${nodeId}/uploads/presign`, {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

// Confirm upload after S3 success
export async function attachNodeUpload(
  nodeId: string,
  payload: {
    slot_key: string;
    filename: string;
    object_key: string;
    content_type: string;
    size_bytes: number;
    etag?: string;
  }
) {
  return api(`/journey/nodes/${nodeId}/uploads/attach`, {
    method: "POST",
    body: JSON.stringify(payload),
  });
}
```

### Upload Component: `NodeAttachmentsSection.tsx`

```tsx
const handleFileSelected = async (slotKey: string, files: FileList | null) => {
  const file = files?.[0];
  if (!file) return;

  setUploadingSlot(slotKey);
  setMessage(null);

  try {
    const contentType = file.type || "application/octet-stream";

    // Step 1: Get presigned URL
    const presign = await presignNodeUpload(nodeId, {
      slot_key: slotKey,
      filename: file.name,
      content_type: contentType,
      size_bytes: file.size,
    });

    // Step 2: Upload directly to S3
    const uploadResp = await fetch(presign.upload_url, {
      method: "PUT",
      headers: { "Content-Type": contentType },
      body: file,
    });

    if (!uploadResp.ok) {
      throw new Error(`Upload failed (${uploadResp.status})`);
    }

    // Step 3: Confirm upload with backend
    const etag = uploadResp.headers.get("ETag")?.replace(/"/g, "");
    await attachNodeUpload(nodeId, {
      slot_key: slotKey,
      filename: file.name,
      object_key: presign.object_key,
      content_type: contentType,
      size_bytes: file.size,
      etag,
    });

    setMessage({ text: "File uploaded successfully", tone: "success" });
    onRefresh?.();
  } catch (error: any) {
    setMessage({ text: error.message || "Upload failed", tone: "error" });
  } finally {
    setUploadingSlot(null);
  }
};
```

---

## Local Development Setup

### Option 1: MinIO (Recommended for Development)

**MinIO** is a lightweight, S3-compatible storage server perfect for local development.

#### 1. Install MinIO

```bash
# macOS
brew install minio/stable/minio

# Linux
wget https://dl.min.io/server/minio/release/linux-amd64/minio
chmod +x minio
sudo mv minio /usr/local/bin/

# Docker
docker run -p 9090:9000 -p 9091:9091 \
  -e "MINIO_ROOT_USER=minioadmin" \
  -e "MINIO_ROOT_PASSWORD=minioadmin" \
  minio/minio server /data --console-address ":9091"
```

#### 2. Start MinIO Server

```bash
# Create data directory
mkdir -p ~/minio-data

# Start server
minio server ~/minio-data --console-address ":9091"

# Output:
# API: http://192.168.1.100:9090
# Console: http://192.168.1.100:9091
# RootUser: minioadmin
# RootPass: minioadmin
```

#### 3. Create Bucket

Visit **http://localhost:9001**, login with `minioadmin` / `minioadmin`:

1. Click **Buckets** → **Create Bucket**
2. Enter name: `phd-portal`
3. Click **Create**

#### 4. Configure CORS

Select bucket → **Configuration** → **CORS** → Add rule:

```json
{
  "CORSRules": [
    {
      "AllowedOrigins": ["http://localhost:5174"],
      "AllowedMethods": ["PUT", "GET", "HEAD"],
      "AllowedHeaders": ["*"],
      "ExposeHeaders": ["ETag"]
    }
  ]
}
```

#### 5. Update `.env`

```bash
S3_ENDPOINT=http://localhost:9000
S3_REGION=us-east-1
S3_BUCKET=phd-portal
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_USE_PATH_STYLE=true
S3_PRESIGN_EXPIRES_MINUTES=15
S3_MAX_FILE_SIZE_MB=100
```

#### 6. Test Upload

```bash
# Restart backend
cd backend
go run cmd/server/main.go

# You should see in logs:
# [S3] Presigning PUT for key=nodes/... bucket=phd-portal expires=15m
```

---

### Option 2: AWS S3 (Production-like)

#### 1. Create S3 Bucket

```bash
# AWS CLI
aws s3 mb s3://phd-portal-dev --region us-east-1
```

#### 2. Configure CORS

Create `cors.json`:

```json
[
  {
    "AllowedOrigins": ["http://localhost:5174"],
    "AllowedMethods": ["PUT", "GET", "HEAD"],
    "AllowedHeaders": ["*"],
    "ExposeHeaders": ["ETag"],
    "MaxAgeSeconds": 3000
  }
]
```

```bash
aws s3api put-bucket-cors --bucket phd-portal-dev --cors-configuration file://cors.json
```

#### 3. Create IAM User

```bash
# Create user
aws iam create-user --user-name phd-portal-uploader

# Attach S3 policy
aws iam attach-user-policy --user-name phd-portal-uploader \
  --policy-arn arn:aws:iam::aws:policy/AmazonS3FullAccess

# Create access key
aws iam create-access-key --user-name phd-portal-uploader
```

#### 4. Update `.env`

```bash
S3_ENDPOINT=                # Leave empty for AWS S3
S3_REGION=us-east-1
S3_BUCKET=phd-portal-dev
S3_ACCESS_KEY=AKIAIOSFODNN7EXAMPLE
S3_SECRET_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
S3_USE_PATH_STYLE=false     # AWS S3 uses virtual-hosted style
```

---

## Production Deployment

### 1. Security Checklist

- [ ] **Never commit credentials** to Git
- [ ] Use **IAM roles** (AWS EC2/ECS) or **service accounts** (K8s) instead of access keys
- [ ] Enable **bucket encryption** at rest
- [ ] Enable **bucket versioning** (optional, for backup)
- [ ] Set **bucket policy** to deny public access
- [ ] Enable **CloudTrail** or S3 access logging
- [ ] Use **VPC endpoints** for private S3 access (AWS)
- [ ] Rotate access keys regularly

### 2. AWS S3 Production Setup

#### Create Production Bucket

```bash
aws s3 mb s3://phd-portal-prod --region us-east-1
```

#### Block Public Access

```bash
aws s3api put-public-access-block \
  --bucket phd-portal-prod \
  --public-access-block-configuration \
    "BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true"
```

#### Enable Encryption

```bash
aws s3api put-bucket-encryption \
  --bucket phd-portal-prod \
  --server-side-encryption-configuration \
    '{"Rules":[{"ApplyServerSideEncryptionByDefault":{"SSEAlgorithm":"AES256"}}]}'
```

#### Create IAM Policy (Least Privilege)

`s3-upload-policy.json`:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["s3:PutObject", "s3:GetObject", "s3:DeleteObject"],
      "Resource": "arn:aws:s3:::phd-portal-prod/*"
    },
    {
      "Effect": "Allow",
      "Action": "s3:ListBucket",
      "Resource": "arn:aws:s3:::phd-portal-prod"
    }
  ]
}
```

```bash
aws iam create-policy --policy-name PHDPortalS3Access \
  --policy-document file://s3-upload-policy.json
```

#### Configure CORS

```bash
aws s3api put-bucket-cors --bucket phd-portal-prod \
  --cors-configuration file://cors-production.json
```

`cors-production.json`:

```json
[
  {
    "AllowedOrigins": ["https://phd-portal.kaznmu.kz"],
    "AllowedMethods": ["PUT", "GET", "HEAD"],
    "AllowedHeaders": ["*"],
    "ExposeHeaders": ["ETag"],
    "MaxAgeSeconds": 3000
  }
]
```

### 3. Environment Variables (Production)

```bash
# Use AWS Secrets Manager or environment variables
S3_ENDPOINT=                # Empty for AWS S3
S3_REGION=us-east-1
S3_BUCKET=phd-portal-prod
S3_ACCESS_KEY=<from-secrets-manager>
S3_SECRET_KEY=<from-secrets-manager>
S3_USE_PATH_STYLE=false
S3_PRESIGN_EXPIRES_MINUTES=10    # Shorter for production
S3_MAX_FILE_SIZE_MB=50           # Adjust as needed
```

### 4. CDN Setup (Optional but Recommended)

Use **CloudFront** to cache downloads and reduce S3 costs:

```bash
# Create CloudFront distribution
aws cloudfront create-distribution \
  --origin-domain-name phd-portal-prod.s3.amazonaws.com \
  --default-root-object index.html
```

Update presigned URLs to use CloudFront domain instead of S3 direct.

---

## Troubleshooting

### Problem: "S3 not configured" error

**Cause**: `S3_BUCKET` is empty or credentials are missing.

**Solution**:

```bash
# Check environment variables
echo $S3_BUCKET
echo $S3_ACCESS_KEY
echo $S3_SECRET_KEY

# Restart backend after setting them
```

---

### Problem: "Access Denied" when uploading

**Causes**:

1. **Wrong credentials**: Access key / secret key mismatch
2. **IAM policy**: User doesn't have `s3:PutObject` permission
3. **Bucket policy**: Denies uploads

**Solution**:

```bash
# Test credentials with AWS CLI
aws s3 ls s3://phd-portal --profile phd

# Check IAM policy
aws iam get-user-policy --user-name phd-portal-uploader --policy-name S3Access

# Check bucket policy
aws s3api get-bucket-policy --bucket phd-portal
```

---

### Problem: CORS error in browser

**Error**: `Access to fetch at '...' from origin 'http://localhost:5174' has been blocked by CORS policy`

**Cause**: Bucket CORS configuration doesn't allow your frontend origin.

**Solution** (MinIO):

1. Go to **http://localhost:9001**
2. Select bucket → **Configuration** → **CORS**
3. Add your frontend URL to `AllowedOrigins`

**Solution** (AWS S3):

```bash
# Update CORS
aws s3api put-bucket-cors --bucket phd-portal \
  --cors-configuration file://cors.json
```

---

### Problem: "Presigned URL expired"

**Cause**: URL is valid for only 15 minutes (default). User took too long to upload.

**Solution**:

1. Increase expiry time:

   ```bash
   S3_PRESIGN_EXPIRES_MINUTES=30
   ```

2. Or implement retry logic in frontend:
   ```tsx
   if (uploadResp.status === 403) {
     // Re-request presigned URL
     const newPresign = await presignNodeUpload(...);
     // Retry upload
   }
   ```

---

### Problem: "File too large"

**Error**: `file size 150000000 bytes exceeds maximum 104857600 bytes`

**Cause**: File exceeds `S3_MAX_FILE_SIZE_MB` (default 100MB).

**Solution**:

```bash
# Increase limit (e.g., to 500MB)
S3_MAX_FILE_SIZE_MB=500
```

Or implement **multipart upload** for files > 100MB (advanced).

---

### Problem: Files uploaded but not showing in app

**Cause**: Frontend didn't call `/uploads/attach` endpoint to save metadata.

**Solution**: Check browser console for errors. Ensure this code runs:

```tsx
await attachNodeUpload(nodeId, {
  slot_key: slotKey,
  filename: file.name,
  object_key: presign.object_key,
  content_type: contentType,
  size_bytes: file.size,
  etag,
});
```

---

## Security Best Practices

### 1. Never Expose Credentials

❌ **WRONG**:

```go
// Hardcoded credentials
accessKey := "AKIAIOSFODNN7EXAMPLE"
```

✅ **CORRECT**:

```go
// Read from environment
accessKey := os.Getenv("S3_ACCESS_KEY")
```

---

### 2. Validate All Inputs

✅ **File size**: Prevent resource exhaustion

```go
if err := services.ValidateFileSize(req.SizeBytes); err != nil {
    return err
}
```

✅ **Content type**: Prevent malicious file uploads

```go
if err := services.ValidateContentType(req.ContentType); err != nil {
    return err
}
```

✅ **Filename**: Sanitize to prevent path traversal

```go
cleanName := sanitizeFilename(filename) // Only alphanumeric, -, _, .
```

---

### 3. Use Short-Lived Presigned URLs

✅ Default: 15 minutes

```bash
S3_PRESIGN_EXPIRES_MINUTES=15
```

Longer URLs increase risk if leaked.

---

### 4. Implement Rate Limiting

Prevent abuse of presign endpoint:

```go
// In middleware
if rateLimiter.Exceeded(userID) {
    c.JSON(429, gin.H{"error": "too many requests"})
    return
}
```

---

### 5. Enable Logging

Track all S3 operations:

```go
log.Printf("[S3] User %s uploaded %s (%d bytes)", userID, objectKey, sizeBytes)
```

---

### 6. Scan Uploaded Files

Integrate antivirus scanning:

```go
// After upload, scan file
if infected, err := antivirusService.Scan(objectKey); infected {
    s3c.DeleteObject(objectKey)
    return errors.New("malicious file detected")
}
```

---

## API Reference

### Backend Endpoints

#### `POST /api/journey/nodes/:nodeId/uploads/presign`

**Request**:

```json
{
  "slot_key": "main_document",
  "filename": "thesis.pdf",
  "content_type": "application/pdf",
  "size_bytes": 1048576
}
```

**Response**:

```json
{
  "upload_url": "https://s3.amazonaws.com/phd-portal/nodes/...",
  "object_key": "nodes/user123/node456/main_document/20251115-abc-thesis.pdf",
  "bucket": "phd-portal",
  "expires_in": 900,
  "max_size_bytes": 104857600,
  "required_headers": {
    "Content-Type": "application/pdf"
  }
}
```

---

#### `POST /api/journey/nodes/:nodeId/uploads/attach`

**Request**:

```json
{
  "slot_key": "main_document",
  "filename": "thesis.pdf",
  "object_key": "nodes/user123/...",
  "content_type": "application/pdf",
  "size_bytes": 1048576,
  "etag": "abc123def456"
}
```

**Response**:

```json
{
  "ok": true,
  "document_id": "doc-uuid-123"
}
```

---

#### `GET /api/documents/:docId/presign-get`

**Response**:

```json
{
  "url": "https://s3.amazonaws.com/phd-portal/nodes/...?X-Amz-Expires=900&..."
}
```

---

## Additional Resources

- **AWS S3 Documentation**: https://docs.aws.amazon.com/s3/
- **MinIO Documentation**: https://min.io/docs/
- **AWS SDK for Go v2**: https://aws.github.io/aws-sdk-go-v2/
- **Presigned URLs Guide**: https://docs.aws.amazon.com/AmazonS3/latest/userguide/using-presigned-url.html

---

## FAQ

### Q: Why use presigned URLs instead of uploading through backend?

**A**: Presigned URLs allow **direct browser-to-S3 uploads**, which:

- ✅ Reduces backend load (no file streaming)
- ✅ Faster uploads (parallel, no backend bottleneck)
- ✅ Lower bandwidth costs (files don't go through backend)
- ✅ Better user experience (progress bars, resume uploads)

---

### Q: Can I use Google Cloud Storage or Azure Blob?

**A**: This implementation uses AWS SDK v2 which is S3-specific. To support GCS/Azure:

1. Create separate client implementations
2. Use unified interface:
   ```go
   type StorageClient interface {
       PresignPut(key string) (string, error)
       PresignGet(key string) (string, error)
   }
   ```
3. Select client based on env var `STORAGE_PROVIDER`

---

### Q: How do I migrate from local storage to S3?

**Steps**:

1. Keep local storage code as fallback
2. Upload new files to S3
3. Create migration script:
   ```go
   // For each file in uploads/
   localPath := "backend/uploads/file.pdf"
   objectKey := "nodes/user123/..."
   s3c.UploadFile(localPath, objectKey)
   // Update DB: storage_path → object_key, bucket → phd-portal
   ```
4. Remove local storage code

---

### Q: How to handle file deletion?

**Current**: Files are not automatically deleted from S3 when records are deleted.

**Solution**: Implement lifecycle policy or soft delete:

```go
func DeleteDocument(docID string) error {
    // 1. Get object_key from DB
    var objectKey string
    db.QueryRow("SELECT object_key FROM document_versions WHERE document_id=$1", docID).Scan(&objectKey)

    // 2. Delete from S3
    s3c.DeleteObject(objectKey)

    // 3. Delete from DB
    db.Exec("DELETE FROM document_versions WHERE document_id=$1", docID)

    return nil
}
```

---

## Support

For issues or questions:

1. Check logs: `backend/logs/app.log`
2. Verify S3 configuration: `echo $S3_BUCKET`
3. Test AWS CLI: `aws s3 ls s3://phd-portal`
4. Open GitHub issue with error logs

---

**Last Updated**: November 15, 2025  
**Version**: 1.0.0  
**Maintainer**: PhD Portal Development Team

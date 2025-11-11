# Direct-to-S3 Upload + “Submit → Review → Notify” (Dev: MinIO, Prod: AWS S3)

This spec is written so an AI agent can implement it and a reviewer can verify it. It covers:
- Local dev with MinIO (S3-compatible) using **presigned PUT** uploads and **presigned GET** downloads
- Drop-in API/DB schema for the submit→review→notify workflow
- Frontend usage (React Router v7 + TanStack Query)
- Security, envs, and a test plan

---

## 1) Goals

1) Students upload `.docx/.pdf` directly to object storage using **presigned PUT**.  
2) The backend stores file metadata in Postgres.  
3) Student clicks **Submit for review** → advisors get a notification and see the doc in their **Review queue**.  
4) Advisor approves or requests changes → student gets notified.  
5) Same Go code runs for both **MinIO (dev/on-prem)** and **AWS S3 (prod)** by switching env vars.

---

## 2) Object Storage (Dev vs Prod)

### Dev / On-prem (MinIO)
- Launch MinIO via Docker Compose (below).
- CORS allows your local frontend origins.
- Bucket: `kaznmu`.

### Production (AWS S3)
- Private bucket (e.g., `kaznmu-phd-docs-prod`), versioning ON, SSE encryption ON.
- Set CORS for your production origin.
- Same backend code path; only env changes.

---

## 3) Environment Variables

Create `.env.dev` (MinIO) and `.env` (prod). The app reads these:

```
# shared
MAX_UPLOAD_MB=30

# dev/on-prem MinIO
STORAGE_DRIVER=minio
S3_ENDPOINT=http://localhost:9000
S3_REGION=us-east-1
S3_BUCKET=kaznmu
S3_ACCESS_KEY=minio
S3_SECRET_KEY=miniosecret
S3_USE_PATH_STYLE=true
S3_SSL=false

# production AWS S3 (example)
# STORAGE_DRIVER=s3
# S3_REGION=eu-central-1
# S3_BUCKET=kaznmu-phd-docs-prod
# S3_ACCESS_KEY=***
# S3_SECRET_KEY=***
# S3_USE_PATH_STYLE=false
# S3_SSL=true
```

---

## 4) Database Schema (Postgres)

```sql
-- roles/users assumed (users.id, users.role ∈ {student, advisor, admin, ...})

CREATE TABLE IF NOT EXISTS advisor_assignments (
  id bigserial PRIMARY KEY,
  student_id bigint NOT NULL REFERENCES users(id),
  advisor_id bigint NOT NULL REFERENCES users(id),
  is_active boolean NOT NULL DEFAULT true,
  created_at timestamptz NOT NULL DEFAULT now(),
  closed_at timestamptz,
  UNIQUE (student_id, advisor_id, is_active)
);

DO $$ BEGIN
  CREATE TYPE document_status AS ENUM ('DRAFT','SUBMITTED','IN_REVIEW','CHANGES_REQUESTED','APPROVED');
EXCEPTION WHEN duplicate_object THEN null; END $$;

CREATE TABLE IF NOT EXISTS documents (
  id bigserial PRIMARY KEY,
  student_id bigint NOT NULL REFERENCES users(id),
  document_type text NOT NULL,
  storage_key text NOT NULL,
  filename text NOT NULL,
  mime_type text NOT NULL,
  size_bytes bigint NOT NULL,
  sha256 text,
  status document_status NOT NULL DEFAULT 'DRAFT',
  version int NOT NULL DEFAULT 1,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS document_reviews (
  id bigserial PRIMARY KEY,
  document_id bigint NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
  advisor_id bigint NOT NULL REFERENCES users(id),
  decision text NOT NULL CHECK (decision IN ('APPROVED','CHANGES_REQUESTED')),
  comment text,
  created_at timestamptz NOT NULL DEFAULT now()
);

DO $$ BEGIN
  CREATE TYPE notif_type AS ENUM ('DOC_SUBMITTED','REVIEW_APPROVED','REVIEW_CHANGES');
  CREATE TYPE notif_status AS ENUM ('UNREAD','READ');
EXCEPTION WHEN duplicate_object THEN null; END $$;

CREATE TABLE IF NOT EXISTS notifications (
  id bigserial PRIMARY KEY,
  user_id bigint NOT NULL REFERENCES users(id),
  type notif_type NOT NULL,
  payload jsonb NOT NULL,
  status notif_status NOT NULL DEFAULT 'UNREAD',
  created_at timestamptz NOT NULL DEFAULT now()
);
```

---

## 5) docker-compose.minio.yml

```yaml
version: "3.9"

services:
  minio:
    image: minio/minio:RELEASE.2025-01-20T00-00-00Z
    command: server /data --console-address ":9001"
    ports:
      - "9000:9000"   # S3 API
      - "9001:9001"   # Web console
    environment:
      MINIO_ROOT_USER: minio
      MINIO_ROOT_PASSWORD: miniosecret
    volumes:
      - ./minio-data:/data

  mc:
    image: minio/mc
    depends_on:
      - minio
    entrypoint: >
      /bin/sh -c "
      mc alias set local http://minio:9000 minio miniosecret &&
      mc mb -p local/kaznmu || true &&
      printf '[{"AllowedOrigins":["http://localhost:5173","http://localhost:3000"],"AllowedMethods":["PUT","GET","HEAD"],"AllowedHeaders":["*"],"ExposeHeaders":["ETag","x-amz-request-id"],"MaxAgeSeconds":3600}]' > /tmp/cors.json &&
      mc cors set local/kaznmu /tmp/cors.json &&
      mc ilm add local/kaznmu --expiry-days 365 || true &&
      sleep 3600
      "
```

---

## 6) Test Plan (Dev with MinIO)

1) Start MinIO (`docker compose -f docker-compose.minio.yml up -d`).  
2) Set `.env.dev` and run backend/frontend locally.  
3) As a **student**:
   - Upload `.docx` → expect HTTP 200 on PUT; doc appears in MinIO console; row in `documents`.
   - Click **Submit** → status `SUBMITTED`; advisors receive `DOC_SUBMITTED`.
4) As an **advisor**:
   - Open **Review queue** → see student doc.
   - Click **Open** → backend returns presigned GET; file opens.
   - **Approve** → student gets `REVIEW_APPROVED`.  
     or **Request changes** → student gets `REVIEW_CHANGES` with comment.
5) Re-upload a new version and re-submit → flow repeats.
6) Switch to prod by setting AWS env vars; redeploy; repeat 3–5 on staging/prod.

---

## 7) Security Checklist

- Enforce MIME & extension server-side.
- Size limit both UI and backend (`MAX_UPLOAD_MB`).
- Private bucket only; **no** public ACL.
- Short presigned TTLs (PUT: 10 min, GET: 2–5 min).
- Store size + optional `sha256`; consider HEAD verify on read if needed.
- (Optional) Virus scan pipeline (enqueue `storage_key` to a worker).
- Audit uploads (user_id, ip, user agent).
- Bucket lifecycle rules (expiry for old versions/temp).

---

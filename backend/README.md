# Backend — PhD Student Portal (Go/Gin/sqlx/PostgreSQL)

This backend serves the PhD Student Portal API: auth, users, checklist, documents, comments, and notifications.

## Stack
- Go + Gin (HTTP)
- JWT (6-month expiry by default, configurable with `JWT_EXP_DAYS`)
- sqlx + PostgreSQL
- golang-migrate (DB migrations)
- SMTP (Mailpit for local emails)

## Key Concepts
- **Roles**: `superadmin`, `admin`, `advisor`, `chair`, `student`
- **Soft removal**: `users.is_active=false`
- **Password policy**: human-readable passphrases (three words + two digits) on initial creation; bcrypt stored; never return plaintext except once at creation.
- **Password reset**: `/auth/forgot` issues a single-use token emailed to the user; `/auth/reset` consumes it.

## Routes

```
POST   /api/auth/login         {email,password} -> {token, role}
POST   /api/auth/forgot        {email}          -> ok
POST   /api/auth/reset         {token,new_password} -> ok
PATCH  /api/me/password        {new_password}   -> ok (self)

# Admin
POST   /api/admin/users        {first_name,last_name,email,role} -> {username,temp_password}
PATCH  /api/admin/users/:id/password {new_password} -> ok (not allowed for superadmin)
PATCH  /api/admin/users/:id/active   {active:bool}   -> ok

# Health
GET    /api/health
```

> NOTE: Add auth middleware and RBAC guards for production. The starter wires endpoints; refine RBAC to read role from JWT claims.

## Storage
- Files: `UPLOAD_DIR` (local) or switch to S3-compatible store later.
- Emails: SMTP (Mailpit during dev).

## Development

1. Create `.env` from `.env.example` and set `DATABASE_URL`.
2. Run migrations:
   ```bash
   make migrate-up
   ```
3. Seed checklist (placeholder):
   ```bash
   make seed
   ```
4. Start server:
   ```bash
   make run
   ```

## Password Reset Flow (Mailpit)
- `POST /api/auth/forgot` with `{ "email": "admin@example.com" }`
- Check Mailpit UI at http://localhost:8025 to get the reset link
- Submit `POST /api/auth/reset` with `{ "token": "...", "new_password": "..." }`

## Notes
- Enforce max upload size & mime-type allowlist when you add the upload endpoints.
- Add per-object authorization checks (students only see their data; advisors/chairs only see assigned students).


## S3-Compatible Uploads
Configure the following to enable pre-signed uploads:
```
S3_ENDPOINT=http://localhost:9000       # e.g., MinIO
S3_REGION=us-east-1
S3_BUCKET=phd-portal
S3_ACCESS_KEY=...
S3_SECRET_KEY=...
S3_USE_PATH_STYLE=true
```
Endpoints:
- `POST /api/documents/:docId/presign` → `{ url, object_key }`
- `POST /api/documents/:docId/versions` (local multipart, fallback)

## Comments & Reviews
- Threading: `comments.parent_id`
- Mentions: `comments.mentions uuid[]`
- Add comment: `POST /api/documents/:docId/comments` with `{ body, parent_id?, mentions? }`
- Approve step: `POST /api/reviews/:id/steps/:stepId/approve` (optional `{ comment, mentions }`)
- Return step: `POST /api/reviews/:id/steps/:stepId/return`

## Users Listing (mentions/autocomplete)
- `GET /api/admin/users?q=` returns up to 50 active users (id, name, email, role).


## Route Guards (Server)
Admin routes are protected by JWT + role checks in middleware. Non-admin routes can be guarded as needed.


## /me Endpoint & Caching
- `GET /api/me` (JWT required) returns current user profile.
- Redis caches the `me:{user_id}` payload for 10 minutes to reduce DB calls.

## Logging
- Minimal structured logging helpers in `internal/logging/` (swap to zap/zerolog later).

## Ports
- Server listens on `APP_PORT` (default **8280**). The root `docker-compose.yml` maps `8280:8280`.


## /me endpoint & Context Hydration
- `GET /api/me` returns the current authenticated user.
- After JWT verification the middleware hydrates `current_user` into Gin context using Redis cache to avoid frequent DB hits.

## Logging
- `middleware.RequestLogger()` prints method, path, status, latency per request (structured key=value).

## Redis Caching
- Set `REDIS_URL=redis://redis:6379/0` in `.env` to enable cache (used by `/me` and available for future use).

## Port
- Backend defaults to **8280**. Update your frontend `VITE_API_URL` accordingly.


## File Previews
- `GET /api/documents/:docId/presign-get` → returns a temporary S3 GET URL for the latest version (if S3 configured).
- `GET /api/documents/versions/:versionId/download` → serves a local file for preview/download (dev fallback).


## Documents Listing
- `GET /api/students/:id/documents` → returns student's documents.

## Previews
- `GET /api/documents/:docId/presign-get` (S3) or `GET /api/documents/versions/:versionId/download` (local).

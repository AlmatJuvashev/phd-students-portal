# Production Configuration Guide

This guide helps you configure the PhD Portal for production deployment.

## Critical Configuration Changes

### 1. SMTP Email Configuration

**Current (Development):**

```env
SMTP_HOST=localhost
SMTP_PORT=1027
SMTP_USER=
SMTP_PASS=
```

**⚠️ This uses Mailpit for local testing. Emails won't be sent in production!**

#### Option A: Gmail

1. Enable 2-Factor Authentication in your Google Account
2. Create an App Password: https://myaccount.google.com/apppasswords
3. Update `.env`:

```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-16-character-app-password
SMTP_FROM="PhD Portal <your-email@gmail.com>"
```

#### Option B: SendGrid

1. Create account at https://sendgrid.com
2. Generate API Key in Settings → API Keys
3. Update `.env`:

```env
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USER=apikey
SMTP_PASS=SG.your-sendgrid-api-key-here
SMTP_FROM="PhD Portal <noreply@yourdomain.com>"
```

#### Option C: AWS SES

1. Verify your domain/email in AWS SES Console
2. Create SMTP credentials
3. Update `.env`:

```env
SMTP_HOST=email-smtp.us-east-1.amazonaws.com
SMTP_PORT=587
SMTP_USER=your-aws-smtp-username
SMTP_PASS=your-aws-smtp-password
SMTP_FROM="PhD Portal <noreply@yourdomain.com>"
```

### 2. Redis Configuration

**Current (Development):**

```env
REDIS_ADDR=localhost:6381
REDIS_PASSWORD=
```

#### Production Setup:

**Option A: Redis Cloud (Recommended)**

1. Create account at https://redis.com/try-free/
2. Get connection details
3. Update `.env`:

```env
REDIS_ADDR=redis-12345.c123.us-east-1-1.ec2.cloud.redislabs.com:12345
REDIS_PASSWORD=your-redis-password-here
```

**Option B: Self-hosted Redis**

1. Install Redis: `apt install redis-server` (Ubuntu)
2. Configure password in `/etc/redis/redis.conf`:
   ```
   requirepass your-strong-password
   ```
3. Update `.env`:

```env
REDIS_ADDR=your-server-ip:6379
REDIS_PASSWORD=your-strong-password
```

**Option C: Docker Compose**

```yaml
redis:
  image: redis:7-alpine
  command: redis-server --requirepass yourpassword
  ports:
    - "6379:6379"
```

```env
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=yourpassword
```

### 3. Database Configuration

**Production PostgreSQL:**

```env
DATABASE_URL=postgres://username:password@hostname:5432/dbname?sslmode=require
```

**Important:** Use `sslmode=require` in production for encrypted connections.

### 4. S3 Storage Configuration

**Current (Development - MinIO):**

```env
S3_ENDPOINT=http://localhost:9090
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_USE_PATH_STYLE=true
```

**Production (AWS S3):**

1. Create S3 bucket in AWS Console
2. Create IAM user with S3 access
3. Update `.env`:

```env
S3_ENDPOINT=
S3_REGION=us-east-1
S3_BUCKET=your-bucket-name
S3_ACCESS_KEY=AKIA...
S3_SECRET_KEY=your-secret-key
S3_USE_PATH_STYLE=false
```

### 5. Security Settings

**Change these values:**

```env
JWT_SECRET=generate-random-64-character-string-here
ADMIN_PASSWORD=strong-unique-password-here
APP_ENV=production
FRONTEND_BASE_URL=https://yourdomain.com
```

## Testing Configuration

### Test Email Sending

```bash
# Start the server
cd backend
go run cmd/server/main.go

# In another terminal, trigger a test email by changing document state
# Check server logs for "Email sent to..." or error messages
```

### Test Redis Connection

```bash
# Test connection manually
redis-cli -h your-redis-host -p 6379 -a your-password ping
# Should return: PONG

# Check server logs on startup:
# ✅ Should NOT see: "Redis connection failed (debouncing disabled)"
```

## Environment Variables Checklist

Before deploying to production, verify:

- [ ] `APP_ENV=production`
- [ ] `JWT_SECRET` changed from default
- [ ] `ADMIN_PASSWORD` changed from default
- [ ] `DATABASE_URL` points to production database with `sslmode=require`
- [ ] `SMTP_HOST` is real SMTP server (not localhost:1027)
- [ ] `SMTP_USER` and `SMTP_PASS` are set
- [ ] `SMTP_FROM` uses your real domain
- [ ] `FRONTEND_BASE_URL` is your production domain (https)
- [ ] `REDIS_ADDR` points to production Redis
- [ ] `REDIS_PASSWORD` is set (if required)
- [ ] `S3_ENDPOINT` is empty (for AWS) or your S3-compatible endpoint
- [ ] `S3_ACCESS_KEY` and `S3_SECRET_KEY` are production credentials

## Common Issues

### Emails not sending

- Check SMTP credentials
- Verify firewall allows outbound port 587/465
- Check server logs for detailed error messages
- Test SMTP connection: `telnet smtp.gmail.com 587`

### Redis connection failed

- Verify Redis is running: `redis-cli ping`
- Check REDIS_ADDR format (host:port)
- Verify REDIS_PASSWORD if authentication is enabled
- System will fall back to allowing all notifications if Redis fails

### S3 upload errors

- Verify bucket exists and is in correct region
- Check IAM permissions include `s3:PutObject`, `s3:GetObject`, `s3:DeleteObject`
- For MinIO, ensure `S3_USE_PATH_STYLE=true`
- For AWS S3, ensure `S3_USE_PATH_STYLE=false`

## Support

For deployment assistance, contact the development team or refer to:

- SMTP issues: Check provider documentation
- Redis issues: https://redis.io/docs/
- AWS S3 issues: https://docs.aws.amazon.com/s3/

# Testing Guide - Document Upload & Review System

## Prerequisites

Before testing, ensure these services are running:

### 1. PostgreSQL Database
```bash
# Check if running
pg_isready -h localhost -p 5435

# If not running (with Docker):
docker-compose up -d postgres
```

### 2. Redis (for debouncing)
```bash
# Check if running
redis-cli -p 6381 ping
# Expected: PONG

# If not running (with Docker):
docker run -d -p 6381:6379 --name phd-redis redis:7-alpine
```

### 3. MinIO (S3-compatible storage)
```bash
# Check if running
curl http://localhost:9090/minio/health/live
# Expected: 200 OK

# If not running (with Docker):
docker run -d \
  -p 9090:9000 \
  -p 9091:9001 \
  --name phd-minio \
  -e "MINIO_ROOT_USER=minioadmin" \
  -e "MINIO_ROOT_PASSWORD=minioadmin" \
  minio/minio server /data --console-address ":9001"

# Access MinIO Console: http://localhost:9091
# Login: minioadmin / minioadmin
# Create bucket: phd-portal
```

### 4. Mailpit (Email testing)
```bash
# Check if running
curl http://localhost:8025
# Expected: HTML response

# If not running (with Docker):
docker run -d \
  -p 1027:1025 \
  -p 8025:8025 \
  --name phd-mailpit \
  axllent/mailpit

# Access Mailpit UI: http://localhost:8025
```

## Starting the Application

### Backend
```bash
cd backend
go run cmd/server/main.go
```

Expected output:
```
INFO: API listening on port 8280
S3 cleanup worker started
```

### Frontend
```bash
cd frontend
npm run dev
```

Expected output:
```
VITE ready in 500ms
Local: http://localhost:5173
```

## Testing Scenarios

### Scenario 1: Email Notifications ✉️

**Test that students receive emails when advisor reviews their work**

1. **Login as Student**
   - Go to http://localhost:5173
   - Login with test student credentials

2. **Upload Document**
   - Navigate to a node that requires document upload
   - Click "Upload file"
   - Select a PDF or image file
   - Submit

3. **Login as Advisor/Admin**
   - Logout and login with admin credentials
   - Find the student's submission in notifications

4. **Review Document**
   - Open the submission
   - Change state to "Approved" or "Changes Requested"
   - Add feedback note (if requesting changes)
   - Submit review

5. **Check Email**
   - Open Mailpit: http://localhost:8025
   - You should see an email sent to the student
   - Verify email contains:
     - Student name
     - Node ID
     - New status
     - Link to document

**Expected Result:** ✅ Email appears in Mailpit within seconds

**Troubleshooting:**
- Check backend logs for "Email sent to..."
- If you see "SMTP not configured", verify `.env` has SMTP_HOST and SMTP_PORT
- Mailpit should be accessible at http://localhost:8025

---

### Scenario 2: Timeline View 📜

**Test that document history displays chronologically**

1. **Upload Multiple Versions**
   - Login as student
   - Upload document version 1
   - Wait a few seconds
   - Upload document version 2
   - Upload document version 3

2. **Verify Timeline**
   - Open the node details page
   - Scroll to "Supporting documents" section
   - You should see:
     - Vertical timeline with gradient line
     - Blue circles for student uploads (oldest at top)
     - Version numbers (1, 2, 3)
     - Timestamps
     - File sizes

3. **Test Advisor Review**
   - Login as advisor
   - Upload reviewed document with comments
   - Logout and login as student again

4. **Verify Timeline Update**
   - Timeline should now show:
     - Green circle for advisor's reviewed document
     - Positioned chronologically after student's upload
     - Advisor name in metadata

**Expected Result:** ✅ Clear conversation-style timeline showing who did what and when

---

### Scenario 3: Document Preview 👁️

**Test in-browser file preview**

1. **Upload Different File Types**
   - Upload a PDF file
   - Upload a JPG/PNG image
   - Upload a DOCX file (optional)

2. **Test PDF Preview**
   - Click the eye icon 👁️ next to PDF file
   - Drawer should slide in from right
   - PDF should display in iframe
   - Click X to close

3. **Test Image Preview**
   - Click eye icon on image file
   - Image should display centered in drawer
   - Should be zoomable/scrollable

4. **Test DOCX Preview**
   - Click eye icon on DOCX
   - Should show warning: "DOCX preview requires publicly accessible URL"
   - This is expected for localhost

**Expected Result:** ✅ PDF and images preview without download

**Note:** DOCX preview only works in production with public domain

---

### Scenario 4: NEW Badges 🔴

**Test unread file indicators**

1. **Setup**
   - Login as student
   - View a node with documents
   - Note the current time

2. **Trigger New Upload (as advisor)**
   - Login as advisor
   - Upload reviewed document
   - Logout

3. **Check Badge**
   - Login as student again
   - Navigate to the same node
   - New advisor upload should have red "NEW" badge
   - Badge should be animated (pulsing)

4. **Mark as Read**
   - Stay on the page for a few seconds
   - Navigate away and come back
   - "NEW" badge should disappear

**Expected Result:** ✅ NEW badge appears only for files uploaded after last view

**How it works:** Uses localStorage to track `last_viewed_node_{nodeId}` timestamp

---

### Scenario 5: S3 Cleanup Worker 🧹

**Test orphaned file deletion**

1. **Create Orphaned File**
   - Manually upload a file to MinIO Console (http://localhost:9091)
   - Bucket: phd-portal
   - Upload any file with name like: `orphan-test-file.pdf`

2. **Wait 24+ Hours** (or modify worker code to run every 1 minute for testing)

3. **Check Cleanup**
   - Check backend logs for: "Found X orphans, deleted Y files"
   - Verify orphan file is deleted from MinIO

**For Quick Testing:**
Modify `backend/internal/worker/cleanup.go`:
```go
ticker := time.NewTicker(1 * time.Minute) // Change from 24 hours
```

**Expected Result:** ✅ Files not in database are deleted after 24 hours

---

### Scenario 6: Notification Debouncing 🔕

**Test that rapid uploads don't spam notifications**

1. **Upload 3 Files Quickly**
   - Login as student
   - Upload file v1
   - Immediately upload file v2
   - Immediately upload file v3

2. **Check Notifications**
   - Login as admin
   - Check notification page

3. **Verify Debouncing**
   - You should see only 1 or 2 notifications (not 3)
   - Check backend logs for: "Redis: notification debounced"

**Note:** This requires Redis to be running. If Redis is down, all notifications will be sent (fallback behavior).

**Expected Result:** ✅ Multiple uploads within 10 minutes = 1 notification

---

## Verification Checklist

After testing, verify:

- [ ] Emails appear in Mailpit (http://localhost:8025)
- [ ] Timeline shows chronological events with colors
- [ ] PDF preview works in drawer
- [ ] Image preview works in drawer
- [ ] NEW badges appear on recent uploads
- [ ] Backend logs show "S3 cleanup worker started"
- [ ] Redis connection successful (no "Redis connection failed" in logs)
- [ ] No compilation errors in frontend or backend

## Common Issues

### "SMTP not configured" in logs
- Check `.env` file has SMTP_HOST and SMTP_PORT
- Verify Mailpit is running on port 1027

### "Redis connection failed"
- Start Redis: `docker run -d -p 6381:6379 redis:7-alpine`
- System will work without Redis (debouncing disabled)

### Timeline not showing
- Check browser console for errors
- Verify API returns `slots` data with `attachments`

### Preview drawer empty
- Check that presigned URLs are being generated
- Verify S3_ENDPOINT in `.env`

## Success Criteria

✅ **System is ready for production if:**
1. All 6 scenarios pass
2. No errors in backend logs
3. No console errors in frontend
4. Emails visible in Mailpit
5. Files downloadable and previewable

## Next Steps

After local testing succeeds:
1. Review `docs/PRODUCTION_CONFIG.md`
2. Configure production SMTP (Gmail/SendGrid/SES)
3. Set up production Redis
4. Deploy to staging environment
5. Repeat all tests in staging

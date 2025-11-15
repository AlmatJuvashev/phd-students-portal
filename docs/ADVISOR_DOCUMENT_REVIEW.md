# Advisor Document Review with Comments

## Overview

Advisors can now upload reviewed documents with comments as an optional step in the document review workflow. This enables a clear correspondence trail between students and advisors.

## Features

- ✅ Upload documents with embedded comments/annotations
- ✅ Tracked history of who reviewed and when
- ✅ Student can see both original and reviewed versions
- ✅ Download reviewed documents with comments
- ✅ Permission verification (only assigned advisors)

## Workflow

### 1. Student Uploads Document

Student submits their document (e.g., dissertation chapter, publication) through the Antiplagiarism node or any other node with file upload slots.

### 2. Advisor Reviews Document

**Option A: Approve/Reject without uploading**

- Advisor reviews the document
- Chooses "Одобрить" (Approve) or "Запросить правки" (Request Fixes)
- Adds optional text note/feedback

**Option B: Upload Reviewed Document with Comments**

- Advisor downloads student's document
- Adds comments/annotations using Word, PDF editor, etc.
- Uploads the reviewed document via new endpoint
- Student receives document with inline comments

### 3. Student Views Feedback

Student sees:

- Original submitted document
- Advisor's text feedback (if any)
- **NEW:** Reviewed document with comments (if uploaded)
- Who reviewed and when
- Status (approved/rejected/submitted)

## API Endpoints

### Upload Reviewed Document

```http
POST /api/admin/attachments/:attachmentId/reviewed-document
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "document_version_id": "uuid-of-uploaded-document-version"
}
```

**Response:**

```json
{
  "ok": true,
  "reviewed_document_version_id": "abc-123-...",
  "reviewed_at": "2025-11-15T14:00:00Z"
}
```

### Get Student Files (includes reviewed documents)

```http
GET /api/admin/students/:studentId/nodes/:nodeId/files
Authorization: Bearer <jwt_token>
```

**Response:**

```json
[
  {
    "attachment_id": "...",
    "filename": "dissertation_chapter1.docx",
    "status": "rejected",
    "review_note": "Please address the methodology section",
    "download_url": "/api/documents/versions/.../download",
    "reviewed_document": {
      "version_id": "...",
      "download_url": "/api/documents/versions/.../download",
      "mime_type": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
      "reviewed_by": "Prof. Smith",
      "reviewed_at": "2025-11-15T14:00:00Z"
    }
  }
]
```

## Database Schema

### Migration 0014: advisor_reviewed_document

Added columns to `node_instance_slot_attachments`:

- `reviewed_document_version_id` (uuid, nullable) - Points to document_versions
- `reviewed_by` (uuid, nullable) - User who uploaded the reviewed document
- `reviewed_at` (timestamptz, nullable) - When the document was uploaded

Index: `idx_slot_attachments_reviewed_doc` for fast queries.

## Frontend Integration

### Displaying Reviewed Documents

```typescript
// Check if reviewed document exists
if (file.reviewed_document) {
  const { download_url, reviewed_by, reviewed_at, mime_type } = file.reviewed_document;

  // Show download button
  <a href={download_url} download>
    📄 Download Reviewed Document (with comments)
  </a>

  // Show metadata
  <p>Reviewed by: {reviewed_by}</p>
  <p>Reviewed at: {formatDate(reviewed_at)}</p>
}
```

### Uploading Reviewed Document

```typescript
// Step 1: Upload document to S3 (same as regular upload)
const presignResponse = await fetch("/api/documents/:docId/presign", {
  method: "POST",
  body: JSON.stringify({
    content_type: file.type,
    size_bytes: file.size,
  }),
});

const { presigned_url, version_id } = await presignResponse.json();

// Step 2: Upload to S3
await fetch(presigned_url, {
  method: "PUT",
  body: file,
  headers: { "Content-Type": file.type },
});

// Step 3: Attach as reviewed document
await fetch(`/api/admin/attachments/${attachmentId}/reviewed-document`, {
  method: "POST",
  headers: {
    Authorization: `Bearer ${token}`,
    "Content-Type": "application/json",
  },
  body: JSON.stringify({ document_version_id: version_id }),
});
```

## Permissions

- **Admin**: Can upload reviewed documents for any student
- **Advisor**: Can upload only for assigned students (verified via `student_advisors` table)
- **Student**: Cannot upload reviewed documents (read-only access)

## Event Tracking

When a reviewed document is uploaded, a `reviewed_document_uploaded` event is logged in `node_events`:

```json
{
  "event_type": "reviewed_document_uploaded",
  "payload": {
    "attachment_id": "...",
    "reviewed_version_id": "..."
  },
  "actor_id": "...",
  "created_at": "..."
}
```

## Use Cases

### 1. Dissertation Chapter Review

- Student uploads dissertation chapter
- Advisor downloads, adds Track Changes comments in Word
- Advisor uploads reviewed document
- Student sees exactly what needs to be changed

### 2. Publication Review

- Student uploads publication draft
- Advisor annotates PDF with comments
- Advisor uploads annotated PDF
- Student addresses specific comments

### 3. Antiplagiarism Report

- Student uploads work for plagiarism check
- Advisor highlights problematic sections
- Advisor uploads marked-up version
- Student can see exact areas of concern

## Migration Guide

### Apply Migration

```bash
cd backend
make migrate-up
```

### Rollback (if needed)

```bash
make migrate-down
```

## Testing

### Manual Test Flow

1. Login as student, upload document to node
2. Login as advisor
3. Navigate to student details page
4. Review document (approve/reject)
5. Upload reviewed document with comments
6. Verify student can see both documents
7. Verify download URLs work for both

### API Test

```bash
# Get attachment ID from student files endpoint
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8280/api/admin/students/$STUDENT_ID/nodes/S1_antiplag/files

# Upload reviewed document
curl -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"document_version_id": "$VERSION_ID"}' \
  http://localhost:8280/api/admin/attachments/$ATTACHMENT_ID/reviewed-document
```

## Future Enhancements

- [ ] Version history for reviewed documents
- [ ] Inline preview of reviewed documents
- [ ] Notification to student when reviewed document is uploaded
- [ ] Comparison view (original vs reviewed side-by-side)
- [ ] Support for multiple review rounds

## Related Documentation

- [S3 Storage Setup](./S3_STORAGE.md)
- [File Upload Testing](./TESTING_FILE_UPLOAD.md)
- [Admin Dashboard Guide](./dashboard_instruction.md)

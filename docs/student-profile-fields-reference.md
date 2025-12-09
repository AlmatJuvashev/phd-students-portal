# Student Profile Fields and Database Schema Reference

## Overview

This document provides a comprehensive reference for all student profile fields and database tables used in the PhD Student Portal, particularly for template prefilling.

---

## Profile Data Sources

The student profile data comes from three main sources:

1. **`users` table** - Basic user information
2. **`profile_submissions` table** - Custom profile form data (S1_profile node)
3. **Publications data** - From S1_publications_list node

---

## Database Schema

### 1. `users` Table

**Core User Information**

| Column | Type | Description | Template Field |
|--------|------|-------------|----------------|
| `id` | UUID | User unique identifier | - |
| `username` | text | Username | - |
| `email` | text | Email address | `student_email` |
| `first_name` | text | First name | Part of `student_full_name` |
| `last_name` | text | Last name | Part of `student_full_name` |
| `role` | user_role | User role (student/advisor/admin) | - |
| `phone` | text | Phone number | `student_phone` |
| `program` | text | PhD program name | `student_program` |
| `department` | text | Department/Кафедра | `student_department` |
| `cohort` | text | Cohort/Year of admission | - |
| `is_active` | boolean | Account status | - |
| `created_at` | timestamptz | Creation timestamp | - |
| `updated_at` | timestamptz | Last update timestamp | - |

**Notes:**
- Email is optional (allows students without email)
- `student_full_name` is typically `first_name + " " + last_name`

---

### 2. `profile_submissions` Table

**Custom Profile Form Data (S1_profile node)**

| Column | Type | Description |
|--------|------|-------------|
| `user_id` | UUID | References `users(id)` |
| `form_data` | JSONB | All profile form fields as JSON |
| `submitted_at` | timestamptz | Submission timestamp |
| `updated_at` | timestamptz | Last update timestamp |

**JSONB Structure (`form_data`):**

The `form_data` column stores a flexible JSON object with fields defined in the S1_profile node form schema. Common fields include:

```json
{
  "specialty": "6D110400 — Менеджмент",
  "dissertation_topic": "...",
  "supervisors": ["Supervisor 1", "Supervisor 2"],
  "iin": "123456789012",
  "birth_year": "1995",
  // ... other custom fields
}
```

---

### 3. `student_advisors` Table

**Student-Advisor Relationship (Many-to-Many)**

| Column | Type | Description |
|--------|------|-------------|
| `student_id` | UUID | References `users(id)` | 
| `advisor_id` | UUID | References `users(id)` |

**Usage:**
- Links students to their scientific advisors/supervisors
- Used to populate `student_supervisors` field
- Fetches advisor names from the `users` table

---

### 4. Publications Data (S1_publications_list node)

**Source:** `node_instances` + `node_instance_form_revisions` tables

**Structure:**
The publications data is stored in `form_data` JSONB column of the latest revision:

```json
{
  "wos_scopus": [
    {
      "title": "...",
      "authors": "...",
      "journal": "...",
      "year": "2023",
      "doi": "10.1234/example",
      "issn_print": "1234-5678",
      "issn_online": "8765-4321",
      "volume_issue": "Vol 10, Issue 3",
      "vol": "10"
    }
  ],
  "kokson": [...],
  "conferences": [...],
  "ip": [...]
}
```

---

## Template Prefill Fields

### Standard Fields (Available in Most Templates)

| Template Field | Source | Description |
|----------------|--------|-------------|
| `student_full_name` | `users.first_name + users.last_name` | Full student name |
| `student_email` | `users.email` | Student email |
| `student_phone` | `users.phone` | Phone number |
| `student_program` | `users.program` | PhD program |
| `student_department` | `users.department` | Department/Кафедра |
| `student_specialty` | `profile_submissions.form_data.specialty` | Specialty code and name |
| `dissertation_topic` | `profile_submissions.form_data.dissertation_topic` | Dissertation topic |
| `student_supervisors` | From `student_advisors` + `users` | Comma-separated advisor names |
| `iin` | `profile_submissions.form_data.iin` | Individual Identification Number (IIN) |
| `birth_year` | `profile_submissions.form_data.birth_year` | Year of birth |

### Date Fields

| Template Field | Source | Description |
|----------------|--------|-------------|
| `day` | Computed from current date | Day of the month |
| `month` | Computed from current date | Month name (локalized) |
| `year` | Computed from current date | Full year |
| `submission_date` | Computed | Formatted submission date |

### Publications Fields (App7 Template)

| Template Field | Source | Description |
|----------------|--------|-------------|
| `publications` | S1_publications_list node | Array of publication objects |
| `publications[].no` | Computed | Publication number (1, 2, 3...) |
| `publications[].title` | Publication data | Article title |
| `publications[].authors` | Publication data | Author list |
| `publications[].journal` | Publication data | Journal name |
| `publications[].year` | Publication data | Publication year |
| `publications[].doi` | Publication data | DOI identifier |
| `publications[].issn_print` | Publication data | Print ISSN |
| `publications[].issn_online` | Publication data | Online ISSN |
| `publications[].volume_issue` | Publication data | Volume and issue |
| `publications[].vol` | Publication data | Volume number |

---

## Backend Data Aggregation

### `GetProfile` Function

**Location:** `backend/internal/handlers/node_submission.go`

**Purpose:** Aggregates all profile data for template prefilling

**Data Flow:**
1. Fetches `profile_submissions.form_data` for the user
2. Fetches publications from `S1_publications_list` node
3. Merges both into a single JSON object
4. Returns combined profile data

**Example Response:**
```json
{
  "specialty": "6D110400 — Менеджмент",
  "dissertation_topic": "...",
  "supervisors": ["..."],
  "iin": "123456789012",
  "birth_year": "1995",
  "wos_scopus": [...],
  "kokson": [...],
  "conferences": [...],
  "ip": [...]
}
```

---

## Frontend Data Builder

### `buildTemplateData` Function

**Location:** `frontend/src/features/docgen/student-template.ts`

**Purpose:** Transforms raw profile data + user data into template-ready format

**Input:**
- User object (from auth context)
- Profile data (from GetProfile API)
- Language/locale

**Output:** `StudentTemplateData` object

**Type Definition:**
```typescript
export type StudentTemplateData = {
  student_full_name: string;
  student_program: string;
  student_specialty: string;
  student_supervisors: string;
  submission_date: string;
  student_email: string;
  student_phone: string;
  dissertation_topic: string;
  student_department: string;
  day: string;
  month: string;
  year: string;
  publications?: Array<{
    no: string;
    title: string;
    authors: string;
    journal: string;
    volume_issue: string;
    vol: string;
    year: string;
    issn_print: string;
    issn_online: string;
    doi: string;
  }>;
};
```

---

## Field Mapping Strategy

### Direct Mapping (from `users` table)
- `first_name`, `last_name` → `student_full_name`
- `email` → `student_email`
- `phone` → `student_phone`
- `program` → `student_program`
- `department` → `student_department`

### Profile Form Mapping (from `profile_submissions.form_data`)
- `specialty` → `student_specialty`
- `dissertation_topic` → `dissertation_topic`
- `iin` → `iin`
- `birth_year` → `birth_year`

### Computed Fields
- `student_supervisors`: Joined from `student_advisors` table
- `day`, `month`, `year`: Extracted from current date
- `submission_date`: Formatted current date
- `publications[].no`: Sequential numbering

---

## Custom Form Fields (S1_profile)

The S1_profile node can have custom fields defined in its form schema. To add new fields:

1. **Backend:** Fields are automatically stored in `profile_submissions.form_data` JSONB
2. **Frontend:** Add field to form schema in playbook definition
3. **Templates:** Reference the field using `[% field_name %]` in DOCX templates
4. **Data Builder:** Update `buildTemplateData` to extract and format the field

---

## Adding New Template Fields

### Step 1: Add to Database (if needed)
If the field is not already in `users` or `profile_submissions`:
- Add column to `users` table via migration
- OR add field to S1_profile form schema (stored in JSONB)

### Step 2: Update Backend
Update `GetProfile` function if aggregation logic is needed

### Step 3: Update Frontend Type
Add field to `StudentTemplateData` type in `student-template.ts`

### Step 4: Update Data Builder
Add extraction/formatting logic in `buildTemplateData` function

### Step 5: Use in Template
Reference field as `[% field_name %]` in DOCX template

---

## Example: Adding "Student ID" Field

### 1. Database Migration
```sql
ALTER TABLE users ADD COLUMN student_id text;
```

### 2. Frontend Type
```typescript
export type StudentTemplateData = {
  // ... existing fields
  student_id: string;
};
```

### 3. Data Builder
```typescript
export function buildTemplateData(user, profileData, lang) {
  return {
    // ... existing fields
    student_id: user.student_id || "",
  };
}
```

### 4. Template Usage
```
Студенческий билет: [% student_id %]
```

---

## Query Examples

### Get Student Full Profile
```sql
SELECT 
  u.first_name,
  u.last_name,
  u.email,
  u.phone,
  u.program,
  u.department,
  p.form_data
FROM users u
LEFT JOIN profile_submissions p ON p.user_id = u.id
WHERE u.id = $1;
```

### Get Student Advisors
```sql
SELECT 
  a.first_name,
  a.last_name
FROM student_advisors sa
JOIN users a ON a.id = sa.advisor_id
WHERE sa.student_id = $1;
```

### Get Publications
```sql
SELECT r.form_data 
FROM node_instances i
JOIN node_instance_form_revisions r 
  ON r.node_instance_id = i.id 
  AND r.rev = i.current_rev
WHERE i.user_id = $1 
  AND i.node_id = 'S1_publications_list';
```

---

**Last Updated:** 2025-11-28

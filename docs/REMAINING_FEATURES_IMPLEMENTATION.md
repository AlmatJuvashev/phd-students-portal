# Remaining Features Implementation Guide

> **Document Version:** 2.0  
> **Created:** January 4, 2026  
> **Last Updated:** January 4, 2026  
> **Purpose:** Comprehensive implementation guide for remaining features to achieve best-in-class education platform

---

## Executive Summary

Based on the complete audit (January 4, 2026):

- **Backend:** ~92% implemented (32 services, 35+ handlers)
- **Frontend:** ~75% implemented (22 modules fully working, 4 partial)
- **Overall Integration:** ~78% complete

This document provides detailed implementation plans for the remaining **~22%** organized into 4 phases.

### Current Implementation Status

| Category         | Implemented | Partial | Missing | Total |
| ---------------- | ----------- | ------- | ------- | ----- |
| Backend Services | 32          | 0       | 2       | 34    |
| Backend Handlers | 35+         | 2       | 3       | 40    |
| Frontend Modules | 22          | 4       | 5       | 31    |
| API Integration  | 85%         | 10%     | 5%      | 100%  |

### Implementation Priority Matrix

| Priority    | Features                                                         | Complexity  | Timeline  |
| ----------- | ---------------------------------------------------------------- | ----------- | --------- |
| 🔴 Critical | Student Course/Assignment Detail Pages                           | High        | 1-2 weeks |
| 🟠 High     | Forums UI, Attendance UI, Teacher Student Tracker                | Medium-High | 2-3 weeks |
| 🟡 Medium   | Rich Editor, Quiz Builder Enhancement, Transcript UI             | Medium      | 2-3 weeks |
| 🟢 Low      | Gamification (Full Stack), AI Tools UI, LTI Admin, Global Search | Variable    | 3-4 weeks |

---

## ✅ Already Implemented (Reference)

### Frontend Modules (Fully Working)

| Module                | Key Files                                                              | Routes                          |
| --------------------- | ---------------------------------------------------------------------- | ------------------------------- |
| `assessments/`        | api.ts, types.ts, AssessmentTaking.tsx, AssessmentResults.tsx          | `/student/assessments/:id/take` |
| `student-portal/`     | api.ts, types.ts, Dashboard, Courses, Assignments, Grades              | `/student/*`                    |
| `teacher/`            | api.ts, types.ts, Dashboard, Courses, CourseDetail, Grading            | `/admin/teacher/*`              |
| `studio/`             | api.ts, types.ts, CourseBuilder, ProgramJourneyBuilder + 10 components | `/admin/studio/*`               |
| `item-bank/`          | api.ts, types.ts, BanksPage, BankItemsPage                             | `/admin/item-banks/*`           |
| `curriculum/`         | api.ts, types.ts, Programs, Courses, ProgramDetail                     | `/admin/programs/*`             |
| `grading/`            | api.ts, types.ts                                                       | (API only, no dedicated pages)  |
| `enrollments/`        | api.ts, types.ts, EnrollmentsPage                                      | `/admin/enrollments`            |
| `calendar/`           | Full module with components                                            | `/calendar`, `/admin/calendar`  |
| `chat/`               | Full module with rooms, messages                                       | Chat integration                |
| `analytics/`          | AnalyticsDashboard                                                     | `/admin/analytics`              |
| `admin/dictionaries/` | Full CRUD for all dictionaries                                         | `/admin/dictionaries`           |
| `superadmin/`         | Tenants, Admins, Logs, Settings, Services                              | `/superadmin/*`                 |
| `journey/`            | Full PhD journey tracking                                              | `/journey`                      |
| `students-monitor/`   | Full monitoring with pages                                             | `/admin/students-monitor`       |
| `scheduler/`          | Full scheduling module                                                 | `/admin/scheduler`              |
| `profile/`            | Profile editing, avatar                                                | `/profile`                      |

### Backend Services (All Operational)

All 32+ services working: `admin`, `ai`, `analytics`, `assessment`, `attendance`, `audit`, `auth`, `authz`, `bulk`, `calendar`, `chat`, `checklist`, `comment`, `contact`, `course_content`, `curriculum`, `dictionary`, `document`, `email`, `forum`, `governance`, `grading`, `item_bank`, `journey`, `lti`, `notification`, `program_builder`, `resource`, `rubric`, `s3`, `scheduler`, `search`, `student`, `superadmin`, `teacher`, `tenant`, `transcript`, `user`.

### Frontend API Clients (Implemented)

```typescript
// assessments/api.ts - ✅ COMPLETE
startAttempt,
  getAssessmentForTaking,
  getAttemptDetails,
  submitResponse,
  completeAttempt,
  listMyAttempts;

// student-portal/api.ts - ✅ COMPLETE (basic)
getStudentDashboard, getStudentCourses, getStudentAssignments, getStudentGrades;

// teacher/api.ts - ✅ COMPLETE
getTeacherDashboard,
  getTeacherCourses,
  getCourseRoster,
  getCourseGradebook,
  getTeacherSubmissions,
  submitGradeForSubmission;

// item-bank/api.ts - ✅ COMPLETE
listBanks,
  createBank,
  updateBank,
  deleteBank,
  listQuestions,
  createQuestion,
  updateQuestion,
  deleteQuestion;

// grading/api.ts - ✅ COMPLETE
listGradingSchemas, createGradingSchema, submitGrade, listStudentGrades;

// studio/api.ts - ✅ COMPLETE
getCourseContent,
  updateCourseContent,
  addModule,
  updateModule,
  deleteModule,
  addLesson,
  addActivity,
  updateActivity;

// curriculum/api.ts - ✅ COMPLETE
getPrograms,
  getProgramVersionMap,
  getProgramVersionNodes,
  createProgramVersionNode,
  updateProgramVersionNode,
  updateProgramVersionMap,
  getProgram,
  createProgram,
  updateProgram,
  deleteProgram,
  getCourses,
  getCourse,
  createCourse,
  updateCourse,
  deleteCourse;
```

---

## Phase 1: Student Experience & Assessment Foundation (Critical)

**Timeline:** 2-3 weeks  
**Goal:** Enable students to view course content and take assessments

### 1.1 Student Course Detail Page

**Current State:** ❌ Missing  
**Location:** `frontend/src/features/student-portal/StudentCourseDetail.tsx`

#### Backend Requirements

The following endpoints may need creation or verification:

```
GET /api/student/courses/:id                 # Course details for student
GET /api/student/courses/:id/modules         # Course modules/lessons
GET /api/student/courses/:id/announcements   # Course announcements
GET /api/student/courses/:id/resources       # Downloadable resources
```

**Check existing:** `course_content_handler.go` may have some of these.

#### Frontend Implementation

```
frontend/src/features/student-portal/
├── StudentCourseDetail.tsx          # Main course view
├── components/
│   ├── CourseModules.tsx            # Accordion of modules
│   ├── ModuleContent.tsx            # Video/Reading/Activity
│   ├── CourseAnnouncements.tsx      # Announcements list
│   ├── CourseResources.tsx          # Downloadable files
│   └── CourseProgress.tsx           # Progress bar
```

#### Key Features

1. **Course Header** - Title, instructor, progress %
2. **Modules Accordion** - Expandable sections with lessons
3. **Content Types:**
   - 📹 Video (embedded player or external link)
   - 📖 Reading (markdown/HTML content)
   - 📝 Activity (link to assessment/assignment)
   - 📎 Resource (downloadable file)
4. **Progress Tracking** - Mark items complete
5. **Announcements** - Course-specific announcements

#### UI Mockup Structure

```tsx
<StudentCourseDetail>
  <CourseHeader>
    <BackButton />
    <CourseTitle />
    <InstructorInfo />
    <ProgressBar percent={75} />
  </CourseHeader>

  <Tabs defaultValue="content">
    <Tab value="content">
      <CourseModules>
        <ModuleAccordion>
          <Lesson type="video" />
          <Lesson type="reading" />
          <Lesson type="activity" />
        </ModuleAccordion>
      </CourseModules>
    </Tab>
    <Tab value="announcements">
      <CourseAnnouncements />
    </Tab>
    <Tab value="resources">
      <CourseResources />
    </Tab>
    <Tab value="grades">
      <CourseGrades />
    </Tab>
  </Tabs>
</StudentCourseDetail>
```

---

### 1.2 Student Assignment Detail Page

**Current State:** ❌ Missing  
**Location:** `frontend/src/features/student-portal/StudentAssignmentDetail.tsx`

#### Backend Requirements

```
GET  /api/student/assignments/:id           # Assignment details
POST /api/student/assignments/:id/submit    # Submit assignment
GET  /api/student/assignments/:id/submission # Get my submission
```

**Existing:** `activity_submissions` table, may need new handler.

#### Frontend Implementation

```
frontend/src/features/student-portal/
├── StudentAssignmentDetail.tsx
├── components/
│   ├── AssignmentInstructions.tsx   # Rich text instructions
│   ├── AssignmentRubric.tsx         # Grading rubric display
│   ├── SubmissionUploader.tsx       # File upload
│   ├── TextSubmission.tsx           # Text/essay input
│   └── SubmissionStatus.tsx         # Status badge
```

#### Key Features

1. **Assignment Info:**

   - Title, description, due date
   - Points possible
   - Submission type (file, text, quiz)
   - Attempts allowed/remaining

2. **Rubric Display:**

   - Criteria and point values
   - Performance levels

3. **Submission Area:**

   - File upload (drag & drop)
   - Rich text editor for essays
   - Link submission option

4. **Submission History:**
   - Previous attempts
   - Grades and feedback

#### Workflow States

```
┌─────────────┐     ┌───────────┐     ┌─────────┐     ┌────────┐
│ Not Started │ ──► │ In Draft  │ ──► │Submitted│ ──► │ Graded │
└─────────────┘     └───────────┘     └─────────┘     └────────┘
                          │                │
                          ▼                ▼
                    [Save Draft]    [View Feedback]
```

---

### 1.3 Assessment Engine UI (Quiz Taking)

**Current State:** Backend exists, frontend ❌ Missing  
**Location:** `frontend/src/features/assessment/`

#### Existing Backend API

From `assessment_handler.go`:

```go
// Attempts
POST /api/assessments/:assessmentId/attempts        // Start attempt
POST /api/assessments/attempts/:id/submit           // Submit attempt
POST /api/assessments/attempts/:id/proctor-event    // Proctoring events

// Item Bank (already working)
GET  /api/item-banks/banks
POST /api/item-banks/banks
GET  /api/item-banks/banks/:bankId/items
POST /api/item-banks/banks/:bankId/items
```

#### Missing Backend Endpoints (Need Implementation)

```go
// Assessment CRUD
POST   /api/assessments                    // Create assessment ⚠️ NOT IMPLEMENTED
GET    /api/assessments/:id                // Get assessment details
PUT    /api/assessments/:id                // Update assessment
DELETE /api/assessments/:id                // Delete assessment
POST   /api/assessments/:id/publish        // Publish assessment

// Student-facing
GET    /api/student/assessments            // List available assessments
GET    /api/student/assessments/:id        // Get assessment to take
GET    /api/assessments/attempts/:id       // Get attempt details
GET    /api/assessments/attempts/:id/results // Get graded results
```

#### Backend Implementation Required

**File:** `backend/internal/handlers/assessment_handler.go`

Add these methods:

```go
// CreateAssessment - POST /api/assessments
func (h *AssessmentHandler) CreateAssessment(c *gin.Context) {
    // Parse AssessmentCreateRequest
    // Validate questions exist in item bank
    // Create assessment record
    // Return created assessment
}

// GetAssessment - GET /api/assessments/:id
func (h *AssessmentHandler) GetAssessment(c *gin.Context) {
    // Get assessment with questions
    // For students: hide correct answers
    // For teachers: include all data
}

// PublishAssessment - POST /api/assessments/:id/publish
func (h *AssessmentHandler) PublishAssessment(c *gin.Context) {
    // Validate assessment has questions
    // Set status = published
    // Set available_from/until dates
}
```

#### Frontend Implementation

```
frontend/src/features/assessment/
├── api.ts                          # Assessment API client
├── types.ts                        # Assessment types
├── AssessmentList.tsx              # Available assessments for student
├── AssessmentTaking.tsx            # Quiz-taking UI
├── AssessmentResults.tsx           # Results & feedback
├── components/
│   ├── QuestionRenderer.tsx        # Renders any question type
│   ├── MCQQuestion.tsx             # Multiple choice
│   ├── MRQQuestion.tsx             # Multiple response
│   ├── TextQuestion.tsx            # Short/long text
│   ├── TrueFalseQuestion.tsx       # True/False
│   ├── OrderingQuestion.tsx        # Drag to order
│   ├── MatrixQuestion.tsx          # Matrix/grid
│   ├── Timer.tsx                   # Countdown timer
│   ├── ProgressIndicator.tsx       # Question progress
│   └── ProctorShield.tsx           # Anti-cheat measures
└── hooks/
    ├── useAssessmentAttempt.ts     # Attempt state management
    └── useProctoring.ts            # Proctoring events
```

#### Quiz Taking UI Flow

```
┌──────────────────────────────────────────────────────────────┐
│  📝 Midterm Exam                              ⏱️ 45:23      │
│  Question 3 of 20                            [Progress Bar]  │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  What is the capital of Kazakhstan?                          │
│                                                              │
│  ○ Almaty                                                    │
│  ● Astana                                                    │
│  ○ Shymkent                                                  │
│  ○ Karaganda                                                 │
│                                                              │
├──────────────────────────────────────────────────────────────┤
│  [◀ Previous]            3/20           [Next ▶] [Flag 🚩]  │
│                                                              │
│  Question Navigator:                                         │
│  [1✓][2✓][3●][4][5][6][7][8][9][10]...                      │
│                                                              │
│                              [Save & Exit]  [Submit All]     │
└──────────────────────────────────────────────────────────────┘
```

#### Key Features

1. **Timer** - Countdown with warnings at 10min, 5min, 1min
2. **Progress** - Visual indicator of completion
3. **Navigation** - Previous/Next, jump to question
4. **Flagging** - Mark questions for review
5. **Auto-save** - Save answers periodically
6. **Proctoring** - Tab switch detection, fullscreen mode

---

### 1.4 Quiz Builder (Full Functionality)

**Current State:** Basic modal in studio/, needs enhancement  
**Location:** Enhance `frontend/src/features/studio/components/QuizBuilderModal.tsx`

#### Supported Question Types

| Type                    | v11 | Backend | Current | Action  |
| ----------------------- | --- | ------- | ------- | ------- |
| Multiple Choice (MCQ)   | ✅  | ✅      | ✅      | OK      |
| Multiple Response (MRQ) | ✅  | ✅      | ✅      | OK      |
| True/False              | ✅  | ✅      | ❌      | Add     |
| Short Text              | ✅  | ✅      | ✅      | OK      |
| Long Text (Essay)       | ✅  | ✅      | ❌      | Add     |
| Ordering                | ✅  | ❌      | ✅      | Backend |
| Matrix/Grid             | ✅  | ❌      | ❌      | Add     |
| Fill in Blank           | ✅  | ❌      | ❌      | Add     |
| Math/LaTeX              | ✅  | ❌      | ❌      | Add     |
| Section Header          | ✅  | ❌      | ❌      | Add     |

#### Enhanced Quiz Builder Structure

```
frontend/src/features/studio/components/
├── QuizBuilderModal.tsx            # Main builder modal
├── quiz-builder/
│   ├── QuizSettings.tsx            # Time limit, attempts, shuffle
│   ├── QuestionList.tsx            # Drag & drop question list
│   ├── QuestionEditor.tsx          # Edit single question
│   ├── QuestionTypes/
│   │   ├── MCQEditor.tsx
│   │   ├── MRQEditor.tsx
│   │   ├── TrueFalseEditor.tsx
│   │   ├── TextEditor.tsx
│   │   ├── EssayEditor.tsx
│   │   ├── OrderingEditor.tsx
│   │   ├── MatrixEditor.tsx
│   │   ├── FillBlankEditor.tsx
│   │   └── SectionHeader.tsx
│   ├── QuestionBankPicker.tsx      # Import from bank
│   ├── AnswerOptions.tsx           # Reusable options editor
│   ├── FeedbackEditor.tsx          # Correct/incorrect feedback
│   ├── PointsInput.tsx             # Points per question
│   └── QuizPreview.tsx             # Preview mode
```

#### Quiz Settings Configuration

```typescript
interface QuizSettings {
  title: string;
  description?: string;

  // Timing
  time_limit_minutes?: number; // null = unlimited
  late_submission: "allow" | "penalize" | "block";
  late_penalty_percent?: number;

  // Attempts
  max_attempts: number; // 0 = unlimited
  attempt_grading: "highest" | "latest" | "average";

  // Display
  shuffle_questions: boolean;
  shuffle_answers: boolean;
  show_one_at_a_time: boolean;
  allow_backtrack: boolean;

  // Feedback
  show_correct_answers: "immediately" | "after_due" | "never";
  show_points_per_question: boolean;

  // Proctoring
  require_proctoring: boolean;
  lock_browser: boolean;
  webcam_required: boolean;

  // Availability
  available_from?: string; // ISO date
  available_until?: string;

  // Grading
  grading_schema_id?: string;
  passing_score?: number; // Percentage
}
```

#### Question Editor Features

```typescript
interface QuizQuestion {
  id: string;
  type: QuestionType;
  stem: string; // Question text (supports HTML/Markdown)
  stem_format: "plain" | "html" | "markdown" | "latex";

  // For MCQ/MRQ
  options?: QuestionOption[];

  // For Ordering
  items_to_order?: string[];
  correct_order?: number[];

  // For Matrix
  rows?: string[];
  columns?: string[];
  correct_matrix?: boolean[][];

  // For Fill in Blank
  blank_answers?: string[]; // Acceptable answers
  case_sensitive?: boolean;

  // Scoring
  points: number;
  partial_credit: boolean;

  // Feedback
  feedback_correct?: string;
  feedback_incorrect?: string;
  feedback_hint?: string;

  // Media
  image_url?: string;
  audio_url?: string;
  video_url?: string;

  // Metadata
  difficulty?: "easy" | "medium" | "hard";
  tags?: string[];
  source_bank_id?: string; // If imported from bank
}

interface QuestionOption {
  id: string;
  text: string;
  is_correct: boolean;
  points?: number; // For partial credit
  feedback?: string; // Option-specific feedback
}
```

#### Question Bank Integration

```tsx
// QuestionBankPicker.tsx
<Dialog>
  <DialogTrigger>
    <Button variant="outline">
      <Database /> Import from Bank
    </Button>
  </DialogTrigger>
  <DialogContent className="max-w-4xl">
    <DialogHeader>
      <DialogTitle>Import Questions</DialogTitle>
    </DialogHeader>

    <div className="grid grid-cols-3 gap-4">
      {/* Left: Bank list */}
      <div className="border rounded-lg p-4">
        <h4>Question Banks</h4>
        <BankList onSelect={setSelectedBank} />
      </div>

      {/* Middle: Questions */}
      <div className="col-span-2 border rounded-lg p-4">
        <div className="flex justify-between mb-4">
          <Input placeholder="Search questions..." />
          <Select>
            <SelectTrigger>Filter by type</SelectTrigger>
            <SelectContent>
              <SelectItem value="mcq">Multiple Choice</SelectItem>
              <SelectItem value="text">Text</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <QuestionList
          questions={bankQuestions}
          selectedIds={selectedQuestionIds}
          onToggle={toggleQuestion}
        />
      </div>
    </div>

    <DialogFooter>
      <Button onClick={importSelected}>
        Import {selectedQuestionIds.length} Questions
      </Button>
    </DialogFooter>
  </DialogContent>
</Dialog>
```

#### Math/LaTeX Support

Add KaTeX for math rendering:

```bash
npm install katex react-katex
```

```tsx
// components/MathRenderer.tsx
import "katex/dist/katex.min.css";
import { InlineMath, BlockMath } from "react-katex";

export const MathRenderer: React.FC<{ content: string }> = ({ content }) => {
  // Parse content for LaTeX delimiters
  // $...$ for inline, $$...$$ for block
  return (
    <div>
      {parseLatex(content).map((part, i) =>
        part.type === "latex-inline" ? (
          <InlineMath key={i} math={part.content} />
        ) : part.type === "latex-block" ? (
          <BlockMath key={i} math={part.content} />
        ) : (
          <span key={i}>{part.content}</span>
        )
      )}
    </div>
  );
};
```

---

## Phase 2: Survey/Form Builder & Teacher Tools (High)

**Timeline:** 2-3 weeks  
**Goal:** Complete assessment tools and teacher tracking

### 2.1 Survey Builder Enhancement

**Current State:** Basic in studio/  
**Location:** `frontend/src/features/studio/components/SurveyBuilderModal.tsx`

#### Additional Survey Question Types

| Type          | Description            | Status    |
| ------------- | ---------------------- | --------- |
| Star Rating   | 1-5 stars              | ✅ Exists |
| NPS (0-10)    | Net Promoter Score     | ✅ Exists |
| Likert Matrix | Agreement scale grid   | ✅ Exists |
| Open Text     | Long text feedback     | ✅ Exists |
| Dropdown      | Single select dropdown | ❌ Add    |
| Ranking       | Rank items in order    | ❌ Add    |
| Slider        | Numeric slider         | ❌ Add    |
| Date/Time     | Date picker            | ❌ Add    |
| File Upload   | Attachment             | ❌ Add    |

#### Survey Settings

```typescript
interface SurveySettings {
  title: string;
  description?: string;

  // Privacy
  anonymous: boolean;
  collect_email: boolean;

  // Display
  show_progress_bar: boolean;
  show_question_numbers: boolean;
  one_question_per_page: boolean;
  allow_back_navigation: boolean;

  // Completion
  thank_you_message: string;
  redirect_url?: string;

  // Availability
  available_from?: string;
  available_until?: string;
  max_responses?: number;

  // Logic
  conditional_logic: ConditionalRule[];
}

interface ConditionalRule {
  if_question_id: string;
  condition: "equals" | "not_equals" | "contains" | "greater_than";
  value: any;
  then_show_question_id: string;
}
```

### 2.2 Form Builder

**Current State:** Basic modal exists  
**Location:** `frontend/src/features/studio/components/FormBuilderModal.tsx`

#### Form Field Types

| Field Type   | Description        | Validation                      |
| ------------ | ------------------ | ------------------------------- |
| Text         | Single line        | Required, min/max length, regex |
| Textarea     | Multi-line         | Required, min/max length        |
| Number       | Numeric input      | Required, min/max, step         |
| Email        | Email address      | Format validation               |
| Phone        | Phone number       | Format validation               |
| Date         | Date picker        | Min/max date                    |
| DateTime     | Date + time        | -                               |
| Select       | Dropdown           | Required                        |
| Multi-Select | Multiple choice    | Min/max selections              |
| Checkbox     | Boolean            | -                               |
| Radio        | Single choice      | Required                        |
| File         | File upload        | Types, max size                 |
| Signature    | Signature pad      | Required                        |
| Dictionary   | Link to dictionary | Auto-populate options           |

#### Dictionary Integration

```tsx
// Connect form fields to dictionaries
<FormFieldEditor field={field} onChange={updateField}>
  {field.type === "select" && (
    <DictionaryPicker
      label="Options Source"
      value={field.dictionary_id}
      onChange={(dictId) =>
        updateField({
          ...field,
          dictionary_id: dictId,
          options: "from_dictionary",
        })
      }
    />
  )}
</FormFieldEditor>
```

### 2.3 Teacher Student Tracker

**Current State:** ❌ Missing  
**Location:** `frontend/src/features/teacher/StudentTracker.tsx`

#### Backend Requirements

```
GET /api/teacher/courses/:id/students         # Students with progress
GET /api/teacher/courses/:id/at-risk          # At-risk students
GET /api/teacher/students/:id/activity        # Student activity log
```

#### UI Components

```
frontend/src/features/teacher/
├── StudentTracker.tsx              # Main tracker page
├── components/
│   ├── StudentProgressTable.tsx    # Table with progress data
│   ├── RiskIndicator.tsx           # At-risk badge
│   ├── EngagementChart.tsx         # Activity over time
│   ├── StudentDetailDrawer.tsx     # Detailed student view
│   └── InterventionActions.tsx     # Send reminder, etc.
```

#### Risk Calculation Logic

```typescript
interface StudentRiskProfile {
  student_id: string;
  student_name: string;

  // Metrics
  overall_progress: number; // 0-100%
  assignments_completed: number;
  assignments_total: number;
  assignments_overdue: number;
  last_activity: string; // ISO date
  days_inactive: number;
  average_grade: number;

  // Risk assessment
  risk_level: "low" | "medium" | "high" | "critical";
  risk_factors: string[]; // ["3 overdue assignments", "14 days inactive"]

  // Recommendations
  suggested_actions: string[];
}
```

#### UI Mockup

```
┌──────────────────────────────────────────────────────────────┐
│  📊 Student Progress Tracker - CS101                         │
│  Showing 45 students                    [Filter ▼] [Export]  │
├──────────────────────────────────────────────────────────────┤
│ ⚠️ 5 students at risk    ✅ 32 on track    📈 8 ahead       │
├──────────────────────────────────────────────────────────────┤
│ Student        │ Progress │ Grade │ Last Active │ Status     │
│────────────────┼──────────┼───────┼─────────────┼────────────│
│ 🔴 John Doe    │ ████░░ 40%│ 65%  │ 14 days ago │ At Risk    │
│ 🟡 Jane Smith  │ █████░ 70%│ 78%  │ 2 days ago  │ Needs Help │
│ 🟢 Bob Wilson  │ ██████ 95%│ 92%  │ Today       │ On Track   │
│ ...            │          │       │             │            │
├──────────────────────────────────────────────────────────────┤
│ Selected: John Doe                                           │
│ ┌──────────────────────────────────────────────────────────┐ │
│ │ Risk Factors:                                            │ │
│ │ • 3 overdue assignments                                  │ │
│ │ • 14 days since last login                               │ │
│ │ • Grade dropped 15% this month                           │ │
│ │                                                          │ │
│ │ Suggested Actions:                                       │ │
│ │ [📧 Send Reminder] [📅 Schedule Meeting] [📝 Add Note]   │ │
│ └──────────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────┘
```

### 2.4 Course Library

**Current State:** ❌ Missing  
**Location:** `frontend/src/features/admin/CourseLibrary.tsx`

A central place for admins to browse all courses and course templates.

#### Features

1. **Browse Courses** - Grid/list view of all courses
2. **Filter & Search** - By program, department, status
3. **Course Templates** - Reusable course structures
4. **Import/Export** - Common course format (LTI, SCORM)
5. **Analytics** - Enrollment stats, completion rates

---

## Phase 3: Rich Editor & Communication (Medium)

**Timeline:** 3-4 weeks  
**Goal:** Enhanced content creation and communication tools

### 3.1 Rich Question Editor

**Location:** `frontend/src/features/item-bank/components/QuestionEditor.tsx`

#### Features

| Feature   | Description                | Implementation         |
| --------- | -------------------------- | ---------------------- |
| Rich Text | Bold, italic, lists, links | TipTap or Lexical      |
| Math      | LaTeX equations            | KaTeX                  |
| Images    | Inline images              | Upload to S3           |
| Audio     | Audio clips                | HTML5 Audio            |
| Video     | Embedded video             | iframe/player          |
| Code      | Syntax highlighting        | Prism.js               |
| Tables    | Data tables                | TipTap table extension |

#### Implementation

```tsx
// QuestionEditor.tsx
import { useEditor, EditorContent } from "@tiptap/react";
import StarterKit from "@tiptap/starter-kit";
import Mathematics from "@tiptap/extension-mathematics";
import Image from "@tiptap/extension-image";

export const QuestionEditor: React.FC = () => {
  const editor = useEditor({
    extensions: [
      StarterKit,
      Mathematics,
      Image.configure({
        uploadImage: async (file) => {
          const url = await uploadToS3(file);
          return url;
        },
      }),
    ],
  });

  return (
    <div>
      <EditorToolbar editor={editor}>
        <ToolbarButton icon={Bold} command="bold" />
        <ToolbarButton icon={Italic} command="italic" />
        <ToolbarButton icon={List} command="bulletList" />
        <ToolbarButton icon={Sigma} command="math" />
        <ToolbarButton icon={Image} command="image" />
        <ToolbarButton icon={Code} command="codeBlock" />
      </EditorToolbar>
      <EditorContent editor={editor} />
    </div>
  );
};
```

### 3.2 Discussion Forums

**Current State:** Backend exists, UI ❌ Missing  
**Location:** `frontend/src/features/forums/`

#### Backend Endpoints (Existing)

```
GET    /api/forums                          # List forums
POST   /api/forums                          # Create forum
GET    /api/forums/:id                      # Get forum
GET    /api/forums/:id/threads              # List threads
POST   /api/forums/:id/threads              # Create thread
GET    /api/forums/threads/:id              # Get thread with posts
POST   /api/forums/threads/:id/posts        # Add post
PUT    /api/forums/posts/:id                # Edit post
DELETE /api/forums/posts/:id                # Delete post
POST   /api/forums/posts/:id/like           # Like post
```

#### Frontend Implementation

```
frontend/src/features/forums/
├── api.ts
├── types.ts
├── ForumList.tsx                   # List of forums
├── ForumDetail.tsx                 # Forum with threads
├── ThreadDetail.tsx                # Thread with posts
├── components/
│   ├── ThreadCard.tsx
│   ├── PostCard.tsx
│   ├── PostEditor.tsx              # Rich text editor
│   ├── ReplyButton.tsx
│   └── LikeButton.tsx
```

### 3.3 Attendance Tracking

**Current State:** Backend exists, UI ❌ Missing  
**Location:** `frontend/src/features/attendance/`

#### Backend Endpoints (Existing)

```
GET  /api/attendance/sessions                # List sessions
POST /api/attendance/sessions                # Create session
GET  /api/attendance/sessions/:id            # Session with records
POST /api/attendance/sessions/:id/check-in   # Student check-in
PUT  /api/attendance/records/:id             # Update record
GET  /api/attendance/student/:id/summary     # Student summary
```

#### UI for Teachers

```tsx
// AttendanceSession.tsx
<Card>
  <CardHeader>
    <CardTitle>CS101 - Lecture 15</CardTitle>
    <CardDescription>January 4, 2026 • 10:00 AM</CardDescription>
  </CardHeader>
  <CardContent>
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Student</TableHead>
          <TableHead>Status</TableHead>
          <TableHead>Check-in Time</TableHead>
          <TableHead>Notes</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {students.map((student) => (
          <TableRow key={student.id}>
            <TableCell>{student.name}</TableCell>
            <TableCell>
              <Select defaultValue={student.status}>
                <SelectItem value="present">✅ Present</SelectItem>
                <SelectItem value="late">🕐 Late</SelectItem>
                <SelectItem value="excused">📝 Excused</SelectItem>
                <SelectItem value="absent">❌ Absent</SelectItem>
              </Select>
            </TableCell>
            <TableCell>{student.check_in_time}</TableCell>
            <TableCell>
              <Input placeholder="Add note..." />
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  </CardContent>
  <CardFooter>
    <Button>Save Attendance</Button>
    <Button variant="outline">Generate QR Code</Button>
  </CardFooter>
</Card>
```

#### QR Code Check-in

Students can scan QR code to check in:

```tsx
// QRCheckIn.tsx (for students)
<Card className="max-w-md mx-auto text-center">
  <CardContent className="pt-6">
    <div className="w-48 h-48 mx-auto bg-slate-100 rounded-xl flex items-center justify-center">
      {/* Camera viewfinder */}
      <QRScanner onScan={handleCheckIn} />
    </div>
    <p className="mt-4 text-sm text-slate-500">
      Scan the QR code displayed by your instructor
    </p>
  </CardContent>
</Card>
```

---

## Phase 4: Advanced Features (Nice-to-Have)

**Timeline:** 3-4 weeks  
**Goal:** Advanced features for best-in-class experience and student engagement

### 4.1 Gamification System (Full Stack)

**Current State:** ❌ Not implemented  
**Backend Required:** Yes (new tables, service, handler)  
**Frontend Required:** Yes (new module)

This is a **complete full-stack feature** requiring both backend and frontend work.

#### 4.1.1 Database Schema

**File:** `backend/db/migrations/XXXXXX_create_gamification_tables.sql`

```sql
-- User XP and Level tracking
CREATE TABLE user_xp (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    total_xp INT DEFAULT 0,
    level INT DEFAULT 1,
    current_streak INT DEFAULT 0,
    longest_streak INT DEFAULT 0,
    last_activity_date DATE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, user_id)
);

-- Badge definitions
CREATE TABLE badges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL,           -- 'first_submission', 'perfect_quiz', etc.
    name VARCHAR(255) NOT NULL,
    description TEXT,
    icon_url VARCHAR(500),
    category VARCHAR(50) DEFAULT 'achievement', -- 'achievement', 'milestone', 'special'
    criteria JSONB NOT NULL,              -- {"type": "count", "event": "quiz_perfect", "threshold": 5}
    xp_reward INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    display_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, code)
);

-- User earned badges
CREATE TABLE user_badges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    badge_id UUID NOT NULL REFERENCES badges(id) ON DELETE CASCADE,
    earned_at TIMESTAMP DEFAULT NOW(),
    notified BOOLEAN DEFAULT false,
    UNIQUE(user_id, badge_id)
);

-- XP transaction log
CREATE TABLE xp_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL,
    xp_amount INT NOT NULL,
    source_type VARCHAR(50),              -- 'assignment', 'quiz', 'forum', 'streak'
    source_id UUID,                       -- Reference to the entity that triggered XP
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Leaderboard cache (refreshed periodically)
CREATE TABLE leaderboard_cache (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    period VARCHAR(20) NOT NULL,          -- 'daily', 'weekly', 'monthly', 'all_time'
    rank INT NOT NULL,
    total_xp INT NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, user_id, period)
);

CREATE INDEX idx_xp_events_user ON xp_events(tenant_id, user_id, created_at DESC);
CREATE INDEX idx_leaderboard_rank ON leaderboard_cache(tenant_id, period, rank);
```

#### 4.1.2 Backend Implementation

**File:** `backend/internal/services/gamification_service.go`

```go
package services

import (
    "context"
    "time"
)

type GamificationService struct {
    db     *sqlx.DB
    notify *NotificationService
}

// XP Event Types and Rewards
var XPRewards = map[string]int{
    "assignment_submitted":   10,
    "assignment_on_time":     5,
    "assignment_early":       10,  // More than 24h before deadline
    "quiz_completed":         10,
    "quiz_passed":            15,
    "quiz_perfect":           25,
    "course_completed":       100,
    "module_completed":       20,
    "daily_login":            5,
    "streak_bonus":           2,   // Per day of streak
    "forum_post":             5,
    "forum_helpful":          10,
    "first_submission":       20,
    "profile_complete":       15,
}

// Level thresholds (XP required for each level)
var LevelThresholds = []int{
    0,      // Level 1
    100,    // Level 2
    300,    // Level 3
    600,    // Level 4
    1000,   // Level 5
    1500,   // Level 6
    2100,   // Level 7
    2800,   // Level 8
    3600,   // Level 9
    4500,   // Level 10
    // ... continues
}



XP Event Configuration
Event Type	Base XP	Multipliers	Description
assignment_submitted	10	-	Submit any assignment
assignment_on_time	+5	-	Bonus for on-time submission
assignment_perfect	+15	-	100% score on assignment
quiz_completed	10	-	Complete any quiz
quiz_passed	+5	-	Pass quiz (≥70%)
quiz_perfect	+25	-	100% on quiz
course_module_completed	15	-	Complete a course module
course_completed	100	-	Complete entire course
daily_login	5	×streak	Daily login (5×streak days)
forum_post	5	-	Create forum post
forum_helpful	+10	-	Post marked helpful
forum_best_answer	+20	-	Answer selected as best
attendance_present	3	-	Attend class
attendance_streak	+2	×week	Weekly attendance streak
first_submission	50	-	First ever submission (one-time)
journey_node_completed	20	-	Complete journey node
Level Progression
Level	XP Required	Title
1	0	Newcomer
2	100	Learner
3	250	Student
4	500	Scholar
5	1000	Advanced Scholar
6	2000	Expert
7	3500	Master
8	5500	Grandmaster
9	8000	Legend
10	12000	Champion
Badge Categories
Academic Badges:

Badge	Criteria	XP Reward	Rarity
First Steps	Complete first assignment	25	Common
Quiz Whiz	Pass 10 quizzes	50	Common
Perfect Score	Get 100% on any assessment	30	Uncommon
Honor Roll	Maintain 90%+ average for a month	100	Rare
Valedictorian	Complete program with highest GPA	500	Legendary
Engagement Badges:

Badge	Criteria	XP Reward	Rarity
Early Bird	Submit 5 assignments before deadline	30	Common
Consistent	7-day login streak	40	Common
Dedicated	30-day login streak	150	Rare
Contributor	Create 10 helpful forum posts	75	Uncommon
Mentor	Have 5 answers marked as best	200	Epic
Milestone Badges:

Badge	Criteria	XP Reward	Rarity
First Course	Complete first course	50	Common
Halfway There	Complete 50% of program	150	Uncommon
Almost Done	Complete 90% of program	300	Rare
Graduate	Complete entire program	1000	Epic


func (s *GamificationService) AwardXP(ctx context.Context, tenantID, userID, eventType string, sourceType string, sourceID *string) error {
    xpAmount := XPRewards[eventType]
    if xpAmount == 0 {
        return nil
    }

    // Record XP event
    _, err := s.db.ExecContext(ctx, `
        INSERT INTO xp_events (tenant_id, user_id, event_type, xp_amount, source_type, source_id)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, tenantID, userID, eventType, xpAmount, sourceType, sourceID)
    if err != nil {
        return err
    }

    // Update user XP and check for level up
    return s.updateUserXP(ctx, tenantID, userID, xpAmount)
}

func (s *GamificationService) CheckAndAwardBadges(ctx context.Context, tenantID, userID string) ([]Badge, error) {
    // Get all badges user doesn't have yet
    // Check criteria for each
    // Award if criteria met
    // Return newly earned badges
}

func (s *GamificationService) GetUserProfile(ctx context.Context, tenantID, userID string) (*GamificationProfile, error) {
    // Return XP, level, badges, rank, streak
}

func (s *GamificationService) GetLeaderboard(ctx context.Context, tenantID, period string, limit int) ([]LeaderboardEntry, error) {
    // Return top users for period
}

func (s *GamificationService) UpdateStreak(ctx context.Context, tenantID, userID string) error {
    // Check last activity, update streak
}
```

**File:** `backend/internal/handlers/gamification_handler.go`

```go
package handlers

type GamificationHandler struct {
    svc *services.GamificationService
}

// GET /api/gamification/profile
func (h *GamificationHandler) GetProfile(c *gin.Context)

// GET /api/gamification/leaderboard?period=weekly&limit=10
func (h *GamificationHandler) GetLeaderboard(c *gin.Context)

// GET /api/gamification/badges
func (h *GamificationHandler) ListBadges(c *gin.Context)

// GET /api/gamification/history?limit=20
func (h *GamificationHandler) GetXPHistory(c *gin.Context)

// Admin endpoints
// POST /api/admin/gamification/badges
func (h *GamificationHandler) CreateBadge(c *gin.Context)

// PUT /api/admin/gamification/badges/:id
func (h *GamificationHandler) UpdateBadge(c *gin.Context)

// POST /api/admin/gamification/award-xp
func (h *GamificationHandler) ManualAwardXP(c *gin.Context)
```

#### 4.1.3 Integration Points

Add XP awards to existing services:

```go
// In assessment_service.go - after quiz completion
func (s *AssessmentService) CompleteAttempt(...) {
    // ... existing logic ...

    // Award XP
    s.gamification.AwardXP(ctx, tenantID, userID, "quiz_completed", "quiz", &quizID)
    if score == maxScore {
        s.gamification.AwardXP(ctx, tenantID, userID, "quiz_perfect", "quiz", &quizID)
    }
    s.gamification.CheckAndAwardBadges(ctx, tenantID, userID)
}

// In node_handler.go - after journey submission
func (h *NodeHandler) SubmitNode(...) {
    // ... existing logic ...

    h.gamification.AwardXP(ctx, tenantID, userID, "assignment_submitted", "node", &nodeID)
    if submittedBeforeDeadline {
        h.gamification.AwardXP(ctx, tenantID, userID, "assignment_on_time", "node", &nodeID)
    }
}

// In auth_handler.go - after login
func (h *AuthHandler) Login(...) {
    // ... existing logic ...

    h.gamification.UpdateStreak(ctx, tenantID, userID)
    h.gamification.AwardXP(ctx, tenantID, userID, "daily_login", "login", nil)
}
```

#### 4.1.4 Frontend Implementation

**File structure:**

```
frontend/src/features/gamification/
├── api.ts
├── types.ts
├── GamificationProfile.tsx         # User's XP, level, badges
├── Leaderboard.tsx                  # Weekly/monthly rankings
├── BadgeShowcase.tsx                # Display earned badges
├── XPHistory.tsx                    # XP transaction history
├── components/
│   ├── XPProgress.tsx               # Level progress bar
│   ├── BadgeCard.tsx                # Single badge display
│   ├── LeaderboardRow.tsx           # Leaderboard entry
│   ├── XPNotification.tsx           # Toast when XP earned
│   ├── LevelUpModal.tsx             # Celebration on level up
│   └── StreakIndicator.tsx          # Daily streak display
└── hooks/
    ├── useGamification.ts           # Main hook
    └── useXPNotifications.ts        # Real-time XP toasts
```

**Types:**

```typescript
// frontend/src/features/gamification/types.ts

export interface GamificationProfile {
  user_id: string;
  total_xp: number;
  level: number;
  xp_to_next_level: number;
  current_streak: number;
  longest_streak: number;
  rank: number;
  badges: Badge[];
  recent_xp: XPEvent[];
}

export interface Badge {
  id: string;
  code: string;
  name: string;
  description: string;
  icon_url: string;
  category: "achievement" | "milestone" | "special";
  xp_reward: number;
  earned_at?: string;
  is_earned: boolean;
}

export interface XPEvent {
  id: string;
  event_type: string;
  xp_amount: number;
  source_type: string;
  created_at: string;
}

export interface LeaderboardEntry {
  rank: number;
  user_id: string;
  user_name: string;
  avatar_url?: string;
  total_xp: number;
  level: number;
}
```

**API:**

```typescript
// frontend/src/features/gamification/api.ts

import { api } from "@/api/client";
import type {
  GamificationProfile,
  Badge,
  LeaderboardEntry,
  XPEvent,
} from "./types";

export const getGamificationProfile = () =>
  api.get<GamificationProfile>("/gamification/profile");

export const getLeaderboard = (
  period: "daily" | "weekly" | "monthly" | "all_time",
  limit = 10
) =>
  api.get<LeaderboardEntry[]>(
    `/gamification/leaderboard?period=${period}&limit=${limit}`
  );

export const getAllBadges = () => api.get<Badge[]>("/gamification/badges");

export const getXPHistory = (limit = 20) =>
  api.get<XPEvent[]>(`/gamification/history?limit=${limit}`);
```

**Main Component:**

```tsx
// frontend/src/features/gamification/GamificationProfile.tsx

import React from "react";
import { useQuery } from "@tanstack/react-query";
import { Trophy, Flame, Star, TrendingUp } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { Badge } from "@/components/ui/badge";
import { getGamificationProfile } from "./api";

export const GamificationProfile: React.FC = () => {
  const { data: profile, isLoading } = useQuery({
    queryKey: ["gamification", "profile"],
    queryFn: getGamificationProfile,
  });

  if (isLoading || !profile) return <div>Loading...</div>;

  const progressPercent = ((profile.total_xp % 100) / 100) * 100; // Simplified

  return (
    <div className="space-y-6">
      {/* Level & XP Card */}
      <Card className="bg-gradient-to-r from-indigo-500 to-purple-600 text-white">
        <CardContent className="pt-6">
          <div className="flex items-center justify-between">
            <div>
              <div className="text-sm opacity-80">Level</div>
              <div className="text-5xl font-black">{profile.level}</div>
            </div>
            <div className="text-right">
              <div className="text-sm opacity-80">Total XP</div>
              <div className="text-3xl font-bold">
                {profile.total_xp.toLocaleString()}
              </div>
            </div>
          </div>
          <div className="mt-4">
            <div className="flex justify-between text-sm mb-1">
              <span>Progress to Level {profile.level + 1}</span>
              <span>{profile.xp_to_next_level} XP needed</span>
            </div>
            <Progress value={progressPercent} className="h-3" />
          </div>
        </CardContent>
      </Card>

      {/* Stats Row */}
      <div className="grid grid-cols-3 gap-4">
        <Card>
          <CardContent className="pt-4 text-center">
            <Flame className="mx-auto text-orange-500" size={24} />
            <div className="text-2xl font-bold mt-2">
              {profile.current_streak}
            </div>
            <div className="text-xs text-muted-foreground">Day Streak</div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-4 text-center">
            <Trophy className="mx-auto text-yellow-500" size={24} />
            <div className="text-2xl font-bold mt-2">#{profile.rank}</div>
            <div className="text-xs text-muted-foreground">Global Rank</div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-4 text-center">
            <Star className="mx-auto text-blue-500" size={24} />
            <div className="text-2xl font-bold mt-2">
              {profile.badges.length}
            </div>
            <div className="text-xs text-muted-foreground">Badges</div>
          </CardContent>
        </Card>
      </div>

      {/* Badges */}
      <Card>
        <CardHeader>
          <CardTitle>Badges Earned</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-4 gap-4">
            {profile.badges.map((badge) => (
              <div key={badge.id} className="text-center">
                <img
                  src={badge.icon_url}
                  alt={badge.name}
                  className="w-16 h-16 mx-auto"
                />
                <div className="text-sm font-medium mt-2">{badge.name}</div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
};
```

#### 4.1.5 XP Event Types (Complete List)

| Event                  | XP  | Trigger                    | Category   |
| ---------------------- | --- | -------------------------- | ---------- |
| `daily_login`          | 5   | First login of the day     | Engagement |
| `streak_3_days`        | 15  | 3-day login streak         | Streak     |
| `streak_7_days`        | 50  | 7-day login streak         | Streak     |
| `streak_30_days`       | 200 | 30-day login streak        | Streak     |
| `assignment_submitted` | 10  | Submit any assignment      | Academic   |
| `assignment_on_time`   | 5   | Submit before deadline     | Academic   |
| `assignment_early`     | 10  | Submit 24h+ early          | Academic   |
| `quiz_completed`       | 10  | Complete any quiz          | Academic   |
| `quiz_passed`          | 15  | Score >= passing threshold | Academic   |
| `quiz_perfect`         | 25  | 100% score                 | Academic   |
| `module_completed`     | 20  | Complete all module items  | Progress   |
| `course_completed`     | 100 | Complete entire course     | Progress   |
| `program_milestone`    | 50  | Complete program stage     | Progress   |
| `forum_post`           | 5   | Create forum post          | Community  |
| `forum_reply`          | 3   | Reply to a thread          | Community  |
| `forum_helpful`        | 10  | Post marked as helpful     | Community  |
| `profile_complete`     | 15  | Fill all profile fields    | Onboarding |
| `first_submission`     | 20  | First ever submission      | Onboarding |
| `peer_review`          | 10  | Complete peer review       | Community  |

#### 4.1.6 Badge Definitions

| Badge            | Code                  | Criteria                | XP Reward |
| ---------------- | --------------------- | ----------------------- | --------- |
| 🌟 First Steps   | `first_submission`    | Submit first assignment | 20        |
| 🔥 On Fire       | `streak_7`            | 7-day streak            | 50        |
| 🏆 Perfect Score | `quiz_perfect_first`  | First 100% quiz         | 25        |
| 📚 Bookworm      | `courses_3`           | Complete 3 courses      | 100       |
| 🎯 Sharpshooter  | `quiz_perfect_5`      | 5 perfect quizzes       | 100       |
| 💬 Contributor   | `forum_posts_10`      | 10 forum posts          | 50        |
| ⭐ Helpful       | `helpful_5`           | 5 helpful answers       | 75        |
| 🚀 Early Bird    | `early_submissions_5` | 5 early submissions     | 50        |
| 🎓 Scholar       | `level_10`            | Reach level 10          | 200       |
| 👑 Master        | `level_25`            | Reach level 25          | 500       |

### 4.2 Global Search

**Backend:** Already exists (`/api/search`)

```
frontend/src/features/search/
├── GlobalSearch.tsx                # Search overlay
├── SearchResults.tsx               # Results display
├── components/
│   ├── SearchInput.tsx             # Command+K trigger
│   ├── ResultCategory.tsx          # Group by type
│   └── SearchFilters.tsx           # Filter by type
```

#### Search UI (Command Palette Style)

```tsx
// GlobalSearch.tsx - triggered by Cmd+K
<Dialog open={open} onOpenChange={setOpen}>
  <DialogContent className="max-w-2xl p-0">
    <Command>
      <CommandInput placeholder="Search courses, students, assignments..." />
      <CommandList>
        <CommandEmpty>No results found.</CommandEmpty>

        <CommandGroup heading="Courses">
          {results.courses.map((course) => (
            <CommandItem onSelect={() => navigate(`/courses/${course.id}`)}>
              <BookOpen className="mr-2" />
              {course.title}
            </CommandItem>
          ))}
        </CommandGroup>

        <CommandGroup heading="Students">
          {results.students.map((student) => (
            <CommandItem>
              <User className="mr-2" />
              {student.name}
            </CommandItem>
          ))}
        </CommandGroup>

        <CommandGroup heading="Assignments">
          {results.assignments.map((assignment) => (
            <CommandItem>
              <FileText className="mr-2" />
              {assignment.title}
            </CommandItem>
          ))}
        </CommandGroup>
      </CommandList>
    </Command>
  </DialogContent>
</Dialog>
```

### 4.3 AI Content Generation

**Backend:** Exists (`/api/ai/*`)

```
GET  /api/ai/generate-questions     # Generate quiz questions
POST /api/ai/summarize              # Summarize content
POST /api/ai/feedback               # Generate feedback
```

#### UI Integration Points

1. **Quiz Builder** - "Generate questions from topic"
2. **Course Content** - "Summarize this reading"
3. **Grading** - "Suggest feedback for this submission"

### 4.4 LTI 1.3 Integration

**Backend:** Exists, UI needed

For integrating external tools (Zoom, Kaltura, etc.)

```
frontend/src/features/admin/lti/
├── LTIToolsPage.tsx                # Manage LTI tools
├── LTIToolForm.tsx                 # Add/edit tool
├── LTILaunchButton.tsx             # Launch tool in course
```

---

## Backend Gaps Summary

### ✅ Already Implemented Backend APIs

Based on audit (January 4, 2026), the following are **fully operational**:

| Category          | Endpoints                                              | Handler                  |
| ----------------- | ------------------------------------------------------ | ------------------------ |
| Assessment Engine | `/api/assessments/*/attempts`, submit, proctor         | assessment_handler.go ✅ |
| Student Portal    | `/api/student/dashboard`, courses, assignments, grades | student_handler.go ✅    |
| Teacher Portal    | `/api/teacher/dashboard`, courses, roster, gradebook   | teacher_handler.go ✅    |
| Item Bank         | `/api/item-banks/*` full CRUD                          | item_bank_handler.go ✅  |
| Grading           | `/api/grading/*` schemas, entries                      | grading_handler.go ✅    |
| Forums            | `/api/forums/*` full CRUD                              | forum_handler.go ✅      |
| Attendance        | `/api/attendance/*` sessions, check-in                 | attendance_handler.go ✅ |
| AI Generation     | `/api/ai/*` generate, summarize                        | ai_handler.go ✅         |
| LTI 1.3           | `/api/lti/*` login, launch, jwks                       | lti_handler.go ✅        |
| Search            | `/api/search` global search                            | search_handler.go ✅     |
| Transcript        | `/api/transcript/:studentId`                           | transcript_handler.go ✅ |

### 🔴 Critical Gaps (Phase 1)

| Endpoint                                   | Handler                   | Status     | Action                          |
| ------------------------------------------ | ------------------------- | ---------- | ------------------------------- |
| `GET /api/student/courses/:id`             | student_handler.go        | ⚠️ Verify  | Check if exists, add if missing |
| `GET /api/student/courses/:id/modules`     | course_content_handler.go | ⚠️ Verify  | May need student-facing wrapper |
| `POST /api/student/assignments/:id/submit` | -                         | ❌ Missing | Create new endpoint             |
| `GET /api/student/assignments/:id`         | student_handler.go        | ⚠️ Verify  | Assignment detail for student   |

### 🟠 High Priority Gaps (Phase 2)

| Endpoint                                 | Handler            | Status     | Action                        |
| ---------------------------------------- | ------------------ | ---------- | ----------------------------- |
| `GET /api/teacher/courses/:id/at-risk`   | teacher_handler.go | ❌ Missing | Add risk calculation endpoint |
| `GET /api/teacher/students/:id/activity` | teacher_handler.go | ❌ Missing | Student activity log          |

### 🟢 New Backend for Phase 4 (Gamification)

| Endpoint                                | Handler                 | Status     | Action                 |
| --------------------------------------- | ----------------------- | ---------- | ---------------------- |
| `GET /api/gamification/profile`         | gamification_handler.go | ❌ Missing | Create full module     |
| `GET /api/gamification/leaderboard`     | gamification_handler.go | ❌ Missing | Create full module     |
| `GET /api/gamification/badges`          | gamification_handler.go | ❌ Missing | Create full module     |
| `GET /api/gamification/history`         | gamification_handler.go | ❌ Missing | Create full module     |
| `POST /api/admin/gamification/badges`   | gamification_handler.go | ❌ Missing | Admin badge management |
| `POST /api/admin/gamification/award-xp` | gamification_handler.go | ❌ Missing | Manual XP award        |

### Frontend Gaps Summary

#### 🔴 Critical (Phase 1)

| Component                         | Location                       | Status    |
| --------------------------------- | ------------------------------ | --------- |
| StudentCourseDetail.tsx           | features/student-portal/       | ❌ Create |
| StudentAssignmentDetail.tsx       | features/student-portal/       | ❌ Create |
| API functions for detail pages    | features/student-portal/api.ts | ❌ Add    |
| Routes `/student/courses/:id`     | routes/index.tsx               | ❌ Add    |
| Routes `/student/assignments/:id` | routes/index.tsx               | ❌ Add    |

#### 🟠 High Priority (Phase 2)

| Component                 | Location             | Status           |
| ------------------------- | -------------------- | ---------------- |
| ForumList.tsx             | features/forums/     | ❌ Create module |
| ForumDetail.tsx           | features/forums/     | ❌ Create module |
| ThreadDetail.tsx          | features/forums/     | ❌ Create module |
| AttendanceSession.tsx     | features/attendance/ | ❌ Create module |
| AttendanceHistory.tsx     | features/attendance/ | ❌ Create module |
| StudentTracker.tsx        | features/teacher/    | ❌ Add           |
| Route `/admin/forums`     | routes/index.tsx     | ❌ Add           |
| Route `/admin/attendance` | routes/index.tsx     | ❌ Add           |

#### 🟡 Medium Priority (Phase 3)

| Component               | Location                       | Status     |
| ----------------------- | ------------------------------ | ---------- |
| GradebookPage.tsx       | features/grading/              | ❌ Create  |
| TranscriptView.tsx      | features/student-portal/       | ❌ Create  |
| Enhanced QuestionEditor | features/item-bank/components/ | ⚠️ Enhance |

#### 🟢 Low Priority (Phase 4)

| Component               | Location               | Status           |
| ----------------------- | ---------------------- | ---------------- |
| GamificationProfile.tsx | features/gamification/ | ❌ Create module |
| Leaderboard.tsx         | features/gamification/ | ❌ Create module |
| BadgeShowcase.tsx       | features/gamification/ | ❌ Create module |
| GlobalSearch.tsx        | features/search/       | ❌ Create module |
| AIToolsPanel.tsx        | features/admin/ai/     | ❌ Create module |
| LTIToolsPage.tsx        | features/admin/lti/    | ❌ Create module |

---

## Testing Strategy

### Unit Tests

```typescript
// Assessment taking
describe("AssessmentTaking", () => {
  it("should start an attempt when component mounts");
  it("should save answer when user selects option");
  it("should auto-save periodically");
  it("should show timer and warn when low");
  it("should submit all answers on completion");
  it("should handle network errors gracefully");
});

// Quiz Builder
describe("QuizBuilder", () => {
  it("should add new question of each type");
  it("should reorder questions via drag and drop");
  it("should import questions from bank");
  it("should validate required fields");
  it("should save quiz to backend");
});
```

### E2E Tests

```typescript
// Playwright tests
test("student can complete quiz", async ({ page }) => {
  await page.goto("/student/courses/cs101");
  await page.click("text=Midterm Quiz");
  await page.click("text=Start Quiz");

  // Answer questions
  for (let i = 0; i < 10; i++) {
    await page.click(`input[value="option_${i}_correct"]`);
    await page.click("text=Next");
  }

  await page.click("text=Submit Quiz");
  await expect(page.locator("text=Your score")).toBeVisible();
});
```

---

## Appendix A: Type Definitions

### Assessment Types

```typescript
// frontend/src/features/assessment/types.ts

export type QuestionType =
  | "MCQ"
  | "MRQ"
  | "TRUE_FALSE"
  | "TEXT"
  | "ESSAY"
  | "ORDERING"
  | "MATRIX"
  | "FILL_BLANK"
  | "LIKERT";

export interface Assessment {
  id: string;
  tenant_id: string;
  title: string;
  description?: string;
  type: "quiz" | "exam" | "survey" | "form";
  status: "draft" | "published" | "closed";
  settings: AssessmentSettings;
  questions: AssessmentQuestion[];
  created_by_id: string;
  created_at: string;
  updated_at: string;
}

export interface AssessmentSettings {
  time_limit_minutes?: number;
  max_attempts: number;
  shuffle_questions: boolean;
  shuffle_answers: boolean;
  show_one_at_a_time: boolean;
  allow_backtrack: boolean;
  show_correct_answers: "immediately" | "after_due" | "never";
  require_proctoring: boolean;
  available_from?: string;
  available_until?: string;
  passing_score?: number;
}

export interface AssessmentQuestion {
  id: string;
  order: number;
  type: QuestionType;
  stem: string;
  stem_format: "plain" | "html" | "markdown";
  options?: QuestionOption[];
  correct_answer?: any;
  points: number;
  feedback_correct?: string;
  feedback_incorrect?: string;
  metadata?: Record<string, any>;
}

export interface QuestionOption {
  id: string;
  text: string;
  is_correct: boolean;
  points?: number;
  feedback?: string;
}

export interface AssessmentAttempt {
  id: string;
  assessment_id: string;
  student_id: string;
  status: "in_progress" | "submitted" | "graded";
  started_at: string;
  submitted_at?: string;
  graded_at?: string;
  time_spent_seconds: number;
  score?: number;
  max_score: number;
  percentage?: number;
  answers: AttemptAnswer[];
}

export interface AttemptAnswer {
  question_id: string;
  answer: any; // Type depends on question type
  is_correct?: boolean;
  points_earned?: number;
  feedback?: string;
  answered_at: string;
}
```

---

## Appendix B: API Client Functions

```typescript
// frontend/src/features/assessment/api.ts

import { api } from "@/api/client";
import type {
  Assessment,
  AssessmentAttempt,
  AssessmentQuestion,
} from "./types";

// Assessments CRUD
export const createAssessment = (data: Partial<Assessment>) =>
  api.post<Assessment>("/assessments", data);

export const getAssessment = (id: string) =>
  api.get<Assessment>(`/assessments/${id}`);

export const updateAssessment = (id: string, data: Partial<Assessment>) =>
  api.put<Assessment>(`/assessments/${id}`, data);

export const deleteAssessment = (id: string) =>
  api.delete(`/assessments/${id}`);

export const publishAssessment = (id: string) =>
  api.post<Assessment>(`/assessments/${id}/publish`);

// Student-facing
export const getAvailableAssessments = () =>
  api.get<Assessment[]>("/student/assessments");

export const getAssessmentForTaking = (id: string) =>
  api.get<Assessment>(`/student/assessments/${id}`);

// Attempts
export const startAttempt = (assessmentId: string) =>
  api.post<AssessmentAttempt>(`/assessments/${assessmentId}/attempts`);

export const saveAnswer = (
  attemptId: string,
  questionId: string,
  answer: any
) =>
  api.post(`/assessments/attempts/${attemptId}/answers/${questionId}`, {
    answer,
  });

export const submitAttempt = (attemptId: string) =>
  api.post<AssessmentAttempt>(`/assessments/attempts/${attemptId}/submit`);

export const getAttemptResults = (attemptId: string) =>
  api.get<AssessmentAttempt>(`/assessments/attempts/${attemptId}/results`);

// Proctoring
export const sendProctorEvent = (attemptId: string, event: ProctorEvent) =>
  api.post(`/assessments/attempts/${attemptId}/proctor-event`, event);
```

---

## Appendix C: Migration Checklist

### Before Starting Phase 1

- [ ] Verify backend assessment endpoints are working
- [ ] Check database tables exist (assessments, attempts, etc.)
- [ ] Ensure S3 configured for file uploads
- [ ] KaTeX installed for math rendering

### Phase 1 Completion Criteria

- [ ] Student can view course detail with modules
- [ ] Student can view assignment and submit work
- [ ] Student can start and complete a quiz
- [ ] Quiz results display correctly
- [ ] Quiz builder creates valid assessments

### Phase 2 Completion Criteria

- [ ] Survey builder functional with all question types
- [ ] Form builder with dictionary integration
- [ ] Teacher can view student progress tracker
- [ ] At-risk students identified correctly

### Phase 3 Completion Criteria

- [ ] Rich text editor working in question editor
- [ ] Math/LaTeX renders in questions
- [ ] Discussion forums functional (UI connected to existing backend)
- [ ] Attendance tracking working (UI connected to existing backend)
- [ ] GradebookPage implemented
- [ ] TranscriptView for students implemented

### Phase 4 Completion Criteria

- [ ] Gamification database tables created
- [ ] Gamification service and handler implemented
- [ ] XP awarded on key actions (login, submission, quiz)
- [ ] Badges system working
- [ ] Leaderboard displaying correctly
- [ ] Global search (Cmd+K) finding all entities
- [ ] AI tools panel accessible to admins
- [ ] LTI tools can be configured and launched

---

## Final Roadmap Summary

### Total Effort Estimate

| Phase   | Duration  | Backend Work | Frontend Work | Priority    |
| ------- | --------- | ------------ | ------------- | ----------- |
| Phase 1 | 1-2 weeks | ~20%         | ~80%          | 🔴 Critical |
| Phase 2 | 2-3 weeks | ~10%         | ~90%          | 🟠 High     |
| Phase 3 | 2-3 weeks | ~5%          | ~95%          | 🟡 Medium   |
| Phase 4 | 3-4 weeks | ~60%         | ~40%          | 🟢 Low      |

### Quick Wins (Can be done in 1-2 days each)

1. **Forums UI** - Backend 100% ready, just needs React components
2. **Attendance UI** - Backend 100% ready, just needs React components
3. **Global Search** - Backend 100% ready, just needs Command palette UI
4. **AI Tools Panel** - Backend 100% ready, just needs admin UI
5. **LTI Admin** - Backend 100% ready, just needs admin UI

### Requires Backend Work

1. **Student Course Detail** - May need new endpoint or verify existing
2. **Student Assignment Submit** - Needs new endpoint
3. **Teacher At-Risk** - Needs new endpoint
4. **Gamification** - Full new module (backend + frontend)

### Success Metrics

After completing all phases:

| Metric              | Current | Target |
| ------------------- | ------- | ------ |
| Backend Coverage    | 92%     | 100%   |
| Frontend Coverage   | 75%     | 100%   |
| Overall Integration | 78%     | 100%   |
| Student Experience  | 60%     | 100%   |
| Teacher Experience  | 90%     | 100%   |
| Admin Experience    | 95%     | 100%   |
| Engagement Features | 0%      | 100%   |

---

**Document maintained by:** Development Team  
**Last audit:** January 4, 2026  
**Next review:** After Phase 1 completion

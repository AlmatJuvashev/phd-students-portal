# Remaining Features Implementation Guide

> **Document Version:** 1.0  
> **Created:** January 4, 2026  
> **Purpose:** Detailed implementation guide for remaining features from v11 and backend

---

## Executive Summary

Based on the audit, **36% of functionality** (37 features) is not yet implemented. This document provides detailed implementation plans organized into 4 phases with estimated timelines.

### Implementation Priority Matrix

| Priority | Features | Complexity | Timeline |
|----------|----------|------------|----------|
| 🔴 Critical | Student Course/Assignment, Quiz Builder | High | 2-3 weeks |
| 🟠 High | Survey/Form Builder, Teacher Tracker | Medium-High | 2-3 weeks |
| 🟡 Medium | Rich Editor, Forums, Attendance | Medium | 3-4 weeks |
| 🟢 Low | Gamification, AI Tools, LTI | Variable | As needed |

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

| Type | v11 | Backend | Current | Action |
|------|-----|---------|---------|--------|
| Multiple Choice (MCQ) | ✅ | ✅ | ✅ | OK |
| Multiple Response (MRQ) | ✅ | ✅ | ✅ | OK |
| True/False | ✅ | ✅ | ❌ | Add |
| Short Text | ✅ | ✅ | ✅ | OK |
| Long Text (Essay) | ✅ | ✅ | ❌ | Add |
| Ordering | ✅ | ❌ | ✅ | Backend |
| Matrix/Grid | ✅ | ❌ | ❌ | Add |
| Fill in Blank | ✅ | ❌ | ❌ | Add |
| Math/LaTeX | ✅ | ❌ | ❌ | Add |
| Section Header | ✅ | ❌ | ❌ | Add |

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
  time_limit_minutes?: number;      // null = unlimited
  late_submission: 'allow' | 'penalize' | 'block';
  late_penalty_percent?: number;
  
  // Attempts
  max_attempts: number;             // 0 = unlimited
  attempt_grading: 'highest' | 'latest' | 'average';
  
  // Display
  shuffle_questions: boolean;
  shuffle_answers: boolean;
  show_one_at_a_time: boolean;
  allow_backtrack: boolean;
  
  // Feedback
  show_correct_answers: 'immediately' | 'after_due' | 'never';
  show_points_per_question: boolean;
  
  // Proctoring
  require_proctoring: boolean;
  lock_browser: boolean;
  webcam_required: boolean;
  
  // Availability
  available_from?: string;          // ISO date
  available_until?: string;
  
  // Grading
  grading_schema_id?: string;
  passing_score?: number;           // Percentage
}
```

#### Question Editor Features

```typescript
interface QuizQuestion {
  id: string;
  type: QuestionType;
  stem: string;                     // Question text (supports HTML/Markdown)
  stem_format: 'plain' | 'html' | 'markdown' | 'latex';
  
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
  blank_answers?: string[];         // Acceptable answers
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
  difficulty?: 'easy' | 'medium' | 'hard';
  tags?: string[];
  source_bank_id?: string;          // If imported from bank
}

interface QuestionOption {
  id: string;
  text: string;
  is_correct: boolean;
  points?: number;                  // For partial credit
  feedback?: string;                // Option-specific feedback
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
import 'katex/dist/katex.min.css';
import { InlineMath, BlockMath } from 'react-katex';

export const MathRenderer: React.FC<{ content: string }> = ({ content }) => {
  // Parse content for LaTeX delimiters
  // $...$ for inline, $$...$$ for block
  return (
    <div>
      {parseLatex(content).map((part, i) => 
        part.type === 'latex-inline' ? (
          <InlineMath key={i} math={part.content} />
        ) : part.type === 'latex-block' ? (
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

| Type | Description | Status |
|------|-------------|--------|
| Star Rating | 1-5 stars | ✅ Exists |
| NPS (0-10) | Net Promoter Score | ✅ Exists |
| Likert Matrix | Agreement scale grid | ✅ Exists |
| Open Text | Long text feedback | ✅ Exists |
| Dropdown | Single select dropdown | ❌ Add |
| Ranking | Rank items in order | ❌ Add |
| Slider | Numeric slider | ❌ Add |
| Date/Time | Date picker | ❌ Add |
| File Upload | Attachment | ❌ Add |

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
  condition: 'equals' | 'not_equals' | 'contains' | 'greater_than';
  value: any;
  then_show_question_id: string;
}
```

### 2.2 Form Builder

**Current State:** Basic modal exists  
**Location:** `frontend/src/features/studio/components/FormBuilderModal.tsx`

#### Form Field Types

| Field Type | Description | Validation |
|------------|-------------|------------|
| Text | Single line | Required, min/max length, regex |
| Textarea | Multi-line | Required, min/max length |
| Number | Numeric input | Required, min/max, step |
| Email | Email address | Format validation |
| Phone | Phone number | Format validation |
| Date | Date picker | Min/max date |
| DateTime | Date + time | - |
| Select | Dropdown | Required |
| Multi-Select | Multiple choice | Min/max selections |
| Checkbox | Boolean | - |
| Radio | Single choice | Required |
| File | File upload | Types, max size |
| Signature | Signature pad | Required |
| Dictionary | Link to dictionary | Auto-populate options |

#### Dictionary Integration

```tsx
// Connect form fields to dictionaries
<FormFieldEditor
  field={field}
  onChange={updateField}
>
  {field.type === 'select' && (
    <DictionaryPicker
      label="Options Source"
      value={field.dictionary_id}
      onChange={(dictId) => updateField({ 
        ...field, 
        dictionary_id: dictId,
        options: 'from_dictionary' 
      })}
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
  overall_progress: number;         // 0-100%
  assignments_completed: number;
  assignments_total: number;
  assignments_overdue: number;
  last_activity: string;            // ISO date
  days_inactive: number;
  average_grade: number;
  
  // Risk assessment
  risk_level: 'low' | 'medium' | 'high' | 'critical';
  risk_factors: string[];           // ["3 overdue assignments", "14 days inactive"]
  
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

| Feature | Description | Implementation |
|---------|-------------|----------------|
| Rich Text | Bold, italic, lists, links | TipTap or Lexical |
| Math | LaTeX equations | KaTeX |
| Images | Inline images | Upload to S3 |
| Audio | Audio clips | HTML5 Audio |
| Video | Embedded video | iframe/player |
| Code | Syntax highlighting | Prism.js |
| Tables | Data tables | TipTap table extension |

#### Implementation

```tsx
// QuestionEditor.tsx
import { useEditor, EditorContent } from '@tiptap/react';
import StarterKit from '@tiptap/starter-kit';
import Mathematics from '@tiptap/extension-mathematics';
import Image from '@tiptap/extension-image';

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
        {students.map(student => (
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

## Phase 4: Nice-to-Have Features (Low Priority)

**Timeline:** As needed  
**Goal:** Advanced features for enhanced experience

### 4.1 Gamification System

**Backend Required:** Yes (new tables and logic)

#### Database Schema

```sql
CREATE TABLE user_xp (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL,
    total_xp INT DEFAULT 0,
    level INT DEFAULT 1,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE badges (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    name VARCHAR(255),
    description TEXT,
    icon_url VARCHAR(500),
    criteria JSONB,
    xp_reward INT DEFAULT 0
);

CREATE TABLE user_badges (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    badge_id UUID NOT NULL,
    earned_at TIMESTAMP
);

CREATE TABLE xp_events (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    event_type VARCHAR(50),  -- 'assignment_complete', 'quiz_perfect', etc.
    xp_amount INT,
    metadata JSONB,
    created_at TIMESTAMP
);
```

#### XP Event Types

| Event | XP | Description |
|-------|-----|-------------|
| Assignment Submitted | 10 | Submit any assignment |
| Assignment On Time | 5 | Submit before deadline |
| Quiz Completed | 10 | Complete any quiz |
| Quiz Perfect Score | 25 | 100% on quiz |
| Course Completed | 100 | Finish all course items |
| Daily Login Streak | 5 | Login daily (multiplied by streak) |
| Forum Post | 5 | Post in forums |
| Helpful Post | 10 | Post marked as helpful |

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
          {results.courses.map(course => (
            <CommandItem onSelect={() => navigate(`/courses/${course.id}`)}>
              <BookOpen className="mr-2" />
              {course.title}
            </CommandItem>
          ))}
        </CommandGroup>
        
        <CommandGroup heading="Students">
          {results.students.map(student => (
            <CommandItem>
              <User className="mr-2" />
              {student.name}
            </CommandItem>
          ))}
        </CommandGroup>
        
        <CommandGroup heading="Assignments">
          {results.assignments.map(assignment => (
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

### Critical (Need for Phase 1)

| Endpoint | Handler | Status |
|----------|---------|--------|
| `POST /api/assessments` | assessment_handler.go | ⚠️ Not Implemented |
| `GET /api/assessments/:id` | assessment_handler.go | ⚠️ Not Implemented |
| `PUT /api/assessments/:id` | assessment_handler.go | ⚠️ Not Implemented |
| `GET /api/student/courses/:id` | student_handler.go | Need to verify |
| `GET /api/student/courses/:id/modules` | course_content_handler.go | Need to verify |
| `POST /api/student/assignments/:id/submit` | - | ⚠️ Needs handler |

### High Priority (Phase 2)

| Endpoint | Handler | Status |
|----------|---------|--------|
| `GET /api/teacher/courses/:id/students` | teacher_handler.go | Need to verify |
| `GET /api/teacher/courses/:id/at-risk` | - | ⚠️ New endpoint |
| `POST /api/surveys` | - | ⚠️ Similar to assessments |
| `POST /api/forms` | - | ⚠️ Similar to assessments |

### Medium Priority (Phase 3)

Most endpoints already exist - just need UI.

---

## Testing Strategy

### Unit Tests

```typescript
// Assessment taking
describe('AssessmentTaking', () => {
  it('should start an attempt when component mounts');
  it('should save answer when user selects option');
  it('should auto-save periodically');
  it('should show timer and warn when low');
  it('should submit all answers on completion');
  it('should handle network errors gracefully');
});

// Quiz Builder
describe('QuizBuilder', () => {
  it('should add new question of each type');
  it('should reorder questions via drag and drop');
  it('should import questions from bank');
  it('should validate required fields');
  it('should save quiz to backend');
});
```

### E2E Tests

```typescript
// Playwright tests
test('student can complete quiz', async ({ page }) => {
  await page.goto('/student/courses/cs101');
  await page.click('text=Midterm Quiz');
  await page.click('text=Start Quiz');
  
  // Answer questions
  for (let i = 0; i < 10; i++) {
    await page.click(`input[value="option_${i}_correct"]`);
    await page.click('text=Next');
  }
  
  await page.click('text=Submit Quiz');
  await expect(page.locator('text=Your score')).toBeVisible();
});
```

---

## Appendix A: Type Definitions

### Assessment Types

```typescript
// frontend/src/features/assessment/types.ts

export type QuestionType = 
  | 'MCQ' 
  | 'MRQ' 
  | 'TRUE_FALSE' 
  | 'TEXT' 
  | 'ESSAY'
  | 'ORDERING'
  | 'MATRIX'
  | 'FILL_BLANK'
  | 'LIKERT';

export interface Assessment {
  id: string;
  tenant_id: string;
  title: string;
  description?: string;
  type: 'quiz' | 'exam' | 'survey' | 'form';
  status: 'draft' | 'published' | 'closed';
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
  show_correct_answers: 'immediately' | 'after_due' | 'never';
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
  stem_format: 'plain' | 'html' | 'markdown';
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
  status: 'in_progress' | 'submitted' | 'graded';
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
  answer: any;               // Type depends on question type
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

import { api } from '@/api/client';
import type { Assessment, AssessmentAttempt, AssessmentQuestion } from './types';

// Assessments CRUD
export const createAssessment = (data: Partial<Assessment>) =>
  api.post<Assessment>('/assessments', data);

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
  api.get<Assessment[]>('/student/assessments');

export const getAssessmentForTaking = (id: string) =>
  api.get<Assessment>(`/student/assessments/${id}`);

// Attempts
export const startAttempt = (assessmentId: string) =>
  api.post<AssessmentAttempt>(`/assessments/${assessmentId}/attempts`);

export const saveAnswer = (attemptId: string, questionId: string, answer: any) =>
  api.post(`/assessments/attempts/${attemptId}/answers/${questionId}`, { answer });

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
- [ ] Discussion forums functional
- [ ] Attendance tracking working

### Phase 4 Completion Criteria

- [ ] Gamification points accumulating
- [ ] Global search finding all entities
- [ ] AI generation producing useful content
- [ ] LTI tools can be configured

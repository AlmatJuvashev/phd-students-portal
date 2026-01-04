# Frontend Migration & Integration Guide

> **Document Version:** 2.0  
> **Updated:** January 3, 2026  
> **Purpose:** Guide for migrating `ui-examples/phd-journey-tracker_v11` components into the working `frontend/` application and connecting them to backend APIs

---

## Executive Summary

Your current **`frontend/`** is a working React application with:

- ✅ Authentication via HttpOnly cookies (`AuthContext`)
- ✅ Multi-tenant support (`X-Tenant-Slug` header)
- ✅ Role-based routing (`ProtectedRoute`, `ServiceProtectedRoute`)
- ✅ React Query for data fetching
- ✅ Feature-based architecture (`features/`)

The **`ui-examples/phd-journey-tracker_v11`** has rich UI components for additional features that use **mock data**. This guide shows how to migrate those components into your working frontend and connect them to the real backend APIs.

---

## 📁 Current Frontend Architecture

```
frontend/src/
├── api/                    # API clients (✅ Use this pattern)
│   ├── client.ts           # Base HTTP client with auth cookies
│   ├── admin.ts            # Admin APIs
│   ├── journey.ts          # Journey/Progress APIs
│   ├── user.ts             # User APIs
│   ├── contacts.ts         # Contacts APIs
│   └── dictionaries.ts     # Dictionary APIs
├── contexts/
│   ├── AuthContext.tsx     # ✅ Auth state (cookie-based)
│   └── TenantServicesContext.tsx  # ✅ Optional services toggle
├── features/               # ✅ Feature modules (preferred location for new features)
│   ├── calendar/           # ✅ Already integrated
│   ├── scheduler/          # ✅ Already integrated
│   ├── analytics/          # ✅ Already integrated
│   ├── chat/               # ✅ Already integrated
│   └── ...
├── routes/index.tsx        # ✅ Central routing with lazy loading
└── layouts/
    ├── AdminLayout.tsx     # ✅ Admin panel layout
    └── SuperadminLayout.tsx
```

### Key Conventions (Follow These)

1. **API Calls** - Use `api/client.ts` functions (`api.get()`, `api.post()`)
2. **Auth** - Use `useAuth()` hook from `AuthContext` (cookie-based, NOT localStorage)
3. **Data Fetching** - Use React Query (`useQuery`, `useMutation`)
4. **New Features** - Add to `features/` directory with `api.ts`, `types.ts`, components
5. **Routes** - Lazy load pages in `routes/index.tsx`
6. **Service Gates** - Wrap optional features in `<ServiceProtectedRoute>`

---

## 📊 Feature Comparison: What to Migrate

| Feature              | Current Frontend                  | v11 Example                     | Backend API                   | Action                |
| -------------------- | --------------------------------- | ------------------------------- | ----------------------------- | --------------------- |
| **Auth/Login**       | ✅ Cookie-based                   | 🔴 N/A                          | ✅ `/api/auth/*`              | No change             |
| **Journey Map**      | ✅ Working                        | ✅ Enhanced UI                  | ✅ `/api/journey/*`           | Optional UI upgrade   |
| **Chat**             | ✅ Working                        | ✅ Different UI                 | ✅ `/api/chat/*`              | No change             |
| **Calendar**         | ✅ `features/calendar/`           | ✅ Different UI                 | ✅ `/api/calendar/*`          | No change             |
| **Scheduler**        | ✅ `features/scheduler/`          | ✅ Enhanced                     | ✅ `/api/scheduler/*`         | Merge UI improvements |
| **Students Monitor** | ✅ `features/students-monitor/`   | ❌ N/A                          | ✅ `/api/admin/*`             | No change             |
| **Dictionaries**     | ✅ `features/admin/dictionaries/` | ❌ N/A                          | ✅ `/api/admin/*`             | No change             |
| **LMS: Programs**    | ❌ Missing                        | ✅ `ProgramsPage`               | ✅ `/api/curriculum/*`        | **Migrate**           |
| **LMS: Courses**     | ❌ Missing                        | ✅ `GlobalCoursesPage`          | ✅ `/api/curriculum/*`        | **Migrate**           |
| **Course Builder**   | ❌ Missing                        | ✅ `CourseBuilder`              | ✅ `/api/course-content/*`    | **Migrate**           |
| **Enrollments**      | ❌ Missing                        | ✅ `EnrollmentsPage`            | ✅ `/api/admin/enrollments/*` | **Migrate**           |
| **Item Bank**        | ❌ Missing                        | ✅ `BanksPage`, `QuestionsPage` | ✅ `/api/item-bank/*`         | **Migrate**           |
| **Quiz Builder**     | ❌ Missing                        | ✅ `QuizBuilder`                | ✅ `/api/assessments/*`       | **Migrate**           |
| **Grading**          | ❌ Missing                        | ✅ `TeacherGradingPage`         | ✅ `/api/grading/*`           | **Migrate**           |
| **Student App**      | ❌ Missing                        | ✅ Full student portal          | ✅ Various APIs               | **Migrate**           |
| **Teacher App**      | ❌ Missing                        | ✅ Full teacher portal          | ✅ Various APIs               | **Migrate**           |

---

## 🚀 Migration Plan by Priority

### Phase 1: LMS Core (Week 1-2)

#### 1.1 Create `features/curriculum/` Module

```
frontend/src/features/curriculum/
├── api.ts              # API calls to /api/curriculum/*
├── types.ts            # TypeScript interfaces
├── ProgramsPage.tsx    # Migrate from v11
├── ProgramDetailPage.tsx
├── CoursesPage.tsx     # Migrate GlobalCoursesPage
├── CourseBuilderPage.tsx
└── components/
    ├── ProgramCard.tsx
    └── CourseCard.tsx
```

**Create `frontend/src/features/curriculum/api.ts`:**

```typescript
import { api } from "@/api/client";
import { Program, Course } from "./types";

// Programs
export const getPrograms = () => api.get<Program[]>("/curriculum/programs");
export const getProgram = (id: string) =>
  api.get<Program>(`/curriculum/programs/${id}`);
export const createProgram = (data: Partial<Program>) =>
  api.post("/curriculum/programs", data);
export const updateProgram = (id: string, data: Partial<Program>) =>
  api.put(`/curriculum/programs/${id}`, data);
export const deleteProgram = (id: string) =>
  api.delete(`/curriculum/programs/${id}`);

// Courses
export const getCourses = (programId?: string) =>
  api.get<Course[]>(
    `/curriculum/courses${programId ? `?program_id=${programId}` : ""}`
  );
export const getCourse = (id: string) =>
  api.get<Course>(`/curriculum/courses/${id}`);
export const createCourse = (data: Partial<Course>) =>
  api.post("/curriculum/courses", data);
export const updateCourse = (id: string, data: Partial<Course>) =>
  api.put(`/curriculum/courses/${id}`, data);
export const deleteCourse = (id: string) =>
  api.delete(`/curriculum/courses/${id}`);
```

**Create `frontend/src/features/curriculum/types.ts`:**

```typescript
export interface Program {
  id: string;
  tenant_id: string;
  title: string;
  code: string;
  description: string;
  type: "bachelor" | "master" | "doctoral" | "certificate";
  total_credits: number;
  duration_semesters: number;
  status: "draft" | "active" | "archived";
  created_at: string;
  updated_at: string;
}

export interface Course {
  id: string;
  tenant_id: string;
  program_id?: string;
  code: string;
  title: string;
  description: string;
  credits: number;
  category: "core" | "elective" | "research";
  prerequisites?: string[];
  status: "draft" | "active" | "archived";
  created_at: string;
  updated_at: string;
}
```

**Migrate v11 `ProgramsPage.tsx` → `frontend/src/features/curriculum/ProgramsPage.tsx`:**

```typescript
// BEFORE (v11 with mock data):
import { getPrograms, Program } from "../../data/opsData";
useEffect(() => {
  setPrograms(getPrograms());
}, []);

// AFTER (with real API):
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { getPrograms, createProgram, deleteProgram } from "./api";
import { Program } from "./types";

export function ProgramsPage() {
  const queryClient = useQueryClient();

  const {
    data: programs = [],
    isLoading,
    error,
  } = useQuery({
    queryKey: ["programs"],
    queryFn: getPrograms,
  });

  const createMutation = useMutation({
    mutationFn: createProgram,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["programs"] }),
  });

  const deleteMutation = useMutation({
    mutationFn: deleteProgram,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["programs"] }),
  });

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error loading programs</div>;

  // ... rest of UI from v11 (keep the JSX layout)
}
```

#### 1.2 Add Routes

**Update `frontend/src/routes/index.tsx`:**

```typescript
// Add lazy imports at top
const ProgramsPage = lazy(() =>
  import("@/features/curriculum/ProgramsPage").then((m) => ({ default: m.ProgramsPage }))
);
const ProgramDetailPage = lazy(() =>
  import("@/features/curriculum/ProgramDetailPage").then((m) => ({ default: m.ProgramDetailPage }))
);
const CoursesPage = lazy(() =>
  import("@/features/curriculum/CoursesPage").then((m) => ({ default: m.CoursesPage }))
);
const CourseBuilderPage = lazy(() =>
  import("@/features/curriculum/CourseBuilderPage").then((m) => ({ default: m.CourseBuilderPage }))
);

// Add to admin children array:
{
  path: "programs",
  element: (
    <ProtectedRoute requiredAnyRole={["admin"]}>
      {WithSuspense(<ProgramsPage />)}
    </ProtectedRoute>
  ),
},
{
  path: "programs/:id",
  element: (
    <ProtectedRoute requiredAnyRole={["admin"]}>
      {WithSuspense(<ProgramDetailPage />)}
    </ProtectedRoute>
  ),
},
{
  path: "courses",
  element: (
    <ProtectedRoute requiredAnyRole={["admin"]}>
      {WithSuspense(<CoursesPage />)}
    </ProtectedRoute>
  ),
},
{
  path: "courses/:id/builder",
  element: (
    <ProtectedRoute requiredAnyRole={["admin"]}>
      {WithSuspense(<CourseBuilderPage />)}
    </ProtectedRoute>
  ),
},
```

---

### Phase 2: Enrollments & Course Content (Week 2)

#### 2.1 Create `features/enrollments/` Module

```
frontend/src/features/enrollments/
├── api.ts
├── types.ts
├── EnrollmentsPage.tsx    # Migrate from v11
└── components/
    └── EnrollmentTable.tsx
```

**Create `frontend/src/features/enrollments/api.ts`:**

```typescript
import { api } from "@/api/client";
import { Enrollment, EnrollmentCreateRequest } from "./types";

export const getEnrollments = (filters?: {
  course_id?: string;
  student_id?: string;
}) => {
  const params = new URLSearchParams();
  if (filters?.course_id) params.append("course_id", filters.course_id);
  if (filters?.student_id) params.append("student_id", filters.student_id);
  return api.get<Enrollment[]>(`/admin/enrollments?${params}`);
};

export const createEnrollment = (data: EnrollmentCreateRequest) =>
  api.post("/admin/enrollments", data);

export const bulkEnroll = (data: {
  course_id: string;
  student_ids: string[];
}) => api.post("/admin/enrollments/bulk", data);

export const dropEnrollment = (id: string) =>
  api.delete(`/admin/enrollments/${id}`);
```

#### 2.2 Create `features/course-content/` Module

```
frontend/src/features/course-content/
├── api.ts
├── types.ts
├── CourseBuilder.tsx      # Migrate from v11
├── ModuleEditor.tsx
├── LessonEditor.tsx
└── ActivityEditor.tsx
```

**Create `frontend/src/features/course-content/api.ts`:**

```typescript
import { api } from "@/api/client";
import { CourseModule, CourseLesson, CourseActivity } from "./types";

// Modules
export const getModules = (courseId: string) =>
  api.get<CourseModule[]>(`/course-content/modules?course_id=${courseId}`);
export const createModule = (data: Partial<CourseModule>) =>
  api.post("/course-content/modules", data);
export const updateModule = (id: string, data: Partial<CourseModule>) =>
  api.put(`/course-content/modules/${id}`, data);
export const reorderModules = (courseId: string, moduleIds: string[]) =>
  api.post(`/course-content/modules/reorder`, {
    course_id: courseId,
    module_ids: moduleIds,
  });

// Lessons
export const getLessons = (moduleId: string) =>
  api.get<CourseLesson[]>(`/course-content/lessons?module_id=${moduleId}`);
export const createLesson = (data: Partial<CourseLesson>) =>
  api.post("/course-content/lessons", data);
export const updateLesson = (id: string, data: Partial<CourseLesson>) =>
  api.put(`/course-content/lessons/${id}`, data);

// Activities
export const getActivities = (lessonId: string) =>
  api.get<CourseActivity[]>(`/course-content/activities?lesson_id=${lessonId}`);
export const createActivity = (data: Partial<CourseActivity>) =>
  api.post("/course-content/activities", data);
```

---

### Phase 3: Item Bank & Assessment (Week 3)

#### 3.1 Create `features/item-bank/` Module

```
frontend/src/features/item-bank/
├── api.ts
├── types.ts
├── BanksPage.tsx          # Migrate from v11
├── QuestionsPage.tsx      # Migrate from v11
├── QuestionEditor.tsx     # Migrate from v11
└── components/
    ├── QuestionCard.tsx
    └── AnswerOptions.tsx
```

**Create `frontend/src/features/item-bank/api.ts`:**

```typescript
import { api } from "@/api/client";
import { QuestionBank, Question, QuestionCreateRequest } from "./types";

// Banks
export const getBanks = () => api.get<QuestionBank[]>("/item-bank/banks");
export const getBank = (id: string) =>
  api.get<QuestionBank>(`/item-bank/banks/${id}`);
export const createBank = (data: Partial<QuestionBank>) =>
  api.post("/item-bank/banks", data);
export const updateBank = (id: string, data: Partial<QuestionBank>) =>
  api.put(`/item-bank/banks/${id}`, data);
export const deleteBank = (id: string) => api.delete(`/item-bank/banks/${id}`);

// Questions
export const getQuestions = (
  bankId: string,
  filters?: { type?: string; difficulty?: string }
) => {
  const params = new URLSearchParams({ bank_id: bankId });
  if (filters?.type) params.append("type", filters.type);
  if (filters?.difficulty) params.append("difficulty", filters.difficulty);
  return api.get<Question[]>(`/item-bank/questions?${params}`);
};
export const getQuestion = (id: string) =>
  api.get<Question>(`/item-bank/questions/${id}`);
export const createQuestion = (data: QuestionCreateRequest) =>
  api.post("/item-bank/questions", data);
export const updateQuestion = (id: string, data: Partial<Question>) =>
  api.put(`/item-bank/questions/${id}`, data);
export const deleteQuestion = (id: string) =>
  api.delete(`/item-bank/questions/${id}`);

// Bulk import
export const importQuestions = (bankId: string, file: File) => {
  const formData = new FormData();
  formData.append("file", file);
  return api.postFormData(`/item-bank/banks/${bankId}/import`, formData);
};
```

#### 3.2 Create `features/assessment/` Module

```
frontend/src/features/assessment/
├── api.ts
├── types.ts
├── QuizBuilder.tsx        # Migrate from v11
├── QuizPreview.tsx        # Migrate from v11
├── SurveyBuilder.tsx      # Migrate from v11
└── components/
    ├── QuestionPicker.tsx
    └── QuizSettings.tsx
```

**Create `frontend/src/features/assessment/api.ts`:**

```typescript
import { api } from "@/api/client";
import { Assessment, AssessmentAttempt, AssessmentSubmission } from "./types";

// Assessments (quizzes, exams, surveys)
export const getAssessments = (courseId?: string) => {
  const params = courseId ? `?course_id=${courseId}` : "";
  return api.get<Assessment[]>(`/assessments${params}`);
};
export const getAssessment = (id: string) =>
  api.get<Assessment>(`/assessments/${id}`);
export const createAssessment = (data: Partial<Assessment>) =>
  api.post("/assessments", data);
export const updateAssessment = (id: string, data: Partial<Assessment>) =>
  api.put(`/assessments/${id}`, data);
export const publishAssessment = (id: string) =>
  api.post(`/assessments/${id}/publish`);

// Student attempts
export const startAttempt = (assessmentId: string) =>
  api.post<AssessmentAttempt>(`/assessments/${assessmentId}/attempts`);
export const submitAnswer = (
  attemptId: string,
  questionId: string,
  answer: any
) =>
  api.post(`/assessments/attempts/${attemptId}/answers`, {
    question_id: questionId,
    answer,
  });
export const finishAttempt = (attemptId: string) =>
  api.post(`/assessments/attempts/${attemptId}/finish`);
export const getAttemptResult = (attemptId: string) =>
  api.get(`/assessments/attempts/${attemptId}/result`);
```

---

### Phase 4: Grading System (Week 3-4)

#### 4.1 Create `features/grading/` Module

```
frontend/src/features/grading/
├── api.ts
├── types.ts
├── TeacherGradingPage.tsx # Migrate from v11
├── StudentGradesPage.tsx  # Migrate from v11
├── GradebookPage.tsx
└── components/
    ├── SubmissionCard.tsx
    ├── GradeInput.tsx
    └── RubricGrader.tsx
```

**Create `frontend/src/features/grading/api.ts`:**

```typescript
import { api } from "@/api/client";
import { Submission, GradeEntry, Gradebook, GradingSchema } from "./types";

// Grading schemas
export const getGradingSchema = (courseId: string) =>
  api.get<GradingSchema>(`/grading/schema?course_id=${courseId}`);
export const updateGradingSchema = (
  courseId: string,
  schema: Partial<GradingSchema>
) => api.put(`/grading/schema/${courseId}`, schema);

// Submissions (for teachers)
export const getPendingSubmissions = (courseId: string) =>
  api.get<Submission[]>(
    `/grading/courses/${courseId}/submissions?status=pending`
  );
export const getSubmission = (id: string) =>
  api.get<Submission>(`/grading/submissions/${id}`);

// Grading
export const submitGrade = (entry: GradeEntry) =>
  api.post("/grading/entries", entry);
export const updateGrade = (entryId: string, data: Partial<GradeEntry>) =>
  api.put(`/grading/entries/${entryId}`, data);

// Gradebook
export const getGradebook = (courseId: string) =>
  api.get<Gradebook>(`/grading/gradebook/${courseId}`);
export const getStudentGrades = () => api.get("/grading/gradebook/my");

// Export
export const exportGradebook = (courseId: string, format: "csv" | "xlsx") =>
  api.get(`/grading/gradebook/${courseId}/export?format=${format}`, {
    responseType: "blob",
  });
```

---

### Phase 5: Student & Teacher Portals (Week 4-5)

#### 5.1 Create `features/student-portal/` Module

The v11 example has a complete student portal. Migrate these pages:

```
frontend/src/features/student-portal/
├── api.ts
├── types.ts
├── StudentDashboard.tsx    # Migrate from v11/student/pages/
├── StudentCourses.tsx
├── StudentCourseDetail.tsx
├── StudentAssignments.tsx
├── StudentGrades.tsx
└── layouts/
    └── StudentLayout.tsx   # Optional: different layout for students
```

**Create `frontend/src/features/student-portal/api.ts`:**

```typescript
import { api } from "@/api/client";
import { MyEnrollment, MyCourse, MyAssignment, MyGrade } from "./types";

// Dashboard
export const getStudentDashboard = () => api.get("/student/dashboard");

// Enrollments
export const getMyEnrollments = () =>
  api.get<MyEnrollment[]>("/lms/enrollments/my");
export const getMyEnrollment = (courseId: string) =>
  api.get<MyEnrollment>(`/lms/enrollments/my/${courseId}`);

// Course content
export const getMyCourseContent = (courseId: string) =>
  api.get<MyCourse>(`/lms/courses/${courseId}/content`);
export const markLessonComplete = (lessonId: string) =>
  api.post(`/lms/lessons/${lessonId}/complete`);

// Assignments
export const getMyAssignments = (
  status?: "pending" | "submitted" | "graded"
) => {
  const params = status ? `?status=${status}` : "";
  return api.get<MyAssignment[]>(`/student/assignments${params}`);
};
export const submitAssignment = (assignmentId: string, data: FormData) =>
  api.postFormData(`/student/assignments/${assignmentId}/submit`, data);

// Grades
export const getMyGrades = () => api.get<MyGrade[]>("/grading/gradebook/my");
export const getMyTranscript = () => api.get("/transcript/student/me");
```

#### 5.2 Create `features/teacher-portal/` Module

```
frontend/src/features/teacher-portal/
├── api.ts
├── types.ts
├── TeacherDashboard.tsx    # Migrate from v11/teacher/pages/
├── TeacherCourses.tsx
├── TeacherCourseDetail.tsx
├── TeacherStudentTracker.tsx
└── layouts/
    └── TeacherLayout.tsx
```

**Create `frontend/src/features/teacher-portal/api.ts`:**

```typescript
import { api } from "@/api/client";
import { TeacherCourse, CourseRoster, StudentProgress } from "./types";

// Dashboard
export const getTeacherDashboard = () => api.get("/teacher/dashboard");

// Courses
export const getMyCourses = () =>
  api.get<TeacherCourse[]>("/lms/courses/teaching");
export const getCourseRoster = (courseId: string) =>
  api.get<CourseRoster>(`/lms/courses/${courseId}/roster`);

// Student tracking
export const getStudentProgress = (courseId: string, studentId: string) =>
  api.get<StudentProgress>(
    `/lms/courses/${courseId}/students/${studentId}/progress`
  );
export const getCourseAnalytics = (courseId: string) =>
  api.get(`/analytics/courses/${courseId}`);

// Attendance
export const recordAttendance = (sessionId: string, attendees: string[]) =>
  api.post(`/attendance/sessions/${sessionId}`, { attendees });
```

---

## 🔄 Migration Checklist for Each v11 Component

When migrating any component from v11, follow this checklist:

### Step 1: Copy the Component File

```bash
cp ui-examples/phd-journey-tracker_v11/admin/pages/ops/ProgramsPage.tsx \
   frontend/src/features/curriculum/ProgramsPage.tsx
```

### Step 2: Update Imports

```typescript
// ❌ BEFORE (v11 imports)
import { getPrograms, Program } from "../../data/opsData";
import { Button } from "../../components/ui/Button";

// ✅ AFTER (current frontend imports)
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { getPrograms, createProgram, Program } from "./api";
import { Button } from "@/components/ui/button"; // shadcn/ui
```

### Step 3: Replace Mock Data with React Query

```typescript
// ❌ BEFORE (mock data)
const [programs, setPrograms] = useState<Program[]>([]);
useEffect(() => {
  setPrograms(getPrograms());
}, []);

// ✅ AFTER (React Query)
const {
  data: programs = [],
  isLoading,
  error,
  refetch,
} = useQuery({
  queryKey: ["programs"],
  queryFn: getPrograms,
});
```

### Step 4: Replace Local State Mutations with useMutation

```typescript
// ❌ BEFORE (local state)
const handleCreate = (data: Program) => {
  setPrograms([...programs, { ...data, id: crypto.randomUUID() }]);
};

// ✅ AFTER (mutation)
const queryClient = useQueryClient();

const createMutation = useMutation({
  mutationFn: createProgram,
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ["programs"] });
    toast.success("Program created successfully");
  },
  onError: (error) => {
    toast.error(error.message);
  },
});

const handleCreate = (data: Partial<Program>) => {
  createMutation.mutate(data);
};
```

### Step 5: Update Navigation

```typescript
// ❌ BEFORE (v11 onNavigate prop)
const { onNavigate } = props;
<Button onClick={() => onNavigate(`/admin/ops/programs/${id}`)}>View</Button>;

// ✅ AFTER (React Router)
import { useNavigate } from "react-router-dom";

const navigate = useNavigate();
<Button onClick={() => navigate(`/admin/programs/${id}`)}>View</Button>;
```

### Step 6: Add Loading and Error States

```typescript
if (isLoading) {
  return (
    <div className="flex items-center justify-center h-64">
      <Loader2 className="h-8 w-8 animate-spin" />
    </div>
  );
}

if (error) {
  return (
    <Alert variant="destructive">
      <AlertDescription>Failed to load data: {error.message}</AlertDescription>
    </Alert>
  );
}
```

### Step 7: Add to Routes

```typescript
// In frontend/src/routes/index.tsx
const ProgramsPage = lazy(() =>
  import("@/features/curriculum/ProgramsPage").then((m) => ({ default: m.ProgramsPage }))
);

// In admin children:
{
  path: "programs",
  element: (
    <ProtectedRoute requiredAnyRole={["admin"]}>
      {WithSuspense(<ProgramsPage />)}
    </ProtectedRoute>
  ),
},
```

---

## 📋 Files to Create Summary

### New API Files

```
frontend/src/features/
├── curriculum/api.ts        # Programs, Courses
├── course-content/api.ts    # Modules, Lessons, Activities
├── enrollments/api.ts       # Enrollments management
├── item-bank/api.ts         # Question banks, Questions
├── assessment/api.ts        # Quizzes, Exams, Surveys
├── grading/api.ts           # Grading, Gradebook
├── student-portal/api.ts    # Student-specific APIs
└── teacher-portal/api.ts    # Teacher-specific APIs
```

### New Route Groups to Add

```typescript
// Admin routes (add to /admin children)
/admin/programs           # ProgramsPage
/admin/programs/:id       # ProgramDetailPage
/admin/courses            # CoursesPage
/admin/courses/:id/builder # CourseBuilderPage
/admin/enrollments        # EnrollmentsPage
/admin/item-bank          # BanksPage
/admin/item-bank/:id      # QuestionsPage
/admin/assessments        # AssessmentsPage
/admin/grading            # GradingAdminPage

// Student routes (consider new /student layout)
/student/dashboard        # StudentDashboard
/student/courses          # StudentCourses
/student/courses/:id      # StudentCourseDetail
/student/assignments      # StudentAssignments
/student/grades           # StudentGrades

// Teacher routes (consider new /teacher layout)
/teacher/dashboard        # TeacherDashboard
/teacher/courses          # TeacherCourses
/teacher/courses/:id      # TeacherCourseDetail
/teacher/grading          # TeacherGradingPage
```

---

## ⚡ Quick Reference: API Endpoints

### Already Implemented in Backend

| Endpoint                           | Method           | Description             |
| ---------------------------------- | ---------------- | ----------------------- |
| `/api/curriculum/programs`         | GET, POST        | List/Create programs    |
| `/api/curriculum/programs/:id`     | GET, PUT, DELETE | CRUD program            |
| `/api/curriculum/courses`          | GET, POST        | List/Create courses     |
| `/api/curriculum/courses/:id`      | GET, PUT, DELETE | CRUD course             |
| `/api/course-content/modules`      | GET, POST        | List/Create modules     |
| `/api/course-content/lessons`      | GET, POST        | List/Create lessons     |
| `/api/course-content/activities`   | GET, POST        | List/Create activities  |
| `/api/admin/enrollments`           | GET, POST        | List/Create enrollments |
| `/api/item-bank/banks`             | GET, POST        | List/Create banks       |
| `/api/item-bank/questions`         | GET, POST        | List/Create questions   |
| `/api/assessments`                 | GET, POST        | List/Create assessments |
| `/api/assessments/:id/attempts`    | POST, GET        | Quiz attempts           |
| `/api/grading/schema`              | GET, PUT         | Grading schema          |
| `/api/grading/entries`             | POST             | Submit grades           |
| `/api/grading/gradebook/:courseId` | GET              | Course gradebook        |
| `/api/scheduler/terms`             | GET, POST        | Academic terms          |
| `/api/scheduler/rooms`             | GET, POST        | Rooms                   |
| `/api/scheduler/sessions`          | GET, POST        | Class sessions          |
| `/api/calendar/events`             | GET, POST        | Calendar events         |

---

## 🎯 Estimated Timeline

| Phase | Features                     | Duration  | Dependencies |
| ----- | ---------------------------- | --------- | ------------ |
| 1     | LMS Core (Programs, Courses) | 1 week    | None         |
| 2     | Enrollments, Course Content  | 1 week    | Phase 1      |
| 3     | Item Bank, Assessment        | 1 week    | Phase 2      |
| 4     | Grading System               | 1 week    | Phase 3      |
| 5     | Student & Teacher Portals    | 1-2 weeks | Phase 4      |

**Total: 5-6 weeks** for complete migration

---

## ✅ Success Criteria

After migration, verify:

1. **All v11 UI pages work** with real backend data
2. **No mock data remains** - all data from APIs
3. **Authentication works** - cookies, not localStorage
4. **Role-based access** - proper route protection
5. **CRUD operations** - create, read, update, delete all work
6. **Error handling** - proper error messages shown
7. **Loading states** - spinners during API calls

## 🏁 Next Steps: Getting Started

1.  **Initialize Features**: Create the `api.ts` and `types.ts` files for the `curriculum` module.
2.  **LMS Foundation**: Migrate the `ProgramsPage` and `CoursesPage` from `ui-examples/phd-journey-tracker_v11`.
3.  **Route Registration**: Add the new pages to `src/routes/index.tsx`.
4.  **Verify Data Flow**: Ensure the real backend APIs are returning the expected data structures.

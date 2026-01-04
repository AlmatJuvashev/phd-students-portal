# Assessment Engine & Quiz Builder - Полная документация

## Обзор

Документ описывает существующую реализацию Assessment Engine (бэкенд) и Quiz Builder (фронтенд v11 и studio/).

---

## 1. Backend Assessment Engine

### 1.1 API Endpoints

#### Assessments Module (`/api/assessments`)

| Method | Endpoint                        | Handler            | Описание                   | Статус             |
| ------ | ------------------------------- | ------------------ | -------------------------- | ------------------ |
| `POST` | `/api/assessments`              | `CreateAssessment` | Создание нового assessment | ⚠️ Not Implemented |
| `POST` | `/api/assessments/:id/attempts` | `StartAttempt`     | Начать попытку сдачи теста | ✅ Implemented     |

#### Attempts Module (`/api/attempts`)

| Method | Endpoint                     | Handler              | Описание                            | Статус         |
| ------ | ---------------------------- | -------------------- | ----------------------------------- | -------------- |
| `POST` | `/api/attempts/:id/response` | `SubmitResponse`     | Отправить ответ на вопрос           | ✅ Implemented |
| `POST` | `/api/attempts/:id/complete` | `CompleteAttempt`    | Завершить попытку и получить оценку | ✅ Implemented |
| `POST` | `/api/attempts/:id/log`      | `LogProctoringEvent` | Логирование событий прокторинга     | ✅ Implemented |

#### Item Bank Module (`/api/item-banks`)

| Method   | Endpoint                                      | Handler      | Описание                |
| -------- | --------------------------------------------- | ------------ | ----------------------- |
| `GET`    | `/api/item-banks/banks`                       | `ListBanks`  | Список банков вопросов  |
| `POST`   | `/api/item-banks/banks`                       | `CreateBank` | Создать банк вопросов   |
| `PUT`    | `/api/item-banks/banks/:bankId`               | `UpdateBank` | Обновить банк           |
| `DELETE` | `/api/item-banks/banks/:bankId`               | `DeleteBank` | Удалить банк            |
| `GET`    | `/api/item-banks/banks/:bankId/items`         | `ListItems`  | Список вопросов в банке |
| `POST`   | `/api/item-banks/banks/:bankId/items`         | `CreateItem` | Создать вопрос          |
| `PUT`    | `/api/item-banks/banks/:bankId/items/:itemId` | `UpdateItem` | Обновить вопрос         |
| `DELETE` | `/api/item-banks/banks/:bankId/items/:itemId` | `DeleteItem` | Удалить вопрос          |

---

### 1.2 Структура данных (Models)

#### QuestionBank

```go
type QuestionBank struct {
    ID             string          `json:"id"`
    TenantID       string          `json:"tenant_id"`
    Title          string          `json:"title"`
    Description    *string         `json:"description,omitempty"`
    Subject        *string         `json:"subject,omitempty"`       // Anatomy, Histology, etc.
    BloomsTaxonomy *BloomsTaxonomy `json:"blooms_taxonomy,omitempty"` // KNOWLEDGE, COMPREHENSION, etc.
    IsPublic       bool            `json:"is_public"`
    CreatedBy      string          `json:"created_by"`
    CreatedAt      time.Time       `json:"created_at"`
    UpdatedAt      time.Time       `json:"updated_at"`
}
```

#### Question

```go
type Question struct {
    ID                string           `json:"id"`
    BankID            string           `json:"bank_id"`
    Type              QuestionType     `json:"type"`              // MCQ, MRQ, TRUE_FALSE, TEXT, LIKERT
    Stem              string           `json:"stem"`              // Question text
    MediaURL          *string          `json:"media_url,omitempty"` // Image/Video
    PointsDefault     float64          `json:"points_default"`
    DifficultyLevel   *DifficultyLevel `json:"difficulty_level,omitempty"` // EASY, MEDIUM, HARD
    LearningOutcomeID *string          `json:"learning_outcome_id,omitempty"`
    Options           []QuestionOption `json:"options,omitempty"` // For MCQ/MRQ
}
```

#### QuestionOption

```go
type QuestionOption struct {
    ID         string  `json:"id"`
    QuestionID string  `json:"question_id"`
    Text       string  `json:"text"`
    IsCorrect  bool    `json:"is_correct"`
    SortOrder  int     `json:"sort_order"`
    Feedback   *string `json:"feedback,omitempty"` // Why correct/incorrect
}
```

#### Assessment

```go
type Assessment struct {
    ID               string           `json:"id"`
    TenantID         string           `json:"tenant_id"`
    CourseOfferingID string           `json:"course_offering_id"`
    Title            string           `json:"title"`
    Description      *string          `json:"description,omitempty"`
    TimeLimitMinutes *int             `json:"time_limit_minutes,omitempty"` // NULL = no limit
    AvailableFrom    *time.Time       `json:"available_from,omitempty"`
    AvailableUntil   *time.Time       `json:"available_until,omitempty"`
    ShuffleQuestions bool             `json:"shuffle_questions"`
    GradingPolicy    GradingPolicy    `json:"grading_policy"`    // AUTOMATIC, MANUAL_REVIEW
    SecuritySettings types.JSONText   `json:"security_settings"` // Proctoring config
    PassingScore     float64          `json:"passing_score"`
    CreatedBy        string           `json:"created_by"`
    Sections         []AssessmentSection `json:"sections,omitempty"`
}
```

#### AssessmentAttempt

```go
type AssessmentAttempt struct {
    ID           string        `json:"id"`
    AssessmentID string        `json:"assessment_id"`
    StudentID    string        `json:"student_id"`
    StartedAt    time.Time     `json:"started_at"`
    FinishedAt   *time.Time    `json:"finished_at,omitempty"`
    Score        float64       `json:"score"`
    Status       AttemptStatus `json:"status"` // IN_PROGRESS, SUBMITTED, GRADED
}
```

#### ItemResponse

```go
type ItemResponse struct {
    ID               string     `json:"id"`
    AttemptID        string     `json:"attempt_id"`
    QuestionID       string     `json:"question_id"`
    SelectedOptionID *string    `json:"selected_option_id,omitempty"` // For MCQ
    TextResponse     *string    `json:"text_response,omitempty"`      // For TEXT
    Score            float64    `json:"score"`
    IsCorrect        bool       `json:"is_correct"`
    GradedAt         *time.Time `json:"graded_at,omitempty"`
}
```

---

### 1.3 Enums

#### QuestionType

| Value        | Описание                                          |
| ------------ | ------------------------------------------------- |
| `MCQ`        | Multiple Choice Question (один правильный)        |
| `MRQ`        | Multiple Response Question (несколько правильных) |
| `TRUE_FALSE` | Правда/Ложь                                       |
| `TEXT`       | Свободный текстовый ответ                         |
| `LIKERT`     | Шкала Лайкерта                                    |

#### DifficultyLevel

| Value    |
| -------- |
| `EASY`   |
| `MEDIUM` |
| `HARD`   |

#### BloomsTaxonomy

| Value           | Описание   |
| --------------- | ---------- |
| `KNOWLEDGE`     | Знание     |
| `COMPREHENSION` | Понимание  |
| `APPLICATION`   | Применение |
| `ANALYSIS`      | Анализ     |
| `SYNTHESIS`     | Синтез     |
| `EVALUATION`    | Оценка     |

#### AttemptStatus

| Value         | Описание                    |
| ------------- | --------------------------- |
| `IN_PROGRESS` | Студент проходит тест       |
| `SUBMITTED`   | Тест сдан, ожидает проверки |
| `GRADED`      | Тест оценён                 |

#### GradingPolicy

| Value           | Описание                       |
| --------------- | ------------------------------ |
| `AUTOMATIC`     | Автоматическая проверка        |
| `MANUAL_REVIEW` | Ручная проверка преподавателем |

---

### 1.4 Workflow Assessment

```
┌─────────────────────────────────────────────────────────────────┐
│                    ASSESSMENT WORKFLOW                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────┐    ┌──────────┐    ┌────────────┐    ┌──────────┐ │
│  │ Create  │───▶│ Publish  │───▶│   Start    │───▶│ Submit   │ │
│  │ Quiz    │    │ (Admin)  │    │  Attempt   │    │ Response │ │
│  └─────────┘    └──────────┘    └────────────┘    └──────────┘ │
│       │                               │                  │      │
│       │                               │                  │      │
│       ▼                               ▼                  ▼      │
│  Question Bank              Check Availability    Save Answer   │
│  + Options                  Check Retake Policy   Auto-Save     │
│  + Sections                 Create Attempt Row    Track Time    │
│                                                                 │
│                                                                 │
│  ┌──────────┐    ┌──────────┐    ┌────────────┐                │
│  │ Complete │───▶│ Grade    │───▶│  Results   │                │
│  │ Attempt  │    │ (Auto)   │    │ + Feedback │                │
│  └──────────┘    └──────────┘    └────────────┘                │
│       │                │                 │                      │
│       │                │                 │                      │
│       ▼                ▼                 ▼                      │
│  Mark Finished   Calculate Score   Show Correct/               │
│  Trigger Grade   Per Question      Incorrect                   │
│                  Total Attempt     Per Question                │
│                                    Feedback                    │
└─────────────────────────────────────────────────────────────────┘
```

---

### 1.5 Proctoring Support

#### SecuritySettings JSON

```json
{
  "full_screen_mode": true,
  "track_tab_switches": true,
  "max_violations": 3,
  "auto_submit_on_limit": true,
  "record_webcam": false
}
```

#### ProctoringEventType

| Value             | Описание                       |
| ----------------- | ------------------------------ |
| `TAB_SWITCH`      | Переключение вкладки           |
| `WINDOW_BLUR`     | Потеря фокуса окна             |
| `FULLSCREEN_EXIT` | Выход из полноэкранного режима |
| `MOUSE_LEAVE`     | Курсор покинул окно            |
| `DEVICE_CHANGE`   | Изменение устройства           |

---

### 1.6 Database Schema (Migration 0097-0098)

```sql
-- Question Banks
CREATE TABLE question_banks (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    subject VARCHAR(100),
    blooms_taxonomy VARCHAR(50),
    is_public BOOLEAN DEFAULT FALSE,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

-- Questions
CREATE TABLE questions (
    id UUID PRIMARY KEY,
    bank_id UUID NOT NULL REFERENCES question_banks(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    stem TEXT NOT NULL,
    media_url VARCHAR(255),
    points_default FLOAT DEFAULT 1.0,
    difficulty_level VARCHAR(50),
    learning_outcome_id UUID
);

-- Question Options
CREATE TABLE question_options (
    id UUID PRIMARY KEY,
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    is_correct BOOLEAN DEFAULT FALSE,
    sort_order INT NOT NULL,
    feedback TEXT
);

-- Assessments
CREATE TABLE assessments (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    course_offering_id UUID NOT NULL REFERENCES course_offerings(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    time_limit_minutes INT,
    available_from TIMESTAMPTZ,
    available_until TIMESTAMPTZ,
    shuffle_questions BOOLEAN DEFAULT FALSE,
    grading_policy VARCHAR(50) DEFAULT 'AUTOMATIC',
    security_settings JSONB DEFAULT '{}',
    passing_score FLOAT DEFAULT 0.0,
    created_by UUID NOT NULL REFERENCES users(id)
);

-- Assessment Sections
CREATE TABLE assessment_sections (
    id UUID PRIMARY KEY,
    assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    title VARCHAR(255),
    instructions TEXT,
    sort_order INT NOT NULL
);

-- Assessment Items (Question Links)
CREATE TABLE assessment_items (
    id UUID PRIMARY KEY,
    assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    section_id UUID REFERENCES assessment_sections(id),
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    points_override FLOAT,
    sort_order INT NOT NULL,
    UNIQUE(assessment_id, question_id)
);

-- Assessment Attempts
CREATE TABLE assessment_attempts (
    id UUID PRIMARY KEY,
    assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    started_at TIMESTAMPTZ DEFAULT NOW(),
    finished_at TIMESTAMPTZ,
    score FLOAT DEFAULT 0.0,
    status VARCHAR(50) DEFAULT 'IN_PROGRESS'
);

-- Item Responses
CREATE TABLE item_responses (
    id UUID PRIMARY KEY,
    attempt_id UUID NOT NULL REFERENCES assessment_attempts(id) ON DELETE CASCADE,
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    selected_option_id UUID REFERENCES question_options(id),
    text_response TEXT,
    score FLOAT DEFAULT 0.0,
    is_correct BOOLEAN DEFAULT FALSE,
    graded_at TIMESTAMPTZ,
    UNIQUE(attempt_id, question_id)
);

-- Proctoring Logs
CREATE TABLE proctoring_logs (
    id UUID PRIMARY KEY,
    attempt_id UUID NOT NULL REFERENCES assessment_attempts(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL,
    occurred_at TIMESTAMPTZ DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'
);
```

---

## 2. v11 Quiz Builder

### 2.1 Расположение файлов

```
ui-examples/phd-journey-tracker_v11/
├── admin/pages/
│   ├── QuizBuilder.tsx      # Full Quiz authoring UI
│   ├── QuizPreview.tsx      # Student-view preview
│   ├── SurveyBuilder.tsx    # Survey authoring
│   └── SurveyPreview.tsx    # Survey preview
└── admin/features/itemBank/
    ├── types.ts             # Question types
    └── api.ts               # API calls
```

### 2.2 Типы вопросов (Quiz - v11)

| Type              | Icon          | Описание                 | Features                           |
| ----------------- | ------------- | ------------------------ | ---------------------------------- |
| `multiple_choice` | CheckSquare   | Single choice (MCQ)      | Options, correct marking, feedback |
| `multi_select`    | CheckSquare   | Multiple correct answers | Partial scoring option             |
| `short_text`      | MessageCircle | Text input answer        | Model answer comparison            |
| `ordering`        | Layers        | Drag-to-order items      | Correct sequence validation        |
| `matrix`          | TableIcon     | Matrix/grid questions    | Rows + Columns                     |
| `section_header`  | Bookmark      | Section divider          | Instructions text                  |
| `page_break`      | -             | Page separator           | Multi-page navigation              |

### 2.3 Типы вопросов (Survey - v11)

| Type              | Icon          | Описание                 |
| ----------------- | ------------- | ------------------------ |
| `rating_stars`    | Star          | 1-5 звёзд                |
| `scale_10`        | Hash          | Шкала 0-10 (NPS)         |
| `likert_matrix`   | TableIcon     | Матрица Лайкерта         |
| `open_feedback`   | MessageSquare | Открытый текстовый ответ |
| `multiple_choice` | CheckSquare   | Выбор из вариантов       |
| `section_header`  | Bookmark      | Заголовок секции         |

### 2.4 QuizBuilder UI Features (v11)

1. **Drag & Drop Reorder** - Перетаскивание вопросов через Framer Motion `Reorder.Group`
2. **Markdown + Math Support** - `$$LaTeX$$` формулы в тексте вопроса
3. **Zen Mode** - Полноэкранный режим фокусировки
4. **AI Distractor Generation** - Placeholder для AI-генерации неправильных ответов
5. **Rich Text Toolbar** - Bold, Italic, Math (Sigma) buttons
6. **Question Import from Bank** - `QuestionPickerDialog` для выбора из Item Bank
7. **Real-time Collaboration** - `AvatarGroup` показывает активных редакторов
8. **Adaptive Feedback** - Separate feedback for correct/incorrect answers
9. **Points Configuration** - Per-question point value
10. **Display Logic** - Conditional question display (planned)

### 2.5 QuizPreview Features (v11)

1. **Timer Display** - Countdown with warning at < 60 seconds
2. **Progress Bar** - Current question / total
3. **Math Rendering** - LaTeX formula display
4. **Option Selection** - Radio button style for MCQ
5. **Results Screen** - Score percentage, time used, pass/fail status
6. **Performance Breakdown** - Per-question feedback with correct answers
7. **Try Again** - Reset and restart quiz

---

## 3. Current Frontend Studio Components

### 3.1 Расположение файлов

```
frontend/src/features/studio/
├── components/
│   ├── QuizBuilderModal.tsx    # Modal quiz editor
│   ├── SurveyBuilderModal.tsx  # Modal survey editor
│   ├── ChecklistBuilderModal.tsx
│   ├── ConfirmTaskBuilderModal.tsx
│   ├── FormBuilderModal.tsx
│   └── MarkdownEditor.tsx
└── types.ts                    # Shared type definitions
```

### 3.2 QuizBuilderModal Capabilities

**Текущие возможности:**

- ✅ Question list with drag & drop reorder
- ✅ Question types: `multiple_choice`, `multi_select`, `short_text`, `ordering`, `section_header`
- ✅ Option editing with correct answer marking
- ✅ Points per question
- ✅ Feedback (correct/incorrect)
- ✅ Quiz config: time limit, passing score, shuffle questions

**Отсутствует по сравнению с v11:**

- ❌ Markdown + Math rendering (`$$LaTeX$$`)
- ❌ Zen Mode
- ❌ AI Distractor generation
- ❌ Rich text toolbar (Bold, Italic, Math)
- ❌ Question import from Item Bank
- ❌ Real-time collaboration indicators
- ❌ Matrix question type
- ❌ Page breaks
- ❌ Display logic (conditional questions)
- ❌ Live preview button

### 3.3 SurveyBuilderModal Capabilities

**Текущие возможности:**

- ✅ Question types: `rating_stars`, `scale_10`, `likert_matrix`, `open_feedback`, `multiple_choice`, `section_header`
- ✅ Drag & drop reorder
- ✅ Required toggle per question
- ✅ Matrix rows/columns editing
- ✅ Config: anonymous mode

**Отсутствует:**

- ❌ Progress bar setting
- ❌ Preview mode
- ❌ Zen Mode

### 3.4 QuizQuestion Type (studio/types.ts)

```typescript
export interface QuizQuestion {
  id: string;
  type: QuestionType;
  text: string;
  subtitle?: string;
  hint?: string;
  points: number;
  feedback_correct?: string;
  feedback_incorrect?: string;
  options?: { id: string; text: string; is_correct: boolean }[];
  correct_order?: string[];
  display_logic?: {
    depends_on_question_id: string;
    condition: "equals" | "not_equals" | "contains";
    value: string;
  };
}
```

---

## 4. GAP Analysis: Что нужно добавить

### 4.1 Backend Gaps

| Feature                           | Status             | Priority | Notes                                |
| --------------------------------- | ------------------ | -------- | ------------------------------------ |
| `CreateAssessment` endpoint       | ❌ Not implemented | HIGH     | Returns 501 Not Implemented          |
| `GetAssessment` endpoint (public) | ❌ Missing         | HIGH     | Need GET /api/assessments/:id        |
| `ListAssessments` endpoint        | ❌ Missing         | HIGH     | Need GET /api/assessments            |
| `UpdateAssessment` endpoint       | ❌ Missing         | MEDIUM   | Need PUT /api/assessments/:id        |
| `DeleteAssessment` endpoint       | ❌ Missing         | LOW      | Need DELETE /api/assessments/:id     |
| `ListAttempts` for student        | ❌ Missing         | MEDIUM   | GET /api/assessments/:id/my-attempts |
| `GetAttemptDetails`               | ❌ Missing         | MEDIUM   | GET /api/attempts/:id                |
| Manual grading endpoints          | ❌ Missing         | LOW      | For TEXT questions                   |
| Retake policy logic               | ❌ Missing         | MEDIUM   | Max attempts, cooldown               |
| Time limit enforcement            | ❌ Missing         | HIGH     | Auto-submit on timeout               |

### 4.2 Frontend Gaps (studio/)

| Feature                     | v11 Has | studio/ Has | Priority |
| --------------------------- | ------- | ----------- | -------- |
| Matrix question type        | ✅      | ❌          | MEDIUM   |
| Math/LaTeX rendering        | ✅      | ❌          | HIGH     |
| Question import from bank   | ✅      | ❌          | HIGH     |
| Zen Mode                    | ✅      | ❌          | LOW      |
| AI Distractor generation    | ✅      | ❌          | LOW      |
| Live preview                | ✅      | ❌          | MEDIUM   |
| Display logic               | ✅      | ❌          | MEDIUM   |
| Page breaks                 | ✅      | ❌          | LOW      |
| Real-time collab indicators | ✅      | ❌          | LOW      |

### 4.3 Student Experience Gaps

| Feature                | Status     | Priority |
| ---------------------- | ---------- | -------- |
| Quiz taking UI         | ❌ Missing | HIGH     |
| Timer component        | ❌ Missing | HIGH     |
| Auto-save responses    | ❌ Missing | MEDIUM   |
| Results page           | ❌ Missing | HIGH     |
| Feedback display       | ❌ Missing | HIGH     |
| Retake functionality   | ❌ Missing | MEDIUM   |
| Proctoring client-side | ❌ Missing | LOW      |

---

## 5. Implementation Recommendations

### 5.1 Phase 1: Complete Backend APIs (Priority: HIGH)

1. Implement `CreateAssessment` in service layer
2. Add CRUD endpoints for assessments
3. Add `GetAttemptDetails` with responses
4. Implement time limit enforcement

### 5.2 Phase 2: Quiz Taking UI (Priority: HIGH)

1. Create `QuizTaker.tsx` component
2. Implement timer with auto-submit
3. Build results display page
4. Connect to `/api/attempts/*` endpoints

### 5.3 Phase 3: Enhanced Builder (Priority: MEDIUM)

1. Migrate v11 QuizBuilder to studio/
2. Add math/LaTeX rendering
3. Implement Question Bank picker dialog
4. Add matrix question type support

### 5.4 Phase 4: Advanced Features (Priority: LOW)

1. Display logic editor
2. Proctoring client integration
3. Real-time collaboration
4. AI-assisted question generation

---

## 6. API Request/Response Examples

### Create Assessment

```http
POST /api/assessments
Content-Type: application/json

{
  "course_offering_id": "uuid",
  "title": "Midterm Exam",
  "description": "Covers chapters 1-5",
  "time_limit_minutes": 60,
  "available_from": "2025-01-15T09:00:00Z",
  "available_until": "2025-01-15T11:00:00Z",
  "shuffle_questions": true,
  "grading_policy": "AUTOMATIC",
  "passing_score": 70,
  "security_settings": {
    "full_screen_mode": true,
    "track_tab_switches": true,
    "max_violations": 3,
    "auto_submit_on_limit": true
  }
}
```

### Start Attempt

```http
POST /api/assessments/:id/attempts

Response:
{
  "id": "attempt-uuid",
  "assessment_id": "assessment-uuid",
  "student_id": "student-uuid",
  "started_at": "2025-01-15T09:05:23Z",
  "status": "IN_PROGRESS",
  "score": 0
}
```

### Submit Response

```http
POST /api/attempts/:id/response
Content-Type: application/json

{
  "question_id": "question-uuid",
  "option_id": "selected-option-uuid"  // For MCQ
}

// OR for text questions:
{
  "question_id": "question-uuid",
  "text_response": "Student's written answer"
}
```

### Complete Attempt

```http
POST /api/attempts/:id/complete

Response:
{
  "id": "attempt-uuid",
  "assessment_id": "assessment-uuid",
  "student_id": "student-uuid",
  "started_at": "2025-01-15T09:05:23Z",
  "finished_at": "2025-01-15T09:45:12Z",
  "score": 85.5,
  "status": "SUBMITTED"
}
```

---

_Документ создан: 4 января 2026_
_Версия: 1.0_

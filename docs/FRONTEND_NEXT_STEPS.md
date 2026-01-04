# Frontend Implementation: Next Steps

> **Документ:** Следующие шаги после аудита  
> **Дата:** 4 января 2026  
> **Статус:** Основано на аудите реализации vs FRONTEND_MIGRATION_GUIDE.md

---

## 📊 Текущий Статус Реализации

| Phase       | Модуль            | Прогресс | Статус                        |
| ----------- | ----------------- | -------- | ----------------------------- |
| **Phase 1** | `curriculum/`     | 60%      | ⚠️ Нет ProgramDetailPage      |
| **Phase 2** | `enrollments/`    | 95%      | ✅ Готово                     |
| **Phase 2** | `course-content/` | 0%       | ℹ️ Реализовано в `studio/`    |
| **Phase 3** | `item-bank/`      | 70%      | ⚠️ Нет update/delete API      |
| **Phase 3** | `assessment/`     | 0%       | ℹ️ Частично в `studio/`       |
| **Phase 4** | `grading/`        | 0%       | ⚠️ Частично в `teacher/`      |
| **Phase 5** | `student-portal/` | 20%      | 🔴 **КРИТИЧНО** — mock данные |
| **Phase 5** | `teacher/`        | 90%      | ✅ Готово                     |
| **Phase 6** | `studio/`         | 85%      | ✅ Готово                     |

**Общий прогресс: ~45%**

---

## 🔴 Критический Пробел: Student Portal

### Проблема

[StudentDashboard.tsx](../frontend/src/features/student-portal/StudentDashboard.tsx) использует **hardcoded mock-данные**:

```typescript
// Текущий код — НЕ работает с API:
const activeProgram = {
  title: user?.program || t("student.dashboard.default_program"),
  progress: 0, // ← hardcoded
  overdue: 0, // ← hardcoded
};
```

### Backend API Статус

| Эндпоинт                             | Статус      | Комментарий                    |
| ------------------------------------ | ----------- | ------------------------------ |
| `GET /student/dashboard`             | ❌ Нет      | Нужно создать                  |
| `GET /student/courses`               | ❌ Нет      | Нужно создать                  |
| `GET /student/assignments`           | ❌ Нет      | Нужно создать                  |
| `GET /student/grades`                | ⚠️ Частично | Есть `/grading/student/:id`    |
| `GET /student/enrollments`           | ⚠️ Частично | Метод в repo есть, handler нет |
| `GET /journey/progress`              | ✅ Есть     | Можно переиспользовать         |
| `GET /grading/transcript/:studentId` | ✅ Есть     | Для транскрипта                |

---

## 📋 План Реализации

### Week 1: Student Portal Backend + Frontend

#### День 1-2: Backend — Student API Handlers

**Создать** `backend/internal/handlers/student_handler.go`:

```go
// Эндпоинты для создания:
GET /api/student/dashboard    // Агрегация: progress + enrollments + upcoming
GET /api/student/courses      // Курсы текущего студента
GET /api/student/assignments  // Активные задания с дедлайнами
GET /api/student/grades       // Оценки текущего студента (self)
```

**Переиспользовать существующее:**

- `LMSRepository.GetStudentEnrollments()` — уже есть
- `GradingService.GetStudentGrades()` — адаптировать для self
- `JourneyService.GetProgress()` — для прогресса программы

#### День 3-4: Frontend — Student Portal API

**Создать** `frontend/src/features/student-portal/api.ts`:

```typescript
// API функции:
export const getStudentDashboard = () =>
  api.get<StudentDashboard>("/student/dashboard");
export const getStudentCourses = () =>
  api.get<StudentCourse[]>("/student/courses");
export const getStudentAssignments = () =>
  api.get<Assignment[]>("/student/assignments");
export const getStudentGrades = () => api.get<GradeEntry[]>("/student/grades");
```

**Создать** `frontend/src/features/student-portal/types.ts`:

```typescript
export interface StudentDashboard {
  program: ProgramProgress;
  upcomingDeadlines: Deadline[];
  recentGrades: GradeEntry[];
  announcements: Announcement[];
}

export interface StudentCourse {
  id: string;
  title: string;
  code: string;
  instructor: string;
  progress: number;
  nextActivity?: Activity;
}
```

#### День 5: Frontend — Подключение к API

**Обновить** `StudentDashboard.tsx`:

- Заменить mock-данные на `useQuery('studentDashboard', getStudentDashboard)`
- Добавить loading/error states
- Подключить реальные данные прогресса

**Создать страницы:**

- `StudentCourses.tsx` — список курсов студента
- `StudentAssignments.tsx` — задания с дедлайнами
- `StudentGrades.tsx` — оценки и транскрипт

#### День 6-7: Роуты и тестирование

**Добавить роуты** в `routes/index.tsx`:

```typescript
{ path: 'my-courses', element: <StudentCourses /> },
{ path: 'my-assignments', element: <StudentAssignments /> },
{ path: 'my-grades', element: <StudentGrades /> },
```

---

### Week 2: Grading Module + Item Bank Completion

#### День 1-3: Grading Module

**Создать** `frontend/src/features/grading/`:

```
grading/
├── api.ts              # getGradebook, getPendingSubmissions, submitGrade
├── types.ts            # GradebookEntry, Submission, GradeRequest
├── GradebookPage.tsx   # Для админов/преподавателей
├── SubmissionQueue.tsx # Очередь на проверку
└── components/
    ├── GradeInput.tsx
    └── RubricGrader.tsx
```

**Backend:** Эндпоинты уже существуют в `/api/grading/*`

#### День 4-5: Item Bank Completion

**Обновить** `frontend/src/features/item-bank/api.ts`:

```typescript
// Добавить недостающие функции:
export const updateBank = (id: string, data: Partial<QuestionBank>) =>
  api.put<QuestionBank>(`/item-bank/banks/${id}`, data);

export const deleteBank = (id: string) => api.delete(`/item-bank/banks/${id}`);

export const updateQuestion = (
  bankId: string,
  id: string,
  data: Partial<Question>
) => api.put<Question>(`/item-bank/banks/${bankId}/questions/${id}`, data);

export const deleteQuestion = (bankId: string, id: string) =>
  api.delete(`/item-bank/banks/${bankId}/questions/${id}`);

export const importQuestions = (bankId: string, file: File) =>
  api.upload(`/item-bank/banks/${bankId}/import`, file);
```

#### День 6-7: Curriculum Completion

**Создать** `frontend/src/features/curriculum/ProgramDetailPage.tsx`:

- Детали программы
- Список курсов в программе
- Статистика enrollments
- Кнопка "Edit in Builder"

---

### Week 3: Polish & Integration Testing

#### День 1-2: Student Layout

**Создать** `frontend/src/layouts/StudentLayout.tsx`:

- Навигация: Dashboard, My Courses, Assignments, Grades, Journey
- Профиль студента в sidebar
- Уведомления о дедлайнах

#### День 3-4: E2E Testing

**Сценарии для тестирования:**

1. Student login → Dashboard с реальными данными
2. Student → My Courses → Course Detail → Activity
3. Student → Assignments → Submit → Check Grade
4. Teacher → Grading Queue → Grade Submission
5. Admin → Enrollments → Enroll Student → Verify in Student Portal

#### День 5-7: Bug Fixes & Documentation

- Исправление найденных багов
- Обновление FRONTEND_MIGRATION_GUIDE.md с отметками ✅
- API документация для новых эндпоинтов

---

## 📁 Файлы для Создания

### Backend (5 файлов)

| Файл                                   | Описание                       |
| -------------------------------------- | ------------------------------ |
| `internal/handlers/student_handler.go` | Student API handlers           |
| `internal/services/student_service.go` | Student business logic         |
| `internal/dto/student_dto.go`          | Student DTOs                   |
| Обновить `cmd/server/routes.go`        | Регистрация /student/\* роутов |
| Обновить `docs/swagger.yaml`           | API документация               |

### Frontend (12 файлов)

| Файл                                             | Описание              |
| ------------------------------------------------ | --------------------- |
| `features/student-portal/api.ts`                 | Student API клиент    |
| `features/student-portal/types.ts`               | Student типы          |
| `features/student-portal/StudentCourses.tsx`     | Страница курсов       |
| `features/student-portal/StudentAssignments.tsx` | Страница заданий      |
| `features/student-portal/StudentGrades.tsx`      | Страница оценок       |
| `features/grading/api.ts`                        | Grading API клиент    |
| `features/grading/types.ts`                      | Grading типы          |
| `features/grading/GradebookPage.tsx`             | Gradebook для админов |
| `features/grading/SubmissionQueue.tsx`           | Очередь проверки      |
| `features/curriculum/ProgramDetailPage.tsx`      | Детали программы      |
| `layouts/StudentLayout.tsx`                      | Layout для студентов  |
| Обновить `routes/index.tsx`                      | Новые роуты           |

---

## ⚠️ Решения по Архитектуре

### 1. Course Content — оставить в `studio/`

**Решение:** НЕ создавать отдельный `course-content/` модуль.

**Причина:** Функционал Course Builder уже полностью реализован в `studio/`:

- `studio/CourseBuilder.tsx`
- `studio/components/ActivityList.tsx`
- `studio/components/ActivityDetails.tsx`
- Quiz/Survey/Form builders как модальные окна

**Действие:** Обновить FRONTEND_MIGRATION_GUIDE — отметить Phase 2 course-content как "Implemented in studio/"

### 2. Assessment — частичная реализация достаточна

**Решение:** НЕ создавать отдельный `assessment/` модуль.

**Причина:** Quiz/Survey builders реализованы в `studio/components/`:

- `QuizBuilderModal.tsx`
- `SurveyBuilderModal.tsx`
- `FormBuilderModal.tsx`

**Действие:** Если понадобятся standalone страницы — импортировать компоненты из studio/

### 3. Grading — выделить в отдельный модуль

**Решение:** Создать `grading/` модуль, вынести логику из `teacher/`.

**Причина:** Grading нужен для:

- Админов (GradebookPage)
- Преподавателей (уже есть TeacherGradingPage)
- Студентов (StudentGrades — просмотр своих оценок)

**Действие:**

1. Создать `grading/api.ts` с общими функциями
2. `teacher/` импортирует из `grading/`
3. `student-portal/` импортирует из `grading/`

---

## 📊 Обновлённая Оценка Прогресса

После выполнения этого плана:

| Phase                    | Текущий | После Week 1 | После Week 2 | После Week 3  |
| ------------------------ | ------- | ------------ | ------------ | ------------- |
| Phase 1 (curriculum)     | 60%     | 60%          | 80%          | 90%           |
| Phase 2 (enrollments)    | 95%     | 95%          | 95%          | 100%          |
| Phase 2 (course-content) | N/A     | N/A          | N/A          | ✅ in studio/ |
| Phase 3 (item-bank)      | 70%     | 70%          | 95%          | 100%          |
| Phase 3 (assessment)     | N/A     | N/A          | N/A          | ✅ in studio/ |
| Phase 4 (grading)        | 0%      | 0%           | 80%          | 95%           |
| Phase 5 (student-portal) | 20%     | **80%**      | 90%          | 100%          |
| Phase 5 (teacher)        | 90%     | 90%          | 95%          | 100%          |
| **Общий**                | **45%** | **60%**      | **80%**      | **95%**       |

---

## 🎯 Критерии Готовности (Definition of Done)

### Student Portal Ready ✓

- [ ] Student Dashboard показывает реальные данные из API
- [ ] Student может видеть свои курсы с прогрессом
- [ ] Student может видеть активные задания с дедлайнами
- [ ] Student может видеть свои оценки и транскрипт
- [ ] Роуты зарегистрированы и работают

### Grading Module Ready ✓

- [ ] GradebookPage показывает оценки по курсу
- [ ] SubmissionQueue показывает работы на проверку
- [ ] Teacher может выставить оценку
- [ ] Student видит оценку после выставления

### Item Bank Complete ✓

- [ ] CRUD операции для Banks работают
- [ ] CRUD операции для Questions работают
- [ ] Import вопросов работает

---

## 📝 Следующий Шаг

**Начать с:** Backend Student API (`student_handler.go`)

Это разблокирует всю работу по Student Portal на фронтенде.

```bash
# Создать файл:
touch backend/internal/handlers/student_handler.go
```

После создания backend handlers — переходить к frontend api.ts и типам.

# Universal Education Portal Blueprint (Market-Competitive)
## Core Platform, University Package, School Package, Prep-Center Package, and a Minimal Roadmap

> Goal: a single **multi-tenant** education platform that can serve universities, schools, and pre-university exam prep centers — without forking the codebase.

---

## 0) Product Positioning (How We Win)

Most “LMS” products either:
- focus on **content delivery** (courses + assignments), or
- focus on **process/workflow** (forms + approvals), or
- focus on **exam prep** (question banks + drills).

A competitive universal platform should combine all three with clear differentiation:

1) **Journey-first learning + workflow**  
   A first-class “journey / checklist / milestone” engine for admissions, compliance, research, and multi-step learning programs (this is your unique advantage).

2) **Standards-based interoperability**  
   So organizations can integrate tools without duct-tape:  
   - LTI 1.3 / LTI Advantage for tool integrations and grade passback. (https://www.imsglobal.org/spec/lti/v1p3, https://www.imsglobal.org/lti-advantage-overview)  
   - OneRoster for rostering and gradebook exchange with SIS. (https://www.imsglobal.org/oneroster-11-introduction)  
   - QTI for assessment interchange (item banks). (https://www.imsglobal.org/spec/qti/v3p0/oview)  
   - SCORM import for legacy e-learning packages. (https://scorm.com/scorm-explained/)  
   - xAPI / Caliper for learning activity telemetry. (https://standards.ieee.org/ieee/9274.1.1/7321/, https://www.imsglobal.org/spec/caliper/v1p2)

3) **Modern “LXP-like” experience**  
   Personalized feed, smart recommendations, searchable knowledge base, multi-language support, mobile-first UI, and fast onboarding.

4) **AI that is grounded and auditable**  
   RAG grounded in org documents; citations; admin controls; safe-by-default prompts; “why this answer” transparency.

---

## 1) Core Platform (Common for All Organization Types)

### 1.1 Multi-tenancy & Organization Model
**Entity: `Organization`**
- `id`, `name`
- `type` = `university | school | prep_center | other`
- branding: `logo_url`, colors, locale defaults
- tenancy: `domain/subdomain`, `data_residency_region` (optional), feature flags in `settings` (JSON)

**Principles**
- Nearly every entity has `organization_id`.
- Admin-only controls for feature flags and packages.
- Strong tenant isolation: queries are always scoped by org.

---

### 1.2 Identity, Access & Governance (Enterprise-grade “table stakes”)
**Authentication**
- Local accounts + OAuth/OIDC SSO (https://openid.net/specs/openid-connect-core-1_0.html)  
- Optional SAML SSO for enterprise/universities (https://docs.oasis-open.org/security/saml/Post2.0/sstc-saml-tech-overview-2.0.html)

**Authorization**
- RBAC + “scopes”:
  - org-wide roles (student/teacher/admin/parent/…)
  - course roles (instructor/TA/student)
  - workflow roles (approver/reviewer/committee_member)

**Audit & compliance**
- Audit log for sensitive actions (grade changes, approvals, access to private docs).
- Data retention policies and export tooling.
- Accessibility baseline: WCAG 2.2 targets for the UI (https://www.w3.org/TR/WCAG22/)

---

### 1.3 Learning Model (Content + Activities + Journeys)
**Core hierarchy**
- `Course` → `Module` → `Lesson` → `Activity`

**Activity types (extensible)**
- `content` (rich text, video links, embeds)
- `quiz` (auto-graded)
- `assignment` (text/files, rubric optional)
- `checklist` (multi-step tasks)
- `milestone` (sign-off gates)
- `survey` (feedback / evaluation)
- `mock_exam` (timed, high-stakes simulation; mostly for prep centers)

**Journey engine (core differentiator)**
- Any program can be represented as a **journey graph**:
  - prerequisites / unlock rules
  - due dates / SLAs
  - required documents
  - approvals (multi-step)
  - “guardrails” (archived nodes, role checks)

---

### 1.4 Assessment Engine (Competitive by default)
A “best-in-market” education platform needs a strong assessment layer, not just a quiz form.

**Question bank**
- `ItemBank` / `Question`
- tagging: topic, difficulty, outcome/competency, exam section
- versioning + review workflow (“draft → reviewed → published”)

**Tests**
- fixed forms and randomized forms
- item selection rules
- time limits, sections, scoring policies

**Interchange**
- QTI import/export for items/tests when feasible (https://www.imsglobal.org/spec/qti/v3p0/oview)

**Grading**
- auto-grading for objective items
- rubrics for subjective work
- moderation workflow (second marker / committee review) as optional package behavior

---

### 1.5 Submissions, Gradebook, Certificates
- `Submission` for any activity requiring student response.
- `Grade` with:
  - scale types (numeric, letter, 1–5, pass/fail, custom)
  - rubric breakdown (optional)
- Gradebook views:
  - per course, per group/class, per student
- Certificates / micro-credentials:
  - configurable templates
  - completion rules (e.g., “finish modules A–D + pass final test”)

---

### 1.6 Collaboration & Communication
To compete with modern platforms, messaging cannot be an afterthought.

- Course channels (announcements + discussions)
- Direct messages
- Group channels (cohorts, classes, prep groups)
- Attachments (integrates with `FileObject`)
- Notification engine:
  - in-app
  - email (optional)
  - push (later; mobile app)

---

### 1.7 Files, Content Library & Versioning
- S3/MinIO-backed `FileObject`
- Content library:
  - reusable modules/lessons across courses
  - templates for assignments/tests
- Versioning:
  - content version history
  - “publish” vs “draft”
  - safe updates without breaking in-progress cohorts

**Standards (optional but valuable)**
- SCORM package import for legacy content (https://scorm.com/scorm-explained/)

---

### 1.8 Scheduling (Core + UI profiles)
Scheduling is a competitive feature for schools and prep centers, and also useful for universities.

- `CalendarEvent` / `ScheduleEntry`
- course sessions, deadlines, exam sessions
- optional attendance tracking (enabled by school/prep profile)

---

### 1.9 Analytics & Learning Telemetry (Core-first, then advanced)
**Core analytics**
- completion rates, drop-off points, time-to-complete
- course and module health metrics

**Telemetry standards (optional, for ecosystems)**
- xAPI for experience tracking (IEEE standard reference: https://standards.ieee.org/ieee/9274.1.1/7321/)
- Caliper for event models and collection (https://www.imsglobal.org/spec/caliper/v1p2)

---

### 1.10 Integrations & APIs (Make it “plug-and-play”)
- Webhooks for key events (submission created, grade posted, milestone approved)
- Public API (scoped tokens) for external dashboards
- LTI 1.3 / Advantage support (tool launch, deep linking, grade passback)  
  (https://www.imsglobal.org/spec/lti/v1p3, https://www.imsglobal.org/lti-advantage-overview)
- OneRoster import/export (CSV + REST where needed)  
  (https://www.imsglobal.org/oneroster-11-introduction)

---

### 1.11 AI Layer (Core)
**AI Assistant**
- org-scoped knowledge base (docs, policies, course materials)
- citations/grounding and safe defaults
- admin prompts and guardrails

**AI as a platform service**
- multiple model providers (cloud + local)
- policy: what data can be sent where (by org settings)
- logs for AI requests (optional; privacy-sensitive)

---

## 2) University Package (Profile: `organization.type = university`)

### 2.1 University Roles (Extensions on RBAC)
- `advisor/scientific_supervisor`
- `program_director/dean`
- `committee_member` (scientific council, dissertation council, ethics, IRB, etc.)
- optional: `department_admin`, `research_office`

### 2.2 Programs, Cohorts, and Complex Journeys
**Entities**
- `Program` (PhD/Master/Residency/etc.)
- `Cohort` (intake year / cycle)
- links:
  - enrollments scoped to cohort + program
  - deadlines derived from cohort rules

**Journeys**
- admissions workflow (doc checks, interviews, decisions)
- dissertation lifecycle:
  - topic approval → proposals → ethics → data collection → publications → pre-defense → defense
- multi-approver gates (advisor + committee)

### 2.3 Research & Compliance Tooling (Differentiator for Universities)
- research project registry (optional)
- ethics approvals and document capture
- publication tracking (manual + optional ORCID integration later)
- anti-plagiarism integration (future module) via LTI or API

### 2.4 University Assessment Extensions
- committee-based grading
- rubric libraries for theses / defenses
- moderation workflow:
  - first marker, second marker, committee finalization
- anomaly detection on scorers (advanced analytics)

### 2.5 University Analytics
- cohort funnel: how many reach each milestone
- advisor dashboards: student statuses and risks
- management dashboards: completion time, compliance issues, research output indicators

---

## 3) School Package (Profile: `organization.type = school`)

### 3.1 School Roles
- `parent`
- `class_teacher` (homeroom)
- `school_director` / `deputy_director`

### 3.2 Classes, Subjects, and Rosters
**Entity**
- `Class` (e.g., “7A”, “10B”)
- class roster management
- optional subject-teacher mapping

**Interoperability**
- OneRoster import/export for rosters and gradebook (https://www.imsglobal.org/oneroster-11-introduction)

### 3.3 Timetable, Attendance, Homework
- timetable view by class
- attendance (optional v1; common requirement)
- homework as `Activity.type = assignment`
- parent notifications for missing work

### 3.4 Gradebook, Student Diary, Report Cards
- configurable grading scales (1–5, 1–10, letter, custom)
- teacher gradebook + student diary views
- report card exports (PDF/Excel)

### 3.5 Parent Portal (Must-have to compete in K-12)
- multi-child support
- progress + grades + homework
- alerts: missing work, low performance, absence patterns

### 3.6 School Analytics
- class-level performance dashboards
- teacher workload
- risk identification (students needing support)

---

## 4) Prep-Center Package (Profile: `organization.type = prep_center`)

Prep centers compete on results, not on “course pages”. The platform must feel like an exam training product.

### 4.1 Groups Instead of “Classes”
Reuse `Class` entity but label it as **Group** in UI:
- “ENT-2026 Math Group #3”
- “IELTS Weekend Group #1”

### 4.2 Exam Profiles & Scoring
**Entities**
- `ExamProfile`
  - exam name (ENT/SAT/IELTS/etc.)
  - sections and weights
  - score conversion rules
- `MockExamAttempt` (could be a specialized submission wrapper)

### 4.3 Mock Exams and Diagnostics (The Core Prep Differentiator)
- timed mock exams (`Activity.type = mock_exam`)
- diagnostic tests (“baseline”, “mid”, “final”)
- question bank with tagging:
  - topics, difficulty, exam section
- analytics:
  - trend lines (score trajectory)
  - topic weaknesses
  - predicted score vs target score

### 4.4 Adaptive Practice (Best-in-market feature)
- personalized daily sets:
  - spaced repetition
  - weakness-based item selection
  - mastery thresholds
- streaks/gamification (optional, configurable)

### 4.5 Parent Role (Optional but Valuable for Paid Prep)
- parents see:
  - attendance + homework completion
  - mock score trend
  - “exam readiness” indicator

### 4.6 Business Module (Optional, Competitive in Private Prep)
This can be a separate paid add-on later:
- packages/subscriptions
- invoices/receipts
- coupons/referrals
- CRM-lite: leads → enrolled students

---

## 5) Capability Matrix (Core vs Packages)

| Layer / Feature | Core Platform | University Package | School Package | Prep-Center Package |
|---|---|---|---|---|
| Multi-tenancy & branding | ✅ | ✅ | ✅ | ✅ |
| RBAC + audit logs | ✅ | ✅ | ✅ | ✅ |
| Courses/modules/lessons/activities | ✅ | ✅ | ✅ | ✅ |
| Journey/checklist engine | ✅ | ✅ (heavy) | ✅ (light) | ✅ (study plans) |
| Assessment engine + item bank | ✅ | ✅ (rubrics, moderation) | ✅ | ✅ (advanced) |
| Gradebook | ✅ | ✅ | ✅ (report cards) | ✅ (exam scoring) |
| Messaging + notifications | ✅ | ✅ | ✅ | ✅ |
| Scheduling | ✅ | ✅ | ✅ (timetable/attendance) | ✅ (sessions, mocks) |
| Documents & templates | ✅ | ✅ (regulations) | ✅ (forms, reports) | ✅ (materials) |
| Standards: LTI / OneRoster / QTI / SCORM | ✅ (platform) | ✅ | ✅ | ✅ |
| AI assistant + RAG | ✅ | ✅ (regulations) | ✅ (school policies) | ✅ (study tips/explanations) |
| Analytics | ✅ | ✅ (program/cohort) | ✅ (class-level) | ✅ (readiness, trends) |

---

## 6) Minimal Roadmap (Competitive and Realistic)

> The goal is to build a *credible* v1 that can win pilots, while keeping the architecture “enterprise-ready”.

### Phase 0 — Foundation (Non-negotiable basics)
- Multi-tenancy (`Organization`) + tenant isolation
- RBAC + audit log
- Files (`FileObject`) + permissions
- Localization (at least EN + RU/KZ-ready framework)
- Accessibility baseline target (WCAG 2.2 principles) (https://www.w3.org/TR/WCAG22/)

**Exit criteria:** multiple organizations can be onboarded safely.

---

### Phase 1 — Core Learning + Assessment v1
- Course/module/lesson/activity
- Submissions + gradebook
- Question bank + quizzes
- Basic analytics (completion, drop-off, averages)
- Simple certificates (completion rules)

**Exit criteria:** platform can run a real course end-to-end.

---

### Phase 2 — Communication + Scheduling v1
- Course channels + DMs + group channels
- Notifications (in-app + email optional)
- Calendar / schedule entries
- Basic “task / deadline” views

**Exit criteria:** daily operations work without WhatsApp/Telegram glue.

---

### Phase 3 — Interoperability (Ecosystem readiness)
- OneRoster import/export (CSV first; REST later as needed) (https://www.imsglobal.org/oneroster-11-introduction)
- LTI 1.3 / Advantage (launch + deep link + grade passback)  
  (https://www.imsglobal.org/spec/lti/v1p3, https://www.imsglobal.org/lti-advantage-overview)
- QTI import/export (start with import) (https://www.imsglobal.org/spec/qti/v3p0/oview)
- SCORM import (optional but high-value for adoption) (https://scorm.com/scorm-explained/)

**Exit criteria:** institutions can integrate existing tools/content quickly.

---

### Phase 4 — AI v1 (Grounded, safe, useful)
- Org-scoped RAG assistant with citations
- Admin prompt controls + “allowed sources” policy
- AI study support for:
  - explaining answers (prep),
  - policy/regulation Q&A (university),
  - “what to do next” planning (journeys)

**Exit criteria:** AI helps more than it confuses; admins can control it.

---

### Phase 5 — Package Releases (University / School / Prep)
**University package**
- Program + cohort + approvals workflows
- research/dissertation journey presets
- advisor dashboards

**School package**
- classes, timetable, parent portal
- report cards exports

**Prep-center package**
- exam profiles + mock exams
- score trends + readiness analytics
- adaptive practice (optional v1.1)

**Exit criteria:** each profile can be sold/piloted with a clean onboarding story.

---

## 7) Competitive Differentiators to Emphasize in Pitch Decks
- Journey-first workflows (admissions → milestones → approvals) = fewer emails, fewer “lost steps”.
- Standards support (LTI/OneRoster/QTI/SCORM) = faster institutional adoption.
- Exam-readiness analytics for prep centers = measurable value.
- Org-scoped AI grounded in internal docs = safer, more trustworthy answers.
- White-label + multi-language + data residency knobs = regional competitiveness.

And yes, if we do standards + great UX, we get to say: “Integrates cleanly” without crossing fingers behind our back.

# 🧠 Project Learnings & Agent Protocol

This document serves as a **compounding learning loop** for all AI agents working on the KazNMU PhD Student Portal.
**Goal:** Distill insights after each task to improve future performance, reduce context switching time, and avoid repeating mistakes.

---

## 🏗️ Project Context

*   **Domain:** PhD Student Management System for a Medical University.
*   **Stack:**
    *   **Backend:** Go (Gin framework), PostgreSQL, JWT Auth.
    *   **Frontend:** React (Vite), TypeScript, TailwindCSS, Radix UI.
    *   **Infrastructure:** Docker Compose, S3 (MinIO/AWS), Mailpit.
*   **Key Concepts:**
    *   **Journey Map:** A graph-based progression system for students (Nodes, Edges).
    *   **RBAC:** Strict role separation (Student, Advisor, Admin).
    *   **Tenancy:** Multi-tenant architecture (via `X-Tenant-Slug` header).

---

## 💡 Critical Learnings & Patterns

### 1. File Uploads (S3 Presigned URLs)
*   **Pattern:** The backend does *not* handle file content directly.
*   **Flow:**
    1.  Frontend requests a **presigned PUT URL** from the backend.
    2.  Frontend uploads the file directly to S3/MinIO using this URL.
    3.  Frontend confirms the upload to the backend (sending `object_key` and `etag`).
*   **Gotcha:** Do not try to send `multipart/form-data` to the backend for file storage. Always use the presign flow.

### 2. Database & Migrations
*   **Tool:** `golang-migrate` is used.
*   **Location:** `backend/db/migrations`.
*   **Protocol:** Always create a new migration file for schema changes. Do not modify existing applied migrations.
*   **Command:** Use `make migrate-up` in `backend/` to apply.

### 3. Frontend State Management
*   **Data Fetching:** Uses `react-query` (or similar hooks pattern).
*   **Auth:** JWTs are stored in `localStorage` (or cookies, verify implementation).
*   **UI Components:** Uses Radix UI primitives. When adding new UI, prefer composing these over raw HTML/CSS.

### 4. Testing
*   **E2E:** Playwright is used for end-to-end testing (`frontend/tests`).
*   **Backend:** Standard Go testing (`go test ./...`).

---

## 🤖 Agent Protocol

**Before starting a task:**
1.  **Read this file** to understand the architectural constraints and patterns.
2.  **Check `task.md`** (if exists) for current progress.

**After completing a task:**
1.  **Update this file** if you encountered a new pattern, a tricky bug, or a significant architectural decision.
2.  **Format:** Add a new entry to "Task History" or "Critical Learnings" if applicable.

---

## 📜 Task History & Insights

| Date | Task | Agent | Key Insight / Learning |
| :--- | :--- | :--- | :--- |
| 2025-12-07 | Audit & Documentation | Antigravity | Established `LEARNINGS.md`. Clarified S3 presign pattern and RBAC structure. Updated `README.md` to be user-centric. |
| 2025-12-07 | Root Makefile | Antigravity | Created root `Makefile` to orchestrate multi-service stack (Backend, Frontend, Infra). Added `kill-ports` utility for common dev friction. |

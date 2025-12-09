# KazNMU PhD Student Portal

Welcome to the official PhD Student Portal for Asfendiyarov Kazakh National Medical University. This application is designed to streamline the doctoral journey, facilitating collaboration between students, advisors, and administrators.

## 🌟 Application Overview

The PhD Student Portal is a comprehensive platform that automates the workflow of doctoral programs. It replaces manual paperwork with a digital, transparent, and efficient process.

**Key Goals:**
*   **Transparency:** Clear visualization of the doctoral journey and requirements.
*   **Efficiency:** Automated document submission, review, and approval flows.
*   **Collaboration:** Integrated communication tools for students and advisors.

---

## 📘 User Guide

### 🎓 For PhD Students
Your dashboard is your command center for your doctoral program.
*   **Journey Map:** View your entire program roadmap. Each "node" represents a task or milestone (e.g., "Submit Research Protocol").
*   **Submitting Documents:** Click on a node to upload required files. You can track the status of your submissions (Submitted, Under Review, Approved, Needs Fixes).
*   **Chat:** Communicate directly with your assigned advisors or administrators through the built-in chat feature. You can create group chats or send direct messages.
*   **Notifications:** Stay updated with real-time alerts for document reviews, deadlines, and new messages.

### 👨‍🏫 For Advisors
Manage your mentees effectively and ensure they stay on track.
*   **Student Oversight:** View a list of your assigned students and their current progress.
*   **Document Review:** Receive notifications when a student submits a document. You can approve it, reject it with feedback, or request specific changes.
*   **Communication:** Use the chat to provide quick guidance or schedule meetings via the calendar integration.

### 🛡️ For Administrators
Oversee the entire doctoral program and manage system settings.
*   **User Management:** Create and manage accounts for students and advisors. Assign advisors to students.
*   **Program Management:** Define the "Journey Map" (Playbook), including milestones, prerequisites, and deadlines.
*   **Reporting:** Access statistics on student progress and program performance.
*   **System Settings:** Configure global settings and manage tenant-specific configurations.

---

## 🚀 Key Features

*   **Interactive Journey Map:** A visual representation of the PhD curriculum, tracking progress step-by-step.
*   **Document Workflow:** Secure upload (S3-backed) and versioned document history.
*   **Real-time Chat:** Integrated messaging with file support.
*   **Calendar & Events:** Scheduling for defenses, exams, and consultations.
*   **Role-Based Access Control (RBAC):** Secure access ensures users only see what they are authorized to see.
*   **Mobile Friendly:** Fully responsive design for access on any device.

---

## 🛠️ Developer & Deployment Info

### Technology Stack
*   **Frontend:** React 18, TypeScript, Vite, TailwindCSS, Radix UI
*   **Backend:** Go 1.21+, Gin, PostgreSQL 14+, JWT Auth
*   **Infrastructure:** Docker Compose, S3 (MinIO/AWS), Mailpit

### Quickstart (Local Development)

1.  **Mailpit (Email Service):**
    ```bash
    cd mailserver && docker compose up -d
    ```

2.  **Database:**
    *   Ensure PostgreSQL is running.
    *   Configure `DATABASE_URL` in `backend/.env`.

3.  **Backend:**
    ```bash
    cd backend
    make migrate-up
    make run
    ```

4.  **Frontend:**
    ```bash
    cd ../frontend
    npm install
    VITE_API_URL=http://localhost:8080/api npm run dev
    ```

### Mock Data
Generate test data (advisors, students, progress) for development:
```bash
./mocks/generate_mock_data.sh
```

### Documentation Links
*   [**System Flow Diagrams**](docs/flow_diagrams.md) - Detailed architectural flows.
*   [**Deployment Guide**](deploy/DEPLOYMENT_GUIDE.md) - Production deployment instructions.
*   [**API Documentation**](backend/README.md) - Backend API reference.

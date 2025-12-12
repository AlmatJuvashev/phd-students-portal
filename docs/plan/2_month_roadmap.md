# 2-Month MVP Roadmap: "The Visual Journey LMS"

## The Challenge
Building a full Moodle replacement in 8 weeks is impossible.
Building a **unique, high-value Visual LMS MVP** is possible.

**Our Goal**: A working application where a Teacher can visually build a "Course Map" (Videos, PDFs, Quizzes) and a Student can traverse it.

---

## Month 1: The Foundation (Builder & Content)

### Week 1-2: The Journey Constructor (Visual Builder)
*   **Goal**: Kill the JSON files. Staff needs a UI.
*   **Dev Focus**:
    *   Integrate `reactflow` library.
    *   Backend API for `POST /playbooks` and `PUT /playbooks/:id`.
    *   UI to Drag-and-Drop Nodes (Start -> Lesson 1 -> Quiz 1 -> End).
*   *Cut*: No complex "undo/redo" history yet. No concurrent editing.

### Week 3-4: LMS Nodes (Content Delivery)
*   **Goal**: Nodes can hold real lesson content.
*   **Dev Focus**:
    *   **Video Node**: Embed YouTube/Vimeo.
    *   **Rich Text Node**: Markdown editor for reading assignments.
    *   **File Node**: PDF Viewer for reading papers.
*   *Cut*: No internal video hosting (use YouTube/Vimeo/Google Drive). No interactive SCORM packages.

---

## Month 2: Assessment & Grading

### Week 5-6: The Quiz Engine (Basic)
*   **Goal**: Auto-graded assessments to "Unlock" the next node.
*   **Dev Focus**:
    *   Schema for Questions (Multiple Choice, True/False).
    *   Database tables: `quizzes`, `quiz_questions`, `quiz_submissions`.
    *   Logic: If `score > 80%` then `Node.state = done`.
*   *Cut*: No essay questions (hard to grade manually). No free-text matching. No "Question Banks".

### Week 7: Gradebook & Tracking
*   **Goal**: Teacher sees who is stuck.
*   **Dev Focus**:
    *   Simple Dashboard: Table of Students vs Nodes Completed.
    *   Calculated "Course Progress %".
*   *Cut*: No weighted averages or complex GPA curves. Pass/Fail or Simple % only.

### Week 8: Polish & Pilot
*   **Goal**: Ready for first real users.
*   **Dev Focus**:
    *   Bug bashing.
    *   Mobile responsiveness check.
    *   "Inviting Students" flow (Email invites).

---

## The "Kill List" (Features we MUST delay to survive)
1.  **LTI Integrations**: No Zoom, no Turnitin, no connection to University main DB. Standalone only.
2.  **Social Features**: No Forums, no peer-to-peer chat.
3.  **Complex Permissions**: Single "Teacher" role, single "Student" role. No "Teaching Assistants" or "Department Heads".
4.  **Attendance**: Not part of the MVP.

## Conclusion
With focused effort (and my help 🤖), we can hit this MVP. It won't be Moodle, but it will be **something Moodle isn't**: A beautiful, graphical learning path builder.

# 6-Month Roadmap: "The Advanced LMS"

## Overview
This roadmap builds upon the 2-Month MVP foundation.
**Goal**: Transform the "Visual Journey MVP" into a production-grade LMS capable of supporting a full academic semester with multiple roles and complex grading.

## Month 1-2: Core MVP (See `2_month_roadmap.md`)
*   **Result**: Visual Builder v1, Basic Content Nodes, Simple Quiz (MCQ), Teacher Dashboard.

---

## Month 3: The "Complex Assessment" Sprint
Focus on **Testing & Content Reuse**.

### Features
1.  **Question Banks**:
    *   Database: `question_bank` (Tag-based: "Math", "Hard", "Calculus").
    *   UI: "Add Random Question from Bank 'Calculus'".
2.  **Essay Grading Interface**:
    *   Student: Rich Text Editor / File Upload.
    *   Teacher: Split-screen view (Essay on left, Rubric/Grade on right).
3.  **Advanced Quiz Types**:
    *   Matching, Ordering, Fill-in-the-Blank.

---

## Month 4: The "Academic Grading" Sprint
Focus on **Data Integrity & Calculation**.

### Features
1.  **Weighted Grading Engine**:
    *   Structure: `Category` (Homework 20%, Midterm 30%, Final 50%).
    *   Logic: Auto-calculate specific node scores into these buckets.
2.  **GPA & Curves**:
    *   **Curves**: "Add +5 points to all students who took Quiz X".
    *   **GPA**: Mapping Percentage -> Letter Grade (A, B, C) -> GPA (4.0, 3.0).
3.  **Gradebook Grid**:
    *   Excel-like view for Teachers to see all grades at once and manually override if needed.

---

## Month 5: The "Hierarchy & Access" Sprint
Focus on **University Structure**.

### Features
1.  **Attendance System**:
    *   **Event Node**: "Lecture - Oct 15".
    *   **Check-in**:
        *   *Manual*: Teacher marks P/A/L.
        *   *Code-based*: Student enters a temporary 4-digit code displayed on projector.
2.  **Advanced RBAC (Role-Based Access Control)**:
    *   **TA (Teaching Assistant)**: Can grade assignments but NOT change course structure.
    *   **Department Head**: Can view all courses in Dept, but not edit.
    *   **Dean**: Read-only access to stats.

---

## Month 6: Polish, Performance & Launch Prep
Focus on **Stability**.

### Features
1.  **Bulk Operations**:
    *   "Enroll 500 students via CSV".
    *   "Clone Course for next Semester".
2.  **Reporting & Analytics**:
    *   "At-Risk Student" reports (Attendance < 80% AND Grade < C).
    *   Department-level aggregate stats.
3.  **Load Testing**: Ensure the system handles 500 concurrent quiz submissions.

---

## Critical Dependencies
1.  **Question Bank Schema**: Must be designed well in Month 3 to support reuse.
2.  **Canvas/Excel Import/Export**: Universities live in CSVs. We need robust import tools in Month 6.

## Conclusion
This 6-month plan is **solid**. It moves from "Toy App" (Month 2) to "Serious Tool" (Month 6). It addresses the "Day-to-Day" needs of a university (Taking Attendance, Grading Essays, curves).

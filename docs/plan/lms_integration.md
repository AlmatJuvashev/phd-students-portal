# LMS Integration Feature Analysis

## Goal Description
Expand the Journey Map to support **Learning & Assignment** nodes. This transforms the portal from purely "Administrative Compliance" to also covering the "Educational" aspect of the PhD program (e.g., specific courses, seminars, mandatory readings).

## Feasibility Analysis
**"Does it fit?"**: **Yes, High Fit.** PhD programs are hybrid.
*   *Current State*: Handles administrative gates (e.g., "Submit Application", "Get Approval").
*   *Future State*: Should handle academic requirements (e.g., "Complete Research Ethics Course", "Read seminal papers").

**Complexity**: **Low to Medium**.
*   We already support:
    *   **Submissions**: `type: "upload"` nodes exist. Homework submission is just a re-labeled upload node.
    *   **Prerequisites**: "You must finish Tutorial A before doing Assignment B" is already solved by our graph engine.
*   **New Requirements**:
    *   **Content Display**: We need a way to show *Video/Text/PDF* content *inside* the node modal before asking for the upload.

## Implementation Plan

### 1. New Node Types / Logic
We can either create strict new types or extend the generic `requirements` object.

**Option A: Dedicated Types**
*   `type: "tutorial"` (Video/Text content, "Mark as Done" button).
*   `type: "assignment"` (Instructions + File Upload).

**Option B: Enhanced "Info" & "Form" Nodes (Recommended)**
Extend the existing `NodeDef` to support **Content Blocks**.

```typescript
type NodeDef = {
  // ...
  content?: {
    description_md?: string; // Markdown text (Lecture notes)
    video_embed_url?: string; // YouTube/Vimeo link
    attachments?: Array<{ title: string; url: string }>; // Reading materials (PDFs)
  };
  // If 'requirements.uploads' exists -> It's an assignment.
  // If no uploads -> It's a "Read/Watch & Confirm" node.
}
```

### 2. UI Updates (Node Modal)
Refactor the `NodeDetails.tsx` (or the new Modal we created) to render this content **above** the action buttons.
*   **Video**: Render an iframe if `video_embed_url` is present.
*   **Reading List**: Render a list of downloadable links with icon (e.g., `[PDF] Research Ethics v2.pdf`).
*   **Markdown**: Render rich text instructions.

### 3. Grading / Feedback (The Hard Part)
*   **Current status**: We have a simple "Administrator" role.
*   **Requirement**: "Teachers" need to grade specifically *their* course's assignments.
*   **Solution**:
    *   Add `owner_role` or `instructor_id` to the Node.
    *   Only that specific Instructor (or generic Admin) can switch the node state from `submitted` -> `done` (Graded).

## Pros & Cons

### Pros
*   **Unified Experience**: Students don't need Moodle/Canvas *and* this Portal. They see *everything* required for their PhD in one map.
*   **Automatic Gating**: "You cannot schedule your Thesis Defense (Admin Node)" until you "Complete Ethics Course (LMS Node)". A separate LMS cannot easily enforce this.

### Cons
*   **Reinventing the Wheel**: LMS platforms (Moodle, Canvas) are huge for a reason. Building a *good* grading interface, quizzes, gradebook, and plagiarism checks is massive work.
*   **Recommendation**: Keep it "Light". Use it for *mandatory distinct milestones* (e.g., "Ethics Exam", "Topic Seminar"). Do NOT try to replace Moodle for day-to-day classwork if the university already has one. Use it for **Program Requirements** only.

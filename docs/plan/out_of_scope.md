# Reality Check: Features Out of Scope (6-Month Timeline)

Even with an aggressive 6-month roadmap, we are building a **Specialized Web LMS**, not an **Enterprise Platform**. Here is what we explicitly **CANNOT** build in this timeframe.

## 1. Integrations (The "Black Hole" of Time)
*   **LTI (Learning Tools Interoperability)**:
    *   *What it is*: The standard that lets Canvas plug into Zoom, Turnitin, Pearson, McGraw-Hill.
    *   *Why NO*: It involves implementing complex XML protocols (LTI 1.3). It takes months just to pass certification.
    *   *Impact*: You cannot use "Turnitin" inside our assignments. You cannot "Launch Zoom" directly from the calendar (must paste links manually).

*   **SIS Real-Time Sync**:
    *   *What it is*: Auto-syncing with Banner/PeopleSoft/PowerSchool every 5 minutes.
    *   *Why NO*: Accessing these legacy enterprise DBs is a nightmare.
    *   *Workaround*: We will use **CSV Imports** (Bulk upload) instead.

## 2. Advanced Content Formats
*   **SCORM / xAPI Player**:
    *   *What it is*: Playing interactive zip files exported from "Articulate Storyline" or "Adobe Captivate".
    *   *Why NO*: Requires a dedicated "SCORM Player Engine" to track every click inside a 3rd party package.
    *   *Impact*: We can only host Videos, PDFs, and our own Quizzes. Not external interactive packages.

*   **Native Video Hosting**:
    *   *What it is*: Uploading raw `.mp4` files directly to our server.
    *   *Why NO*: Transcoding (making 1080p, 720p, 480p versions) is expensive and hard.
    *   *Workaround*: Use YouTube (Unlisted), Vimeo, or Google Drive links.

## 3. Specialized Engines
*   **Plagiarism Detection**:
    *   *Why NO*: Requires a database of billions of websites to compare against. Impossible to build from scratch.

*   **Offline Mobile App**:
    *   *What it is*: Students downloading quizzes to take on the bus without verification.
    *   *Why NO*: Requires native iOS/Android development and complex "Conflict Resolution" sync logic. We are building a **Mobile-Responsive Web App** (works great in browser, but needs internet).

## 4. Complex Social Features
*   **Discussion Forums (Advanced)**:
    *   *Why NO*: We can do simple comments. But "Threaded views", "Email replies", "Moderation queues", and "Upvoting" is a separate product (like Reddit).

## Summary
In 6 months, you get a **Sleek, Modern, Standalone LMS**.
You do **NOT** get an "Integration Hub" that talks to every legacy software in the university.

# PhD Student Portal - Application Audit

This document outlines the current state of the application and identifies missing features categorized by priority.

## 1. Required Features (Critical Gaps)

These features are essential for a complete and functional doctoral journey management system.

| Feature | Description | Justification |
| :--- | :--- | :--- |
| **Scheduling & Calendar** | A dedicated calendar view for students and advisors. | Students need to book meetings with advisors, track defense deadlines, and see university academic events. Currently, no scheduling system exists. |
| **In-App Notifications** | A notification center (bell icon) for real-time updates. | Email notifications are implemented, but in-app notifications are crucial for immediate feedback on document status, chat messages, and system alerts. |
| **Advanced Analytics** | Aggregate dashboards for Admins and Chairs. | The current `students-monitor` lists students. Admins need high-level metrics: "How many students are delayed?", "Average time to complete Stage 1", "Advisor load balancing". |
| **Granular RBAC** | Refined permissions beyond basic role checks. | Ensure strict data isolation. For example, Advisors should only see *their* students' detailed sensitive data. Chairs should see their department's data. |
| **Mobile Optimization** | Comprehensive mobile responsiveness audit. | The journey map and complex tables must be usable on mobile devices, as students often check status on the go. |

## 2. Enhance Application (Usability & Efficiency)

These features would significantly improve the user experience and operational efficiency.

| Feature | Description | Benefit |
| :--- | :--- | :--- |
| **Global Search** | Search bar to find students, documents, and chat messages. | As data grows, finding specific information becomes difficult without a unified search. |
| **Granular Document Feedback** | Ability to comment on specific parts of a submission. | Currently, rejection is likely binary. Advisors need to provide specific feedback (e.g., "Fix the bibliography") directly on the document record. |
| **University Integrations** | SSO (LDAP/SAML) and Library System integration. | Reduces login friction and automates verification steps (e.g., automatically checking if a thesis was deposited in the library). |
| **Offline/PWA Support** | Basic offline access to the journey map. | Allows students to check requirements even with poor internet connectivity. |
| **Activity Audit Log** | Visible history of actions for each student. | A timeline view for students to see "Who changed my status and when?" for transparency. |

## 3. Nice to Have (Future Proofing)

These features add "delight" and advanced capabilities but are not immediately critical.

| Feature | Description | Value Add |
| :--- | :--- | :--- |
| **AI Assistant** | Chatbot trained on the `playbook.json` and regulations. | Instantly answers FAQs like "What documents do I need for defense?" reducing admin workload. |
| **Multi-Language UI** | Full localization (KZ/RU/EN) for the entire interface. | While the playbook supports languages, the app shell (menus, buttons) should also be fully localized. |
| **Dark Mode** | System-wide dark mode support. | Improves accessibility and user preference. |
| **Peer Support Groups** | Interest-based forums or groups. | Fosters a community among PhD students beyond their immediate cohort. |
| **Portfolio Generation** | Auto-generate a CV/Portfolio from uploaded data. | Helps students prepare for post-doc applications by aggregating their publications and achievements. |

## Summary of Next Steps

1.  **Prioritize "Required" bucket:** Begin design and implementation of the Scheduling and Notification systems.
2.  **Refine "Enhance" items:** Assess technical feasibility of SSO and Global Search.
3.  **Backlog "Nice to Have":** Keep for future roadmap discussions.

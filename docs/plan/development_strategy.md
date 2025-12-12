# Development Strategy: Vertical vs. Iterative

## The Dilemma
*   **Approach A (Vertical / Fully Featured)**: Build the "Perfect Quiz Engine" (Month 1-2). Then the "Perfect Gradebook" (Month 3-4).
*   **Approach B (Horizontal / Iterative)**: Build a "Simple Quiz" + "Simple Gradebook" (Month 1-2). Then upgrade both (Month 3-4).

## The Recommendation: Approach B (Iterative)

**Verdict**: You **MUST** choose **Approach B (Iterative)** for an LMS.

### Why? The "Ecosystem" Problem
An LMS is not a collection of isolated tools; it is an *interconnected ecosystem*.
*   A **Quiz** is useless without a **Gradebook** to store the score.
*   A **Gradebook** is useless without a **Roster** (Students).
*   A **Roster** is useless without a **Course**.

If you use **Approach A** (Vertical):
*   By Month 2, you have an incredible Quiz engine (Supports essays, matching, banks).
*   *But you can't save the grade anywhere* because you haven't built the Gradebook yet.
*   *Result*: You have zero deliverable value until Month 6.

If you use **Approach B** (Iterative):
*   By Month 2, you have:
    *   A simple Quiz (True/False only).
    *   A simple Gradebook (Pass/Fail only).
*   *Result*: **You have a working product**. A Pilot group can actually take a course. Use feedback to decide if you *really* need "Essay Questions" or if "Multiple Choice" is enough.

## Practical Example: The "Quiz" Feature

| Phase | Vertical Approach (Don't do this) | Iterative Approach (Do this) |
| :--- | :--- | :--- |
| **Month 2** | Quiz with 10 question types, time limits, IP restrictions, Question Banks. | Quiz with **only** Multiple Choice. Direct save to DB. |
| **Month 4** | (Still fixing bugs in complex quiz logic). Start Gradebook. | Add "Essay" type. Add "Weighted Grades" to Gradebook. |
| **Month 6** | Connect Quiz to Gradebook. Realize data structures don't match. **Panic.** | Content reuse. Question Banks. Analytics. |

## Conclusion
**"Polishing" is a trap.** exact requirements often change once users touch the software.
Build the **Skeleton** first (The flow of data from Start -> Quiz -> Grade -> Report).
Then add **Muscle** (Features) to that skeleton.

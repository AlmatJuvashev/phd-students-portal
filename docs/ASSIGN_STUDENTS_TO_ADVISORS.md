# Agent Specification: Assign or Reassign Advisor to Student

## Purpose

This agent manages assignment and reassignment of advisors to students within the PhD portal.  
It allows authorized users (admins or program staff) to:

- Assign an advisor to a student for the first time.
- Self-assign as the advisor.
- Reassign an existing student to a new advisor.
- Notify relevant parties and maintain a complete audit trail.

---

## Inputs

| Field         | Type            | Required                  | Description                                              |
| ------------- | --------------- | ------------------------- | -------------------------------------------------------- |
| `student_id`  | integer         | ✅                        | Target student ID.                                       |
| `actor_id`    | integer         | ✅                        | The admin or staff user performing the action.           |
| `advisor_id`  | integer or null | ⚙️                        | ID of advisor to assign; optional if `self_assign=true`. |
| `self_assign` | boolean         | Optional (default: false) | When true, assigns the current actor as the advisor.     |
| `reason`      | string          | Optional                  | Reason for (re)assignment, used in audit record.         |

---

## Authorization Rules

1. Only users with roles **`admin`**, **`secretary`**, or **`program_manager`** may assign or reassign advisors.
2. A student can have only **one active primary advisor** at any given time.
3. The process must be **idempotent**—if the advisor is already assigned, the agent should return `no_op=true`.
4. Each operation (assign, reassign, no-op) must be logged in the audit trail and trigger notifications.

---

## Data Model / Endpoints

The agent expects or can trigger the following API resources:

| Method                                          | Endpoint                                         | Purpose |
| ----------------------------------------------- | ------------------------------------------------ | ------- |
| `GET /api/users/:id`                            | Retrieve user details and roles.                 |
| `GET /api/students/:student_id/advisors/active` | Fetch the currently assigned advisor (if any).   |
| `POST /api/advisor-assignments`                 | Create a new assignment record.                  |
| `PATCH /api/advisor-assignments/:id/close`      | Close a previous advisor assignment.             |
| `POST /api/notifications`                       | Create notifications for affected users.         |
| `POST /api/audit`                               | Record audit entries (optional but recommended). |

---

## Validation Sequence

1. **Authorize** the actor (`actor_id`).
2. **Validate** that the student exists and is active.
3. Resolve `advisor_id`:
   - If `self_assign=true`, use `actor_id`.
   - Otherwise require `advisor_id` to be provided.
4. Ensure the target advisor exists and has a valid role (`advisor`, `faculty`, `admin`).
5. Retrieve the student’s current advisor assignment.

---

## Decision Logic

| Scenario                      | Action                                                                                   | Result |
| ----------------------------- | ---------------------------------------------------------------------------------------- | ------ |
| No current advisor            | Create new assignment → Notify student & advisor → Audit `ASSIGN_ADVISOR`.               |
| Same advisor already assigned | Return `no_op=true` → Audit `ASSIGN_ADVISOR_NOOP`.                                       |
| Different advisor assigned    | Close old assignment → Create new → Notify all three parties → Audit `REASSIGN_ADVISOR`. |

---

## Notifications

The agent should generate these notification events:

| Type                              | Recipient        | Payload Example                                                             |
| --------------------------------- | ---------------- | --------------------------------------------------------------------------- |
| `ADVISOR_ASSIGNED_STUDENT`        | Student          | `{ "student_id":123, "advisor_id":456, "by":7 }`                            |
| `ADVISOR_ASSIGNED_ADVISOR`        | New advisor      | `{ "student_id":123, "advisor_id":456, "by":7 }`                            |
| `ADVISOR_REASSIGNED_PREV_ADVISOR` | Previous advisor | `{ "student_id":123, "prev_advisor_id":789, "new_advisor_id":456, "by":7 }` |

---

## Audit Log Format

| Field       | Example                                                           |
| ----------- | ----------------------------------------------------------------- |
| `action`    | `"ASSIGN_ADVISOR"`, `"REASSIGN_ADVISOR"`, `"ASSIGN_ADVISOR_NOOP"` |
| `entity`    | `"student"`                                                       |
| `entity_id` | `<student_id>`                                                    |
| `old`       | `{ "advisor_id": 789 }`                                           |
| `new`       | `{ "advisor_id": 456 }`                                           |
| `reason`    | `"Load balancing"`                                                |
| `actor_id`  | `<admin_id>`                                                      |

---

## Responses

### ✅ Success

```json
{
  "ok": true,
  "no_op": false,
  "student_id": 123,
  "advisor_id": 456,
  "previous_advisor_id": 789,
  "assignment_id": 1122,
  "message": "Advisor reassigned successfully."
}
```

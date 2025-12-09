# Journey Nodes Logic & Acceptance Criteria

This document outlines the logic for each node in the PhD Student Journey map, specifically focusing on nodes that require document submissions and acceptance.

## Roles
- **Student**: Initiates the process, uploads documents, and confirms tasks.
- **Advisor**: Reviews and approves specific documents (e.g., reviews).
- **Secretary (DC/SC)**: Manages the process, accepts documents, and moves the student to the next stage.
- **Admin**: System administrator with full access.

## Node Logic Table

The following table details the actions required for each node, including uploads and specific logic.

| Node ID | Title (RU) | Who Can Complete | Requirements (Uploads/Fields) | Logic / Next Steps |
| :--- | :--- | :--- | :--- | :--- |
| **S1_profile** | Профиль докторанта | Student | **Fields:** Full Name, Specialty, Program, Graduation Date, Advisors, Dissertation Form | **Next:** `S1_text_ready` <br> **Outcome:** Profile filled. |
| **S1_text_ready** | Текст диссертации подготовлен | Student | **Upload:** `dissertation_draft_file` (DOCX/PDF) | **Next:** `S0_antiplagiat` <br> **Logic:** Confirms text formatting per standard. |
| **S0_antiplagiat** | Справка на антиплагиат | Student | **Upload:** `antiplagiat_cert` (PDF/IMG) | **Next:** `S1_publications_list` <br> **Logic:** Confirms uniqueness certificate upload. |
| **S1_publications_list** | Список публикаций | Student | **Upload:** `signed_pub_list` (PDF/IMG) <br> **Fields:** Checkboxes for thesis, reviews, app7, sc_extract, primary_docs, lcb_defense | **Next:** `E1_apply_omid` <br> **Logic:** Confirms publication list and related docs. |
| **E1_apply_omid** | Заявка в ОМиД | Student | **Upload:** `omid_application` (PDF) | **Next:** `NK_package` <br> **Logic:** Confirms application submission to OMiD. |
| **E3_hearing_nk** | Заслушивание НК | Student | **Upload:** `sc_protocol_extract_nk` (PDF/IMG) <br> **Fields:** Hearing happened? Remarks exist? Plan prepared? Remarks resolved? | **Next:** `RP1_overview_actualization` OR `D1_normokontrol_ncste` <br> **Condition:** `rp_required` (if >3 years). |
| **RP1_overview_actualization** | Актуализация Research Proposal | Student | **Fields:** Notes on OMiD letter, SC date, Student prep | **Next:** `RP2_sc_hearing_prep` <br> **Condition:** Only if `rp_required`. |
| **RP2_sc_hearing_prep** | Заслушивание RP в НК | Student | **Upload:** `sc_protocol_extract` (PDF/IMG) <br> **Fields:** Checkboxes for date, presentation, publist, presented, recommended, reviewers | **Next:** `RP3_pre_expertise_application` |
| **RP3_pre_expertise_application** | Заявление в ОМиД (пред. эксп.) | Student | **Action:** Confirm submission | **Next:** `NK_package` |
| **D1_normokontrol_ncste** | НЦГНТЭ: нормоконтроль | Student | **Uploads:** `dissertation_docx` (DOCX), `ncste_receipt` (PDF/IMG) | **Next:** `IV_rector_application` |
| **IV_rector_application** | Заявление ректору | Student | **Upload:** `rector_letter` (PDF) | **Next:** `IV3_publication_certificate_ncste` |
| **IV3_publication_certificate_ncste** | НЦГНТЭ: справка о публикациях | Student | **Upload:** `ncste_publication_certificate` (PDF) | **Next:** `D2_apply_to_ds` |
| **D2_apply_to_ds** | Пакет документов в ДС | Student | **Fields:** Checkboxes for App 5, App 6, SC extract, Dissertation print, App 7, App 8, AC extract, Transcript | **Next:** `V1_reinstatement_package` |
| **V1_reinstatement_package** | Восстановление на защиту | Student | **Fields:** Checkboxes for Rector letter, DC letter, Transcript, ID copy, Topic approvals, Pre-reviews, Fee (if applicable) | **Next:** `A1_post_acceptance_overview` |
| **A1_post_acceptance_overview** | После принятия документов ДС | Student | **Fields:** Notes on Orders, Notification, Publication, Letters | **Next:** `VI1_post_defense_overview` |
| **VI1_post_defense_overview** | После защиты | Student | **Fields:** Notes on Orders, Site publication | **Next:** `VI2_library_deposits` |
| **VI2_library_deposits** | Сдача в библиотеки | Student | **Fields:** Checkboxes for Hardbound copies, CD abstracts, CD hard case, Salem letters, Delivery to NAL, NL, KazNMU Lib, Receipts | **Next:** `VI3_ncste_state_registration` |
| **VI3_ncste_state_registration** | НЦГНТЭ: госрегистрация | Student | **Fields:** Checkboxes for Print materials (Unbound thesis, Abstracts, UKD/DEK, ID, Pub list) and CD materials (Word thesis, Abstracts, UKD/DEK, PDF Pub list), Salem letter, Submission | **Next:** `VI_attestation_file` |
| **VI_attestation_file** | Аттестационное дело | Student | **Fields:** Checkboxes for Inventory, Rector app, DC app, Cover letter, Student info, Personal sheet, Diplomas, Transcript copy, Topic extract, Pub list copies, Abstracts, SC extract, Ethics conclusion, NCSTE antiplag, Uni antiplag, Advisor reviews, Reviewer reports, Vote protocol, Ballots, Attendance list, Defense protocol, Video, Reg card, NCSTE reg card, Library certs (KazNMU, NL, NAL), Final thesis CD | **Next:** `END` |

## Document Acceptance Logic

For nodes requiring uploads (e.g., `S1_text_ready`, `S0_antiplagiat`, `E1_apply_omid`), the following general logic applies:

1.  **Student Uploads**: The student uploads the required file(s) matching the specified MIME types (PDF, DOCX, Images).
2.  **System Validation**: The system checks if the file is present and matches the allowed types.
3.  **Confirmation**: The student must explicitly confirm the action (e.g., clicking "Confirm upload" or answering "Yes").
4.  **State Transition**: Upon confirmation, the node state changes to `completed`, and the user is allowed to proceed to the `next` node.

*Note: Currently, most nodes are self-confirmed by the Student. Admin/Secretary verification steps are implicit in the process (e.g., "Wait for confirmation" instructions) or handled outside the digital flow, unless specific "Approval" nodes are implemented in future versions.*

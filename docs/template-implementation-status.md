# Template Prefill Implementation Status

## Overview

This document tracks the implementation status of template prefills using `[% %]` delimiters across all nodes in the PhD Student Portal.

## Template Delimiter Migration

All templates have been migrated from `{{ }}` to `[% %]` delimiters to prevent XML parsing errors with `docxtemplater`.

## Templates by Node

### 1. S0_omid_application - OMiD Application

| Template | Languages | Delimiter Status | Prefill Status | Notes |
|----------|-----------|------------------|----------------|-------|
| `Zayavlenie_v_OMiD` | RU, KZ, EN | ✅ `[% %]` | ✅ Implemented | Generated via `create_clean_template.js` |

**Related Files:**
- Script: `frontend/scripts/create_clean_template.js`
- Templates: `frontend/public/templates/Zayavlenie_v_OMiD_{ru,kz,en}.docx`

**Prefill Fields:**
- Student name, specialty, program, dissertation topic
- Date fields (day, month, year)

---

### 2. S0_bioethics_letter - Bioethics Letter (LCB Request)

| Template | Languages | Delimiter Status | Prefill Status | Notes |
|----------|-----------|------------------|----------------|-------|
| `Letter_to_LCB` | RU, KZ, EN | ✅ `[% %]` | ✅ Implemented | Generated via `create_clean_lcb_template.js` |

**Related Files:**
- Script: `frontend/scripts/create_clean_lcb_template.js`
- Templates: `frontend/public/templates/Letter_to_LCB_{ru,kz,en}.docx`

**Prefill Fields:**
- Student name, specialty, dissertation topic
- Date fields (day, month, year)

---

### 3. S0_normocontrol_request - NCSTE Normocontrol Request

| Template | Languages | Delimiter Status | Prefill Status | Notes |
|----------|-----------|------------------|----------------|-------|
| `normocontrol_letter` | RU | ✅ `[% %]` | ✅ Implemented | Fixed via `fix_static_templates.js` |

**Related Files:**
- Fix Script: `frontend/scripts/fix_static_templates.js`
- Template: `frontend/public/templates/normocontrol_letter.docx`

**Prefill Fields:**
- Student name, specialty, dissertation topic
- University, department information
- Date fields

---

### 4. S1_publications_list - Publications List

| Template | Languages | Delimiter Status | Prefill Status | Notes |
|----------|-----------|------------------|----------------|-------|
| `app7` | RU, KZ, EN | ✅ `[% %]` | ✅ Implemented | Generated dynamically via `app7-templated.ts` |

**Related Files:**
- Generator: `frontend/src/features/docgen/app7-templated.ts`
- Templates: `frontend/public/templates/app7.{ru,kz,en}.docx`

**Prefill Fields:**
- Publications data (WOS/Scopus, KOKSON, conferences, IP)
- Student name, specialty, dissertation topic
- Computed publication counts

**Features:**
- Dynamic generation with publication list
- Supports all three languages
- Validation for DOI format and required fields

---

### 5. S2_ncste_publication_certificate - NCSTE Publication Certificate

| Template | Languages | Delimiter Status | Prefill Status | Notes |
|----------|-----------|------------------|----------------|-------|
| `letter_to_ncste_publication_certificate` | RU | ✅ `[% %]` | ✅ Implemented | Generated via `create_ncgnt_pub_letter.js` |

**Related Files:**
- Script: `frontend/scripts/create_ncgnt_pub_letter.js`
- Template: `frontend/public/templates/letter_to_ncste_publication_certificate_ru.docx`

**Prefill Fields:**
- Student name, birth year, IIN
- Specialty, dissertation topic
- Publications list with DOI links
- Date fields

---

### 6. NK_package - Defense Package (Приложения 4-9)

| Template | Languages | Delimiter Status | Prefill Status | Notes |
|----------|-----------|------------------|----------------|-------|
| `letter_to_rector_request_defense` (App 4) | RU, KZ, EN | ✅ `[% %]` | ✅ Implemented | Generated via `fix_rector_letter.js` |
| `letter_to_dissertation_head_request_defence` (App 5) | RU | ✅ `[% %]` | ✅ Implemented | New template replaces old `tpl_app5_ru.docx` |
| `tpl_app6_ru` (App 6) | RU | ❓ Not checked | ⏳ To be verified | Manual template |
| `tpl_app7_ru` (App 7) | RU | ❓ Not checked | ⏳ To be verified | See S1_publications_list |
| `tpl_app8_ru` (App 8) | RU | ❓ Not checked | ⏳ To be verified | Manual template |
| `tpl_app9_ru` (App 9) | RU | ❓ Not checked | ⏳ To be verified | Manual template |

**Related Files:**
- Script: `frontend/scripts/fix_rector_letter.js` (App 4)
- Template: `frontend/public/templates/letter_to_dissertation_head_request_defence_ru.docx` (App 5)
- Templates: `frontend/public/templates/tpl_app{5,6,7,8,9}_ru.docx` (legacy references)

**Prefill Fields (App 4 - Rector Letter):**
- Student name, specialty, program, dissertation topic
- Date fields
- Appendix list

**Prefill Fields (App 5 - Defense Letter to Dissertation Council):**
- Student name, specialty, dissertation topic
- Date fields
- Council/committee information

---

### 7. Reinstatement Templates

| Template | Languages | Delimiter Status | Prefill Status | Notes |
|----------|-----------|------------------|----------------|-------|
| `letter_to_rector_request_reinstatement` | RU, KZ, EN | ❓ Not checked | ⏳ To be verified | Static templates |
| `letter_to_dissertation_head_request_reinstatement` | RU, KZ, EN | ❓ Not checked | ⏳ To be verified | Static templates |

**Related Files:**
- Templates: `frontend/public/templates/letter_to_rector_request_reinstatement_{ru,kz,en}.docx`
- Templates: `frontend/public/templates/letter_to_dissertation_head_request_reinstatement_{ru,kz,en}.docx`

**Status:** These templates may still need delimiter migration and prefill implementation.

---

## Implementation Summary

### ✅ Fully Implemented (8 templates)
1. **Zayavlenie_v_OMiD** (RU, KZ, EN) - OMiD Application
2. **Letter_to_LCB** (RU, KZ, EN) - Bioethics Letter
3. **normocontrol_letter** (RU) - NCSTE Normocontrol
4. **app7** (RU, KZ, EN) - Publications List
5. **letter_to_ncste_publication_certificate** (RU) - NCSTE Pub Certificate
6. **letter_to_rector_request_defense** (RU, KZ, EN) - Rector Defense Letter
7. **tpl_app5_ru** (RU) - Appendix 5
8. **tpl_app7_ru** (RU) - Appendix 7 (same as app7)

### ⏳ To Be Verified (7 templates)
1. **tpl_app6_ru** (RU) - Appendix 6 (Advisor Reviews)
2. **tpl_app8_ru** (RU) - Appendix 8 (Abstracts)
3. **tpl_app9_ru** (RU) - Appendix 9 (Bioethics Conclusion)
4. **letter_to_rector_request_reinstatement** (RU, KZ, EN)
5. **letter_to_dissertation_head_request_reinstatement** (RU, KZ, EN)

---

## Backend Integration

### Profile Data Source

All template prefills pull data from the `GetProfile` function in `backend/internal/handlers/node_submission.go`:

**Data Sources:**
- User profile data (name, IIN, specialty, etc.)
- Stage data (dissertation topic, program)
- Publications data from `S1_publications_list` node

### Frontend Template Data Builder

The `buildTemplateData` function in `frontend/src/features/docgen/student-template.ts` aggregates:
- User information
- Profile snapshot
- Publications list
- Computed fields (dates, publication counts)

---

## Recent Changes

### Delimiter Standardization (Completed)
- ✅ Changed `docxtemplater` delimiters from `{{ }}` to `[% %]`
- ✅ Updated all generation scripts
- ✅ Fixed static templates (`normocontrol_letter.docx`, `tpl_app5_ru.docx`)
- ✅ Verified no remaining `{{` delimiters in templates

### Form Validation (Completed)
- ✅ Implemented DOI format validation (regex: `^10\.\d{4,9}/[-._;()/:A-Z0-9]+$`)
- ✅ Required field validation for publications
- ✅ Visual error feedback (red outlines) without data loss
- ✅ Persistent form storage using `localStorage`
- ✅ Auto-clear validation errors on form change

---

## Next Steps

1. **Verify Remaining Templates:**
   - Check `tpl_app6_ru.docx`, `tpl_app8_ru.docx`, `tpl_app9_ru.docx` for `[% %]` delimiters
   - Check reinstatement templates for delimiter migration

2. **Implement Prefills for Remaining Templates:**
   - Create generation scripts or manual prefill implementations
   - Test with real student data

3. **Testing:**
   - End-to-end testing of all document generation flows
   - Verify multilingual support (RU, KZ, EN)
   - Validate prefilled data accuracy

---

## Scripts Reference

| Script | Purpose | Templates Modified |
|--------|---------|-------------------|
| `create_clean_template.js` | Generate OMiD Application | `Zayavlenie_v_OMiD_{ru,kz,en}.docx` |
| `create_clean_lcb_template.js` | Generate LCB Request | `Letter_to_LCB_{ru,kz,en}.docx` |
| `fix_rector_letter.js` | Generate Rector Defense Letter | `letter_to_rector_request_defense_{ru,kz,en}.docx` |
| `create_ncgnt_pub_letter.js` | Generate NCSTE Pub Certificate | `letter_to_ncste_publication_certificate_ru.docx` |
| `fix_static_templates.js` | Fix static templates | `normocontrol_letter.docx`, `tpl_app5_ru.docx` |
| `scan_templates.js` | Scan for old delimiters | N/A (utility) |

---

## Known Issues

- None currently tracked

---

**Last Updated:** 2025-11-28
**Status:** In Progress

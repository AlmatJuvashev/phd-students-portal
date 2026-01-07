-- =============================================
-- DEMO UNIVERSITY: Complete Curriculum Data
-- =============================================

DO $$
DECLARE
    v_tenant_id uuid;
    v_term_fall_id uuid := gen_random_uuid();
    v_term_spring_id uuid := gen_random_uuid();
    v_dept_ph_id uuid;
    v_prog_phd_id uuid;
    v_course_ph501_id uuid := gen_random_uuid();
    v_course_epi601_id uuid := gen_random_uuid();
    v_offering_ph501_id uuid := gen_random_uuid();
    v_offering_epi601_id uuid := gen_random_uuid();
    v_instructor_johnson_id uuid;
    v_instructor_williams_id uuid;
    v_bank_ph501_id uuid := gen_random_uuid();
    v_building_id uuid := gen_random_uuid();
    v_room_id uuid := gen_random_uuid();
    
    -- Question IDs
    v_q_mcq_id uuid := gen_random_uuid();
    v_q_tf_id uuid := gen_random_uuid();
BEGIN
    -- 1. Get Tenant and References
    SELECT id INTO v_tenant_id FROM tenants WHERE name = 'Demo University' LIMIT 1;
    IF v_tenant_id IS NULL THEN
        -- Fallback to the default one if name doesn't match
        SELECT id INTO v_tenant_id FROM tenants LIMIT 1;
    END IF;

    -- Get Departments (assuming created by previous migrations)
    SELECT id INTO v_dept_ph_id FROM departments WHERE tenant_id = v_tenant_id LIMIT 1;
    -- Get Instructors
    SELECT id INTO v_instructor_johnson_id FROM users WHERE username = 'dr.johnson' LIMIT 1;
    SELECT id INTO v_instructor_williams_id FROM users WHERE username = 'dr.williams' LIMIT 1;
    
    -- If basics don't exist, we skip or they will fail constraints. 
    -- Assuming environment is set up from previous "demo data" steps.

    -- 2. Academic Terms
    INSERT INTO academic_terms (id, name, code, start_date, end_date, tenant_id)
    VALUES
      (v_term_fall_id, 'Fall 2025', 'FA25', '2025-09-01', '2025-12-15', v_tenant_id),
      (v_term_spring_id, 'Spring 2026', 'SP26', '2026-01-15', '2026-05-15', v_tenant_id)
    ON CONFLICT DO NOTHING;

    -- 3. Courses (Sample)
    INSERT INTO courses (id, code, title, credits, description, department_id, tenant_id)
    VALUES
      (v_course_ph501_id, 'PH501', '{"en":"Research Methods in Public Health"}', 3, '{"en":"Introduction to research design and methods."}', v_dept_ph_id, v_tenant_id),
      (v_course_epi601_id, 'EPI601', '{"en":"Advanced Epidemiology"}', 4, '{"en":"Methods for epidemiological studies."}', v_dept_ph_id, v_tenant_id)
    ON CONFLICT DO NOTHING;

    -- 4. Course Offerings
    INSERT INTO course_offerings (id, course_id, term_id, section, status, max_capacity, tenant_id)
    VALUES
      (v_offering_ph501_id, v_course_ph501_id, v_term_fall_id, 'A', 'ACTIVE', 30, v_tenant_id),
      (v_offering_epi601_id, v_course_epi601_id, v_term_fall_id, 'A', 'ACTIVE', 25, v_tenant_id)
    ON CONFLICT DO NOTHING;

    -- 5. Course Staff
    IF v_instructor_johnson_id IS NOT NULL THEN
        INSERT INTO course_staff (id, course_offering_id, user_id, role, is_primary)
        VALUES (gen_random_uuid(), v_offering_ph501_id, v_instructor_johnson_id, 'INSTRUCTOR', true);
    END IF;
    
    IF v_instructor_williams_id IS NOT NULL THEN
        INSERT INTO course_staff (id, course_offering_id, user_id, role, is_primary)
        VALUES (gen_random_uuid(), v_offering_epi601_id, v_instructor_williams_id, 'INSTRUCTOR', true);
    END IF;

    -- 6. Student Enrollments
    INSERT INTO course_enrollments (id, course_offering_id, student_id, status, enrolled_at)
    SELECT 
      gen_random_uuid(),
      v_offering_ph501_id,
      id,
      'ENROLLED',
      NOW()
    FROM users
    WHERE role = 'student'
    LIMIT 20
    ON CONFLICT DO NOTHING;

    -- 7. Item Banks and Questions (Only if instructor exists)
    IF v_instructor_johnson_id IS NOT NULL THEN
        INSERT INTO question_banks (id, title, description, tenant_id, created_by)
        VALUES
          (v_bank_ph501_id, 'Research Methods Question Bank', 'Questions for PH501', v_tenant_id, v_instructor_johnson_id)
        ON CONFLICT DO NOTHING;

        -- 8. Questions
        INSERT INTO questions (id, bank_id, type, stem, points_default)
        VALUES
          (v_q_mcq_id, v_bank_ph501_id, 'MCQ', 'What is the gold standard study design for causality?', 10),
          (v_q_tf_id, v_bank_ph501_id, 'TRUE_FALSE', 'Observational studies can establish causation.', 5)
        ON CONFLICT DO NOTHING;
        
        -- 8b. Options
        INSERT INTO question_options (question_id, text, is_correct, sort_order) VALUES
          (v_q_mcq_id, 'Case-control', false, 1),
          (v_q_mcq_id, 'RCT', true, 2),
          (v_q_mcq_id, 'Cross-sectional', false, 3),
          (v_q_tf_id, 'True', false, 1),
          (v_q_tf_id, 'False', true, 2)
        ON CONFLICT DO NOTHING;
    END IF;

    -- 9. Buildings & Rooms
    INSERT INTO buildings (id, name, address, tenant_id)
    VALUES (v_building_id, 'Health Sciences Building', '123 Medical Center Dr', v_tenant_id)
    ON CONFLICT DO NOTHING;

    INSERT INTO rooms (id, building_id, name, capacity, type)
    VALUES (v_room_id, v_building_id, 'Lecture Hall A', 100, 'LECTURE')
    ON CONFLICT DO NOTHING;

    -- 10. Class Sessions
    INSERT INTO class_sessions (id, course_offering_id, title, date, start_time, end_time, room_id, type)
    VALUES
      (gen_random_uuid(), v_offering_ph501_id, 'Research Methods Lecture 1', '2025-09-10', '09:00', '10:30', v_room_id, 'LECTURE'),
      (gen_random_uuid(), v_offering_ph501_id, 'Research Methods Lecture 2', '2025-09-12', '09:00', '10:30', v_room_id, 'LECTURE');

END $$;

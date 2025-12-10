-- Migration: Populate missing profile data for demo students
-- This fixes the issue where the S1 Profile node is empty because the user record lacks these details.

-- Update single demo student 1 for immediate testing
UPDATE users 
SET 
  program = 'PhD in Public Health',
  specialty = 'Public Health Policy',
  department = 'Department of Health Policy',
  cohort = 'Cohort 2024'
WHERE username = 'demo.student1';

-- Update other students with generic data so they aren't empty
UPDATE users 
SET 
  program = 'PhD in Public Health',
  specialty = 'Epidemiology',
  department = 'Department of Epidemiology',
  cohort = 'Cohort 2023'
WHERE username LIKE 'demo.student%' AND username != 'demo.student1';

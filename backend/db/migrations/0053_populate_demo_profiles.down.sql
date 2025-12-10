-- Revert profile population
UPDATE users 
SET 
  program = NULL,
  specialty = NULL,
  department = NULL,
  cohort = NULL
WHERE username LIKE 'demo.student%';

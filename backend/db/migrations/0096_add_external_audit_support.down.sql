DROP INDEX IF EXISTS idx_audit_access_tokens_user;
DROP INDEX IF EXISTS idx_curriculum_change_log_entity;
DROP INDEX IF EXISTS idx_curriculum_change_log_tenant;
DROP INDEX IF EXISTS idx_learning_outcomes_course;
DROP INDEX IF EXISTS idx_learning_outcomes_program;
DROP INDEX IF EXISTS idx_learning_outcomes_tenant;

DROP TABLE IF EXISTS audit_access_tokens;
DROP TABLE IF EXISTS curriculum_change_log;
DROP TABLE IF EXISTS outcome_assessments;
DROP TABLE IF EXISTS learning_outcomes;

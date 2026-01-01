DROP TABLE IF EXISTS cohorts;
DROP TABLE IF EXISTS journey_node_definitions;
DROP TABLE IF EXISTS journey_maps;
DROP TABLE IF EXISTS courses;

ALTER TABLE programs DROP COLUMN IF NOT EXISTS duration_months;
ALTER TABLE programs DROP COLUMN IF NOT EXISTS credits;
ALTER TABLE programs DROP COLUMN IF NOT EXISTS description;
ALTER TABLE programs DROP COLUMN IF NOT EXISTS title;
-- We do NOT drop tenant_id as it might be required by other migrations or future logic, strictly speaking we should look at dependencies but for down migration it is safer to leave it or cascade carefully.
-- But to be strict reverse:
-- ALTER TABLE programs DROP COLUMN IF NOT EXISTS tenant_id; 

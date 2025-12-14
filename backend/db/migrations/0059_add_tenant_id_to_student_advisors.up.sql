-- Add tenant_id to student_advisors
ALTER TABLE student_advisors ADD COLUMN IF NOT EXISTS tenant_id uuid REFERENCES tenants(id);
UPDATE student_advisors SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
ALTER TABLE student_advisors ALTER COLUMN tenant_id SET NOT NULL;
CREATE INDEX idx_student_advisors_tenant ON student_advisors(tenant_id);

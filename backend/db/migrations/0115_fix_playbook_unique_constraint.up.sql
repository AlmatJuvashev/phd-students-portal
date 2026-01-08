-- Fix unique constraint on playbook_versions to be tenant-aware
ALTER TABLE playbook_versions DROP CONSTRAINT IF EXISTS playbook_versions_checksum_key;
ALTER TABLE playbook_versions ADD CONSTRAINT playbook_versions_checksum_tenant_key UNIQUE (checksum, tenant_id);

-- Revert unique constraint on playbook_versions
ALTER TABLE playbook_versions DROP CONSTRAINT IF EXISTS playbook_versions_checksum_tenant_key;
ALTER TABLE playbook_versions ADD CONSTRAINT playbook_versions_checksum_key UNIQUE (checksum);

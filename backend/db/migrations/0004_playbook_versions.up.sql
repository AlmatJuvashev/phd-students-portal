CREATE TABLE IF NOT EXISTS playbook_versions (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  version text NOT NULL,
  checksum text NOT NULL UNIQUE,
  raw_json jsonb NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

-- Track which version is active (single row table referencing playbook_versions)
CREATE TABLE IF NOT EXISTS playbook_active_version (
  id boolean PRIMARY KEY DEFAULT TRUE,
  playbook_version_id uuid NOT NULL REFERENCES playbook_versions(id) ON DELETE CASCADE,
  updated_at timestamptz NOT NULL DEFAULT now()
);

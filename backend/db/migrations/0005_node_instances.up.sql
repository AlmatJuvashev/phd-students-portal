DO $$
BEGIN
    CREATE TYPE slot_multiplicity AS ENUM ('single', 'multi');
EXCEPTION
    WHEN duplicate_object THEN
        NULL;
END
$$;

CREATE TABLE IF NOT EXISTS node_instances (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  playbook_version_id uuid NOT NULL REFERENCES playbook_versions(id) ON DELETE CASCADE,
  node_id text NOT NULL,
  state text NOT NULL DEFAULT 'active',
  opened_at timestamptz NOT NULL DEFAULT now(),
  submitted_at timestamptz,
  updated_at timestamptz NOT NULL DEFAULT now(),
  locale text,
  current_rev int NOT NULL DEFAULT 0,
  UNIQUE (user_id, playbook_version_id, node_id)
);

CREATE TABLE IF NOT EXISTS node_instance_form_revisions (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  node_instance_id uuid NOT NULL REFERENCES node_instances(id) ON DELETE CASCADE,
  rev int NOT NULL,
  form_data jsonb NOT NULL,
  edited_by uuid NOT NULL REFERENCES users(id),
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (node_instance_id, rev)
);

CREATE TABLE IF NOT EXISTS node_instance_slots (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  node_instance_id uuid NOT NULL REFERENCES node_instances(id) ON DELETE CASCADE,
  slot_key text NOT NULL,
  required boolean NOT NULL DEFAULT false,
  multiplicity slot_multiplicity NOT NULL DEFAULT 'single',
  mime_whitelist text[] NOT NULL DEFAULT ARRAY[]::text[],
  UNIQUE (node_instance_id, slot_key)
);

CREATE TABLE IF NOT EXISTS node_instance_slot_attachments (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  slot_id uuid NOT NULL REFERENCES node_instance_slots(id) ON DELETE CASCADE,
  document_version_id uuid NOT NULL REFERENCES document_versions(id) ON DELETE CASCADE,
  filename text NOT NULL,
  size_bytes int NOT NULL,
  attached_by uuid NOT NULL REFERENCES users(id),
  attached_at timestamptz NOT NULL DEFAULT now(),
  is_active boolean NOT NULL DEFAULT true
);

CREATE INDEX IF NOT EXISTS idx_node_instances_user_node ON node_instances(user_id, playbook_version_id, node_id);
CREATE INDEX IF NOT EXISTS idx_node_instance_slots_instance_key ON node_instance_slots(node_instance_id, slot_key);
CREATE INDEX IF NOT EXISTS idx_node_instance_slot_attachments_slot ON node_instance_slot_attachments(slot_id) WHERE is_active;

CREATE TABLE IF NOT EXISTS node_outcomes (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  node_instance_id uuid NOT NULL REFERENCES node_instances(id) ON DELETE CASCADE,
  outcome_value text NOT NULL,
  decided_by uuid NOT NULL REFERENCES users(id),
  note text,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS node_events (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  node_instance_id uuid NOT NULL REFERENCES node_instances(id) ON DELETE CASCADE,
  event_type text NOT NULL,
  payload jsonb NOT NULL DEFAULT '{}'::jsonb,
  actor_id uuid NOT NULL REFERENCES users(id),
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_node_events_instance_created_at ON node_events(node_instance_id, created_at DESC);

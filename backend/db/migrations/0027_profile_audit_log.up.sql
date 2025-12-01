CREATE TABLE profile_audit_log (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  field_name text NOT NULL,
  old_value text,
  new_value text,
  changed_at timestamptz NOT NULL DEFAULT now(),
  changed_by uuid REFERENCES users(id)
);

CREATE INDEX idx_profile_audit_user ON profile_audit_log(user_id);
CREATE INDEX idx_profile_audit_date ON profile_audit_log(changed_at);

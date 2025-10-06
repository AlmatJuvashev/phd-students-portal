CREATE TABLE IF NOT EXISTS node_state_transitions (
  from_state text NOT NULL,
  to_state text NOT NULL,
  allowed_roles text[] NOT NULL,
  PRIMARY KEY (from_state, to_state)
);

INSERT INTO node_state_transitions(from_state, to_state, allowed_roles) VALUES
  ('active','submitted', ARRAY['student']),
  ('submitted','needs_fixes', ARRAY['advisor','secretary','chair','admin']),
  ('submitted','done', ARRAY['advisor','secretary','chair','admin']),
  ('needs_fixes','submitted', ARRAY['student']),
  ('done','submitted', ARRAY['admin'])
ON CONFLICT (from_state, to_state)
DO UPDATE SET allowed_roles = EXCLUDED.allowed_roles;

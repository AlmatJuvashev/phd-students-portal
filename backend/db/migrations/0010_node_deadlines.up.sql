CREATE TABLE IF NOT EXISTS node_deadlines (
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  node_id text NOT NULL,
  due_at timestamptz NOT NULL,
  note text,
  created_by uuid NOT NULL REFERENCES users(id),
  created_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, node_id)
);

CREATE INDEX IF NOT EXISTS idx_node_deadlines_due ON node_deadlines(due_at);


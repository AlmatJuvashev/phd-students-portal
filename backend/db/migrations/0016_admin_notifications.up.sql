-- Create admin notifications table
CREATE TABLE IF NOT EXISTS admin_notifications (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  student_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  node_id text NOT NULL,
  node_instance_id uuid REFERENCES node_instances(id) ON DELETE CASCADE,
  event_type text NOT NULL, -- 'document_submitted', 'document_uploaded', 'state_changed', etc.
  is_read boolean NOT NULL DEFAULT false,
  message text NOT NULL,
  metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now()
);

-- Indexes for efficient queries
CREATE INDEX IF NOT EXISTS idx_notifications_unread ON admin_notifications(is_read, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_student ON admin_notifications(student_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_instance ON admin_notifications(node_instance_id);

-- Trigger function to create notifications from node_events
CREATE OR REPLACE FUNCTION create_admin_notification_from_event()
RETURNS TRIGGER AS $$
DECLARE
  v_student_id uuid;
  v_node_id text;
  v_student_name text;
  v_message text;
BEGIN
  -- Get student_id and node_id from node_instance
  SELECT ni.user_id, ni.node_id
  INTO v_student_id, v_node_id
  FROM node_instances ni
  WHERE ni.id = NEW.node_instance_id;

  -- Get student name
  SELECT first_name || ' ' || last_name
  INTO v_student_name
  FROM users
  WHERE id = v_student_id;

  -- Create notification based on event type
  IF NEW.event_type = 'state_changed' AND (NEW.payload->>'new_state') = 'under_review' THEN
    v_message := v_student_name || ' отправил(а) документы на проверку в узле ' || v_node_id;
    INSERT INTO admin_notifications (student_id, node_id, node_instance_id, event_type, message, metadata)
    VALUES (v_student_id, v_node_id, NEW.node_instance_id, 'document_submitted', v_message, NEW.payload);
  
  ELSIF NEW.event_type = 'attachment_uploaded' THEN
    v_message := v_student_name || ' загрузил(а) документ в узел ' || v_node_id;
    INSERT INTO admin_notifications (student_id, node_id, node_instance_id, event_type, message, metadata)
    VALUES (v_student_id, v_node_id, NEW.node_instance_id, 'document_uploaded', v_message, NEW.payload);
  
  ELSIF NEW.event_type = 'form_updated' AND (NEW.payload->>'submitted') = 'true' THEN
    v_message := v_student_name || ' обновил(а) форму в узле ' || v_node_id;
    INSERT INTO admin_notifications (student_id, node_id, node_instance_id, event_type, message, metadata)
    VALUES (v_student_id, v_node_id, NEW.node_instance_id, 'form_submitted', v_message, NEW.payload);
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger on node_events
CREATE TRIGGER trigger_create_admin_notification
  AFTER INSERT ON node_events
  FOR EACH ROW
  EXECUTE FUNCTION create_admin_notification_from_event();

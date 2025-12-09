-- Rollback to previous trigger version (buggy one from 0016)
-- This is for rollback purposes only

CREATE OR REPLACE FUNCTION create_admin_notification_from_event()
RETURNS TRIGGER AS $$
DECLARE
  student_name TEXT;
  node_label TEXT;
  message TEXT;
BEGIN
  -- Get student name
  SELECT COALESCE(first_name || ' ' || last_name, email) 
  INTO student_name 
  FROM users 
  WHERE id = NEW.user_id;

  -- Determine message based on event type (old buggy version)
  IF NEW.event_type = 'state_changed' THEN
    IF NEW.payload->>'new_state' IN ('submitted', 'under_review') THEN
      node_label := COALESCE(NEW.node_id, 'узел');
      message := student_name || ' отправил(а) узел ' || node_label || ' на проверку';
      
      INSERT INTO admin_notifications (node_instance_id, event_id, message, is_read, created_at)
      VALUES (NEW.node_instance_id, NEW.id, message, FALSE, NEW.created_at);
    END IF;
    
  ELSIF NEW.event_type = 'attachment_uploaded' THEN
    node_label := COALESCE(NEW.node_id, 'узел');
    message := student_name || ' загрузил(а) документ в узел ' || node_label;
    
    INSERT INTO admin_notifications (node_instance_id, event_id, message, is_read, created_at)
    VALUES (NEW.node_instance_id, NEW.id, message, FALSE, NEW.created_at);
    
  ELSIF NEW.event_type = 'form_updated' THEN
    node_label := COALESCE(NEW.node_id, 'узел');
    message := student_name || ' обновил(а) форму в узле ' || node_label;
    
    INSERT INTO admin_notifications (node_instance_id, event_id, message, is_read, created_at)
    VALUES (NEW.node_instance_id, NEW.id, message, FALSE, NEW.created_at);
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

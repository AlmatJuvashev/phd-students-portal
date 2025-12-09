-- Restore old trigger function with event_id
-- This is a down migration, so we restore the previous state

CREATE OR REPLACE FUNCTION create_admin_notification_from_event()
RETURNS TRIGGER AS $$
DECLARE
  student_name TEXT;
  node_label TEXT;
  message TEXT;
  v_user_id UUID;
  v_node_id TEXT;
BEGIN
  SELECT user_id, node_id 
  INTO v_user_id, v_node_id
  FROM node_instances 
  WHERE id = NEW.node_instance_id;

  SELECT COALESCE(first_name || ' ' || last_name, email) 
  INTO student_name 
  FROM users 
  WHERE id = v_user_id;

  IF NEW.event_type = 'state_changed' THEN
    IF NEW.payload->>'to' IN ('submitted', 'under_review') THEN
      node_label := COALESCE(v_node_id, 'узел');
      message := student_name || ' отправил(а) узел ' || node_label || ' на проверку';
      
      INSERT INTO admin_notifications (node_instance_id, event_id, message, is_read, created_at)
      VALUES (NEW.node_instance_id, NEW.id, message, FALSE, NEW.created_at);
    END IF;
    
  ELSIF NEW.event_type = 'file_attached' THEN
    node_label := COALESCE(v_node_id, 'узел');
    message := student_name || ' загрузил(а) документ в узел ' || node_label;
    
    INSERT INTO admin_notifications (node_instance_id, event_id, message, is_read, created_at)
    VALUES (NEW.node_instance_id, NEW.id, message, FALSE, NEW.created_at);
    
  ELSIF NEW.event_type = 'form_updated' THEN
    node_label := COALESCE(v_node_id, 'узел');
    message := student_name || ' обновил(а) форму в узле ' || node_label;
    
    INSERT INTO admin_notifications (node_instance_id, event_id, message, is_read, created_at)
    VALUES (NEW.node_instance_id, NEW.id, message, FALSE, NEW.created_at);
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

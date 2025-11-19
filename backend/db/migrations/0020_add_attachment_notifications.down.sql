-- Revert to previous trigger without attachment_uploaded

CREATE OR REPLACE FUNCTION create_admin_notification_from_event()
RETURNS TRIGGER AS $$
DECLARE
  student_name TEXT;
  node_label TEXT;
  message TEXT;
  v_user_id UUID;
  v_node_id TEXT;
BEGIN
  -- Get user_id and node_id from node_instances
  SELECT user_id, node_id 
  INTO v_user_id, v_node_id
  FROM node_instances 
  WHERE id = NEW.node_instance_id;

  -- Get student name
  SELECT COALESCE(first_name || ' ' || last_name, email) 
  INTO student_name 
  FROM users 
  WHERE id = v_user_id;

  -- Determine message based on event type
  IF NEW.event_type = 'state_changed' THEN
    -- State changed events - check 'to' field in payload
    IF NEW.payload->>'to' IN ('submitted', 'under_review') THEN
      node_label := COALESCE(v_node_id, 'узел');
      message := student_name || ' отправил(а) узел ' || node_label || ' на проверку';
      
      INSERT INTO admin_notifications (student_id, node_id, node_instance_id, event_type, message, metadata)
      VALUES (v_user_id, v_node_id, NEW.node_instance_id, NEW.event_type, message, NEW.payload);
    END IF;
    
  ELSIF NEW.event_type = 'form_updated' THEN
    -- Form update events
    node_label := COALESCE(v_node_id, 'узел');
    message := student_name || ' обновил(а) форму в узле ' || node_label;
    
    INSERT INTO admin_notifications (student_id, node_id, node_instance_id, event_type, message, metadata)
    VALUES (v_user_id, v_node_id, NEW.node_instance_id, NEW.event_type, message, NEW.payload);
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_admin_notification_on_event ON node_events;
DROP TRIGGER IF EXISTS trigger_create_admin_notification ON node_events;
-- Migration skipped to preserve notification trigger
-- The application code does not yet handle admin_notifications insertion manually.
-- Once the Go backend handles this, we can uncomment the lines below.

-- DROP TRIGGER IF EXISTS trigger_admin_notification_on_event ON node_events;
-- DROP TRIGGER IF EXISTS trigger_create_admin_notification ON node_events;
-- DROP TRIGGER IF EXISTS on_event_create_notification ON events;
-- DROP FUNCTION IF EXISTS create_admin_notification_from_event();

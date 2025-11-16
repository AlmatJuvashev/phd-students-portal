-- Drop trigger and function
DROP TRIGGER IF EXISTS trigger_create_admin_notification ON node_events;
DROP FUNCTION IF EXISTS create_admin_notification_from_event();

-- Drop indexes
DROP INDEX IF EXISTS idx_notifications_instance;
DROP INDEX IF EXISTS idx_notifications_student;
DROP INDEX IF EXISTS idx_notifications_unread;

-- Drop table
DROP TABLE IF EXISTS admin_notifications;

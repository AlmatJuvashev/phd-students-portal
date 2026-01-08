-- Enum values cannot be removed safely.
-- Removing metadata is optional but safer to leave in case of references.
DELETE FROM role_metadata WHERE role_name IN ('ta', 'chairman', 'hr', 'facility_manager');

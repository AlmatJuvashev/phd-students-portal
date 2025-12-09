-- Rollback advisor reviewed document fields
DROP INDEX IF EXISTS idx_slot_attachments_reviewed_doc;

ALTER TABLE node_instance_slot_attachments
    DROP COLUMN IF EXISTS reviewed_document_version_id,
    DROP COLUMN IF EXISTS reviewed_by,
    DROP COLUMN IF EXISTS reviewed_at;

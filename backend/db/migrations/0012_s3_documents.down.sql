DROP INDEX IF EXISTS idx_slot_attachments_version;
DROP INDEX IF EXISTS idx_slot_attachments_status;

ALTER TABLE node_instance_slot_attachments
    DROP COLUMN IF EXISTS review_note,
    DROP COLUMN IF EXISTS approved_at,
    DROP COLUMN IF EXISTS approved_by,
    DROP COLUMN IF EXISTS status;

ALTER TABLE document_versions
    DROP COLUMN IF EXISTS metadata,
    DROP COLUMN IF EXISTS checksum,
    DROP COLUMN IF EXISTS etag,
    DROP COLUMN IF EXISTS object_key,
    DROP COLUMN IF EXISTS bucket;

DROP TYPE IF EXISTS attachment_review_status;

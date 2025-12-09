DO $$
BEGIN
    CREATE TYPE attachment_review_status AS ENUM ('submitted', 'approved', 'rejected');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END$$;

ALTER TABLE document_versions
    ADD COLUMN IF NOT EXISTS bucket text,
    ADD COLUMN IF NOT EXISTS object_key text,
    ADD COLUMN IF NOT EXISTS etag text,
    ADD COLUMN IF NOT EXISTS checksum text,
    ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;

UPDATE document_versions
SET object_key = COALESCE(object_key, storage_path)
WHERE object_key IS NULL;

ALTER TABLE node_instance_slot_attachments
    ADD COLUMN IF NOT EXISTS status attachment_review_status NOT NULL DEFAULT 'submitted',
    ADD COLUMN IF NOT EXISTS approved_by uuid REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS approved_at timestamptz,
    ADD COLUMN IF NOT EXISTS review_note text;

CREATE INDEX IF NOT EXISTS idx_slot_attachments_status ON node_instance_slot_attachments(status)
    WHERE is_active;
CREATE INDEX IF NOT EXISTS idx_slot_attachments_version ON node_instance_slot_attachments(document_version_id);

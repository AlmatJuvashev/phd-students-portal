-- Add ability for advisors to upload reviewed/commented documents
ALTER TABLE node_instance_slot_attachments
    ADD COLUMN IF NOT EXISTS reviewed_document_version_id uuid REFERENCES document_versions(id),
    ADD COLUMN IF NOT EXISTS reviewed_by uuid REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS reviewed_at timestamptz;

CREATE INDEX IF NOT EXISTS idx_slot_attachments_reviewed_doc ON node_instance_slot_attachments(reviewed_document_version_id)
    WHERE reviewed_document_version_id IS NOT NULL;

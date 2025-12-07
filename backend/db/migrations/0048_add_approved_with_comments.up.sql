-- Add approved_with_comments to attachment_review_status enum
-- This allows approving documents with minor feedback that doesn't block the student

ALTER TYPE attachment_review_status ADD VALUE IF NOT EXISTS 'approved_with_comments' AFTER 'approved';

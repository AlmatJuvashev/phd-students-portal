-- Remove approved_with_comments from attachment_review_status enum
-- NOTE: PostgreSQL doesn't support removing values from enums directly
-- This would require recreating the enum which is complex
-- The safest down migration is to do nothing

-- No action needed

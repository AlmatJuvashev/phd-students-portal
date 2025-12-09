-- Add indexes for pagination performance

-- Index for sorting users by last_name (most common sort)
CREATE INDEX IF NOT EXISTS idx_users_last_name_active 
ON users(last_name, first_name) 
WHERE is_active = true;

-- Index for filtering by role
CREATE INDEX IF NOT EXISTS idx_users_role_active 
ON users(role, is_active);

-- Index for profile_submissions subqueries
CREATE INDEX IF NOT EXISTS idx_profile_submissions_user_submitted 
ON profile_submissions(user_id, submitted_at DESC);

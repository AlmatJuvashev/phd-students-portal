-- Restore password reset tokens table (rollback)
CREATE TABLE password_reset_tokens (
  email text PRIMARY KEY,
  token text NOT NULL,
  expires_at timestamptz NOT NULL
);

-- Irreversible action for safety, or simple no-op since we can't easily know which were false before.
-- We will just do nothing as rolling back "fixing data" is usually not desired.
SELECT 1;

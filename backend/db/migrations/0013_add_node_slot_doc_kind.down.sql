-- Down migration: removing enum values in Postgres enum types is non-trivial
-- and not safely reversible in place without reconstructing the type.
-- Leave as no-op to avoid accidental data loss. If you need to remove the
-- value, create a careful migration that recreates the type and updates
-- dependent columns.
/* no-op */

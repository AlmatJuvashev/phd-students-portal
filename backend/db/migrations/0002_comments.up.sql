-- comments threading & mentions
ALTER TABLE comments ADD COLUMN parent_id uuid NULL;
ALTER TABLE comments ADD COLUMN mentions uuid[] NULL;

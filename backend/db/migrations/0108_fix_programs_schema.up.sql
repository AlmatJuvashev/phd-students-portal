UPDATE programs SET title = '{}'::jsonb WHERE title IS NULL;
UPDATE programs SET description = '{}'::jsonb WHERE description IS NULL;

ALTER TABLE programs ALTER COLUMN title SET DEFAULT '{}'::jsonb;
ALTER TABLE programs ALTER COLUMN title SET NOT NULL;

ALTER TABLE programs ALTER COLUMN description SET DEFAULT '{}'::jsonb;
ALTER TABLE programs ALTER COLUMN description SET NOT NULL;

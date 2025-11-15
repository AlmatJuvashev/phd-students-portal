DO $$
BEGIN
    -- Add node_slot kind for documents created from journey node slots
    ALTER TYPE doc_kind ADD VALUE 'node_slot';
EXCEPTION
    WHEN duplicate_object THEN NULL;
END$$;

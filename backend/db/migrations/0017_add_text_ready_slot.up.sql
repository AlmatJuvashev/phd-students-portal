-- Add upload slot for “Текст диссертации подготовлен” node to existing instances
INSERT INTO node_instance_slots (
    node_instance_id,
    slot_key,
    required,
    multiplicity,
    mime_whitelist
)
SELECT
    ni.id,
    'dissertation_draft_file',
    true,
    'single',
    ARRAY[
        'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
        'application/msword',
        'application/pdf'
    ]::text[]
FROM node_instances ni
LEFT JOIN node_instance_slots s
    ON s.node_instance_id = ni.id
   AND s.slot_key = 'dissertation_draft_file'
WHERE ni.node_id = 'S1_text_ready'
  AND s.id IS NULL;

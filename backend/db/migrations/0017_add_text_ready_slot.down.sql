-- Remove dissertation draft upload slots for S1_text_ready node
DELETE FROM node_instance_slots sis
USING node_instances ni
WHERE sis.node_instance_id = ni.id
  AND ni.node_id = 'S1_text_ready'
  AND sis.slot_key = 'dissertation_draft_file';

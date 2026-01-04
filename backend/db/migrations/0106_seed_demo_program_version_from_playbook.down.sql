-- Remove demo program version seeded from playbook.json

DELETE FROM programs
WHERE id = 'dd200009-0000-0000-0009-000000000009'::uuid;

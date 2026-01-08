-- Add new roles to user_role enum
ALTER TYPE user_role ADD VALUE IF NOT EXISTS 'ta';
ALTER TYPE user_role ADD VALUE IF NOT EXISTS 'chairman';
ALTER TYPE user_role ADD VALUE IF NOT EXISTS 'hr';
ALTER TYPE user_role ADD VALUE IF NOT EXISTS 'facility_manager';

-- Add metadata for new roles
INSERT INTO role_metadata (role_name, category, description) VALUES
('ta', 'teaching', 'Teaching Assistant'),
('chairman', 'academic', 'Department Chairman'),
('hr', 'administrative', 'Human Resources Manager'),
('facility_manager', 'administrative', 'Facility & Resources Manager')
ON CONFLICT (role_name) DO NOTHING;

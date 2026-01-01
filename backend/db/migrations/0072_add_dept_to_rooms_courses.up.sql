-- Add department_id to courses
ALTER TABLE courses
  ADD COLUMN department_id UUID REFERENCES departments(id) ON DELETE SET NULL;

-- Add department_id and floor to rooms
ALTER TABLE rooms
  ADD COLUMN department_id UUID REFERENCES departments(id) ON DELETE SET NULL,
  ADD COLUMN floor INTEGER DEFAULT 1;

-- Create index for faster filtering
CREATE INDEX idx_courses_department ON courses(department_id);
CREATE INDEX idx_rooms_department ON rooms(department_id);

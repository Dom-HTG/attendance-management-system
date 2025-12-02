-- Add metadata fields to events table for analytics
-- These fields support lecturer and admin dashboard analytics

-- Add lecturer_id to track who created the event
ALTER TABLE events ADD COLUMN IF NOT EXISTS lecturer_id BIGINT;

-- Add course information
ALTER TABLE events ADD COLUMN IF NOT EXISTS course_code VARCHAR(20);
ALTER TABLE events ADD COLUMN IF NOT EXISTS course_name TEXT;

-- Add department for analytics grouping
ALTER TABLE events ADD COLUMN IF NOT EXISTS department VARCHAR(100);

-- Add foreign key constraint to lecturers table
ALTER TABLE events ADD CONSTRAINT fk_events_lecturer 
    FOREIGN KEY (lecturer_id) REFERENCES lecturers(id);

-- Add indexes for performance on analytics queries
CREATE INDEX IF NOT EXISTS idx_events_lecturer_id ON events(lecturer_id);
CREATE INDEX IF NOT EXISTS idx_events_department ON events(department);
CREATE INDEX IF NOT EXISTS idx_events_course_code ON events(course_code);
CREATE INDEX IF NOT EXISTS idx_events_start_time ON events(start_time DESC);
CREATE INDEX IF NOT EXISTS idx_events_created_at ON events(created_at DESC);

COMMENT ON COLUMN events.lecturer_id IS 'ID of the lecturer who created this event';
COMMENT ON COLUMN events.course_code IS 'Course code (e.g., CSC301)';
COMMENT ON COLUMN events.course_name IS 'Full course name';
COMMENT ON COLUMN events.department IS 'Department offering the course';

-- Analytics Indexes Migration
-- These indexes optimize analytics query performance

-- Index on student_id + marked_time for fast student query lookups
CREATE INDEX IF NOT EXISTS idx_user_attendances_student_marked 
ON user_attendances(student_id, marked_time DESC);

-- Index on event_id + status for course/event analytics
CREATE INDEX IF NOT EXISTS idx_user_attendances_event_status 
ON user_attendances(event_id, status);

-- Index on marked_time for temporal queries
CREATE INDEX IF NOT EXISTS idx_user_attendances_marked_time 
ON user_attendances(marked_time DESC);

-- Composite index for department analytics (if students table has department)
CREATE INDEX IF NOT EXISTS idx_students_department 
ON students(department);

-- Composite index for lecturer department queries
CREATE INDEX IF NOT EXISTS idx_lecturers_department 
ON lecturers(department);

-- Index on event dates for temporal analysis
CREATE INDEX IF NOT EXISTS idx_events_start_end_time 
ON events(start_time, end_time);

-- Index for real-time dashboard queries
CREATE INDEX IF NOT EXISTS idx_user_attendances_created_at 
ON user_attendances(created_at DESC);

-- COMMENT: These indexes are automatically created by GORM migration if using embedded model timestamps.
-- Ensure that user_attendances table has proper indexes on foreign keys:
CREATE INDEX IF NOT EXISTS idx_user_attendances_student_id 
ON user_attendances(student_id);

CREATE INDEX IF NOT EXISTS idx_user_attendances_event_id 
ON user_attendances(event_id);

-- Migration: Add Admin, AuditLog, and SystemSettings tables
-- Created: 2025-12-01
-- Description: Adds admin functionality with audit logging and system settings

-- Create admin table
CREATE TABLE IF NOT EXISTS admins (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    department VARCHAR(100),
    role VARCHAR(20) DEFAULT 'admin',
    is_super_admin BOOLEAN DEFAULT FALSE,
    active BOOLEAN DEFAULT TRUE
);

CREATE INDEX IF NOT EXISTS idx_admins_email ON admins(email);
CREATE INDEX IF NOT EXISTS idx_admins_deleted_at ON admins(deleted_at);

-- Create audit_logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_type VARCHAR(20) NOT NULL,
    user_id INTEGER NOT NULL,
    user_email VARCHAR(255),
    action VARCHAR(50) NOT NULL,
    resource_type VARCHAR(50),
    resource_id INTEGER,
    details TEXT,
    ip_address VARCHAR(45),
    user_agent TEXT
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_timestamp ON audit_logs(timestamp);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user ON audit_logs(user_type, user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_deleted_at ON audit_logs(deleted_at);

-- Create system_settings table
CREATE TABLE IF NOT EXISTS system_settings (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    setting_key VARCHAR(100) UNIQUE NOT NULL,
    setting_value TEXT NOT NULL,
    data_type VARCHAR(20) NOT NULL,
    description TEXT,
    updated_by INTEGER
);

CREATE INDEX IF NOT EXISTS idx_system_settings_key ON system_settings(setting_key);
CREATE INDEX IF NOT EXISTS idx_system_settings_deleted_at ON system_settings(deleted_at);

-- Insert default admin user (password: Admin@2024)
-- Password will need to be updated after running migration
-- Use: UPDATE admins SET password = '<bcrypt_hash>' WHERE email = 'admin@fupre.edu.ng';
-- Generate hash with: echo -n "Admin@2024" | htpasswd -bnBC 10 "" | tr -d ':\n'
INSERT INTO admins (first_name, last_name, email, password, department, is_super_admin, active)
VALUES (
    'System',
    'Administrator',
    'admin@fupre.edu.ng',
    '$2a$10$YourHashWillBeGeneratedDuringSeeding',
    'Administration',
    TRUE,
    TRUE
)
ON CONFLICT (email) DO UPDATE SET
    password = EXCLUDED.password,
    updated_at = CURRENT_TIMESTAMP;

-- Insert default system settings
INSERT INTO system_settings (setting_key, setting_value, data_type, description) VALUES
('qr_code_validity_minutes', '30', 'number', 'QR code validity duration in minutes'),
('attendance_grace_period_minutes', '15', 'number', 'Grace period for late attendance'),
('low_attendance_threshold', '75', 'number', 'Minimum attendance percentage threshold'),
('academic_year', '2024/2025', 'string', 'Current academic year'),
('semester', 'First Semester', 'string', 'Current semester'),
('require_email_verification', 'false', 'boolean', 'Require email verification for new users'),
('allow_student_self_registration', 'false', 'boolean', 'Allow students to self-register'),
('max_events_per_day_per_lecturer', '5', 'number', 'Maximum events a lecturer can create per day')
ON CONFLICT (setting_key) DO NOTHING;

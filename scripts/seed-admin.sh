#!/bin/bash

# Admin User Seeder
# Creates the default admin user in the database

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}  Admin User Seeder${NC}"
echo -e "${BLUE}======================================${NC}"
echo ""

# Wait for database to be ready
echo -e "${YELLOW}Waiting for database connection...${NC}"
sleep 3

echo -e "${BLUE}Creating admin user in database...${NC}"

# Use pre-generated bcrypt hash for "Admin@2024"
# Generated using: bcrypt.GenerateFromPassword([]byte("Admin@2024"), 10)
HASH='$2a$10$LjllxwJnSpMqd4j3TaBZ2uZjoUV0Db7A6UKm6G33MDqUfQUMu4yqC'

# Insert admin user into database
docker exec attendance-management-db psql -U postgres -d attendance-management -c "
INSERT INTO admins (created_at, updated_at, first_name, last_name, email, password, department, role, is_super_admin, active)
VALUES (
    NOW(),
    NOW(),
    'System',
    'Administrator',
    'admin@fupre.edu.ng',
    '${HASH}',
    'Administration',
    'admin',
    TRUE,
    TRUE
)
ON CONFLICT (email) DO UPDATE SET
    password = EXCLUDED.password,
    updated_at = NOW();
" > /dev/null 2>&1

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Admin user created successfully${NC}"
    echo -e "  Email: admin@fupre.edu.ng"
    echo -e "  Password: Admin@2024"
    echo -e "  Role: admin"
else
    echo -e "${RED}✗ Failed to create admin user${NC}"
    exit 1
fi

# Seed system settings
echo ""
echo -e "${BLUE}Seeding system settings...${NC}"
docker exec attendance-management-db psql -U postgres -d attendance-management -c "
INSERT INTO system_settings (created_at, updated_at, setting_key, setting_value, data_type, description) VALUES 
(NOW(), NOW(), 'qr_code_validity_minutes', '30', 'number', 'QR code validity duration in minutes'),
(NOW(), NOW(), 'attendance_grace_period_minutes', '15', 'number', 'Grace period for late attendance'),
(NOW(), NOW(), 'low_attendance_threshold', '75', 'number', 'Minimum attendance percentage threshold'),
(NOW(), NOW(), 'academic_year', '2024/2025', 'string', 'Current academic year'),
(NOW(), NOW(), 'semester', 'First Semester', 'string', 'Current semester'),
(NOW(), NOW(), 'require_email_verification', 'false', 'boolean', 'Require email verification for new users'),
(NOW(), NOW(), 'allow_student_self_registration', 'false', 'boolean', 'Allow students to self-register'),
(NOW(), NOW(), 'max_events_per_day_per_lecturer', '5', 'number', 'Maximum events a lecturer can create per day')
ON CONFLICT (setting_key) DO NOTHING;
" > /dev/null 2>&1

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ System settings seeded${NC}"
else
    echo -e "${YELLOW}⚠ System settings may already exist${NC}"
fi

echo ""
echo -e "${GREEN}Admin user and settings are ready!${NC}"

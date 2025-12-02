# Admin Dashboard API Guide

This document provides comprehensive documentation for the Admin Dashboard API endpoints.

## Table of Contents
- [Authentication](#authentication)
- [User Management](#user-management)
- [Event Management](#event-management)
- [Analytics](#analytics)
- [System Settings](#system-settings)
- [Audit Logs](#audit-logs)
- [Testing](#testing)

## Authentication

### Admin Login
Login as an administrator to receive a JWT token valid for 7 days.

**Endpoint:** `POST /api/auth/login-admin`

**Request Body:**
```json
{
  "email": "admin@fupre.edu.ng",
  "password": "Admin@2024"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGc...",
    "user": {
      "id": 1,
      "first_name": "System",
      "last_name": "Administrator",
      "email": "admin@fupre.edu.ng",
      "role": "admin",
      "department": "Administration",
      "created_at": "2025-12-02T00:00:00Z"
    }
  }
}
```

**Default Credentials:**
- Email: `admin@fupre.edu.ng`
- Password: `Admin@2024`

---

## User Management

All user management endpoints require admin authentication. Include the JWT token in the Authorization header:
```
Authorization: Bearer <your_token>
```

### List All Students

**Endpoint:** `GET /api/admin/students`

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10, max: 100)
- `department` (optional): Filter by department
- `search` (optional): Search by name or matric number

**Example Request:**
```bash
curl -H "Authorization: Bearer <token>" \
  "http://localhost:2754/api/admin/students?page=1&limit=10&department=Computer%20Science"
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "students": [
      {
        "id": 1,
        "matric_number": "FUP/CSC/21/1234",
        "first_name": "John",
        "last_name": "Doe",
        "email": "john.doe@fupre.edu.ng",
        "department": "Computer Science",
        "total_attendance": 15,
        "total_possible": 20,
        "attendance_rate": 75.0,
        "status": "active",
        "created_at": "2025-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "current_page": 1,
      "per_page": 10,
      "total_pages": 5,
      "total_items": 50
    }
  }
}
```

### List All Lecturers

**Endpoint:** `GET /api/admin/lecturers`

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10, max: 100)
- `department` (optional): Filter by department

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "lecturers": [
      {
        "id": 1,
        "staff_id": "FUP/STAFF/001",
        "first_name": "Dr. Jane",
        "last_name": "Smith",
        "email": "jane.smith@fupre.edu.ng",
        "department": "Computer Science",
        "total_events": 10,
        "total_students_reached": 150,
        "status": "active",
        "created_at": "2025-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "current_page": 1,
      "per_page": 10,
      "total_pages": 2,
      "total_items": 15
    }
  }
}
```

### Get Student Details

**Endpoint:** `GET /api/admin/users/student/:user_id`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "matric_number": "FUP/CSC/21/1234",
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@fupre.edu.ng",
    "department": "Computer Science",
    "total_attendance": 15,
    "total_possible": 20,
    "attendance_rate": 75.0,
    "events_attended": [
      {
        "event_id": 1,
        "event_name": "Data Structures Lecture",
        "date": "2025-11-01",
        "time": "10:00 AM",
        "checked_in_at": "2025-11-01T10:05:00Z"
      }
    ],
    "status": "active",
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

### Get Lecturer Details

**Endpoint:** `GET /api/admin/users/lecturer/:user_id`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "staff_id": "FUP/STAFF/001",
    "first_name": "Dr. Jane",
    "last_name": "Smith",
    "email": "jane.smith@fupre.edu.ng",
    "department": "Computer Science",
    "total_events": 10,
    "total_students_reached": 150,
    "events_created": [
      {
        "event_id": 1,
        "event_name": "Data Structures Lecture",
        "date": "2025-11-01",
        "time": "10:00 AM",
        "total_students": 25,
        "attendance_count": 20
      }
    ],
    "status": "active",
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

### Update User Status

**Endpoint:** `PATCH /api/admin/users/:user_type/:user_id/status`

**Path Parameters:**
- `user_type`: Either "student" or "lecturer"
- `user_id`: The user's ID

**Request Body:**
```json
{
  "active": false,
  "reason": "Account suspended due to violation of terms"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "User status updated successfully",
  "data": {
    "user_id": 1,
    "user_type": "student",
    "active": false,
    "updated_at": "2025-12-02T03:12:00Z"
  }
}
```

### Delete User

**Endpoint:** `DELETE /api/admin/users/:user_type/:user_id`

**Path Parameters:**
- `user_type`: Either "student" or "lecturer"
- `user_id`: The user's ID

**Response (200 OK):**
```json
{
  "success": true,
  "message": "User deleted successfully"
}
```

**Note:** This performs a soft delete. The user record is marked as deleted but remains in the database.

---

## Event Management

### List All Events

**Endpoint:** `GET /api/admin/events`

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10)
- `department` (optional): Filter by department
- `lecturer_id` (optional): Filter by lecturer ID
- `status` (optional): Filter by status (active/expired)
- `date_from` (optional): Filter events from date (YYYY-MM-DD)
- `date_to` (optional): Filter events until date (YYYY-MM-DD)

**Example Request:**
```bash
curl -H "Authorization: Bearer <token>" \
  "http://localhost:2754/api/admin/events?status=active&department=Computer%20Science"
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "events": [
      {
        "event_id": 1,
        "event_name": "Data Structures Lecture",
        "lecturer_name": "Dr. Jane Smith",
        "lecturer_email": "jane.smith@fupre.edu.ng",
        "department": "Computer Science",
        "event_date": "2025-11-01",
        "event_time": "10:00 AM",
        "venue": "Room 101",
        "expected_students": 25,
        "actual_attendance": 20,
        "attendance_rate": 80.0,
        "status": "expired",
        "created_at": "2025-10-30T00:00:00Z"
      }
    ],
    "pagination": {
      "current_page": 1,
      "per_page": 10,
      "total_pages": 3,
      "total_items": 25
    }
  }
}
```

### Delete Event

**Endpoint:** `DELETE /api/admin/events/:event_id`

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Event deleted successfully. All associated attendance records have been removed."
}
```

**Note:** Deleting an event also removes all associated attendance records (cascade delete).

---

## Analytics

### Get Attendance Trends

Retrieve attendance trends over time with flexible period grouping.

**Endpoint:** `GET /api/admin/trends`

**Query Parameters:**
- `period` (optional): Grouping period - "daily", "weekly", or "monthly" (default: "daily")
- `date_from` (optional): Start date (YYYY-MM-DD)
- `date_to` (optional): End date (YYYY-MM-DD)

**Example Request:**
```bash
curl -H "Authorization: Bearer <token>" \
  "http://localhost:2754/api/admin/trends?period=weekly&date_from=2025-01-01&date_to=2025-12-31"
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "trends": [
      {
        "period": "2025-W01",
        "total_events": 15,
        "total_attendance": 300,
        "average_attendance_rate": 78.5,
        "unique_students": 120
      },
      {
        "period": "2025-W02",
        "total_events": 18,
        "total_attendance": 340,
        "average_attendance_rate": 82.3,
        "unique_students": 125
      }
    ]
  }
}
```

### Get Low Attendance Students

Identify students with attendance rates below a specified threshold.

**Endpoint:** `GET /api/admin/low-attendance`

**Query Parameters:**
- `threshold` (optional): Attendance percentage threshold (default: 75)
- `limit` (optional): Maximum number of results (default: 50)

**Example Request:**
```bash
curl -H "Authorization: Bearer <token>" \
  "http://localhost:2754/api/admin/low-attendance?threshold=60&limit=10"
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "students": [
      {
        "student_id": 5,
        "matric_number": "FUP/CSC/21/5678",
        "first_name": "Alice",
        "last_name": "Johnson",
        "email": "alice.johnson@fupre.edu.ng",
        "department": "Computer Science",
        "total_events": 20,
        "attended": 10,
        "missed": 10,
        "attendance_rate": 50.0
      }
    ],
    "threshold": 60,
    "total_at_risk": 8
  }
}
```

---

## System Settings

### Get System Settings

Retrieve all configurable system settings.

**Endpoint:** `GET /api/admin/settings`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "qr_code_validity_minutes": 30,
    "grace_period_minutes": 15,
    "attendance_threshold_percentage": 75,
    "academic_year": "2024/2025",
    "current_semester": "First",
    "email_verification_required": false,
    "allow_self_registration": true,
    "max_events_per_lecturer": 50
  }
}
```

### Update System Settings

**Endpoint:** `PATCH /api/admin/settings`

**Request Body:**
```json
{
  "qr_code_validity_minutes": 45,
  "attendance_threshold_percentage": 80,
  "current_semester": "Second"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Settings updated successfully",
  "data": {
    "updated_settings": ["qr_code_validity_minutes", "attendance_threshold_percentage", "current_semester"],
    "updated_at": "2025-12-02T03:15:00Z"
  }
}
```

**Available Settings:**
- `qr_code_validity_minutes`: Duration QR codes remain valid (integer)
- `grace_period_minutes`: Late check-in grace period (integer)
- `attendance_threshold_percentage`: Minimum required attendance (integer, 0-100)
- `academic_year`: Current academic year (string, e.g., "2024/2025")
- `current_semester`: Current semester (string, e.g., "First", "Second")
- `email_verification_required`: Require email verification (boolean)
- `allow_self_registration`: Allow new user self-registration (boolean)
- `max_events_per_lecturer`: Maximum events per lecturer (integer)

---

## Audit Logs

### Get Audit Logs

Retrieve a history of all admin actions for compliance and monitoring.

**Endpoint:** `GET /api/admin/audit-logs`

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 50)
- `action` (optional): Filter by action type (login/create/update/delete)
- `user_email` (optional): Filter by admin email
- `date_from` (optional): Start date (YYYY-MM-DD)
- `date_to` (optional): End date (YYYY-MM-DD)

**Example Request:**
```bash
curl -H "Authorization: Bearer <token>" \
  "http://localhost:2754/api/admin/audit-logs?action=delete&limit=10"
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "logs": [
      {
        "log_id": 17,
        "timestamp": "2025-12-02T03:12:19Z",
        "user_type": "admin",
        "user_email": "admin@fupre.edu.ng",
        "action": "delete",
        "resource_type": "event",
        "resource_id": 1,
        "details": "Deleted event with ID 1 and all associated attendance records",
        "ip_address": "172.18.0.1"
      },
      {
        "log_id": 16,
        "timestamp": "2025-12-02T03:12:11Z",
        "user_type": "admin",
        "user_email": "admin@fupre.edu.ng",
        "action": "delete",
        "resource_type": "student",
        "resource_id": 1,
        "details": "Deleted student with ID 1",
        "ip_address": "172.18.0.1"
      }
    ],
    "pagination": {
      "current_page": 1,
      "per_page": 50,
      "total_pages": 1,
      "total_items": 12
    }
  }
}
```

---

## Testing

### Automated Test Suite

A comprehensive test suite is provided to verify all admin endpoints.

**Run the test suite:**
```bash
bash scripts/test-admin-api.sh
```

**Test Coverage:**
1. ✅ Admin login authentication
2. ✅ List students with pagination
3. ✅ List lecturers with pagination
4. ✅ Get student details
5. ✅ Get lecturer details
6. ✅ List all events with filters
7. ✅ Filter expired events
8. ✅ Get attendance trends
9. ✅ Get low attendance students
10. ✅ Get system settings
11. ✅ Update system settings
12. ✅ Get audit logs
13. ✅ Search students by keyword
14. ✅ Filter lecturers by department
15. ✅ Unauthorized access protection

**Expected Output:**
```
======================================
  Test Results
======================================
  Total Tests: 15
  Passed: 15
  Failed: 0

All tests passed! ✓
```

### Manual Testing Examples

**Example 1: Get admin token and list students**
```bash
# Login as admin
TOKEN=$(curl -s http://localhost:2754/api/auth/login-admin \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@fupre.edu.ng","password":"Admin@2024"}' \
  | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)

# List students
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:2754/api/admin/students?page=1&limit=5"
```

**Example 2: Update student status**
```bash
curl -X PATCH http://localhost:2754/api/admin/users/student/1/status \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"active":false,"reason":"Account suspended"}'
```

**Example 3: Get attendance trends**
```bash
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:2754/api/admin/trends?period=weekly&date_from=2025-01-01"
```

**Example 4: Update system settings**
```bash
curl -X PATCH http://localhost:2754/api/admin/settings \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "qr_code_validity_minutes": 45,
    "attendance_threshold_percentage": 80,
    "current_semester": "Second"
  }'
```

---

## Error Handling

All endpoints return consistent error responses:

**401 Unauthorized:**
```json
{
  "error": "invalid or expired token"
}
```

**400 Bad Request:**
```json
{
  "success": false,
  "message": "Invalid user ID"
}
```

**403 Forbidden:**
```json
{
  "error": "insufficient permissions. admin role required"
}
```

**404 Not Found:**
```json
{
  "success": false,
  "message": "User not found"
}
```

**500 Internal Server Error:**
```json
{
  "success": false,
  "message": "Failed to retrieve data"
}
```

---

## Security Considerations

1. **Authentication Required**: All admin endpoints require a valid JWT token
2. **Role-Based Access**: Only users with "admin" role can access these endpoints
3. **Token Expiration**: Admin tokens expire after 7 days
4. **Audit Logging**: All admin actions are logged with IP address and user agent
5. **Soft Deletes**: User and event deletions are soft deletes (recoverable)
6. **Rate Limiting**: Consider implementing rate limiting in production
7. **IP Whitelisting**: Consider restricting admin access to specific IP ranges

---

## Database Schema

### Admin Table
```sql
CREATE TABLE admins (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    department VARCHAR(100),
    role VARCHAR(50) DEFAULT 'admin',
    is_super_admin BOOLEAN DEFAULT false,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

### Audit Log Table
```sql
CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_type VARCHAR(50),
    user_id INTEGER,
    user_email VARCHAR(255),
    action VARCHAR(50),
    resource_type VARCHAR(50),
    resource_id INTEGER,
    details TEXT,
    ip_address VARCHAR(45),
    user_agent TEXT
);
```

### System Settings Table
```sql
CREATE TABLE system_settings (
    id SERIAL PRIMARY KEY,
    setting_key VARCHAR(100) UNIQUE NOT NULL,
    setting_value TEXT,
    data_type VARCHAR(20),
    description TEXT,
    updated_by INTEGER REFERENCES admins(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## Postman Collection

Import the Postman collection from `postman_collection.json` for easy API testing. The collection includes:
- Pre-configured admin authentication
- All admin endpoints with example requests
- Environment variables for base URL and tokens
- Test scripts for response validation

---

## Support

For issues or questions:
- Email: admin@fupre.edu.ng
- Documentation: See `/docs` folder
- API Reference: See `docs/API_REFERENCE.md`

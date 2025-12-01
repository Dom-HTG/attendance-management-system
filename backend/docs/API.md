# API Reference

This file lists the HTTP API endpoints provided by the Attendance Management System. Each endpoint includes method, path, auth requirements, request JSON and example response. Use the root README for quick start and environment details.

Base URL
```
http://localhost:2754
```

1) Student Registration
- Method: POST
- Path: /api/auth/register-student
- Auth: none
- Request JSON:
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@student.edu",
  "password": "securePassword123",
  "matric_number": "STU-2024-001"
}
```
- Success (201):
```json
{ "message": "Student registered successfully", "student_id": 1, "email": "john.doe@student.edu" }
```

2) Lecturer Registration
- Method: POST
- Path: /api/auth/register-lecturer
- Auth: none
- Request JSON:
```json
{
  "first_name": "Jane",
  "last_name": "Smith",
  "email": "jane.smith@lecturer.edu",
  "password": "securePassword123",
  "department": "Computer Science",
  "staff_id": "STAFF-2024-001"
}
```
- Success (201):
```json
{ "message": "Lecturer registered successfully", "lecturer_id": 1, "email": "jane.smith@lecturer.edu" }
```

3) Student Login
- Method: POST
- Path: /api/auth/login-student
- Auth: none
- Request JSON:
```json
{ "email": "john.doe@student.edu", "password": "securePassword123" }
```
- Success (200):
```json
{ "message": "Login successful", "token": "<JWT>", "user_id": 1, "role": "student" }
```

4) Lecturer Login
- Method: POST
- Path: /api/auth/login-lecturer
- Auth: none
- Request JSON:
```json
{ "email": "jane.smith@lecturer.edu", "password": "securePassword123" }
```
- Success (200):
```json
{ "message": "Login successful", "token": "<JWT>", "user_id": 1, "role": "lecturer" }
```

5) Generate QR Code (Lecturer only)
- Method: POST
- Path: /api/lecturer/qrcode/generate
- Auth: Bearer JWT (role=lecturer)
- Request headers: Authorization: Bearer <token>
- Request JSON:
```json
{
  "course_name": "Introduction to Programming",
  "course_code": "CS101",
  "department": "Computer Science",
  "venue": "Room 201",
  "start_time": "2025-11-28T10:00:00Z",
  "end_time": "2025-11-28T11:00:00Z"
}
```
- Success (201): returns event_id, qr_token (UUID), qr_code_data (base64 PNG):
```json
{
  "message": "QR code generated successfully",
  "event_id": 1,
  "qr_token": "550e8400-e29b-41d4-a716-446655440000",
  "qr_code_data": "<base64-png>"
}
```

6) Student Check-In (Scan QR)
- Method: POST
- Path: /api/attendance/check-in
- Auth: Bearer JWT (role=student)
- Request JSON:
```json
{ "qr_token": "550e8400-e29b-41d4-a716-446655440000" }
```
- Success (200):
```json
{ "message": "Check-in successful", "status": "present", "student_id": 1, "marked_time": "2025-11-28T10:15:00Z" }
```

7) Get Student Attendance Records
- Method: GET
- Path: /api/attendance/student/records
- Auth: Bearer JWT (role=student)
- Success (200): returns attendance_records array with marked times and status

8) Get Event Attendance Records (Lecturer)
- Method: GET
- Path: /api/attendance/{event_id}
- Auth: Bearer JWT (role=lecturer)
- Success (200): returns attendance_records with student details for the event

Errors and status codes
- 400 Bad Request: invalid input, missing fields, invalid QR token, event not in correct time window
- 401 Unauthorized: missing or invalid token
- 403 Forbidden: insufficient role permissions
- 404 Not Found: event/user not found
- 409 Conflict: duplicate check-in

Notes
- JWT tokens expire after configured duration (default ~60 minutes). Re-login to obtain a fresh token.
- QR codes are represented as base64-encoded PNG; the important field for check-in is `qr_token` (UUID).
- Times use RFC3339 formatting (e.g., 2025-11-28T10:00:00Z).

For integration examples and sample client snippets, see `../docs/INTEGRATION.md`.

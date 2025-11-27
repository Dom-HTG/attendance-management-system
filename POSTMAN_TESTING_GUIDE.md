# Attendance Management System - Postman Collection Quick Guide

## 📌 Quick Setup

### 1. Import Endpoints

Add these endpoints to your Postman collection:

#### Authentication (Already Exists)
- `POST /api/auth/register-lecturer`
- `POST /api/auth/register-student`
- `POST /api/auth/login-lecturer`
- `POST /api/auth/login-student`

#### QR Code & Attendance (NEW)
- `POST /api/lecturer/qrcode/generate`
- `POST /api/attendance/check-in`
- `GET /api/attendance/:event_id`
- `GET /api/attendance/student/records`

---

## 🧪 Testing Workflow

### Step 1: Register & Login

```
1. Register Lecturer
   POST http://localhost:2754/api/auth/register-lecturer
   Body:
   {
     "first_name": "John",
     "last_name": "Doe",
     "email": "prof_john@university.edu",
     "password": "ProffPassword123",
     "department": "Computer Science",
     "staff_id": "PROF-2024-001"
   }

2. Login Lecturer (Save Token)
   POST http://localhost:2754/api/auth/login-lecturer
   Body:
   {
     "email": "prof_john@university.edu",
     "password": "ProffPassword123"
   }
   
   Response includes access_token
   Save this in Postman variable: {{lecturer_token}}

3. Register Student
   POST http://localhost:2754/api/auth/register-student
   Body:
   {
     "first_name": "Jane",
     "last_name": "Smith",
     "email": "jane_smith@student.edu",
     "password": "StudentPass123",
     "matric_number": "STU-2024-001"
   }

4. Login Student (Save Token)
   POST http://localhost:2754/api/auth/login-student
   Body:
   {
     "email": "jane_smith@student.edu",
     "password": "StudentPass123"
   }
   
   Response includes access_token
   Save this in Postman variable: {{student_token}}
```

---

## 🎓 Test QR Generation

### Request
```
POST http://localhost:2754/api/lecturer/qrcode/generate

Headers:
Authorization: Bearer {{lecturer_token}}
Content-Type: application/json

Body:
{
  "course_name": "Database Systems",
  "course_code": "CS201",
  "start_time": "2025-11-27T14:00:00Z",
  "end_time": "2025-11-27T15:00:00Z",
  "venue": "Computer Lab 301",
  "department": "Computer Science"
}
```

### Response (201 Created)
```json
{
  "message": "QR code generated successfully",
  "event_id": 1,
  "qr_token": "550e8400-e29b-41d4-a716-446655440000",
  "qr_code": "iVBORw0KGgoAAAANSUhEUgAAAgAAAAIACAIAAAB7GkOtAAACh...",
  "course_name": "Database Systems",
  "course_code": "CS201",
  "start_time": "2025-11-27T14:00:00Z",
  "end_time": "2025-11-27T15:00:00Z",
  "venue": "Computer Lab 301",
  "department": "Computer Science",
  "created_by": "John Doe",
  "created_at": "2025-11-27T13:45:00Z",
  "expires_at": "2025-11-27T15:00:00Z"
}
```

**📝 Save `qr_token` in Postman variable: {{qr_token}}**

---

## ✅ Test Student Check-In

### Request
```
POST http://localhost:2754/api/attendance/check-in

Headers:
Authorization: Bearer {{student_token}}
Content-Type: application/json

Body:
{
  "qr_token": "{{qr_token}}"
}
```

### Response (200 OK)
```json
{
  "message": "Check-in successful",
  "status": "present",
  "student_id": 2,
  "student_name": "jane_smith@student.edu",
  "matric_number": "",
  "course_name": "Database Systems (CS201)",
  "course_code": "CS201",
  "marked_time": "2025-11-27T14:15:30Z"
}
```

---

## 📊 View Attendance Records

### Request
```
GET http://localhost:2754/api/attendance/1

Headers:
Authorization: Bearer {{lecturer_token}}
```

### Response (200 OK)
```json
{
  "message": "Attendance records retrieved successfully",
  "event_id": 1,
  "course_name": "Database Systems (CS201)",
  "course_code": "CS201",
  "department": "",
  "start_time": "2025-11-27T14:00:00Z",
  "end_time": "2025-11-27T15:00:00Z",
  "venue": "Computer Lab 301",
  "created_by": "",
  "total_present": 1,
  "attendance_records": [
    {
      "id": 1,
      "student_id": 2,
      "student_name": "Jane Smith",
      "matric_number": "STU-2024-001",
      "status": "present",
      "marked_time": "2025-11-27T14:15:30Z"
    }
  ],
  "generated_at": "2025-11-27T14:20:00Z"
}
```

---

## 📈 Student Attendance History

### Request
```
GET http://localhost:2754/api/attendance/student/records

Headers:
Authorization: Bearer {{student_token}}
```

### Response (200 OK)
```json
{
  "message": "Student attendance records retrieved successfully",
  "student_id": 2,
  "student_name": "",
  "matric_number": "",
  "total_events": 1,
  "total_present": 1,
  "attendance_records": [
    {
      "id": 1,
      "student_id": 2,
      "student_name": "Jane Smith",
      "matric_number": "STU-2024-001",
      "status": "present",
      "marked_time": "2025-11-27T14:15:30Z"
    }
  ],
  "generated_at": "2025-11-27T14:20:00Z"
}
```

---

## 🔄 Complete Test Sequence

1. ✅ Register Lecturer
2. ✅ Login Lecturer (save token)
3. ✅ Register Student
4. ✅ Login Student (save token)
5. ✅ Generate QR Code (save token)
6. ✅ Check-In with QR
7. ✅ View Event Attendance
8. ✅ View Student Attendance History

---

## 🚨 Common Errors & Solutions

### Error 1: `authorization header missing`
**Problem:** No JWT token provided
**Solution:** Add header: `Authorization: Bearer {token}`

### Error 2: `access denied. only lecturers are allowed`
**Problem:** Using student token for lecturer endpoint
**Solution:** Use lecturer token for `/qrcode/generate`

### Error 3: `you have already checked in for this event`
**Problem:** Student checking in twice
**Solution:** This is expected! Use different student or different event

### Error 4: `qr code not found or invalid`
**Problem:** Invalid QR token
**Solution:** Copy exact token from generation response

### Error 5: `event has ended`
**Problem:** Current time is after event end time
**Solution:** Create event with future times, or adjust server time

### Error 6: `invalid token`
**Problem:** Expired or invalid JWT token
**Solution:** Login again to get new token

---

## 📝 Postman Environment Setup

Create a new Postman Environment with these variables:

```json
{
  "name": "Attendance-Test",
  "values": [
    {
      "key": "base_url",
      "value": "http://localhost:2754"
    },
    {
      "key": "lecturer_token",
      "value": ""
    },
    {
      "key": "student_token",
      "value": ""
    },
    {
      "key": "qr_token",
      "value": ""
    },
    {
      "key": "event_id",
      "value": ""
    }
  ]
}
```

---

## 🎯 Request Templates

### Template: Generate QR Code
```
POST {{base_url}}/api/lecturer/qrcode/generate
Headers: Authorization: Bearer {{lecturer_token}}

{
  "course_name": "Course Name",
  "course_code": "CS101",
  "start_time": "2025-11-27T10:00:00Z",
  "end_time": "2025-11-27T11:00:00Z",
  "venue": "Room 101",
  "department": "Computer Science"
}
```

### Template: Check-In
```
POST {{base_url}}/api/attendance/check-in
Headers: Authorization: Bearer {{student_token}}

{
  "qr_token": "{{qr_token}}"
}
```

### Template: Get Event Attendance
```
GET {{base_url}}/api/attendance/{{event_id}}
Headers: Authorization: Bearer {{lecturer_token}}
```

### Template: Get Student History
```
GET {{base_url}}/api/attendance/student/records
Headers: Authorization: Bearer {{student_token}}
```

---

## ✨ Tips & Tricks

### Tip 1: Automatic Token Extraction
Add a test script to Postman's login request:
```javascript
var jsonData = pm.response.json();
pm.environment.set("lecturer_token", jsonData.access_token);
```

### Tip 2: Timestamp Management
For testing with current time:
- Use: `new Date().toISOString()` in JavaScript
- Or set times to be in the future

### Tip 3: Multiple Students
Create multiple students with different matriculation numbers to test bulk attendance

### Tip 4: Organize Collection
```
Attendance-Management
├── Auth
│   ├── Register Lecturer
│   ├── Login Lecturer
│   ├── Register Student
│   └── Login Student
├── QR Code
│   └── Generate QR Code
├── Attendance
│   ├── Check-In
│   ├── Get Event Attendance
│   └── Get Student Records
```

---

## 🔍 Debugging Tips

### Check Request
- Verify all required fields in body
- Confirm Authorization header format
- Check base URL is correct

### Check Response
- Look at status code (200, 201, 400, 401, etc.)
- Read error message for details
- Check response data for actual errors

### Check Server
```bash
# View server logs
docker-compose logs app

# Check database
docker-compose exec postgres psql -U postgres -d attendance-management
```

---

## 📊 Data Verification Queries

To verify data in PostgreSQL:

```sql
-- Check events created
SELECT id, event_name, qr_code_token, start_time, end_time FROM events;

-- Check attendance records
SELECT ua.id, ua.student_id, ua.status, ua.marked_time
FROM user_attendances ua
JOIN attendances a ON ua.attendance_id = a.id;

-- Check specific event attendance
SELECT s.first_name, s.last_name, ua.status, ua.marked_time
FROM user_attendances ua
JOIN students s ON ua.student_id = s.id
WHERE ua.attendance_id = 1
ORDER BY ua.marked_time;

-- Check student history
SELECT e.event_name, ua.status, ua.marked_time
FROM user_attendances ua
JOIN attendances a ON ua.attendance_id = a.id
JOIN events e ON a.event_id = e.id
WHERE ua.student_id = 1
ORDER BY ua.marked_time DESC;
```

---

## 🚀 Performance Testing

To test with multiple students:

1. Create 10+ students (Register Student endpoint)
2. Generate 1 QR code
3. Check-in all students simultaneously (use Postman Collection Runner)
4. View results: `GET /api/attendance/1`

Expected: All check-ins succeed within 1 second

---

## 🎓 Learning Path

1. **Basic**: Generate QR → Check-In → View Records
2. **Intermediate**: Test error cases (duplicate, invalid token, expired event)
3. **Advanced**: Multiple students, concurrent check-ins, performance testing

---

**Happy Testing! 🚀**

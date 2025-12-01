# Postman Collection - Updated with Analytics & Nigerian Data

## Overview

The Postman collection has been completely updated to include all 25+ analytics endpoints and uses Nigerian names with proper matric number formatting.

## What's New

### 1. **Nigerian Names & Matric Numbers**

All dummy data now uses Nigerian names and matric format `COS/7XXX/XXX`:

- **Student 1**: Chioma Okafor - `COS/7452/234`
- **Student 2**: Adeyemi Oluwaseun - `COS/7381/156`
- **Student 3**: Folake Adebayo - `COS/7629/487`
- **Lecturer**: Dr. Adekunle Afolabi

### 2. **Complete Collection Structure**

The collection now has **9 main sections** with 45+ individual endpoints:

#### Section 1: Authentication (6 endpoints)
- 1.1 Register Student (Nigerian)
- 1.2 Register Student 2 (Nigerian)
- 1.3 Register Student 3 (Nigerian)
- 1.4 Register Lecturer (Nigerian)
- 1.5 Student Login (auto-saves token)
- 1.6 Lecturer Login (auto-saves token)

#### Section 2: QR Code (1 endpoint)
- 2.1 Generate QR Code for Event (auto-saves event_id and qr_token)

#### Section 3: Attendance - Student (2 endpoints)
- 3.1 Check-In (Scan QR Code)
- 3.2 Get Student Attendance Records

#### Section 4: Attendance - Lecturer (1 endpoint)
- 4.1 Get Event Attendance Records

#### Section 5: Analytics - Student (2 endpoints)
- 5.1 Get Student Metrics
- 5.2 Get Student Insights

#### Section 6: Analytics - Lecturer (3 endpoints)
- 6.1 Get Lecturer Course Metrics
- 6.2 Get Course Performance Details
- 6.3 Get Lecturer Insights

#### Section 7: Analytics - Admin (3 endpoints)
- 7.1 Get Admin Overview
- 7.2 Get Department Metrics
- 7.3 Get Real-Time Dashboard

#### Section 8: Analytics - Temporal & Advanced (7 endpoints)
- 8.1 Get Temporal Analytics
- 8.2 Detect Anomalies
- 8.3 Predict Student Attendance
- 8.4 Predict Course Attendance
- 8.5 Get Benchmark Comparison
- 8.6 Get Chart Data - Line Trend
- 8.7 Get Chart Data - Bar Comparison

#### Section 9: Error Tests (8 endpoints)
- E1. Missing Required Fields
- E2. Wrong Password on Login
- E3. Access Protected Route Without Token
- E4. Student Cannot Generate QR Code
- E5. Invalid QR Token on Check-In
- E6. Access Analytics Without Token
- E7. Student Access Other Student Analytics
- E8. Student Cannot Access Admin Dashboard

## How to Use

### Step 1: Import Collection
1. Open Postman
2. Click **Import** → **Upload File**
3. Select `postman_collection.json`
4. Click **Import**

### Step 2: Setup Environment Variables
The collection uses these variables (auto-populated on successful login):
- `base_url`: `http://localhost:2754`
- `student_token`: Auto-saved from login
- `student_id`: Auto-saved from login
- `lecturer_token`: Auto-saved from login
- `lecturer_id`: Auto-saved from login
- `event_id`: Auto-saved from QR generation
- `qr_token`: Auto-saved from QR generation
- `course_code`: `CS101` (default)
- `department`: `Computer Science` (default)

### Step 3: Run Test Flow

**Recommended order for testing:**

1. **Register Students** (1.1, 1.2, 1.3)
   - Creates 3 student accounts
   
2. **Register Lecturer** (1.4)
   - Creates lecturer account
   
3. **Student Login** (1.5)
   - Saves student_token and student_id
   
4. **Lecturer Login** (1.6)
   - Saves lecturer_token and lecturer_id
   
5. **Generate QR Code** (2.1)
   - Saves event_id and qr_token
   
6. **Check-In** (3.1)
   - Student uses QR token to sign in
   
7. **Get Attendance Records** (3.2, 4.1)
   - View attendance data
   
8. **Analytics Endpoints** (5.1 → 8.7)
   - Test all analytics features
   - Student endpoints: 5.1, 5.2, 8.1-8.7
   - Lecturer endpoints: 6.1-6.3, 8.3-8.7
   - Admin endpoints: 7.1-7.3, 8.4, 8.5-8.7

9. **Error Tests** (9.1 → 9.8)
   - Verify error handling and authorization

## Key Features

### Auto-Token Saving
Successful login responses automatically save tokens using test scripts:
```javascript
if (pm.response.code === 200) {
    var jsonData = pm.response.json();
    pm.environment.set('student_token', jsonData.data.access_token);
    pm.environment.set('student_id', jsonData.data.user_id);
}
```

### Auto-Event ID Saving
QR code generation automatically saves event details:
```javascript
if (pm.response.code === 201) {
    var jsonData = pm.response.json();
    pm.environment.set('event_id', jsonData.data.event_id);
    pm.environment.set('qr_token', jsonData.data.qr_token);
}
```

### Variable Substitution
All requests use Postman variables:
- `{{base_url}}` - API base URL
- `{{student_token}}` - Student authentication token
- `{{lecturer_token}}` - Lecturer authentication token
- `{{student_id}}` - Student ID for analytics
- `{{event_id}}` - Event ID for attendance
- `{{qr_token}}` - QR token for check-in
- `{{course_code}}` - Course code for queries
- `{{department}}` - Department name for queries

## Dummy Data Reference

### Students
| Name | Email | Matric | Password |
|------|-------|--------|----------|
| Chioma Okafor | chioma.okafor@student.edu | COS/7452/234 | SecurePass123! |
| Adeyemi Oluwaseun | adeyemi.oluwaseun@student.edu | COS/7381/156 | SecurePass123! |
| Folake Adebayo | folake.adebayo@student.edu | COS/7629/487 | SecurePass123! |

### Lecturer
| Name | Email | Staff ID | Password |
|------|-------|----------|----------|
| Dr. Adekunle Afolabi | adekunle.afolabi@lecturer.edu | CS-STAFF-001 | SecurePass123! |

## Testing Analytics

### Student Analytics Flow
1. Login as student (1.5)
2. Get student metrics (5.1)
3. Get student insights (5.2)
4. Get predictions (8.3)
5. Get benchmark comparison (8.5)
6. Get chart data (8.6)

### Lecturer Analytics Flow
1. Login as lecturer (1.6)
2. Generate QR code (2.1)
3. Get course metrics (6.1)
4. Get course performance (6.2)
5. Get lecturer insights (6.3)
6. Get course predictions (8.4)

### Admin Analytics Flow
1. Login as lecturer (1.6) - Has admin privileges
2. Get admin overview (7.1)
3. Get department metrics (7.2)
4. Get real-time dashboard (7.3)
5. Get temporal analytics (8.1)
6. Detect anomalies (8.2)

## Error Testing

All error tests demonstrate proper authorization and validation:
- **E1**: Missing fields validation
- **E2**: Invalid credentials
- **E3**: Missing authentication token
- **E4**: Role-based access control (RBAC)
- **E5**: Invalid input validation
- **E6**: Protected endpoint without token
- **E7**: Student privacy enforcement
- **E8**: Role hierarchy enforcement

## Notes

- All timestamps are in UTC format (ISO 8601)
- All analytics endpoints require Bearer token authentication
- Students can only view their own analytics
- Lecturers can view their courses and student analytics
- Admins (lecturer role) can view university-wide analytics
- Timestamps use: `2025-12-01T10:00:00Z` format
- Department: `Computer Science`
- Course Code: `CS101`

## Next Steps

1. Start Docker containers: `docker compose up --build`
2. Import this collection in Postman
3. Run login endpoints to populate variables
4. Execute analytics endpoints in order
5. Check responses for proper data structure
6. Run error tests to verify security

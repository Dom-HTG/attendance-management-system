# Quick Testing Guide - Postman Collection

## 🚀 Quick Start (5 minutes)

### Prerequisites
- Docker running: `docker compose up --build`
- Postman installed
- Collection imported

### Test Flow
```
1. Run 1.1 → 1.4  (Register all users)
2. Run 1.5 & 1.6  (Login → auto-saves tokens)
3. Run 2.1        (Generate QR → auto-saves event_id & qr_token)
4. Run 3.1        (Student check-in with QR)
5. Run 5.1        (Get Student Metrics)
6. Run 8.1 - 8.7  (Test all analytics)
```

## 📋 Test Checklist

### Authentication ✓
- [ ] 1.1 Student Register - Chioma Okafor (COS/7452/234)
- [ ] 1.2 Student Register - Adeyemi (COS/7381/156)
- [ ] 1.3 Student Register - Folake (COS/7629/487)
- [ ] 1.4 Lecturer Register - Dr. Adekunle
- [ ] 1.5 Student Login - Token auto-saved
- [ ] 1.6 Lecturer Login - Token auto-saved

### QR & Attendance ✓
- [ ] 2.1 Generate QR Code - event_id auto-saved
- [ ] 3.1 Student Check-In - Uses {{qr_token}}
- [ ] 3.2 Get Attendance Records (Student)
- [ ] 4.1 Get Attendance Records (Lecturer)

### Student Analytics ✓
- [ ] 5.1 Get Student Metrics - Shows rate, trends, engagement
- [ ] 5.2 Get Student Insights - AI recommendations

### Lecturer Analytics ✓
- [ ] 6.1 Get Course Metrics - All courses
- [ ] 6.2 Get Course Performance - CS101 details
- [ ] 6.3 Get Lecturer Insights - Performance insights

### Admin Analytics ✓
- [ ] 7.1 Get Admin Overview - University-wide stats
- [ ] 7.2 Get Department Metrics - Computer Science
- [ ] 7.3 Get Real-Time Dashboard - Active sessions

### Advanced Analytics ✓
- [ ] 8.1 Temporal Analytics - Day-of-week patterns
- [ ] 8.2 Detect Anomalies - Duplicate check-ins
- [ ] 8.3 Student Predictions - Attendance forecast
- [ ] 8.4 Course Predictions - Course forecast
- [ ] 8.5 Benchmark Comparison - Peer ranking
- [ ] 8.6 Chart Data (Line) - Trend visualization
- [ ] 8.7 Chart Data (Bar) - Comparison data

### Security & Errors ✓
- [ ] E1 Missing Fields - 400 Bad Request
- [ ] E2 Wrong Password - 401 Unauthorized
- [ ] E3 No Token - 401 Unauthorized
- [ ] E4 Student QR Gen - 403 Forbidden
- [ ] E5 Invalid Token - 400 Bad Request
- [ ] E6 Analytics No Token - 401 Unauthorized
- [ ] E7 Student Other ID - 403 Forbidden
- [ ] E8 Student Admin Access - 403 Forbidden

## 🔑 Key Variables

After running logins, these will be auto-populated:
```
student_token = JWT token for student
student_id = Auto-saved student ID
lecturer_token = JWT token for lecturer
lecturer_id = Auto-saved lecturer ID
event_id = Auto-saved from QR generation
qr_token = Auto-saved QR token
```

## 📊 Expected Responses

### Success (200/201)
```json
{
  "success": true,
  "message": "...",
  "data": { ... }
}
```

### Error (400/401/403)
```json
{
  "success": false,
  "error_message": "...",
  "error": "..."
}
```

## 🧪 Test Scenarios

### Scenario 1: Full Student Flow
1. Register & login (1.1 → 1.5)
2. Scan QR code (2.1 → 3.1)
3. View metrics (5.1 → 5.2)
4. Get predictions (8.3)

### Scenario 2: Lecturer Analytics
1. Register & login lecturer (1.4 → 1.6)
2. Generate QR (2.1)
3. View course metrics (6.1 → 6.2)
4. Get insights (6.3)

### Scenario 3: Admin Dashboard
1. Login as lecturer (1.6) - Has admin privileges
2. Overview (7.1)
3. Department (7.2)
4. Real-time (7.3)

### Scenario 4: Security Testing
1. Run all E* tests (9.1 → 9.8)
2. Verify proper error codes
3. Confirm authorization blocks

## 🆘 Troubleshooting

### "Invalid token" error
- Re-run 1.5 or 1.6 to get fresh token
- Token expires after 1 hour

### "Student not found" on analytics
- Make sure you ran 1.1 first
- Check {{student_id}} variable is set

### "Event not found" on check-in
- Run 2.1 first to generate QR
- Check {{event_id}} is populated

### Variables not auto-saving
- Make sure response is 200/201
- Check Postman test scripts are enabled
- Run login again

## 💡 Tips

- **Run in order**: Follow 1.1 → 9.8 for best results
- **Check variables**: Click "Environment" → Look at variable values
- **Read descriptions**: Each endpoint has explanation
- **Use raw preview**: View actual JSON response
- **Export results**: Use "Runner" to test entire flow

## 📝 Nigerian Matric Format

Pattern: `COS/7XXX/XXX` (Computer Science / Year 2024 / Student Number)

Examples used:
- `COS/7452/234` - Chioma
- `COS/7381/156` - Adeyemi  
- `COS/7629/487` - Folake

## ⏱️ Estimated Times

| Task | Time |
|------|------|
| Register all users | 1 min |
| Login & generate QR | 1 min |
| Student analytics (5.1-5.2) | 1 min |
| All analytics endpoints | 3 min |
| Error tests | 2 min |
| **Total** | **~8 minutes** |

## 🎯 Success Criteria

✅ All endpoints return 2xx for valid requests
✅ All error tests return expected 4xx codes
✅ Variables auto-populate after login
✅ Analytics data contains Nigerian names
✅ Matric numbers in format COS/7XXX/XXX
✅ Timestamps in ISO 8601 format
✅ All responses contain proper JSON structure

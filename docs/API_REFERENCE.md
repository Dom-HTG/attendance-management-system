# API Reference

Base URL: `http://localhost:2754`

## Authentication

**Register Student**
```
POST /api/auth/register-student
Body: { first_name, last_name, email, password, matric_number }
Response: { success, message, data }
```

**Register Lecturer**
```
POST /api/auth/register-lecturer
Body: { first_name, last_name, email, password, department, staff_id }
Response: { success, message, data }
```

**Student Login**
```
POST /api/auth/login-student
Body: { email, password }
Response: { success, message, data: { access_token, user } }
```

**Lecturer Login**
```
POST /api/auth/login-lecturer
Body: { email, password }
Response: { success, message, data: { access_token, user } }
```

## Attendance

**Generate QR Code** (Lecturer)
```
POST /api/lecturer/qrcode/generate
Auth: Bearer <token>
Body: { course_name, course_code, start_time, end_time, venue, department }
Response: { message, event_id, qr_token, qr_code, course_name, created_by, expires_at }
```

**Check In** (Student)
```
POST /api/attendance/check-in
Auth: Bearer <token>
Body: { qr_token }
Response: { message, status, student_id, student_name, matric_number, course_name, course_code, marked_time }
```

**Get Event Attendance** (Lecturer)
```
GET /api/attendance/{event_id}
Auth: Bearer <token>
Response: { message, event_id, course_name, department, attendance_records[], total_present }
```

**Get Student Records** (Student)
```
GET /api/attendance/student/records
Auth: Bearer <token>
Response: { message, student_id, student_name, matric_number, attendance_records[], total_events }
```

## Analytics - Student

**Student Metrics**
```
GET /api/analytics/student/{student_id}
Auth: Bearer <token>
Response: { overall_attendance_rate, total_sessions, per_course_rates[], attendance_trend[], engagement_score }
```

**Student Insights**
```
GET /api/analytics/student/{student_id}/insights
Auth: Bearer <token>
Response: { insights[], recommendations[], strengths[], areas_for_improvement[] }
```

## Analytics - Lecturer

**Course Metrics**
```
GET /api/analytics/lecturer/courses
Auth: Bearer <token> (lecturer)
Response: { courses: [{ course_code, average_attendance_rate, total_sessions, at_risk_students }] }
```

**Course Performance**
```
GET /api/analytics/lecturer/course/{course_code}
Auth: Bearer <token> (lecturer)
Response: { course_code, attendance_summary, session_breakdown[], top_performers[], at_risk_students[] }
```

**Lecturer Insights**
```
GET /api/analytics/lecturer/insights
Auth: Bearer <token> (lecturer)
Response: { overall_summary, course_comparisons[], recommendations[], trends[] }
```

## Analytics - Admin

**Admin Overview**
```
GET /api/analytics/admin/overview
Auth: Bearer <token> (lecturer/admin)
Response: { university_wide_metrics, department_summaries[], recent_trends, at_risk_summary }
```

**Department Metrics**
```
GET /api/analytics/admin/department/{department}
Auth: Bearer <token> (lecturer/admin)
Response: { department, metrics, course_breakdown[], lecturer_performance[], student_engagement }
```

**Real-Time Dashboard**
```
GET /api/analytics/admin/realtime
Auth: Bearer <token> (lecturer/admin)
Response: { active_sessions[], recent_checkins[], current_attendance_rate, system_health }
```

## Analytics - Advanced

**Temporal Analytics**
```
GET /api/analytics/temporal?timeframe={week|month|semester}&group_by={day|week|month}
Auth: Bearer <token>
Response: { timeframe, attendance_patterns[], peak_times[], day_of_week_analysis[] }
```

**Anomaly Detection**
```
GET /api/analytics/anomalies?type={duplicate|suspicious|pattern}&timeframe={week|month}
Auth: Bearer <token> (lecturer/admin)
Response: { anomalies: [{ type, severity, description, timestamp, affected_entities }] }
```

**Predict Student Attendance**
```
GET /api/analytics/predictions/student/{student_id}?weeks=4
Auth: Bearer <token>
Response: { student_id, predictions: [{ week, predicted_rate, confidence, factors[] }] }
```

**Predict Course Attendance**
```
GET /api/analytics/predictions/course/{course_code}?weeks=4
Auth: Bearer <token> (lecturer)
Response: { course_code, predictions: [{ week, predicted_rate, confidence, trends[] }] }
```

**Benchmark Comparison**
```
GET /api/analytics/benchmark?entity_type={student|course}&entity_id={id}&comparison_group={department|university}
Auth: Bearer <token>
Response: { entity_info, benchmark_metrics, percentile_ranking, comparison_data[] }
```

**Chart Data**
```
GET /api/analytics/charts/{chart_type}?params...
Auth: Bearer <token>
chart_type: line|bar|pie|heatmap
Response: { chart_type, data_points[], labels[], metadata }
```

## Status Codes

- **200** OK
- **201** Created
- **400** Bad Request (invalid input)
- **401** Unauthorized (missing/invalid token)
- **403** Forbidden (insufficient permissions)
- **404** Not Found
- **409** Conflict (duplicate entry)
- **500** Internal Server Error

## Notes

- All timestamps use RFC3339 format (e.g., `2025-12-01T10:00:00Z`)
- JWT tokens expire after 60 minutes
- QR codes valid within event time window
- Role-based access: student, lecturer, admin

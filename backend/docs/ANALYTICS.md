# Analytics API Documentation

This document provides detailed reference for all analytics endpoints and their usage.

## Overview

The analytics system provides actionable insights for students, lecturers, and administrators across multiple dimensions:

- **Student Analytics**: Individual performance, trends, engagement scores, at-risk detection
- **Lecturer Analytics**: Course performance, attendance trends, class insights
- **Admin Analytics**: University-wide metrics, department comparisons, real-time dashboards
- **Temporal Analytics**: Time-based patterns, day-of-week analysis, seasonal trends
- **Predictive Analytics**: Forecasted attendance, risk indicators
- **Anomaly Detection**: Unusual patterns, suspected fraud, duplicate check-ins
- **Benchmarking**: Peer comparisons, percentile ranking
- **Visualization**: Chart data for frontend rendering

All analytics endpoints require authentication (Bearer token).

---

## 1. Student Analytics Endpoints

### Get Student Metrics
**Endpoint**: `GET /api/analytics/student/{student_id}`

**Authorization**: 
- Students can view only their own metrics
- Lecturers and admins can view any student's metrics

**Response**:
```json
{
  "success": true,
  "message": "Student metrics retrieved successfully",
  "data": {
    "student_id": 1,
    "student_name": "John Doe",
    "matric_number": "MAT001",
    "overall_attendance_rate": 85.5,
    "total_sessions": 20,
    "total_present": 17,
    "total_absent": 3,
    "total_late": 2,
    "attendance_streak": 5,
    "late_checkin_frequency": 2,
    "class_average_comparison": 5.5,
    "at_risk_status": false,
    "engagement_score": 82.3,
    "per_course_rates": [
      {
        "course_code": "CS101",
        "course_name": "Computer Networks",
        "attendance_rate": 90.0,
        "sessions_attended": 9,
        "total_sessions": 10,
        "department": "CSE"
      }
    ],
    "attendance_trend": [
      {
        "period": "2025-45",
        "attendance_rate": 80.0,
        "sessions_attended": 4,
        "total_sessions": 5,
        "average_checkin_time_minutes": 3
      }
    ],
    "generated_at": "2025-11-29T15:30:00Z"
  }
}
```

### Get Student Insights
**Endpoint**: `GET /api/analytics/student/{student_id}/insights`

**Response**:
```json
{
  "success": true,
  "message": "Student insights generated successfully",
  "data": {
    "entity_type": "student",
    "entity_id": 1,
    "entity_name": "John Doe",
    "summary": "John Doe has excellent attendance at 85.5%. Keep up the consistency and engagement with courses.",
    "key_takeaways": [
      "Excellent attendance record",
      "High engagement score"
    ],
    "trends": [
      {
        "trend": "Attendance improving",
        "explanation": "Recent weeks show improvement in attendance rates",
        "timeframe": "past 4 weeks"
      }
    ],
    "recommendations": [],
    "generated_at": "2025-11-29T15:30:00Z"
  }
}
```

---

## 2. Lecturer Analytics Endpoints

### Get Lecturer Course Metrics
**Endpoint**: `GET /api/analytics/lecturer/courses`

**Authorization**: Lecturer role required

**Response**:
```json
{
  "success": true,
  "message": "Lecturer course metrics retrieved successfully",
  "data": {
    "lecturer_id": 5,
    "lecturer_name": "Dr. Ahmed Hassan",
    "department": "Computer Science",
    "total_courses": 3,
    "average_attendance": 78.5,
    "qr_generated_count": 45,
    "course_metrics": [
      {
        "course_code": "CS101",
        "course_name": "Computer Networks",
        "attendance_average": 82.0,
        "session_count": 15,
        "student_count": 45,
        "most_attended_session": {
          "event_id": 101,
          "event_name": "CS101 Lecture 1",
          "start_time": "2025-11-01T10:00:00Z",
          "attendance_rate": 95.0,
          "students_present": 43,
          "total_enrolled": 45
        },
        "least_attended_session": {
          "event_id": 110,
          "event_name": "CS101 Lecture 10",
          "start_time": "2025-11-20T10:00:00Z",
          "attendance_rate": 60.0,
          "students_present": 27,
          "total_enrolled": 45
        },
        "average_checkin_time_minutes": 2
      }
    ],
    "generated_at": "2025-11-29T15:30:00Z"
  }
}
```

### Get Course Performance
**Endpoint**: `GET /api/analytics/lecturer/course/{course_code}`

**Authorization**: Lecturer role required

**Response**:
```json
{
  "success": true,
  "message": "Course performance retrieved successfully",
  "data": {
    "course_code": "CS101",
    "course_name": "Computer Networks",
    "lecturer_name": "Dr. Ahmed Hassan",
    "department": "Computer Science",
    "student_count": 45,
    "overall_attendance_rate": 82.0,
    "attendance_distribution": {
      "range_0_to_20": 0,
      "range_20_to_40": 2,
      "range_40_to_60": 5,
      "range_60_to_80": 18,
      "range_80_to_100": 20
    },
    "students_at_risk": 7,
    "average_checkin_time_minutes": 2,
    "late_arrivals_count": 12,
    "session_duration_vs_attendance": [
      {
        "duration_minutes": 60,
        "attendance_rate": 85.0,
        "session_count": 10
      }
    ],
    "generated_at": "2025-11-29T15:30:00Z"
  }
}
```

### Get Lecturer Insights
**Endpoint**: `GET /api/analytics/lecturer/insights`

**Authorization**: Lecturer role required

**Response**: Similar structure to student insights with lecturer-specific recommendations

---

## 3. Admin Analytics Endpoints

### Get Admin Overview
**Endpoint**: `GET /api/analytics/admin/overview`

**Authorization**: Lecturer/Admin role required

**Response**:
```json
{
  "success": true,
  "message": "Admin overview retrieved successfully",
  "data": {
    "overall_attendance_rate": 78.5,
    "total_active_sessions": 12,
    "total_students": 500,
    "total_lecturers": 30,
    "department_comparison": [
      {
        "department_name": "Computer Science",
        "attendance_rate": 82.0,
        "student_count": 150,
        "lecturer_count": 10,
        "course_count": 25,
        "average_checkin_time_minutes": 2
      }
    ],
    "top_performing_courses": [...],
    "lowest_performing_courses": [...],
    "lecturer_performance": [
      {
        "lecturer_id": 5,
        "lecturer_name": "Dr. Ahmed Hassan",
        "department": "Computer Science",
        "average_class_attendance": 82.0,
        "courses_managed": 3,
        "qr_sessions_created": 45,
        "efficiency_score": 88.5
      }
    ],
    "generated_at": "2025-11-29T15:30:00Z"
  }
}
```

### Get Department Metrics
**Endpoint**: `GET /api/analytics/admin/department/{department}`

**Authorization**: Lecturer/Admin role required

**Query Parameters**: None

**Response**:
```json
{
  "success": true,
  "message": "Department metrics retrieved successfully",
  "data": {
    "department_name": "Computer Science",
    "overall_attendance_rate": 82.0,
    "student_count": 150,
    "lecturer_count": 10,
    "course_count": 25,
    "attendance_trend": [...],
    "course_enrollment_vs_attendance": [
      {
        "course_code": "CS101",
        "course_name": "Computer Networks",
        "enrolled": 45,
        "actual_attended": 37,
        "attendance_rate": 82.2
      }
    ],
    "lecturer_efficiency": [...],
    "student_engagement_by_year": [...],
    "venue_utilization": [...]
  }
}
```

### Get Real-Time Dashboard
**Endpoint**: `GET /api/analytics/admin/realtime`

**Authorization**: Lecturer/Admin role required

**Response**:
```json
{
  "success": true,
  "message": "Real-time dashboard retrieved successfully",
  "data": {
    "active_sessions_now": 5,
    "total_checkins_today": 234,
    "average_attendance_today": 81.5,
    "ongoing_sessions": [
      {
        "event_id": 101,
        "course_name": "CS101 - Lecture",
        "lecturer": "Dr. Ahmed Hassan",
        "venue": "Room 201",
        "start_time": "2025-11-29T14:00:00Z",
        "checkins_count": 40,
        "students_enrolled": 45,
        "attendance_rate": 88.9
      }
    ],
    "system_usage_stats": {
      "total_api_calls_today": 1250,
      "qr_codes_generated_today": 15,
      "checkins_processed_today": 234,
      "active_users_today": 150
    },
    "generated_at": "2025-11-29T15:30:00Z"
  }
}
```

---

## 4. Temporal Analytics

### Get Temporal Analytics
**Endpoint**: `GET /api/analytics/temporal?start_date={start}&end_date={end}&granularity={granularity}`

**Query Parameters**:
- `start_date` (required): RFC3339 format, e.g., `2025-11-01T00:00:00Z`
- `end_date` (required): RFC3339 format, e.g., `2025-11-30T23:59:59Z`
- `granularity` (optional): `daily`, `weekly`, or `monthly` (default: `weekly`)

**Response**:
```json
{
  "success": true,
  "message": "Temporal analytics retrieved successfully",
  "data": {
    "granularity": "weekly",
    "start_date": "2025-11-01T00:00:00Z",
    "end_date": "2025-11-30T23:59:59Z",
    "attendance_heatmap": [
      {
        "day_of_week": "Monday",
        "time_slot": "10:00",
        "attendance_rate": 85.5,
        "session_count": 12,
        "avg_checkin_time_minutes": 2
      }
    ],
    "day_of_week_analysis": [
      {
        "day_of_week": "Monday",
        "attendance_rate": 85.0,
        "session_count": 30,
        "average_present": 27,
        "time_slots": [...]
      }
    ],
    "seasonal_trends": [...],
    "holiday_impact": [...]
  }
}
```

---

## 5. Anomaly Detection

### Detect Anomalies
**Endpoint**: `GET /api/analytics/anomalies`

**Authorization**: Any authenticated user

**Response**:
```json
{
  "success": true,
  "message": "Anomaly detection completed",
  "data": {
    "anomaly_count": 3,
    "critical_anomalies": 1,
    "anomalies": [
      {
        "id": 1,
        "type": "duplicate_checkin",
        "severity": "high",
        "description": "Multiple check-ins detected for this event",
        "student_id": 42,
        "student_name": "Jane Smith",
        "event_id": 150,
        "course_name": "CS101",
        "detection_time": "2025-11-29T15:20:00Z",
        "recommended_action": "Review for possible QR code sharing or technical glitch"
      }
    ],
    "generated_at": "2025-11-29T15:30:00Z"
  }
}
```

---

## 6. Predictive Analytics

### Predict Student Attendance
**Endpoint**: `GET /api/analytics/predictions/student/{student_id}`

**Response**:
```json
{
  "success": true,
  "message": "Attendance prediction generated successfully",
  "data": {
    "entity_type": "student",
    "entity_id": 1,
    "entity_name": "John Doe",
    "forecasted_attendance": 85.0,
    "current_attendance": 85.0,
    "confidence_level": 65.0,
    "risk_factors": [],
    "recommended_actions": [
      "Continue current attendance pattern for best outcomes"
    ],
    "generated_at": "2025-11-29T15:30:00Z"
  }
}
```

### Predict Course Attendance
**Endpoint**: `GET /api/analytics/predictions/course/{course_code}`

**Response**: Similar structure with course-level predictions

---

## 7. Benchmarking

### Get Benchmark Comparison
**Endpoint**: `GET /api/analytics/benchmark?entity_type={type}&entity_id={id}`

**Query Parameters**:
- `entity_type` (required): `student`, `course`, or `department`
- `entity_id` (required): ID of the entity

**Response**:
```json
{
  "success": true,
  "message": "Benchmark comparison retrieved successfully",
  "data": {
    "entity_type": "student",
    "entity_id": 1,
    "entity_name": "John Doe",
    "performance_value": 85.5,
    "peer_average": 78.0,
    "peer_std_dev": 8.5,
    "percentile_rank": 109.6,
    "performance_vs_peers": "above",
    "historical_comparison": {
      "current_semester": 85.5,
      "previous_semester": 82.0,
      "change_percent": 4.27,
      "trend_direction": "up"
    },
    "goal_tracking": {
      "target_attendance": 80.0,
      "actual_attendance": 85.5,
      "goal_met_status": "on_track",
      "days_remaining_in_period": 45
    }
  }
}
```

---

## 8. Visualization / Chart Data

### Get Chart Data
**Endpoint**: `GET /api/analytics/charts/{chart_type}?entity_type={type}&entity_id={id}`

**Path Parameters**:
- `chart_type`: `line_trend`, `bar_comparison`, `pie_distribution`, `heatmap`, `scatter_correlation`

**Query Parameters**:
- `entity_type` (required): `student`, `course`, or `department`
- `entity_id` (required): ID of the entity

**Response Examples**:

#### Line Chart (Trend)
```json
{
  "success": true,
  "message": "Chart data retrieved successfully",
  "data": {
    "chart_type": "line_trend",
    "title": "Attendance Trend",
    "description": "Attendance rate over time",
    "labels": ["2025-43", "2025-44", "2025-45", "2025-46"],
    "data_points": {
      "datasets": [
        {
          "label": "Attendance Rate",
          "data": [80.0, 82.0, 85.5, 84.0],
          "color": "#2ecc71"
        }
      ]
    }
  }
}
```

#### Bar Chart (Comparison)
```json
{
  "success": true,
  "message": "Chart data retrieved successfully",
  "data": {
    "chart_type": "bar_comparison",
    "title": "Course Comparison",
    "description": "Attendance rates across courses",
    "labels": ["CS101", "CS102", "CS201", "CS202"],
    "data_points": {
      "labels": ["CS101", "CS102", "CS201", "CS202"],
      "datasets": [
        {
          "label": "Attendance %",
          "data": [90.0, 85.5, 78.0, 72.5],
          "color": "#3498db"
        }
      ]
    }
  }
}
```

---

## Performance Targets

- **Single-entity queries** (student/course metrics): <500ms
- **Bulk queries** (admin overview, department metrics): <2s
- **Real-time dashboard**: Updates every 30 seconds
- **Prediction generation**: <1s

---

## Error Responses

All analytics endpoints use standard HTTP status codes:

- `200 OK`: Request successful
- `400 Bad Request`: Invalid parameters
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Insufficient permissions
- `500 Internal Server Error`: Server-side error

Error response format:
```json
{
  "success": false,
  "error_message": "Failed to retrieve student metrics",
  "error": "student not found"
}
```

---

## Authentication

All endpoints require a Bearer token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

Obtain a token by logging in:
```bash
POST /api/auth/login-student
POST /api/auth/login-lecturer
```

---

## Rate Limiting & Caching

Currently, no rate limiting is enforced. Frequently accessed metrics are cached and invalidated on new check-ins.

Future enhancements:
- Redis-based distributed caching
- Rate limiting per user/role
- Batch API for multiple entity queries

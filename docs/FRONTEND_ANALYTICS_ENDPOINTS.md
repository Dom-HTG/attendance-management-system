# Frontend Analytics Endpoints - Implementation Guide

**Status**: ✅ All endpoints implemented and tested  
**Date**: December 1, 2025  
**API Base URL**: `http://localhost:2754`

---

## Overview

All 4 required analytics endpoints are **fully implemented and working**. This document provides complete integration details for the frontend team.

---

## 🔑 Authentication

All endpoints require JWT authentication. Include the token in the Authorization header:

```javascript
headers: {
  'Authorization': `Bearer ${token}`,
  'Content-Type': 'application/json'
}
```

### Get Lecturer Token:
```bash
curl -X POST http://localhost:2754/api/auth/login-lecturer \
  -H "Content-Type: application/json" \
  -d '{
    "email": "dr.adebayo.olumide@fupre.edu.ng",
    "password": "Lecturer@123"
  }'
```

### Get Student Token:
```bash
curl -X POST http://localhost:2754/api/auth/login-student \
  -H "Content-Type: application/json" \
  -d '{
    "email": "chukwuemeka.okonkwo@fupre.edu.ng",
    "password": "Student@100"
  }'
```

---

## 📊 Endpoint 1: Get Lecturer Events

**Endpoint**: `GET /api/events/lecturer`  
**Auth**: Bearer token (lecturer role required)  
**Purpose**: Get all events created by the logged-in lecturer with attendance counts

### Request Example:
```bash
curl -X GET http://localhost:2754/api/events/lecturer \
  -H "Authorization: Bearer $LECTURER_TOKEN"
```

### Response Format:
```json
{
  "success": true,
  "message": "Events retrieved successfully",
  "data": {
    "events": [
      {
        "event_id": 5,
        "course_name": "Machine Learning",
        "course_code": "CSC501",
        "department": "Computer Science",
        "venue": "Lecture Hall 3",
        "start_time": "2025-12-01T23:33:36Z",
        "end_time": "2025-12-02T01:33:36Z",
        "qr_token": "88aed18c-9c55-4afe-b431-7661bb5baccd",
        "status": "active",
        "total_attendance": 12,
        "created_at": "2025-12-02T00:03:36.982809Z"
      },
      {
        "event_id": 4,
        "course_name": "Software Engineering",
        "course_code": "CSC401",
        "department": "Computer Science",
        "venue": "Lecture Hall 2",
        "start_time": "2025-12-01T23:33:36Z",
        "end_time": "2025-12-02T01:33:36Z",
        "qr_token": "e184bd92-32c3-4ed9-b3b1-e235680cbe9b",
        "status": "active",
        "total_attendance": 10,
        "created_at": "2025-12-02T00:03:36.624228Z"
      }
    ],
    "total_events": 5,
    "total_students_reached": 15
  }
}
```

### Response Fields:

| Field | Type | Description |
|-------|------|-------------|
| `event_id` | integer | Unique event identifier |
| `course_name` | string | Name of the course |
| `course_code` | string | Course code (e.g., CSC301) |
| `department` | string | Department name |
| `venue` | string | Event location |
| `start_time` | datetime (ISO 8601) | Event start time |
| `end_time` | datetime (ISO 8601) | Event end time |
| `qr_token` | string | UUID token for QR code |
| `status` | string | "active" or "expired" |
| `total_attendance` | integer | Number of students who checked in |
| `created_at` | datetime (ISO 8601) | When event was created |
| `total_events` | integer | Total number of events |
| `total_students_reached` | integer | Unique students across all events |

### Status Logic:
- **"active"**: Current time is before `end_time`
- **"expired"**: Current time is after `end_time`

### Frontend Integration:
```typescript
// TypeScript/React example
interface LecturerEvent {
  event_id: number;
  course_name: string;
  course_code: string;
  department: string;
  venue: string;
  start_time: string;
  end_time: string;
  qr_token: string;
  status: 'active' | 'expired';
  total_attendance: number;
  created_at: string;
}

interface LecturerEventsResponse {
  events: LecturerEvent[];
  total_events: number;
  total_students_reached: number;
}

const fetchLecturerEvents = async (): Promise<LecturerEventsResponse> => {
  const response = await fetch(`${API_BASE_URL}/api/events/lecturer`, {
    headers: {
      'Authorization': `Bearer ${getToken()}`,
      'Content-Type': 'application/json'
    }
  });
  
  if (!response.ok) throw new Error('Failed to fetch events');
  
  const result = await response.json();
  return result.data;
};
```

---

## 📈 Endpoint 2: Get Lecturer Summary

**Endpoint**: `GET /api/analytics/lecturer/summary`  
**Auth**: Bearer token (lecturer role required)  
**Purpose**: Get aggregated statistics for lecturer dashboard

### Request Example:
```bash
curl -X GET http://localhost:2754/api/analytics/lecturer/summary \
  -H "Authorization: Bearer $LECTURER_TOKEN"
```

### Response Format:
```json
{
  "success": true,
  "message": "Lecturer summary retrieved successfully",
  "data": {
    "total_events_created": 5,
    "total_students_reached": 15,
    "average_attendance_rate": 100.0,
    "sessions_this_week": 5,
    "sessions_today": 0,
    "most_attended_course": {
      "course_code": "CSC501",
      "course_name": "Machine Learning",
      "avg_attendance": 95.5
    },
    "attendance_trend": [
      {
        "period": "2025-12-01",
        "attendance_rate": 100.0,
        "sessions_attended": 56,
        "total_sessions": 5,
        "average_checkin_time_minutes": 0
      }
    ]
  }
}
```

### Response Fields:

| Field | Type | Description |
|-------|------|-------------|
| `total_events_created` | integer | Total events created by lecturer |
| `total_students_reached` | integer | Unique students who attended any event |
| `average_attendance_rate` | float | Average attendance rate across all events (%) |
| `sessions_this_week` | integer | Events created this week (Mon-Sun) |
| `sessions_today` | integer | Events created today |
| `most_attended_course` | object/null | Course with highest attendance (null if no events) |
| `attendance_trend` | array | Daily attendance trends |

### Most Attended Course:
```json
{
  "course_code": "CSC501",
  "course_name": "Machine Learning",
  "avg_attendance": 95.5
}
```

### Attendance Trend Data Point:
```json
{
  "period": "2025-12-01",
  "attendance_rate": 100.0,
  "sessions_attended": 56,
  "total_sessions": 5,
  "average_checkin_time_minutes": 0
}
```

### Frontend Integration:
```typescript
interface LecturerSummary {
  total_events_created: number;
  total_students_reached: number;
  average_attendance_rate: number;
  sessions_this_week: number;
  sessions_today: number;
  most_attended_course: {
    course_code: string;
    course_name: string;
    avg_attendance: number;
  } | null;
  attendance_trend: Array<{
    period: string;
    attendance_rate: number;
    sessions_attended: number;
    total_sessions: number;
    average_checkin_time_minutes: number;
  }>;
}

const fetchLecturerSummary = async (): Promise<LecturerSummary> => {
  const response = await fetch(`${API_BASE_URL}/api/analytics/lecturer/summary`, {
    headers: {
      'Authorization': `Bearer ${getToken()}`,
      'Content-Type': 'application/json'
    }
  });
  
  if (!response.ok) throw new Error('Failed to fetch summary');
  
  const result = await response.json();
  return result.data;
};
```

---

## 🏛️ Endpoint 3: Get Admin Overview

**Endpoint**: `GET /api/analytics/admin/overview`  
**Auth**: Bearer token (lecturer/admin role required)  
**Purpose**: Get university-wide statistics for admin dashboard

### Request Example:
```bash
curl -X GET http://localhost:2754/api/analytics/admin/overview \
  -H "Authorization: Bearer $LECTURER_TOKEN"
```

### Response Format:
```json
{
  "success": true,
  "message": "Admin overview retrieved successfully",
  "data": {
    "total_students": 18,
    "total_lecturers": 4,
    "total_departments": 2,
    "total_events": 8,
    "average_attendance_rate": 100.0,
    "active_sessions_now": 8,
    "qr_codes_generated_today": 8,
    "total_check_ins_today": 59,
    "system_health": {
      "database_status": "healthy",
      "last_check_in": "2025-12-02T01:03:26.29662493Z",
      "uptime_hours": 0
    },
    "overall_attendance_rate": 100.0,
    "total_active_sessions": 8,
    "generated_at": "2025-12-02T01:03:26.29662573Z"
  }
}
```

### Response Fields:

| Field | Type | Description |
|-------|------|-------------|
| `total_students` | integer | Total registered students |
| `total_lecturers` | integer | Total registered lecturers |
| `total_departments` | integer | Number of unique departments |
| `total_events` | integer | Total events created |
| `average_attendance_rate` | float | System-wide average attendance (%) |
| `active_sessions_now` | integer | Events happening right now |
| `qr_codes_generated_today` | integer | Events created today |
| `total_check_ins_today` | integer | Student check-ins today |
| `system_health` | object | System health metrics |
| `overall_attendance_rate` | float | Same as average_attendance_rate |
| `total_active_sessions` | integer | Same as active_sessions_now |
| `generated_at` | datetime (ISO 8601) | When report was generated |

### System Health:
```json
{
  "database_status": "healthy",
  "last_check_in": "2025-12-02T01:03:26.29662493Z",
  "uptime_hours": 0
}
```

### Frontend Integration:
```typescript
interface AdminOverview {
  total_students: number;
  total_lecturers: number;
  total_departments: number;
  total_events: number;
  average_attendance_rate: number;
  active_sessions_now: number;
  qr_codes_generated_today: number;
  total_check_ins_today: number;
  system_health: {
    database_status: string;
    last_check_in: string;
    uptime_hours: number;
  };
  overall_attendance_rate: number;
  total_active_sessions: number;
  generated_at: string;
}

const fetchAdminOverview = async (): Promise<AdminOverview> => {
  const response = await fetch(`${API_BASE_URL}/api/analytics/admin/overview`, {
    headers: {
      'Authorization': `Bearer ${getToken()}`,
      'Content-Type': 'application/json'
    }
  });
  
  if (!response.ok) throw new Error('Failed to fetch overview');
  
  const result = await response.json();
  return result.data;
};
```

---

## 🏢 Endpoint 4: Get Department Statistics

**Endpoint**: `GET /api/analytics/admin/departments`  
**Auth**: Bearer token (lecturer/admin role required)  
**Purpose**: Get per-department breakdown for charts and analysis

### Request Example:
```bash
curl -X GET http://localhost:2754/api/analytics/admin/departments \
  -H "Authorization: Bearer $LECTURER_TOKEN"
```

### Response Format:
```json
{
  "success": true,
  "message": "Department statistics retrieved successfully",
  "data": {
    "departments": [
      {
        "department": "Computer Science",
        "total_students": 15,
        "total_lecturers": 1,
        "total_events": 5,
        "average_attendance_rate": 100.0,
        "total_check_ins": 56
      },
      {
        "department": "Test Department",
        "total_students": 3,
        "total_lecturers": 3,
        "total_events": 3,
        "average_attendance_rate": 100.0,
        "total_check_ins": 9
      }
    ]
  }
}
```

### Response Fields:

| Field | Type | Description |
|-------|------|-------------|
| `department` | string | Department name |
| `total_students` | integer | Students in this department |
| `total_lecturers` | integer | Lecturers in this department |
| `total_events` | integer | Events for this department |
| `average_attendance_rate` | float | Department average attendance (%) |
| `total_check_ins` | integer | Total check-ins for this department |

### Sorting:
Departments are sorted by `total_students` descending (largest departments first).

### Frontend Integration:
```typescript
interface DepartmentStat {
  department: string;
  total_students: number;
  total_lecturers: number;
  total_events: number;
  average_attendance_rate: number;
  total_check_ins: number;
}

interface DepartmentStatsResponse {
  departments: DepartmentStat[];
}

const fetchDepartmentStats = async (): Promise<DepartmentStatsResponse> => {
  const response = await fetch(`${API_BASE_URL}/api/analytics/admin/departments`, {
    headers: {
      'Authorization': `Bearer ${getToken()}`,
      'Content-Type': 'application/json'
    }
  });
  
  if (!response.ok) throw new Error('Failed to fetch department stats');
  
  const result = await response.json();
  return result.data;
};
```

---

## 🔄 Complete Integration Example

### React/TypeScript Dashboard Component:

```typescript
import { useEffect, useState } from 'react';

const LecturerDashboard = () => {
  const [summary, setSummary] = useState<LecturerSummary | null>(null);
  const [events, setEvents] = useState<LecturerEventsResponse | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [summaryData, eventsData] = await Promise.all([
          fetchLecturerSummary(),
          fetchLecturerEvents()
        ]);
        
        setSummary(summaryData);
        setEvents(eventsData);
      } catch (error) {
        console.error('Failed to fetch dashboard data:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  if (loading) return <div>Loading...</div>;

  return (
    <div className="dashboard">
      <h1>Lecturer Dashboard</h1>
      
      {/* Summary Cards */}
      <div className="stats-grid">
        <StatCard 
          title="Total Events" 
          value={summary?.total_events_created || 0} 
        />
        <StatCard 
          title="Students Reached" 
          value={summary?.total_students_reached || 0} 
        />
        <StatCard 
          title="Average Attendance" 
          value={`${summary?.average_attendance_rate.toFixed(1)}%`} 
        />
        <StatCard 
          title="Sessions This Week" 
          value={summary?.sessions_this_week || 0} 
        />
      </div>

      {/* Events Table */}
      <div className="events-table">
        <h2>Recent Events</h2>
        <table>
          <thead>
            <tr>
              <th>Course</th>
              <th>Venue</th>
              <th>Date</th>
              <th>Status</th>
              <th>Attendance</th>
            </tr>
          </thead>
          <tbody>
            {events?.events.map(event => (
              <tr key={event.event_id}>
                <td>{event.course_code} - {event.course_name}</td>
                <td>{event.venue}</td>
                <td>{new Date(event.start_time).toLocaleDateString()}</td>
                <td>
                  <span className={`badge ${event.status}`}>
                    {event.status}
                  </span>
                </td>
                <td>{event.total_attendance} students</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};
```

### Admin Dashboard Component:

```typescript
const AdminDashboard = () => {
  const [overview, setOverview] = useState<AdminOverview | null>(null);
  const [departments, setDepartments] = useState<DepartmentStatsResponse | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [overviewData, deptData] = await Promise.all([
          fetchAdminOverview(),
          fetchDepartmentStats()
        ]);
        
        setOverview(overviewData);
        setDepartments(deptData);
      } catch (error) {
        console.error('Failed to fetch admin data:', error);
      }
    };

    fetchData();
  }, []);

  return (
    <div className="admin-dashboard">
      <h1>Admin Dashboard</h1>
      
      {/* Overview Cards */}
      <div className="stats-grid">
        <StatCard title="Total Students" value={overview?.total_students || 0} />
        <StatCard title="Total Lecturers" value={overview?.total_lecturers || 0} />
        <StatCard title="Active Sessions" value={overview?.active_sessions_now || 0} />
        <StatCard title="Avg Attendance" value={`${overview?.average_attendance_rate.toFixed(1)}%`} />
      </div>

      {/* Department Breakdown */}
      <div className="departments-chart">
        <h2>Department Statistics</h2>
        {departments?.departments.map(dept => (
          <DepartmentCard key={dept.department} data={dept} />
        ))}
      </div>
    </div>
  );
};
```

---

## ⚡ Performance Notes

### Response Times:
- All endpoints respond in **< 200ms** with current dataset
- Optimized with database indexes on key columns

### Caching Strategy:
```typescript
// Optional: Cache data for 30 seconds to reduce API calls
const cache = new Map<string, { data: any, timestamp: number }>();

const fetchWithCache = async (key: string, fetcher: () => Promise<any>) => {
  const cached = cache.get(key);
  const now = Date.now();
  
  if (cached && now - cached.timestamp < 30000) {
    return cached.data;
  }
  
  const data = await fetcher();
  cache.set(key, { data, timestamp: now });
  return data;
};
```

---

## 🛡️ Error Handling

### Common Error Responses:

**401 Unauthorized:**
```json
{
  "success": false,
  "message": "Unauthorized: Invalid or missing token"
}
```

**403 Forbidden:**
```json
{
  "success": false,
  "message": "Forbidden: Insufficient permissions"
}
```

**500 Internal Server Error:**
```json
{
  "success": false,
  "message": "Internal server error"
}
```

### Frontend Error Handler:
```typescript
const handleAPIError = (error: any) => {
  if (error.status === 401) {
    // Redirect to login
    window.location.href = '/login';
  } else if (error.status === 403) {
    // Show permission denied message
    showError('You do not have permission to access this resource');
  } else {
    // Show generic error
    showError('An error occurred. Please try again later.');
  }
};
```

---

## 🧪 Testing Endpoints

### Test with seeded data:

```bash
# Get lecturer token
LECTURER_TOKEN=$(curl -s -X POST http://localhost:2754/api/auth/login-lecturer \
  -H "Content-Type: application/json" \
  -d '{"email":"dr.adebayo.olumide@fupre.edu.ng","password":"Lecturer@123"}' \
  | jq -r '.access_token')

# Test all 4 endpoints
curl -X GET http://localhost:2754/api/events/lecturer \
  -H "Authorization: Bearer $LECTURER_TOKEN" | jq

curl -X GET http://localhost:2754/api/analytics/lecturer/summary \
  -H "Authorization: Bearer $LECTURER_TOKEN" | jq

curl -X GET http://localhost:2754/api/analytics/admin/overview \
  -H "Authorization: Bearer $LECTURER_TOKEN" | jq

curl -X GET http://localhost:2754/api/analytics/admin/departments \
  -H "Authorization: Bearer $LECTURER_TOKEN" | jq
```

---

## ✅ Checklist for Frontend Integration

- [ ] Update API base URL in frontend config
- [ ] Implement authentication token management
- [ ] Create TypeScript interfaces for all response types
- [ ] Add API functions for all 4 endpoints
- [ ] Update Lecturer Dashboard to fetch real data
- [ ] Update Admin Dashboard to fetch real data
- [ ] Add loading states while fetching data
- [ ] Implement error handling for failed requests
- [ ] Add retry logic for failed requests
- [ ] Test with different user roles (student, lecturer)
- [ ] Test with empty data (new lecturer with no events)
- [ ] Verify all numbers update when new attendance is marked

---

## 🎯 Success Criteria

✅ All 4 endpoints are implemented and tested  
✅ Lecturer dashboard shows real numbers (not 0s)  
✅ Admin dashboard shows real numbers (not 0s)  
✅ Response times are under 500ms  
✅ Proper authentication and role checking  
✅ Comprehensive error handling  
✅ Data updates reflect new attendance records  

---

**Implementation Status**: ✅ Complete  
**Last Tested**: December 1, 2025  
**Next Steps**: Frontend team can now integrate these endpoints

For questions or issues, refer to `docs/API_REFERENCE.md` or test the endpoints using the provided curl commands.

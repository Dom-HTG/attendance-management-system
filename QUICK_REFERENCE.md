# Analytics Endpoints - Quick Reference

**API Base**: `http://localhost:2754`  
**Status**: ‚úÖ All endpoints working

---

## Ì¥ê Get Auth Token

```bash
# Lecturer
curl -X POST http://localhost:2754/api/auth/login-lecturer \
  -H "Content-Type: application/json" \
  -d '{"email":"dr.adebayo.olumide@fupre.edu.ng","password":"Lecturer@123"}'
```

---

## Ì≥ä Lecturer Endpoints

### 1. Get All Events
```bash
GET /api/events/lecturer
Authorization: Bearer {token}
```
**Returns**: List of events with attendance counts

### 2. Get Dashboard Summary
```bash
GET /api/analytics/lecturer/summary
Authorization: Bearer {token}
```
**Returns**: Total events, students reached, avg attendance, sessions this week

---

## ÌøõÔ∏è Admin Endpoints

### 3. Get System Overview
```bash
GET /api/analytics/admin/overview
Authorization: Bearer {token}
```
**Returns**: Total students, lecturers, events, avg attendance, active sessions

### 4. Get Department Stats
```bash
GET /api/analytics/admin/departments
Authorization: Bearer {token}
```
**Returns**: Per-department breakdown with attendance rates

---

## ÌæØ TypeScript Interfaces

```typescript
// Lecturer Events
interface LecturerEventsResponse {
  events: Array<{
    event_id: number;
    course_name: string;
    course_code: string;
    venue: string;
    status: 'active' | 'expired';
    total_attendance: number;
  }>;
  total_events: number;
  total_students_reached: number;
}

// Lecturer Summary
interface LecturerSummary {
  total_events_created: number;
  total_students_reached: number;
  average_attendance_rate: number;
  sessions_this_week: number;
  sessions_today: number;
}

// Admin Overview
interface AdminOverview {
  total_students: number;
  total_lecturers: number;
  total_events: number;
  average_attendance_rate: number;
  active_sessions_now: number;
  qr_codes_generated_today: number;
  total_check_ins_today: number;
}

// Department Stats
interface DepartmentStats {
  departments: Array<{
    department: string;
    total_students: number;
    total_lecturers: number;
    total_events: number;
    average_attendance_rate: number;
    total_check_ins: number;
  }>;
}
```

---

## Ì∫Ä Quick Integration

```typescript
const API_BASE = 'http://localhost:2754';

// Fetch lecturer events
const events = await fetch(`${API_BASE}/api/events/lecturer`, {
  headers: { 'Authorization': `Bearer ${token}` }
}).then(r => r.json());

// Fetch lecturer summary
const summary = await fetch(`${API_BASE}/api/analytics/lecturer/summary`, {
  headers: { 'Authorization': `Bearer ${token}` }
}).then(r => r.json());

// Fetch admin overview
const overview = await fetch(`${API_BASE}/api/analytics/admin/overview`, {
  headers: { 'Authorization': `Bearer ${token}` }
}).then(r => r.json());

// Fetch department stats
const departments = await fetch(`${API_BASE}/api/analytics/admin/departments`, {
  headers: { 'Authorization': `Bearer ${token}` }
}).then(r => r.json());
```

---

## Ì≥ñ Full Documentation

- **Complete Guide**: `docs/FRONTEND_ANALYTICS_ENDPOINTS.md`
- **Status Report**: `ANALYTICS_STATUS.md`
- **API Reference**: `docs/API_REFERENCE.md`

---

## ‚úÖ Test Credentials

**Lecturer**:
- Email: `dr.adebayo.olumide@fupre.edu.ng`
- Password: `Lecturer@123`

**Students**: See `seed-login-credentials.txt` (15 students available)

---

**Status**: Ìø¢ Production Ready | **Response Times**: < 200ms

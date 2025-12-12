# QR Code Generation - Frontend Guide

## Overview
This guide explains everything the frontend needs to know about QR code generation for attendance tracking in the Attendance Management System.

---

## **Endpoint**
```
POST /api/lecturer/qrcode/generate
```

## **Authentication**
- Requires JWT Bearer token (Lecturer role only)
- Header: `Authorization: Bearer <lecturer_token>`

---

## **Request Body**
```json
{
  "course_code": "CSC301",
  "course_name": "Data Structures & Algorithms",
  "department": "Computer Science",
  "venue": "Lecture Hall 3",
  "start_time": "2025-12-02T10:00:00Z",
  "end_time": "2025-12-02T12:00:00Z"
}
```

### **Field Requirements:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `course_code` | string | Yes | Course identifier (e.g., "CSC301") |
| `course_name` | string | Yes | Full course name |
| `department` | string | Yes | Department name |
| `venue` | string | Yes | Location of the event |
| `start_time` | string | Yes | ISO 8601 datetime format (UTC) |
| `end_time` | string | Yes | ISO 8601 datetime format (UTC), must be after start_time |

---

## **Success Response (201 Created)**
```json
{
  "success": true,
  "message": "QR code generated successfully",
  "data": {
    "event_id": 5,
    "qr_code_base64": "iVBORw0KGgoAAAANSUhEUgAAAZAAAAGQCAIAAAAP3aGbAAAG...",
    "qr_token": "QR_1733108691_8a7b9c3d",
    "course_code": "CSC301",
    "course_name": "Data Structures & Algorithms",
    "venue": "Lecture Hall 3",
    "start_time": "2025-12-02T10:00:00Z",
    "end_time": "2025-12-02T12:00:00Z",
    "status": "active"
  }
}
```

### **Response Fields:**
- `event_id`: Database ID for the event (use for tracking/analytics)
- `qr_code_base64`: PNG image encoded in base64 (256x256 pixels)
- `qr_token`: The actual token students will scan (format: `QR_{timestamp}_{hash}`)
- `course_code`, `course_name`, `venue`: Echo of request data
- `start_time`, `end_time`: Validity window for the QR code
- `status`: Current status - "active" or "expired"

---

## **Frontend Display**

### **1. Display QR Code Image**
```typescript
// The qr_code_base64 is a PNG image encoded in base64
const qrCodeSrc = `data:image/png;base64,${response.data.qr_code_base64}`;

// Display in React/Vue/Angular
<img src={qrCodeSrc} alt="Event QR Code" className="qr-code" />
```

### **2. Important Data to Store**
```typescript
interface EventData {
  event_id: number;        // For tracking and analytics
  qr_token: string;        // For manual entry fallback
  qr_code_base64: string;  // For display
  course_code: string;
  course_name: string;
  venue: string;
  start_time: string;      // Display validity window
  end_time: string;        // Display validity window
  status: "active" | "expired";
}
```

### **3. Display Event Details**
```typescript
// Format for UI display
const eventDetails = {
  "Course": `${response.data.course_code} - ${response.data.course_name}`,
  "Venue": response.data.venue,
  "Valid From": formatDateTime(response.data.start_time),
  "Valid Until": formatDateTime(response.data.end_time),
  "Status": response.data.status,
  "Token": response.data.qr_token  // For manual entry
};
```

---

## **QR Code Token Format**
- **Format**: `QR_{timestamp}_{random_hash}`
- **Example**: `QR_1733108691_8a7b9c3d`
- **Purpose**: Students scan this token to check in
- **Uniqueness**: Each event has a unique QR token

---

## **Validity Rules**
1. ✅ QR code is **active** only between `start_time` and `end_time`
2. ⏰ After `end_time`, status automatically becomes "expired"
3. 🚫 Students can only check in during the active window
4. 🔑 Each event has a unique `qr_token` that cannot be reused

---

## **Student Check-In Process**

### **Endpoint for Students**
```
POST /api/attendance/check-in
Authorization: Bearer <student_token>
```

### **Request Body**
```json
{
  "qr_token": "QR_1733108691_8a7b9c3d"
}
```

### **Success Response**
```json
{
  "success": true,
  "message": "Attendance marked successfully",
  "data": {
    "attendance_id": 123,
    "student_name": "John Doe",
    "event_name": "Data Structures & Algorithms",
    "checked_in_at": "2025-12-02T10:05:00Z"
  }
}
```

---

## **Error Responses**

### **400 - Validation Error**
```json
{
  "success": false,
  "error_message": "Invalid request data",
  "details": {
    "course_code": "this field is required"
  }
}
```

### **400 - Invalid Time Range**
```json
{
  "success": false,
  "error_message": "end_time must be after start_time"
}
```

### **401 - Unauthorized**
```json
{
  "error": "invalid or expired token"
}
```

### **403 - Wrong Role**
```json
{
  "error": "insufficient permissions. lecturer role required"
}
```

### **404 - Event Not Found** (for check-in)
```json
{
  "success": false,
  "message": "Invalid or expired QR token"
}
```

### **409 - Already Checked In** (for check-in)
```json
{
  "success": false,
  "message": "You have already checked in for this event"
}
```

---

## **Frontend Implementation Example**

### **TypeScript/React Example**
```typescript
import axios from 'axios';

interface QRCodeRequest {
  course_code: string;
  course_name: string;
  department: string;
  venue: string;
  start_time: string;
  end_time: string;
}

interface QRCodeResponse {
  success: boolean;
  message: string;
  data: {
    event_id: number;
    qr_code_base64: string;
    qr_token: string;
    course_code: string;
    course_name: string;
    venue: string;
    start_time: string;
    end_time: string;
    status: string;
  };
}

async function generateQRCode(token: string): Promise<QRCodeResponse> {
  const payload: QRCodeRequest = {
    course_code: "CSC301",
    course_name: "Data Structures",
    department: "Computer Science",
    venue: "Lecture Theatre 3",
    start_time: new Date(Date.now() + 5 * 60000).toISOString(),  // 5 min from now
    end_time: new Date(Date.now() + 125 * 60000).toISOString()   // 2 hours 5 min
  };

  try {
    const response = await axios.post<QRCodeResponse>(
      'http://localhost:2754/api/lecturer/qrcode/generate',
      payload,
      {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      }
    );

    return response.data;
  } catch (error) {
    if (axios.isAxiosError(error)) {
      console.error('QR Generation failed:', error.response?.data);
      throw new Error(error.response?.data?.error_message || 'Failed to generate QR code');
    }
    throw error;
  }
}

// Usage in React Component
function QRCodeGenerator() {
  const [qrData, setQRData] = useState<QRCodeResponse['data'] | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleGenerate = async () => {
    setLoading(true);
    setError(null);
    
    try {
      const token = localStorage.getItem('lecturer_token');
      if (!token) throw new Error('Not authenticated');
      
      const response = await generateQRCode(token);
      setQRData(response.data);
      
      // Optional: Start countdown timer
      startCountdown(response.data.end_time);
      
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <button onClick={handleGenerate} disabled={loading}>
        {loading ? 'Generating...' : 'Generate QR Code'}
      </button>
      
      {error && <div className="error">{error}</div>}
      
      {qrData && (
        <div className="qr-display">
          <img 
            src={`data:image/png;base64,${qrData.qr_code_base64}`}
            alt="Event QR Code"
            className="qr-code-image"
          />
          
          <div className="event-details">
            <h3>{qrData.course_code} - {qrData.course_name}</h3>
            <p>Venue: {qrData.venue}</p>
            <p>Valid: {new Date(qrData.start_time).toLocaleString()} - 
                      {new Date(qrData.end_time).toLocaleString()}</p>
            <p>Status: <span className={qrData.status}>{qrData.status}</span></p>
            <p>Token: <code>{qrData.qr_token}</code></p>
          </div>
        </div>
      )}
    </div>
  );
}
```

---

## **Best Practices**

### **1. Time Validation**
```typescript
// Validate times on frontend before submission
function validateTimes(startTime: Date, endTime: Date): boolean {
  if (endTime <= startTime) {
    alert('End time must be after start time');
    return false;
  }
  
  if (startTime < new Date()) {
    alert('Start time cannot be in the past');
    return false;
  }
  
  return true;
}
```

### **2. Timezone Handling**
```typescript
// Always use UTC timestamps (ISO 8601 format)
const startTime = new Date().toISOString(); // "2025-12-02T10:00:00.000Z"

// For display to users, convert to local time
const displayTime = new Date(startTime).toLocaleString();
```

### **3. Countdown Timer**
```typescript
function startCountdown(endTime: string) {
  const interval = setInterval(() => {
    const now = new Date().getTime();
    const end = new Date(endTime).getTime();
    const distance = end - now;
    
    if (distance < 0) {
      clearInterval(interval);
      updateStatus('expired');
      return;
    }
    
    const hours = Math.floor(distance / (1000 * 60 * 60));
    const minutes = Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60));
    const seconds = Math.floor((distance % (1000 * 60)) / 1000);
    
    updateTimer(`${hours}h ${minutes}m ${seconds}s`);
  }, 1000);
}
```

### **4. Status Indicator**
```css
/* Visual indicator for QR code status */
.status.active {
  color: #22c55e;
  font-weight: bold;
}

.status.expired {
  color: #ef4444;
  font-weight: bold;
}
```

### **5. Download QR Code**
```typescript
function downloadQRCode(base64Data: string, filename: string) {
  const link = document.createElement('a');
  link.href = `data:image/png;base64,${base64Data}`;
  link.download = `${filename}.png`;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
}
```

### **6. Manual Token Entry**
```typescript
// Provide manual token input for students who can't scan
<div className="manual-entry">
  <p>Can't scan? Enter token manually:</p>
  <code className="token-display">{qrData.qr_token}</code>
  <button onClick={() => navigator.clipboard.writeText(qrData.qr_token)}>
    Copy Token
  </button>
</div>
```

### **7. Refresh Option**
```typescript
// Allow lecturer to view/manage all their events
async function fetchLecturerEvents(token: string) {
  const response = await axios.get(
    'http://localhost:2754/api/events/lecturer',
    {
      headers: { 'Authorization': `Bearer ${token}` }
    }
  );
  return response.data;
}
```

---

## **Additional Related Endpoints**

### **Get Lecturer's Events**
```
GET /api/events/lecturer
Authorization: Bearer <lecturer_token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "events": [
      {
        "event_id": 1,
        "course_code": "CSC301",
        "course_name": "Data Structures",
        "venue": "LT3",
        "event_date": "2025-12-02",
        "event_time": "10:00 AM",
        "status": "active",
        "qr_token": "QR_1733108691_8a7b9c3d",
        "total_attendance": 45,
        "created_at": "2025-12-01T08:00:00Z"
      }
    ]
  }
}
```

### **Get Event Attendance**
```
GET /api/attendance/{event_id}
Authorization: Bearer <lecturer_token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "event": {
      "event_id": 1,
      "course_name": "Data Structures",
      "venue": "LT3",
      "date": "2025-12-02"
    },
    "attendance": [
      {
        "student_id": 10,
        "matric_number": "FUPRE/2021/10000",
        "first_name": "Chukwuemeka",
        "last_name": "Okonkwo",
        "checked_in_at": "2025-12-02T10:05:00Z"
      }
    ],
    "total_attendance": 45
  }
}
```

---

## **QR Code Specifications**

- **Format**: PNG image
- **Size**: 256x256 pixels
- **Encoding**: Base64
- **Content**: QR token in format `QR_{timestamp}_{hash}`
- **Error Correction**: Medium level (M)
- **Ready to use**: No additional processing needed

---

## **Testing Credentials**

From `seed-login-credentials.txt`:

**Lecturer:**
- Email: `dr.adebayo.olumide@fupre.edu.ng`
- Password: `Lecturer@123`

**Test Student:**
- Email: `chukwuemeka.okonkwo@fupre.edu.ng`
- Password: `Student@100`

**API Base URL:**
```
http://localhost:2754
```

---

## **Common Issues & Solutions**

### **Issue: QR Code Not Displaying**
```typescript
// Ensure you're using the correct data URI format
const correctFormat = `data:image/png;base64,${base64String}`;
// NOT: `data:image/png,${base64String}` (missing base64)
```

### **Issue: Token Expired Immediately**
```typescript
// Check your time ranges - start_time should be in the future
const startTime = new Date(Date.now() + 5 * 60000).toISOString(); // 5 minutes from now
```

### **Issue: CORS Error**
```typescript
// Backend should have CORS enabled
// If testing locally, ensure both frontend and backend are running
```

### **Issue: 403 Forbidden**
```typescript
// Ensure you're using a LECTURER token, not a student token
// Check token format: "Bearer <token>" not just "<token>"
```

---

## **Summary Checklist**

Frontend needs to handle:
- ✅ JWT authentication with lecturer token
- ✅ Form validation (all fields required, end_time > start_time)
- ✅ Display QR code as base64 PNG image
- ✅ Show event details (course, venue, times)
- ✅ Display countdown timer to expiration
- ✅ Status indicator (active/expired)
- ✅ Provide manual token for students without QR scanner
- ✅ Download QR code option
- ✅ Error handling for all API responses
- ✅ Refresh/fetch lecturer's existing events

---

**For more information, see:**
- Main API Documentation: `docs/API_REFERENCE.md`
- Admin Dashboard Guide: `docs/ADMIN_GUIDE.md`
- Quick Start Guide: `docs/QUICK_START.md`

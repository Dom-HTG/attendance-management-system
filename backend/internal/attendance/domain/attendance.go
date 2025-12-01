package attendance

// Request DTOs

// GenerateQRCodeDTO represents the request to generate a QR code for an event.
type GenerateQRCodeDTO struct {
	CourseName string `json:"course_name" binding:"required"`
	CourseCode string `json:"course_code" binding:"required"`
	StartTime  string `json:"start_time" binding:"required"` // ISO 8601 format: 2025-11-27T10:00:00Z
	EndTime    string `json:"end_time" binding:"required"`   // ISO 8601 format: 2025-11-27T11:00:00Z
	Venue      string `json:"venue" binding:"required"`
	Department string `json:"department" binding:"required"`
}

// ScanQRCodeDTO represents the request when a student scans a QR code.
type ScanQRCodeDTO struct {
	QRToken string `json:"qr_token" binding:"required"`
}

// Response DTOs

// GenerateQRCodeResponse represents the response when a QR code is generated.
type GenerateQRCodeResponse struct {
	Message    string `json:"message"`
	EventID    int    `json:"event_id"`
	QRToken    string `json:"qr_token"`
	QRCodeData string `json:"qr_code"` // Base64 encoded PNG image
	CourseName string `json:"course_name"`
	CourseCode string `json:"course_code"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	Venue      string `json:"venue"`
	Department string `json:"department"`
	CreatedBy  string `json:"created_by"` // Lecturer name
	CreatedAt  string `json:"created_at"`
	ExpiresAt  string `json:"expires_at"` // When QR code expires
}

// CheckInResponse represents the response when a student checks in.
type CheckInResponse struct {
	Message      string `json:"message"`
	Status       string `json:"status"` // "present"
	StudentID    int    `json:"student_id"`
	StudentName  string `json:"student_name"`
	MatricNumber string `json:"matric_number"`
	CourseName   string `json:"course_name"`
	CourseCode   string `json:"course_code"`
	MarkedTime   string `json:"marked_time"`
}

// AttendanceRecordResponse represents a single attendance record.
type AttendanceRecordResponse struct {
	ID           int    `json:"id"`
	StudentID    int    `json:"student_id"`
	StudentName  string `json:"student_name"`
	MatricNumber string `json:"matric_number"`
	Status       string `json:"status"` // "present"
	MarkedTime   string `json:"marked_time"`
}

// EventAttendanceResponse represents attendance records for an entire event.
type EventAttendanceResponse struct {
	Message           string                     `json:"message"`
	EventID           int                        `json:"event_id"`
	CourseName        string                     `json:"course_name"`
	CourseCode        string                     `json:"course_code"`
	Department        string                     `json:"department"`
	StartTime         string                     `json:"start_time"`
	EndTime           string                     `json:"end_time"`
	Venue             string                     `json:"venue"`
	CreatedBy         string                     `json:"created_by"` // Lecturer name
	TotalPresent      int                        `json:"total_present"`
	AttendanceRecords []AttendanceRecordResponse `json:"attendance_records"`
	GeneratedAt       string                     `json:"generated_at"`
}

// StudentAttendanceResponse represents attendance history for a student.
type StudentAttendanceResponse struct {
	Message           string                     `json:"message"`
	StudentID         int                        `json:"student_id"`
	StudentName       string                     `json:"student_name"`
	MatricNumber      string                     `json:"matric_number"`
	TotalEvents       int                        `json:"total_events"`
	TotalPresent      int                        `json:"total_present"`
	AttendanceRecords []AttendanceRecordResponse `json:"attendance_records"`
	GeneratedAt       string                     `json:"generated_at"`
}

// Error Response
type ErrorResponse struct {
	Error      string `json:"error"`
	StatusCode int    `json:"status_code"`
}

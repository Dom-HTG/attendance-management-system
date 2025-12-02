package domain

import "time"

// LoginRequest represents admin login credentials
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the admin login response
type LoginResponse struct {
	AccessToken string    `json:"access_token"`
	User        AdminUser `json:"user"`
}

// AdminUser represents admin user information in responses
type AdminUser struct {
	ID         uint      `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	Role       string    `json:"role"`
	Department string    `json:"department"`
	CreatedAt  time.Time `json:"created_at"`
}

// StudentListResponse represents paginated student list
type StudentListResponse struct {
	Students   []StudentSummary `json:"students"`
	Pagination PaginationInfo   `json:"pagination"`
}

// StudentSummary represents student information in list views
type StudentSummary struct {
	StudentID           uint       `json:"student_id"`
	FirstName           string     `json:"first_name"`
	LastName            string     `json:"last_name"`
	Email               string     `json:"email"`
	MatricNumber        string     `json:"matric_number"`
	Department          string     `json:"department"`
	TotalEventsAttended int        `json:"total_events_attended"`
	AttendanceRate      float64    `json:"attendance_rate"`
	CreatedAt           time.Time  `json:"created_at"`
	LastAttendance      *time.Time `json:"last_attendance,omitempty"`
}

// LecturerListResponse represents paginated lecturer list
type LecturerListResponse struct {
	Lecturers  []LecturerSummary `json:"lecturers"`
	Pagination PaginationInfo    `json:"pagination"`
}

// LecturerSummary represents lecturer information in list views
type LecturerSummary struct {
	LecturerID            uint       `json:"lecturer_id"`
	FirstName             string     `json:"first_name"`
	LastName              string     `json:"last_name"`
	Email                 string     `json:"email"`
	StaffID               string     `json:"staff_id"`
	Department            string     `json:"department"`
	TotalEventsCreated    int        `json:"total_events_created"`
	TotalStudentsReached  int        `json:"total_students_reached"`
	AverageAttendanceRate float64    `json:"average_attendance_rate"`
	CreatedAt             time.Time  `json:"created_at"`
	LastEvent             *time.Time `json:"last_event,omitempty"`
}

// UserDetailResponse represents detailed user information
type UserDetailResponse struct {
	UserID           uint               `json:"user_id"`
	UserType         string             `json:"user_type"`
	FirstName        string             `json:"first_name"`
	LastName         string             `json:"last_name"`
	Email            string             `json:"email"`
	MatricNumber     string             `json:"matric_number,omitempty"`
	StaffID          string             `json:"staff_id,omitempty"`
	Department       string             `json:"department"`
	CreatedAt        time.Time          `json:"created_at"`
	Statistics       UserStatistics     `json:"statistics"`
	RecentAttendance []AttendanceRecord `json:"recent_attendance"`
}

// UserStatistics represents user performance statistics
type UserStatistics struct {
	TotalEventsAttended  int        `json:"total_events_attended,omitempty"`
	TotalEventsAvailable int        `json:"total_events_available,omitempty"`
	AttendanceRate       float64    `json:"attendance_rate"`
	FirstAttendance      *time.Time `json:"first_attendance,omitempty"`
	LastAttendance       *time.Time `json:"last_attendance,omitempty"`
	TotalEventsCreated   int        `json:"total_events_created,omitempty"`
	TotalStudentsReached int        `json:"total_students_reached,omitempty"`
}

// AttendanceRecord represents an attendance record
type AttendanceRecord struct {
	EventID    uint      `json:"event_id"`
	CourseCode string    `json:"course_code"`
	CourseName string    `json:"course_name"`
	MarkedTime time.Time `json:"marked_time"`
	Status     string    `json:"status"`
}

// EventListResponse represents paginated event list
type EventListResponse struct {
	Events     []EventSummary `json:"events"`
	Pagination PaginationInfo `json:"pagination"`
}

// EventSummary represents event information in list views
type EventSummary struct {
	EventID         uint      `json:"event_id"`
	CourseCode      string    `json:"course_code"`
	CourseName      string    `json:"course_name"`
	LecturerName    string    `json:"lecturer_name"`
	LecturerEmail   string    `json:"lecturer_email"`
	Department      string    `json:"department"`
	Venue           string    `json:"venue"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	Status          string    `json:"status"`
	TotalAttendance int       `json:"total_attendance"`
	QRGeneratedAt   time.Time `json:"qr_generated_at"`
}

// PaginationInfo represents pagination metadata
type PaginationInfo struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	TotalItems   int `json:"total_items"`
	ItemsPerPage int `json:"items_per_page"`
}

// UpdateStatusRequest represents user status update request
type UpdateStatusRequest struct {
	Active bool   `json:"active"`
	Reason string `json:"reason,omitempty"`
}

// TrendsResponse represents attendance trends data
type TrendsResponse struct {
	Period string       `json:"period"`
	Trends []TrendPoint `json:"trends"`
}

// TrendPoint represents a single trend data point
type TrendPoint struct {
	PeriodLabel     string  `json:"period_label"`
	TotalEvents     int     `json:"total_events"`
	TotalAttendance int     `json:"total_attendance"`
	AttendanceRate  float64 `json:"attendance_rate"`
	UniqueStudents  int     `json:"unique_students"`
}

// LowAttendanceResponse represents students with low attendance
type LowAttendanceResponse struct {
	Threshold      float64         `json:"threshold"`
	StudentsAtRisk []StudentAtRisk `json:"students_at_risk"`
	TotalAtRisk    int             `json:"total_at_risk"`
}

// StudentAtRisk represents a student with low attendance
type StudentAtRisk struct {
	StudentID       uint       `json:"student_id"`
	StudentName     string     `json:"student_name"`
	MatricNumber    string     `json:"matric_number"`
	Department      string     `json:"department"`
	Email           string     `json:"email"`
	AttendanceRate  float64    `json:"attendance_rate"`
	EventsAttended  int        `json:"events_attended"`
	EventsAvailable int        `json:"events_available"`
	LastAttendance  *time.Time `json:"last_attendance,omitempty"`
}

// SystemSettingsResponse represents system settings
type SystemSettingsResponse struct {
	QRCodeValidityMinutes        int    `json:"qr_code_validity_minutes"`
	AttendanceGracePeriodMinutes int    `json:"attendance_grace_period_minutes"`
	LowAttendanceThreshold       int    `json:"low_attendance_threshold"`
	RequireEmailVerification     bool   `json:"require_email_verification"`
	AllowStudentSelfRegistration bool   `json:"allow_student_self_registration"`
	MaxEventsPerDayPerLecturer   int    `json:"max_events_per_day_per_lecturer"`
	AcademicYear                 string `json:"academic_year"`
	Semester                     string `json:"semester"`
}

// UpdateSettingsRequest represents system settings update request
type UpdateSettingsRequest struct {
	QRCodeValidityMinutes        *int    `json:"qr_code_validity_minutes,omitempty"`
	AttendanceGracePeriodMinutes *int    `json:"attendance_grace_period_minutes,omitempty"`
	LowAttendanceThreshold       *int    `json:"low_attendance_threshold,omitempty"`
	RequireEmailVerification     *bool   `json:"require_email_verification,omitempty"`
	AllowStudentSelfRegistration *bool   `json:"allow_student_self_registration,omitempty"`
	MaxEventsPerDayPerLecturer   *int    `json:"max_events_per_day_per_lecturer,omitempty"`
	AcademicYear                 *string `json:"academic_year,omitempty"`
	Semester                     *string `json:"semester,omitempty"`
}

// AuditLogResponse represents paginated audit logs
type AuditLogResponse struct {
	Logs       []AuditLogEntry `json:"logs"`
	Pagination PaginationInfo  `json:"pagination"`
}

// AuditLogEntry represents a single audit log entry
type AuditLogEntry struct {
	LogID        uint      `json:"log_id"`
	Timestamp    time.Time `json:"timestamp"`
	UserType     string    `json:"user_type"`
	UserEmail    string    `json:"user_email"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resource_type"`
	ResourceID   *int      `json:"resource_id,omitempty"`
	Details      string    `json:"details"`
	IPAddress    string    `json:"ip_address"`
}

package entities

import (
	"time"

	"gorm.io/gorm"
)

// Student represents a student in the system.
type Student struct {
	gorm.Model
	FirstName    string `gorm:"column:first_name;not null"`
	LastName     string `gorm:"column:last_name;not null"`
	Email        string `gorm:"column:email;uniqueIndex;not null"`
	Role         string `gorm:"column:role;default:'student'"`
	Password     string `gorm:"column:password;not null"`
	MatricNumber string `gorm:"uniqueIndex;not null;column:matric_number;type:varchar(50)"`
}

// Lecturer represents a lecturer in the system.
type Lecturer struct {
	gorm.Model
	FirstName  string `gorm:"column:first_name;not null"`
	LastName   string `gorm:"column:last_name;not null"`
	Email      string `gorm:"column:email;uniqueIndex;not null"`
	Role       string `gorm:"column:role;default:'lecturer'"`
	Password   string `gorm:"column:password;not null"`
	Department string `gorm:"column:department;not null"`
	StaffID    string `gorm:"uniqueIndex;column:staff_id;not null;type:varchar(50)"`
}

// Event represents a class session.
type Event struct {
	gorm.Model
	EventName   string           `gorm:"column:event_name"`
	StartTime   time.Time        `gorm:"column:start_time"`
	EndTime     time.Time        `gorm:"column:end_time"`
	Venue       string           `gorm:"column:venue"`
	QRCodeToken string           `gorm:"column:qr_code_token"` // Unique token used to generate the QR code
	LecturerID  *int             `gorm:"column:lecturer_id"`   // ID of lecturer who created the event
	CourseCode  string           `gorm:"column:course_code"`
	CourseName  string           `gorm:"column:course_name"`
	Department  string           `gorm:"column:department"`
	Records     []UserAttendance `gorm:"foreignKey:EventID"`
}

// UserAttendance represents an individual attendance record.
type UserAttendance struct {
	gorm.Model
	EventID    int       `gorm:"index;column:event_id;not null"` // References the Event's ID
	Event      Event     `gorm:"foreignKey:EventID;references:ID"`
	StudentID  int       `gorm:"index;column:student_id;not null"` // References the Student's ID
	Student    Student   `gorm:"foreignKey:StudentID;references:ID"`
	Status     string    `gorm:"column:status;default:'present'"`
	MarkedTime time.Time `gorm:"column:marked_time;not null"`
}

// Admin represents an administrator in the system.
type Admin struct {
	gorm.Model
	FirstName    string `gorm:"column:first_name;not null"`
	LastName     string `gorm:"column:last_name;not null"`
	Email        string `gorm:"column:email;uniqueIndex;not null"`
	Password     string `gorm:"column:password;not null"`
	Department   string `gorm:"column:department"`
	Role         string `gorm:"column:role;default:'admin'"`
	IsSuperAdmin bool   `gorm:"column:is_super_admin;default:false"`
	Active       bool   `gorm:"column:active;default:true"`
}

// AuditLog represents a system audit log entry.
type AuditLog struct {
	gorm.Model
	Timestamp    time.Time `gorm:"column:timestamp;default:CURRENT_TIMESTAMP;index"`
	UserType     string    `gorm:"column:user_type;not null;index"` // 'admin', 'lecturer', 'student'
	UserID       int       `gorm:"column:user_id;not null;index"`
	UserEmail    string    `gorm:"column:user_email"`
	Action       string    `gorm:"column:action;not null;index"` // 'login', 'create', 'update', 'delete', 'export'
	ResourceType string    `gorm:"column:resource_type"`
	ResourceID   *int      `gorm:"column:resource_id"`
	Details      string    `gorm:"column:details;type:text"`
	IPAddress    string    `gorm:"column:ip_address;type:varchar(45)"`
	UserAgent    string    `gorm:"column:user_agent;type:text"`
}

// SystemSettings represents system configuration settings.
type SystemSettings struct {
	gorm.Model
	SettingKey   string `gorm:"column:setting_key;uniqueIndex;not null"`
	SettingValue string `gorm:"column:setting_value;type:text;not null"`
	DataType     string `gorm:"column:data_type;not null"` // 'string', 'number', 'boolean'
	Description  string `gorm:"column:description;type:text"`
	UpdatedBy    *int   `gorm:"column:updated_by"`
}

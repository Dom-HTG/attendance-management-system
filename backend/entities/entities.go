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
	EventName   string    `gorm:"column:event_name"`
	StartTime   time.Time `gorm:"column:start_time"`
	EndTime     time.Time `gorm:"column:end_time"`
	Venue       string    `gorm:"column:venue"`
	QRCodeToken string    `gorm:"column:qr_code_token"` // Unique token used to generate the QR code
}

// Attendance represents the overall attendance record for an event.
type Attendance struct {
	gorm.Model
	EventID int              `gorm:"index;column:event_id"`
	Event   Event            `gorm:"foreignKey:EventID;references:ID"`
	Records []UserAttendance `gorm:"foreignKey:AttendanceID"`
}

// UserAttendance represents an individual attendance record.
type UserAttendance struct {
	gorm.Model
	AttendanceID int       `gorm:"index;column:attendance_id"`
	StudentID    int       `gorm:"index;column:student_id"` // References the Student's ID
	Student      Student   `gorm:"foreignKey:StudentID;references:ID"`
	Status       string    `gorm:"column:status"`      // [student, lecturer]
	MarkedTime   time.Time `gorm:"column:marked_time"` // The time when attendance was recorded
}

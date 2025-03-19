package entities

import (
	"time"
)

// Student represents a student in the system.
type Student struct {
	ID           int    `gorm:"primaryKey"`
	FirstName    string `gorm:"column:first_name"`
	LastName     string `gorm:"column:last_name"`
	Email        string `gorm:"column:email;uniqueIndex;not null"`
	Role         string `gorm:"column:role"`
	MatricNumber int    `gorm:"uniqueIndex;not null;column:matric_number"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Lecturer represents a lecturer in the system.
type Lecturer struct {
	ID         int    `gorm:"primaryKey"`
	FirstName  string `gorm:"column:first_name"`
	LastName   string `gorm:"column:last_name"`
	Email      string `gorm:"column:email;uniqueIndex;not null"`
	Role       string `gorm:"column:role"`
	Department string `gorm:"column:department;not null"`
	StaffID    int    `gorm:"uniqueIndex;column:staff_id"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Event represents a class session.
type Event struct {
	ID          int       `gorm:"primaryKey"`
	EventName   string    `gorm:"column:event_name"`
	StartTime   time.Time `gorm:"column:start_time"`
	EndTime     time.Time `gorm:"column:end_time"`
	Venue       string    `gorm:"column:venue"`
	QRCodeToken string    `gorm:"column:qr_code_token"` // Unique token used to generate the QR code
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Attendance represents the overall attendance record for an event.
type Attendance struct {
	ID        int              `gorm:"primaryKey"`
	EventID   int              `gorm:"index;column:event_id"`
	Event     Event            `gorm:"foreignKey:EventID;references:ID"`
	Records   []UserAttendance `gorm:"foreignKey:AttendanceID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserAttendance represents an individual attendance record.
type UserAttendance struct {
	ID           int       `gorm:"primaryKey"`
	AttendanceID int       `gorm:"index;column:attendance_id"`
	StudentID    int       `gorm:"index;column:student_id"` // References the Student's ID
	Student      Student   `gorm:"foreignKey:StudentID;references:ID"`
	Status       string    `gorm:"column:status"`      // [student, lecturer]
	MarkedTime   time.Time `gorm:"column:marked_time"` // The time when attendance was recorded
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

package repository

import (
	"errors"

	"github.com/Dom-HTG/attendance-management-system/entities"
	"gorm.io/gorm"
)

// AttendanceRepoInterface defines the repository interface for attendance operations.
type AttendanceRepoInterface interface {
	// Event operations
	CreateEvent(event *entities.Event) error
	GetEventByQRToken(qrToken string) (*entities.Event, error)
	GetEventByID(eventID int) (*entities.Event, error)

	// Attendance operations
	CreateAttendanceRecord(attendanceRecord *entities.UserAttendance) error
	GetAttendanceByEventID(eventID int) ([]*entities.UserAttendance, error)
	GetStudentAttendance(studentID int) ([]*entities.UserAttendance, error)
	CheckIfStudentMarkedAttendance(eventID, studentID int) (bool, error)
	GetEventWithAttendanceRecords(eventID int) (*entities.Event, []*entities.UserAttendance, error)

	// Database access for PDF export
	DB() *gorm.DB
}

// AttendanceRepo implements the AttendanceRepoInterface.
type AttendanceRepo struct {
	db *gorm.DB
}

// NewAttendanceRepo returns a new instance of AttendanceRepo.
func NewAttendanceRepo(db *gorm.DB) *AttendanceRepo {
	return &AttendanceRepo{
		db: db,
	}
}

// CreateEvent creates a new event in the database.
func (ar *AttendanceRepo) CreateEvent(event *entities.Event) error {
	if err := ar.db.Create(event).Error; err != nil {
		return errors.New("failed to create event: " + err.Error())
	}
	return nil
}

// GetEventByQRToken retrieves an event by its QR token.
func (ar *AttendanceRepo) GetEventByQRToken(qrToken string) (*entities.Event, error) {
	event := &entities.Event{}
	if err := ar.db.Where("qr_code_token = ?", qrToken).First(event).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("qr code not found or invalid")
		}
		return nil, errors.New("failed to retrieve event: " + err.Error())
	}
	return event, nil
}

// GetEventByID retrieves an event by its ID.
func (ar *AttendanceRepo) GetEventByID(eventID int) (*entities.Event, error) {
	event := &entities.Event{}
	if err := ar.db.Where("id = ?", eventID).First(event).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("event not found")
		}
		return nil, errors.New("failed to retrieve event: " + err.Error())
	}
	return event, nil
}

// CreateAttendanceRecord creates a new attendance record for a student.
func (ar *AttendanceRepo) CreateAttendanceRecord(attendanceRecord *entities.UserAttendance) error {
	if err := ar.db.Create(attendanceRecord).Error; err != nil {
		return errors.New("failed to create attendance record: " + err.Error())
	}
	return nil
}

// GetAttendanceByEventID retrieves all attendance records for a specific event.
func (ar *AttendanceRepo) GetAttendanceByEventID(eventID int) ([]*entities.UserAttendance, error) {
	var records []*entities.UserAttendance
	if err := ar.db.Where("event_id = ?", eventID).
		Preload("Student").
		Order("marked_time ASC").
		Find(&records).Error; err != nil {
		return nil, errors.New("failed to retrieve attendance records: " + err.Error())
	}
	return records, nil
}

// GetStudentAttendance retrieves all attendance records for a specific student.
func (ar *AttendanceRepo) GetStudentAttendance(studentID int) ([]*entities.UserAttendance, error) {
	var records []*entities.UserAttendance
	if err := ar.db.Where("student_id = ?", studentID).
		Preload("Student").
		Order("marked_time DESC").
		Find(&records).Error; err != nil {
		return nil, errors.New("failed to retrieve student attendance: " + err.Error())
	}
	return records, nil
}

// CheckIfStudentMarkedAttendance checks if a student has already marked attendance for an event.
func (ar *AttendanceRepo) CheckIfStudentMarkedAttendance(eventID, studentID int) (bool, error) {
	var count int64
	if err := ar.db.
		Where("event_id = ? AND student_id = ?", eventID, studentID).
		Model(&entities.UserAttendance{}).
		Count(&count).Error; err != nil {
		return false, errors.New("failed to check attendance status: " + err.Error())
	}
	return count > 0, nil
}

// GetEventWithAttendanceRecords retrieves an event along with all its attendance records.
func (ar *AttendanceRepo) GetEventWithAttendanceRecords(eventID int) (*entities.Event, []*entities.UserAttendance, error) {
	event, err := ar.GetEventByID(eventID)
	if err != nil {
		return nil, nil, err
	}

	records, err := ar.GetAttendanceByEventID(eventID)
	if err != nil {
		return nil, nil, err
	}

	return event, records, nil
}

// DB returns the underlying database connection for PDF export operations
func (ar *AttendanceRepo) DB() *gorm.DB {
	return ar.db
}

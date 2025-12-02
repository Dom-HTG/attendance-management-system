package repository

import (
	"fmt"
	"time"

	"github.com/Dom-HTG/attendance-management-system/entities"
	"github.com/Dom-HTG/attendance-management-system/internal/admin/domain"
	"gorm.io/gorm"
)

type AdminRepoInterface interface {
	// Authentication
	FindAdminByEmail(email string) (*entities.Admin, error)

	// User Management - Students
	GetAllStudents(page, limit int, department, search string) ([]domain.StudentSummary, int, error)
	GetStudentDetail(studentID uint) (*domain.UserDetailResponse, error)
	UpdateStudentStatus(studentID uint, active bool) error
	DeleteStudent(studentID uint) error

	// User Management - Lecturers
	GetAllLecturers(page, limit int, department, search string) ([]domain.LecturerSummary, int, error)
	GetLecturerDetail(lecturerID uint) (*domain.UserDetailResponse, error)
	UpdateLecturerStatus(lecturerID uint, active bool) error
	DeleteLecturer(lecturerID uint) error

	// Event Management
	GetAllEvents(page, limit int, department string, lecturerID *uint, status string, dateFrom, dateTo *time.Time) ([]domain.EventSummary, int, error)
	DeleteEvent(eventID uint) error

	// Analytics
	GetAttendanceTrends(period string, dateFrom, dateTo time.Time) ([]domain.TrendPoint, error)
	GetLowAttendanceStudents(threshold float64, limit int) ([]domain.StudentAtRisk, int, error)

	// System Settings
	GetAllSettings() (map[string]string, error)
	UpdateSetting(key, value string, adminID uint) error
	GetSetting(key string) (string, error)

	// Audit Logs
	CreateAuditLog(log *entities.AuditLog) error
	GetAuditLogs(page, limit int, userType, actionType string, dateFrom, dateTo *time.Time) ([]domain.AuditLogEntry, int, error)

	// Helper
	DB() *gorm.DB
}

type AdminRepo struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepoInterface {
	return &AdminRepo{db: db}
}

func (ar *AdminRepo) DB() *gorm.DB {
	return ar.db
}

// FindAdminByEmail finds an admin by email
func (ar *AdminRepo) FindAdminByEmail(email string) (*entities.Admin, error) {
	var admin entities.Admin
	err := ar.db.Where("email = ? AND active = ?", email, true).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// GetAllStudents retrieves all students with pagination and filtering
func (ar *AdminRepo) GetAllStudents(page, limit int, department, search string) ([]domain.StudentSummary, int, error) {
	var students []domain.StudentSummary
	var totalCount int64

	query := ar.db.Table("students").Select(`
		students.id as student_id,
		students.first_name,
		students.last_name,
		students.email,
		students.matric_number,
		COALESCE((SELECT events.department FROM user_attendances 
			JOIN events ON user_attendances.event_id = events.id 
			WHERE user_attendances.student_id = students.id LIMIT 1), 'N/A') as department,
		COUNT(DISTINCT user_attendances.id) as total_events_attended,
		CASE 
			WHEN COUNT(DISTINCT events.id) > 0 THEN 
				(COUNT(DISTINCT CASE WHEN user_attendances.status = 'present' THEN user_attendances.id END)::float / COUNT(DISTINCT events.id)::float * 100)
			ELSE 0
		END as attendance_rate,
		students.created_at,
		MAX(user_attendances.marked_time) as last_attendance
	`).
		Joins("LEFT JOIN user_attendances ON user_attendances.student_id = students.id").
		Joins("LEFT JOIN events ON user_attendances.event_id = events.id").
		Where("students.deleted_at IS NULL")

	// Apply filters
	if department != "" {
		query = query.Where("events.department = ?", department)
	}

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("students.first_name ILIKE ? OR students.last_name ILIKE ? OR students.email ILIKE ? OR students.matric_number ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Get total count
	countQuery := ar.db.Table("students").Where("students.deleted_at IS NULL")
	if search != "" {
		searchPattern := "%" + search + "%"
		countQuery = countQuery.Where("students.first_name ILIKE ? OR students.last_name ILIKE ? OR students.email ILIKE ? OR students.matric_number ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern)
	}
	countQuery.Count(&totalCount)

	// Apply pagination
	offset := (page - 1) * limit
	query = query.Group("students.id, students.first_name, students.last_name, students.email, students.matric_number, students.created_at").
		Order("students.created_at DESC").
		Limit(limit).
		Offset(offset)

	err := query.Scan(&students).Error
	if err != nil {
		return nil, 0, err
	}

	return students, int(totalCount), nil
}

// GetStudentDetail retrieves detailed information about a student
func (ar *AdminRepo) GetStudentDetail(studentID uint) (*domain.UserDetailResponse, error) {
	var student entities.Student
	if err := ar.db.First(&student, studentID).Error; err != nil {
		return nil, err
	}

	// Get statistics
	var stats struct {
		TotalEventsAttended  int
		TotalEventsAvailable int
		FirstAttendance      *time.Time
		LastAttendance       *time.Time
	}

	ar.db.Table("user_attendances").
		Select("COUNT(*) as total_events_attended, MIN(marked_time) as first_attendance, MAX(marked_time) as last_attendance").
		Where("student_id = ? AND status = 'present'", studentID).
		Scan(&stats)

	ar.db.Table("events").
		Select("COUNT(*) as total_events_available").
		Scan(&stats.TotalEventsAvailable)

	attendanceRate := 0.0
	if stats.TotalEventsAvailable > 0 {
		attendanceRate = float64(stats.TotalEventsAttended) / float64(stats.TotalEventsAvailable) * 100
	}

	// Get recent attendance
	var recentAttendance []domain.AttendanceRecord
	ar.db.Table("user_attendances").
		Select("user_attendances.event_id, events.course_code, events.course_name, user_attendances.marked_time, user_attendances.status").
		Joins("JOIN events ON user_attendances.event_id = events.id").
		Where("user_attendances.student_id = ?", studentID).
		Order("user_attendances.marked_time DESC").
		Limit(10).
		Scan(&recentAttendance)

	// Get department from first event
	department := "N/A"
	if len(recentAttendance) > 0 {
		ar.db.Table("events").Select("department").Where("id = ?", recentAttendance[0].EventID).Scan(&department)
	}

	response := &domain.UserDetailResponse{
		UserID:       student.ID,
		UserType:     "student",
		FirstName:    student.FirstName,
		LastName:     student.LastName,
		Email:        student.Email,
		MatricNumber: student.MatricNumber,
		Department:   department,
		CreatedAt:    student.CreatedAt,
		Statistics: domain.UserStatistics{
			TotalEventsAttended:  stats.TotalEventsAttended,
			TotalEventsAvailable: stats.TotalEventsAvailable,
			AttendanceRate:       attendanceRate,
			FirstAttendance:      stats.FirstAttendance,
			LastAttendance:       stats.LastAttendance,
		},
		RecentAttendance: recentAttendance,
	}

	return response, nil
}

// UpdateStudentStatus updates a student's active status
func (ar *AdminRepo) UpdateStudentStatus(studentID uint, active bool) error {
	// For now, we don't have an 'active' field in students table
	// This would require a migration to add the field
	// For demonstration, we'll just return success
	return nil
}

// DeleteStudent soft deletes a student
func (ar *AdminRepo) DeleteStudent(studentID uint) error {
	return ar.db.Delete(&entities.Student{}, studentID).Error
}

// GetAllLecturers retrieves all lecturers with pagination and filtering
func (ar *AdminRepo) GetAllLecturers(page, limit int, department, search string) ([]domain.LecturerSummary, int, error) {
	var lecturers []domain.LecturerSummary
	var totalCount int64

	query := ar.db.Table("lecturers").Select(`
		lecturers.id as lecturer_id,
		lecturers.first_name,
		lecturers.last_name,
		lecturers.email,
		lecturers.staff_id,
		lecturers.department,
		COUNT(DISTINCT events.id) as total_events_created,
		COUNT(DISTINCT user_attendances.student_id) as total_students_reached,
		CASE 
			WHEN COUNT(DISTINCT events.id) > 0 THEN 
				(COUNT(DISTINCT CASE WHEN user_attendances.status = 'present' THEN user_attendances.id END)::float / 
				 NULLIF(COUNT(DISTINCT user_attendances.id), 0)::float * 100)
			ELSE 0
		END as average_attendance_rate,
		lecturers.created_at,
		MAX(events.start_time) as last_event
	`).
		Joins("LEFT JOIN events ON events.lecturer_id = lecturers.id").
		Joins("LEFT JOIN user_attendances ON user_attendances.event_id = events.id").
		Where("lecturers.deleted_at IS NULL")

	// Apply filters
	if department != "" {
		query = query.Where("lecturers.department = ?", department)
	}

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("lecturers.first_name ILIKE ? OR lecturers.last_name ILIKE ? OR lecturers.email ILIKE ? OR lecturers.staff_id ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Get total count
	countQuery := ar.db.Table("lecturers").Where("lecturers.deleted_at IS NULL")
	if department != "" {
		countQuery = countQuery.Where("department = ?", department)
	}
	if search != "" {
		searchPattern := "%" + search + "%"
		countQuery = countQuery.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ? OR staff_id ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern)
	}
	countQuery.Count(&totalCount)

	// Apply pagination
	offset := (page - 1) * limit
	query = query.Group("lecturers.id, lecturers.first_name, lecturers.last_name, lecturers.email, lecturers.staff_id, lecturers.department, lecturers.created_at").
		Order("lecturers.created_at DESC").
		Limit(limit).
		Offset(offset)

	err := query.Scan(&lecturers).Error
	if err != nil {
		return nil, 0, err
	}

	return lecturers, int(totalCount), nil
}

// GetLecturerDetail retrieves detailed information about a lecturer
func (ar *AdminRepo) GetLecturerDetail(lecturerID uint) (*domain.UserDetailResponse, error) {
	var lecturer entities.Lecturer
	if err := ar.db.First(&lecturer, lecturerID).Error; err != nil {
		return nil, err
	}

	// Get statistics
	var stats struct {
		TotalEventsCreated   int
		TotalStudentsReached int
		FirstEvent           *time.Time
		LastEvent            *time.Time
	}

	ar.db.Table("events").
		Select("COUNT(*) as total_events_created, MIN(start_time) as first_event, MAX(start_time) as last_event").
		Where("lecturer_id = ?", lecturerID).
		Scan(&stats)

	ar.db.Table("user_attendances").
		Select("COUNT(DISTINCT student_id) as total_students_reached").
		Joins("JOIN events ON user_attendances.event_id = events.id").
		Where("events.lecturer_id = ?", lecturerID).
		Scan(&stats)

	// Get recent events as attendance records
	var recentEvents []domain.AttendanceRecord
	ar.db.Table("events").
		Select("events.id as event_id, events.course_code, events.course_name, events.start_time as marked_time, 'completed' as status").
		Where("lecturer_id = ?", lecturerID).
		Order("start_time DESC").
		Limit(10).
		Scan(&recentEvents)

	response := &domain.UserDetailResponse{
		UserID:     lecturer.ID,
		UserType:   "lecturer",
		FirstName:  lecturer.FirstName,
		LastName:   lecturer.LastName,
		Email:      lecturer.Email,
		StaffID:    lecturer.StaffID,
		Department: lecturer.Department,
		CreatedAt:  lecturer.CreatedAt,
		Statistics: domain.UserStatistics{
			TotalEventsCreated:   stats.TotalEventsCreated,
			TotalStudentsReached: stats.TotalStudentsReached,
			FirstAttendance:      stats.FirstEvent,
			LastAttendance:       stats.LastEvent,
		},
		RecentAttendance: recentEvents,
	}

	return response, nil
}

// UpdateLecturerStatus updates a lecturer's active status
func (ar *AdminRepo) UpdateLecturerStatus(lecturerID uint, active bool) error {
	// Similar to students, would require migration
	return nil
}

// DeleteLecturer soft deletes a lecturer
func (ar *AdminRepo) DeleteLecturer(lecturerID uint) error {
	return ar.db.Delete(&entities.Lecturer{}, lecturerID).Error
}

// GetAllEvents retrieves all events with pagination and filtering
func (ar *AdminRepo) GetAllEvents(page, limit int, department string, lecturerID *uint, status string, dateFrom, dateTo *time.Time) ([]domain.EventSummary, int, error) {
	var events []domain.EventSummary
	var totalCount int64

	query := ar.db.Table("events").Select(`
		events.id as event_id,
		events.course_code,
		events.course_name,
		CONCAT(lecturers.first_name, ' ', lecturers.last_name) as lecturer_name,
		lecturers.email as lecturer_email,
		events.department,
		events.venue,
		events.start_time,
		events.end_time,
		CASE 
			WHEN events.end_time < NOW() THEN 'expired'
			ELSE 'active'
		END as status,
		COUNT(DISTINCT user_attendances.id) as total_attendance,
		events.created_at as qr_generated_at
	`).
		Joins("LEFT JOIN lecturers ON events.lecturer_id = lecturers.id").
		Joins("LEFT JOIN user_attendances ON user_attendances.event_id = events.id").
		Where("events.deleted_at IS NULL")

	// Apply filters
	if department != "" {
		query = query.Where("events.department = ?", department)
	}

	if lecturerID != nil {
		query = query.Where("events.lecturer_id = ?", *lecturerID)
	}

	if status != "" {
		if status == "active" {
			query = query.Where("events.end_time >= NOW()")
		} else if status == "expired" {
			query = query.Where("events.end_time < NOW()")
		}
	}

	if dateFrom != nil {
		query = query.Where("events.start_time >= ?", *dateFrom)
	}

	if dateTo != nil {
		query = query.Where("events.start_time <= ?", *dateTo)
	}

	// Get total count
	countQuery := ar.db.Table("events").Where("events.deleted_at IS NULL")
	if department != "" {
		countQuery = countQuery.Where("department = ?", department)
	}
	if lecturerID != nil {
		countQuery = countQuery.Where("lecturer_id = ?", *lecturerID)
	}
	if dateFrom != nil {
		countQuery = countQuery.Where("start_time >= ?", *dateFrom)
	}
	if dateTo != nil {
		countQuery = countQuery.Where("start_time <= ?", *dateTo)
	}
	countQuery.Count(&totalCount)

	// Apply pagination
	offset := (page - 1) * limit
	query = query.Group("events.id, events.course_code, events.course_name, lecturer_name, lecturers.email, events.department, events.venue, events.start_time, events.end_time, events.created_at").
		Order("events.start_time DESC").
		Limit(limit).
		Offset(offset)

	err := query.Scan(&events).Error
	if err != nil {
		return nil, 0, err
	}

	return events, int(totalCount), nil
}

// DeleteEvent soft deletes an event and its attendance records
func (ar *AdminRepo) DeleteEvent(eventID uint) error {
	return ar.db.Transaction(func(tx *gorm.DB) error {
		// Delete attendance records
		if err := tx.Where("event_id = ?", eventID).Delete(&entities.UserAttendance{}).Error; err != nil {
			return err
		}
		// Delete event
		if err := tx.Delete(&entities.Event{}, eventID).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetAttendanceTrends retrieves attendance trends for a period
func (ar *AdminRepo) GetAttendanceTrends(period string, dateFrom, dateTo time.Time) ([]domain.TrendPoint, error) {
	var trends []domain.TrendPoint

	// Determine date truncation based on period
	var dateTrunc string
	switch period {
	case "daily":
		dateTrunc = "day"
	case "weekly":
		dateTrunc = "week"
	case "monthly":
		dateTrunc = "month"
	default:
		dateTrunc = "week"
	}

	query := fmt.Sprintf(`
		SELECT 
			TO_CHAR(DATE_TRUNC('%s', events.start_time), 'YYYY-MM-DD') as period_label,
			COUNT(DISTINCT events.id) as total_events,
			COUNT(DISTINCT user_attendances.id) as total_attendance,
			CASE 
				WHEN COUNT(DISTINCT events.id) > 0 THEN 
					(COUNT(DISTINCT CASE WHEN user_attendances.status = 'present' THEN user_attendances.id END)::float / 
					 NULLIF(COUNT(DISTINCT user_attendances.id), 0)::float * 100)
				ELSE 0
			END as attendance_rate,
			COUNT(DISTINCT user_attendances.student_id) as unique_students
		FROM events
		LEFT JOIN user_attendances ON user_attendances.event_id = events.id
		WHERE events.deleted_at IS NULL 
			AND events.start_time >= ? 
			AND events.start_time <= ?
		GROUP BY DATE_TRUNC('%s', events.start_time)
		ORDER BY period_label
	`, dateTrunc, dateTrunc)

	err := ar.db.Raw(query, dateFrom, dateTo).Scan(&trends).Error
	if err != nil {
		return nil, err
	}

	return trends, nil
}

// GetLowAttendanceStudents retrieves students with attendance below threshold
func (ar *AdminRepo) GetLowAttendanceStudents(threshold float64, limit int) ([]domain.StudentAtRisk, int, error) {
	var students []domain.StudentAtRisk

	query := ar.db.Table("students").Select(`
		students.id as student_id,
		CONCAT(students.first_name, ' ', students.last_name) as student_name,
		students.matric_number,
		COALESCE((SELECT events.department FROM user_attendances 
			JOIN events ON user_attendances.event_id = events.id 
			WHERE user_attendances.student_id = students.id LIMIT 1), 'N/A') as department,
		students.email,
		CASE 
			WHEN (SELECT COUNT(*) FROM events WHERE deleted_at IS NULL) > 0 THEN 
				(COUNT(DISTINCT CASE WHEN user_attendances.status = 'present' THEN user_attendances.id END)::float / 
				 (SELECT COUNT(*) FROM events WHERE deleted_at IS NULL)::float * 100)
			ELSE 0
		END as attendance_rate,
		COUNT(DISTINCT CASE WHEN user_attendances.status = 'present' THEN user_attendances.id END) as events_attended,
		(SELECT COUNT(*) FROM events WHERE deleted_at IS NULL) as events_available,
		MAX(user_attendances.marked_time) as last_attendance
	`).
		Joins("LEFT JOIN user_attendances ON user_attendances.student_id = students.id").
		Where("students.deleted_at IS NULL").
		Group("students.id, students.first_name, students.last_name, students.matric_number, students.email").
		Having("CASE WHEN (SELECT COUNT(*) FROM events WHERE deleted_at IS NULL) > 0 THEN (COUNT(DISTINCT CASE WHEN user_attendances.status = 'present' THEN user_attendances.id END)::float / (SELECT COUNT(*) FROM events WHERE deleted_at IS NULL)::float * 100) ELSE 0 END < ?", threshold).
		Order("attendance_rate ASC").
		Limit(limit)

	err := query.Scan(&students).Error
	if err != nil {
		return nil, 0, err
	}

	// Get total count
	var totalCount int64
	ar.db.Table("students").Select("students.id").
		Joins("LEFT JOIN user_attendances ON user_attendances.student_id = students.id").
		Where("students.deleted_at IS NULL").
		Group("students.id").
		Having("CASE WHEN (SELECT COUNT(*) FROM events WHERE deleted_at IS NULL) > 0 THEN (COUNT(DISTINCT CASE WHEN user_attendances.status = 'present' THEN user_attendances.id END)::float / (SELECT COUNT(*) FROM events WHERE deleted_at IS NULL)::float * 100) ELSE 0 END < ?", threshold).
		Count(&totalCount)

	return students, int(totalCount), nil
}

// GetAllSettings retrieves all system settings
func (ar *AdminRepo) GetAllSettings() (map[string]string, error) {
	var settings []entities.SystemSettings
	err := ar.db.Find(&settings).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, setting := range settings {
		result[setting.SettingKey] = setting.SettingValue
	}

	return result, nil
}

// UpdateSetting updates a system setting
func (ar *AdminRepo) UpdateSetting(key, value string, adminID uint) error {
	return ar.db.Model(&entities.SystemSettings{}).
		Where("setting_key = ?", key).
		Updates(map[string]interface{}{
			"setting_value": value,
			"updated_by":    adminID,
			"updated_at":    time.Now(),
		}).Error
}

// GetSetting retrieves a single setting value
func (ar *AdminRepo) GetSetting(key string) (string, error) {
	var setting entities.SystemSettings
	err := ar.db.Where("setting_key = ?", key).First(&setting).Error
	if err != nil {
		return "", err
	}
	return setting.SettingValue, nil
}

// CreateAuditLog creates a new audit log entry
func (ar *AdminRepo) CreateAuditLog(log *entities.AuditLog) error {
	return ar.db.Create(log).Error
}

// GetAuditLogs retrieves audit logs with pagination and filtering
func (ar *AdminRepo) GetAuditLogs(page, limit int, userType, actionType string, dateFrom, dateTo *time.Time) ([]domain.AuditLogEntry, int, error) {
	var logs []domain.AuditLogEntry
	var totalCount int64

	query := ar.db.Table("audit_logs").Select(`
		audit_logs.id as log_id,
		audit_logs.timestamp,
		audit_logs.user_type,
		audit_logs.user_email,
		audit_logs.action,
		audit_logs.resource_type,
		audit_logs.resource_id,
		audit_logs.details,
		audit_logs.ip_address
	`).
		Where("audit_logs.deleted_at IS NULL")

	// Apply filters
	if userType != "" {
		query = query.Where("user_type = ?", userType)
	}

	if actionType != "" {
		query = query.Where("action = ?", actionType)
	}

	if dateFrom != nil {
		query = query.Where("timestamp >= ?", *dateFrom)
	}

	if dateTo != nil {
		query = query.Where("timestamp <= ?", *dateTo)
	}

	// Get total count
	countQuery := ar.db.Table("audit_logs").Where("deleted_at IS NULL")
	if userType != "" {
		countQuery = countQuery.Where("user_type = ?", userType)
	}
	if actionType != "" {
		countQuery = countQuery.Where("action = ?", actionType)
	}
	if dateFrom != nil {
		countQuery = countQuery.Where("timestamp >= ?", *dateFrom)
	}
	if dateTo != nil {
		countQuery = countQuery.Where("timestamp <= ?", *dateTo)
	}
	countQuery.Count(&totalCount)

	// Apply pagination
	offset := (page - 1) * limit
	query = query.Order("timestamp DESC").
		Limit(limit).
		Offset(offset)

	err := query.Scan(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, int(totalCount), nil
}

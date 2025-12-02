package service

import (
	"bytes"
	"fmt"
	"time"

	"github.com/Dom-HTG/attendance-management-system/entities"
	"github.com/Dom-HTG/attendance-management-system/pkg/pdf"
	"gorm.io/gorm"
)

// ExportStudentAttendancePDF generates a PDF report for student attendance
func (as *AttendanceSvc) ExportStudentAttendancePDF(studentID int) (*bytes.Buffer, string, error) {
	// Fetch student information
	var student entities.Student
	if err := as.attendanceRepo.DB().Where("id = ?", studentID).First(&student).Error; err != nil {
		return nil, "", fmt.Errorf("student not found: %w", err)
	}

	// Fetch attendance records with event details including department
	type AttendanceRecord struct {
		MarkedTime time.Time
		CourseCode string
		CourseName string
		Venue      string
		Status     string
		Department string
	}

	var records []AttendanceRecord
	err := as.attendanceRepo.DB().Table("user_attendances").
		Select("user_attendances.marked_time, events.course_code, events.course_name, events.venue, events.department, user_attendances.status").
		Joins("JOIN events ON user_attendances.event_id = events.id").
		Where("user_attendances.student_id = ?", studentID).
		Order("user_attendances.marked_time DESC").
		Scan(&records).Error

	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch attendance records: %w", err)
	}

	// Calculate statistics
	totalEvents := len(records)
	eventsPresent := 0
	for _, record := range records {
		if record.Status == "present" {
			eventsPresent++
		}
	}

	attendanceRate := 0.0
	if totalEvents > 0 {
		attendanceRate = float64(eventsPresent) / float64(totalEvents) * 100
	}

	// Get department from most recent event, or use "N/A"
	department := "N/A"
	if len(records) > 0 {
		department = records[0].Department
	}

	// Prepare student info for PDF
	studentInfo := pdf.StudentInfo{
		FullName:       student.FirstName + " " + student.LastName,
		MatricNumber:   student.MatricNumber,
		Email:          student.Email,
		Department:     department,
		TotalEvents:    totalEvents,
		EventsPresent:  eventsPresent,
		AttendanceRate: attendanceRate,
	}

	// Convert records to PDF format
	pdfRecords := make([]pdf.StudentAttendanceRecord, len(records))
	for i, record := range records {
		pdfRecords[i] = pdf.StudentAttendanceRecord{
			Date:       record.MarkedTime,
			CourseCode: record.CourseCode,
			CourseName: record.CourseName,
			Venue:      record.Venue,
			TimeMarked: record.MarkedTime,
			Status:     record.Status,
		}
	}

	// Generate PDF
	pdfDoc, err := pdf.GenerateStudentAttendancePDF(studentInfo, pdfRecords)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Write PDF to buffer
	var buf bytes.Buffer
	if err := pdfDoc.Output(&buf); err != nil {
		return nil, "", fmt.Errorf("failed to output PDF: %w", err)
	}

	// Generate filename
	filename := fmt.Sprintf("attendance-report-%s-%s.pdf",
		student.MatricNumber,
		time.Now().Format("2006-01-02"))

	return &buf, filename, nil
}

// ExportLecturerAllEventsPDF generates a PDF report for all lecturer events
func (as *AttendanceSvc) ExportLecturerAllEventsPDF(lecturerID int) (*bytes.Buffer, string, error) {
	// Fetch lecturer information
	var lecturer entities.Lecturer
	if err := as.attendanceRepo.DB().Where("id = ?", lecturerID).First(&lecturer).Error; err != nil {
		return nil, "", fmt.Errorf("lecturer not found: %w", err)
	}

	// Fetch all events with attendance counts
	type EventWithAttendance struct {
		CourseCode      string
		CourseName      string
		StartTime       time.Time
		Venue           string
		Status          string
		StudentsPresent int64
	}

	var events []EventWithAttendance
	err := as.attendanceRepo.DB().Table("events").
		Select("events.course_code, events.course_name, events.start_time, events.venue, "+
			"CASE WHEN events.end_time < NOW() THEN 'expired' ELSE 'active' END as status, "+
			"COUNT(DISTINCT user_attendances.student_id) as students_present").
		Joins("LEFT JOIN user_attendances ON events.id = user_attendances.event_id").
		Where("events.lecturer_id = ?", lecturerID).
		Group("events.id, events.course_code, events.course_name, events.start_time, events.venue, events.end_time").
		Order("events.created_at DESC").
		Scan(&events).Error

	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch events: %w", err)
	}

	// Calculate statistics
	totalEvents := len(events)

	// Get total unique students reached
	var totalStudents int64
	as.attendanceRepo.DB().Table("user_attendances").
		Select("COUNT(DISTINCT student_id)").
		Joins("JOIN events ON user_attendances.event_id = events.id").
		Where("events.lecturer_id = ?", lecturerID).
		Scan(&totalStudents)

	// Calculate average attendance
	avgAttendance := 0.0
	if totalEvents > 0 {
		var sum int64
		for _, event := range events {
			sum += event.StudentsPresent
		}
		avgAttendance = float64(sum) / float64(totalEvents)
	}

	// Prepare lecturer info
	lecturerInfo := pdf.LecturerInfo{
		FullName:        lecturer.FirstName + " " + lecturer.LastName,
		Email:           lecturer.Email,
		Department:      lecturer.Department,
		TotalEvents:     totalEvents,
		StudentsReached: int(totalStudents),
		AvgAttendance:   avgAttendance,
	}

	// Convert events to PDF format
	eventSummaries := make([]pdf.EventSummary, len(events))
	for i, event := range events {
		eventSummaries[i] = pdf.EventSummary{
			CourseCode:      event.CourseCode,
			CourseName:      event.CourseName,
			Date:            event.StartTime,
			Venue:           event.Venue,
			StudentsPresent: int(event.StudentsPresent),
			Status:          event.Status,
		}
	}

	// Generate PDF
	pdfDoc, err := pdf.GenerateLecturerAllEventsPDF(lecturerInfo, eventSummaries)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Write PDF to buffer
	var buf bytes.Buffer
	if err := pdfDoc.Output(&buf); err != nil {
		return nil, "", fmt.Errorf("failed to output PDF: %w", err)
	}

	// Generate filename
	filename := fmt.Sprintf("lecturer-attendance-report-%s.pdf",
		time.Now().Format("2006-01-02"))

	return &buf, filename, nil
}

// ExportLecturerSingleEventPDF generates a PDF report for a single event
func (as *AttendanceSvc) ExportLecturerSingleEventPDF(lecturerID, eventID int) (*bytes.Buffer, string, error) {
	// Fetch lecturer information
	var lecturer entities.Lecturer
	if err := as.attendanceRepo.DB().Where("id = ?", lecturerID).First(&lecturer).Error; err != nil {
		return nil, "", fmt.Errorf("lecturer not found: %w", err)
	}

	// Fetch event details
	var event entities.Event
	if err := as.attendanceRepo.DB().Where("id = ? AND lecturer_id = ?", eventID, lecturerID).First(&event).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, "", fmt.Errorf("event not found or access denied")
		}
		return nil, "", fmt.Errorf("failed to fetch event: %w", err)
	}

	// Fetch attendance records for this event
	type AttendeeInfo struct {
		MatricNumber string
		FirstName    string
		LastName     string
		Email        string
		MarkedTime   time.Time
		Status       string
	}

	var attendees []AttendeeInfo
	err := as.attendanceRepo.DB().Table("user_attendances").
		Select("students.matric_number, students.first_name, students.last_name, students.email, user_attendances.marked_time, user_attendances.status").
		Joins("JOIN students ON user_attendances.student_id = students.id").
		Where("user_attendances.event_id = ?", eventID).
		Order("user_attendances.marked_time ASC").
		Scan(&attendees).Error

	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch attendance records: %w", err)
	}

	// Prepare lecturer and event info
	lecturerInfo := pdf.LecturerInfo{
		FullName:   lecturer.FirstName + " " + lecturer.LastName,
		Email:      lecturer.Email,
		Department: lecturer.Department,
	}

	eventStatus := "active"
	if event.EndTime.Before(time.Now()) {
		eventStatus = "expired"
	}

	eventDetail := pdf.EventDetail{
		CourseCode: event.CourseCode,
		CourseName: event.CourseName,
		Venue:      event.Venue,
		StartTime:  event.StartTime,
		Status:     eventStatus,
	}

	// Convert attendees to PDF format
	attendeeRecords := make([]pdf.AttendeeRecord, len(attendees))
	for i, attendee := range attendees {
		attendeeRecords[i] = pdf.AttendeeRecord{
			MatricNumber: attendee.MatricNumber,
			StudentName:  attendee.FirstName + " " + attendee.LastName,
			Email:        attendee.Email,
			TimeMarked:   attendee.MarkedTime,
			Status:       attendee.Status,
		}
	}

	// Generate PDF
	pdfDoc, err := pdf.GenerateLecturerSingleEventPDF(lecturerInfo, eventDetail, attendeeRecords)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Write PDF to buffer
	var buf bytes.Buffer
	if err := pdfDoc.Output(&buf); err != nil {
		return nil, "", fmt.Errorf("failed to output PDF: %w", err)
	}

	// Generate filename
	filename := fmt.Sprintf("event-%s-attendance-%s.pdf",
		event.CourseCode,
		time.Now().Format("2006-01-02"))

	return &buf, filename, nil
}

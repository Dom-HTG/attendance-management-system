package service

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/Dom-HTG/attendance-management-system/entities"
	attendance "github.com/Dom-HTG/attendance-management-system/internal/attendance/domain"
	"github.com/Dom-HTG/attendance-management-system/internal/attendance/repository"
	authDomain "github.com/Dom-HTG/attendance-management-system/internal/auth/domain"
	authRepo "github.com/Dom-HTG/attendance-management-system/internal/auth/repository"
	"github.com/Dom-HTG/attendance-management-system/pkg/middleware"
	"github.com/Dom-HTG/attendance-management-system/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AttendanceSvcInterface defines the service interface for attendance operations.
type AttendanceSvcInterface interface {
	GenerateQRCode(ctx *gin.Context)
	CheckIn(ctx *gin.Context)
	GetEventAttendance(ctx *gin.Context)
	GetStudentAttendance(ctx *gin.Context)
}

// AttendanceSvc implements the AttendanceSvcInterface.
type AttendanceSvc struct {
	attendanceRepo repository.AttendanceRepoInterface
	authRepo       authRepo.AuthRepoInterface
}

// NewAttendanceSvc returns a new instance of AttendanceSvc.
func NewAttendanceSvc(attendanceRepo repository.AttendanceRepoInterface, authRepo authRepo.AuthRepoInterface) *AttendanceSvc {
	return &AttendanceSvc{
		attendanceRepo: attendanceRepo,
		authRepo:       authRepo,
	}
}

// GenerateQRCode handles QR code generation for lecturers.
// Only lecturers can generate QR codes for events.
func (as *AttendanceSvc) GenerateQRCode(ctx *gin.Context) {
	// Get user ID and role from context (set by AuthMiddleware and RoleMiddleware)
	_, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "user id not found in context",
		})
		return
	}

	// Bind request body
	var req attendance.GenerateQRCodeDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Parse start and end times
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid start_time format. expected RFC3339 format (e.g., 2025-11-27T10:00:00Z)",
		})
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid end_time format. expected RFC3339 format (e.g., 2025-11-27T11:00:00Z)",
		})
		return
	}

	// Validate that end time is after start time
	if endTime.Before(startTime) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "end_time must be after start_time",
		})
		return
	}

	// Generate a unique QR token
	qrToken := uuid.New().String()

	// Get lecturer ID from context
	lecturerID, _ := middleware.GetUserIDFromContext(ctx)

	// Create the event
	event := &entities.Event{
		EventName:   fmt.Sprintf("%s (%s)", req.CourseName, req.CourseCode),
		StartTime:   startTime,
		EndTime:     endTime,
		Venue:       req.Venue,
		QRCodeToken: qrToken,
		LecturerID:  &lecturerID,
		CourseCode:  req.CourseCode,
		CourseName:  req.CourseName,
		Department:  req.Department,
	}

	if err := as.attendanceRepo.CreateEvent(event); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to create event",
			"details": err.Error(),
		})
		return
	}

	// Generate QR code with the token as data
	qrCodeData, err := utils.GenerateQRCodePNG(qrToken, 256)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to generate QR code",
			"details": err.Error(),
		})
		return
	}

	// Get lecturer information for the response
	lecturerEmail, ok := middleware.GetUserEmailFromContext(ctx)
	if !ok {
		lecturerEmail = "unknown@fupre.edu"
	}

	lecturerResp, err := as.authRepo.FindLecturerByEmail(lecturerEmail)
	if err != nil {
		// If we can't find lecturer details, just use a placeholder
		lecturerResp = &authDomain.LecturerResponse{
			FirstName: "Unknown",
			LastName:  "Lecturer",
		}
	}

	lecturerName := fmt.Sprintf("%s %s", lecturerResp.FirstName, lecturerResp.LastName)

	// Prepare response
	response := attendance.GenerateQRCodeResponse{
		Message:    "QR code generated successfully",
		EventID:    int(event.ID),
		QRToken:    qrToken,
		QRCodeData: qrCodeData,
		CourseName: req.CourseName,
		CourseCode: req.CourseCode,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		Venue:      req.Venue,
		Department: req.Department,
		CreatedBy:  lecturerName,
		CreatedAt:  event.CreatedAt.Format(time.RFC3339),
		ExpiresAt:  endTime.Format(time.RFC3339),
	}

	ctx.JSON(http.StatusCreated, response)
}

// CheckIn handles student check-in when they scan a QR code.
func (as *AttendanceSvc) CheckIn(ctx *gin.Context) {
	// Get user ID and role from context
	studentID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "user id not found in context",
		})
		return
	}

	// Bind request body
	var req attendance.ScanQRCodeDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate QR token
	if !utils.ValidateQRCodeToken(req.QRToken) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid QR token",
		})
		return
	}

	// Get the event by QR token
	event, err := as.attendanceRepo.GetEventByQRToken(req.QRToken)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check if the event is still active (within time range)
	now := time.Now()
	if now.Before(event.StartTime) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":      "event has not started yet",
			"start_time": event.StartTime.Format(time.RFC3339),
		})
		return
	}

	if now.After(event.EndTime) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":    "event has ended",
			"end_time": event.EndTime.Format(time.RFC3339),
		})
		return
	}

	// Check if the student has already marked attendance for this event
	alreadyMarked, err := as.attendanceRepo.CheckIfStudentMarkedAttendance(int(event.ID), studentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to check attendance status",
			"details": err.Error(),
		})
		return
	}

	if alreadyMarked {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "you have already checked in for this event",
		})
		return
	}

	// Create attendance record
	attendanceRecord := &entities.UserAttendance{
		EventID:    int(event.ID),
		StudentID:  studentID,
		Status:     "present",
		MarkedTime: now,
	}

	if err := as.attendanceRepo.CreateAttendanceRecord(attendanceRecord); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to record attendance",
			"details": err.Error(),
		})
		return
	}

	// Get student information from database
	userEmail, _ := middleware.GetUserEmailFromContext(ctx)
	student, err := as.authRepo.FindStudentByEmail(userEmail)
	if err != nil {
		// Fallback to basic info if we can't find the student
		student = &authDomain.StudentResponse{
			FirstName:    "Unknown",
			LastName:     "Student",
			MatricNumber: "N/A",
		}
	}

	studentFullName := fmt.Sprintf("%s %s", student.FirstName, student.LastName)

	// Extract course code from event name (format: "CourseName (CourseCode)")
	courseCode := ""
	if len(event.EventName) > 0 {
		// Try to extract course code from parentheses
		if start := len(event.EventName) - 1; start > 0 {
			for i := len(event.EventName) - 1; i >= 0; i-- {
				if event.EventName[i] == '(' && i+1 < len(event.EventName) {
					end := len(event.EventName)
					if event.EventName[end-1] == ')' {
						courseCode = event.EventName[i+1 : end-1]
						break
					}
				}
			}
		}
	}

	response := attendance.CheckInResponse{
		Message:      "Check-in successful",
		Status:       "present",
		StudentID:    studentID,
		StudentName:  studentFullName,
		MatricNumber: student.MatricNumber,
		CourseName:   event.EventName,
		CourseCode:   courseCode,
		MarkedTime:   now.Format(time.RFC3339),
	}

	ctx.JSON(http.StatusOK, response)
}

// GetEventAttendance retrieves attendance records for a specific event.
// Only lecturers can access this endpoint.
func (as *AttendanceSvc) GetEventAttendance(ctx *gin.Context) {
	// Get event ID from URL parameter
	eventID, ok := ctx.Params.Get("event_id")
	if !ok || eventID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "event_id parameter is required",
		})
		return
	}

	// Convert to integer
	var eventIDInt int
	if _, err := fmt.Sscanf(eventID, "%d", &eventIDInt); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "event_id must be a valid integer",
		})
		return
	}

	// Get event and attendance records
	event, records, err := as.attendanceRepo.GetEventWithAttendanceRecords(eventIDInt)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Build attendance records response
	attendanceRecords := []attendance.AttendanceRecordResponse{}
	for _, record := range records {
		attendanceRecords = append(attendanceRecords, attendance.AttendanceRecordResponse{
			ID:           int(record.ID),
			StudentID:    record.StudentID,
			StudentName:  fmt.Sprintf("%s %s", record.Student.FirstName, record.Student.LastName),
			MatricNumber: record.Student.MatricNumber,
			Status:       record.Status,
			MarkedTime:   record.MarkedTime.Format(time.RFC3339),
		})
	}

	// Extract course code from event name (format: "CourseName (CourseCode)")
	courseCode := ""
	if len(event.EventName) > 0 {
		for i := len(event.EventName) - 1; i >= 0; i-- {
			if event.EventName[i] == '(' && i+1 < len(event.EventName) {
				end := len(event.EventName)
				if event.EventName[end-1] == ')' {
					courseCode = event.EventName[i+1 : end-1]
					break
				}
			}
		}
	}

	// Get lecturer information
	lecturerEmail, _ := middleware.GetUserEmailFromContext(ctx)
	lecturer, err := as.authRepo.FindLecturerByEmail(lecturerEmail)
	lecturerName := "Unknown Lecturer"
	department := ""
	if err == nil {
		lecturerName = fmt.Sprintf("%s %s", lecturer.FirstName, lecturer.LastName)
		department = lecturer.Department
	}

	response := attendance.EventAttendanceResponse{
		Message:           "Attendance records retrieved successfully",
		EventID:           int(event.ID),
		CourseName:        event.EventName,
		CourseCode:        courseCode,
		Department:        department,
		StartTime:         event.StartTime.Format(time.RFC3339),
		EndTime:           event.EndTime.Format(time.RFC3339),
		Venue:             event.Venue,
		CreatedBy:         lecturerName,
		TotalPresent:      len(attendanceRecords),
		AttendanceRecords: attendanceRecords,
		GeneratedAt:       time.Now().Format(time.RFC3339),
	}

	ctx.JSON(http.StatusOK, response)
}

// GetStudentAttendance retrieves attendance history for a specific student.
func (as *AttendanceSvc) GetStudentAttendance(ctx *gin.Context) {
	// Get user ID and role from context
	studentID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "user id not found in context",
		})
		return
	}

	// Get student attendance records
	records, err := as.attendanceRepo.GetStudentAttendance(studentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to retrieve attendance records",
			"details": err.Error(),
		})
		return
	}

	// Build attendance records response
	attendanceRecords := []attendance.AttendanceRecordResponse{}
	for _, record := range records {
		attendanceRecords = append(attendanceRecords, attendance.AttendanceRecordResponse{
			ID:           int(record.ID),
			StudentID:    record.StudentID,
			StudentName:  fmt.Sprintf("%s %s", record.Student.FirstName, record.Student.LastName),
			MatricNumber: record.Student.MatricNumber,
			Status:       record.Status,
			MarkedTime:   record.MarkedTime.Format(time.RFC3339),
		})
	}

	// Get student information from database
	userEmail, _ := middleware.GetUserEmailFromContext(ctx)
	student, err := as.authRepo.FindStudentByEmail(userEmail)
	studentFullName := "Unknown Student"
	matricNumber := "N/A"
	if err == nil {
		studentFullName = fmt.Sprintf("%s %s", student.FirstName, student.LastName)
		matricNumber = student.MatricNumber
	}

	response := attendance.StudentAttendanceResponse{
		Message:           "Student attendance records retrieved successfully",
		StudentID:         studentID,
		StudentName:       studentFullName,
		MatricNumber:      matricNumber,
		TotalEvents:       len(attendanceRecords),
		TotalPresent:      len(attendanceRecords), // Assuming all records are "present"
		AttendanceRecords: attendanceRecords,
		GeneratedAt:       time.Now().Format(time.RFC3339),
	}

	ctx.JSON(http.StatusOK, response)
}

// ExportStudentAttendancePDFHandler handles PDF export requests for student attendance
func (as *AttendanceSvc) ExportStudentAttendancePDFHandler(ctx *gin.Context) {
	// Get student ID from context
	studentID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized: user ID not found",
		})
		return
	}

	// Generate PDF
	pdfBuffer, filename, err := as.ExportStudentAttendancePDF(studentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate PDF report",
			"error":   err.Error(),
		})
		return
	}

	// Set headers for PDF download
	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	ctx.Header("Content-Length", fmt.Sprintf("%d", pdfBuffer.Len()))

	// Send PDF
	ctx.Data(http.StatusOK, "application/pdf", pdfBuffer.Bytes())
}

// ExportLecturerAttendancePDFHandler handles PDF export requests for lecturer attendance
func (as *AttendanceSvc) ExportLecturerAttendancePDFHandler(ctx *gin.Context) {
	// Get lecturer ID from context
	lecturerID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized: user ID not found",
		})
		return
	}

	// Check if event_id parameter is provided
	eventIDStr := ctx.Query("event_id")

	var pdfBuffer *bytes.Buffer
	var filename string
	var err error

	if eventIDStr != "" {
		// Export single event
		eventID := 0
		if _, scanErr := fmt.Sscanf(eventIDStr, "%d", &eventID); scanErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid event_id parameter",
			})
			return
		}

		pdfBuffer, filename, err = as.ExportLecturerSingleEventPDF(lecturerID, eventID)
	} else {
		// Export all events
		pdfBuffer, filename, err = as.ExportLecturerAllEventsPDF(lecturerID)
	}

	if err != nil {
		statusCode := http.StatusInternalServerError
		message := "Failed to generate PDF report"

		if err.Error() == "event not found or access denied" {
			statusCode = http.StatusNotFound
			message = "Event not found or access denied"
		}

		ctx.JSON(statusCode, gin.H{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
		return
	}

	// Set headers for PDF download
	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	ctx.Header("Content-Length", fmt.Sprintf("%d", pdfBuffer.Len()))

	// Send PDF
	ctx.Data(http.StatusOK, "application/pdf", pdfBuffer.Bytes())
}

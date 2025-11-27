package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Dom-HTG/attendance-management-system/entities"
	attendance "github.com/Dom-HTG/attendance-management-system/internal/attendance/domain"
	"github.com/Dom-HTG/attendance-management-system/internal/attendance/repository"
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
	lecturerID, ok := middleware.GetUserIDFromContext(ctx)
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

	// Create the event
	event := &entities.Event{
		EventName:   fmt.Sprintf("%s (%s)", req.CourseName, req.CourseCode),
		StartTime:   startTime,
		EndTime:     endTime,
		Venue:       req.Venue,
		QRCodeToken: qrToken,
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
	lecturer, err := as.authRepo.FindLecturerByEmail("")
	if err != nil {
		// If we can't find lecturer details, just use the ID
		lecturer = &entities.Lecturer{
			FirstName: "Unknown",
			LastName:  "Lecturer",
		}
	}

	lecturerName := fmt.Sprintf("%s %s", lecturer.FirstName, lecturer.LastName)

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
		AttendanceID: int(event.ID),
		StudentID:    studentID,
		Status:       "present",
		MarkedTime:   now,
	}

	if err := as.attendanceRepo.CreateAttendanceRecord(attendanceRecord); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to record attendance",
			"details": err.Error(),
		})
		return
	}

	// Get student information
	// For now, we'll fetch from context or use placeholder
	userEmail, _ := middleware.GetUserEmailFromContext(ctx)

	response := attendance.CheckInResponse{
		Message:      "Check-in successful",
		Status:       "present",
		StudentID:    studentID,
		StudentName:  userEmail, // In a real app, fetch full name from database
		MatricNumber: "",        // Would need to fetch from student record
		CourseName:   event.EventName,
		CourseCode:   "", // Would need to extract or store separately
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

	response := attendance.EventAttendanceResponse{
		Message:           "Attendance records retrieved successfully",
		EventID:           int(event.ID),
		CourseName:        event.EventName,
		CourseCode:        "", // Would need to parse from EventName or store separately
		Department:        "",
		StartTime:         event.StartTime.Format(time.RFC3339),
		EndTime:           event.EndTime.Format(time.RFC3339),
		Venue:             event.Venue,
		CreatedBy:         "",
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

	response := attendance.StudentAttendanceResponse{
		Message:           "Student attendance records retrieved successfully",
		StudentID:         studentID,
		StudentName:       "", // Would need to fetch from database
		MatricNumber:      "", // Would need to fetch from database
		TotalEvents:       len(attendanceRecords),
		TotalPresent:      len(attendanceRecords), // Assuming all records are "present"
		AttendanceRecords: attendanceRecords,
		GeneratedAt:       time.Now().Format(time.RFC3339),
	}

	ctx.JSON(http.StatusOK, response)
}

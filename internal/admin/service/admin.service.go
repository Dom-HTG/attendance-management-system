package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Dom-HTG/attendance-management-system/entities"
	"github.com/Dom-HTG/attendance-management-system/internal/admin/domain"
	"github.com/Dom-HTG/attendance-management-system/internal/admin/repository"
	"github.com/Dom-HTG/attendance-management-system/pkg/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminServiceInterface interface {
	// Authentication
	LoginAdmin(ctx *gin.Context)

	// User Management
	GetAllStudents(ctx *gin.Context)
	GetAllLecturers(ctx *gin.Context)
	GetUserDetail(ctx *gin.Context)
	UpdateUserStatus(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)

	// Event Management
	GetAllEvents(ctx *gin.Context)
	DeleteEvent(ctx *gin.Context)

	// Analytics
	GetAttendanceTrends(ctx *gin.Context)
	GetLowAttendanceStudents(ctx *gin.Context)

	// Settings
	GetSystemSettings(ctx *gin.Context)
	UpdateSystemSettings(ctx *gin.Context)

	// Audit Logs
	GetAuditLogs(ctx *gin.Context)

	// Helper
	LogAction(userType string, userID int, userEmail, action, resourceType string, resourceID *int, details, ipAddress, userAgent string) error
}

type AdminSvc struct {
	adminRepo repository.AdminRepoInterface
}

func NewAdminService(adminRepo repository.AdminRepoInterface) AdminServiceInterface {
	return &AdminSvc{adminRepo: adminRepo}
}

// LoginAdmin handles admin authentication
func (as *AdminSvc) LoginAdmin(ctx *gin.Context) {
	var req domain.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "Email and password are required",
		})
		return
	}

	// Find admin by email
	admin, err := as.adminRepo.FindAdminByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(401, gin.H{
				"success": false,
				"message": "Invalid email or password",
			})
			return
		}
		ctx.JSON(500, gin.H{
			"success": false,
			"message": "An error occurred during login",
		})
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password))
	if err != nil {
		ctx.JSON(401, gin.H{
			"success": false,
			"message": "Invalid email or password",
		})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(int(admin.ID), admin.Email, admin.Role, 60*24*7) // 7 days
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"message": "Failed to generate authentication token",
		})
		return
	}

	// Log admin login
	as.LogAction("admin", int(admin.ID), admin.Email, "login", "admin", nil, "Admin logged in successfully", ctx.ClientIP(), ctx.GetHeader("User-Agent"))

	// Return response
	ctx.JSON(200, gin.H{
		"success": true,
		"message": "Admin login successful",
		"data": domain.LoginResponse{
			AccessToken: token,
			User: domain.AdminUser{
				ID:         admin.ID,
				FirstName:  admin.FirstName,
				LastName:   admin.LastName,
				Email:      admin.Email,
				Role:       admin.Role,
				Department: admin.Department,
				CreatedAt:  admin.CreatedAt,
			},
		},
	})
}

// GetAllStudents retrieves all students with pagination
func (as *AdminSvc) GetAllStudents(ctx *gin.Context) {
	// Get query parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "50"))
	department := ctx.Query("department")
	search := ctx.Query("search")

	// Fetch students
	students, totalCount, err := as.adminRepo.GetAllStudents(page, limit, department, search)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"message": "Failed to retrieve students",
		})
		return
	}

	// Calculate pagination
	totalPages := (totalCount + limit - 1) / limit

	ctx.JSON(200, gin.H{
		"success": true,
		"data": domain.StudentListResponse{
			Students: students,
			Pagination: domain.PaginationInfo{
				CurrentPage:  page,
				TotalPages:   totalPages,
				TotalItems:   totalCount,
				ItemsPerPage: limit,
			},
		},
	})
}

// GetAllLecturers retrieves all lecturers with pagination
func (as *AdminSvc) GetAllLecturers(ctx *gin.Context) {
	// Get query parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "50"))
	department := ctx.Query("department")
	search := ctx.Query("search")

	// Fetch lecturers
	lecturers, totalCount, err := as.adminRepo.GetAllLecturers(page, limit, department, search)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"message": "Failed to retrieve lecturers",
		})
		return
	}

	// Calculate pagination
	totalPages := (totalCount + limit - 1) / limit

	ctx.JSON(200, gin.H{
		"success": true,
		"data": domain.LecturerListResponse{
			Lecturers: lecturers,
			Pagination: domain.PaginationInfo{
				CurrentPage:  page,
				TotalPages:   totalPages,
				TotalItems:   totalCount,
				ItemsPerPage: limit,
			},
		},
	})
}

// GetUserDetail retrieves detailed information about a user
func (as *AdminSvc) GetUserDetail(ctx *gin.Context) {
	// Determine user type from URL path
	path := ctx.FullPath()
	var userType string
	if strings.Contains(path, "/student/") {
		userType = "student"
	} else if strings.Contains(path, "/lecturer/") {
		userType = "lecturer"
	} else {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "Invalid user type in URL",
		})
		return
	}

	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	var userDetail *domain.UserDetailResponse

	switch userType {
	case "student":
		userDetail, err = as.adminRepo.GetStudentDetail(uint(userID))
	case "lecturer":
		userDetail, err = as.adminRepo.GetLecturerDetail(uint(userID))
	default:
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "Invalid user type. Must be 'student' or 'lecturer'",
		})
		return
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(404, gin.H{
				"success": false,
				"message": "User not found",
			})
			return
		}
		ctx.JSON(500, gin.H{
			"success": false,
			"message": "Failed to retrieve user details",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"data":    userDetail,
	})
}

// UpdateUserStatus updates a user's active status
func (as *AdminSvc) UpdateUserStatus(ctx *gin.Context) {
	// Determine user type from URL path
	path := ctx.FullPath()
	var userType string
	if strings.Contains(path, "/student/") {
		userType = "student"
	} else if strings.Contains(path, "/lecturer/") {
		userType = "lecturer"
	} else {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "Invalid user type in URL",
		})
		return
	}

	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	var req domain.UpdateStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	switch userType {
	case "student":
		err = as.adminRepo.UpdateStudentStatus(uint(userID), req.Active)
	case "lecturer":
		err = as.adminRepo.UpdateLecturerStatus(uint(userID), req.Active)
	default:
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "Invalid user type",
		})
		return
	}

	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"message": "Failed to update user status",
		})
		return
	}

	// Get admin info from context
	adminID := ctx.GetInt("id")
	adminEmail := ctx.GetString("email")

	// Log the action
	resourceID := int(userID)
	details := fmt.Sprintf("Updated %s status to active=%v. Reason: %s", userType, req.Active, req.Reason)
	as.LogAction("admin", adminID, adminEmail, "update", userType, &resourceID, details, ctx.ClientIP(), ctx.GetHeader("User-Agent"))

	ctx.JSON(200, gin.H{
		"success": true,
		"message": "User status updated successfully",
		"data": gin.H{
			"user_id":    userID,
			"user_type":  userType,
			"active":     req.Active,
			"updated_at": time.Now(),
		},
	})
}

// DeleteUser deletes a user
func (as *AdminSvc) DeleteUser(ctx *gin.Context) {
	// Determine user type from URL path
	path := ctx.FullPath()
	var userType string
	if strings.Contains(path, "/student/") {
		userType = "student"
	} else if strings.Contains(path, "/lecturer/") {
		userType = "lecturer"
	} else {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "Invalid user type in URL",
		})
		return
	}

	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	switch userType {
	case "student":
		err = as.adminRepo.DeleteStudent(uint(userID))
	case "lecturer":
		err = as.adminRepo.DeleteLecturer(uint(userID))
	default:
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "Invalid user type",
		})
		return
	}

	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"message": "Failed to delete user",
		})
		return
	}

	// Get admin info from context
	adminID := ctx.GetInt("id")
	adminEmail := ctx.GetString("email")

	// Log the action
	resourceID := int(userID)
	details := fmt.Sprintf("Deleted %s with ID %d", userType, userID)
	as.LogAction("admin", adminID, adminEmail, "delete", userType, &resourceID, details, ctx.ClientIP(), ctx.GetHeader("User-Agent"))

	ctx.JSON(200, gin.H{
		"success": true,
		"message": "User deleted successfully",
	})
}

// GetAllEvents retrieves all events with filtering
func (as *AdminSvc) GetAllEvents(ctx *gin.Context) {
	// Get query parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "50"))
	department := ctx.Query("department")
	status := ctx.Query("status")

	var lecturerID *uint
	if lid := ctx.Query("lecturer_id"); lid != "" {
		if id, err := strconv.ParseUint(lid, 10, 32); err == nil {
			uid := uint(id)
			lecturerID = &uid
		}
	}

	var dateFrom, dateTo *time.Time
	if df := ctx.Query("date_from"); df != "" {
		if t, err := time.Parse("2006-01-02", df); err == nil {
			dateFrom = &t
		}
	}
	if dt := ctx.Query("date_to"); dt != "" {
		if t, err := time.Parse("2006-01-02", dt); err == nil {
			dateTo = &t
		}
	}

	// Fetch events
	events, totalCount, err := as.adminRepo.GetAllEvents(page, limit, department, lecturerID, status, dateFrom, dateTo)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"message": "Failed to retrieve events",
		})
		return
	}

	// Calculate pagination
	totalPages := (totalCount + limit - 1) / limit

	ctx.JSON(200, gin.H{
		"success": true,
		"data": domain.EventListResponse{
			Events: events,
			Pagination: domain.PaginationInfo{
				CurrentPage:  page,
				TotalPages:   totalPages,
				TotalItems:   totalCount,
				ItemsPerPage: limit,
			},
		},
	})
}

// DeleteEvent deletes an event
func (as *AdminSvc) DeleteEvent(ctx *gin.Context) {
	eventID, err := strconv.ParseUint(ctx.Param("event_id"), 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "Invalid event ID",
		})
		return
	}

	err = as.adminRepo.DeleteEvent(uint(eventID))
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"message": "Failed to delete event",
		})
		return
	}

	// Get admin info from context
	adminID := ctx.GetInt("id")
	adminEmail := ctx.GetString("email")

	// Log the action
	resourceID := int(eventID)
	details := fmt.Sprintf("Deleted event with ID %d and all associated attendance records", eventID)
	as.LogAction("admin", adminID, adminEmail, "delete", "event", &resourceID, details, ctx.ClientIP(), ctx.GetHeader("User-Agent"))

	ctx.JSON(200, gin.H{
		"success": true,
		"message": "Event deleted successfully. All associated attendance records have been removed.",
	})
}

// GetAttendanceTrends retrieves attendance trends
func (as *AdminSvc) GetAttendanceTrends(ctx *gin.Context) {
	period := ctx.DefaultQuery("period", "weekly")

	// Default to last 30 days
	dateTo := time.Now()
	dateFrom := dateTo.AddDate(0, 0, -30)

	if df := ctx.Query("date_from"); df != "" {
		if t, err := time.Parse("2006-01-02", df); err == nil {
			dateFrom = t
		}
	}
	if dt := ctx.Query("date_to"); dt != "" {
		if t, err := time.Parse("2006-01-02", dt); err == nil {
			dateTo = t
		}
	}

	trends, err := as.adminRepo.GetAttendanceTrends(period, dateFrom, dateTo)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"message": "Failed to retrieve attendance trends",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"data": domain.TrendsResponse{
			Period: period,
			Trends: trends,
		},
	})
}

// GetLowAttendanceStudents retrieves students with low attendance
func (as *AdminSvc) GetLowAttendanceStudents(ctx *gin.Context) {
	threshold, _ := strconv.ParseFloat(ctx.DefaultQuery("threshold", "75"), 64)
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "50"))

	students, totalCount, err := as.adminRepo.GetLowAttendanceStudents(threshold, limit)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"message": "Failed to retrieve low attendance students",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"data": domain.LowAttendanceResponse{
			Threshold:      threshold,
			StudentsAtRisk: students,
			TotalAtRisk:    totalCount,
		},
	})
}

// GetSystemSettings retrieves system settings
func (as *AdminSvc) GetSystemSettings(ctx *gin.Context) {
	settings, err := as.adminRepo.GetAllSettings()
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"message": "Failed to retrieve system settings",
		})
		return
	}

	// Parse settings into response struct
	response := domain.SystemSettingsResponse{}

	if val, ok := settings["qr_code_validity_minutes"]; ok {
		if i, err := strconv.Atoi(val); err == nil {
			response.QRCodeValidityMinutes = i
		}
	}
	if val, ok := settings["attendance_grace_period_minutes"]; ok {
		if i, err := strconv.Atoi(val); err == nil {
			response.AttendanceGracePeriodMinutes = i
		}
	}
	if val, ok := settings["low_attendance_threshold"]; ok {
		if i, err := strconv.Atoi(val); err == nil {
			response.LowAttendanceThreshold = i
		}
	}
	if val, ok := settings["max_events_per_day_per_lecturer"]; ok {
		if i, err := strconv.Atoi(val); err == nil {
			response.MaxEventsPerDayPerLecturer = i
		}
	}
	if val, ok := settings["require_email_verification"]; ok {
		response.RequireEmailVerification = val == "true"
	}
	if val, ok := settings["allow_student_self_registration"]; ok {
		response.AllowStudentSelfRegistration = val == "true"
	}
	if val, ok := settings["academic_year"]; ok {
		response.AcademicYear = val
	}
	if val, ok := settings["semester"]; ok {
		response.Semester = val
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"data":    response,
	})
}

// UpdateSystemSettings updates system settings
func (as *AdminSvc) UpdateSystemSettings(ctx *gin.Context) {
	var req domain.UpdateSettingsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	// Get admin ID from context
	adminID := uint(ctx.GetInt("id"))
	updatedFields := []string{}

	// Update each setting if provided
	if req.QRCodeValidityMinutes != nil {
		if err := as.adminRepo.UpdateSetting("qr_code_validity_minutes", strconv.Itoa(*req.QRCodeValidityMinutes), adminID); err == nil {
			updatedFields = append(updatedFields, "qr_code_validity_minutes")
		}
	}
	if req.AttendanceGracePeriodMinutes != nil {
		if err := as.adminRepo.UpdateSetting("attendance_grace_period_minutes", strconv.Itoa(*req.AttendanceGracePeriodMinutes), adminID); err == nil {
			updatedFields = append(updatedFields, "attendance_grace_period_minutes")
		}
	}
	if req.LowAttendanceThreshold != nil {
		if err := as.adminRepo.UpdateSetting("low_attendance_threshold", strconv.Itoa(*req.LowAttendanceThreshold), adminID); err == nil {
			updatedFields = append(updatedFields, "low_attendance_threshold")
		}
	}
	if req.MaxEventsPerDayPerLecturer != nil {
		if err := as.adminRepo.UpdateSetting("max_events_per_day_per_lecturer", strconv.Itoa(*req.MaxEventsPerDayPerLecturer), adminID); err == nil {
			updatedFields = append(updatedFields, "max_events_per_day_per_lecturer")
		}
	}
	if req.RequireEmailVerification != nil {
		val := "false"
		if *req.RequireEmailVerification {
			val = "true"
		}
		if err := as.adminRepo.UpdateSetting("require_email_verification", val, adminID); err == nil {
			updatedFields = append(updatedFields, "require_email_verification")
		}
	}
	if req.AllowStudentSelfRegistration != nil {
		val := "false"
		if *req.AllowStudentSelfRegistration {
			val = "true"
		}
		if err := as.adminRepo.UpdateSetting("allow_student_self_registration", val, adminID); err == nil {
			updatedFields = append(updatedFields, "allow_student_self_registration")
		}
	}
	if req.AcademicYear != nil {
		if err := as.adminRepo.UpdateSetting("academic_year", *req.AcademicYear, adminID); err == nil {
			updatedFields = append(updatedFields, "academic_year")
		}
	}
	if req.Semester != nil {
		if err := as.adminRepo.UpdateSetting("semester", *req.Semester, adminID); err == nil {
			updatedFields = append(updatedFields, "semester")
		}
	}

	// Log the action
	adminEmail := ctx.GetString("email")
	details := fmt.Sprintf("Updated system settings: %v", updatedFields)
	as.LogAction("admin", int(adminID), adminEmail, "update", "settings", nil, details, ctx.ClientIP(), ctx.GetHeader("User-Agent"))

	ctx.JSON(200, gin.H{
		"success": true,
		"message": "Settings updated successfully",
		"data": gin.H{
			"updated_fields": updatedFields,
			"updated_at":     time.Now(),
		},
	})
}

// GetAuditLogs retrieves audit logs
func (as *AdminSvc) GetAuditLogs(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "50"))
	userType := ctx.Query("user_type")
	actionType := ctx.Query("action_type")

	var dateFrom, dateTo *time.Time
	if df := ctx.Query("date_from"); df != "" {
		if t, err := time.Parse("2006-01-02", df); err == nil {
			dateFrom = &t
		}
	}
	if dt := ctx.Query("date_to"); dt != "" {
		if t, err := time.Parse("2006-01-02", dt); err == nil {
			dateTo = &t
		}
	}

	logs, totalCount, err := as.adminRepo.GetAuditLogs(page, limit, userType, actionType, dateFrom, dateTo)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"message": "Failed to retrieve audit logs",
		})
		return
	}

	totalPages := (totalCount + limit - 1) / limit

	ctx.JSON(200, gin.H{
		"success": true,
		"data": domain.AuditLogResponse{
			Logs: logs,
			Pagination: domain.PaginationInfo{
				CurrentPage:  page,
				TotalPages:   totalPages,
				TotalItems:   totalCount,
				ItemsPerPage: limit,
			},
		},
	})
}

// LogAction is a helper function to log admin actions
func (as *AdminSvc) LogAction(userType string, userID int, userEmail, action, resourceType string, resourceID *int, details, ipAddress, userAgent string) error {
	log := &entities.AuditLog{
		Timestamp:    time.Now(),
		UserType:     userType,
		UserID:       userID,
		UserEmail:    userEmail,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Details:      details,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
	}

	return as.adminRepo.CreateAuditLog(log)
}

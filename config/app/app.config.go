package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Dom-HTG/attendance-management-system/config/database"
	adminRepo "github.com/Dom-HTG/attendance-management-system/internal/admin/repository"
	adminSvc "github.com/Dom-HTG/attendance-management-system/internal/admin/service"
	analyticsHandler "github.com/Dom-HTG/attendance-management-system/internal/analytics/handler"
	analyticsRepo "github.com/Dom-HTG/attendance-management-system/internal/analytics/repository"
	analyticsSvc "github.com/Dom-HTG/attendance-management-system/internal/analytics/service"
	attendanceRepo "github.com/Dom-HTG/attendance-management-system/internal/attendance/repository"
	attendanceSvc "github.com/Dom-HTG/attendance-management-system/internal/attendance/service"
	authRepo "github.com/Dom-HTG/attendance-management-system/internal/auth/repository"
	authSvc "github.com/Dom-HTG/attendance-management-system/internal/auth/service"
	"github.com/Dom-HTG/attendance-management-system/pkg/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Application struct {
	DB  database.DbConfig
	App AppConfig
}

type AppConfig struct {
	Port string
}

type Handlers struct {
	AuthHandler       *authSvc.AuthSvc
	AttendanceHandler *attendanceSvc.AttendanceSvc
	AnalyticsHandler  *analyticsHandler.AnalyticsHandler
	AdminHandler      adminSvc.AdminServiceInterface
}

// Mount method mounts the application routes and midddlewares to the gin engine.
func (app *Application) Mount(handler *Handlers) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(cors.New(buildCORSConfig()))

	// Auth routes.
	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/register-student", handler.AuthHandler.RegisterStudent)   // Registers new student.
		authRoutes.POST("/register-lecturer", handler.AuthHandler.RegisterLecturer) // Registers new lecturer.
		authRoutes.POST("/login-student", handler.AuthHandler.LoginStudent)         // Logs in student.
		authRoutes.POST("/login-lecturer", handler.AuthHandler.LoginLecturer)       // Logs in lecturer.
		authRoutes.POST("/login-admin", handler.AdminHandler.LoginAdmin)            // Logs in admin.
		authRoutes.POST("/forgot-password")                                         // Sends reset password email.
		authRoutes.POST("/logout")                                                  // Logs out user.
		authRoutes.POST("/refresh-token")                                           // Refresh access token..
	}

	// Student routes.
	studentRoutes := router.Group("/api/student")
	{
		studentRoutes.GET("/:id") // Retrieve student by id.
		studentRoutes.PUT("/:id") // Update student data by id.
	}

	// Lecturer routes.
	lecturerRoutes := router.Group("/api/lecturer")
	lecturerRoutes.Use(middleware.AuthMiddleware())
	lecturerRoutes.Use(middleware.RoleMiddleware("lecturer"))
	{
		lecturerRoutes.GET("/:id")                                                        // Retrieve lecturer by id.
		lecturerRoutes.PUT("/:id")                                                        // Update lecturer data by id.
		lecturerRoutes.POST("/qrcode/generate", handler.AttendanceHandler.GenerateQRCode) // Generate new QR Code.
	}

	// Events routes (NEW - Frontend requirement)
	eventsRoutes := router.Group("/api/events")
	eventsRoutes.Use(middleware.AuthMiddleware())
	{
		eventsRoutes.GET("/lecturer", middleware.RoleMiddleware("lecturer"), handler.AnalyticsHandler.GetLecturerEvents) // Get all lecturer events
	}

	// Attendance routes.
	attendanceRoutes := router.Group("/api/attendance")
	{
		attendanceRoutes.POST("/check-in", middleware.AuthMiddleware(), middleware.RoleMiddleware("student"), handler.AttendanceHandler.CheckIn)                    // Checks in user [marks user as present].
		attendanceRoutes.GET("/:event_id", middleware.AuthMiddleware(), middleware.RoleMiddleware("lecturer"), handler.AttendanceHandler.GetEventAttendance)        // Retrieves attendance record for an event.
		attendanceRoutes.GET("/student/records", middleware.AuthMiddleware(), middleware.RoleMiddleware("student"), handler.AttendanceHandler.GetStudentAttendance) // Retrieves student attendance history.
		attendanceRoutes.POST("/report")                                                                                                                            // Generates detailed attendance report for individual user.
	}

	// PDF Export routes
	router.GET("/api/student/attendance/export/pdf", middleware.AuthMiddleware(), middleware.RoleMiddleware("student"), handler.AttendanceHandler.ExportStudentAttendancePDFHandler)    // Export student attendance as PDF
	router.GET("/api/lecturer/attendance/export/pdf", middleware.AuthMiddleware(), middleware.RoleMiddleware("lecturer"), handler.AttendanceHandler.ExportLecturerAttendancePDFHandler) // Export lecturer attendance as PDF

	// Analytics routes - all require authentication
	analyticsRoutes := router.Group("/api/analytics")
	analyticsRoutes.Use(middleware.AuthMiddleware())
	{
		// Student analytics
		analyticsRoutes.GET("/student/:student_id", handler.AnalyticsHandler.GetStudentMetrics)           // Get student metrics
		analyticsRoutes.GET("/student/:student_id/insights", handler.AnalyticsHandler.GetStudentInsights) // Get student insights

		// Lecturer analytics (lecturer role required)
		lecturerAnalytics := analyticsRoutes.Group("")
		lecturerAnalytics.Use(middleware.RoleMiddleware("lecturer"))
		{
			// NEW - Frontend requirements
			lecturerAnalytics.GET("/lecturer/summary", handler.AnalyticsHandler.GetLecturerSummary) // Get lecturer summary

			// EXISTING
			lecturerAnalytics.GET("/lecturer/courses", handler.AnalyticsHandler.GetLecturerCourseMetrics)                 // Get lecturer course metrics
			lecturerAnalytics.GET("/lecturer/course/:course_code", handler.AnalyticsHandler.GetLecturerCoursePerformance) // Get course performance
			lecturerAnalytics.GET("/lecturer/insights", handler.AnalyticsHandler.GetLecturerInsights)                     // Get lecturer insights
		}

		// Admin analytics (admin/lecturer role required for now)
		adminAnalytics := analyticsRoutes.Group("")
		adminAnalytics.Use(middleware.RoleMiddleware("lecturer")) // Can be extended to "admin" role
		{
			// NEW - Frontend requirements
			adminAnalytics.GET("/admin/overview", handler.AnalyticsHandler.GetAdminOverviewNew)   // Get admin overview (NEW)
			adminAnalytics.GET("/admin/departments", handler.AnalyticsHandler.GetDepartmentStats) // Get department stats (NEW)

			// EXISTING
			adminAnalytics.GET("/admin/department/:department", handler.AnalyticsHandler.GetDepartmentMetrics) // Get department metrics
			adminAnalytics.GET("/admin/realtime", handler.AnalyticsHandler.GetRealTimeDashboard)               // Get real-time dashboard
		}

		// Temporal, anomaly, prediction, benchmarking, and chart endpoints (all authenticated users)
		analyticsRoutes.GET("/temporal", handler.AnalyticsHandler.GetTemporalAnalytics)                            // Get temporal analytics
		analyticsRoutes.GET("/anomalies", handler.AnalyticsHandler.DetectAnomalies)                                // Detect anomalies
		analyticsRoutes.GET("/predictions/student/:student_id", handler.AnalyticsHandler.PredictStudentAttendance) // Predict student attendance
		analyticsRoutes.GET("/predictions/course/:course_code", handler.AnalyticsHandler.PredictCourseAttendance)  // Predict course attendance
		analyticsRoutes.GET("/benchmark", handler.AnalyticsHandler.GetBenchmarkComparison)                         // Get benchmark comparison
		analyticsRoutes.GET("/charts/:chart_type", handler.AnalyticsHandler.GetChartData)                          // Get chart data
	}

	// Admin routes - all require admin role
	adminRoutes := router.Group("/api/admin")
	adminRoutes.Use(middleware.AuthMiddleware())
	adminRoutes.Use(middleware.RoleMiddleware("admin"))
	{
		// User Management - Students
		adminRoutes.GET("/students", handler.AdminHandler.GetAllStudents)                          // Get all students with pagination
		adminRoutes.GET("/users/student/:user_id", handler.AdminHandler.GetUserDetail)             // Get student details
		adminRoutes.PATCH("/users/student/:user_id/status", handler.AdminHandler.UpdateUserStatus) // Update student status
		adminRoutes.DELETE("/users/student/:user_id", handler.AdminHandler.DeleteUser)             // Delete student

		// User Management - Lecturers
		adminRoutes.GET("/lecturers", handler.AdminHandler.GetAllLecturers)                         // Get all lecturers with pagination
		adminRoutes.GET("/users/lecturer/:user_id", handler.AdminHandler.GetUserDetail)             // Get lecturer details
		adminRoutes.PATCH("/users/lecturer/:user_id/status", handler.AdminHandler.UpdateUserStatus) // Update lecturer status
		adminRoutes.DELETE("/users/lecturer/:user_id", handler.AdminHandler.DeleteUser)             // Delete lecturer

		// Event Management
		adminRoutes.GET("/events", handler.AdminHandler.GetAllEvents)             // Get all events with filtering
		adminRoutes.DELETE("/events/:event_id", handler.AdminHandler.DeleteEvent) // Delete event

		// Analytics
		adminRoutes.GET("/trends", handler.AdminHandler.GetAttendanceTrends)              // Get attendance trends
		adminRoutes.GET("/low-attendance", handler.AdminHandler.GetLowAttendanceStudents) // Get low attendance students

		// System Settings
		adminRoutes.GET("/settings", handler.AdminHandler.GetSystemSettings)      // Get system settings
		adminRoutes.PATCH("/settings", handler.AdminHandler.UpdateSystemSettings) // Update system settings

		// Audit Logs
		adminRoutes.GET("/audit-logs", handler.AdminHandler.GetAuditLogs) // Get audit logs
	}

	return router
}

func buildCORSConfig() cors.Config {
	cfg := cors.Config{
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}

	allowedOrigins := readEnvList("CORS_ALLOW_ORIGINS")
	allowedHeaders := readEnvList("CORS_ALLOW_HEADERS")
	exposedHeaders := readEnvList("CORS_EXPOSE_HEADERS")
	allowCredentials := parseBoolEnv(os.Getenv("CORS_ALLOW_CREDENTIALS"))

	if len(allowedHeaders) > 0 {
		cfg.AllowHeaders = allowedHeaders
	}

	if len(exposedHeaders) > 0 {
		cfg.ExposeHeaders = exposedHeaders
	}

	switch {
	case len(allowedOrigins) == 0:
		if allowCredentials {
			log.Println("CORS_ALLOW_CREDENTIALS=true but no explicit origins configured; disabling credential sharing")
			allowCredentials = false
		}
		cfg.AllowAllOrigins = true
	case len(allowedOrigins) == 1 && allowedOrigins[0] == "*":
		if allowCredentials {
			log.Println("CORS_ALLOW_ORIGINS contains '*', disabling credential sharing for safety")
			allowCredentials = false
		}
		cfg.AllowAllOrigins = true
	default:
		cfg.AllowOrigins = allowedOrigins
		for _, origin := range allowedOrigins {
			if strings.Contains(origin, "*") {
				cfg.AllowWildcard = true
				break
			}
		}
	}

	if maxAge := strings.TrimSpace(os.Getenv("CORS_MAX_AGE")); maxAge != "" {
		if duration, err := time.ParseDuration(maxAge); err == nil {
			cfg.MaxAge = duration
		} else {
			log.Printf("invalid CORS_MAX_AGE %q: %v", maxAge, err)
		}
	}

	cfg.AllowCredentials = allowCredentials
	return cfg
}

func readEnvList(key string) []string {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value != "" {
			result = append(result, value)
		}
	}
	return result
}

func parseBoolEnv(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "t", "yes", "y":
		return true
	}
	return false
}

func (app *Application) InjectDependencies(db *gorm.DB) *Handlers {
	// auth
	authRepoInstance := authRepo.NewAuthRepo(db)
	authSvcInstance := authSvc.NewAuthSvc(authRepoInstance)

	// attendance
	attendanceRepoInstance := attendanceRepo.NewAttendanceRepo(db)
	attendanceSvcInstance := attendanceSvc.NewAttendanceSvc(attendanceRepoInstance, authRepoInstance)

	// analytics
	analyticsRepoInstance := analyticsRepo.NewAnalyticsRepo(db)
	analyticsSvcInstance := analyticsSvc.NewAnalyticsService(analyticsRepoInstance)
	analyticsHandlerInstance := analyticsHandler.NewAnalyticsHandler(analyticsSvcInstance)

	// admin
	adminRepoInstance := adminRepo.NewAdminRepository(db)
	adminSvcInstance := adminSvc.NewAdminService(adminRepoInstance)

	return &Handlers{
		AuthHandler:       authSvcInstance,
		AttendanceHandler: attendanceSvcInstance,
		AnalyticsHandler:  analyticsHandlerInstance,
		AdminHandler:      adminSvcInstance,
	}
}

func (app *Application) Start(router *gin.Engine) error {

	port := os.Getenv("APP_PORT")
	// default to :2754 if not provided
	if port == "" {
		port = ":2754"
	}
	if err := router.Run(port); err != nil {
		return err
	}
	// trim leading ':' so we don't print 'http://localhost::2754'
	displayPort := strings.TrimPrefix(port, ":")
	fmt.Printf("Application started locally at http://localhost:%s\n", displayPort)

	return nil
}

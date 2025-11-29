package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/Dom-HTG/attendance-management-system/config/database"
	analyticsHandler "github.com/Dom-HTG/attendance-management-system/internal/analytics/handler"
	analyticsRepo "github.com/Dom-HTG/attendance-management-system/internal/analytics/repository"
	analyticsSvc "github.com/Dom-HTG/attendance-management-system/internal/analytics/service"
	attendanceRepo "github.com/Dom-HTG/attendance-management-system/internal/attendance/repository"
	attendanceSvc "github.com/Dom-HTG/attendance-management-system/internal/attendance/service"
	authRepo "github.com/Dom-HTG/attendance-management-system/internal/auth/repository"
	authSvc "github.com/Dom-HTG/attendance-management-system/internal/auth/service"
	"github.com/Dom-HTG/attendance-management-system/pkg/middleware"
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
}

// Mount method mounts the application routes and midddlewares to the gin engine.
func (app *Application) Mount(handler *Handlers) *gin.Engine {
	router := gin.Default()

	// Base middleware stack.
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Auth routes.
	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/register-student", handler.AuthHandler.RegisterStudent)   // Registers new student.
		authRoutes.POST("/register-lecturer", handler.AuthHandler.RegisterLecturer) // Registers new lecturer.
		authRoutes.POST("/login-student", handler.AuthHandler.LoginStudent)         // Logs in student.
		authRoutes.POST("/login-lecturer", handler.AuthHandler.LoginLecturer)       // Logs in lecturer.
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

	// Attendance routes.
	attendanceRoutes := router.Group("/api/attendance")
	{
		attendanceRoutes.POST("/check-in", middleware.AuthMiddleware(), middleware.RoleMiddleware("student"), handler.AttendanceHandler.CheckIn)                    // Checks in user [marks user as present].
		attendanceRoutes.GET("/:event_id", middleware.AuthMiddleware(), middleware.RoleMiddleware("lecturer"), handler.AttendanceHandler.GetEventAttendance)        // Retrieves attendance record for an event.
		attendanceRoutes.GET("/student/records", middleware.AuthMiddleware(), middleware.RoleMiddleware("student"), handler.AttendanceHandler.GetStudentAttendance) // Retrieves student attendance history.
		attendanceRoutes.POST("/report")                                                                                                                            // Generates detailed attendance report for individual user.
	}

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
			lecturerAnalytics.GET("/lecturer/courses", handler.AnalyticsHandler.GetLecturerCourseMetrics)                 // Get lecturer course metrics
			lecturerAnalytics.GET("/lecturer/course/:course_code", handler.AnalyticsHandler.GetLecturerCoursePerformance) // Get course performance
			lecturerAnalytics.GET("/lecturer/insights", handler.AnalyticsHandler.GetLecturerInsights)                     // Get lecturer insights
		}

		// Admin analytics (admin/lecturer role required for now)
		adminAnalytics := analyticsRoutes.Group("")
		adminAnalytics.Use(middleware.RoleMiddleware("lecturer")) // Can be extended to "admin" role
		{
			adminAnalytics.GET("/admin/overview", handler.AnalyticsHandler.GetAdminOverview)                   // Get admin overview
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

	return router
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

	return &Handlers{
		AuthHandler:       authSvcInstance,
		AttendanceHandler: attendanceSvcInstance,
		AnalyticsHandler:  analyticsHandlerInstance,
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

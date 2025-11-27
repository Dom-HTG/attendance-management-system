package config

import (
	"fmt"
	"os"

	"github.com/Dom-HTG/attendance-management-system/config/database"
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

	return router
}

func (app *Application) InjectDependencies(db *gorm.DB) *Handlers {
	// auth
	authRepoInstance := authRepo.NewAuthRepo(db)
	authSvcInstance := authSvc.NewAuthSvc(authRepoInstance)

	// attendance
	attendanceRepoInstance := attendanceRepo.NewAttendanceRepo(db)
	attendanceSvcInstance := attendanceSvc.NewAttendanceSvc(attendanceRepoInstance, authRepoInstance)

	return &Handlers{
		AuthHandler:       authSvcInstance,
		AttendanceHandler: attendanceSvcInstance,
	}
}

func (app *Application) Start(router *gin.Engine) error {

	port := os.Getenv("APP_PORT")
	if err := router.Run(port); err != nil {
		return err
	}
	fmt.Printf("Application started locally at http://localhost:%d\n", port)

	return nil
}

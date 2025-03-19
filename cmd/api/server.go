package main

import (
	"fmt"
	"os"

	"github.com/Dom-HTG/attendance-management-system/database"
	"github.com/gin-gonic/gin"
)

type Application struct {
	db  database.DbConfig
	app appConfig
}

type appConfig struct {
	port string
}

// Mount method mounts the application routes and midddlewares to the gin engine.
func (app *Application) Mount() *gin.Engine {
	router := gin.Default()

	// Base middleware stack.
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Auth routes.
	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/register")        // Registers new user.
		authRoutes.POST("/login")           // Logs in user.
		authRoutes.POST("/forgot-password") // Sends reset password email.
		authRoutes.POST("/logout")          // Logs out user.
		authRoutes.POST("/refresh-token")   // Refresh access token..
	}

	// User routes.
	userRoutes := router.Group("/api/user")
	{
		userRoutes.GET("/{id}") // Retrieve user by id.
		userRoutes.PUT("/{id}") // Update user by id.
	}

	// Attendance routes.
	attendanceRoutes := router.Group("/api/attendance")
	{
		attendanceRoutes.GET("/")          // Retrieves attendance record.
		attendanceRoutes.POST("/check-in") // Checks in user [marks user as present].
		attendanceRoutes.POST("/report")   // Generates detailed attendance report for individual user.
	}

	return router
}

func (app *Application) Start(router *gin.Engine) error {

	port := os.Getenv("APP_PORT")
	if err := router.Run(port); err != nil {
		return err
	}
	fmt.Printf("Application started locally at http://localhost:%d\n", port)

	return nil
}

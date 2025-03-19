package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

type Application struct {
	db  dbConfig
	app appConfig
}

type dbConfig struct {
	dsn           string
	maxOpenConns  int
	maxIdleConns  int
	maxIdleTimout string
}

type appConfig struct {
	port string
}

// Mount method mounts the application routes to the gin engine.
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

func Start(router *gin.Engine) (*Application, error) {
	app := &Application{
		db: dbConfig{
			dsn:           os.Getenv("DATABASE_DSN"),
			maxOpenConns:  10,
			maxIdleConns:  5,
			maxIdleTimout: "1m",
		},
		app: appConfig{
			port: "8080",
		},
	}

	port := os.Getenv("APP_PORT")
	if err := router.Run(port); err != nil {
		return nil, err
	}
	fmt.Printf("Application started locally at http://localhost:%d\n", port)

	return app, nil
}

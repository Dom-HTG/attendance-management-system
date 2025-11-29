package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	config "github.com/Dom-HTG/attendance-management-system/config/app"
	"github.com/Dom-HTG/attendance-management-system/config/database"
	"github.com/Dom-HTG/attendance-management-system/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	// initialize structured logger (JSON) writing to logs/app.log + stdout
	logger.Init("logs/app.log", logger.LogrusLevel())

	// Load environment variables from .env file if running locally (not in Docker).
	// In Docker, env vars are provided via docker-compose.yml or --env-file.
	if _, err := os.Stat("cmd/api/app.env"); err == nil {
		if err := godotenv.Load("cmd/api/app.env"); err != nil {
			log.Printf("Warning: could not load .env file: %v\n", err)
		}
	}

	// Build DSN string.
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	openConn, _ := strconv.Atoi(os.Getenv("POOL_MAX_OPEN_CONN"))
	idleConn, _ := strconv.Atoi(os.Getenv("POOL_MAX_IDLE_CONN"))
	connTimeout := os.Getenv("POOL_MAX_CONN_TIMEOUT")

	// this configuration will be passed into the handlers.
	app := &config.Application{
		DB: database.DbConfig{
			DSN:           dsn,
			MaxOpenConns:  openConn,
			MaxIdleConns:  idleConn,
			MaxIdleTimout: connTimeout,
		},
		App: config.AppConfig{
			Port: os.Getenv("APP_PORT"),
		},
	}

	// Start the database connection.
	db, er := app.DB.Start()
	if er != nil {
		logger.Errorf("Error starting the database: %v", er)
		log.Fatalf("Error starting the database: %v", er)
	}

	// Inject dependencies.
	handler := app.InjectDependencies(db)

	// Register routes and middlewares.
	router := app.Mount(handler)

	// Start HTTP server using net/http.Server so we can shut it down gracefully.
	addr := app.App.Port
	if addr == "" {
		addr = ":2754"
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Start server in a goroutine.
	go func() {
		logger.Infof("Application started locally at http://localhost:%s", addrTrim(addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("HTTP server ListenAndServe error: %v", err)
			log.Fatalf("HTTP server ListenAndServe error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutdown signal received, shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
	} else {
		logger.Info("Server stopped gracefully")
	}

	// Close the underlying database connection.
	if sqlDB, err := db.DB(); err == nil {
		if cerr := sqlDB.Close(); cerr != nil {
			logger.Errorf("Error closing database: %v", cerr)
		} else {
			logger.Info("Database connection closed")
		}
	}
}

// addrTrim removes a leading ':' from an address string for cleaner printing.
func addrTrim(a string) string {
	if len(a) > 0 && a[0] == ':' {
		return a[1:]
	}
	return a
}

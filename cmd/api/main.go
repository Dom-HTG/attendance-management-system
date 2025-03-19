package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Dom-HTG/attendance-management-system/config/app"
	"github.com/Dom-HTG/attendance-management-system/config/database"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables.
	err := godotenv.Load("cmd/api/app.env")
	if err != nil {
		fmt.Errorf(err.Error())
		panic("Error loading environment variables..")
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
	_, er := app.DB.Start()
	if er != nil {
		log.Fatalf("Error starting the database: %v", er)
	}

	// Register routes and middlewares.
	router := app.Mount()

	// Start server.
	if errr := app.Start(router); errr != nil {
		log.Fatalf("Error starting the server: %v", errr)

	}
}

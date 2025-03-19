package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Dom-HTG/attendance-management-system/database"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables.
	err := godotenv.Load("app.env")
	if err != nil {
		fmt.Errorf(err.Error())
		panic("Error loading environment variables..")
	}

	// this configuration will be passed into the handlers.
	app := &Application{
		db: database.DbConfig{
			DSN:           os.Getenv("DATABASE_DSN"),
			MaxOpenConns:  10,
			MaxIdleConns:  5,
			MaxIdleTimout: "1m",
		},
		app: appConfig{
			port: ":8080",
		},
	}

	// Start the database connection.
	_, er := app.db.Start()
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

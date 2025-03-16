package main

import (
	"fmt"
	"os"

	"github.com/Dom-HTG/attendance-management-system/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables.
	err := godotenv.Load("app.env")
	if err != nil {
		panic("Error loading environment variables..")
	}

	// Start gin server.
	router := gin.Default()

	// Register routes.
	routes.RegisterRoutes(router)

}

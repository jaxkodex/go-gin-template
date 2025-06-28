package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jaxkodex/go-gin-template/src/database"
	"github.com/jaxkodex/go-gin-template/src/routes"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to PostgreSQL
	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}
	defer database.DB.Close(context.Background())

	// Initialize Gin router
	r := gin.Default()

	// Register all routes
	routes.RegisterRoutes(r)

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8000"
	}
	log.Fatal(r.Run(":" + port))
}

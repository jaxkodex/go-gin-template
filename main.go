package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jaxkodex/go-gin-template/src/api"
	"github.com/jaxkodex/go-gin-template/src/config"
	"github.com/jaxkodex/go-gin-template/src/database"
	"github.com/jaxkodex/go-gin-template/src/middleware"
	"github.com/jaxkodex/go-gin-template/src/routes"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file if present.
	// A missing file is not fatal — production containers rely on real env vars.
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, relying on environment variables")
	}

	// Load typed application config.
	cfg := config.Load()

	// Connect to PostgreSQL.
	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}
	defer database.DB.Close()

	// Initialize authentication middleware.
	authMiddleware, err := middleware.NewAuth(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize auth middleware: %v", err)
	}

	// Initialize Gin router.
	r := gin.Default()

	// Register all routes.
	routes.RegisterRoutes(r, authMiddleware)

	// Register generated OpenAPI routes.
	api.RegisterHandlers(r, api.NewServer())

	// Start server.
	log.Fatal(r.Run(":" + cfg.ServerPort))
}

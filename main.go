package main

import (
	"log"
	"os"

	"go-postgres-api/internal/config"
	"go-postgres-api/internal/database"
	"go-postgres-api/internal/middleware"
	"go-postgres-api/internal/models"
	"go-postgres-api/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables with explicit path
	if err := godotenv.Load("d:\\Projects\\go-postgres-api\\.env"); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Initialize configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate models
	err = db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.AuthLog{},
		&models.TokenBlacklist{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Get the underlying SQL DB to set up connection pool parameters
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB: %v", err)
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// Initialize Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(middleware.CORSMiddleware())

	// Set up routes
	routes.SetupRoutes(router, cfg)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on http://localhost:%s", port)
	log.Fatal(router.Run(":" + port))
}

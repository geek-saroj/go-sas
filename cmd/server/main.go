package main

import (
	"log"

	"sas-pro/config"
	"sas-pro/internal/routes"
	"sas-pro/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	if err := database.Connect(cfg.DSN()); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Auto migrate models
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}


	// Create Gin router
	router := gin.Default()

	// Setup routes
	routes.Setup(router)

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
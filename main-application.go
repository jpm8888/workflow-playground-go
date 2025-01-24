package main

import (
	"app/auth/internal/handlers"
	"app/auth/internal/repository"
	"app/auth/internal/workflow"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Database connection
	dsn := "host=localhost user=postgres password=postgres dbname=auth port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate database schema
	err = db.AutoMigrate(&workflow.Workflow{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Create repository
	repo := repository.NewWorkflowRepository(db)

	// Setup Gin router
	r := gin.Default()

	// Setup workflow routes
	handlers.SetupRoutes(r, repo)

	// Start the server
	if err := r.Run(":9012"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

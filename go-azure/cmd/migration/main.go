package main

import (
	"go-azure/config"
	"go-azure/migrations"
	"go-azure/utils"
	"os"
)

func main() {
	// Initialize logger
	utils.InitLogger()
	logger := utils.GetLogger()
	logger.Info("Starting application")

	// Load configuration
	cfg := config.LoadConfig()
	
	// Initialize database
	db, err := utils.InitDatabase(cfg)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize database")
	}

	// Run migrations
	if err := migrations.Migrate(db); err != nil {
		logger.WithError(err).Fatal("Failed to run migrations")
	}

	// Seed database (only in development environment)
	if os.Getenv("APP_ENV") != "production" {
		if err := migrations.Seed(db); err != nil {
			logger.WithError(err).Fatal("Failed to seed database")
		}
	}
}

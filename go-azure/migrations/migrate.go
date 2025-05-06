package migrations

import (
	"go-azure/models"
	"log"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Migrate runs database migrations
func Migrate(db *gorm.DB) error {
	logrus.Info("Running database migrations")

	// Check if the database exists
	dbName := "go_azure"
	// err := db.Exec("USE " + dbName).Error
	err2 := db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName).Error
	if err2 != nil {
		logrus.WithError(err2).Error("Failed to create database")
	}
	log.Printf("Database %s is ready.", dbName)
	// logrus.Info("Database %s is ready.", dbName)

	// Auto migrate models
	err := db.AutoMigrate(
		&models.User{},
		&models.Task{},
		&models.SocialMediaPost{},
		&models.SocialMediaComments{},
		&models.SocialMediaLikes{},
	)
	if err != nil {
		logrus.WithError(err).Error("Failed to run migrations")
		return err
	}

	logrus.Info("Database migrations completed successfully")
	return nil
}

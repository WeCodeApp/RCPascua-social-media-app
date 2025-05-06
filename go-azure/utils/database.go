package utils

import (
	"fmt"
	"time"

	"go-azure/config"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

// InitDatabase initializes the database connection
func InitDatabase(cfg *config.Config) (*gorm.DB, error) {
	// Create DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	// Configure GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Connect to database
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		logrus.WithError(err).Error("Failed to connect to database")
		return nil, err
	}

	// Configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		logrus.WithError(err).Error("Failed to get database connection")
		return nil, err
	}

	// Ping the database to ensure connectivity
	if err := sqlDB.Ping(); err != nil {
		logrus.WithError(err).Error("Failed to ping database")
		return nil, err
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logrus.Info("Database connection established")
	return DB, nil
}

// GetDB returns the database connection
func GetDB() *gorm.DB {
	return DB
}

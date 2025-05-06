package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

// InitLogger initializes the logger
func InitLogger() {
	// Set logger output format
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Set output to stdout
	logrus.SetOutput(os.Stdout)

	// Set log level
	logrus.SetLevel(logrus.InfoLevel)

	logrus.Info("Logger initialized")
}

// GetLogger returns a logger instance
func GetLogger() *logrus.Logger {
	return logrus.StandardLogger()
}
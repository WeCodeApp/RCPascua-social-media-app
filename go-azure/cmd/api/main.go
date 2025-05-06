package main

import (
	"go-azure/config"
	"go-azure/controllers"
	"go-azure/middleware"
	"go-azure/services"
	"go-azure/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize logger
	utils.InitLogger()
	logger := utils.GetLogger()
	logger.Info("Starting application")

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	_, err := utils.InitDatabase(cfg)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize database")
	}

	// Initialize services
	authService := services.NewAuthService(cfg)
	taskService := services.NewTaskService()
	socialMediaService := services.NewSocialMediaService()

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Initialize controllers
	authController := controllers.NewAuthController(authService, cfg)
	taskController := controllers.NewTaskController(taskService, authMiddleware)
	socialMediaController := controllers.NewSocialMediaController(socialMediaService, authMiddleware)

	// Initialize router
	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", cfg.AppURL)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Add middleware for logging
	router.Use(func(c *gin.Context) {
		// Log request
		logger.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		}).Info("Request received")

		c.Next()

		// Log response
		logger.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"latency_ms": c.Writer.Size(),
		}).Info("Response sent")
	})

	// Register routes
	authController.RegisterRoutes(router)
	taskController.RegisterRoutes(router)
	socialMediaController.RegisterRoutes(router)

	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Start server
	logger.WithFields(logrus.Fields{
		"port": cfg.Port,
	}).Info("Server starting")

	if err := router.Run(cfg.Host + ":" + cfg.Port); err != nil {
		logger.WithError(err).Fatal("Failed to start server")
	}

	// Use:
	//if err := router.RunTLS("192.168.5.143:"+cfg.Port, "cert.pem", "key.pem"); err != nil {
	//	logger.WithError(err).Fatal("Failed to start HTTPS server")
	//}

}

package middleware

import (
	"net/http"
	"strings"

	"go-azure/services"
	"go-azure/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AuthMiddleware is a middleware for JWT authentication
type AuthMiddleware struct {
	authService *services.AuthService
	logger      *logrus.Logger
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(authService *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		logger:      utils.GetLogger(),
	}
}

// RequireAuth is a middleware that requires JWT authentication
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			m.logger.Warn("Missing authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		// Check if the header has the Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			m.logger.Warn("Invalid authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token
		claims, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			m.logger.WithError(err).Warn("Invalid token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Set user ID in context
		userID, ok := claims["user_id"].(string)
		if !ok {
			m.logger.Warn("User ID not found in token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Set user info in context; this is now the session variables
		c.Set("user_id", userID)
		c.Set("email", claims["email"])
		c.Set("name", claims["name"])

		m.logger.WithFields(logrus.Fields{
			"user_id": userID,
		}).Info("User authenticated")

		c.Next()
	}
}

package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"go-azure/config"
	"go-azure/services"
	"go-azure/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AuthController handles authentication endpoints
type AuthController struct {
	authService *services.AuthService
	logger      *logrus.Logger
	config      *config.Config
}

// NewAuthController creates a new AuthController
func NewAuthController(authService *services.AuthService, config *config.Config) *AuthController {
	return &AuthController{
		authService: authService,
		logger:      utils.GetLogger(),
		config:      config,
	}
}

// RegisterRoutes registers the routes for the AuthController
func (c *AuthController) RegisterRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.GET("/microsoft", c.MicrosoftLogin)
		auth.GET("/microsoft/callback", c.MicrosoftCallback)
		auth.POST("/signout", c.SignOut)
		auth.GET("/me", c.GetCurrentUser)
	}
}

// MicrosoftLogin returns Microsoft OAuth login URL
func (c *AuthController) MicrosoftLogin(ctx *gin.Context) {
	// Generate state for CSRF protection
	state, err := c.authService.GenerateState()
	if err != nil {
		c.logger.WithError(err).Error("Failed to generate state")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate login"})
		return
	}

	// Store state in cookie
	ctx.SetCookie("oauth_state", state, 3600, "/", "", false, true)

	// Get Microsoft login URL
	loginURL := c.authService.GetMicrosoftLoginURL(state)

	// Return the login URL as JSON
	ctx.JSON(http.StatusOK, gin.H{"login_url": loginURL})
}

// MicrosoftCallback handles the callback from Microsoft OAuth
func (c *AuthController) MicrosoftCallback(ctx *gin.Context) {
	////Get state from cookie
	//stateCookie, err := ctx.Cookie("oauth_state")
	//if err != nil {
	//	c.logger.WithError(err).Error("Failed to get state cookie")
	//	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
	//	return
	//}
	//
	////Verify state
	//state := ctx.Query("state")
	//if state != stateCookie {
	//	c.logger.Error("State mismatch")
	//	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
	//	return
	//}

	// Get code
	code := ctx.Query("code")
	if code == "" {
		c.logger.Error("Code not found")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}

	// Exchange code for token
	tokenDetails, user, err := c.authService.HandleMicrosoftCallback(code)
	if err != nil {
		c.logger.WithError(err).Error("Failed to handle Microsoft callback")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate"})
		return
	}

	// Clear state cookie
	//ctx.SetCookie("oauth_state", "", -1, "/", "", false, true)

	// Log successful login
	c.logger.WithFields(logrus.Fields{
		"user_id": user.ID,
		"name":    user.Name,
		"email":   user.Email,
	}).Info("User logged in")

	// Convert token and user to JSON
	tokenJSON, err := json.Marshal(tokenDetails)
	if err != nil {
		c.logger.WithError(err).Error("Failed to marshal token")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process authentication"})
		return
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		c.logger.WithError(err).Error("Failed to marshal user")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process authentication"})
		return
	}

	// Build redirect URL with token and user data as query parameters
	redirectURL, err := url.Parse(fmt.Sprintf("%s/login", c.config.AppURL))
	if err != nil {
		c.logger.WithError(err).Error("Failed to parse frontend URL")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process authentication"})
		return
	}

	query := redirectURL.Query()
	query.Set("token", string(tokenJSON))
	query.Set("user", string(userJSON))
	redirectURL.RawQuery = query.Encode()

	// Redirect to frontend application with token and user data
	ctx.Redirect(http.StatusTemporaryRedirect, redirectURL.String())
}

// SignOut handles user sign out
func (c *AuthController) SignOut(ctx *gin.Context) {
	// In a stateless JWT authentication system, the client is responsible for
	// discarding the token. The server doesn't need to do anything special.
	// However, in a real application, you might want to blacklist the token.

	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("user_id")
	if exists {
		c.logger.WithFields(logrus.Fields{
			"user_id": userID,
		}).Info("User signed out")
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully signed out"})
}

func (c *AuthController) GetCurrentUser(ctx *gin.Context) {
	// Get user info from context (set by auth middleware)
	userID := ctx.GetString("user_id")
	name := ctx.GetString("name")
	email := ctx.GetString("email")

	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Log successful login
	c.logger.WithFields(logrus.Fields{
		"user_id": userID,
		"name":    name,
		"email":   email,
	}).Info("User logged in")

	ctx.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"name":    name,
		"email":   email,
	})
}

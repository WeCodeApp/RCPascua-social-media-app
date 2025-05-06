package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"go-azure/config"
	"go-azure/models"
	"go-azure/utils"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
	"gorm.io/gorm"
)

// AuthService handles authentication operations
type AuthService struct {
	config *config.Config
	logger *logrus.Logger
	db     *gorm.DB
}

// NewAuthService creates a new AuthService
func NewAuthService(config *config.Config) *AuthService {
	return &AuthService{
		config: config,
		logger: utils.GetLogger(),
		db:     utils.GetDB(),
	}
}

// GetMicrosoftOAuthConfig returns the OAuth2 config for Microsoft
func (s *AuthService) GetMicrosoftOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     s.config.MicrosoftClientID,
		ClientSecret: s.config.MicrosoftClientSecret,
		RedirectURL:  s.config.MicrosoftRedirectURI,
		Scopes:       []string{"openid", "profile", "email", "offline_access", "User.Read"},
		Endpoint:     microsoft.AzureADEndpoint(s.config.MicrosoftTenantID),
	}
}

// GenerateState generates a random state string for OAuth
func (s *AuthService) GenerateState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// GetMicrosoftLoginURL returns the URL for Microsoft login
func (s *AuthService) GetMicrosoftLoginURL(state string) string {
	return s.GetMicrosoftOAuthConfig().AuthCodeURL(state)
}

// HandleMicrosoftCallback handles the callback from Microsoft OAuth
func (s *AuthService) HandleMicrosoftCallback(code string) (*models.TokenDetails, *models.User, error) {
	// Exchange code for token
	oauth2Config := s.GetMicrosoftOAuthConfig()
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		s.logger.WithError(err).Error("Failed to exchange code for token")
		return nil, nil, err
	}

	// Get user info
	userInfo, err := s.getUserInfo(token.AccessToken)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get user info")
		return nil, nil, err
	}

	// Check if user exists in database
	var user models.User

	result := s.db.Where("email = ?", userInfo["userPrincipalName"].(string)).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Create new user
			user = models.User{
				ID:        uuid.New().String(),
				Email:     userInfo["userPrincipalName"].(string),
				Name:      userInfo["displayName"].(string),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			// Save user to database
			if err := s.db.Create(&user).Error; err != nil {
				s.logger.WithError(err).Error("Failed to create user")
				return nil, nil, errors.New("failed to create user")
			}

			s.logger.WithFields(logrus.Fields{
				"user_id": user.ID,
				"email":   user.Email,
			}).Info("New user created")
		} else {
			s.logger.WithError(result.Error).Error("Failed to query user")
			return nil, nil, errors.New("failed to query user")
		}
	} else {
		// Update user information
		user.Name = userInfo["displayName"].(string)
		user.UpdatedAt = time.Now()

		if err := s.db.Save(&user).Error; err != nil {
			s.logger.WithError(err).Error("Failed to update user")
			return nil, nil, errors.New("failed to update user")
		}

		s.logger.WithFields(logrus.Fields{
			"user_id": user.ID,
			"name":    user.Name,
			"email":   user.Email,
		}).Info("Existing user updated")
	}

	// Generate JWT token
	tokenDetails, err := utils.GenerateToken(user.ID, user.Email, user.Name, s.config.JWTSecret, s.config.JWTExpirationMinutes)
	if err != nil {
		s.logger.WithError(err).Error("Failed to generate JWT token")
		return nil, nil, err
	}

	return tokenDetails, &user, nil
}

// getUserInfo gets user information from Microsoft Graph API
func (s *AuthService) getUserInfo(accessToken string) (map[string]interface{}, error) {
	// Create request
	req, err := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/me", nil)
	if err != nil {
		return nil, err
	}

	// Add authorization header
	req.Header.Add("Authorization", "Bearer "+accessToken)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get user info: " + resp.Status)
	}

	// Parse response
	var userInfo map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

// ValidateToken validates a JWT token
func (s *AuthService) ValidateToken(tokenString string) (map[string]interface{}, error) {
	claims, err := utils.ValidateToken(tokenString, s.config.JWTSecret)
	if err != nil {
		return nil, err
	}

	// Convert claims to map
	result := make(map[string]interface{})
	for key, value := range claims {
		result[key] = value
	}

	return result, nil
}

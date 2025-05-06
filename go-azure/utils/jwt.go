package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go-azure/models"
)

// Custom claims struct
type CustomClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

// GenerateToken generates a new JWT token for a user
func GenerateToken(userID string, email string, name string, secret string, expirationMinutes int) (*models.TokenDetails, error) {
	// Create token details
	expiresAt := time.Now().Add(time.Minute * time.Duration(expirationMinutes))
	td := &models.TokenDetails{
		TokenType:    "Bearer",
		ExpiresIn:    int64(expirationMinutes * 60), // Convert to seconds
		ExpiresAt:    expiresAt,
		RefreshToken: uuid.New().String(),
	}

	// Create claims with registered claims for better security
	claims := CustomClaims{
		UserID: userID,
		Email:  email,
		Name:   name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-azure-api",
			Subject:   userID,
			ID:        uuid.New().String(),
			Audience:  []string{"go-azure-api-users"},
		},
	}

	// Create token with custom claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	var err error
	td.AccessToken, err = token.SignedString([]byte(secret))
	if err != nil {
		logrus.WithError(err).Error("Failed to sign JWT token")
		return nil, err
	}

	return td, nil
}

// ValidateToken validates a JWT token
func ValidateToken(tokenString string, secret string) (jwt.MapClaims, error) {
	// Parse token with custom claims
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		logrus.WithError(err).Error("Failed to parse JWT token")
		return nil, err
	}

	// Validate token
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Get claims
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Convert to map for compatibility with existing code
	claimsMap := jwt.MapClaims{
		"user_id": claims.UserID,
		"email":   claims.Email,
		"name":    claims.Name,
		"exp":     claims.ExpiresAt.Time.Unix(),
		"iat":     claims.IssuedAt.Time.Unix(),
		"sub":     claims.Subject,
		"jti":     claims.ID,
		"iss":     claims.Issuer,
		"aud":     claims.Audience,
	}

	return claimsMap, nil
}

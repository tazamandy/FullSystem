// services/jwt_service.go
package services

import (
	"attendance-system/models"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Import types from models for convenience
type JWTClaims = models.JWTClaims
type RefreshTokenClaims = models.RefreshTokenClaims

var (
	ErrNoSecretKey     = models.ErrNoSecretKey
	ErrInvalidToken    = models.ErrInvalidToken
	AccessTokenExpiry  = models.AccessTokenExpiry
	RefreshTokenExpiry = models.RefreshTokenExpiry
)

var jwtSecret string

// init loads JWT secret from environment
func init() {
	jwtSecret = os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-change-in-production" // Fallback for development only
	}
}

// GenerateAccessToken creates a new JWT access token
func GenerateAccessToken(user models.User) (string, error) {
	if jwtSecret == "" {
		return "", ErrNoSecretKey
	}

	now := time.Now().UTC()
	claims := JWTClaims{
		StudentID: user.StudentID,
		Email:     user.Email,
		Role:      user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   user.StudentID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// GenerateRefreshToken creates a new JWT refresh token
func GenerateRefreshToken(user models.User) (string, error) {
	if jwtSecret == "" {
		return "", ErrNoSecretKey
	}

	now := time.Now().UTC()
	claims := RefreshTokenClaims{
		StudentID: user.StudentID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(RefreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   user.StudentID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, nil
}

// VerifyAccessToken verifies and parses an access token
func VerifyAccessToken(tokenString string) (*JWTClaims, error) {
	if tokenString == "" {
		return nil, ErrInvalidToken
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// VerifyRefreshToken verifies and parses a refresh token
func VerifyRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	if tokenString == "" {
		return nil, ErrInvalidToken
	}

	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GeneratePasswordResetToken creates a JWT token for password reset verification
func GeneratePasswordResetToken(email, code string) (string, error) {
	if jwtSecret == "" {
		return "", ErrNoSecretKey
	}

	claims := models.PasswordResetTokenClaims{
		Email: email,
		Code:  code,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(models.PasswordResetTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign password reset token: %w", err)
	}

	return tokenString, nil
}

// VerifyPasswordResetToken verifies and parses a password reset token
func VerifyPasswordResetToken(tokenString string) (*models.PasswordResetTokenClaims, error) {
	if tokenString == "" {
		return nil, ErrInvalidToken
	}

	token, err := jwt.ParseWithClaims(tokenString, &models.PasswordResetTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	}, jwt.WithLeeway(1*time.Minute))

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*models.PasswordResetTokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GenerateEmailVerificationToken creates a JWT token for email verification
func GenerateEmailVerificationToken(email, code string) (string, error) {
	if jwtSecret == "" {
		return "", ErrNoSecretKey
	}

	now := time.Now().UTC()
	claims := models.EmailVerificationTokenClaims{
		Email: email,
		Code:  code,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(models.EmailVerificationTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign email verification token: %w", err)
	}

	return tokenString, nil
}

// VerifyEmailVerificationToken verifies and parses an email verification token
func VerifyEmailVerificationToken(tokenString string) (*models.EmailVerificationTokenClaims, error) {
	if tokenString == "" {
		return nil, ErrInvalidToken
	}

	token, err := jwt.ParseWithClaims(tokenString, &models.EmailVerificationTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	}, jwt.WithLeeway(1*time.Minute))

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*models.EmailVerificationTokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

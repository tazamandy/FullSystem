// models/jwt.go
// JWT token claims and configuration

package models

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims contains JWT token claims with expiration
type JWTClaims struct {
	StudentID string `json:"student_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	jwt.RegisteredClaims
}

// RefreshTokenClaims contains refresh token claims
type RefreshTokenClaims struct {
	StudentID string `json:"student_id"`
	jwt.RegisteredClaims
}

// PasswordResetTokenClaims contains password reset token claims
type PasswordResetTokenClaims struct {
	Email string `json:"email"`
	Code  string `json:"code"`
	jwt.RegisteredClaims
}

// EmailVerificationTokenClaims contains email verification token claims
type EmailVerificationTokenClaims struct {
	Email string `json:"email"`
	Code  string `json:"code"`
	jwt.RegisteredClaims
}

const (
	// AccessTokenExpiry is the expiration time for access tokens (15 minutes)
	AccessTokenExpiry = 15 * time.Minute
	// RefreshTokenExpiry is the expiration time for refresh tokens (7 days)
	RefreshTokenExpiry = 7 * 24 * time.Hour
	// PasswordResetTokenExpiry is the expiration time for password reset tokens (15 minutes)
	PasswordResetTokenExpiry = 15 * time.Minute
	// EmailVerificationTokenExpiry is the expiration time for email verification tokens (15 minutes)
	EmailVerificationTokenExpiry = 15 * time.Minute
)

var (
	// JWT Errors
	ErrInvalidToken = errors.New("invalid or expired token")
	ErrNoSecretKey  = errors.New("JWT_SECRET not configured")
)

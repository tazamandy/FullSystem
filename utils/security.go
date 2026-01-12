// utils/security.go
package utils

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	// BcryptCost is the cost factor for password hashing
	BcryptCost = 12
	// VerificationCodeLength is the length of verification codes
	VerificationCodeLength = 6
)

// HashPassword securely hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashed), nil
}

// ComparePassword compares a password with a hash
func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// GenerateSecureCode generates a cryptographically secure random code
func GenerateSecureCode(length int) (string, error) {
	if length <= 0 || length > 10 {
		return "", errors.New("code length must be between 1 and 10")
	}

	// Calculate the maximum value for the given length
	maxValue := int64(1)
	for i := 0; i < length; i++ {
		maxValue *= 10
	}

	// Generate random bytes
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random code: %w", err)
	}

	// Convert to uint64 and take modulo safely (avoid signed overflow)
	randomValue := binary.BigEndian.Uint64(b)
	codeNum := randomValue % uint64(maxValue)

	// Format with leading zeros
	format := fmt.Sprintf("%%0%dd", length)
	return fmt.Sprintf(format, int(codeNum)), nil
}

// GenerateVerificationCode generates a 6-digit verification code or returns error
func GenerateVerificationCode() (string, error) {
	return GenerateSecureCode(VerificationCodeLength)
}

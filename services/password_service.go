// services/password_service.go
package services

import (
	"attendance-system/connection"
	"attendance-system/models"
	"attendance-system/utils"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Constants
const (
	emailQuery = "email = ?"
	codeQuery  = "email = ? AND code = ? AND used = ? AND expires_at > ?"
)

func ForgotPassword(email string) (string, string, error) {
	// Sanitize input
	email = utils.SanitizeEmail(email)

	// Check if email exists in users table
	var user models.User
	if err := connection.DB.Where(emailQuery, email).First(&user).Error; err != nil {
		// For security, don't reveal if email exists
		return "", "", nil
	}

	// Generate secure 6-digit code using crypto/rand
	code, err := utils.GenerateVerificationCode()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate verification code: %w", err)
	}

	// Save to password_resets table
	reset := models.PasswordReset{
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(15 * time.Minute), // 15 minutes only
		Used:      false,
	}

	if err := connection.DB.Omit("id").Create(&reset).Error; err != nil {
		return "", "", fmt.Errorf("failed to create reset record: %v", err)
	}

	// Generate JWT token for password reset flow
	token, err := GeneratePasswordResetToken(email, code)
	if err != nil {
		return code, "", fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Send email with code
	content := fmt.Sprintf(`<p>You requested to reset your password.</p>
		<p><strong>Reset code:</strong></p>
		<p style="font-size:22px; font-weight:700; letter-spacing:2px;">%s</p>
		<p>This code will expire in 15 minutes.</p>
		<p>If you didn't request this, please ignore this email.</p>
	`, code)

	footer := `<p class="muted">If you didn't request this change, contact support immediately.</p>`
	htmlBody := BuildHTMLEmail("Password reset code", "Password Reset Code", content, footer)

	if err := SendEmail(email, "Password Reset Code - Attendance System", htmlBody); err != nil {
		// Log error without exposing sensitive information to client
		fmt.Printf("Failed to send reset email to %s: %v\n", email, err)
	}

	return code, token, nil
}

func VerifyResetCode(email, code string) error {
	var reset models.PasswordReset
	if err := connection.DB.Where(codeQuery, email, code, false, time.Now()).First(&reset).Error; err != nil {
		return errors.New("invalid or expired reset code")
	}
	return nil
}

// VerifyResetCodeDB is an alias for VerifyResetCode for controller usage
func VerifyResetCodeDB(email, code string) error {
	return VerifyResetCode(email, code)
}

func ResetPasswordWithCode(email, code, newPassword string) error {
	// Sanitize inputs
	email = utils.SanitizeEmail(email)
	code = strings.TrimSpace(code) // Just trim, don't sanitize the code itself

	// Verify code first - try exact match first
	var reset models.PasswordReset
	if err := connection.DB.Where(codeQuery, email, code, false, time.Now()).First(&reset).Error; err != nil {
		// If exact match fails, fetch all valid codes and compare with trimming
		var allResets []models.PasswordReset
		if err := connection.DB.Where("email = ? AND used = ? AND expires_at > ?", email, false, time.Now()).Find(&allResets).Error; err != nil {
			return errors.New("invalid or expired reset code")
		}

		// Check if any code matches (with trimming)
		found := false
		for _, r := range allResets {
			if strings.TrimSpace(r.Code) == code {
				reset = r
				found = true
				break
			}
		}

		if !found {
			return errors.New("invalid or expired reset code")
		}
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Find user
	var user models.User
	if err := connection.DB.Where(emailQuery, email).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	// Update password
	if err := connection.DB.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		return fmt.Errorf("failed to update password: %v", err)
	}

	// Mark code as used
	if err := connection.DB.Model(&reset).Update("used", true).Error; err != nil {
		// Log error without exposing sensitive information
		fmt.Printf("Failed to mark code as used\n")
	}

	// Send confirmation email
	content := `<p>Your password has been successfully reset.</p>
		<p>You can now log in with your new password.</p>
		<p>If you didn't make this change, contact support immediately.</p>`
	footer := `<p class="muted">If you did not initiate this change, please contact support.</p>`
	htmlBody := BuildHTMLEmail("Password changed", "Password Changed", content, footer)

	if err := SendEmail(email, "Password Changed - Attendance System", htmlBody); err != nil {
		// Log error without exposing sensitive information to client
		fmt.Printf("Failed to send confirmation email to %s: %v\n", email, err)
	}

	return nil
}

// Resend reset code
func ResendResetCode(email string) (string, string, error) {
	// Delete any existing unused codes for this email
	connection.DB.Where("email = ? AND used = ?", email, false).Delete(&models.PasswordReset{})

	// Generate and send new code
	return ForgotPassword(email)
}

// utils/validation.go
package utils

import (
	"regexp"
	"strings"
)

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidatePassword validates password strength
func ValidatePassword(password string) (bool, string) {
	if len(password) < 8 {
		return false, "Password must be at least 8 characters long"
	}
	if len(password) > 128 {
		return false, "Password must be less than 128 characters"
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasUpper || !hasLower || !hasNumber {
		return false, "Password must contain at least one uppercase letter, one lowercase letter, and one number"
	}
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_=+\[\]{};:'",.<>?/\\|~` + "`" + `]`).MatchString(password)
	if !hasSpecial {
		return false, "Password must contain at least one special character (!@#$%^&*() etc)"
	}
	return true, ""
}

// SanitizeString removes potentially dangerous characters and SQL patterns
func SanitizeString(input string) string {
	// Remove null bytes and control characters
	input = strings.ReplaceAll(input, "\x00", "")
	input = strings.TrimSpace(input)
	// Remove potential SQL injection patterns
	input = strings.ReplaceAll(input, "'", "")
	input = strings.ReplaceAll(input, "\"", "")
	input = strings.ReplaceAll(input, ";", "")
	input = strings.ReplaceAll(input, "--", "")
	input = strings.ReplaceAll(input, "/*", "")
	input = strings.ReplaceAll(input, "*/", "")
	// Block SQL keywords (case-insensitive check)
	upper := strings.ToUpper(input)
	sqlKeywords := []string{"UNION", "SELECT", "INSERT", "UPDATE", "DELETE", "DROP", "CREATE", "ALTER", "EXEC", "EXECUTE"}
	for _, keyword := range sqlKeywords {
		if strings.Contains(upper, keyword) {
			input = strings.ReplaceAll(input, keyword, "")
			input = strings.ReplaceAll(strings.ToLower(input), strings.ToLower(keyword), "")
		}
	}
	// Collapse multiple consecutive hyphens into a single hyphen
	reHyphen := regexp.MustCompile("-{2,}")
	input = reHyphen.ReplaceAllString(input, "-")
	return input
}

// SanitizeEmail sanitizes email input
func SanitizeEmail(email string) string {
	email = strings.ToLower(strings.TrimSpace(email))
	return SanitizeString(email)
}

// SanitizeStudentID sanitizes student ID input
func SanitizeStudentID(studentID string) string {
	studentID = strings.ToUpper(strings.TrimSpace(studentID))
	return SanitizeString(studentID)
}

// ValidateStudentID validates student ID format
func ValidateStudentID(studentID string) bool {
	// Allow alphanumeric and hyphens, 3-50 characters
	studentIDRegex := regexp.MustCompile(`^[A-Za-z0-9\-]{3,50}$`)
	return studentIDRegex.MatchString(studentID)
}

// ValidateBase64Image validates base64 image format
func ValidateBase64Image(base64Str string) bool {
	// Check if it's a valid base64 image data URL
	imageRegex := regexp.MustCompile(`^data:image\/(jpeg|jpg|png|gif|webp);base64,[A-Za-z0-9+/=]+$`)
	return imageRegex.MatchString(base64Str)
}

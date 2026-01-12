// services/login_service.go
package services

import (
	"attendance-system/connection"
	"attendance-system/models"
	"attendance-system/utils"
	"errors"
	"strings"
)

const errInvalidCredentials = "invalid credentials"

func LoginService(req models.LoginRequest) error {
	// Keep original to detect email vs student id
	original := req.StudentID

	// If the user passed an email in the StudentID field, delegate to email login
	if strings.Contains(original, "@") || utils.ValidateEmail(original) {
		return LoginByEmailService(original, req.Password)
	}

	// Sanitize input as a StudentID and continue
	req.StudentID = utils.SanitizeStudentID(req.StudentID)

	var user models.User

	// Find user by StudentID
	err := connection.DB.Where("student_id = ?", req.StudentID).First(&user).Error
	if err != nil {
		return errors.New(errInvalidCredentials)
	}

	// Check if user is verified
	if !user.IsVerified {
		return errors.New("account not verified")
	}

	if err := utils.ComparePassword(user.Password, req.Password); err != nil {
		return errors.New(errInvalidCredentials)
	}

	return nil
}

// Optional: Add login by email service
func LoginByEmailService(email, password string) error {
	// Sanitize input
	email = utils.SanitizeEmail(email)

	var user models.User

	// Find user by Email
	err := connection.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return errors.New(errInvalidCredentials)
	}

	// Check if user is verified
	if !user.IsVerified {
		return errors.New("account not verified")
	}

	if err := utils.ComparePassword(user.Password, password); err != nil {
		return errors.New(errInvalidCredentials)
	}

	return nil
}

// GetUserByEmail loads a user by email into provided pointer
func GetUserByEmail(email string, out *models.User) error {
	email = utils.SanitizeEmail(email)
	return connection.DB.Where("email = ?", email).First(out).Error
}

// GetUserByStudentID loads a user by student id into provided pointer
func GetUserByStudentID(studentID string, out *models.User) error {
	return connection.DB.Where("student_id = ?", studentID).First(out).Error
}

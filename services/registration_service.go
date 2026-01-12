// services/registration_service.go
package services

import (
	"attendance-system/connection"
	"attendance-system/models"
	"attendance-system/utils"
	"errors"
	"fmt"
	"time"
)

func RegisterService(req models.RegisterRequest) (string, string, error) {
	// Sanitize inputs
	req.Email = utils.SanitizeEmail(req.Email)
	if req.StudentID != "" {
		req.StudentID = utils.SanitizeStudentID(req.StudentID)
	}

	if err := validateRegistrationInput(req); err != nil {
		return "", "", err
	}

	studentID, err := generateOrValidateStudentID(req.StudentID)
	if err != nil {
		return "", "", err
	}

	if err := checkEmailDuplicates(req.Email); err != nil {
		return "", "", err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return "", "", err
	}

	verificationCode, err := utils.GenerateVerificationCode()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate verification code: %w", err)
	}

	_, err = createPendingUser(req, studentID, hashedPassword, verificationCode)
	if err != nil {
		return "", "", err
	}

	if err := sendVerificationEmail(req.Email, studentID, verificationCode); err != nil {
		// Log error without exposing sensitive information
		fmt.Printf("Failed to send verification email\n")
	}

	// Generate email verification token
	token, err := GenerateEmailVerificationToken(req.Email, verificationCode)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate verification token: %w", err)
	}

	return studentID, token, nil
}

func validateRegistrationInput(req models.RegisterRequest) error {
	if req.Email == "" || req.Password == "" {
		return errors.New("email and password are required")
	}

	// Validate email format
	if !utils.ValidateEmail(req.Email) {
		return errors.New("invalid email format")
	}

	// Validate password strength
	if valid, msg := utils.ValidatePassword(req.Password); !valid {
		return errors.New(msg)
	}

	// Validate student ID if provided
	if req.StudentID != "" && !utils.ValidateStudentID(req.StudentID) {
		return errors.New("invalid student ID format")
	}

	return nil
}

func generateOrValidateStudentID(inputStudentID string) (string, error) {
	if inputStudentID == "" {
		return generateUniqueStudentID()
	}
	return validateExistingStudentID(inputStudentID)
}

func generateUniqueStudentID() (string, error) {
	currentTime := time.Now()
	datePart := currentTime.Format("060102")

	for i := 0; i < 10; i++ {
		randomPart, err := utils.GenerateSecureCode(4)
		if err != nil {
			return "", fmt.Errorf("failed to generate student ID: %w", err)
		}
		studentID := datePart + "-" + randomPart

		if isStudentIDUnique(studentID) {
			return studentID, nil
		}
	}

	return "", errors.New("failed to generate unique student ID")
}

func isStudentIDUnique(studentID string) bool {
	var existingUser models.User
	var existingPending models.PendingUser

	userResult := connection.DB.Where("student_id = ?", studentID).First(&existingUser)
	pendingResult := connection.DB.Where("student_id = ?", studentID).First(&existingPending)

	return userResult.Error != nil && pendingResult.Error != nil
}

func validateExistingStudentID(studentID string) (string, error) {
	if !isStudentIDUnique(studentID) {
		return "", errors.New("student_id already exists")
	}
	return studentID, nil
}

func checkEmailDuplicates(email string) error {
	var existingUser models.User
	var existingPending models.PendingUser

	if connection.DB.Where("email = ?", email).First(&existingUser).Error == nil {
		return errors.New("email already registered")
	}

	if connection.DB.Where("email = ?", email).First(&existingPending).Error == nil {
		return errors.New("email already pending verification")
	}

	return nil
}

func createPendingUser(req models.RegisterRequest, studentID, hashedPassword, verificationCode string) (*models.PendingUser, error) {
	pending := &models.PendingUser{
		StudentID:        studentID,
		Email:            req.Email,
		Password:         hashedPassword,
		Username:         req.Username,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		MiddleName:       req.MiddleName,
		Course:           req.Course,
		YearLevel:        req.YearLevel,
		Section:          req.Section,
		Department:       req.Department,
		College:          req.College,
		ContactNumber:    req.ContactNumber,
		Address:          req.Address,
		VerificationCode: verificationCode,
		ExpiresAt:        time.Now().Add(30 * time.Minute),
	}

	if err := connection.DB.Omit("id").Create(pending).Error; err != nil {
		return nil, fmt.Errorf("failed to save pending user: %v", err)
	}

	return pending, nil
}

func sendVerificationEmail(email, studentID, verificationCode string) error {
	// Build modern HTML email
	content := fmt.Sprintf(`<p>Thank you for registering with the Attendance System.</p>
		<p><strong>Student ID:</strong> <code>%s</code></p>
		<p><strong>Verification code:</strong></p>
		<p style="font-size:22px; font-weight:700; letter-spacing:2px;">%s</p>
		<p>This code will expire in 30 minutes. Enter it on the verification page to complete your registration.</p>
	`, studentID, verificationCode)

	footer := `<p class="muted">If you did not register, please ignore this message.</p>`
	htmlBody := BuildHTMLEmail("Verify your email", "Email Verification", content, footer)

	return SendEmail(email, "Verification Code - Attendance System", htmlBody)
}

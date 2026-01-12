package controller

import (
	"attendance-system/models"
	"attendance-system/services"
	"attendance-system/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	req := new(models.LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": models.ErrInvalidRequest})
	}

	if req.StudentID == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": models.ErrFieldsRequired})
	}

	// Fetch user profile
	// Determine lookup key (email vs student id)
	var user models.User
	if strings.Contains(req.StudentID, "@") {
		// login by email
		if err := services.GetUserByEmail(req.StudentID, &user); err != nil {
			return c.Status(401).JSON(fiber.Map{"error": models.ErrInvalidCredentials})
		}
	} else {
		sid := utils.SanitizeStudentID(req.StudentID)
		if err := services.GetUserByStudentID(sid, &user); err != nil {
			return c.Status(401).JSON(fiber.Map{"error": models.ErrInvalidStudentID})
		}
	}

	// Verify password
	if err := utils.ComparePassword(user.Password, req.Password); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": models.ErrInvalidStudentID})
	}

	// Check if verified
	if !user.IsVerified {
		return c.Status(403).JSON(fiber.Map{"error": models.ErrEmailNotVerified})
	}

	accessToken, err := services.GenerateAccessToken(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": models.ErrFailedGenerateAccessToken})
	}

	refreshToken, err := services.GenerateRefreshToken(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": models.ErrFailedGenerateRefreshToken})
	}

	return c.JSON(fiber.Map{
		"message":       models.SuccessLoginFacultyAdmin,
		"student_id":    user.StudentID,
		"email":         user.Email,
		"role":          user.Role,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// LoginByEmail adds login by email endpoint
func LoginByEmail(c *fiber.Ctx) error {
	type EmailLoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	req := new(EmailLoginRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": models.ErrInvalidRequest})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": models.ErrEmailPasswordRequired})
	}

	// Get user by email
	var user models.User
	if err := services.GetUserByEmail(req.Email, &user); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": models.ErrInvalidCredentials})
	}

	// Verify password
	if err := utils.ComparePassword(user.Password, req.Password); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": models.ErrInvalidCredentials})
	}

	// Check if verified
	if !user.IsVerified {
		return c.Status(403).JSON(fiber.Map{"error": models.ErrEmailNotVerified})
	}

	accessToken, err := services.GenerateAccessToken(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": models.ErrFailedGenerateAccessToken})
	}

	refreshToken, err := services.GenerateRefreshToken(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": models.ErrFailedGenerateRefreshToken})
	}

	return c.JSON(fiber.Map{
		"message":       models.SuccessLoginFacultyAdmin,
		"email":         user.Email,
		"student_id":    user.StudentID,
		"role":          user.Role,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// RefreshToken generates a new access token from a valid refresh token
func RefreshToken(c *fiber.Ctx) error {
	type RefreshRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	req := new(RefreshRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": models.ErrInvalidRequest})
	}

	if req.RefreshToken == "" {
		return c.Status(400).JSON(fiber.Map{"error": models.ErrTokenRequired})
	}

	// Verify refresh token
	claims, err := services.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": models.ErrRefreshTokenExpired})
	}

	// Get user
	var user models.User
	if err := services.GetUserByStudentID(claims.StudentID, &user); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": models.ErrUserNotFound})
	}

	// Generate new access token
	accessToken, err := services.GenerateAccessToken(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": models.ErrFailedGenerateAccessToken})
	}

	return c.JSON(fiber.Map{
		"message":      models.SuccessTokenRefreshed,
		"access_token": accessToken,
	})
}

// GetProfile returns the authenticated user's profile
func GetProfile(c *fiber.Ctx) error {
	// Get user from context (set by RequireAuth middleware)
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": models.ErrUnauthorized})
	}

	return c.JSON(fiber.Map{
		"message":    models.SuccessProfileFetched,
		"student_id": user.StudentID,
		"email":      user.Email,
		"role":       user.Role,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
	})
}

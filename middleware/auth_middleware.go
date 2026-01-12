package middleware

import (
	"attendance-system/connection"
	"attendance-system/models"
	"attendance-system/services"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// extractStudentID extracts student ID from Authorization header (supports both JWT and legacy Bearer format)
func extractStudentID(c *fiber.Ctx) (string, string, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		// Try alternative header
		authHeader = c.Get("X-Student-ID")
		if authHeader == "" {
			return "", "", fiber.NewError(fiber.StatusUnauthorized, "Authorization required")
		}
		// Set as Bearer format for consistency
		authHeader = "Bearer " + authHeader
	}

	// Extract from "Bearer {token}"
	parts := strings.Split(authHeader, " ")
	var token string
	if len(parts) == 2 && parts[0] == "Bearer" {
		token = parts[1]
	} else if len(parts) == 1 {
		token = parts[0]
	} else {
		return "", "", fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization format. Use: Bearer {JWT_token}")
	}

	if token == "" {
		return "", "", fiber.NewError(fiber.StatusUnauthorized, "Token is required")
	}

	// Try JWT first
	claims, err := services.VerifyAccessToken(token)
	if err == nil {
		// Valid JWT token
		return claims.StudentID, claims.Role, nil
	}

	// Fallback: treat as legacy student ID format (for backward compatibility)
	// This should be deprecated in future versions
	return token, "", nil
}

func RequireSuperAdmin(c *fiber.Ctx) error {
	studentID, _, err := extractStudentID(c)
	if err != nil {
		return err
	}

	// Find user
	var user models.User
	if err := connection.DB.Where("student_id = ?", studentID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Check if user is superadmin
	if user.Role != models.RoleSuperAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied - superadmin role required",
		})
	}

	// Store user info in context for later use
	c.Locals("user", user)

	return c.Next()
}

// RequireAuth - JWT authentication middleware (checks if user exists and is verified)
func RequireAuth(c *fiber.Ctx) error {
	studentID, _, err := extractStudentID(c)
	if err != nil {
		return err
	}

	// Find user
	var user models.User
	if err := connection.DB.Where("student_id = ?", studentID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Check if user is verified
	if !user.IsVerified {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Account not verified. Please verify your email first.",
		})
	}

	// Store user info in context
	c.Locals("user", user)

	return c.Next()
}

// RequireAdmin - middleware for organization/event managers (admin + superadmin)
func RequireAdmin(c *fiber.Ctx) error {
	studentID, _, err := extractStudentID(c)
	if err != nil {
		return err
	}

	// Find user
	var user models.User
	if err := connection.DB.Where("student_id = ?", studentID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Allow both admin and superadmin
	if user.Role != models.RoleAdmin && user.Role != models.RoleSuperAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied - admin role required",
		})
	}

	// Ensure verified
	if !user.IsVerified {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Account not verified. Please verify your email first.",
		})
	}

	c.Locals("user", user)
	return c.Next()
}

// RequireFacultyOrAdmin - Requires faculty, admin, or superadmin role
func RequireFacultyOrAdmin(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	validRoles := map[string]bool{
		models.RoleSuperAdmin: true,
		models.RoleAdmin:      true,
		models.RoleFaculty:    true,
	}

	if !validRoles[user.Role] {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only faculty, admin, or superadmin can perform this action",
		})
	}

	return c.Next()
}

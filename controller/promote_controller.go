package controller

import (
	"attendance-system/connection"
	"attendance-system/models"
	"attendance-system/services"

	"github.com/gofiber/fiber/v2"
)

func GetAllUsers(c *fiber.Ctx) error {
	var users []models.User
	if err := connection.DB.Find(&users).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Users retrieved successfully",
		"status":  "success",
		"users":   users,
	})
}

func PromoteUser(c *fiber.Ctx) error {
	req := new(models.PromoteRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	if req.StudentID == "" || req.Role == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "student_id and role are required",
		})
	}

	// Validate role
	validRoles := map[string]bool{
		"student": true,
		"faculty": true,
		"admin":   true,
	}
	if !validRoles[req.Role] {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid role. Must be 'student', 'faculty', or 'admin'",
		})
	}

	// Get superadmin user from context
	superadmin, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Check if user exists
	var user models.User
	if err := connection.DB.Where("student_id = ?", req.StudentID).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Student not found",
		})
	}

	// Update the user's role
	if err := connection.DB.Model(&user).Update("role", req.Role).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update user role",
		})
	}

	// Log the action for audit trail
	go services.LogAuditAction(
		services.AuditUserPromoted,
		superadmin.StudentID,
		req.StudentID,
		"Promoted to "+req.Role,
		c.IP(),
	)

	return c.JSON(fiber.Map{
		"message":    "User promoted successfully",
		"status":     "success",
		"student_id": req.StudentID,
		"new_role":   req.Role,
	})
}

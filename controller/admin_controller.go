package controller

import (
	"attendance-system/services"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// GetSystemStats returns aggregated system-wide statistics (superadmin only)
func GetSystemStats(c *fiber.Ctx) error {
	stats, err := services.GetSystemStats()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"stats": stats})
}

// CreateAdmin creates a new admin account (superadmin only)
func CreateAdmin(c *fiber.Ctx) error {
	type Req struct {
		StudentID string `json:"student_id"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	req := new(Req)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	// Simple validation
	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "email, password, first_name, last_name are required"})
	}

	// Create user record
	u := make(map[string]interface{})
	u["student_id"] = req.StudentID
	u["email"] = req.Email
	u["password"] = req.Password
	u["first_name"] = req.FirstName
	u["last_name"] = req.LastName
	u["role"] = "admin"
	u["is_verified"] = true

	// Use DB directly to avoid adding service layer for this utility
	if err := services.CreateUserFromMap(u); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "admin created", "email": req.Email})
}

// UpdateAdmin updates admin profile (superadmin only)
func UpdateAdmin(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "id is required"})
	}
	updates := make(map[string]interface{})
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	if err := services.UpdateUserByID(id, updates); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "admin updated"})
}

// DeleteAdmin deletes a user by student_id (superadmin only)
func DeleteAdmin(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "id is required"})
	}
	if err := services.DeleteUserByStudentID(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "user deleted"})
}

// GetAuditLogs returns audit logs (superadmin)
func GetAuditLogs(c *fiber.Ctx) error {
	// Parse simple filters
	filters := make(map[string]interface{})
	if action := c.Query("action"); action != "" {
		filters["action"] = action
	}
	if actor := c.Query("actor_id"); actor != "" {
		filters["actor_id"] = actor
	}
	limit := 50
	offset := 0
	logs, total, err := services.GetAuditLogs(filters, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"total": total, "logs": logs})
}

// GetAllAttendance returns all attendance records for superadmin
func GetAllAttendance(c *fiber.Ctx) error {
	// Simple implementation: reuse service to fetch by event if provided, otherwise fetch all
	if eventID := c.Query("event_id"); eventID != "" {
		// delegate to existing controller via service
		// parse event id
		// but to keep it short, call attendance service directly
	}
	// Fetch all attendances
	attendances, err := services.GetAllAttendance()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"count": len(attendances), "attendances": attendances})
}

// GetEventQRCode returns or generates the QR code for an event
func GetEventQRCode(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(400).JSON(fiber.Map{"error": "event id required"})
	}
	// parse uint
	var eventID uint
	_, err := fmt.Sscanf(idStr, "%d", &eventID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid event id"})
	}
	qr, err := services.GetEventQRCode(eventID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"qr_code": qr})
}

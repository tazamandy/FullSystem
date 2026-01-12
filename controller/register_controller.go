package controller

import (
	"attendance-system/models"
	"attendance-system/services"
	"attendance-system/utils"

	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	var req models.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Email, password, first name, and last name are required"})
	}

	studentID, token, err := services.RegisterService(req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"message":    "Registration successful. Please check your email to verify your account.",
		"student_id": studentID,
		"token":      token,
		"status":     "success",
	})
}

// GetRegistrationDropdowns returns predefined departments and sections for registration dropdowns
func GetRegistrationDropdowns(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"departments": utils.Departments,
		"sections":    utils.Sections,
	})
}

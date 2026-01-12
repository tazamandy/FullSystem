package controller

import (
	"attendance-system/models"
	"attendance-system/services"
	"attendance-system/utils"

	"github.com/gofiber/fiber/v2"
)

const (
	ErrInvalidRequest    = "Invalid request"
	ErrEmailRequired     = "Email is required"
	ErrCodeRequired      = "Code is required"
	ErrPasswordRequired  = "Password is required"
	ErrAllFieldsRequired = "All fields are required"
	SuccessResetCodeSent = "If your email is registered, you will receive a reset code."
	SuccessCodeValid     = "Code is valid"
	SuccessPasswordReset = "Password reset successful"
	SuccessNewCodeSent   = "New code sent if email is registered"
)

func ForgotPassword(c *fiber.Ctx) error {
	req := new(models.ResetPasswordRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": ErrInvalidRequest})
	}

	if req.Email == "" {
		return c.Status(400).JSON(fiber.Map{"error": ErrEmailRequired})
	}

	// Always return success (security), but get token if email exists
	_, token, err := services.ForgotPassword(req.Email)
	if err != nil {
		// Still return success message for security, but with empty token
		return c.JSON(fiber.Map{
			"message": SuccessResetCodeSent,
			"status":  "success",
			"token":   "",
		})
	}

	return c.JSON(fiber.Map{
		"message": SuccessResetCodeSent,
		"status":  "success",
		"token":   token,
	})
}

// VerifyResetCode verifies the reset code using bearer token
// Request Header: Authorization: Bearer <token_from_forgot_password>
// Request: { "code": "123456" }  <- code received in email
// Response: { "message": "Code is valid", "status": "success", "token": "new_token_for_reset_password" }
func VerifyResetCode(c *fiber.Ctx) error {
	type VerifyRequest struct {
		Code string `json:"code"`
	}

	req := new(VerifyRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": ErrInvalidRequest})
	}

	if req.Code == "" {
		return c.Status(400).JSON(fiber.Map{"error": ErrCodeRequired})
	}

	// Get bearer token from Authorization header
	auth := c.Get("Authorization")
	if auth == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Authorization header required"})
	}

	// Extract token from "Bearer <token>"
	const bearerPrefix = "Bearer "
	if len(auth) < len(bearerPrefix) || auth[:len(bearerPrefix)] != bearerPrefix {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid authorization format. Use: Bearer <token>"})
	}

	token := auth[len(bearerPrefix):]

	// Verify the JWT token
	claims, err := services.VerifyPasswordResetToken(token)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	// Verify the code in database with the email from token
	err = services.VerifyResetCodeDB(claims.Email, req.Code)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid or expired code"})
	}

	// Code is valid - generate a new token for /reset-password endpoint
	newToken, err := services.GeneratePasswordResetToken(claims.Email, req.Code)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{
		"message": SuccessCodeValid,
		"status":  "success",
		"email":   claims.Email,
		"token":   newToken,
	})
}

func ResetPassword(c *fiber.Ctx) error {
	type ResetPasswordRequest struct {
		NewPassword        string `json:"new_password"`
		ConfirmNewPassword string `json:"confirm_new_password"`
	}

	req := new(ResetPasswordRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": ErrInvalidRequest})
	}

	if req.NewPassword == "" || req.ConfirmNewPassword == "" {
		return c.Status(400).JSON(fiber.Map{"error": ErrAllFieldsRequired})
	}

	// Validate passwords match
	if req.NewPassword != req.ConfirmNewPassword {
		return c.Status(400).JSON(fiber.Map{"error": "Passwords do not match"})
	}

	if valid, msg := utils.ValidatePassword(req.NewPassword); !valid {
		return c.Status(400).JSON(fiber.Map{"error": msg})
	}

	// Get bearer token from Authorization header
	auth := c.Get("Authorization")
	if auth == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Authorization header required"})
	}

	// Extract token from "Bearer <token>"
	const bearerPrefix = "Bearer "
	if len(auth) < len(bearerPrefix) || auth[:len(bearerPrefix)] != bearerPrefix {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid authorization format. Use: Bearer <token>"})
	}

	token := auth[len(bearerPrefix):]

	// Verify the JWT token
	claims, err := services.VerifyPasswordResetToken(token)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	// Reset password using email from token (code already verified in /verify-reset-code)
	err = services.ResetPasswordWithCode(claims.Email, claims.Code, req.NewPassword)
	if err != nil {
		if err.Error() == "invalid or expired reset code" {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid or expired code"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to reset password"})
	}

	return c.JSON(fiber.Map{
		"message": SuccessPasswordReset,
		"status":  "success",
	})
}

func ResendCode(c *fiber.Ctx) error {
	type Request struct {
		Email string `json:"email"`
	}

	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": ErrInvalidRequest})
	}

	if req.Email == "" {
		return c.Status(400).JSON(fiber.Map{"error": ErrEmailRequired})
	}

	_, token, err := services.ResendResetCode(req.Email)
	if err != nil {
	}

	return c.JSON(fiber.Map{
		"message": SuccessNewCodeSent,
		"status":  "success",
		"token":   token,
	})
}

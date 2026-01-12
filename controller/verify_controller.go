package controller

import (
	"attendance-system/connection"
	"attendance-system/models"
	"attendance-system/services"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"
)

func VerifyEmail(c *fiber.Ctx) error {
	// Accept multiple verification methods for convenience:
	// 1) Authorization: Bearer <token> (preferred)
	// 2) JSON body: { "token": "...", "code": "..." }
	// 3) Fallback JSON body: { "email": "...", "code": "..." } (less secure)

	// Parse request data (accept code, token or email)
	var req struct {
		Code  string `json:"code"`
		Token string `json:"token,omitempty"`
		Email string `json:"email,omitempty"`
	}

	// Support GET with query parameters for easier testing (e.g. Postman GET)
	if c.Method() == "GET" {
		req.Code = c.Query("code")
		req.Token = c.Query("token")
		req.Email = c.Query("email")
		// If no query params provided, try parsing JSON body for clients that send body with GET (Postman)
		if req.Code == "" && req.Token == "" && req.Email == "" {
			if err := c.BodyParser(&req); err != nil {
				return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
			}
		}
	} else {
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}
	}

	if req.Code == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Code is required"})
	}

	// Determine token: prefer Authorization header, then body token
	authHeader := c.Get("Authorization")
	var token string
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			token = parts[1]
		}
	}
	if token == "" && req.Token != "" {
		token = req.Token
	}

	var pending models.PendingUser
	verifiedByToken := false
	var claims *models.EmailVerificationTokenClaims

	if token != "" {
		// Try to verify token. If token is invalid/expired we will fall back
		// to email+code verification when the client provided `email`.
		tokenClaims, err := services.VerifyEmailVerificationToken(token)
		if err == nil {
			claims = tokenClaims
			// Verify the code matches the token claims
			if strings.TrimSpace(req.Code) != strings.TrimSpace(claims.Code) {
				return c.Status(400).JSON(fiber.Map{"error": "Invalid verification code"})
			}

			// Find pending user by email from claims
			if err := connection.DB.Where("email = ?", claims.Email).First(&pending).Error; err != nil {
				return c.Status(400).JSON(fiber.Map{"error": "Registration not found"})
			}

			// Check if code expired (compare in UTC)
			if time.Now().UTC().After(pending.ExpiresAt.UTC()) {
				connection.DB.Delete(&pending)
				return c.Status(400).JSON(fiber.Map{"error": "Verification code has expired. Please register again."})
			}

			verifiedByToken = true
		} else {
			// Token invalid/expired. If no email provided, return 401; otherwise fall through to email+code fallback.
			if req.Email == "" {
				return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired token"})
			}
		}
	}

	if !verifiedByToken {
		// No valid token — verify by email + code (fallback)
		if req.Email == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Authorization header or token is required (or provide email+code)"})
		}

		if err := connection.DB.Where("email = ?", req.Email).First(&pending).Error; err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Registration not found"})
		}

		// Validate code and expiry (compare in UTC)
		if strings.TrimSpace(req.Code) != strings.TrimSpace(pending.VerificationCode) {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid verification code"})
		}
		if time.Now().UTC().After(pending.ExpiresAt.UTC()) {
			connection.DB.Delete(&pending)
			return c.Status(400).JSON(fiber.Map{"error": "Verification code has expired. Please register again."})
		}
	}

	// Generate QR code
	qrContent := fmt.Sprintf("student:%s", pending.StudentID)
	qrCodePNG, err := qrcode.Encode(qrContent, qrcode.Medium, 256)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate QR code"})
	}

	qrCodeBase64 := "data:image/png;base64," + base64.StdEncoding.EncodeToString(qrCodePNG)

	// Move to User table
	user := models.User{
		StudentID:     pending.StudentID,
		Email:         pending.Email,
		Password:      pending.Password,
		Username:      pending.Username,
		FirstName:     pending.FirstName,
		LastName:      pending.LastName,
		MiddleName:    pending.MiddleName,
		Course:        pending.Course,
		YearLevel:     pending.YearLevel,
		Section:       pending.Section,
		Department:    pending.Department,
		College:       pending.College,
		ContactNumber: pending.ContactNumber,
		Address:       pending.Address,

		QRCodeData: qrCodeBase64,
		IsVerified: true,
		VerifiedAt: time.Now(),
	}

	// Ensure ID is zero so DB assigns it (defensive against client-provided IDs)
	user.ID = 0
	if err := connection.DB.Omit("id").Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user account"})
	}

	// Delete from pending table
	connection.DB.Delete(&pending)

	return c.JSON(fiber.Map{
		"message": "Email verification successful",
		"status":  "success",
	})
}

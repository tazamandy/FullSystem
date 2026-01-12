// middleware/security.go
package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// SecurityHeaders adds security headers to responses
func SecurityHeaders(c *fiber.Ctx) error {
	// Prevent clickjacking
	c.Set("X-Frame-Options", "DENY")
	// Prevent MIME type sniffing
	c.Set("X-Content-Type-Options", "nosniff")
	// Enable XSS protection
	c.Set("X-XSS-Protection", "1; mode=block")
	// Enforce HTTPS in production (set via env)
	c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	// Content Security Policy
	c.Set("Content-Security-Policy", "default-src 'self'")
	// Referrer Policy
	c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
	// Permissions Policy
	c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

	return c.Next()
}

package middleware

import (
	"time"

	"attendance-system/logging"

	"github.com/gofiber/fiber/v2"
)

// RequestLogger middleware logs all HTTP requests with IP, method, path, status, and duration
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Get client IP (handle proxies)
		clientIP := c.IP()
		if forwarded := c.Get("X-Forwarded-For"); forwarded != "" {
			clientIP = forwarded
		}
		if realIP := c.Get("X-Real-IP"); realIP != "" {
			clientIP = realIP
		}

		// Process request
		err := c.Next()

		// Log after request
		duration := time.Since(start)
		status := c.Response().StatusCode()

		// Use our custom logger
		logging.LogRequest(clientIP, c.Method(), c.Path(), status, duration)

		return err
	}
}

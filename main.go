package main

import (
	"attendance-system/API"
	"attendance-system/connection"
	"attendance-system/logging"
	"attendance-system/middleware"
	"attendance-system/seeder"
	"attendance-system/services"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load env
	_ = godotenv.Load()

	// Init logger
	if err := logging.InitLogger(); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
	}

	// Database
	connection.Connect()

	// Seed default admin
	seeder.SeedSuperAdmin()

	// Background job
	go startEventStatusChecker()

	// Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			msg := "An error occurred"
			if code < 500 {
				msg = err.Error()
			}

			return c.Status(code).JSON(fiber.Map{
				"error": msg,
			})
		},
		BodyLimit: 10 * 1024 * 1024, // 10MB
	})

	// ---------------- MIDDLEWARE ----------------

	// Security headers
	app.Use(middleware.SecurityHeaders)

	// Rate limiting
	app.Use(middleware.RateLimit)

	// CORS (Flutter Web SAFE)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // DEV ONLY
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		MaxAge:       300,
	}))

	// Request logger
	app.Use(middleware.RequestLogger())

	// ---------------- ROUTES ----------------

	API.AuthRoutes(app)
	API.EventRoutes(app)
	API.AttendanceRoutes(app)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "attendance-backend",
			"port":    "3000",
		})
	})

	// Root info
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"remote_ip":  c.IP(),
			"user_agent": c.Get("User-Agent"),
			"host":       c.Get("Host"),
		})
	})

	// ---------------- START SERVER ----------------

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("\n🚀 Server running on port %s\n", port)
	fmt.Printf("Health: http://localhost:%s/health\n", port)

	defer logging.Logger.Sync()
	if err := app.Listen(":" + port); err != nil {
		panic(err)
	}
}

// Background job
func startEventStatusChecker() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	// Run once at startup
	services.CheckAndUpdateCompletedEvents()

	for range ticker.C {
		services.CheckAndUpdateCompletedEvents()
	}
}
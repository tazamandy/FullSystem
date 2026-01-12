package API

import (
	"attendance-system/controller"
	"attendance-system/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	// PUBLIC routes - NO authentication required
	// Registration routes
	app.Post("/register", controller.Register)
	// New/expected route for verification
	app.Post("/reg/verify", controller.VerifyEmail)
	// Backwards-compatible routes used by older docs/tests
	app.Post("/verify", controller.VerifyEmail)
	app.Get("/verify", controller.VerifyEmail)

	// Get registration dropdown options (departments, sections)
	app.Get("/registration-dropdowns", controller.GetRegistrationDropdowns)

	// Login routes
	app.Post("/login", controller.Login)
	app.Post("/refresh-token", controller.RefreshToken)

	// Password Reset Routes - PUBLIC (no authentication required)
	// Direct root level routes (easier for frontend)
	app.Post("/forgot-password", controller.ForgotPassword)
	app.Post("/verify-reset-code", controller.VerifyResetCode)
	app.Post("/reset-password", controller.ResetPassword)
	app.Post("/resend-code", controller.ResendCode)

	// Also support /fgtp prefix group for backward compatibility
	fgtp := app.Group("/fgtp")
	{
		fgtp.Post("/forgot-password", controller.ForgotPassword)
		fgtp.Post("/verify-reset-code", controller.VerifyResetCode)
		fgtp.Post("/reset-password", controller.ResetPassword)
		fgtp.Post("/resend-code", controller.ResendCode)
	}

	// Protected routes (require authentication)
	protected := app.Group("", middleware.RequireAuth)
	{
		protected.Get("/profile", controller.GetProfile)
	}

	adminRoutes := app.Group("/admin", middleware.RequireAuth, middleware.RequireSuperAdmin)
	{
		adminRoutes.Get("/users", controller.GetAllUsers)
		adminRoutes.Get("/stats", controller.GetSystemStats)
		adminRoutes.Post("/promote", controller.PromoteUser)
	}

	// Admin-level (organization managers) routes
	// adminManager group reserved for admin-level routes; currently no endpoints defined.
}

func EventRoutes(app *fiber.App) {
	// Public routes
	app.Get("/events/creation-dropdowns", controller.GetEventCreationDropdowns)

	events := app.Group("/events", middleware.RequireAuth)
	{
		events.Get("/", controller.GetAllEvents)
		events.Get("/my-events", controller.GetMyEvents)
		events.Get("/:id", controller.GetEvent)
	}

	eventsProtected := app.Group("/events", middleware.RequireAuth, middleware.RequireFacultyOrAdmin)
	{
		eventsProtected.Post("/", controller.CreateEvent)
		eventsProtected.Put("/:id", controller.UpdateEvent)
		eventsProtected.Delete("/:id", controller.DeleteEvent)
	}
}

func AttendanceRoutes(app *fiber.App) {
	attendance := app.Group("/attendance", middleware.RequireAuth)
	{
		attendance.Post("/mark", controller.MarkAttendance)
		attendance.Get("/my-attendance", controller.GetMyAttendance)
		attendance.Get("/stats", controller.GetAttendanceStats)
	}

	// Specific route for event attendance - must be after /forgot-password and public routes
	attendanceByEvent := app.Group("/events/:event_id/attendance", middleware.RequireAuth)
	{
		attendanceByEvent.Get("/", controller.GetAttendanceByEvent)
	}

	attendanceAdmin := app.Group("/attendance", middleware.RequireAuth, middleware.RequireFacultyOrAdmin)
	{
		attendanceAdmin.Put("/:id/status", controller.UpdateAttendanceStatus)
	}
}

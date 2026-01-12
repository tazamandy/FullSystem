// controller/attendance_controller.go
package controller

import (
	"attendance-system/connection"
	"attendance-system/models"
	"attendance-system/services"
	"attendance-system/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// getUserAndCheckPermissions retrieves the user and checks if they can mark attendance for the given studentID
func getUserAndCheckPermissions(c *fiber.Ctx, targetStudentID string) (models.User, error) {
	// Get user from context
	user, ok := c.Locals("user").(models.User)
	if !ok {
		studentID := c.Get(utils.HeaderStudentID)
		if studentID == "" {
			return models.User{}, fiber.NewError(401, utils.ErrUnauthorized)
		}
		var dbUser models.User
		if err := connection.DB.Where("student_id = ?", studentID).First(&dbUser).Error; err != nil {
			return models.User{}, fiber.NewError(401, utils.ErrUserNotFound)
		}
		user = dbUser
	}

	// Only allow marking own attendance unless admin/faculty
	if targetStudentID != "" && targetStudentID != user.StudentID {
		if user.Role != "superadmin" && user.Role != "admin" && user.Role != "faculty" {
			return models.User{}, fiber.NewError(403, "You can only mark your own attendance")
		}
	}

	return user, nil
}

// MarkAttendance marks attendance for a student (check-in or check-out)
func MarkAttendance(c *fiber.Ctx) error {
	req := new(models.AttendanceRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request format"})
	}

	if req.EventID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "event_id is required"})
	}

	if req.Action == "" {
		return c.Status(400).JSON(fiber.Map{"error": "action is required. Use 'check_in' or 'check_out'"})
	}

	// If studentID not provided, use current user (but we need to get user first)
	user, err := getUserAndCheckPermissions(c, req.StudentID)
	if err != nil {
		return err
	}

	// If studentID not provided, use current user
	if req.StudentID == "" {
		req.StudentID = user.StudentID
	}

	attendance, err := services.MarkAttendance(*req, user.StudentID, user.Role)
	if err != nil {
		// Map service-level access denial to HTTP 403
		if err == services.ErrEventAccessDenied {
			return c.Status(403).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"message":    "Attendance marked successfully",
		"attendance": attendance,
	})
}

// GetAttendanceByEvent retrieves all attendance for an event
func GetAttendanceByEvent(c *fiber.Ctx) error {
	eventIDStr := c.Params("event_id")
	eventID, err := strconv.ParseUint(eventIDStr, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": utils.ErrInvalidEventID})
	}

	attendances, err := services.GetAttendanceByEvent(uint(eventID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"attendances": attendances,
		"count":       len(attendances),
	})
}

// GetMyAttendance retrieves attendance records for current student
func GetMyAttendance(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		studentID := c.Get(utils.HeaderStudentID)
		if studentID == "" {
			return c.Status(401).JSON(fiber.Map{"error": utils.ErrUnauthorized})
		}
		user.StudentID = studentID
	}

	filters := make(map[string]interface{})
	if eventID := c.Query("event_id"); eventID != "" {
		if id, err := strconv.ParseUint(eventID, 10, 32); err == nil {
			filters["event_id"] = uint(id)
		}
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if startDate := c.Query("start_date"); startDate != "" {
		filters["start_date"] = startDate
	}
	if endDate := c.Query("end_date"); endDate != "" {
		filters["end_date"] = endDate
	}

	attendances, err := services.GetAttendanceByStudent(user.StudentID, filters)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"attendances": attendances,
		"count":       len(attendances),
	})
}

// GetAttendanceStats retrieves attendance statistics
func GetAttendanceStats(c *fiber.Ctx) error {
	studentID := c.Query("student_id")
	eventIDStr := c.Query("event_id")

	var eventID *uint
	if eventIDStr != "" {
		if id, err := strconv.ParseUint(eventIDStr, 10, 32); err == nil {
			idUint := uint(id)
			eventID = &idUint
		}
	}

	// If no studentID provided, try to get from context
	if studentID == "" {
		user, ok := c.Locals("user").(models.User)
		if ok {
			studentID = user.StudentID
		}
	}

	stats, err := services.GetAttendanceStats(studentID, eventID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"stats": stats,
	})
}

// UpdateAttendanceStatus updates attendance status (admin/faculty only)
func UpdateAttendanceStatus(c *fiber.Ctx) error {
	attendanceIDStr := c.Params("id")
	attendanceID, err := strconv.ParseUint(attendanceIDStr, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": utils.ErrInvalidAttendanceID})
	}

	type UpdateRequest struct {
		Status string `json:"status"`
		Notes  string `json:"notes"`
	}

	req := new(UpdateRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request format"})
	}

	if req.Status == "" {
		return c.Status(400).JSON(fiber.Map{"error": "status is required"})
	}

	user, ok := c.Locals("user").(models.User)
	if !ok {
		studentID := c.Get(utils.HeaderStudentID)
		if studentID == "" {
			return c.Status(401).JSON(fiber.Map{"error": utils.ErrUnauthorized})
		}
		user.StudentID = studentID
	}

	attendance, err := services.UpdateAttendanceStatus(uint(attendanceID), req.Status, req.Notes, user.StudentID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message":    "Attendance status updated successfully",
		"attendance": attendance,
	})
}

// controller/event_controller.go
package controller

import (
	"attendance-system/connection"
	"attendance-system/models"
	"attendance-system/services"
	"attendance-system/utils"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// parseTaggedCourses parses tagged courses from form value if not provided in JSON
func parseTaggedCourses(req *models.EventRequest, c *fiber.Ctx) {
	if len(req.TaggedCourses) == 0 {
		if raw := c.FormValue("tagged_courses"); raw != "" {
			parts := strings.Split(raw, ",")
			for _, p := range parts {
				p = strings.ToUpper(strings.TrimSpace(p))
				if p != "" {
					req.TaggedCourses = append(req.TaggedCourses, p)
				}
			}
		}
	}
}

// getUserFromContext retrieves the user from context or fallback
func getUserFromContext(c *fiber.Ctx) (models.User, error) {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		// Fallback: try to get from query/header (for backward compatibility)
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
	return user, nil
}

// CreateEvent creates a new event
func CreateEvent(c *fiber.Ctx) error {
	req := new(models.EventRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	parseTaggedCourses(req, c)

	// Validate required fields
	if req.Title == "" || req.EventDate == "" || req.StartTime == "" || req.EndTime == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Title, event_date, start_time, and end_time are required"})
	}

	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	event, err := services.CreateEvent(*req, user.StudentID, user.Role)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Event created successfully",
		"event":   event,
	})
}

// GetEvent retrieves a single event
func GetEvent(c *fiber.Ctx) error {
	eventIDStr := c.Params("id")
	eventID, err := strconv.ParseUint(eventIDStr, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": utils.ErrInvalidEventID})
	}

	event, err := services.GetEvent(uint(eventID))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"event": event,
	})
}

// GetAllEvents retrieves all events with optional filters
func GetAllEvents(c *fiber.Ctx) error {
	filters := make(map[string]interface{})

	if course := c.Query("course"); course != "" {
		filters["course"] = course
	}
	if section := c.Query("section"); section != "" {
		filters["section"] = section
	}
	if yearLevel := c.Query("year_level"); yearLevel != "" {
		filters["year_level"] = yearLevel
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if isActive := c.Query("is_active"); isActive != "" {
		filters["is_active"] = isActive == "true"
	}

	events, err := services.GetAllEvents(filters)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"events": events,
		"count":  len(events),
	})
}

// GetMyEvents retrieves events for the current student
func GetMyEvents(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		studentID := c.Get("X-Student-ID")
		if studentID == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
		}
		user.StudentID = studentID
	}

	// Admins/faculty/superadmin should see all events (management view).
	var events []models.Event
	var err error
	if user.Role == "superadmin" || user.Role == "admin" || user.Role == "faculty" {
		events, err = services.GetAllEvents(map[string]interface{}{})
	} else {
		events, err = services.GetEventsByStudent(user.StudentID)
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"events": events,
		"count":  len(events),
	})
}

// UpdateEvent updates an event
func UpdateEvent(c *fiber.Ctx) error {
	eventIDStr := c.Params("id")
	eventID, err := strconv.ParseUint(eventIDStr, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": utils.ErrInvalidEventID})
	}

	req := new(models.EventRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request format"})
	}

	parseTaggedCourses(req, c)

	user, ok := c.Locals("user").(models.User)
	if !ok {
		studentID := c.Get(utils.HeaderStudentID)
		if studentID == "" {
			return c.Status(401).JSON(fiber.Map{"error": utils.ErrUnauthorized})
		}
		user.StudentID = studentID
	}

	event, err := services.UpdateEvent(uint(eventID), *req, user.StudentID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Event updated successfully",
		"event":   event,
	})
}

// DeleteEvent deletes an event (soft delete)
func DeleteEvent(c *fiber.Ctx) error {
	eventIDStr := c.Params("id")
	eventID, err := strconv.ParseUint(eventIDStr, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": utils.ErrInvalidEventID})
	}

	user, ok := c.Locals("user").(models.User)
	if !ok {
		studentID := c.Get(utils.HeaderStudentID)
		if studentID == "" {
			return c.Status(401).JSON(fiber.Map{"error": utils.ErrUnauthorized})
		}
		user.StudentID = studentID
	}

	if err := services.DeleteEvent(uint(eventID), user.StudentID); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Event deleted successfully",
	})
}

// GetEventCreationDropdowns returns predefined sections and departments for event creation
func GetEventCreationDropdowns(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"sections":    utils.Sections,
		"departments": utils.Departments,
	})
}

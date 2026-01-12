// services/event_service.go
package services

import (
	"attendance-system/connection"
	"attendance-system/logging"
	"attendance-system/models"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/skip2/go-qrcode"
	"go.uber.org/zap"
)

const (
	base64Prefix     = "data:image/png;base64,"
	errEventNotFound = "event not found"
	studentWhere     = "student_id = ?"
	sectionWhere     = "section = ?"
)

// CreateEvent creates a new event
func CreateEvent(req models.EventRequest, createdBy, createdByRole string) (*models.Event, error) {
	// Parse dates and produce start/end datetimes
	eventDate, startDateTime, endDateTime, err := parseEventDateTimes(req)
	if err != nil {
		return nil, err
	}

	// Log the parsed times for debugging
	logging.Logger.Info("Event times parsed",
		zap.String("request_start", req.StartTime),
		zap.String("request_end", req.EndTime),
		zap.Time("parsed_start", startDateTime),
		zap.Time("parsed_end", endDateTime),
	)

	// Generate QR code for event
	qrCodeBase64, err := generateEventQRCode(createdBy)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %v", err)
	}
	// Build event object
	event := &models.Event{
		Title:         req.Title,
		Description:   req.Description,
		EventDate:     eventDate,
		StartTime:     startDateTime,
		EndTime:       endDateTime,
		Location:      req.Location,
		Course:        req.Course,
		Section:       req.Section,
		YearLevel:     req.YearLevel,
		Department:    req.Department,
		College:       req.College,
		CreatedBy:     createdBy,
		CreatedByRole: createdByRole,
		Status:        "scheduled",
		IsActive:      true,
		QRCodeData:    qrCodeBase64,
	}

	// Normalize and set tagged courses (helper handles trimming/uppercasing)
	setTaggedCoursesFromRequest(event, req)

	// Ensure ID is zero so DB assigns it
	event.ID = 0

	// Persist event (handles duplicate-key sequence resync + retry)
	if err := persistEvent(event); err != nil {
		return nil, err
	}

	// If event has tagged courses, update student QR codes to event-specific
	if len(req.TaggedCourses) > 0 {
		go updateStudentQRCodesForEvent(event.ID, req.TaggedCourses, req.YearLevel, req.Section)
	} else if req.Course != "" && req.YearLevel != "" {
		// Fallback to single course for backward compatibility
		go updateStudentQRCodesForEvent(event.ID, []string{req.Course}, req.YearLevel, req.Section)
	}

	return event, nil
}

// persistEvent inserts an event while omitting client-provided ID, and retries once
// after resyncing the sequence if a duplicate-key error occurs.
func persistEvent(event *models.Event) error {
	if err := connection.DB.Omit("id").Create(event).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			// Resync sequence for events.id using pg_get_serial_sequence
			_ = connection.DB.Exec("SELECT setval(pg_get_serial_sequence('events','id'), (SELECT COALESCE(MAX(id),1) FROM events))")
			// Retry create once (omit id again explicitly)
			err2 := connection.DB.Omit("id").Create(event).Error
			if err2 == nil {
				return nil
			}
			return fmt.Errorf("failed to create event after sequence fix: %v", err2)
		}
		return fmt.Errorf("failed to create event: %v", err)
	}
	return nil
}

// parseEventDateTimes parses EventDate, StartTime and EndTime from request and returns
// the event date plus combined start and end datetimes.
// Supports both HH:MM format (local time) and ISO 8601 format
func parseEventDateTimes(req models.EventRequest) (time.Time, time.Time, time.Time, error) {
	eventDate, err := time.Parse("2006-01-02", req.EventDate)
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, errors.New("invalid event_date format. Use YYYY-MM-DD")
	}

	// Try parsing start_time as ISO 8601 first, then fallback to HH:MM
	var startDateTime time.Time
	if t, err := time.Parse("2006-01-02T15:04:05Z07:00", req.StartTime); err == nil {
		startDateTime = t
	} else if t, err := time.Parse("2006-01-02T15:04:05", req.StartTime); err == nil {
		startDateTime = t
	} else if t, err := time.Parse("15:04", req.StartTime); err == nil {
		// HH:MM format - combine with event date using local timezone
		startDateTime = time.Date(eventDate.Year(), eventDate.Month(), eventDate.Day(),
			t.Hour(), t.Minute(), 0, 0, time.Local)
	} else {
		return time.Time{}, time.Time{}, time.Time{}, errors.New("invalid start_time format. Use HH:MM, YYYY-MM-DDTHH:MM:SS, or ISO 8601")
	}

	// Try parsing end_time as ISO 8601 first, then fallback to HH:MM
	var endDateTime time.Time
	if t, err := time.Parse("2006-01-02T15:04:05Z07:00", req.EndTime); err == nil {
		endDateTime = t
	} else if t, err := time.Parse("2006-01-02T15:04:05", req.EndTime); err == nil {
		endDateTime = t
	} else if t, err := time.Parse("15:04", req.EndTime); err == nil {
		// HH:MM format - combine with event date using local timezone
		endDateTime = time.Date(eventDate.Year(), eventDate.Month(), eventDate.Day(),
			t.Hour(), t.Minute(), 0, 0, time.Local)
	} else {
		return time.Time{}, time.Time{}, time.Time{}, errors.New("invalid end_time format. Use HH:MM, YYYY-MM-DDTHH:MM:SS, or ISO 8601")
	}

	if endDateTime.Before(startDateTime) || endDateTime.Equal(startDateTime) {
		return time.Time{}, time.Time{}, time.Time{}, errors.New("end_time must be after start_time")
	}

	return eventDate, startDateTime, endDateTime, nil
}

// generateEventQRCode produces a base64 PNG QR code string for the event creator.
func generateEventQRCode(createdBy string) (string, error) {
	qrContent := fmt.Sprintf("event:%d:%s", time.Now().Unix(), createdBy)
	qrCodePNG, err := qrcode.Encode(qrContent, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}
	return base64Prefix + base64.StdEncoding.EncodeToString(qrCodePNG), nil
}

// GetEventQRCode returns the QR code for an event, generating and persisting one if missing.
func GetEventQRCode(eventID uint) (string, error) {
	var event models.Event
	if err := connection.DB.First(&event, eventID).Error; err != nil {
		return "", fmt.Errorf("event not found")
	}

	if event.QRCodeData != "" {
		return event.QRCodeData, nil
	}

	qrCodeBase64, err := generateEventQRCode(event.CreatedBy)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %v", err)
	}

	if err := connection.DB.Model(&event).Update("qr_code_data", qrCodeBase64).Error; err != nil {
		return "", fmt.Errorf("failed to save QR code: %v", err)
	}

	return qrCodeBase64, nil
}

// updateStudentQRCodesForEvent updates QR codes for students matching event criteria
func updateStudentQRCodesForEvent(eventID uint, courses []string, yearLevel, section string) {
	var students []models.User

	// Build query for multiple courses
	query := connection.DB.Where("course IN ?", courses)

	// Add year level filter if specified
	if yearLevel != "" {
		query = query.Where("year_level = ?", yearLevel)
	}

	// Add section filter if specified
	if section != "" {
		query = query.Where("section = ?", section)
	}

	if err := query.Find(&students).Error; err != nil {
		// Log error without exposing sensitive information
		return
	}

	// Get event details
	var event models.Event
	if err := connection.DB.First(&event, eventID).Error; err != nil {
		// Log error without exposing sensitive information
		return
	}

	// Generate event-specific QR code for each student
	for _, student := range students {
		// Skip if student already has an active event
		if student.ActiveEventID != nil {
			continue
		}

		// Save original QR code data
		originalQRCodeData := student.QRCodeData
		originalQRType := student.QRType

		// Generate event-specific QR code
		qrContent := fmt.Sprintf("event:%d:student:%s", eventID, student.StudentID)
		qrCodePNG, err := qrcode.Encode(qrContent, qrcode.Medium, 256)
		if err != nil {
			// Log error without exposing sensitive information
			continue
		}
		qrCodeBase64 := base64Prefix + base64.StdEncoding.EncodeToString(qrCodePNG)

		// Build updates map for event-specific QR code
		activeEventID := eventID
		updates := map[string]interface{}{
			"original_qr_code_data": originalQRCodeData,
			"original_qr_type":      originalQRType,
			"qr_code_data":          qrCodeBase64,
			"qr_type":               fmt.Sprintf("event:%d", eventID),
			"active_event_id":       activeEventID,
			"qr_generated_at":       time.Now(),
		}
		if err := connection.DB.Model(&models.User{}).Where(studentWhere, student.StudentID).Updates(updates).Error; err != nil {
			// Log error without exposing sensitive information
		}
	}
}

// RevertStudentQRCodesForEvent reverts QR codes back to student_id after event ends
func RevertStudentQRCodesForEvent(eventID uint) error {
	var students []models.User
	if err := connection.DB.Where("active_event_id = ?", eventID).Find(&students).Error; err != nil {
		return fmt.Errorf("error finding students for event: %v", err)
	}

	for _, student := range students {
		// Restore original QR code
		if student.OriginalQRCodeData != "" {
			student.QRCodeData = student.OriginalQRCodeData
			student.QRType = student.OriginalQRType
		} else {
			// If no original, generate new student_id QR code
			qrContent := fmt.Sprintf("student:%s", student.StudentID)
			qrCodePNG, err := qrcode.Encode(qrContent, qrcode.Medium, 256)
			if err != nil {
				// Log error without exposing sensitive information
				continue
			}
			qrCodeBase64 := base64Prefix + base64.StdEncoding.EncodeToString(qrCodePNG)
			student.QRCodeData = qrCodeBase64
			student.QRType = "student_id"
		}

		// Clear event-specific data
		now := time.Now()
		updates := map[string]interface{}{
			"qr_code_data":          student.QRCodeData,
			"qr_type":               student.QRType,
			"active_event_id":       nil,
			"original_qr_code_data": "",
			"original_qr_type":      "",
			"qr_generated_at":       now,
		}
		if err := connection.DB.Model(&models.User{}).Where(studentWhere, student.StudentID).Updates(updates).Error; err != nil {
			// Log error without exposing sensitive information
		}
	}

	return nil
}

// CheckAndUpdateCompletedEvents checks for completed events and reverts QR codes
func CheckAndUpdateCompletedEvents() {
	now := time.Now()
	var events []models.Event

	// Find events that have ended but status is still ongoing/scheduled
	if err := connection.DB.Where("end_time < ? AND status IN ? AND is_active = ?",
		now, []string{"scheduled", "ongoing"}, true).Find(&events).Error; err != nil {
		// Log error without exposing sensitive information
		return
	}

	for _, event := range events {
		// Update event status
		if err := connection.DB.Model(&models.Event{}).Where("id = ?", event.ID).Updates(map[string]interface{}{"status": "completed"}).Error; err != nil {
			// Log error without exposing sensitive information
			continue
		}

		// Revert student QR codes back to student_id
		if err := RevertStudentQRCodesForEvent(event.ID); err != nil {
			// Log error without exposing sensitive information
		}
	}

	// Also check for events that should be marked as "ongoing"
	var ongoingEvents []models.Event
	if err := connection.DB.Where("start_time <= ? AND end_time >= ? AND status = ? AND is_active = ?",
		now, now, "scheduled", true).Find(&ongoingEvents).Error; err == nil {
		for _, event := range ongoingEvents {
			if err := connection.DB.Model(&models.Event{}).Where("id = ?", event.ID).Updates(map[string]interface{}{"status": "ongoing"}).Error; err != nil {
				// Log error if needed
			}
		}
	}
}

// GetEvent retrieves an event by ID
func GetEvent(eventID uint) (*models.Event, error) {
	var event models.Event
	if err := connection.DB.First(&event, eventID).Error; err != nil {
		return nil, errors.New(errEventNotFound)
	}
	// Hide description until 24 hours before the event start_time
	if time.Now().Before(event.StartTime.Add(-24 * time.Hour)) {
		event.Description = ""
	}

	return &event, nil
}

// GetAllEvents retrieves all events with optional filters
func GetAllEvents(filters map[string]interface{}) ([]models.Event, error) {
	var events []models.Event
	query := connection.DB

	// Apply filters
	if course, ok := filters["course"].(string); ok && course != "" {
		query = query.Where("course = ?", course)
	}
	if section, ok := filters["section"].(string); ok && section != "" {
		query = query.Where(sectionWhere, section)
	}
	if yearLevel, ok := filters["year_level"].(string); ok && yearLevel != "" {
		query = query.Where("year_level = ?", yearLevel)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if isActive, ok := filters["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}

	if err := query.Order("event_date DESC, start_time DESC").Find(&events).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch events: %v", err)
	}

	// Hide description for events that are more than 24 hours away
	now := time.Now()
	for i := range events {
		if now.Before(events[i].StartTime.Add(-24 * time.Hour)) {
			events[i].Description = ""
		}
	}

	return events, nil
}

// UpdateEvent updates an existing event
func UpdateEvent(eventID uint, req models.EventRequest, updatedBy string) (*models.Event, error) {
	var event models.Event
	if err := connection.DB.First(&event, eventID).Error; err != nil {
		return nil, errors.New(errEventNotFound)
	}

	// Check permissions (only creator or admin can update)
	if err := ensureUpdatePermission(&event, updatedBy); err != nil {
		return nil, err
	}

	if err := applyEventUpdates(&event, req); err != nil {
		return nil, err
	}

	if err := connection.DB.Save(&event).Error; err != nil {
		return nil, fmt.Errorf("failed to update event: %v", err)
	}

	// If tagged courses were updated, revert existing QR codes and generate new ones
	if len(req.TaggedCourses) > 0 {
		go func() {
			// First revert existing QR codes
			RevertStudentQRCodesForEvent(eventID)
			// Then generate new ones for updated courses
			updateStudentQRCodesForEvent(eventID, req.TaggedCourses, req.YearLevel, req.Section)
		}()
	}

	return &event, nil
}

// ensureUpdatePermission returns nil when updatedBy is allowed to modify event.
func ensureUpdatePermission(event *models.Event, updatedBy string) error {
	if event.CreatedBy == updatedBy {
		return nil
	}
	var user models.User
	if err := connection.DB.Where(studentWhere, updatedBy).First(&user).Error; err != nil {
		return errors.New("unauthorized")
	}
	if user.Role != "superadmin" && user.Role != "admin" && user.Role != "faculty" {
		return errors.New("unauthorized: only event creator, admin, or faculty can update")
	}
	return nil
}

// applyEventUpdates applies provided non-empty fields from req to the event.
func applyEventUpdates(event *models.Event, req models.EventRequest) error {
	if err := applyTimeUpdates(event, req); err != nil {
		return err
	}
	applyOtherUpdates(event, req)
	return nil
}

func applyTimeUpdates(event *models.Event, req models.EventRequest) error {
	if req.EventDate != "" {
		if eventDate, err := time.Parse("2006-01-02", req.EventDate); err == nil {
			event.EventDate = eventDate
		}
	}

	if req.StartTime != "" {
		if startTime, err := time.Parse("15:04", req.StartTime); err == nil {
			event.StartTime = time.Date(event.EventDate.Year(), event.EventDate.Month(), event.EventDate.Day(),
				startTime.Hour(), startTime.Minute(), 0, 0, time.Local)
		}
	}

	if req.EndTime != "" {
		if endTime, err := time.Parse("15:04", req.EndTime); err == nil {
			event.EndTime = time.Date(event.EventDate.Year(), event.EventDate.Month(), event.EventDate.Day(),
				endTime.Hour(), endTime.Minute(), 0, 0, time.Local)
		}
	}

	return nil
}

func applyOtherUpdates(event *models.Event, req models.EventRequest) {
	if req.Title != "" {
		event.Title = req.Title
	}
	if req.Description != "" {
		event.Description = req.Description
	}
	if req.Location != "" {
		event.Location = req.Location
	}
	if req.Course != "" {
		event.Course = req.Course
	}
	if req.Section != "" {
		event.Section = req.Section
	}
	if req.YearLevel != "" {
		event.YearLevel = req.YearLevel
	}
	if req.Department != "" {
		event.Department = req.Department
	}
	if req.College != "" {
		event.College = req.College
	}
	// Update tagged courses if provided
	setTaggedCoursesFromRequest(event, req)
}

// setTaggedCoursesFromRequest normalizes and sets tagged courses on the event from the request.
func setTaggedCoursesFromRequest(event *models.Event, req models.EventRequest) {
	if len(req.TaggedCourses) == 0 {
		return
	}
	var cleaned []string
	for _, c := range req.TaggedCourses {
		c = strings.ToUpper(strings.TrimSpace(c))
		if c != "" {
			cleaned = append(cleaned, c)
		}
	}
	if len(cleaned) > 0 {
		event.TaggedCoursesCSV = strings.Join(cleaned, ",")
		event.TaggedCourses = cleaned
	} else {
		event.TaggedCoursesCSV = ""
		event.TaggedCourses = nil
	}
}

// DeleteEvent deletes an event (soft delete by setting is_active to false)
func DeleteEvent(eventID uint, deletedBy string) error {
	var event models.Event
	if err := connection.DB.First(&event, eventID).Error; err != nil {
		return errors.New(errEventNotFound)
	}

	// Check permissions
	if event.CreatedBy != deletedBy {
		var user models.User
		if err := connection.DB.Where(studentWhere, deletedBy).First(&user).Error; err != nil {
			return errors.New("unauthorized")
		}
		if user.Role != "superadmin" && user.Role != "admin" {
			return errors.New("unauthorized: only event creator or admin can delete")
		}
	}

	// Soft delete
	event.IsActive = false
	event.Status = "cancelled"
	if err := connection.DB.Save(&event).Error; err != nil {
		return fmt.Errorf("failed to delete event: %v", err)
	}

	// Revert student QR codes back to original when event is deleted
	go RevertStudentQRCodesForEvent(eventID)

	return nil
}

// GetEventsByStudent retrieves events relevant to a student
func GetEventsByStudent(studentID string) ([]models.Event, error) {
	var user models.User
	if err := connection.DB.Where(studentWhere, studentID).First(&user).Error; err != nil {
		return nil, errors.New("student not found")
	}

	// Return all active events but compute whether the student is allowed to enter
	var events []models.Event
	query := connection.DB.Where("is_active = ?", true)
	if err := query.Order("event_date DESC, start_time DESC").Find(&events).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch events: %v", err)
	}

	populateTaggedCoursesAndAllowed(events, &user)

	return events, nil
}

// populateTaggedCoursesAndAllowed fills transient TaggedCourses and Allowed fields
// for a slice of events given a user. This keeps the logic out of GetEventsByStudent
// and reduces its cognitive complexity.
func populateTaggedCoursesAndAllowed(events []models.Event, user *models.User) {
	for i := range events {
		events[i].TaggedCourses = parseTaggedCoursesCSV(events[i].TaggedCoursesCSV)
		events[i].Allowed = isUserAllowedForEvent(events[i], user)
	}
}

// parseTaggedCoursesCSV converts a CSV string to a normalized slice of course tags.
func parseTaggedCoursesCSV(csv string) []string {
	if csv == "" {
		return nil
	}
	parts := strings.Split(csv, ",")
	var trimmed []string
	for _, p := range parts {
		p = strings.ToUpper(strings.TrimSpace(p))
		if p != "" {
			trimmed = append(trimmed, p)
		}
	}
	if len(trimmed) == 0 {
		return nil
	}
	return trimmed
}

// isUserAllowedForEvent returns whether the given user may enter the event.
func isUserAllowedForEvent(event models.Event, user *models.User) bool {
	// No tags => open to all
	if len(event.TaggedCourses) == 0 {
		return true
	}
	// Student without course info cannot enter tagged events
	if user.Course == "" {
		return false
	}
	userCourse := strings.ToUpper(strings.TrimSpace(user.Course))
	for _, c := range event.TaggedCourses {
		if c == userCourse {
			return true
		}
	}
	return false
}

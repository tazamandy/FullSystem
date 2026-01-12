package services

import (
	"attendance-system/connection"
	"attendance-system/models"
	"attendance-system/utils"
	"fmt"
)

// CreateUserFromMap creates a new user from a generic map (used by admin tools).
func CreateUserFromMap(data map[string]interface{}) error {
	var user models.User

	if v, ok := data["student_id"].(string); ok {
		user.StudentID = utils.SanitizeStudentID(v)
	}
	if v, ok := data["email"].(string); ok {
		user.Email = utils.SanitizeEmail(v)
	}
	if v, ok := data["first_name"].(string); ok {
		user.FirstName = v
	}
	if v, ok := data["last_name"].(string); ok {
		user.LastName = v
	}
	if v, ok := data["role"].(string); ok {
		user.Role = v
	} else {
		user.Role = "student"
	}
	if v, ok := data["is_verified"].(bool); ok {
		user.IsVerified = v
	} else {
		user.IsVerified = false
	}

	if p, ok := data["password"].(string); ok && p != "" {
		hashed, err := utils.HashPassword(p)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = hashed
	} else {
		return fmt.Errorf("password is required")
	}

	// Ensure required fields
	if user.Email == "" {
		return fmt.Errorf("email is required")
	}

	// Prevent duplicate email
	var existing models.User
	if err := connection.DB.Where("email = ?", user.Email).First(&existing).Error; err == nil {
		return fmt.Errorf("email already exists")
	}

	// Create user
	if err := connection.DB.Omit("id").Create(&user).Error; err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil
}

// UpdateUserByID updates a user by student_id
func UpdateUserByID(studentID string, updates map[string]interface{}) error {
	if studentID == "" {
		return fmt.Errorf("student id is required")
	}
	if _, ok := updates["password"]; ok {
		if pw, ok2 := updates["password"].(string); ok2 && pw != "" {
			hashed, err := utils.HashPassword(pw)
			if err != nil {
				return fmt.Errorf("failed to hash password: %w", err)
			}
			updates["password"] = hashed
		} else {
			delete(updates, "password")
		}
	}
	if err := connection.DB.Model(&models.User{}).Where("student_id = ?", studentID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}
	return nil
}

// DeleteUserByStudentID deletes a user by student id
func DeleteUserByStudentID(studentID string) error {
	if studentID == "" {
		return fmt.Errorf("student id is required")
	}
	if err := connection.DB.Where("student_id = ?", studentID).Delete(&models.User{}).Error; err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}

// GetAllAttendance returns all attendance records (for superadmin)
func GetAllAttendance() ([]models.Attendance, error) {
	var attendances []models.Attendance
	if err := connection.DB.Preload("Event").Preload("Student").Order("marked_at DESC").Find(&attendances).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch attendance: %v", err)
	}
	return attendances, nil
}

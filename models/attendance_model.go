// models/attendance_model.go
package models

import "time"

// Attendance represents a single attendance record for a student in an event
type Attendance struct {
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	EventID   uint   `json:"event_id" gorm:"not null;index"`
	StudentID string `json:"student_id" gorm:"not null;type:varchar(255);index"`

	// Attendance details
	Status       string    `json:"status" gorm:"type:varchar(50);default:'present'"` // present, absent, late, excused
	MarkedAt     time.Time `json:"marked_at" gorm:"not null"`
	MarkedBy     string    `json:"marked_by" gorm:"type:varchar(255)"`     // StudentID of who marked (self or admin)
	MarkedByRole string    `json:"marked_by_role" gorm:"type:varchar(50)"` // student, admin, faculty

	// Method of marking
	Method string `json:"method" gorm:"type:varchar(50);default:'qr_scan'"` // qr_scan, manual, api

	// Location data (if available)
	Latitude  float64 `json:"latitude,omitempty" gorm:"type:decimal(10,8)"`
	Longitude float64 `json:"longitude,omitempty" gorm:"type:decimal(11,8)"`

	// Notes
	Notes string `json:"notes,omitempty" gorm:"type:text"`

	// Check-in/Check-out tracking
	CheckInTime    *time.Time `json:"check_in_time,omitempty" gorm:"type:timestamp"`
	CheckOutTime   *time.Time `json:"check_out_time,omitempty" gorm:"type:timestamp"`
	CheckInStatus  string     `json:"check_in_status,omitempty" gorm:"type:varchar(50)"`  // early, on_time, late
	CheckOutStatus string     `json:"check_out_status,omitempty" gorm:"type:varchar(50)"` // early, on_time, late

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Event   Event `json:"event,omitempty" gorm:"foreignKey:EventID"`
	Student User  `json:"student,omitempty" gorm:"foreignKey:StudentID;references:StudentID"`
}

// AttendanceRequest for marking attendance
type AttendanceRequest struct {
	EventID   uint    `json:"event_id"`
	StudentID string  `json:"student_id,omitempty"`
	Status    string  `json:"status,omitempty"`
	Method    string  `json:"method,omitempty"`
	Action    string  `json:"action"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Notes     string  `json:"notes,omitempty"`
}

// AttendanceStats for reporting
type AttendanceStats struct {
	TotalEvents    int     `json:"total_events"`
	PresentCount   int     `json:"present_count"`
	AbsentCount    int     `json:"absent_count"`
	LateCount      int     `json:"late_count"`
	ExcusedCount   int     `json:"excused_count"`
	AttendanceRate float64 `json:"attendance_rate"`
}

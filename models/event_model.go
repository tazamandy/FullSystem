// models/event_model.go
package models

import "time"

// Event represents a class, meeting, or activity that requires attendance
type Event struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string    `json:"title" gorm:"not null;type:varchar(255)"`
	Description string    `json:"description" gorm:"type:text"`
	EventDate   time.Time `json:"event_date" gorm:"not null"`
	StartTime   time.Time `json:"start_time" gorm:"not null"`
	EndTime     time.Time `json:"end_time" gorm:"not null"`
	Location    string    `json:"location" gorm:"type:varchar(255)"`
	Course      string    `json:"course" gorm:"type:varchar(100)"`
	Section     string    `json:"section" gorm:"type:varchar(50)"`
	YearLevel   string    `json:"year_level" gorm:"type:varchar(50)"`
	Department  string    `json:"department" gorm:"type:varchar(100)"`
	College     string    `json:"college" gorm:"type:varchar(100)"`

	// Event creator/owner
	CreatedBy     string `json:"created_by" gorm:"not null;type:varchar(255)"` // StudentID of creator
	CreatedByRole string `json:"created_by_role" gorm:"type:varchar(50);default:'faculty'"`

	// Status
	Status   string `json:"status" gorm:"type:varchar(50);default:'scheduled'"` // scheduled, ongoing, completed, cancelled
	IsActive bool   `json:"is_active" gorm:"default:true"`

	// QR Code for this event
	QRCodeData string `json:"qr_code_data,omitempty" gorm:"type:text"`
	// TaggedCoursesCSV stores comma-separated course tags (e.g., "CS,IT,EE").
	// Use `TaggedCourses` (transient) for JSON response.
	TaggedCoursesCSV string `json:"-" gorm:"type:text;column:tagged_courses"`
	// Transient fields (not persisted by GORM)
	TaggedCourses []string `json:"tagged_courses,omitempty" gorm:"-"`
	Allowed       bool     `json:"allowed,omitempty" gorm:"-"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Attendances []Attendance `json:"attendances,omitempty" gorm:"foreignKey:EventID"`
}

// EventRequest for creating/updating events
type EventRequest struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	EventDate      string `json:"event_date"` // ISO 8601 format
	StartTime      string `json:"start_time"` // ISO 8601 format
	EndTime        string `json:"end_time"`   // ISO 8601 format
	Location       string `json:"location"`
	Course         string `json:"course"`
	Section        string `json:"section"`
	YearLevel      string `json:"year_level"`
	Department     string `json:"department"`
	College        string `json:"college"`
	TaggedCourses []string `json:"tagged_courses,omitempty"`
}

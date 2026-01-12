// models/verification.go
package models

import "time"

type PendingUser struct {
	ID               uint   `gorm:"primaryKey;autoIncrement"`
	StudentID        string `gorm:"uniqueIndex;not null"`
	Email            string `gorm:"uniqueIndex;not null"`
	Password         string `gorm:"not null"`
	Username         string `gorm:"not null"`
	FirstName        string `gorm:"not null"`
	LastName         string `gorm:"not null"`
	MiddleName       string
	Course           string
	YearLevel        string
	Section          string
	Department       string
	College          string
	ContactNumber    string
	Address          string
	VerificationCode string    `gorm:"not null"`
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	ExpiresAt        time.Time
}

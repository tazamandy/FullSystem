// models/password_reset.go
package models

import "time"

type PasswordReset struct {
    ID        uint      `gorm:"primaryKey;autoIncrement"`
    Email     string    `gorm:"not null;index"`
    Code      string    `gorm:"not null"`
    ExpiresAt time.Time `gorm:"not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    Used      bool      `gorm:"default:false"`
}
// services/audit_service.go
package services

import (
	"attendance-system/connection"
	"attendance-system/logging"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AuditLog represents an audit trail entry
type AuditLog struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Action    string    `gorm:"not null;type:varchar(255)" json:"action"`   // e.g., "USER_PROMOTED", "PASSWORD_CHANGED", "EVENT_DELETED"
	ActorID   string    `gorm:"not null;type:varchar(255)" json:"actor_id"` // Who did the action
	TargetID  string    `gorm:"type:varchar(255)" json:"target_id"`         // Who/what it was done to
	Details   string    `gorm:"type:text" json:"details"`                   // JSON details of action
	IPAddress string    `gorm:"type:varchar(255)" json:"ip_address"`        // Request IP
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// LogAuditAction logs an admin action for security audit trail
func LogAuditAction(action, actorID, targetID, details, ipAddress string) error {
	if actorID == "" || action == "" {
		return fmt.Errorf("actor_id and action are required")
	}

	auditLog := AuditLog{
		Action:    action,
		ActorID:   actorID,
		TargetID:  targetID,
		Details:   details,
		IPAddress: ipAddress,
	}

	if err := connection.DB.Create(&auditLog).Error; err != nil {
		logging.Logger.Error("Failed to log audit action",
			zap.String("action", action),
			zap.String("actor_id", actorID),
			zap.Error(err),
		)
		return err
	}

	logging.Logger.Info("Audit action logged",
		zap.String("action", action),
		zap.String("actor_id", actorID),
		zap.String("target_id", targetID),
		zap.String("ip_address", ipAddress),
	)

	return nil
}

// GetAuditLogs retrieves audit logs with filtering
func GetAuditLogs(filters map[string]interface{}, limit int, offset int) ([]AuditLog, int64, error) {
	var auditLogs []AuditLog
	var total int64

	query := connection.DB
	if action, ok := filters["action"].(string); ok && action != "" {
		query = query.Where("action = ?", action)
	}
	if actorID, ok := filters["actor_id"].(string); ok && actorID != "" {
		query = query.Where("actor_id = ?", actorID)
	}
	if targetID, ok := filters["target_id"].(string); ok && targetID != "" {
		query = query.Where("target_id = ?", targetID)
	}
	if startDate, ok := filters["start_date"].(time.Time); ok && !startDate.IsZero() {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate, ok := filters["end_date"].(time.Time); ok && !endDate.IsZero() {
		query = query.Where("created_at <= ?", endDate)
	}

	if err := query.Model(&AuditLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit == 0 {
		limit = 50 // default limit
	}
	if limit > 500 {
		limit = 500 // max limit
	}

	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&auditLogs).Error; err != nil {
		return nil, 0, err
	}

	return auditLogs, total, nil
}

// Audit action constants
const (
	AuditUserPromoted       = "USER_PROMOTED"
	AuditPasswordChanged    = "PASSWORD_CHANGED"
	AuditPasswordReset      = "PASSWORD_RESET"
	AuditEventCreated       = "EVENT_CREATED"
	AuditEventUpdated       = "EVENT_UPDATED"
	AuditEventDeleted       = "EVENT_DELETED"
	AuditUserVerified       = "USER_VERIFIED"
	AuditUserRegistered     = "USER_REGISTERED"
	AuditAttendanceMarked   = "ATTENDANCE_MARKED"
	AuditAttendanceUpdated  = "ATTENDANCE_UPDATED"
	AuditAdminAccessAttempt = "ADMIN_ACCESS_ATTEMPT"
)

// TableName ensures the audit_logs table is used in queries
func (AuditLog) TableName() string {
	return "audit_logs"
}

// MigrateAuditLog creates the audit_logs table if it doesn't exist
func MigrateAuditLog() error {
	return connection.DB.AutoMigrate(&AuditLog{})
}

// DeleteOldAuditLogs removes audit logs older than specified duration
// This should be run periodically (e.g., monthly via cron job)
func DeleteOldAuditLogs(olderThanDays int) error {
	cutoffDate := time.Now().AddDate(0, 0, -olderThanDays)
	result := connection.DB.Where("created_at < ?", cutoffDate).Delete(&AuditLog{})
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	logging.Logger.Info("Old audit logs deleted",
		zap.Int64("rows_affected", result.RowsAffected),
		zap.Time("cutoff_date", cutoffDate),
	)

	return nil
}

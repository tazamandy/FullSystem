package services

import (
	"attendance-system/connection"
	"attendance-system/models"
	"fmt"
	"time"
)

var StartTime time.Time

func init() {
	StartTime = time.Now().UTC()
}

type SystemStats struct {
	TotalUsers     int64            `json:"total_users"`
	RoleCounts     map[string]int64 `json:"role_counts"`
	TotalEvents    int64            `json:"total_events"`
	UptimeSeconds  int64            `json:"uptime_seconds"`
	DBSizeBytes    int64            `json:"db_size_bytes"`
	DBSizeMB       float64          `json:"db_size_mb"`
	DBUsagePercent float64          `json:"db_usage_percent"`
}

// GetSystemStats returns aggregated system statistics for superadmin dashboards
func GetSystemStats() (*SystemStats, error) {
	var stats SystemStats

	// Total users
	if err := connection.DB.Model(&models.User{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	// Counts by role
	stats.RoleCounts = make(map[string]int64)
	type roleCountRow struct {
		Role  string
		Count int64
	}
	var rows []roleCountRow
	if err := connection.DB.Model(&models.User{}).Select("role, COUNT(*) as count").Group("role").Scan(&rows).Error; err == nil {
		for _, r := range rows {
			stats.RoleCounts[r.Role] = r.Count
		}
	}

	// Total events
	if err := connection.DB.Model(&models.Event{}).Count(&stats.TotalEvents).Error; err != nil {
		return nil, fmt.Errorf("failed to count events: %w", err)
	}

	// Uptime
	stats.UptimeSeconds = int64(time.Since(StartTime).Seconds())

	// Try to get DB size (Postgres). If it fails, leave defaults (0)
	var dbSizeBytes int64
	err := connection.DB.Raw("SELECT pg_database_size(current_database())").Scan(&dbSizeBytes).Error
	if err == nil {
		stats.DBSizeBytes = dbSizeBytes
		stats.DBSizeMB = float64(dbSizeBytes) / 1024.0 / 1024.0
		// If there's no configured max, we can't compute a meaningful percent. Leave 0.
		stats.DBUsagePercent = 0
	}

	return &stats, nil
}

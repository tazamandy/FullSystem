//seeder/seed.go

package seeder

import (
	"attendance-system/connection"
	"attendance-system/models"
	"attendance-system/utils"
	"log"
	"time"
)

func SeedSuperAdmin() {
	db := connection.DB

	var count int64
	db.Model(&models.User{}).Where("role = ?", "superadmin").Count(&count)

	if count > 0 {
		log.Println("SuperAdmin already exists. Skipping seed.")
		return
	}

	hashedPassword, err := utils.HashPassword("superadmin123")
	if err != nil {
		log.Println("Failed to hash password for superadmin:", err)
		return
	}

	superadmin := models.User{
		StudentID:  "SUPERADMIN",
		Username:   "superadmin",
		Email:      "superadmin@example.com",
		Password:   string(hashedPassword),
		Role:       "superadmin",
		IsVerified: true,
		FirstName:  "Super",
		LastName:   "Admin",
		CreatedAt:  time.Now(),
		VerifiedAt: time.Now(),
	}

	// Ensure ID is zero so DB assigns it
	superadmin.ID = 0
	if err := db.Omit("id").Create(&superadmin).Error; err != nil {
		log.Println("Failed to create SuperAdmin:", err)
		return
	}

	log.Println("Default SuperAdmin created successfully!")
}

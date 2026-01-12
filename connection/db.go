// connection/db.go
package connection

import (
	"log"
	"os"

	"attendance-system/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	godotenv.Load()

	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"

	// Disable automatic foreign key constraint creation during AutoMigrate.
	// This avoids invalid FK creation (e.g., referencing non-unique columns)
	// and gives more control over DB schema in production.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	// Auto migrate ALL tables
	// Note: AutoMigrate can cause issues with existing constraints.
	// Only uncomment if you're starting with a fresh database.
	// db.AutoMigrate(
	// 	&models.User{},
	// 	&models.PendingUser{},
	// 	&models.PasswordReset{},
	// 	&models.Event{},
	// 	&models.Attendance{},
	// )

	// Create audit_logs table if it doesn't exist
	if !db.Migrator().HasTable("audit_logs") {
		db.Exec(`
			CREATE TABLE audit_logs (
				id SERIAL PRIMARY KEY,
				action VARCHAR(255) NOT NULL,
				actor_id VARCHAR(255) NOT NULL,
				target_id VARCHAR(255),
				details TEXT,
				ip_address VARCHAR(255),
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)
		`)
		// Create indexes for audit_logs
		db.Exec("CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action)")
		db.Exec("CREATE INDEX IF NOT EXISTS idx_audit_logs_actor_id ON audit_logs(actor_id)")
		db.Exec("CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at)")
		db.Exec("CREATE INDEX IF NOT EXISTS idx_audit_logs_target_id ON audit_logs(target_id)")
		log.Println("audit_logs table created successfully!")
	}

	// Ensure `tagged_courses` column exists for older databases that may
	// not have been migrated to include the new field. Prefer GORM Migrator
	// which will use the current DB connection and respects permissions.
	if !db.Migrator().HasColumn(&models.Event{}, "tagged_courses") {
		// Try to add column using GORM migrator (field name)
		if err := db.Migrator().AddColumn(&models.Event{}, "TaggedCoursesCSV"); err != nil {
			log.Printf("Failed to add column tagged_courses: %v", err)
			// Fallback: attempt a safe raw ALTER TABLE (IF NOT EXISTS)
			if execErr := db.Exec("ALTER TABLE events ADD COLUMN IF NOT EXISTS tagged_courses text").Error; execErr != nil {
				log.Printf("Fallback ALTER TABLE failed: %v", execErr)
			}
		}
	}

	DB = db
	log.Println("Database connected!")
}

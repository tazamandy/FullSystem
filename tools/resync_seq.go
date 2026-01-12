package main

import (
	"attendance-system/connection"
	"log"
)

func main() {
	// Connect using same logic as app (reads .env)
	connection.Connect()
	db := connection.DB

	queries := []string{
		"SELECT setval(pg_get_serial_sequence('events','id'), (SELECT COALESCE(MAX(id),1) FROM events));",
		"SELECT setval(pg_get_serial_sequence('users','id'), (SELECT COALESCE(MAX(id),1) FROM users));",
	}

	for _, q := range queries {
		if err := db.Exec(q).Error; err != nil {
			log.Fatalf("failed to exec '%s': %v", q, err)
		} else {
			log.Println("executed:", q)
		}
	}

	log.Println("sequence resync completed")
}

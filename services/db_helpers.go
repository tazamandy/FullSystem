// services/db_helpers.go
package services

import "attendance-system/connection"

// CreateWithoutID inserts a record while omitting the `id` column so
// the database assigns the primary key via its sequence. This helps
// avoid accidental inserts that include an explicit `id` value.
func CreateWithoutID(model interface{}) error {
	return connection.DB.Omit("id").Create(model).Error
}

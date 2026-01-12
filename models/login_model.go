// models/login_model.go

package models

type LoginRequest struct {
	StudentID string `json:"student_id" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

// models/resetpassword_model.go

package models

type ResetPasswordRequest struct {
	Email       string `json:"email"`
	Code        string `json:"code"`
	NewPassword string `json:"new_password"`
}

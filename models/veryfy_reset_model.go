// models/veryfy_reset_model.go

package models

type Request struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
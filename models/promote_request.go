// models/promote_request.go

package models


type PromoteRequest struct {
	StudentID string `json:"student_id"`
	Role      string `json:"role"`
}

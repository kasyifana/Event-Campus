package response

import (
	"event-campus-backend/internal/domain"
	"time"

	"github.com/google/uuid"
)

// AttendanceResponse represents attendance data
type AttendanceResponse struct {
	ID             uuid.UUID `json:"id"`
	RegistrationID uuid.UUID `json:"registration_id"`
	UserName       string    `json:"user_name"`
	UserEmail      string    `json:"user_email"`
	MarkedAt       time.Time `json:"marked_at"`
	Notes          *string   `json:"notes,omitempty"`
}

// ToAttendanceResponse converts domain.Attendance to response
func ToAttendanceResponse(att *domain.Attendance) AttendanceResponse {
	resp := AttendanceResponse{
		ID:             att.ID,
		RegistrationID: att.RegistrationID,
		MarkedAt:       att.MarkedAt,
		Notes:          att.Notes,
	}

	if att.UserName != nil {
		resp.UserName = *att.UserName
	}

	if att.UserEmail != nil {
		resp.UserEmail = *att.UserEmail
	}

	return resp
}

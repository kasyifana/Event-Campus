package domain

import (
	"time"

	"github.com/google/uuid"
)

// Attendance represents attendance record for an event
type Attendance struct {
	ID             uuid.UUID `json:"id" db:"id"`
	EventID        uuid.UUID `json:"event_id" db:"event_id"`
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	RegistrationID uuid.UUID `json:"registration_id" db:"registration_id"`
	MarkedAt       time.Time `json:"marked_at" db:"marked_at"`
	MarkedBy       uuid.UUID `json:"marked_by" db:"marked_by"`
	Notes          *string   `json:"notes,omitempty" db:"notes"`

	// Additional fields for joined queries
	UserName   *string `json:"user_name,omitempty" db:"user_name"`
	UserEmail  *string `json:"user_email,omitempty" db:"user_email"`
	EventTitle *string `json:"event_title,omitempty" db:"event_title"`
}

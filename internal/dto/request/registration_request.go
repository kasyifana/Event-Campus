package request

import "github.com/google/uuid"

// RegisterEventRequest represents event registration request
type RegisterEventRequest struct {
	EventID uuid.UUID `json:"event_id" binding:"required"`
}

// CancelRegistrationRequest represents registration cancellation request
type CancelRegistrationRequest struct {
	RegistrationID uuid.UUID `json:"registration_id" binding:"required"`
}

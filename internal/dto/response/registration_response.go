package response

import (
	"event-campus-backend/internal/domain"
	"time"

	"github.com/google/uuid"
)

// RegistrationResponse represents registration data
type RegistrationResponse struct {
	ID           uuid.UUID  `json:"id"`
	EventID      uuid.UUID  `json:"event_id"`
	EventTitle   string     `json:"event_title"`
	EventDate    time.Time  `json:"event_date"`
	Status       string     `json:"status"`
	RegisteredAt time.Time  `json:"registered_at"`
	CancelledAt  *time.Time `json:"cancelled_at,omitempty"`
	ReminderSent bool       `json:"reminder_sent"`
}

// MyRegistrationsResponse represents user's registrations grouped by status
type MyRegistrationsResponse struct {
	Upcoming  []RegistrationResponse `json:"upcoming"`
	Waitlist  []RegistrationResponse `json:"waitlist"`
	Past      []RegistrationResponse `json:"past"`
	Cancelled []RegistrationResponse `json:"cancelled"`
}

// ToRegistrationResponse converts domain.Registration to RegistrationResponse
func ToRegistrationResponse(reg *domain.Registration) RegistrationResponse {
	resp := RegistrationResponse{
		ID:           reg.ID,
		EventID:      reg.EventID,
		Status:       reg.Status,
		RegisteredAt: reg.RegisteredAt,
		CancelledAt:  reg.CancelledAt,
		ReminderSent: reg.ReminderSent,
	}

	if reg.EventTitle != nil {
		resp.EventTitle = *reg.EventTitle
	}

	if reg.EventDate != nil {
		resp.EventDate = *reg.EventDate
	}

	return resp
}

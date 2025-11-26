package domain

import (
	"time"

	"github.com/google/uuid"
)

// Registration statuses
const (
	RegistrationStatusRegistered = "registered"
	RegistrationStatusWaitlist   = "waitlist"
	RegistrationStatusCancelled  = "cancelled"
	RegistrationStatusAttended   = "attended"
)

// Registration represents a user's registration to an event
type Registration struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	EventID      uuid.UUID  `json:"event_id" db:"event_id"`
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	Status       string     `json:"status" db:"status"`
	RegisteredAt time.Time  `json:"registered_at" db:"registered_at"`
	CancelledAt  *time.Time `json:"cancelled_at,omitempty" db:"cancelled_at"`
	ReminderSent bool       `json:"reminder_sent" db:"reminder_sent"`

	// Additional fields for joined queries
	EventTitle *string    `json:"event_title,omitempty" db:"event_title"`
	EventDate  *time.Time `json:"event_date,omitempty" db:"event_date"`
	UserName   *string    `json:"user_name,omitempty" db:"user_name"`
	UserEmail  *string    `json:"user_email,omitempty" db:"user_email"`
	UserPhone  *string    `json:"user_phone,omitempty" db:"user_phone"`
}

// IsRegistered checks if registration is in registered status
func (r *Registration) IsRegistered() bool {
	return r.Status == RegistrationStatusRegistered
}

// IsWaitlist checks if registration is in waitlist
func (r *Registration) IsWaitlist() bool {
	return r.Status == RegistrationStatusWaitlist
}

// IsCancelled checks if registration is cancelled
func (r *Registration) IsCancelled() bool {
	return r.Status == RegistrationStatusCancelled
}

// IsAttended checks if user has attended
func (r *Registration) IsAttended() bool {
	return r.Status == RegistrationStatusAttended
}

// CanCancel checks if registration can be cancelled
func (r *Registration) CanCancel() bool {
	return r.Status == RegistrationStatusRegistered || r.Status == RegistrationStatusWaitlist
}

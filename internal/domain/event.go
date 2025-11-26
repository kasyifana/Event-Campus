package domain

import (
	"time"

	"github.com/google/uuid"
)

// Event categories
const (
	CategorySeminar  = "seminar"
	CategoryWorkshop = "workshop"
	CategoryLomba    = "lomba"
	CategoryKonser   = "konser"
)

// Event types
const (
	EventTypeOnline  = "online"
	EventTypeOffline = "offline"
)

// Event statuses
const (
	StatusDraft     = "draft"
	StatusPublished = "published"
	StatusOngoing   = "ongoing"
	StatusCompleted = "completed"
	StatusCancelled = "cancelled"
)

// Event represents an event in the system
type Event struct {
	ID                   uuid.UUID `json:"id" db:"id"`
	OrganizerID          uuid.UUID `json:"organizer_id" db:"organizer_id"`
	Title                string    `json:"title" db:"title"`
	Description          string    `json:"description" db:"description"`
	Category             string    `json:"category" db:"category"`
	EventType            string    `json:"event_type" db:"event_type"`
	Location             *string   `json:"location,omitempty" db:"location"`
	ZoomLink             *string   `json:"zoom_link,omitempty" db:"zoom_link"`
	PosterPath           *string   `json:"poster_path,omitempty" db:"poster_path"`
	StartDate            time.Time `json:"start_date" db:"start_date"`
	EndDate              time.Time `json:"end_date" db:"end_date"`
	RegistrationDeadline time.Time `json:"registration_deadline" db:"registration_deadline"`
	MaxParticipants      int       `json:"max_participants" db:"max_participants"`
	CurrentParticipants  int       `json:"current_participants" db:"current_participants"`
	IsUIIOnly            bool      `json:"is_uii_only" db:"is_uii_only"`
	Status               string    `json:"status" db:"status"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`

	// Additional fields for joined queries
	OrganizerName *string `json:"organizer_name,omitempty" db:"organizer_name"`
}

// IsFull checks if event is at capacity
func (e *Event) IsFull() bool {
	return e.CurrentParticipants >= e.MaxParticipants
}

// HasCapacity checks if event has available slots
func (e *Event) HasCapacity() bool {
	return e.CurrentParticipants < e.MaxParticipants
}

// CanRegister checks if registration is still open
func (e *Event) CanRegister() bool {
	now := time.Now()
	return e.Status == StatusPublished &&
		now.Before(e.RegistrationDeadline) &&
		now.Before(e.StartDate)
}

// IsOnline checks if event is online
func (e *Event) IsOnline() bool {
	return e.EventType == EventTypeOnline
}

// IsOffline checks if event is offline
func (e *Event) IsOffline() bool {
	return e.EventType == EventTypeOffline
}

// HasStarted checks if event has started
func (e *Event) HasStarted() bool {
	return time.Now().After(e.StartDate) || time.Now().Equal(e.StartDate)
}

// HasEnded checks if event has ended
func (e *Event) HasEnded() bool {
	return time.Now().After(e.EndDate)
}

// AvailableSlots returns number of available slots
func (e *Event) AvailableSlots() int {
	return e.MaxParticipants - e.CurrentParticipants
}

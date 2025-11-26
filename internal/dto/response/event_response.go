package response

import (
	"event-campus-backend/internal/domain"
	"time"

	"github.com/google/uuid"
)

// EventResponse represents event data with organizer info
type EventResponse struct {
	ID                   uuid.UUID `json:"id"`
	OrganizerID          uuid.UUID `json:"organizer_id"`
	OrganizerName        string    `json:"organizer_name"`
	Title                string    `json:"title"`
	Description          string    `json:"description"`
	Category             string    `json:"category"`
	EventType            string    `json:"event_type"`
	Location             *string   `json:"location,omitempty"`
	ZoomLink             *string   `json:"zoom_link,omitempty"`
	PosterPath           *string   `json:"poster_path,omitempty"`
	PosterURL            *string   `json:"poster_url,omitempty"`
	StartDate            time.Time `json:"start_date"`
	EndDate              time.Time `json:"end_date"`
	RegistrationDeadline time.Time `json:"registration_deadline"`
	MaxParticipants      int       `json:"max_participants"`
	CurrentParticipants  int       `json:"current_participants"`
	AvailableSlots       int       `json:"available_slots"`
	IsUIIOnly            bool      `json:"is_uii_only"`
	Status               string    `json:"status"`
	IsFull               bool      `json:"is_full"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// EventListResponse represents list of events
type EventListResponse struct {
	Events []EventResponse `json:"events"`
}

// EventRegistrantResponse represents registrant info
type EventRegistrantResponse struct {
	RegistrationID uuid.UUID  `json:"registration_id"`
	UserID         uuid.UUID  `json:"user_id"`
	Name           string     `json:"name"`
	Email          string     `json:"email"`
	PhoneNumber    string     `json:"phone_number"`
	Status         string     `json:"status"`
	RegisteredAt   time.Time  `json:"registered_at"`
	CancelledAt    *time.Time `json:"cancelled_at,omitempty"`
}

// ToEventResponse converts domain.Event to EventResponse
func ToEventResponse(event *domain.Event, baseURL string) EventResponse {
	resp := EventResponse{
		ID:                   event.ID,
		OrganizerID:          event.OrganizerID,
		Title:                event.Title,
		Description:          event.Description,
		Category:             event.Category,
		EventType:            event.EventType,
		Location:             event.Location,
		ZoomLink:             event.ZoomLink,
		PosterPath:           event.PosterPath,
		StartDate:            event.StartDate,
		EndDate:              event.EndDate,
		RegistrationDeadline: event.RegistrationDeadline,
		MaxParticipants:      event.MaxParticipants,
		CurrentParticipants:  event.CurrentParticipants,
		AvailableSlots:       event.AvailableSlots(),
		IsUIIOnly:            event.IsUIIOnly,
		Status:               event.Status,
		IsFull:               event.IsFull(),
		CreatedAt:            event.CreatedAt,
		UpdatedAt:            event.UpdatedAt,
	}

	if event.OrganizerName != nil {
		resp.OrganizerName = *event.OrganizerName
	}

	// Generate poster URL if path exists
	if event.PosterPath != nil && *event.PosterPath != "" {
		posterURL := baseURL + "/files/" + *event.PosterPath
		resp.PosterURL = &posterURL
	}

	return resp
}

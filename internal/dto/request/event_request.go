package request

import "time"

// CreateEventRequest represents event creation request
type CreateEventRequest struct {
	Title                string    `json:"title" binding:"required"`
	Description          string    `json:"description" binding:"required"`
	Category             string    `json:"category" binding:"required,oneof=seminar workshop lomba konser"`
	EventType            string    `json:"event_type" binding:"required,oneof=online offline"`
	Location             *string   `json:"location,omitempty"`
	ZoomLink             *string   `json:"zoom_link,omitempty"`
	StartDate            time.Time `json:"start_date" binding:"required"`
	EndDate              time.Time `json:"end_date" binding:"required"`
	RegistrationDeadline time.Time `json:"registration_deadline" binding:"required"`
	MaxParticipants      int       `json:"max_participants" binding:"required,min=1"`
	IsUIIOnly            bool      `json:"is_uii_only"`
	Status               string    `json:"status" binding:"required,oneof=draft published"`
}

// UpdateEventRequest represents event update request
type UpdateEventRequest struct {
	Title                *string    `json:"title,omitempty"`
	Description          *string    `json:"description,omitempty"`
	Category             *string    `json:"category,omitempty" binding:"omitempty,oneof=seminar workshop lomba konser"`
	EventType            *string    `json:"event_type,omitempty" binding:"omitempty,oneof=online offline"`
	Location             *string    `json:"location,omitempty"`
	ZoomLink             *string    `json:"zoom_link,omitempty"`
	StartDate            *time.Time `json:"start_date,omitempty"`
	EndDate              *time.Time `json:"end_date,omitempty"`
	RegistrationDeadline *time.Time `json:"registration_deadline,omitempty"`
	MaxParticipants      *int       `json:"max_participants,omitempty" binding:"omitempty,min=1"`
	IsUIIOnly            *bool      `json:"is_uii_only,omitempty"`
	Status               *string    `json:"status,omitempty" binding:"omitempty,oneof=draft published ongoing completed cancelled"`
}

// EventFilterRequest represents event filtering parameters
type EventFilterRequest struct {
	Category  *string    `form:"category,omitempty" binding:"omitempty,oneof=seminar workshop lomba konser"`
	Status    *string    `form:"status,omitempty" binding:"omitempty,oneof=draft published ongoing completed cancelled"`
	IsUIIOnly *bool      `form:"is_uii_only,omitempty"`
	StartDate *time.Time `form:"start_date,omitempty"`
	EndDate   *time.Time `form:"end_date,omitempty"`
	Page      int        `form:"page,default=1" binding:"omitempty,min=1"`
	Limit     int        `form:"limit,default=20" binding:"omitempty,min=1,max=100"`
}

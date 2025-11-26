package usecase

import (
	"context"
	"event-campus-backend/internal/domain"
	"event-campus-backend/internal/dto/request"
	"event-campus-backend/internal/dto/response"
	"event-campus-backend/internal/repository"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// EventUsecase defines interface for event business logic
type EventUsecase interface {
	CreateEvent(ctx context.Context, organizerID uuid.UUID, req *request.CreateEventRequest, posterPath *string) (*domain.Event, error)
	GetEvent(ctx context.Context, id uuid.UUID) (*response.EventResponse, error)
	GetAllEvents(ctx context.Context, filters map[string]interface{}) ([]response.EventResponse, error)
	GetMyEvents(ctx context.Context, organizerID uuid.UUID) ([]response.EventResponse, error)
	UpdateEvent(ctx context.Context, organizerID uuid.UUID, eventID uuid.UUID, req *request.UpdateEventRequest, posterPath *string) error
	DeleteEvent(ctx context.Context, organizerID uuid.UUID, eventID uuid.UUID) error
	PublishEvent(ctx context.Context, organizerID uuid.UUID, eventID uuid.UUID) error
}

type eventUsecase struct {
	eventRepo repository.EventRepository
	userRepo  repository.UserRepository
	baseURL   string
}

// NewEventUsecase creates a new event usecase
func NewEventUsecase(
	eventRepo repository.EventRepository,
	userRepo repository.UserRepository,
	baseURL string,
) EventUsecase {
	return &eventUsecase{
		eventRepo: eventRepo,
		userRepo:  userRepo,
		baseURL:   baseURL,
	}
}

func (u *eventUsecase) CreateEvent(ctx context.Context, organizerID uuid.UUID, req *request.CreateEventRequest, posterPath *string) (*domain.Event, error) {
	// Validate dates
	if req.StartDate.Before(time.Now()) {
		return nil, fmt.Errorf("start date must be in the future")
	}

	if req.EndDate.Before(req.StartDate) {
		return nil, fmt.Errorf("end date must be after start date")
	}

	if req.RegistrationDeadline.After(req.StartDate) {
		return nil, fmt.Errorf("registration deadline must be before start date")
	}

	// Validate capacity
	if req.MaxParticipants <= 0 {
		return nil, fmt.Errorf("max participants must be greater than 0")
	}

	// Validate event type and location/zoom link
	if req.EventType == domain.EventTypeOffline && (req.Location == nil || *req.Location == "") {
		return nil, fmt.Errorf("location is required for offline events")
	}

	if req.EventType == domain.EventTypeOnline && (req.ZoomLink == nil || *req.ZoomLink == "") {
		return nil, fmt.Errorf("zoom link is required for online events")
	}

	// Create event
	event := &domain.Event{
		OrganizerID:          organizerID,
		Title:                req.Title,
		Description:          req.Description,
		Category:             req.Category,
		EventType:            req.EventType,
		Location:             req.Location,
		ZoomLink:             req.ZoomLink,
		PosterPath:           posterPath,
		StartDate:            req.StartDate,
		EndDate:              req.EndDate,
		RegistrationDeadline: req.RegistrationDeadline,
		MaxParticipants:      req.MaxParticipants,
		IsUIIOnly:            req.IsUIIOnly,
		Status:               domain.StatusDraft,
	}

	if err := u.eventRepo.Create(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return event, nil
}

func (u *eventUsecase) GetEvent(ctx context.Context, id uuid.UUID) (*response.EventResponse, error) {
	event, err := u.eventRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("event not found")
	}

	// Get organizer info
	organizer, err := u.userRepo.GetByID(ctx, event.OrganizerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get organizer: %w", err)
	}

	resp := response.ToEventResponse(event, u.baseURL)
	resp.OrganizerName = organizer.FullName

	return &resp, nil
}

func (u *eventUsecase) GetAllEvents(ctx context.Context, filters map[string]interface{}) ([]response.EventResponse, error) {
	// Only show published events unless admin/organizer
	if _, ok := filters["status"]; !ok {
		filters["status"] = domain.StatusPublished
	}

	events, err := u.eventRepo.GetAll(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	var responses []response.EventResponse
	for _, event := range events {
		resp := response.ToEventResponse(&event, u.baseURL)

		// Get organizer name
		organizer, err := u.userRepo.GetByID(ctx, event.OrganizerID)
		if err == nil {
			resp.OrganizerName = organizer.FullName
		}

		responses = append(responses, resp)
	}

	return responses, nil
}

func (u *eventUsecase) GetMyEvents(ctx context.Context, organizerID uuid.UUID) ([]response.EventResponse, error) {
	events, err := u.eventRepo.GetByOrganizer(ctx, organizerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	var responses []response.EventResponse
	for _, event := range events {
		responses = append(responses, response.ToEventResponse(&event, u.baseURL))
	}

	return responses, nil
}

func (u *eventUsecase) UpdateEvent(ctx context.Context, organizerID uuid.UUID, eventID uuid.UUID, req *request.UpdateEventRequest, posterPath *string) error {
	// Get event
	event, err := u.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return fmt.Errorf("event not found")
	}

	// Check ownership
	if event.OrganizerID != organizerID {
		return fmt.Errorf("you don't have permission to update this event")
	}

	// Can't update completed or cancelled events
	if event.Status == domain.StatusCompleted || event.Status == domain.StatusCancelled {
		return fmt.Errorf("cannot update completed or cancelled events")
	}

	// Validate dates if provided
	if req.StartDate != nil {
		if req.StartDate.Before(time.Now()) && event.Status == domain.StatusDraft {
			return fmt.Errorf("start date must be in the future")
		}
		event.StartDate = *req.StartDate
	}

	if req.EndDate != nil {
		if req.EndDate.Before(event.StartDate) {
			return fmt.Errorf("end date must be after start date")
		}
		event.EndDate = *req.EndDate
	}

	if req.RegistrationDeadline != nil {
		if req.RegistrationDeadline.After(event.StartDate) {
			return fmt.Errorf("registration deadline must be before start date")
		}
		event.RegistrationDeadline = *req.RegistrationDeadline
	}

	// Update fields
	if req.Title != nil {
		event.Title = *req.Title
	}
	if req.Description != nil {
		event.Description = *req.Description
	}
	if req.Category != nil {
		event.Category = *req.Category
	}
	if req.EventType != nil {
		event.EventType = *req.EventType
	}
	if req.Location != nil {
		event.Location = req.Location
	}
	if req.ZoomLink != nil {
		event.ZoomLink = req.ZoomLink
	}
	if posterPath != nil {
		event.PosterPath = posterPath
	}
	if req.MaxParticipants != nil {
		// Can't reduce below current participants
		if *req.MaxParticipants < event.CurrentParticipants {
			return fmt.Errorf("cannot reduce max participants below current participants count")
		}
		event.MaxParticipants = *req.MaxParticipants
	}
	if req.IsUIIOnly != nil {
		event.IsUIIOnly = *req.IsUIIOnly
	}

	// Update event
	if err := u.eventRepo.Update(ctx, event); err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	return nil
}

func (u *eventUsecase) DeleteEvent(ctx context.Context, organizerID uuid.UUID, eventID uuid.UUID) error {
	// Get event
	event, err := u.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return fmt.Errorf("event not found")
	}

	// Check ownership
	if event.OrganizerID != organizerID {
		return fmt.Errorf("you don't have permission to delete this event")
	}

	// Can only delete draft events or cancel published ones
	if event.Status != domain.StatusDraft {
		// Change status to cancelled instead of deleting
		if err := u.eventRepo.UpdateStatus(ctx, eventID, domain.StatusCancelled); err != nil {
			return fmt.Errorf("failed to cancel event: %w", err)
		}
		return nil
	}

	// Delete draft event
	if err := u.eventRepo.Delete(ctx, eventID); err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	return nil
}

func (u *eventUsecase) PublishEvent(ctx context.Context, organizerID uuid.UUID, eventID uuid.UUID) error {
	// Get event
	event, err := u.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return fmt.Errorf("event not found")
	}

	// Check ownership
	if event.OrganizerID != organizerID {
		return fmt.Errorf("you don't have permission to publish this event")
	}

	// Can only publish draft events
	if event.Status != domain.StatusDraft {
		return fmt.Errorf("event is not in draft status")
	}

	// Validate event has required data
	if event.PosterPath == nil || *event.PosterPath == "" {
		return fmt.Errorf("event must have a poster before publishing")
	}

	// Update status to published
	if err := u.eventRepo.UpdateStatus(ctx, eventID, domain.StatusPublished); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

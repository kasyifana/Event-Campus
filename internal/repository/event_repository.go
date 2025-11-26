package repository

import (
	"context"
	"database/sql"
	"event-campus-backend/internal/domain"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// EventRepository defines interface for event data access
type EventRepository interface {
	Create(ctx context.Context, event *domain.Event) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	GetAll(ctx context.Context, filters map[string]interface{}) ([]domain.Event, error)
	GetByOrganizer(ctx context.Context, organizerID uuid.UUID) ([]domain.Event, error)
	Update(ctx context.Context, event *domain.Event) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	IncrementParticipants(ctx context.Context, id uuid.UUID) error
	DecrementParticipants(ctx context.Context, id uuid.UUID) error
}

type eventRepository struct {
	db *sql.DB
}

// NewEventRepository creates a new event repository
func NewEventRepository(db *sql.DB) EventRepository {
	return &eventRepository{
		db: db,
	}
}

func (r *eventRepository) Create(ctx context.Context, event *domain.Event) error {
	// Generate ID if not set
	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}

	// Set timestamps
	now := time.Now()
	event.CreatedAt = now
	event.UpdatedAt = now

	// Set default values
	if event.Status == "" {
		event.Status = domain.StatusDraft
	}
	if event.CurrentParticipants == 0 {
		event.CurrentParticipants = 0
	}

	query := `
		INSERT INTO events (
			id, organizer_id, title, description, category, event_type,
			location, zoom_link, poster_path, start_date, end_date,
			registration_deadline, max_participants, current_participants,
			is_uii_only, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`

	_, err := r.db.ExecContext(ctx, query,
		event.ID,
		event.OrganizerID,
		event.Title,
		event.Description,
		event.Category,
		event.EventType,
		event.Location,
		event.ZoomLink,
		event.PosterPath,
		event.StartDate,
		event.EndDate,
		event.RegistrationDeadline,
		event.MaxParticipants,
		event.CurrentParticipants,
		event.IsUIIOnly,
		event.Status,
		event.CreatedAt,
		event.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	return nil
}

func (r *eventRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	query := `
		SELECT id, organizer_id, title, description, category, event_type,
		       location, zoom_link, poster_path, start_date, end_date,
		       registration_deadline, max_participants, current_participants,
		       is_uii_only, status, created_at, updated_at
		FROM events
		WHERE id = $1
	`

	var event domain.Event
	var location, zoomLink, posterPath sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&event.ID,
		&event.OrganizerID,
		&event.Title,
		&event.Description,
		&event.Category,
		&event.EventType,
		&location,
		&zoomLink,
		&posterPath,
		&event.StartDate,
		&event.EndDate,
		&event.RegistrationDeadline,
		&event.MaxParticipants,
		&event.CurrentParticipants,
		&event.IsUIIOnly,
		&event.Status,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("event not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	if location.Valid {
		s := location.String
		event.Location = &s
	}
	if zoomLink.Valid {
		s := zoomLink.String
		event.ZoomLink = &s
	}
	if posterPath.Valid {
		s := posterPath.String
		event.PosterPath = &s
	}

	return &event, nil
}

func (r *eventRepository) GetAll(ctx context.Context, filters map[string]interface{}) ([]domain.Event, error) {
	query := `
		SELECT id, organizer_id, title, description, category, event_type,
		       location, zoom_link, poster_path, start_date, end_date,
		       registration_deadline, max_participants, current_participants,
		       is_uii_only, status, created_at, updated_at
		FROM events
		WHERE 1=1
	`

	var args []interface{}
	argCount := 1

	// Apply filters
	if category, ok := filters["category"].(string); ok && category != "" {
		query += fmt.Sprintf(" AND category = $%d", argCount)
		args = append(args, category)
		argCount++
	}

	if status, ok := filters["status"].(string); ok && status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
		argCount++
	}

	if eventType, ok := filters["event_type"].(string); ok && eventType != "" {
		query += fmt.Sprintf(" AND event_type = $%d", argCount)
		args = append(args, eventType)
		argCount++
	}

	if search, ok := filters["search"].(string); ok && search != "" {
		query += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+search+"%")
		argCount++
	}

	query += " ORDER BY start_date DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var event domain.Event
		var location, zoomLink, posterPath sql.NullString

		err := rows.Scan(
			&event.ID,
			&event.OrganizerID,
			&event.Title,
			&event.Description,
			&event.Category,
			&event.EventType,
			&location,
			&zoomLink,
			&posterPath,
			&event.StartDate,
			&event.EndDate,
			&event.RegistrationDeadline,
			&event.MaxParticipants,
			&event.CurrentParticipants,
			&event.IsUIIOnly,
			&event.Status,
			&event.CreatedAt,
			&event.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}

		if location.Valid {
			s := location.String
			event.Location = &s
		}
		if zoomLink.Valid {
			s := zoomLink.String
			event.ZoomLink = &s
		}
		if posterPath.Valid {
			s := posterPath.String
			event.PosterPath = &s
		}

		events = append(events, event)
	}

	return events, nil
}

func (r *eventRepository) GetByOrganizer(ctx context.Context, organizerID uuid.UUID) ([]domain.Event, error) {
	query := `
		SELECT id, organizer_id, title, description, category, event_type,
		       location, zoom_link, poster_path, start_date, end_date,
		       registration_deadline, max_participants, current_participants,
		       is_uii_only, status, created_at, updated_at
		FROM events
		WHERE organizer_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, organizerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var event domain.Event
		var location, zoomLink, posterPath sql.NullString

		err := rows.Scan(
			&event.ID,
			&event.OrganizerID,
			&event.Title,
			&event.Description,
			&event.Category,
			&event.EventType,
			&location,
			&zoomLink,
			&posterPath,
			&event.StartDate,
			&event.EndDate,
			&event.RegistrationDeadline,
			&event.MaxParticipants,
			&event.CurrentParticipants,
			&event.IsUIIOnly,
			&event.Status,
			&event.CreatedAt,
			&event.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}

		if location.Valid {
			s := location.String
			event.Location = &s
		}
		if zoomLink.Valid {
			s := zoomLink.String
			event.ZoomLink = &s
		}
		if posterPath.Valid {
			s := posterPath.String
			event.PosterPath = &s
		}

		events = append(events, event)
	}

	return events, nil
}

func (r *eventRepository) Update(ctx context.Context, event *domain.Event) error {
	event.UpdatedAt = time.Now()

	query := `
		UPDATE events
		SET title = $1, description = $2, category = $3, event_type = $4,
		    location = $5, zoom_link = $6, poster_path = $7,
		    start_date = $8, end_date = $9, registration_deadline = $10,
		    max_participants = $11, is_uii_only = $12, status = $13,
		    updated_at = $14
		WHERE id = $15
	`

	result, err := r.db.ExecContext(ctx, query,
		event.Title,
		event.Description,
		event.Category,
		event.EventType,
		event.Location,
		event.ZoomLink,
		event.PosterPath,
		event.StartDate,
		event.EndDate,
		event.RegistrationDeadline,
		event.MaxParticipants,
		event.IsUIIOnly,
		event.Status,
		event.UpdatedAt,
		event.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("event not found")
	}

	return nil
}

func (r *eventRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM events WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("event not found")
	}

	return nil
}

func (r *eventRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `
		UPDATE events
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update event status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("event not found")
	}

	return nil
}

func (r *eventRepository) IncrementParticipants(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE events
		SET current_participants = current_participants + 1, updated_at = $1
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to increment participants: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("event not found")
	}

	return nil
}

func (r *eventRepository) DecrementParticipants(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE events
		SET current_participants = GREATEST(current_participants - 1, 0), updated_at = $1
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to decrement participants: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("event not found")
	}

	return nil
}

package repository

import (
	"context"
	"database/sql"
	"event-campus-backend/internal/domain"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// RegistrationRepository defines interface for registration data access
type RegistrationRepository interface {
	Create(ctx context.Context, registration *domain.Registration) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Registration, error)
	GetByUserAndEvent(ctx context.Context, userID, eventID uuid.UUID) (*domain.Registration, error)
	GetByEvent(ctx context.Context, eventID uuid.UUID, status string) ([]domain.Registration, error)
	GetByUser(ctx context.Context, userID uuid.UUID) ([]domain.Registration, error)
	Update(ctx context.Context, registration *domain.Registration) error
	Cancel(ctx context.Context, id uuid.UUID) error
	GetWaitlistByEvent(ctx context.Context, eventID uuid.UUID) ([]domain.Registration, error)
	PromoteFromWaitlist(ctx context.Context, eventID uuid.UUID) (*domain.Registration, error)
	CountByEventAndStatus(ctx context.Context, eventID uuid.UUID, status string) (int, error)
}

type registrationRepository struct {
	db *sql.DB
}

// NewRegistrationRepository creates a new registration repository
func NewRegistrationRepository(db *sql.DB) RegistrationRepository {
	return &registrationRepository{
		db: db,
	}
}

func (r *registrationRepository) Create(ctx context.Context, registration *domain.Registration) error {
	// Generate ID if not set
	if registration.ID == uuid.Nil {
		registration.ID = uuid.New()
	}

	// Set registered timestamp
	registration.RegisteredAt = time.Now()

	query := `
		INSERT INTO registrations (id, event_id, user_id, status, registered_at, reminder_sent)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		registration.ID,
		registration.EventID,
		registration.UserID,
		registration.Status,
		registration.RegisteredAt,
		registration.ReminderSent,
	)

	if err != nil {
		return fmt.Errorf("failed to create registration: %w", err)
	}

	return nil
}

func (r *registrationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Registration, error) {
	query := `
		SELECT id, event_id, user_id, status, registered_at, cancelled_at, reminder_sent
		FROM registrations
		WHERE id = $1
	`

	var registration domain.Registration
	var cancelledAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&registration.ID,
		&registration.EventID,
		&registration.UserID,
		&registration.Status,
		&registration.RegisteredAt,
		&cancelledAt,
		&registration.ReminderSent,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("registration not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get registration: %w", err)
	}

	if cancelledAt.Valid {
		registration.CancelledAt = &cancelledAt.Time
	}

	return &registration, nil
}

func (r *registrationRepository) GetByUserAndEvent(ctx context.Context, userID, eventID uuid.UUID) (*domain.Registration, error) {
	query := `
		SELECT id, event_id, user_id, status, registered_at, cancelled_at, reminder_sent
		FROM registrations
		WHERE user_id = $1 AND event_id = $2
		ORDER BY registered_at DESC
		LIMIT 1
	`

	var registration domain.Registration
	var cancelledAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, userID, eventID).Scan(
		&registration.ID,
		&registration.EventID,
		&registration.UserID,
		&registration.Status,
		&registration.RegisteredAt,
		&cancelledAt,
		&registration.ReminderSent,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Not found is not an error
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get registration: %w", err)
	}

	if cancelledAt.Valid {
		registration.CancelledAt = &cancelledAt.Time
	}

	return &registration, nil
}

func (r *registrationRepository) GetByEvent(ctx context.Context, eventID uuid.UUID, status string) ([]domain.Registration, error) {
	var query string
	var args []interface{}

	if status == "" {
		query = `
			SELECT id, event_id, user_id, status, registered_at, cancelled_at, reminder_sent
			FROM registrations
			WHERE event_id = $1
			ORDER BY registered_at ASC
		`
		args = append(args, eventID)
	} else {
		query = `
			SELECT id, event_id, user_id, status, registered_at, cancelled_at, reminder_sent
			FROM registrations
			WHERE event_id = $1 AND status = $2
			ORDER BY registered_at ASC
		`
		args = append(args, eventID, status)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get registrations: %w", err)
	}
	defer rows.Close()

	var registrations []domain.Registration
	for rows.Next() {
		var registration domain.Registration
		var cancelledAt sql.NullTime

		err := rows.Scan(
			&registration.ID,
			&registration.EventID,
			&registration.UserID,
			&registration.Status,
			&registration.RegisteredAt,
			&cancelledAt,
			&registration.ReminderSent,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan registration: %w", err)
		}

		if cancelledAt.Valid {
			registration.CancelledAt = &cancelledAt.Time
		}

		registrations = append(registrations, registration)
	}

	return registrations, nil
}

func (r *registrationRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]domain.Registration, error) {
	query := `
		SELECT id, event_id, user_id, status, registered_at, cancelled_at, reminder_sent
		FROM registrations
		WHERE user_id = $1
		ORDER BY registered_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get registrations: %w", err)
	}
	defer rows.Close()

	var registrations []domain.Registration
	for rows.Next() {
		var registration domain.Registration
		var cancelledAt sql.NullTime

		err := rows.Scan(
			&registration.ID,
			&registration.EventID,
			&registration.UserID,
			&registration.Status,
			&registration.RegisteredAt,
			&cancelledAt,
			&registration.ReminderSent,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan registration: %w", err)
		}

		if cancelledAt.Valid {
			registration.CancelledAt = &cancelledAt.Time
		}

		registrations = append(registrations, registration)
	}

	return registrations, nil
}

func (r *registrationRepository) Update(ctx context.Context, registration *domain.Registration) error {
	query := `
		UPDATE registrations
		SET status = $1, cancelled_at = $2, reminder_sent = $3
		WHERE id = $4
	`

	result, err := r.db.ExecContext(ctx, query,
		registration.Status,
		registration.CancelledAt,
		registration.ReminderSent,
		registration.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update registration: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("registration not found")
	}

	return nil
}

func (r *registrationRepository) Cancel(ctx context.Context, id uuid.UUID) error {
	now := time.Now()

	query := `
		UPDATE registrations
		SET status = $1, cancelled_at = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, domain.RegistrationStatusCancelled, now, id)
	if err != nil {
		return fmt.Errorf("failed to cancel registration: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("registration not found")
	}

	return nil
}

func (r *registrationRepository) GetWaitlistByEvent(ctx context.Context, eventID uuid.UUID) ([]domain.Registration, error) {
	return r.GetByEvent(ctx, eventID, domain.RegistrationStatusWaitlist)
}

func (r *registrationRepository) PromoteFromWaitlist(ctx context.Context, eventID uuid.UUID) (*domain.Registration, error) {
	// Get first person in waitlist (FIFO)
	waitlist, err := r.GetWaitlistByEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if len(waitlist) == 0 {
		return nil, nil // No one in waitlist
	}

	// Promote first person
	firstInLine := &waitlist[0]
	firstInLine.Status = domain.RegistrationStatusRegistered

	if err := r.Update(ctx, firstInLine); err != nil {
		return nil, fmt.Errorf("failed to promote from waitlist: %w", err)
	}

	return firstInLine, nil
}

func (r *registrationRepository) CountByEventAndStatus(ctx context.Context, eventID uuid.UUID, status string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM registrations
		WHERE event_id = $1 AND status = $2
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, eventID, status).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count registrations: %w", err)
	}

	return count, nil
}

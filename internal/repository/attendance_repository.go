package repository

import (
	"context"
	"database/sql"
	"event-campus-backend/internal/domain"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AttendanceRepository defines interface for attendance data access
type AttendanceRepository interface {
	Create(ctx context.Context, attendance *domain.Attendance) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Attendance, error)
	GetByEventAndUser(ctx context.Context, eventID, userID uuid.UUID) (*domain.Attendance, error)
	GetByEvent(ctx context.Context, eventID uuid.UUID) ([]domain.Attendance, error)
	Update(ctx context.Context, attendance *domain.Attendance) error
	BulkCreate(ctx context.Context, attendances []domain.Attendance) error
	CountByEvent(ctx context.Context, eventID uuid.UUID) (int, error)
}

type attendanceRepository struct {
	db *sql.DB
}

// NewAttendanceRepository creates a new attendance repository
func NewAttendanceRepository(db *sql.DB) AttendanceRepository {
	return &attendanceRepository{
		db: db,
	}
}

func (r *attendanceRepository) Create(ctx context.Context, attendance *domain.Attendance) error {
	// Generate ID if not set
	if attendance.ID == uuid.Nil {
		attendance.ID = uuid.New()
	}

	// Set marked timestamp
	attendance.MarkedAt = time.Now()

	query := `
		INSERT INTO attendances (id, event_id, user_id, registration_id, marked_at, marked_by, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		attendance.ID,
		attendance.EventID,
		attendance.UserID,
		attendance.RegistrationID,
		attendance.MarkedAt,
		attendance.MarkedBy,
		attendance.Notes,
	)

	if err != nil {
		return fmt.Errorf("failed to create attendance: %w", err)
	}

	return nil
}

func (r *attendanceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Attendance, error) {
	query := `
		SELECT id, event_id, user_id, registration_id, marked_at, marked_by, notes
		FROM attendances
		WHERE id = $1
	`

	var attendance domain.Attendance
	var notes sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&attendance.ID,
		&attendance.EventID,
		&attendance.UserID,
		&attendance.RegistrationID,
		&attendance.MarkedAt,
		&attendance.MarkedBy,
		&notes,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("attendance not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get attendance: %w", err)
	}

	if notes.Valid {
		s := notes.String
		attendance.Notes = &s
	}

	return &attendance, nil
}

func (r *attendanceRepository) GetByEventAndUser(ctx context.Context, eventID, userID uuid.UUID) (*domain.Attendance, error) {
	query := `
		SELECT id, event_id, user_id, registration_id, marked_at, marked_by, notes
		FROM attendances
		WHERE event_id = $1 AND user_id = $2
	`

	var attendance domain.Attendance
	var notes sql.NullString

	err := r.db.QueryRowContext(ctx, query, eventID, userID).Scan(
		&attendance.ID,
		&attendance.EventID,
		&attendance.UserID,
		&attendance.RegistrationID,
		&attendance.MarkedAt,
		&attendance.MarkedBy,
		&notes,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Not found is not an error
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get attendance: %w", err)
	}

	if notes.Valid {
		s := notes.String
		attendance.Notes = &s
	}

	return &attendance, nil
}

func (r *attendanceRepository) GetByEvent(ctx context.Context, eventID uuid.UUID) ([]domain.Attendance, error) {
	query := `
		SELECT id, event_id, user_id, registration_id, marked_at, marked_by, notes
		FROM attendances
		WHERE event_id = $1
		ORDER BY marked_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendances: %w", err)
	}
	defer rows.Close()

	var attendances []domain.Attendance
	for rows.Next() {
		var attendance domain.Attendance
		var notes sql.NullString

		err := rows.Scan(
			&attendance.ID,
			&attendance.EventID,
			&attendance.UserID,
			&attendance.RegistrationID,
			&attendance.MarkedAt,
			&attendance.MarkedBy,
			&notes,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan attendance: %w", err)
		}

		if notes.Valid {
			s := notes.String
			attendance.Notes = &s
		}

		attendances = append(attendances, attendance)
	}

	return attendances, nil
}

func (r *attendanceRepository) Update(ctx context.Context, attendance *domain.Attendance) error {
	query := `
		UPDATE attendances
		SET notes = $1, marked_at = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query,
		attendance.Notes,
		attendance.MarkedAt,
		attendance.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update attendance: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("attendance not found")
	}

	return nil
}

func (r *attendanceRepository) BulkCreate(ctx context.Context, attendances []domain.Attendance) error {
	if len(attendances) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO attendances (id, event_id, user_id, registration_id, marked_at, marked_by, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, attendance := range attendances {
		if attendance.ID == uuid.Nil {
			attendance.ID = uuid.New()
		}
		if attendance.MarkedAt.IsZero() {
			attendance.MarkedAt = time.Now()
		}

		_, err := stmt.ExecContext(ctx,
			attendance.ID,
			attendance.EventID,
			attendance.UserID,
			attendance.RegistrationID,
			attendance.MarkedAt,
			attendance.MarkedBy,
			attendance.Notes,
		)
		if err != nil {
			return fmt.Errorf("failed to insert attendance: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *attendanceRepository) CountByEvent(ctx context.Context, eventID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM attendances
		WHERE event_id = $1
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, eventID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count attendances: %w", err)
	}

	return count, nil
}

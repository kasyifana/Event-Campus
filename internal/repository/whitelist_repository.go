package repository

import (
	"context"
	"database/sql"
	"event-campus-backend/internal/domain"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// WhitelistRepository defines interface for whitelist data access
type WhitelistRepository interface {
	Create(ctx context.Context, request *domain.WhitelistRequest) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.WhitelistRequest, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.WhitelistRequest, error)
	GetPendingRequests(ctx context.Context) ([]domain.WhitelistRequest, error)
	GetAllRequests(ctx context.Context, status string) ([]domain.WhitelistRequest, error)
	Update(ctx context.Context, request *domain.WhitelistRequest) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, adminNotes string, reviewedBy uuid.UUID) error
}

type whitelistRepository struct {
	db *sql.DB
}

// NewWhitelistRepository creates a new whitelist repository
func NewWhitelistRepository(db *sql.DB) WhitelistRepository {
	return &whitelistRepository{
		db: db,
	}
}

func (r *whitelistRepository) Create(ctx context.Context, request *domain.WhitelistRequest) error {
	// Generate ID if not set
	if request.ID == uuid.Nil {
		request.ID = uuid.New()
	}

	// Set submitted timestamp
	request.SubmittedAt = time.Now()
	request.Status = domain.WhitelistStatusPending

	query := `
		INSERT INTO whitelist_requests (id, user_id, organization_name, document_path, status, submitted_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		request.ID,
		request.UserID,
		request.OrganizationName,
		request.DocumentPath,
		request.Status,
		request.SubmittedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create whitelist request: %w", err)
	}

	return nil
}

func (r *whitelistRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.WhitelistRequest, error) {
	query := `
		SELECT id, user_id, organization_name, document_path, status, admin_notes, 
		       submitted_at, reviewed_at, reviewed_by
		FROM whitelist_requests
		WHERE id = $1
	`

	var request domain.WhitelistRequest
	var reviewedAt sql.NullTime
	var reviewedBy uuid.NullUUID
	var adminNotes sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&request.ID,
		&request.UserID,
		&request.OrganizationName,
		&request.DocumentPath,
		&request.Status,
		&adminNotes,
		&request.SubmittedAt,
		&reviewedAt,
		&reviewedBy,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("whitelist request not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get whitelist request: %w", err)
	}

	if adminNotes.Valid {
		s := adminNotes.String
		request.AdminNotes = &s
	}
	if reviewedAt.Valid {
		request.ReviewedAt = &reviewedAt.Time
	}
	if reviewedBy.Valid {
		request.ReviewedBy = &reviewedBy.UUID
	}

	return &request, nil
}

func (r *whitelistRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.WhitelistRequest, error) {
	query := `
		SELECT id, user_id, organization_name, document_path, status, admin_notes, 
		       submitted_at, reviewed_at, reviewed_by
		FROM whitelist_requests
		WHERE user_id = $1
		ORDER BY submitted_at DESC
		LIMIT 1
	`

	var request domain.WhitelistRequest
	var reviewedAt sql.NullTime
	var reviewedBy uuid.NullUUID
	var adminNotes sql.NullString

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&request.ID,
		&request.UserID,
		&request.OrganizationName,
		&request.DocumentPath,
		&request.Status,
		&adminNotes,
		&request.SubmittedAt,
		&reviewedAt,
		&reviewedBy,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No request found is not an error
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get whitelist request: %w", err)
	}

	if adminNotes.Valid {
		s := adminNotes.String
		request.AdminNotes = &s
	}
	if reviewedAt.Valid {
		request.ReviewedAt = &reviewedAt.Time
	}
	if reviewedBy.Valid {
		request.ReviewedBy = &reviewedBy.UUID
	}

	return &request, nil
}

func (r *whitelistRepository) GetPendingRequests(ctx context.Context) ([]domain.WhitelistRequest, error) {
	return r.GetAllRequests(ctx, domain.WhitelistStatusPending)
}

func (r *whitelistRepository) GetAllRequests(ctx context.Context, status string) ([]domain.WhitelistRequest, error) {
	var query string
	var args []interface{}

	if status == "" {
		query = `
			SELECT id, user_id, organization_name, document_path, status, admin_notes, 
			       submitted_at, reviewed_at, reviewed_by
			FROM whitelist_requests
			ORDER BY submitted_at DESC
		`
	} else {
		query = `
			SELECT id, user_id, organization_name, document_path, status, admin_notes, 
			       submitted_at, reviewed_at, reviewed_by
			FROM whitelist_requests
			WHERE status = $1
			ORDER BY submitted_at DESC
		`
		args = append(args, status)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get whitelist requests: %w", err)
	}
	defer rows.Close()

	var requests []domain.WhitelistRequest
	for rows.Next() {
		var request domain.WhitelistRequest
		var reviewedAt sql.NullTime
		var reviewedBy uuid.NullUUID
		var adminNotes sql.NullString

		err := rows.Scan(
			&request.ID,
			&request.UserID,
			&request.OrganizationName,
			&request.DocumentPath,
			&request.Status,
			&adminNotes,
			&request.SubmittedAt,
			&reviewedAt,
			&reviewedBy,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan whitelist request: %w", err)
		}

		if adminNotes.Valid {
			s := adminNotes.String
			request.AdminNotes = &s
		}
		if reviewedAt.Valid {
			request.ReviewedAt = &reviewedAt.Time
		}
		if reviewedBy.Valid {
			request.ReviewedBy = &reviewedBy.UUID
		}

		requests = append(requests, request)
	}

	return requests, nil
}

func (r *whitelistRepository) Update(ctx context.Context, request *domain.WhitelistRequest) error {
	query := `
		UPDATE whitelist_requests
		SET organization_name = $1, document_path = $2, status = $3, 
		    admin_notes = $4, reviewed_at = $5, reviewed_by = $6
		WHERE id = $7
	`

	result, err := r.db.ExecContext(ctx, query,
		request.OrganizationName,
		request.DocumentPath,
		request.Status,
		request.AdminNotes,
		request.ReviewedAt,
		request.ReviewedBy,
		request.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update whitelist request: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("whitelist request not found")
	}

	return nil
}

func (r *whitelistRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, adminNotes string, reviewedBy uuid.UUID) error {
	now := time.Now()

	query := `
		UPDATE whitelist_requests
		SET status = $1, admin_notes = $2, reviewed_at = $3, reviewed_by = $4
		WHERE id = $5
	`

	result, err := r.db.ExecContext(ctx, query, status, adminNotes, now, reviewedBy, id)
	if err != nil {
		return fmt.Errorf("failed to update whitelist status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("whitelist request not found")
	}

	return nil
}

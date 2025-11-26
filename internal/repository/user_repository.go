package repository

import (
	"context"
	"database/sql"
	"event-campus-backend/internal/domain"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// UserRepository defines interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	UpdateRole(ctx context.Context, userID uuid.UUID, role string, isApproved bool) error
}

// postgresUserRepository implements UserRepository with PostgreSQL
type postgresUserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository with PostgreSQL
func NewUserRepository(db *sql.DB) UserRepository {
	return &postgresUserRepository{
		db: db,
	}
}

func (r *postgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	// Check if email already exists
	existingUser, _ := r.GetByEmail(ctx, user.Email)
	if existingUser != nil {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}

	// Generate ID if not set
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Insert into database
	query := `
		INSERT INTO users (id, email, password_hash, full_name, phone_number, role, is_uii_civitas, is_approved, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.PhoneNumber,
		user.Role,
		user.IsUIICivitas,
		user.IsApproved,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *postgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, phone_number, role, is_uii_civitas, is_approved, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.PhoneNumber,
		&user.Role,
		&user.IsUIICivitas,
		&user.IsApproved,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *postgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, phone_number, role, is_uii_civitas, is_approved, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.PhoneNumber,
		&user.Role,
		&user.IsUIICivitas,
		&user.IsApproved,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *postgresUserRepository) Update(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users
		SET email = $1, password_hash = $2, full_name = $3, phone_number = $4, 
		    role = $5, is_uii_civitas = $6, is_approved = $7, updated_at = $8
		WHERE id = $9
	`

	result, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.PhoneNumber,
		user.Role,
		user.IsUIICivitas,
		user.IsApproved,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *postgresUserRepository) UpdateRole(ctx context.Context, userID uuid.UUID, role string, isApproved bool) error {
	query := `
		UPDATE users
		SET role = $1, is_approved = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := r.db.ExecContext(ctx, query, role, isApproved, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

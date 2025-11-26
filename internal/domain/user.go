package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// User roles
const (
	RoleMahasiswa  = "mahasiswa"
	RoleOrganisasi = "organisasi"
	RoleAdmin      = "admin"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	FullName     string    `json:"full_name" db:"full_name"`
	PhoneNumber  string    `json:"phone_number" db:"phone_number"`
	Role         string    `json:"role" db:"role"`
	IsUIICivitas bool      `json:"is_uii_civitas" db:"is_uii_civitas"`
	IsApproved   bool      `json:"is_approved" db:"is_approved"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// IsUIIEmail checks if email is from UII domain
func IsUIIEmail(email string) bool {
	return strings.HasSuffix(strings.ToLower(email), "uii.ac.id")
}

// IsMahasiswa checks if user is a mahasiswa
func (u *User) IsMahasiswa() bool {
	return u.Role == RoleMahasiswa
}

// IsOrganisasi checks if user is an organisasi
func (u *User) IsOrganisasi() bool {
	return u.Role == RoleOrganisasi
}

// IsAdmin checks if user is an admin
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// CanCreateEvent checks if user can create events
func (u *User) CanCreateEvent() bool {
	return (u.Role == RoleOrganisasi && u.IsApproved) || u.Role == RoleAdmin
}

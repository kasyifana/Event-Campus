package domain

import (
	"time"

	"github.com/google/uuid"
)

// Whitelist request statuses
const (
	WhitelistStatusPending  = "pending"
	WhitelistStatusApproved = "approved"
	WhitelistStatusRejected = "rejected"
)

// WhitelistRequest represents a request to become an organisasi
type WhitelistRequest struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	UserID           uuid.UUID  `json:"user_id" db:"user_id"`
	OrganizationName string     `json:"organization_name" db:"organization_name"`
	DocumentPath     string     `json:"document_path" db:"document_path"`
	Status           string     `json:"status" db:"status"`
	AdminNotes       *string    `json:"admin_notes,omitempty" db:"admin_notes"`
	SubmittedAt      time.Time  `json:"submitted_at" db:"submitted_at"`
	ReviewedAt       *time.Time `json:"reviewed_at,omitempty" db:"reviewed_at"`
	ReviewedBy       *uuid.UUID `json:"reviewed_by,omitempty" db:"reviewed_by"`

	// Additional fields for joined queries
	UserName     *string `json:"user_name,omitempty" db:"user_name"`
	UserEmail    *string `json:"user_email,omitempty" db:"user_email"`
	ReviewerName *string `json:"reviewer_name,omitempty" db:"reviewer_name"`
}

// IsPending checks if request is pending
func (w *WhitelistRequest) IsPending() bool {
	return w.Status == WhitelistStatusPending
}

// IsApproved checks if request is approved
func (w *WhitelistRequest) IsApproved() bool {
	return w.Status == WhitelistStatusApproved
}

// IsRejected checks if request is rejected
func (w *WhitelistRequest) IsRejected() bool {
	return w.Status == WhitelistStatusRejected
}

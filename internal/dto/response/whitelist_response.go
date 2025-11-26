package response

import (
	"event-campus-backend/internal/domain"
	"time"

	"github.com/google/uuid"
)

// WhitelistRequestResponse represents whitelist request data
type WhitelistRequestResponse struct {
	ID               uuid.UUID  `json:"id"`
	UserID           uuid.UUID  `json:"user_id"`
	UserName         string     `json:"user_name"`
	UserEmail        string     `json:"user_email"`
	OrganizationName string     `json:"organization_name"`
	DocumentPath     string     `json:"document_path"`
	DocumentURL      string     `json:"document_url"`
	Status           string     `json:"status"`
	AdminNotes       *string    `json:"admin_notes,omitempty"`
	SubmittedAt      time.Time  `json:"submitted_at"`
	ReviewedAt       *time.Time `json:"reviewed_at,omitempty"`
	ReviewedBy       *uuid.UUID `json:"reviewed_by,omitempty"`
	ReviewerName     *string    `json:"reviewer_name,omitempty"`
}

// ToWhitelistRequestResponse converts domain.WhitelistRequest to response
func ToWhitelistRequestResponse(req *domain.WhitelistRequest, baseURL string) WhitelistRequestResponse {
	resp := WhitelistRequestResponse{
		ID:               req.ID,
		UserID:           req.UserID,
		OrganizationName: req.OrganizationName,
		DocumentPath:     req.DocumentPath,
		DocumentURL:      baseURL + "/files/" + req.DocumentPath,
		Status:           req.Status,
		AdminNotes:       req.AdminNotes,
		SubmittedAt:      req.SubmittedAt,
		ReviewedAt:       req.ReviewedAt,
		ReviewedBy:       req.ReviewedBy,
	}

	if req.UserName != nil {
		resp.UserName = *req.UserName
	}

	if req.UserEmail != nil {
		resp.UserEmail = *req.UserEmail
	}

	if req.ReviewerName != nil {
		resp.ReviewerName = req.ReviewerName
	}

	return resp
}

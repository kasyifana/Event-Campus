package response

import (
	"event-campus-backend/internal/domain"
	"time"

	"github.com/google/uuid"
)

// LoginResponse represents login response
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// UserResponse represents sanitized user data
type UserResponse struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	PhoneNumber  string    `json:"phone_number"`
	Role         string    `json:"role"`
	IsUIICivitas bool      `json:"is_uii_civitas"`
	IsApproved   bool      `json:"is_approved"`
	CreatedAt    time.Time `json:"created_at"`
}

// ToUserResponse converts domain.User to UserResponse
func ToUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:           user.ID,
		Email:        user.Email,
		FullName:     user.FullName,
		PhoneNumber:  user.PhoneNumber,
		Role:         user.Role,
		IsUIICivitas: user.IsUIICivitas,
		IsApproved:   user.IsApproved,
		CreatedAt:    user.CreatedAt,
	}
}

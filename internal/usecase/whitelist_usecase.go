package usecase

import (
	"context"
	"event-campus-backend/internal/domain"
	"event-campus-backend/internal/dto/request"
	"event-campus-backend/internal/dto/response"
	"event-campus-backend/internal/repository"
	"event-campus-backend/internal/utils"
	"fmt"

	"github.com/google/uuid"
)

// WhitelistUsecase defines interface for whitelist business logic
type WhitelistUsecase interface {
	SubmitRequest(ctx context.Context, userID uuid.UUID, req *request.SubmitWhitelistRequest, documentPath string) error
	GetMyRequest(ctx context.Context, userID uuid.UUID) (*response.WhitelistRequestResponse, error)
	GetAllRequests(ctx context.Context, status string) ([]response.WhitelistRequestResponse, error)
	GetPendingRequests(ctx context.Context) ([]response.WhitelistRequestResponse, error)
	ReviewRequest(ctx context.Context, requestID uuid.UUID, reviewerID uuid.UUID, req *request.ReviewWhitelistRequest) error
}

type whitelistUsecase struct {
	whitelistRepo repository.WhitelistRepository
	userRepo      repository.UserRepository
	emailSender   *utils.EmailSender
	baseURL       string
}

// NewWhitelistUsecase creates a new whitelist usecase
func NewWhitelistUsecase(
	whitelistRepo repository.WhitelistRepository,
	userRepo repository.UserRepository,
	emailSender *utils.EmailSender,
	baseURL string,
) WhitelistUsecase {
	return &whitelistUsecase{
		whitelistRepo: whitelistRepo,
		userRepo:      userRepo,
		emailSender:   emailSender,
		baseURL:       baseURL,
	}
}

func (u *whitelistUsecase) SubmitRequest(ctx context.Context, userID uuid.UUID, req *request.SubmitWhitelistRequest, documentPath string) error {
	// Get user
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Check if user is mahasiswa
	if !user.IsMahasiswa() {
		return fmt.Errorf("only mahasiswa can submit whitelist request")
	}

	// Check if user already has pending request
	existingRequest, err := u.whitelistRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check existing request: %w", err)
	}

	if existingRequest != nil && existingRequest.Status == domain.WhitelistStatusPending {
		return fmt.Errorf("you already have a pending request")
	}

	// Create whitelist request
	whitelistRequest := &domain.WhitelistRequest{
		UserID:           userID,
		OrganizationName: req.OrganizationName,
		DocumentPath:     documentPath,
		Status:           domain.WhitelistStatusPending,
	}

	if err := u.whitelistRepo.Create(ctx, whitelistRequest); err != nil {
		return fmt.Errorf("failed to create whitelist request: %w", err)
	}

	return nil
}

func (u *whitelistUsecase) GetMyRequest(ctx context.Context, userID uuid.UUID) (*response.WhitelistRequestResponse, error) {
	request, err := u.whitelistRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get request: %w", err)
	}

	if request == nil {
		return nil, nil
	}

	// Get user info
	user, err := u.userRepo.GetByID(ctx, request.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	resp := response.ToWhitelistRequestResponse(request, u.baseURL)
	resp.UserName = user.FullName
	resp.UserEmail = user.Email

	return &resp, nil
}

func (u *whitelistUsecase) GetAllRequests(ctx context.Context, status string) ([]response.WhitelistRequestResponse, error) {
	requests, err := u.whitelistRepo.GetAllRequests(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get requests: %w", err)
	}

	var responses []response.WhitelistRequestResponse
	for _, req := range requests {
		user, err := u.userRepo.GetByID(ctx, req.UserID)
		if err != nil {
			continue // Skip if user not found
		}

		resp := response.ToWhitelistRequestResponse(&req, u.baseURL)
		resp.UserName = user.FullName
		resp.UserEmail = user.Email

		responses = append(responses, resp)
	}

	return responses, nil
}

func (u *whitelistUsecase) GetPendingRequests(ctx context.Context) ([]response.WhitelistRequestResponse, error) {
	return u.GetAllRequests(ctx, domain.WhitelistStatusPending)
}

func (u *whitelistUsecase) ReviewRequest(ctx context.Context, requestID uuid.UUID, reviewerID uuid.UUID, req *request.ReviewWhitelistRequest) error {
	// Get whitelist request
	whitelistRequest, err := u.whitelistRepo.GetByID(ctx, requestID)
	if err != nil {
		return fmt.Errorf("request not found")
	}

	// Check if already reviewed
	if whitelistRequest.Status != domain.WhitelistStatusPending {
		return fmt.Errorf("request has already been reviewed")
	}

	// Get user
	user, err := u.userRepo.GetByID(ctx, whitelistRequest.UserID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Determine status
	var newStatus string
	if req.Approved {
		newStatus = domain.WhitelistStatusApproved
	} else {
		newStatus = domain.WhitelistStatusRejected
	}

	// Update whitelist request status
	adminNotesStr := ""
	if req.AdminNotes != nil {
		adminNotesStr = *req.AdminNotes
	}

	if err := u.whitelistRepo.UpdateStatus(ctx, requestID, newStatus, adminNotesStr, reviewerID); err != nil {
		return fmt.Errorf("failed to update request: %w", err)
	}

	// If approved, update user role to organisasi
	if req.Approved {
		if err := u.userRepo.UpdateRole(ctx, user.ID, domain.RoleOrganisasi, true); err != nil {
			return fmt.Errorf("failed to update user role: %w", err)
		}

		// Send approval email
		if u.emailSender != nil {
			if err := u.emailSender.SendWhitelistApproval(user.Email, user.FullName, whitelistRequest.OrganizationName); err != nil {
				// Log error but don't fail
				fmt.Printf("Failed to send approval email: %v\n", err)
			}
		}
	} else {
		// Send rejection email
		if u.emailSender != nil {
			if err := u.emailSender.SendWhitelistRejection(user.Email, user.FullName, whitelistRequest.OrganizationName); err != nil {
				// Log error but don't fail
				fmt.Printf("Failed to send rejection email: %v\n", err)
			}
		}
	}

	return nil
}

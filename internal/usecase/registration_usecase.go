package usecase

import (
	"context"
	"event-campus-backend/internal/domain"
	"event-campus-backend/internal/repository"
	"event-campus-backend/internal/utils"
	"fmt"

	"github.com/google/uuid"
)

// RegistrationUsecase defines interface for registration business logic
type RegistrationUsecase interface {
	RegisterForEvent(ctx context.Context, userID, eventID uuid.UUID) error
	CancelRegistration(ctx context.Context, userID, registrationID uuid.UUID) error
	GetMyRegistrations(ctx context.Context, userID uuid.UUID) ([]domain.Registration, error)
	GetEventRegistrations(ctx context.Context, organizerID, eventID uuid.UUID) ([]domain.Registration, error)
}

type registrationUsecase struct {
	registrationRepo repository.RegistrationRepository
	eventRepo        repository.EventRepository
	userRepo         repository.UserRepository
	emailSender      *utils.EmailSender
}

// NewRegistrationUsecase creates a new registration usecase
func NewRegistrationUsecase(
	registrationRepo repository.RegistrationRepository,
	eventRepo repository.EventRepository,
	userRepo repository.UserRepository,
	emailSender *utils.EmailSender,
) RegistrationUsecase {
	return &registrationUsecase{
		registrationRepo: registrationRepo,
		eventRepo:        eventRepo,
		userRepo:         userRepo,
		emailSender:      emailSender,
	}
}

func (u *registrationUsecase) RegisterForEvent(ctx context.Context, userID, eventID uuid.UUID) error {
	// Get event
	event, err := u.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return fmt.Errorf("event not found")
	}

	// Check if event can accept registrations
	if !event.CanRegister() {
		return fmt.Errorf("registration is closed for this event")
	}

	// Get user
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Check if event is UII only
	if event.IsUIIOnly && !user.IsUIICivitas {
		return fmt.Errorf("this event is only for UII civitas")
	}

	// Check if already registered
	existingReg, err := u.registrationRepo.GetByUserAndEvent(ctx, userID, eventID)
	if err != nil {
		return fmt.Errorf("failed to check existing registration: %w", err)
	}

	if existingReg != nil && !existingReg.IsCancelled() {
		if existingReg.IsRegistered() {
			return fmt.Errorf("you are already registered for this event")
		}
		if existingReg.IsWaitlist() {
			return fmt.Errorf("you are already in the waitlist for this event")
		}
	}

	// Determine registration status based on capacity
	status := domain.RegistrationStatusRegistered
	if event.IsFull() {
		status = domain.RegistrationStatusWaitlist
	}

	// Create registration
	registration := &domain.Registration{
		EventID:      eventID,
		UserID:       userID,
		Status:       status,
		ReminderSent: false,
	}

	if err := u.registrationRepo.Create(ctx, registration); err != nil {
		return fmt.Errorf("failed to create registration: %w", err)
	}

	// Increment participant count if registered (not waitlisted)
	if status == domain.RegistrationStatusRegistered {
		if err := u.eventRepo.IncrementParticipants(ctx, eventID); err != nil {
			return fmt.Errorf("failed to update participant count: %w", err)
		}

		// Send confirmation email
		if u.emailSender != nil {
			zoomLink := ""
			if event.ZoomLink != nil {
				zoomLink = *event.ZoomLink
			}
			if err := u.emailSender.SendRegistrationConfirmation(user.Email, user.FullName, event.Title, event.StartDate, zoomLink); err != nil {
				// Log error but don't fail
				fmt.Printf("Failed to send confirmation email: %v\n", err)
			}
		}
	} else {
		// Send waitlist notification
		if u.emailSender != nil {
			waitlistCount, _ := u.registrationRepo.CountByEventAndStatus(ctx, eventID, domain.RegistrationStatusWaitlist)
			if err := u.emailSender.SendWaitlistNotification(user.Email, user.FullName, event.Title, waitlistCount); err != nil {
				// Log error but don't fail
				fmt.Printf("Failed to send waitlist notification: %v\n", err)
			}
		}
	}

	return nil
}

func (u *registrationUsecase) CancelRegistration(ctx context.Context, userID, registrationID uuid.UUID) error {
	// Get registration
	registration, err := u.registrationRepo.GetByID(ctx, registrationID)
	if err != nil {
		return fmt.Errorf("registration not found")
	}

	// Check ownership
	if registration.UserID != userID {
		return fmt.Errorf("you don't have permission to cancel this registration")
	}

	// Check if can cancel
	if !registration.CanCancel() {
		return fmt.Errorf("cannot cancel this registration")
	}

	// Get event
	event, err := u.eventRepo.GetByID(ctx, registration.EventID)
	if err != nil {
		return fmt.Errorf("event not found")
	}

	// Get user
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	wasRegistered := registration.IsRegistered()

	// Cancel registration
	if err := u.registrationRepo.Cancel(ctx, registrationID); err != nil {
		return fmt.Errorf("failed to cancel registration: %w", err)
	}

	// If was registered (not waitlisted), decrement count and promote from waitlist
	if wasRegistered {
		if err := u.eventRepo.DecrementParticipants(ctx, registration.EventID); err != nil {
			return fmt.Errorf("failed to update participant count: %w", err)
		}

		// Promote first person from waitlist
		promoted, err := u.registrationRepo.PromoteFromWaitlist(ctx, registration.EventID)
		if err != nil {
			fmt.Printf("Failed to promote from waitlist: %v\n", err)
		}

		if promoted != nil {
			// Increment participant count for promoted person
			if err := u.eventRepo.IncrementParticipants(ctx, registration.EventID); err != nil {
				fmt.Printf("Failed to increment participants for promoted registration: %v\n", err)
			}

			// Send promotion email
			promotedUser, err := u.userRepo.GetByID(ctx, promoted.UserID)
			if err == nil && u.emailSender != nil {
				zoomLink := ""
				if event.ZoomLink != nil {
					zoomLink = *event.ZoomLink
				}
				if err := u.emailSender.SendWaitlistPromotion(promotedUser.Email, promotedUser.FullName, event.Title, event.StartDate, zoomLink); err != nil {
					fmt.Printf("Failed to send promotion email: %v\n", err)
				}
			}
		}
	}

	// Send cancellation email
	if u.emailSender != nil {
		if err := u.emailSender.SendCancellationConfirmation(user.Email, user.FullName, event.Title); err != nil {
			fmt.Printf("Failed to send cancellation email: %v\n", err)
		}
	}

	return nil
}

func (u *registrationUsecase) GetMyRegistrations(ctx context.Context, userID uuid.UUID) ([]domain.Registration, error) {
	registrations, err := u.registrationRepo.GetByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get registrations: %w", err)
	}

	return registrations, nil
}

func (u *registrationUsecase) GetEventRegistrations(ctx context.Context, organizerID, eventID uuid.UUID) ([]domain.Registration, error) {
	// Get event
	event, err := u.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("event not found")
	}

	// Check ownership
	if event.OrganizerID != organizerID {
		return nil, fmt.Errorf("you don't have permission to view registrations for this event")
	}

	registrations, err := u.registrationRepo.GetByEvent(ctx, eventID, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get registrations: %w", err)
	}

	return registrations, nil
}

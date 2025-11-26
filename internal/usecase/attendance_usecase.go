package usecase

import (
	"context"
	"event-campus-backend/internal/domain"
	"event-campus-backend/internal/repository"
	"fmt"

	"github.com/google/uuid"
)

// AttendanceUsecase defines interface for attendance business logic
type AttendanceUsecase interface {
	MarkAttendance(ctx context.Context, organizerID, eventID, userID uuid.UUID, notes *string) error
	BulkMarkAttendance(ctx context.Context, organizerID, eventID uuid.UUID, userIDs []uuid.UUID) error
	GetEventAttendance(ctx context.Context, organizerID, eventID uuid.UUID) ([]domain.Attendance, error)
}

type attendanceUsecase struct {
	attendanceRepo   repository.AttendanceRepository
	eventRepo        repository.EventRepository
	registrationRepo repository.RegistrationRepository
	userRepo         repository.UserRepository
}

// NewAttendanceUsecase creates a new attendance usecase
func NewAttendanceUsecase(
	attendanceRepo repository.AttendanceRepository,
	eventRepo repository.EventRepository,
	registrationRepo repository.RegistrationRepository,
	userRepo repository.UserRepository,
) AttendanceUsecase {
	return &attendanceUsecase{
		attendanceRepo:   attendanceRepo,
		eventRepo:        eventRepo,
		registrationRepo: registrationRepo,
		userRepo:         userRepo,
	}
}

func (u *attendanceUsecase) MarkAttendance(ctx context.Context, organizerID, eventID, userID uuid.UUID, notes *string) error {
	// Get event
	event, err := u.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return fmt.Errorf("event not found")
	}

	// Check ownership
	if event.OrganizerID != organizerID {
		return fmt.Errorf("you don't have permission to mark attendance for this event")
	}

	// Check if event has started
	if !event.HasStarted() {
		return fmt.Errorf("cannot mark attendance before event starts")
	}

	// Get registration
	registration, err := u.registrationRepo.GetByUserAndEvent(ctx, userID, eventID)
	if err != nil {
		return fmt.Errorf("failed to get registration: %w", err)
	}

	if registration == nil {
		return fmt.Errorf("user is not registered for this event")
	}

	if !registration.IsRegistered() {
		return fmt.Errorf("user registration is not active")
	}

	// Check if already marked
	existingAttendance, err := u.attendanceRepo.GetByEventAndUser(ctx, eventID, userID)
	if err != nil {
		return fmt.Errorf("failed to check existing attendance: %w", err)
	}

	if existingAttendance != nil {
		return fmt.Errorf("attendance already marked for this user")
	}

	// Create attendance record
	attendance := &domain.Attendance{
		EventID:        eventID,
		UserID:         userID,
		RegistrationID: registration.ID,
		MarkedBy:       organizerID,
		Notes:          notes,
	}

	if err := u.attendanceRepo.Create(ctx, attendance); err != nil {
		return fmt.Errorf("failed to mark attendance: %w", err)
	}

	// Update registration status to attended
	registration.Status = domain.RegistrationStatusAttended
	if err := u.registrationRepo.Update(ctx, registration); err != nil {
		// Log error but don't fail
		fmt.Printf("Failed to update registration status: %v\n", err)
	}

	return nil
}

func (u *attendanceUsecase) BulkMarkAttendance(ctx context.Context, organizerID, eventID uuid.UUID, userIDs []uuid.UUID) error {
	// Get event
	event, err := u.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return fmt.Errorf("event not found")
	}

	// Check ownership
	if event.OrganizerID != organizerID {
		return fmt.Errorf("you don't have permission to mark attendance for this event")
	}

	// Check if event has started
	if !event.HasStarted() {
		return fmt.Errorf("cannot mark attendance before event starts")
	}

	var attendances []domain.Attendance
	var registrationsToUpdate []domain.Registration

	for _, userID := range userIDs {
		// Get registration
		registration, err := u.registrationRepo.GetByUserAndEvent(ctx, userID, eventID)
		if err != nil || registration == nil || !registration.IsRegistered() {
			// Skip invalid registrations
			continue
		}

		// Check if already marked
		existingAttendance, _ := u.attendanceRepo.GetByEventAndUser(ctx, eventID, userID)
		if existingAttendance != nil {
			// Skip already marked
			continue
		}

		// Add to batch
		attendance := domain.Attendance{
			ID:             uuid.New(),
			EventID:        eventID,
			UserID:         userID,
			RegistrationID: registration.ID,
			MarkedBy:       organizerID,
		}
		attendances = append(attendances, attendance)

		// Update registration status
		registration.Status = domain.RegistrationStatusAttended
		registrationsToUpdate = append(registrationsToUpdate, *registration)
	}

	if len(attendances) == 0 {
		return fmt.Errorf("no valid attendances to mark")
	}

	// Bulk create attendances
	if err := u.attendanceRepo.BulkCreate(ctx, attendances); err != nil {
		return fmt.Errorf("failed to bulk mark attendance: %w", err)
	}

	// Update registration statuses
	for _, reg := range registrationsToUpdate {
		if err := u.registrationRepo.Update(ctx, &reg); err != nil {
			fmt.Printf("Failed to update registration status for user %s: %v\n", reg.UserID, err)
		}
	}

	return nil
}

func (u *attendanceUsecase) GetEventAttendance(ctx context.Context, organizerID, eventID uuid.UUID) ([]domain.Attendance, error) {
	// Get event
	event, err := u.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("event not found")
	}

	// Check ownership
	if event.OrganizerID != organizerID {
		return nil, fmt.Errorf("you don't have permission to view attendance for this event")
	}

	attendances, err := u.attendanceRepo.GetByEvent(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendances: %w", err)
	}

	return attendances, nil
}

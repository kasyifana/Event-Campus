package scheduler

import (
	"context"
	"event-campus-backend/internal/domain"
	"event-campus-backend/internal/repository"
	"event-campus-backend/internal/utils"
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

// Scheduler manages automated tasks
type Scheduler struct {
	cron             *cron.Cron
	eventRepo        repository.EventRepository
	registrationRepo repository.RegistrationRepository
	userRepo         repository.UserRepository
	emailSender      *utils.EmailSender
}

// NewScheduler creates a new scheduler
func NewScheduler(
	eventRepo repository.EventRepository,
	registrationRepo repository.RegistrationRepository,
	userRepo repository.UserRepository,
	emailSender *utils.EmailSender,
) *Scheduler {
	return &Scheduler{
		cron:             cron.New(),
		eventRepo:        eventRepo,
		registrationRepo: registrationRepo,
		userRepo:         userRepo,
		emailSender:      emailSender,
	}
}

// Start starts all scheduled tasks
func (s *Scheduler) Start() error {
	// Run H-1 reminder every day at 09:00 AM
	_, err := s.cron.AddFunc("0 9 * * *", s.SendH1Reminders)
	if err != nil {
		return fmt.Errorf("failed to add H-1 reminder job: %w", err)
	}

	// Run event status updater every hour
	_, err = s.cron.AddFunc("0 * * * *", s.UpdateEventStatuses)
	if err != nil {
		return fmt.Errorf("failed to add status updater job: %w", err)
	}

	s.cron.Start()
	log.Println("✅ Scheduler started successfully")
	log.Println("  - H-1 Reminder: Daily at 09:00 AM")
	log.Println("  - Event Status Updater: Hourly")

	return nil
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	s.cron.Stop()
	log.Println("Scheduler stopped")
}

// SendH1Reminders sends reminder emails for events starting tomorrow
func (s *Scheduler) SendH1Reminders() {
	ctx := context.Background()
	log.Println("🔔 Running H-1 reminder job...")

	// Get all published events
	events, err := s.eventRepo.GetAll(ctx, map[string]interface{}{
		"status": domain.StatusPublished,
	})
	if err != nil {
		log.Printf("Failed to get events for reminders: %v", err)
		return
	}

	remindersSent := 0

	for _, event := range events {
		// Check if event starts tomorrow (within 24-48 hours)
		hoursUntilStart := time.Until(event.StartDate).Hours()
		if hoursUntilStart < 24 || hoursUntilStart >= 48 {
			continue // Not tomorrow
		}

		// Get registered participants
		registrations, err := s.registrationRepo.GetByEvent(ctx, event.ID, domain.RegistrationStatusRegistered)
		if err != nil {
			log.Printf("Failed to get registrations for event %s: %v", event.Title, err)
			continue
		}

		// Send reminder to each participant
		for _, reg := range registrations {
			// Skip if reminder already sent
			if reg.ReminderSent {
				continue
			}

			// Get user
			user, err := s.userRepo.GetByID(ctx, reg.UserID)
			if err != nil {
				log.Printf("Failed to get user %s: %v", reg.UserID, err)
				continue
			}

			// Prepare zoom link (reveal on H-1)
			zoomLink := ""
			if event.ZoomLink != nil && *event.ZoomLink != "" {
				zoomLink = *event.ZoomLink
			}

			location := ""
			if event.Location != nil {
				location = *event.Location
			}

			// Send reminder email
			if s.emailSender != nil {
				err := s.emailSender.SendReminderEmail(
					user.Email,
					user.FullName,
					event.Title,
					event.StartDate,
					location,
					&zoomLink,
					reg.ID.String(),
				)
				if err != nil {
					log.Printf("Failed to send reminder to %s: %v", user.Email, err)
					continue
				}
			}

			// Mark reminder as sent
			reg.ReminderSent = true
			if err := s.registrationRepo.Update(ctx, &reg); err != nil {
				log.Printf("Failed to update registration reminder status: %v", err)
			}

			remindersSent++
		}
	}

	log.Printf("✅ H-1 reminders sent: %d", remindersSent)
}

// UpdateEventStatuses updates event statuses based on current time
func (s *Scheduler) UpdateEventStatuses() {
	ctx := context.Background()
	log.Println("🔄 Running event status updater...")

	// Get all non-completed/cancelled events
	filters := map[string]interface{}{}
	events, err := s.eventRepo.GetAll(ctx, filters)
	if err != nil {
		log.Printf("Failed to get events: %v", err)
		return
	}

	updated := 0

	for _, event := range events {
		var newStatus string
		shouldUpdate := false

		switch event.Status {
		case domain.StatusPublished:
			// Check if event has started
			if event.HasStarted() && !event.HasEnded() {
				newStatus = domain.StatusOngoing
				shouldUpdate = true
			}

		case domain.StatusOngoing:
			// Check if event has ended
			if event.HasEnded() {
				newStatus = domain.StatusCompleted
				shouldUpdate = true
			}
		}

		if shouldUpdate {
			if err := s.eventRepo.UpdateStatus(ctx, event.ID, newStatus); err != nil {
				log.Printf("Failed to update event %s status: %v", event.Title, err)
			} else {
				log.Printf("  Updated event '%s': %s → %s", event.Title, event.Status, newStatus)
				updated++
			}
		}
	}

	log.Printf("✅ Event statuses updated: %d", updated)
}

// RunNow runs specific job immediately (for testing)
func (s *Scheduler) RunH1RemindersNow() {
	s.SendH1Reminders()
}

func (s *Scheduler) RunStatusUpdaterNow() {
	s.UpdateEventStatuses()
}

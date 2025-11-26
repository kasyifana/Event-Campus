package utils

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

var (
	ErrInvalidEmail           = errors.New("invalid email format")
	ErrInvalidPhone           = errors.New("invalid phone number format")
	ErrInvalidEventDates      = errors.New("invalid event dates")
	ErrRegistrationAfterStart = errors.New("registration deadline must be before event start")
	ErrEndBeforeStart         = errors.New("event end date must be after start date")
	ErrPastDate               = errors.New("event dates must be in the future")
	ErrInvalidZoomLink        = errors.New("zoom link is required for online events")
	ErrInvalidLocation        = errors.New("location is required for offline events")
)

// IsUIIEmail checks if email is from UII domain
func IsUIIEmail(email string) bool {
	return strings.HasSuffix(strings.ToLower(email), "uii.ac.id")
}

// IsValidEmail validates email format
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsValidPhone validates Indonesian phone number format
func IsValidPhone(phone string) bool {
	// Remove common separators
	cleaned := strings.ReplaceAll(phone, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")

	// Indonesia phone format: 08xx, +628xx, or 628xx
	phoneRegex := regexp.MustCompile(`^(\+62|62|0)[0-9]{9,12}$`)
	return phoneRegex.MatchString(cleaned)
}

// IsValidEventDate validates event dates
func IsValidEventDate(startDate, endDate, registrationDeadline time.Time) error {
	now := time.Now()

	// Check if dates are in the future
	if startDate.Before(now) {
		return ErrPastDate
	}

	// Check if end date is after start date
	if endDate.Before(startDate) || endDate.Equal(startDate) {
		return ErrEndBeforeStart
	}

	// Check if registration deadline is before start date
	if registrationDeadline.After(startDate) || registrationDeadline.Equal(startDate) {
		return ErrRegistrationAfterStart
	}

	return nil
}

// ValidateEventType validates event type specific requirements
func ValidateEventType(eventType string, location, zoomLink *string) error {
	if eventType == "online" {
		if zoomLink == nil || *zoomLink == "" {
			return ErrInvalidZoomLink
		}
	} else if eventType == "offline" {
		if location == nil || *location == "" {
			return ErrInvalidLocation
		}
	}
	return nil
}

// NormalizePhone normalizes phone number to standard format
func NormalizePhone(phone string) string {
	// Remove separators
	cleaned := strings.ReplaceAll(phone, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")

	// Convert to +62 format
	if strings.HasPrefix(cleaned, "0") {
		cleaned = "+62" + cleaned[1:]
	} else if strings.HasPrefix(cleaned, "62") {
		cleaned = "+" + cleaned
	}

	return cleaned
}

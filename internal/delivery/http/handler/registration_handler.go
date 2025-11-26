package handler

import (
	"event-campus-backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RegistrationHandler handles registration endpoints
type RegistrationHandler struct {
	registrationUsecase usecase.RegistrationUsecase
}

// NewRegistrationHandler creates a new registration handler
func NewRegistrationHandler(registrationUsecase usecase.RegistrationUsecase) *RegistrationHandler {
	return &RegistrationHandler{
		registrationUsecase: registrationUsecase,
	}
}

// RegisterForEvent handles event registration
func (h *RegistrationHandler) RegisterForEvent(c *gin.Context) {
	// Get user ID from context
	userIDStr, _ := c.Get("userID")
	userID, _ := uuid.Parse(userIDStr.(string))

	// Get event ID from URL
	eventIDStr := c.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid event ID",
		})
		return
	}

	if err := h.registrationUsecase.RegisterForEvent(c.Request.Context(), userID, eventID); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Failed to register for event",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(201, gin.H{
		"success": true,
		"message": "Registration successful",
	})
}

// CancelRegistration handles registration cancellation
func (h *RegistrationHandler) CancelRegistration(c *gin.Context) {
	// Get user ID from context
	userIDStr, _ := c.Get("userID")
	userID, _ := uuid.Parse(userIDStr.(string))

	// Get registration ID from URL
	registrationIDStr := c.Param("id")
	registrationID, err := uuid.Parse(registrationIDStr)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid registration ID",
		})
		return
	}

	if err := h.registrationUsecase.CancelRegistration(c.Request.Context(), userID, registrationID); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Failed to cancel registration",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Registration cancelled successfully",
	})
}

// GetMyRegistrations gets user's registrations
func (h *RegistrationHandler) GetMyRegistrations(c *gin.Context) {
	// Get user ID from context
	userIDStr, _ := c.Get("userID")
	userID, _ := uuid.Parse(userIDStr.(string))

	registrations, err := h.registrationUsecase.GetMyRegistrations(c.Request.Context(), userID)
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to get registrations",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Registrations retrieved successfully",
		"data":    registrations,
	})
}

// GetEventRegistrations gets event's registrations (organizer only)
func (h *RegistrationHandler) GetEventRegistrations(c *gin.Context) {
	// Get organizer ID from context
	organizerIDStr, _ := c.Get("userID")
	organizerID, _ := uuid.Parse(organizerIDStr.(string))

	// Get event ID from URL
	eventIDStr := c.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid event ID",
		})
		return
	}

	registrations, err := h.registrationUsecase.GetEventRegistrations(c.Request.Context(), organizerID, eventID)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Failed to get registrations",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Registrations retrieved successfully",
		"data":    registrations,
	})
}

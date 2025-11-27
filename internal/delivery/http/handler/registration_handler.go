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
// @Summary Register for an event
// @Description Register authenticated user for a specific event
// @Tags Registrations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID (UUID)"
// @Success 201 {object} map[string]interface{} "Registration successful"
// @Failure 400 {object} map[string]interface{} "Invalid event ID or registration failed"
// @Router /events/{id}/register [post]
func (h *RegistrationHandler) RegisterForEvent(c *gin.Context) {
	// Get user ID from context
	userIDInterface, _ := c.Get("userID")
	userID, _ := userIDInterface.(uuid.UUID)

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

	registration, err := h.registrationUsecase.RegisterForEvent(c.Request.Context(), userID, eventID)
	if err != nil {
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
		"data": gin.H{
			"registration_id": registration.ID,
			"status":          registration.Status,
		},
	})
}

// CancelRegistration handles registration cancellation
// @Summary Cancel event registration
// @Description Cancel user's registration for an event
// @Tags Registrations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Registration ID (UUID)"
// @Success 200 {object} map[string]interface{} "Registration cancelled successfully"
// @Failure 400 {object} map[string]interface{} "Invalid ID or cancellation failed"
// @Router /registrations/{id} [delete]
func (h *RegistrationHandler) CancelRegistration(c *gin.Context) {
	// Get user ID from context
	userIDInterface, _ := c.Get("userID")
	userID, _ := userIDInterface.(uuid.UUID)

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
// @Summary Get my registrations
// @Description Get all registrations for authenticated user
// @Tags Registrations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Registrations retrieved successfully"
// @Failure 500 {object} map[string]interface{} "Failed to get registrations"
// @Router /registrations/my [get]
func (h *RegistrationHandler) GetMyRegistrations(c *gin.Context) {
	// Get user ID from context
	userIDInterface, _ := c.Get("userID")
	userID, _ := userIDInterface.(uuid.UUID)

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
// @Summary Get event registrations
// @Description Get all registrations for a specific event (organizer only)
// @Tags Registrations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID (UUID)"
// @Success 200 {object} map[string]interface{} "Registrations retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid event ID or failed to get registrations"
// @Router /events/{id}/registrations [get]
func (h *RegistrationHandler) GetEventRegistrations(c *gin.Context) {
	// Get organizer ID from context
	organizerIDInterface, _ := c.Get("userID")
	organizerID, _ := organizerIDInterface.(uuid.UUID)

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

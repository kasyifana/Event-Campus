package handler

import (
	"event-campus-backend/internal/dto/request"
	"event-campus-backend/internal/usecase"
	"event-campus-backend/internal/utils"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// EventHandler handles event endpoints
type EventHandler struct {
	eventUsecase usecase.EventUsecase
	fileUploader *utils.FileUploader
}

// NewEventHandler creates a new event handler
func NewEventHandler(eventUsecase usecase.EventUsecase, fileUploader *utils.FileUploader) *EventHandler {
	return &EventHandler{
		eventUsecase: eventUsecase,
		fileUploader: fileUploader,
	}
}

// CreateEvent handles event creation
// @Summary Create a new event
// @Description Create a new event (organizer only). Poster can be uploaded separately.
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateEventRequest true "Event details"
// @Success 201 {object} map[string]interface{} "Event created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or failed to create"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /events [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {
	// Get organizer ID from context
	organizerIDInterface, _ := c.Get("userID")
	organizerID, _ := organizerIDInterface.(uuid.UUID)

	// Parse JSON request
	var req request.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	// Create event (poster can be uploaded separately)
	event, err := h.eventUsecase.CreateEvent(c.Request.Context(), organizerID, &req, nil)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Failed to create event",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(201, gin.H{
		"success": true,
		"message": "Event created successfully",
		"data": gin.H{
			"id": event.ID,
		},
	})
}

// GetAllEvents gets list of events with filters
// @Summary Get all events
// @Description Get list of events with optional filters (category, status, event_type, search)
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category query string false "Filter by category"
// @Param status query string false "Filter by status (draft/published/ongoing/completed/cancelled)"
// @Param event_type query string false "Filter by event type (online/offline/hybrid)"
// @Param search query string false "Search by name or description"
// @Success 200 {object} map[string]interface{} "Events retrieved successfully"
// @Failure 500 {object} map[string]interface{} "Failed to get events"
// @Router /events [get]
func (h *EventHandler) GetAllEvents(c *gin.Context) {
	filters := make(map[string]interface{})

	// Get query parameters
	if category := c.Query("category"); category != "" {
		filters["category"] = category
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if eventType := c.Query("event_type"); eventType != "" {
		filters["event_type"] = eventType
	}
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}

	events, err := h.eventUsecase.GetAllEvents(c.Request.Context(), filters)
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to get events",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Events retrieved successfully",
		"data":    events,
	})
}

// GetEvent gets event detail
// @Summary Get event by ID
// @Description Get detailed information about a specific event
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID (UUID)"
// @Success 200 {object} map[string]interface{} "Event retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid event ID"
// @Failure 404 {object} map[string]interface{} "Event not found"
// @Router /events/{id} [get]
func (h *EventHandler) GetEvent(c *gin.Context) {
	eventIDStr := c.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid event ID",
		})
		return
	}

	event, err := h.eventUsecase.GetEvent(c.Request.Context(), eventID)
	if err != nil {
		c.JSON(404, gin.H{
			"success": false,
			"message": "Event not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Event retrieved successfully",
		"data":    event,
	})
}

// GetMyEvents gets organizer's events
// @Summary Get my events
// @Description Get all events created by authenticated organizer
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Events retrieved successfully"
// @Failure 500 {object} map[string]interface{} "Failed to get events"
// @Router /events/my-events [get]
func (h *EventHandler) GetMyEvents(c *gin.Context) {
	organizerIDInterface, _ := c.Get("userID")
	organizerID, _ := organizerIDInterface.(uuid.UUID)

	events, err := h.eventUsecase.GetMyEvents(c.Request.Context(), organizerID)
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to get events",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Events retrieved successfully",
		"data":    events,
	})
}

// UpdateEvent handles event update
// @Summary Update event
// @Description Update event details (organizer only)
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID (UUID)"
// @Param request body request.UpdateEventRequest true "Updated event details"
// @Success 200 {object} map[string]interface{} "Event updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or update failed"
// @Router /events/{id} [put]
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	organizerIDInterface, _ := c.Get("userID")
	organizerID, _ := organizerIDInterface.(uuid.UUID)

	eventIDStr := c.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid event ID",
		})
		return
	}

	var req request.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	if err := h.eventUsecase.UpdateEvent(c.Request.Context(), organizerID, eventID, &req, nil); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Failed to update event",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Event updated successfully",
	})
}

// UploadPoster handles poster upload
// @Summary Upload event poster
// @Description Upload poster image for an event (organizer only)
// @Tags Events
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID (UUID)"
// @Param poster formData file true "Poster image (JPG/PNG)"
// @Success 200 {object} map[string]interface{} "Poster uploaded successfully"
// @Failure 400 {object} map[string]interface{} "Invalid file or upload failed"
// @Router /events/{id}/poster [post]
func (h *EventHandler) UploadPoster(c *gin.Context) {
	organizerIDInterface, _ := c.Get("userID")
	organizerID, _ := organizerIDInterface.(uuid.UUID)

	eventIDStr := c.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid event ID",
		})
		return
	}

	// Get uploaded file
	file, err := c.FormFile("poster")
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Poster file is required",
			"error":   err.Error(),
		})
		return
	}

	// Validate file extension
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid file type. Only JPG and PNG are allowed",
		})
		return
	}

	// Save poster
	posterPath, err := h.fileUploader.SavePoster(file)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Failed to upload poster",
			"error":   err.Error(),
		})
		return
	}

	// Update event with poster path
	req := &request.UpdateEventRequest{}
	if err := h.eventUsecase.UpdateEvent(c.Request.Context(), organizerID, eventID, req, &posterPath); err != nil {
		// Delete uploaded file if update fails
		h.fileUploader.DeleteFile(posterPath)

		c.JSON(400, gin.H{
			"success": false,
			"message": "Failed to update event with poster",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Poster uploaded successfully",
		"data": gin.H{
			"poster_path": posterPath,
		},
	})
}

// DeleteEvent handles event deletion
// @Summary Delete event
// @Description Delete an event (organizer only)
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID (UUID)"
// @Success 200 {object} map[string]interface{} "Event deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid event ID or deletion failed"
// @Router /events/{id} [delete]
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	organizerIDInterface, _ := c.Get("userID")
	organizerID, _ := organizerIDInterface.(uuid.UUID)

	eventIDStr := c.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid event ID",
		})
		return
	}

	if err := h.eventUsecase.DeleteEvent(c.Request.Context(), organizerID, eventID); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Failed to delete event",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Event deleted successfully",
	})
}

// PublishEvent handles event publishing
// @Summary Publish event
// @Description Publish a draft event to make it visible to users
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID (UUID)"
// @Success 200 {object} map[string]interface{} "Event published successfully"
// @Failure 400 {object} map[string]interface{} "Invalid event ID or publish failed"
// @Router /events/{id}/publish [post]
func (h *EventHandler) PublishEvent(c *gin.Context) {
	organizerIDInterface, _ := c.Get("userID")
	organizerID, _ := organizerIDInterface.(uuid.UUID)

	eventIDStr := c.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid event ID",
		})
		return
	}

	if err := h.eventUsecase.PublishEvent(c.Request.Context(), organizerID, eventID); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Failed to publish event",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Event published successfully",
	})
}

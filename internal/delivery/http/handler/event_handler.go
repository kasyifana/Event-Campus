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

package handler

import (
	"event-campus-backend/internal/dto/request"
	"event-campus-backend/internal/usecase"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AttendanceHandler handles attendance endpoints
type AttendanceHandler struct {
	attendanceUsecase usecase.AttendanceUsecase
}

// NewAttendanceHandler creates a new attendance handler
func NewAttendanceHandler(attendanceUsecase usecase.AttendanceUsecase) *AttendanceHandler {
	return &AttendanceHandler{
		attendanceUsecase: attendanceUsecase,
	}
}

// MarkAttendance handles single attendance marking
func (h *AttendanceHandler) MarkAttendance(c *gin.Context) {
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

	// Parse request
	var req request.MarkAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	if err := h.attendanceUsecase.MarkAttendance(c.Request.Context(), organizerID, eventID, userID, req.Notes); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Failed to mark attendance",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Attendance marked successfully",
	})
}

// BulkMarkAttendance handles bulk attendance marking
func (h *AttendanceHandler) BulkMarkAttendance(c *gin.Context) {
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

	// Parse request
	var req request.BulkMarkAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	// Parse user IDs
	var userIDs []uuid.UUID
	for _, idStr := range req.UserIDs {
		userID, err := uuid.Parse(idStr)
		if err != nil {
			continue // Skip invalid IDs
		}
		userIDs = append(userIDs, userID)
	}

	if len(userIDs) == 0 {
		c.JSON(400, gin.H{
			"success": false,
			"message": "No valid user IDs provided",
		})
		return
	}

	if err := h.attendanceUsecase.BulkMarkAttendance(c.Request.Context(), organizerID, eventID, userIDs); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Failed to bulk mark attendance",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": fmt.Sprintf("Bulk attendance marked successfully for %d users", len(userIDs)),
	})
}

// GetEventAttendance gets attendance list for event
func (h *AttendanceHandler) GetEventAttendance(c *gin.Context) {
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

	attendances, err := h.attendanceUsecase.GetEventAttendance(c.Request.Context(), organizerID, eventID)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Failed to get attendance",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Attendance retrieved successfully",
		"data":    attendances,
	})
}

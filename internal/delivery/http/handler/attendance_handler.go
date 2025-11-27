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
// @Summary Mark single attendance
// @Description Mark attendance for a single user at an event (organizer only)
// @Tags Attendance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID (UUID)"
// @Param request body request.MarkAttendanceRequest true "Attendance details"
// @Success 200 {object} map[string]interface{} "Attendance marked successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or failed to mark"
// @Router /events/{id}/attendance [post]
func (h *AttendanceHandler) MarkAttendance(c *gin.Context) {
	// Get organizer ID from context
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
// @Summary Mark bulk attendance
// @Description Mark attendance for multiple users at once (organizer only)
// @Tags Attendance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID (UUID)"
// @Param request body request.BulkMarkAttendanceRequest true "List of user IDs"
// @Success 200 {object} map[string]interface{} "Bulk attendance marked successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or failed to mark"
// @Router /events/{id}/attendance/bulk [post]
func (h *AttendanceHandler) BulkMarkAttendance(c *gin.Context) {
	// Get organizer ID from context
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
// @Summary Get event attendance
// @Description Get attendance list for a specific event (organizer only)
// @Tags Attendance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID (UUID)"
// @Success 200 {object} map[string]interface{} "Attendance retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid event ID or failed to get attendance"
// @Router /events/{id}/attendance [get]
func (h *AttendanceHandler) GetEventAttendance(c *gin.Context) {
	// Get organizer ID from context
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

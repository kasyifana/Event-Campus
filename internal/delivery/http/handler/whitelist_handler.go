package handler

import (
	"event-campus-backend/internal/dto/request"
	"event-campus-backend/internal/usecase"
	"event-campus-backend/internal/utils"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// WhitelistHandler handles whitelist endpoints
type WhitelistHandler struct {
	whitelistUsecase usecase.WhitelistUsecase
	fileUploader     *utils.FileUploader
}

// NewWhitelistHandler creates a new whitelist handler
func NewWhitelistHandler(whitelistUsecase usecase.WhitelistUsecase, fileUploader *utils.FileUploader) *WhitelistHandler {
	return &WhitelistHandler{
		whitelistUsecase: whitelistUsecase,
		fileUploader:     fileUploader,
	}
}

// SubmitRequest handles whitelist request submission
func (h *WhitelistHandler) SubmitRequest(c *gin.Context) {
	// Get user ID from context
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{
			"success": false,
			"message": "Unauthorized",
			"error":   "User ID not found in context",
		})
		return
	}

	// Auth middleware stores userID as uuid.UUID, not string
	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid user ID",
			"error":   "User ID type mismatch",
		})
		return
	}

	// Get organization name from form
	orgName := c.PostForm("organization_name")
	if orgName == "" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   "organization_name is required",
		})
		return
	}

	// Get uploaded file
	file, err := c.FormFile("document")
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   "document file is required",
		})
		return
	}

	// Validate file extension
	ext := filepath.Ext(file.Filename)
	if ext != ".pdf" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid file type",
			"error":   "Only PDF files are allowed",
		})
		return
	}

	// Save document
	documentPath, err := h.fileUploader.SaveDocument(file)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Failed to upload document",
			"error":   err.Error(),
		})
		return
	}

	// Create request
	req := &request.SubmitWhitelistRequest{
		OrganizationName: orgName,
	}

	// Submit whitelist request
	whitelistReq, err := h.whitelistUsecase.SubmitRequest(c.Request.Context(), userID, req, documentPath)
	if err != nil {
		// Delete uploaded file if request fails
		h.fileUploader.DeleteFile(documentPath)

		c.JSON(400, gin.H{
			"success": false,
			"message": "Failed to submit request",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(201, gin.H{
		"success": true,
		"message": "Whitelist request submitted successfully",
		"data": gin.H{
			"id": whitelistReq.ID,
		},
	})
}

// GetMyRequest gets current user's whitelist request
func (h *WhitelistHandler) GetMyRequest(c *gin.Context) {
	// Get user ID from context
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	request, err := h.whitelistUsecase.GetMyRequest(c.Request.Context(), userID)
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to get request",
			"error":   err.Error(),
		})
		return
	}

	if request == nil {
		c.JSON(404, gin.H{
			"success": false,
			"message": "No whitelist request found",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Request retrieved successfully",
		"data":    request,
	})
}

// GetAllRequests gets all whitelist requests (admin only)
func (h *WhitelistHandler) GetAllRequests(c *gin.Context) {
	status := c.Query("status")

	requests, err := h.whitelistUsecase.GetAllRequests(c.Request.Context(), status)
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to get requests",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Requests retrieved successfully",
		"data":    requests,
	})
}

// ReviewRequest handles whitelist request review (admin only)
func (h *WhitelistHandler) ReviewRequest(c *gin.Context) {
	// Get request ID from URL
	requestIDStr := c.Param("id")
	requestID, err := uuid.Parse(requestIDStr)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid request ID",
		})
		return
	}

	// Get reviewer ID from context
	reviewerIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	reviewerID, ok := reviewerIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid reviewer ID",
		})
		return
	}

	// Parse request body
	var req request.ReviewWhitelistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	// Review request
	if err := h.whitelistUsecase.ReviewRequest(c.Request.Context(), requestID, reviewerID, &req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Failed to review request",
			"error":   err.Error(),
		})
		return
	}

	var message string
	if req.Approved {
		message = "Request approved successfully"
	} else {
		message = "Request rejected"
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": message,
	})
}

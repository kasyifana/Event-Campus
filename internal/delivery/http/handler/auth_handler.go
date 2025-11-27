package handler

import (
	"event-campus-backend/internal/dto/request"
	"event-campus-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user account with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body request.RegisterRequest true "Registration details"
// @Success 200 {object} map[string]interface{} "Registration successful"
// @Failure 400 {object} map[string]interface{} "Invalid request or registration failed"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	resp, err := h.authUsecase.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Registration failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Registration successful",
		"data":    resp,
	})
}

// Login handles user login
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "Login successful with JWT token"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Authentication failed"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	resp, err := h.authUsecase.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(401, gin.H{
			"success": false,
			"message": "Login failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Login successful",
		"data":    resp,
	})
}

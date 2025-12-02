package router

import (
	"event-campus-backend/internal/delivery/http/handler"
	"event-campus-backend/internal/delivery/http/middleware"

	_ "event-campus-backend/docs" // Swagger docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Router holds all HTTP handlers
type Router struct {
	authHandler         *handler.AuthHandler
	whitelistHandler    *handler.WhitelistHandler
	eventHandler        *handler.EventHandler
	registrationHandler *handler.RegistrationHandler
	attendanceHandler   *handler.AttendanceHandler
	jwtSecret           string
	corsOrigins         []string
}

// NewRouter creates a new router
func NewRouter(
	authHandler *handler.AuthHandler,
	whitelistHandler *handler.WhitelistHandler,
	eventHandler *handler.EventHandler,
	registrationHandler *handler.RegistrationHandler,
	attendanceHandler *handler.AttendanceHandler,
	jwtSecret string,
	corsOrigins []string,
) *Router {
	return &Router{
		authHandler:         authHandler,
		whitelistHandler:    whitelistHandler,
		eventHandler:        eventHandler,
		registrationHandler: registrationHandler,
		attendanceHandler:   attendanceHandler,
		jwtSecret:           jwtSecret,
		corsOrigins:         corsOrigins,
	}
}

// Setup sets up all routes
func (r *Router) Setup() *gin.Engine {
	router := gin.Default()

	// CORS middleware
	router.Use(middleware.CORSMiddleware(r.corsOrigins))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Event Campus API is running",
		})
	})

	// Swagger documentation
	// Access via: /docs/index.html
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Authentication routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", r.authHandler.Register)
			auth.POST("/login", r.authHandler.Login)
		}

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(r.jwtSecret))
		{
			// User routes
			// @Summary Get user profile
			// @Description Get authenticated user's profile information
			// @Tags User
			// @Accept json
			// @Produce json
			// @Security BearerAuth
			// @Success 200 {object} map[string]interface{} "Profile retrieved successfully"
			// @Router /profile [get]
			protected.GET("/profile", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"success": true,
					"message": "Profile retrieved",
					"data": gin.H{
						"user_id": c.GetString("userID"),
						"email":   c.GetString("userEmail"),
						"role":    c.GetString("userRole"),
					},
				})
			})

			// Whitelist routes
			whitelist := protected.Group("/whitelist")
			{
				// Mahasiswa can submit request and check their own
				whitelist.POST("/request", r.whitelistHandler.SubmitRequest)
				whitelist.GET("/my-request", r.whitelistHandler.GetMyRequest)

				// Admin only routes
				whitelist.GET("/requests", middleware.RequireAdmin(), r.whitelistHandler.GetAllRequests)
				whitelist.PATCH("/:id/review", middleware.RequireAdmin(), r.whitelistHandler.ReviewRequest)
			}

			// Event routes
			events := protected.Group("/events")
			{
				// Public routes (anyone authenticated can view)
				events.GET("", r.eventHandler.GetAllEvents)
				events.GET("/:id", r.eventHandler.GetEvent)

				// Organisasi & Admin routes
				events.POST("", middleware.RequireOrganisasi(), r.eventHandler.CreateEvent)
				events.GET("/my-events", middleware.RequireOrganisasi(), r.eventHandler.GetMyEvents)
				events.PUT("/:id", middleware.RequireOrganisasi(), r.eventHandler.UpdateEvent)
				events.POST("/:id/poster", middleware.RequireOrganisasi(), r.eventHandler.UploadPoster)
				events.DELETE("/:id", middleware.RequireOrganisasi(), r.eventHandler.DeleteEvent)
				events.POST("/:id/publish", middleware.RequireOrganisasi(), r.eventHandler.PublishEvent)
				events.POST("/:id/reminders", middleware.RequireOrganisasi(), r.eventHandler.SendReminders)

				// Registration routes
				events.POST("/:id/register", r.registrationHandler.RegisterForEvent)
				events.GET("/:id/registrations", middleware.RequireOrganisasi(), r.registrationHandler.GetEventRegistrations)

				// Attendance routes
				events.POST("/:id/attendance", middleware.RequireOrganisasi(), r.attendanceHandler.MarkAttendance)
				events.POST("/:id/attendance/bulk", middleware.RequireOrganisasi(), r.attendanceHandler.BulkMarkAttendance)
				events.GET("/:id/attendance", middleware.RequireOrganisasi(), r.attendanceHandler.GetEventAttendance)
			}

			// Registration routes
			registrations := protected.Group("/registrations")
			{
				registrations.GET("/my", r.registrationHandler.GetMyRegistrations)
				registrations.DELETE("/:id", r.registrationHandler.CancelRegistration)
			}
		}
	}

	// Serve uploaded files
	router.Static("/files", "./storage")

	return router
}

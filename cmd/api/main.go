package main

import (
	"database/sql"
	"event-campus-backend/internal/config"
	"event-campus-backend/internal/delivery/http/handler"
	"event-campus-backend/internal/delivery/http/router"
	"event-campus-backend/internal/repository"
	"event-campus-backend/internal/scheduler"
	"event-campus-backend/internal/usecase"
	"event-campus-backend/internal/utils"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println("🚀 Starting Event Campus API...")
	log.Printf("📍 Environment: %s", cfg.Server.Env)

	// Initialize PostgreSQL connection
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.PostgreSQL.Host,
		cfg.PostgreSQL.Port,
		cfg.PostgreSQL.User,
		cfg.PostgreSQL.Password,
		cfg.PostgreSQL.Database,
		cfg.PostgreSQL.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("✅ PostgreSQL connected to Supabase")

	// Run database migrations
	log.Println("🔄 Running database migrations...")
	if err := repository.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories with database connection
	userRepo := repository.NewUserRepository(db)
	whitelistRepo := repository.NewWhitelistRepository(db)
	eventRepo := repository.NewEventRepository(db)
	registrationRepo := repository.NewRegistrationRepository(db)
	attendanceRepo := repository.NewAttendanceRepository(db)

	// Parse JWT expiration
	jwtExpiration, err := time.ParseDuration(cfg.JWT.Expiration)
	if err != nil {
		log.Fatalf("Invalid JWT expiration: %v", err)
	}

	// Initialize email sender
	emailSender := utils.NewEmailSender(
		cfg.Email.SMTPHost,
		cfg.Email.SMTPPort,
		cfg.Email.SMTPUser,
		cfg.Email.SMTPPassword,
	)

	// Initialize file uploader
	fileUploader := utils.NewFileUploader(cfg.Upload.Path, cfg.Upload.MaxSize)

	// Initialize use cases
	authUsecase := usecase.NewAuthUsecase(userRepo, cfg.JWT.Secret, jwtExpiration)
	whitelistUsecase := usecase.NewWhitelistUsecase(
		whitelistRepo,
		userRepo,
		emailSender,
		fmt.Sprintf("http://localhost:%s", cfg.Server.Port),
	)
	eventUsecase := usecase.NewEventUsecase(
		eventRepo,
		userRepo,
		fmt.Sprintf("http://localhost:%s", cfg.Server.Port),
	)
	registrationUsecase := usecase.NewRegistrationUsecase(
		registrationRepo,
		eventRepo,
		userRepo,
		emailSender,
	)
	attendanceUsecase := usecase.NewAttendanceUsecase(
		attendanceRepo,
		eventRepo,
		registrationRepo,
		userRepo,
	)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUsecase)
	whitelistHandler := handler.NewWhitelistHandler(whitelistUsecase, fileUploader)
	eventHandler := handler.NewEventHandler(eventUsecase, fileUploader)
	registrationHandler := handler.NewRegistrationHandler(registrationUsecase)
	attendanceHandler := handler.NewAttendanceHandler(attendanceUsecase)

	// Setup router
	r := router.NewRouter(
		authHandler,
		whitelistHandler,
		eventHandler,
		registrationHandler,
		attendanceHandler,
		cfg.JWT.Secret,
		cfg.CORS.AllowedOrigins,
	)

	// Initialize and start scheduler
	sched := scheduler.NewScheduler(
		eventRepo,
		registrationRepo,
		userRepo,
		emailSender,
	)

	if err := sched.Start(); err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}
	defer sched.Stop()

	// Setup Gin engine
	ginRouter := r.Setup()

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("✅ Server running on http://localhost%s", addr)
	log.Println("📚 API Documentation: http://localhost" + addr + "/health")
	log.Println("")
	log.Println("Available endpoints:")
	log.Println("  POST /api/v1/auth/register - User registration")
	log.Println("  POST /api/v1/auth/login - User login")
	log.Println("  GET  /api/v1/profile - Get user profile (protected)")
	log.Println("  GET  /health - Health check")
	log.Println("")

	if err := ginRouter.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

package usecase

import (
	"context"
	"event-campus-backend/internal/domain"
	"event-campus-backend/internal/dto/request"
	"event-campus-backend/internal/dto/response"
	"event-campus-backend/internal/repository"
	"event-campus-backend/internal/utils"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// AuthUsecase defines interface for authentication business logic
type AuthUsecase interface {
	Register(ctx context.Context, req *request.RegisterRequest) (*response.LoginResponse, error)
	Login(ctx context.Context, req *request.LoginRequest) (*response.LoginResponse, error)
}

type authUsecase struct {
	userRepo      repository.UserRepository
	jwtSecret     string
	jwtExpiration time.Duration
}

// NewAuthUsecase creates a new authentication usecase
func NewAuthUsecase(userRepo repository.UserRepository, jwtSecret string, jwtExpiration time.Duration) AuthUsecase {
	return &authUsecase{
		userRepo:      userRepo,
		jwtSecret:     jwtSecret,
		jwtExpiration: jwtExpiration,
	}
}

func (u *authUsecase) Register(ctx context.Context, req *request.RegisterRequest) (*response.LoginResponse, error) {
	// Validate email format
	if !utils.IsValidEmail(req.Email) {
		return nil, fmt.Errorf("invalid email format")
	}

	// Validate phone number
	if !utils.IsValidPhone(req.PhoneNumber) {
		return nil, fmt.Errorf("invalid phone number format")
	}

	// Check if user already exists
	existingUser, _ := u.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Detect UII civitas
	isUIICivitas := domain.IsUIIEmail(req.Email)

	// Create user
	user := &domain.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		PhoneNumber:  utils.NormalizePhone(req.PhoneNumber),
		Role:         domain.RoleMahasiswa,
		IsUIICivitas: isUIICivitas,
		IsApproved:   false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save user
	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, u.jwtSecret, u.jwtExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Return response
	return &response.LoginResponse{
		Token: token,
		User:  response.ToUserResponse(user),
	}, nil
}

func (u *authUsecase) Login(ctx context.Context, req *request.LoginRequest) (*response.LoginResponse, error) {
	// Get user by email
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, u.jwtSecret, u.jwtExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Return response
	return &response.LoginResponse{
		Token: token,
		User:  response.ToUserResponse(user),
	}, nil
}

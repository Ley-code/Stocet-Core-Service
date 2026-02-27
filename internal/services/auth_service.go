package services

import (
	"errors"
	"stocet-core-service/internal/config"
	"stocet-core-service/internal/models"
	"stocet-core-service/internal/repositories"
	"stocet-core-service/internal/utils"

	"github.com/google/uuid"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignupRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func (s *AuthService) Login(req LoginRequest) (*AuthResponse, error) {
	// Validate email
	if !utils.ValidateEmail(req.Email) {
		return nil, errors.New("invalid email format")
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) Signup(req SignupRequest) (*AuthResponse, error) {
	// Validate email
	if !utils.ValidateEmail(req.Email) {
		return nil, errors.New("invalid email format")
	}

	// Validate password
	if !utils.ValidatePassword(req.Password) {
		return nil, errors.New("password must be at least 8 characters with at least one letter and one number")
	}

	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := &models.User{
		ID:                 uuid.New().String(),
		Email:              req.Email,
		PasswordHash:       hashedPassword,
		Role:               config.USER_ROLE,
		SubscriptionStatus: config.SUBSCRIPTION_FREE,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

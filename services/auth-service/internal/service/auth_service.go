package service

import (
	"errors"
	"time"

	"github.com/vidhyashekar/cloudcart/services/auth-service/internal/auth"
	"github.com/vidhyashekar/cloudcart/services/auth-service/internal/model"
	"github.com/vidhyashekar/cloudcart/services/auth-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(request RegisterRequest) (*RegisterResponse, error)
	FindUserByEmail(email string) (*model.User, error)
	Login(request LoginRequest) (*LoginResponse, error)
	RefreshAccessToken(refreshToken string) (*RefreshTokenResponse, error)
}

type authService struct {
	userRepository repository.UserRepository
	jwtManager     *auth.JWTManager
}

func NewAuthService(userRepository repository.UserRepository, jwtManager *auth.JWTManager) AuthService {
	return &authService{
		userRepository: userRepository,
		jwtManager:     jwtManager,
	}
}

func (s *authService) Register(request RegisterRequest) (*RegisterResponse, error) {
	// Check whether email already exists
	existingUser, err := s.userRepository.FindByEmail(request.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Get default CUSTOMER role
	role, err := s.userRepository.FindRoleByName("CUSTOMER")
	if err != nil {
		return nil, err
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	// Request mapping to model.User
	user := &model.User{
		FirstName:    request.FirstName,
		LastName:     request.LastName,
		Email:        request.Email,
		PasswordHash: string(passwordHash),
		RoleID:       role.ID,
	}

	// Create user in database
	if err := s.userRepository.Create(user); err != nil {
		return nil, err
	}

	return &RegisterResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}, nil
}

// Login authenticates a user and returns a JWT token if successful
func (s *authService) Login(request LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepository.FindByEmail(request.Email)
	if err != nil {
		return nil, errors.New("Invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(request.Password),
	); err != nil {
		return nil, errors.New("Invalid email or password")
	}

	token, err := s.jwtManager.GenerateToken(
		user.ID,
		user.Email,
		user.Role.Name,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}
	refreshTokenModel := &model.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := s.userRepository.CreateRefreshToken(refreshTokenModel); err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	}, nil
}

// FindUserByEmail retrieves a user by their email address.
func (s *authService) FindUserByEmail(email string) (*model.User, error) {
	return s.userRepository.FindByEmail(email)
}

// RefreshAccessToken generates a new access token using the provided refresh token.
func (s *authService) RefreshAccessToken(refreshToken string) (*RefreshTokenResponse, error) {
	token, err := s.userRepository.FindRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if time.Now().After(token.ExpiresAt) {
		return nil, errors.New("refresh token expired")
	}

	accessToken, err := s.jwtManager.GenerateToken(
		token.User.ID,
		token.User.Email,
		token.User.Role.Name,
	)
	if err != nil {
		return nil, err
	}

	return &RefreshTokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
	}, nil
}

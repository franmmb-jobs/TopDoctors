package application

import (
	"errors"
	"log/slog"
	"topdoctors/internal/domain"
	"topdoctors/internal/infrastructure/config"
)

type AuthService struct {
	userRepo domain.UserRepository
	support  domain.Support
	cfg      *config.Config
}

func NewAuthService(userRepo domain.UserRepository, support domain.Support, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		cfg:      cfg,
		support:  support,
	}
}

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		slog.Warn("Login failed: user not found", "username", username)
		return "", domain.ErrInvalidCredentials
	}

	err = s.support.CompareHashPassword(password, user.Password)
	if err != nil {
		slog.Warn("Login failed: incorrect password", "username", username)
		return "", domain.ErrInvalidCredentials
	}

	tokenString, err := s.support.GenerateToken(user, s.cfg.Api.JWTSecret)
	if err != nil {
		slog.Error("Token generation failed", "username", username, "error", err)
		return "", err
	}

	slog.Info("Login successful", "username", username)
	return tokenString, nil
}

func (s *AuthService) Register(username, password string) error {
	// Check if user already exists
	existingUser, _ := s.userRepo.GetByUsername(username)
	if existingUser != nil {
		slog.Warn("Registration failed: username already taken", "username", username)
		return errors.New("username already taken")
	}

	hashedPassword, err := s.support.GenerateHashPassword(password)
	if err != nil {
		slog.Error("Password hashing failed", "username", username, "error", err)
		return err
	}

	id, err := s.support.CreateNewID()
	if err != nil {
		slog.Error("ID creation failed during registration", "error", err)
		return err
	}

	user := &domain.User{
		ID:       id,
		Username: username,
		Password: string(hashedPassword),
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		slog.Error("User creation in repository failed", "username", username, "error", err)
		return err
	}

	slog.Info("User registered successfully", "username", username)
	return nil
}

func (s *AuthService) ValidateToken(tokenString string) error {
	err := s.support.ValidateToken(tokenString, s.cfg.Api.JWTSecret)
	if err != nil {
		return err
	}

	return nil
}

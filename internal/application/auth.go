package application

import (
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
		return "", domain.ErrInvalidCredentials
	}

	err = s.support.CompareHashPassword(password, user.Password)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	tokenString, err := s.support.GenerateToken(user, s.cfg.Api.JWTSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) Register(username, password string) error {
	hashedPassword, err := s.support.GenerateHashPassword(password)
	if err != nil {
		return err
	}

	id, err := s.support.CreateNewID()
	if err != nil {
		return err
	}

	user := &domain.User{
		ID:       id,
		Username: username,
		Password: string(hashedPassword),
	}

	return s.userRepo.CreateUser(user)
}

func (s *AuthService) ValidateToken(tokenString string) error {
	err := s.support.ValidateToken(tokenString, s.cfg.Api.JWTSecret)
	if err != nil {
		return err
	}

	return nil
}

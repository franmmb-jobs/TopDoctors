package application

import (
	"errors"
	"testing"
	"topdoctors/internal/domain"
	"topdoctors/internal/infrastructure/config"
	"topdoctors/internal/mocks"

	"go.uber.org/mock/gomock"
)

func TestAuthService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockSupport := mocks.NewMockSupport(ctrl)
	cfg := &config.Config{Api: config.ApiConfig{JWTSecret: "test-secret"}}
	service := NewAuthService(mockRepo, mockSupport, cfg)

	t.Run("successful login", func(t *testing.T) {
		user := &domain.User{Username: "doctor", Password: "hashed-password"}
		mockRepo.EXPECT().GetByUsername("doctor").Return(user, nil)
		mockSupport.EXPECT().CompareHashPassword("password123", "hashed-password").Return(nil)
		mockSupport.EXPECT().GenerateToken(user, "test-secret").Return("valid-token", nil)

		token, err := service.Login("doctor", "password123")
		if err != nil {
			t.Errorf("Login() unexpected error = %v", err)
		}
		if token != "valid-token" {
			t.Errorf("Login() expected token 'valid-token', got %s", token)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.EXPECT().GetByUsername("unknown").Return(nil, errors.New("not found"))

		_, err := service.Login("unknown", "password123")
		if !errors.Is(err, domain.ErrInvalidCredentials) {
			t.Errorf("Login() expected ErrInvalidCredentials, got %v", err)
		}
	})

	t.Run("incorrect password", func(t *testing.T) {
		user := &domain.User{Username: "doctor", Password: "hashed-password"}
		mockRepo.EXPECT().GetByUsername("doctor").Return(user, nil)
		mockSupport.EXPECT().CompareHashPassword("wrong", "hashed-password").Return(errors.New("wrong"))

		_, err := service.Login("doctor", "wrong")
		if !errors.Is(err, domain.ErrInvalidCredentials) {
			t.Errorf("Login() expected ErrInvalidCredentials, got %v", err)
		}
	})
}

func TestAuthService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockSupport := mocks.NewMockSupport(ctrl)
	cfg := &config.Config{}
	service := NewAuthService(mockRepo, mockSupport, cfg)

	t.Run("successful registration", func(t *testing.T) {
		mockRepo.EXPECT().GetByUsername("newuser").Return(nil, errors.New("not found"))
		mockSupport.EXPECT().GenerateHashPassword("password123").Return("hashed-password", nil)
		mockSupport.EXPECT().CreateNewID().Return("user-id", nil)
		mockRepo.EXPECT().CreateUser(gomock.Any()).Return(nil)

		err := service.Register("newuser", "password123")
		if err != nil {
			t.Errorf("Register() unexpected error = %v", err)
		}
	})

	t.Run("username already taken", func(t *testing.T) {
		mockRepo.EXPECT().GetByUsername("existinguser").Return(&domain.User{Username: "existinguser"}, nil)

		err := service.Register("existinguser", "password123")
		if err == nil {
			t.Error("Register() expected error for existing user, got nil")
		}
		if err.Error() != "username already taken" {
			t.Errorf("Register() expected error 'username already taken', got %v", err)
		}
	})

	t.Run("hashing failure", func(t *testing.T) {
		mockRepo.EXPECT().GetByUsername("newuser").Return(nil, errors.New("not found"))
		mockSupport.EXPECT().GenerateHashPassword("password123").Return("", errors.New("hash error"))

		err := service.Register("newuser", "password123")
		if err == nil {
			t.Error("Register() expected error, got nil")
		}
	})
}

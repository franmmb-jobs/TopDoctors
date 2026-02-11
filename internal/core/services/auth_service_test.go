package services_test

import (
	"errors"
	"testing"
	"topdoctors/internal/config"
	"topdoctors/internal/core/domain"
	"topdoctors/internal/core/ports/mocks"
	"topdoctors/internal/core/services"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	cfg := &config.Config{JWTSecret: "testsecret"}
	authService := services.NewAuthService(mockUserRepo, cfg)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	user := &domain.User{
		ID:       1,
		Username: "testuser",
		Password: string(hashedPassword),
	}

	t.Run("Success", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByUsername("testuser").Return(user, nil)

		token, err := authService.Login("testuser", "password")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if token == "" {
			t.Error("expected token, got empty string")
		}

		// Verify token
		parsedToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte("testsecret"), nil
		})
		if !parsedToken.Valid {
			t.Error("generated token is invalid")
		}
	})

	t.Run("Invalid Password", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByUsername("testuser").Return(user, nil)

		_, err := authService.Login("testuser", "wrongpassword")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if err.Error() != "invalid credentials" {
			t.Errorf("expected 'invalid credentials', got '%v'", err)
		}
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByUsername("unknown").Return(nil, errors.New("not found"))

		_, err := authService.Login("unknown", "password")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

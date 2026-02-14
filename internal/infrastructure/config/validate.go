package config

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidateConfig validates the configuration struct
func ValidateConfig(cfg *Config) error {
	slog.Info("Validating configuration")

	// Create a new validator instance
	validate := validator.New()

	// Register custom validation functions if needed
	validate.RegisterValidation("isdb", tieneExtensionDB)

	// Validate the configuration struct
	if err := validate.Struct(cfg); err != nil {
		// Handle validation errors
		validationErrors := err.(validator.ValidationErrors)
		var errorMessages []string

		for _, e := range validationErrors {
			// Create user-friendly error messages
			msg := fmt.Sprintf("field '%s' failed validation on '%s' tag", e.Field(), e.Tag())
			errorMessages = append(errorMessages, msg)
		}

		return fmt.Errorf("configuration validation failed: %s", strings.Join(errorMessages, ", "))
	}

	slog.Info("Configuration validated successfully")
	return nil
}

func tieneExtensionDB(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	return strings.HasSuffix(field, ".db")
}

package config

import (
	"os"
	"testing"
)

func TestLoadConfig_FromFile(t *testing.T) {
	// Set test environment
	os.Setenv("APP_ENV", "test")
	defer os.Unsetenv("APP_ENV")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify values from file
	if cfg.Api.Port != "8010" {
		t.Errorf("Expected Api.Port to be 8010, got %s", cfg.Api.Port)
	}
}

func TestLoadConfig_FromEnv(t *testing.T) {

	// Set environment variable to override config file
	expectedPort := "9000"
	os.Setenv("API_PORT", expectedPort)
	defer os.Unsetenv("API_PORT")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Api.Port != expectedPort {
		t.Errorf("Expected Api.Port to be overridden to %s, got %s", expectedPort, cfg.Api.Port)
	}
}

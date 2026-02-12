package config

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"topdoctors/internal/helpers"

	"github.com/spf13/viper"
)

type Config struct {
	Logs     LogsConfig     `mapstructure:"logs"`
	Database DatabaseConfig `mapstructure:"database"`
	Api      ApiConfig      `mapstructure:"api"`
}

type LogsConfig struct {
	Level string `mapstructure:"level"`
}

type DatabaseConfig struct {
	Type     string `mapstructure:"type"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
}

type ApiConfig struct {
	Port      string `mapstructure:"port"`
	JWTSecret string `mapstructure:"jwt_secret"`
}

const defaultTestConfigPath = "configs/config.test.yml"

func LoadConfig() (*Config, error) {
	v := viper.New()

	// -----------------
	// Get config path
	var cfgPath string

	if Flags.InTestEnv {
		// If test env and no explicit path, use test config
		basePath := helpers.GetProjectRoot()
		cfgPath = filepath.Join(basePath, defaultTestConfigPath)
		slog.Info("Using test config", "path", cfgPath)
	} else {
		// If not test env and no env var, set config path from flag
		cfgPath = Flags.Config
	}

	// Normalize path independently of the OS separators
	cfgPath = filepath.Clean(filepath.FromSlash(cfgPath))

	// -----------------
	// Viper Setup

	// Environment variables
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv() // This will look for environment variables like LOGS_LEVEL, DATABASE_USER, etc.

	// Load from file if exists
	var fileConfigExist bool
	if _, err := os.Stat(cfgPath); err == nil {
		slog.Info("Loading config from file", "path", cfgPath)
		fileConfigExist = true
	} else {
		slog.Info("Config file not found, loading from env and defaults only", "path", cfgPath)
		fileConfigExist = false
	}

	if !fileConfigExist {
		slog.Info("Config file not found", "path", cfgPath)
	}

	if fileConfigExist {
		v.SetConfigFile(cfgPath)
		if err := v.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

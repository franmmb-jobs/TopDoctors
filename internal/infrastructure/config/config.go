package config

import (
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"

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
	DSN      string `mapstructure:"dsn"`
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
		basePath := getProjectRoot()
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

// GetProjectRoot devuelve la ruta absoluta de la raíz del proyecto
func getProjectRoot() string {
	// 1. Intentamos obtener la ruta del archivo actual que se está ejecutando
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)

	// 2. Subimos por el árbol de directorios buscando go.mod
	for {
		if _, err := os.Stat(filepath.Join(basePath, "go.mod")); err == nil {
			return basePath
		}

		parent := filepath.Dir(basePath)
		if parent == basePath {
			// Hemos llegado a la raíz del sistema de archivos sin encontrar go.mod
			return ""
		}
		basePath = parent
	}
}

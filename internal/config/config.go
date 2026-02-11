package config

import (
	"flag"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	// Standard configuration values come from YAML
	Port string `yaml:"port" env:"PORT" env-default:"8080"`

	// Secrets are defined in the struct but populated from .env or ENV vars
	Database struct {
		User     string `yaml:"user" env:"DB_USER" env-required:"true"`
		Password string `yaml:"password" env:"DB_PASSWORD" env-required:"true"`
		Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	} `yaml:"database"`

	JWTSecret string `yaml:"jwt_secret" env:"JWT_SECRET" env-required:"true"`
}

func LoadConfig() (*Config, error) {

	// Definimos el flag: nombre, valor por defecto y descripción
	configFile := flag.String("config", "config.yml", "Ruta al archivo de configuración YAML")

	flag.Parse()

	var cfg Config

	// 1. Read YAML (base configuration)
	// 2. Overwrite with Environment Variables (secrets)
	err := cleanenv.ReadConfig(*configFile, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

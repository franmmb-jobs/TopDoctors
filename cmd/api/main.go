package main

import (
	"log/slog"
	"os"
	"topdoctors/internal/application"
	"topdoctors/internal/infrastructure/config"
	httpinfra "topdoctors/internal/infrastructure/http"
	"topdoctors/internal/infrastructure/persistence"
	"topdoctors/internal/infrastructure/shared"
	"topdoctors/pkg/logger"
)

// @title TopDoctors API
// @version 1.0
// @description API for managing clinic patients and medical diagnostics.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and then your personal token.

func main() {
	// Load Config
	cfg, errLoadCfg := config.LoadConfig()
	if errLoadCfg != nil {
		slog.Error("Failed to load configuration", "error", errLoadCfg)
		os.Exit(1)
	}
	slog.Info("Configuration loaded successfully")

	// Configure Logger
	logger.SetConfig(logger.Config{
		Level: &cfg.Logs.Level,
	})

	// Initialize Repository (Infrastructure)
	repo, err := persistence.NewGormRepository(
		persistence.Config{
			DSN: cfg.Database.DSN,
		},
	)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	slog.Info("Connected to database successfully")

	// Initialize Support (Infrastructure)
	support := shared.NewSupport()

	// Initialize Application Services (Application)
	app := application.NewApplication(
		repo,
		repo,
		support,
		cfg,
	)
	slog.Info("Application services initialized")

	// Initialize Handler (Adapter)
	h := httpinfra.NewHttpHandler(app, cfg)

	// Initialize and Start Server (Infrastructure)
	server := httpinfra.NewServer(cfg, h)
	if err := server.Start(); err != nil {
		slog.Error("Server stopped unexpectedly", "error", err)
		os.Exit(1)
	}
}

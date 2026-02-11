package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"topdoctors/internal/adapters/handler"
	"topdoctors/internal/adapters/repository"
	"topdoctors/internal/config"
	"topdoctors/internal/core/services"
)

func main() {
	// Load Config
	cfg, errLoadCfg := config.LoadConfig()
	if errLoadCfg != nil {
		slog.Error("Failed to load configuration", "error", errLoadCfg)
		//It consider it a failure
		os.Exit(1)
	}

	// Initialize Repository (Infrastructure)
	repo, err := repository.NewGormRepository()
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
	}

	// Initialize Services (Application)
	authService := services.NewAuthService(repo, cfg)
	patientService := services.NewPatientService(repo)
	diagnosisService := services.NewDiagnosisService(repo, repo)

	// Initialize Handler (Adapter)
	h := handler.NewHttpHandler(authService, patientService, diagnosisService, cfg)

	// Router setup
	mux := http.NewServeMux()

	// Public Routes
	mux.HandleFunc("POST /login", h.Login)
	mux.HandleFunc("POST /register", h.Register)      // Optional
	mux.HandleFunc("POST /patients", h.CreatePatient) // Helper for verification

	// Protected Routes
	mux.Handle("GET /diagnostics", h.AuthMiddleware(http.HandlerFunc(h.GetDiagnostics)))
	mux.Handle("POST /diagnostics", h.AuthMiddleware(http.HandlerFunc(h.CreateDiagnosis)))

	// Start Server
	fmt.Printf("Starting server on port %s...\n", cfg.Port)
	errStartServer := http.ListenAndServe(":"+cfg.Port, mux)
	if errStartServer != nil {
		slog.Error("Server stopped unexpectedly", "error", errStartServer)
		os.Exit(1)
	}
}

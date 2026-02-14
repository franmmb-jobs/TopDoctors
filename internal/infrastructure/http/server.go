package http

import (
	"log/slog"
	"net/http"
	_ "topdoctors/docs"
	"topdoctors/internal/infrastructure/config"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Server struct {
	cfg     *config.Config
	handler *HttpHandler
}

func NewServer(cfg *config.Config, h *HttpHandler) *Server {
	return &Server{
		cfg:     cfg,
		handler: h,
	}
}

func (s *Server) Start() error {
	mux := s.setupRouter()

	slog.Info("Starting server", "port", s.cfg.Api.Port)

	return http.ListenAndServe(":"+s.cfg.Api.Port, mux)
}

func (s *Server) GetHandler() http.Handler {
	return s.setupRouter()
}

func (s *Server) setupRouter() *http.ServeMux {
	slog.Debug("Setting up router")
	mux := http.NewServeMux()
	h := s.handler

	// Public Routes
	mux.HandleFunc("POST /login", h.Login)
	mux.HandleFunc("POST /register", h.Register)

	// Protected Routes
	mux.Handle("GET /diagnostics", h.AuthMiddleware(http.HandlerFunc(h.GetDiagnostics)))
	mux.Handle("POST /diagnostics", h.AuthMiddleware(http.HandlerFunc(h.CreateDiagnosis)))
	mux.Handle("POST /patients", h.AuthMiddleware(http.HandlerFunc(h.CreatePatient)))

	// Swagger UI
	mux.Handle("GET /swagger/", httpSwagger.WrapHandler)

	return mux
}

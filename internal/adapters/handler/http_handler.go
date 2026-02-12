package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"topdoctors/internal/core/domain"
	"topdoctors/internal/core/ports"
	"topdoctors/internal/infrastructure/config"

	"github.com/golang-jwt/jwt/v5"
)

type HttpHandler struct {
	authService      ports.AuthService
	patientService   ports.PatientService
	diagnosisService ports.DiagnosisService
	cfg              *config.Config
}

func NewHttpHandler(
	authService ports.AuthService,
	patientService ports.PatientService,
	diagnosisService ports.DiagnosisService,
	cfg *config.Config,
) *HttpHandler {
	return &HttpHandler{
		authService:      authService,
		patientService:   patientService,
		diagnosisService: diagnosisService,
		cfg:              cfg,
	}
}

// Login Handler
func (h *HttpHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// Register Handler (Optional, for creating users)
func (h *HttpHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.authService.Register(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Create Diagnosis Handler
func (h *HttpHandler) CreateDiagnosis(w http.ResponseWriter, r *http.Request) {
	var diagnosis domain.Diagnosis
	if err := json.NewDecoder(r.Body).Decode(&diagnosis); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.diagnosisService.CreateDiagnosis(&diagnosis)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Get Diagnostics Handler
func (h *HttpHandler) GetDiagnostics(w http.ResponseWriter, r *http.Request) {
	patientName := r.URL.Query().Get("patient_name")
	date := r.URL.Query().Get("date")

	diagnostics, err := h.diagnosisService.GetDiagnostics(patientName, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(diagnostics)
}

// Create Patient Handler (For seeding)
func (h *HttpHandler) CreatePatient(w http.ResponseWriter, r *http.Request) {
	var patient domain.Patient
	if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.patientService.CreatePatient(&patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Auth Middleware
func (h *HttpHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(h.cfg.Api.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

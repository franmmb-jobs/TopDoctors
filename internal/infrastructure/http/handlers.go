package http

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"
	"topdoctors/internal/application"
	"topdoctors/internal/infrastructure/config"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	userIDKey contextKey = "user_id"
)

type HttpHandler struct {
	app *application.Application
	cfg *config.Config
}

func NewHttpHandler(app *application.Application, cfg *config.Config) *HttpHandler {
	return &HttpHandler{
		app: app,
		cfg: cfg,
	}
}

// Login handles user authentication
// @Summary User login
// @Description Authenticate a user and return a JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body LoginRequest true "Login Credentials"
// @Success 200 {object} LoginResponse
// @Failure 401 {string} string "Invalid credentials"
// @Router /login [post]
func (h *HttpHandler) Login(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Login request received")
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode login request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.app.Auth().Login(req.Username, req.Password)
	if err != nil {
		slog.Warn("Invalid login attempt", "username", req.Username)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	slog.Info("User logged in successfully", "username", req.Username)
	response := LoginResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Register handles user registration
// @Summary Register user
// @Description Register a new user in the system
// @Tags Auth
// @Accept json
// @Produce json
// @Param register body RegisterRequest true "Registration Info"
// @Success 201 {string} string "Created"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /register [post]
func (h *HttpHandler) Register(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Register request received")
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode register request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.app.Auth().Register(req.Username, req.Password)
	if err != nil {
		slog.Error("Failed to register user", "username", req.Username, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("User registered successfully", "username", req.Username)
	w.WriteHeader(http.StatusCreated)
}

// CreateDiagnosis handles the creation of a new medical diagnosis
// @Summary Create diagnosis
// @Description Add a new diagnosis to a patient
// @Tags Diagnostics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param diagnosis body CreateDiagnosisRequest true "Diagnosis Info"
// @Success 201 {string} string "Created"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /diagnostics [post]
func (h *HttpHandler) CreateDiagnosis(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Create diagnosis request received")
	var req CreateDiagnosisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode create diagnosis request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse date
	var diagnosisDate time.Time
	if req.Date != "" {
		parsedDate, err := time.Parse(time.RFC3339, req.Date)
		if err != nil {
			slog.Warn("Invalid date format in diagnosis request", "date", req.Date)
			http.Error(w, "Invalid date format, use ISO 8601", http.StatusBadRequest)
			return
		}
		diagnosisDate = parsedDate
	} else {
		diagnosisDate = time.Now()
	}

	// Map to domain
	diagnosis := toDiagnosisDomain(req)
	diagnosis.Date = diagnosisDate

	err := h.app.Patient().CreateDiagnosis(&diagnosis)
	if err != nil {
		slog.Error("Failed to create diagnosis", "patient_id", req.PatientID, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Diagnosis created successfully", "patient_id", req.PatientID)
	w.WriteHeader(http.StatusCreated)
}

// GetDiagnostics searches for diagnostics based on filters
// @Summary Search diagnostics
// @Description Retrieve a list of diagnostics filtering by patient name and/or date range
// @Tags Diagnostics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param patient_name query string false "Filter by patient name"
// @Param date_start query string false "Filter by start date (YYYY-MM-DD)"
// @Param date_end query string false "Filter by end date (YYYY-MM-DD)"
// @Success 200 {array} DiagnosisResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /diagnostics [get]
func (h *HttpHandler) GetDiagnostics(w http.ResponseWriter, r *http.Request) {
	patientName := r.URL.Query().Get("patient_name")
	dateStart := r.URL.Query().Get("date_start")
	dateEnd := r.URL.Query().Get("date_end")

	slog.Debug("Get diagnostics request received", "patient_name", patientName, "date_start", dateStart, "date_end", dateEnd)

	var parsedPatientName *string
	if patientName != "" {
		parsedPatientName = &patientName
	}

	var parsedDateStart *time.Time
	if dateStart != "" {
		d, err := time.Parse("2006-01-02", dateStart)
		if err == nil {
			parsedDateStart = &d
		} else {
			slog.Warn("Invalid date_start format", "date", dateStart)
			http.Error(w, "Invalid date_start format", http.StatusBadRequest)
			return
		}
	}

	var parsedDateEnd *time.Time
	if dateEnd != "" {
		d, err := time.Parse("2006-01-02", dateEnd)
		if err == nil {
			parsedDateEnd = &d
		} else {
			slog.Warn("Invalid date_end format", "date", dateEnd)
			http.Error(w, "Invalid date_end format", http.StatusBadRequest)
			return
		}
	}

	if parsedPatientName == nil && parsedDateStart == nil && parsedDateEnd == nil {
		slog.Warn("Get diagnostics request missing parameters")
		http.Error(w, "At least one parameter is required", http.StatusBadRequest)
		return
	}

	diagnostics, err := h.app.Patient().GetDiagnostics(parsedPatientName, parsedDateStart, parsedDateEnd)
	if err != nil {
		slog.Error("Failed to get diagnostics", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Map to DTOs
	response := toDiagnosisResponseList(diagnostics)

	slog.Info("Diagnostics retrieved successfully", "count", len(response))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreatePatient handles the registration of a new patient
// @Summary Create patient
// @Description Record a new patient in the system
// @Tags Patients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param patient body CreatePatientRequest true "Patient Info"
// @Success 201 {object} PatientResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /patients [post]
func (h *HttpHandler) CreatePatient(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Create patient request received")
	var req CreatePatientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode create patient request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Map to domain
	patient := toPatientDomain(req)

	err := h.app.Patient().CreatePatient(&patient)
	if err != nil {
		slog.Error("Failed to create patient", "name", req.Name, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Patient created successfully", "patient_id", patient.ID, "name", patient.Name)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toPatientResponse(patient))
}

// Auth Middleware
func (h *HttpHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			slog.Warn("Unauthorized request: missing Authorization header", "path", r.URL.Path)
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			slog.Warn("Unauthorized request: invalid Authorization header format", "path", r.URL.Path)
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
			slog.Warn("Unauthorized request: invalid token", "path", r.URL.Path, "error", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			slog.Error("Failed to extract claims from token", "path", r.URL.Path)
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			slog.Error("Subject claim missing or not a string", "path", r.URL.Path)
			http.Error(w, "Invalid token subject", http.StatusUnauthorized)
			return
		}

		// Inject user_id into context
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		r = r.WithContext(ctx)

		slog.Debug("Authorized request", "path", r.URL.Path, "user_id", userID)
		next.ServeHTTP(w, r)
	})
}

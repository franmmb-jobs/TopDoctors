package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"topdoctors/internal/adapters/handler"
	"topdoctors/internal/adapters/repository"
	"topdoctors/internal/config"
	"topdoctors/internal/core/services"
)

func TestAPI_Flow(t *testing.T) {
	// Use a temporary file DB
	dbFile := "test_integration.db"
	os.Remove(dbFile)
	defer os.Remove(dbFile)

	// Initialize Dependencies
	repo, err := repository.NewGormRepositoryWithDSN(dbFile)
	if err != nil {
		t.Fatalf("Failed to init repo: %v", err)
	}

	cfg := &config.Config{JWTSecret: "test-secret", Port: "8080"}

	authService := services.NewAuthService(repo, cfg)
	patientService := services.NewPatientService(repo)
	diagnosisService := services.NewDiagnosisService(repo, repo)

	h := handler.NewHttpHandler(authService, patientService, diagnosisService, cfg)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", h.Register)
	mux.HandleFunc("POST /login", h.Login)
	mux.HandleFunc("POST /patients", h.CreatePatient)
	mux.Handle("POST /diagnostics", h.AuthMiddleware(http.HandlerFunc(h.CreateDiagnosis)))
	mux.Handle("GET /diagnostics", h.AuthMiddleware(http.HandlerFunc(h.GetDiagnostics)))

	server := httptest.NewServer(mux)
	defer server.Close()

	client := server.Client()

	// 1. Register User
	registerPayload := `{"username": "doc", "password": "password"}`
	resp, err := client.Post(server.URL+"/register", "application/json", bytes.NewBufferString(registerPayload))
	if err != nil {
		t.Fatalf("Failed to register: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected 201 Created for register, got %d", resp.StatusCode)
	}

	// 2. Login
	loginPayload := `{"username": "doc", "password": "password"}`
	resp, err = client.Post(server.URL+"/login", "application/json", bytes.NewBufferString(loginPayload))
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK for login, got %d", resp.StatusCode)
	}

	var loginResp map[string]string
	json.NewDecoder(resp.Body).Decode(&loginResp)
	token := loginResp["token"]
	if token == "" {
		t.Fatal("Token is empty")
	}

	// 3. Create Patient
	patientPayload := `{"name": "Jane Doe", "dni": "98765432B", "email": "hane@example.com"}`
	resp, err = client.Post(server.URL+"/patients", "application/json", bytes.NewBufferString(patientPayload))
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Errorf("Failed to create patient: %v, status: %d", err, resp.StatusCode)
	}
	// We assume ID 1

	// 4. Create Diagnosis
	diagnosisPayload := `{"patient_id": 1, "diagnosis": "Fever", "date": "2023-11-01T10:00:00Z"}`
	req, _ := http.NewRequest("POST", server.URL+"/diagnostics", bytes.NewBufferString(diagnosisPayload))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Errorf("Failed to create diagnosis: %v, status: %d", err, resp.StatusCode)
	}

	// 5. Get Diagnostics
	req, _ = http.NewRequest("GET", server.URL+"/diagnostics?patient_name=Jane", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Errorf("Failed to get diagnostics: %v, status: %d", err, resp.StatusCode)
	}
}

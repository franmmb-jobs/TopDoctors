package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"topdoctors/internal/application"
	"topdoctors/internal/infrastructure/config"
	httpinfra "topdoctors/internal/infrastructure/http"
	"topdoctors/internal/infrastructure/persistence"
	"topdoctors/internal/infrastructure/shared"
	"topdoctors/pkg/logger"
)

func TestAPI_Flow(t *testing.T) {

	// Load Config
	cfg, errLoadCfg := config.LoadConfig()
	if errLoadCfg != nil {
		t.Fatalf("Failed to load configuration: %v", errLoadCfg)
	}

	// Configure Logger
	logger.SetConfig(logger.Config{
		Level: &cfg.Logs.Level,
	})

	// Use a temporary file DB
	dbFile := cfg.Database.DSN
	os.Remove(dbFile) // Clean before start

	// Initialize Dependencies
	repo, err := persistence.NewGormRepository(
		persistence.Config{DSN: dbFile},
	)
	if err != nil {
		t.Fatalf("Failed to init repo: %v", err)
	}
	defer func() {
		repo.Close()
		os.Remove(dbFile)
	}()

	support := shared.NewSupport()
	// Initialize Application Services
	app := application.NewApplication(repo, repo, support, cfg)

	h := httpinfra.NewHttpHandler(app, cfg)

	// Initialize Server and get Handler
	var baseURL string
	var client *http.Client

	if externalURL := os.Getenv("API_URL"); externalURL != "" {
		baseURL = externalURL
		client = http.DefaultClient
		slog.Info("Using external API for integration tests", "url", baseURL)
	} else {
		serverAPI := httpinfra.NewServer(cfg, h)
		router := serverAPI.GetHandler()
		server := httptest.NewServer(router)
		defer server.Close()
		baseURL = server.URL
		client = server.Client()
		slog.Info("Using local httptest server for integration tests", "url", baseURL)
	}

	// 1. Register User
	registerPayload := `{"username": "doc", "password": "password"}`
	resp, err := client.Post(baseURL+"/register", "application/json", bytes.NewBufferString(registerPayload))
	if err != nil {
		t.Fatalf("Failed to register: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected 201 Created for register, got %d", resp.StatusCode)
	}

	// 2. Login
	loginPayload := `{"username": "doc", "password": "password"}`
	resp, err = client.Post(baseURL+"/login", "application/json", bytes.NewBufferString(loginPayload))
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
	patientPayload := `{"name": "Jane Doe", "dni": "11111111H", "email": "hane@example.com"}`
	req, _ := http.NewRequest("POST", baseURL+"/patients", bytes.NewBufferString(patientPayload))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("Failed to create patient: %v, status: %d, body: %s", err, resp.StatusCode, string(body))
	}
	var patientResp httpinfra.PatientResponse
	json.NewDecoder(resp.Body).Decode(&patientResp)
	patientID := patientResp.ID
	if patientID == "" {
		t.Fatal("Patient ID is empty")
	}

	// 4. Create Diagnosis
	diagnosisPayload := `{"patient_id": "` + patientID + `", "diagnosis": "Fever", "date": "2023-11-01T10:00:00Z"}`
	req, _ = http.NewRequest("POST", baseURL+"/diagnostics", bytes.NewBufferString(diagnosisPayload))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("Failed to create diagnosis: %v, status: %d, body: %s", err, resp.StatusCode, string(body))
	}

	// 5. Get Diagnostics
	req, _ = http.NewRequest("GET", baseURL+"/diagnostics?patient_name=Jane", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Errorf("Failed to get diagnostics: %v, status: %d", err, resp.StatusCode)
	}
}

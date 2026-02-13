package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"topdoctors/internal/application"
	"topdoctors/internal/infrastructure/config"
	httpinfra "topdoctors/internal/infrastructure/http"
	"topdoctors/internal/infrastructure/persistence"
	"topdoctors/internal/infrastructure/shared"
)

func TestAPI_Flow(t *testing.T) {
	// Use a temporary file DB
	dbFile := "test_integration.db"
	os.Remove(dbFile)
	defer os.Remove(dbFile)

	// Initialize Dependencies
	repo, err := persistence.NewGormRepository(
		persistence.Config{DSN: dbFile},
	)
	if err != nil {
		t.Fatalf("Failed to init repo: %v", err)
	}

	cfg := &config.Config{Api: config.ApiConfig{JWTSecret: "test-secret", Port: "8080"}}

	support := shared.NewSupport()
	// Initialize Application Services
	app := application.NewApplication(repo, repo, support, cfg)

	h := httpinfra.NewHttpHandler(app, cfg)

	// Initialize Server and get Handler
	serverAPI := httpinfra.NewServer(cfg, h)

	router := serverAPI.GetHandler()

	server := httptest.NewServer(router)

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
	req, _ := http.NewRequest("POST", server.URL+"/diagnostics", bytes.NewBufferString(diagnosisPayload))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("Failed to create diagnosis: %v, status: %d, body: %s", err, resp.StatusCode, string(body))
	}

	// 5. Get Diagnostics
	req, _ = http.NewRequest("GET", server.URL+"/diagnostics?patient_name=Jane", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Errorf("Failed to get diagnostics: %v, status: %d", err, resp.StatusCode)
	}
}

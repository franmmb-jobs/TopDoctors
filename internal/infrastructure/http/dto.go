package http

import (
	"time"
	"topdoctors/internal/domain"
)

// Request DTOs

type LoginRequest struct {
	Username string `json:"username" example:"doctor_admin"`
	Password string `json:"password" example:"secret123"`
}

type RegisterRequest struct {
	Username string `json:"username" example:"new_doctor"`
	Password string `json:"password" example:"secure_password"`
}

type CreatePatientRequest struct {
	Name    string `json:"name" example:"Maria Garcia"`
	DNI     string `json:"dni" example:"12345678X"`
	Email   string `json:"email" example:"maria@example.com"`
	Phone   string `json:"phone" example:"+34600123456"`
	Address string `json:"address" example:"Calle Mayor 1, Madrid"`
}

type CreateDiagnosisRequest struct {
	PatientID    string `json:"patient_id" example:"01HMGNBPJNX0G2BZXJ7XW1RHPR"`
	Diagnosis    string `json:"diagnosis" example:"Gripe común"`
	Prescription string `json:"prescription" example:"Ibuprofeno 600mg cada 8h"`
	Date         string `json:"date" example:"2026-02-13T10:00:00Z"` // ISO 8601 format
}

// Response DTOs

type LoginResponse struct {
	Token string `json:"token" example:"string"`
}

type PatientResponse struct {
	ID        string    `json:"id" example:"01HMGNBPJNX0G2BZXJ7XW1RHPR"`
	Name      string    `json:"name" example:"Maria Garcia"`
	DNI       string    `json:"dni" example:"12345678X"`
	Email     string    `json:"email" example:"maria@example.com"`
	Phone     string    `json:"phone" example:"+34600123456"`
	Address   string    `json:"address" example:"Calle Mayor 1, Madrid"`
	CreatedAt time.Time `json:"created_at" example:"2026-02-13T18:23:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2026-02-13T18:23:00Z"`
}

type DiagnosisResponse struct {
	ID           string          `json:"id" example:"01HMGNBPJNX0G2BZXJ7XW1RHPR"`
	PatientID    string          `json:"patient_id" example:"01HMGNBPJNX0G2BZXJ7XW1RHPR"`
	Patient      PatientResponse `json:"patient,omitempty"`
	Diagnosis    string          `json:"diagnosis" example:"Fiebre alta y tos persistente"`
	Prescription string          `json:"prescription" example:"Paracetamol 1g cada 8 horas"`
	Date         time.Time       `json:"date" example:"2026-02-13T18:23:00Z"`
	CreatedAt    time.Time       `json:"created_at" example:"2026-02-13T18:23:00Z"`
}

// Mappers: Domain -> DTO

func toPatientResponse(p domain.Patient) PatientResponse {
	return PatientResponse{
		ID:      p.ID,
		Name:    p.Name,
		DNI:     p.DNI,
		Email:   p.Email,
		Phone:   p.Phone,
		Address: p.Address,
	}
}

func toDiagnosisResponse(d domain.Diagnosis) DiagnosisResponse {
	return DiagnosisResponse{
		ID:           d.ID,
		PatientID:    d.PatientID,
		Patient:      toPatientResponse(d.Patient),
		Diagnosis:    d.Diagnosis,
		Prescription: d.Prescription,
		Date:         d.Date,
	}
}

func toDiagnosisResponseList(diagnostics []domain.Diagnosis) []DiagnosisResponse {
	result := make([]DiagnosisResponse, len(diagnostics))
	for i, d := range diagnostics {
		result[i] = toDiagnosisResponse(d)
	}
	return result
}

// Mappers: DTO -> Domain

func toPatientDomain(req CreatePatientRequest) domain.Patient {
	return domain.Patient{
		Name:    req.Name,
		DNI:     req.DNI,
		Email:   req.Email,
		Phone:   req.Phone,
		Address: req.Address,
	}
}

func toDiagnosisDomain(req CreateDiagnosisRequest) domain.Diagnosis {
	// Date parsing will be handled in the handler
	return domain.Diagnosis{
		PatientID:    req.PatientID,
		Diagnosis:    req.Diagnosis,
		Prescription: req.Prescription,
	}
}

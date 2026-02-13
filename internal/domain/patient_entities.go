package domain

import (
	"errors"
	"time"
)

var (
	ErrEmptyPatientID     = errors.New("patient ID cannot be empty")
	ErrEmptyDNI           = errors.New("patient DNI cannot be empty")
	ErrEmptyName          = errors.New("patient name cannot be empty")
	ErrEmptyEmail         = errors.New("patient email cannot be empty")
	ErrEmptyDiagnosisID   = errors.New("diagnosis ID cannot be empty")
	ErrEmptyDiagnosisText = errors.New("diagnosis text cannot be empty")
	ErrEmptyPatientFK     = errors.New("patient ID is required for diagnosis")
	ErrEmptyDate          = errors.New("diagnosis date is required")
)

// Patient represents a patient in the system
type Patient struct {
	ID        string
	Name      string
	DNI       string
	Email     string
	Phone     string
	Address   string
	Diagnosis []Diagnosis
}

// Validate ensures the patient's domain invariants are met
func (p *Patient) Validate() error {
	if p.ID == "" {
		return ErrEmptyPatientID
	}
	if p.Name == "" {
		return ErrEmptyName
	}
	if p.DNI == "" {
		return ErrEmptyDNI
	}
	if p.Email == "" {
		return ErrEmptyEmail
	}
	if p.Diagnosis != nil {
		for _, d := range p.Diagnosis {
			if err := d.Validate(); err != nil {
				return err
			}
		}
	}
	return nil
}

// Diagnosis represents a medical diagnosis
type Diagnosis struct {
	ID           string
	PatientID    string
	Patient      Patient
	Diagnosis    string
	Prescription string
	Date         time.Time
}

// Validate ensures the diagnosis domain invariants are met
func (d *Diagnosis) Validate() error {
	if d.PatientID == "" {
		return ErrEmptyPatientFK
	}
	if d.Diagnosis == "" {
		return ErrEmptyDiagnosisText
	}
	if d.Date.IsZero() {
		return ErrEmptyDate
	}
	return nil
}

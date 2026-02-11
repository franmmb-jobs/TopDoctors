package ports

import (
	"time"
	"topdoctors/internal/core/domain"
)

// PatientRepository interface
type PatientRepository interface {
	CreatePatient(patient *domain.Patient) error
	GetByID(id uint) (*domain.Patient, error)
	GetByDNI(dni string) (*domain.Patient, error)
	// Add other methods as needed
}

// DiagnosisRepository interface
type DiagnosisRepository interface {
	CreateDiagnosis(diagnosis *domain.Diagnosis) error
	GetByPatientID(patientID uint) ([]domain.Diagnosis, error)
	GetByDateRange(startDate, endDate time.Time) ([]domain.Diagnosis, error)
	GetByPatientName(name string) ([]domain.Diagnosis, error)
	Search(patientName string, date *time.Time) ([]domain.Diagnosis, error)
}

// UserRepository interface (for Auth)
type UserRepository interface {
	GetByUsername(username string) (*domain.User, error)
	CreateUser(user *domain.User) error
}

// AuthService interface
type AuthService interface {
	Login(username, password string) (string, error) // Returns token
	Register(username, password string) error
}

// DiagnosisService interface
type DiagnosisService interface {
	CreateDiagnosis(diagnosis *domain.Diagnosis) error
	GetDiagnostics(patientName string, date string) ([]domain.Diagnosis, error)
}

// PatientService interface
type PatientService interface {
	CreatePatient(patient *domain.Patient) error
	GetPatient(dni string) (*domain.Patient, error)
}

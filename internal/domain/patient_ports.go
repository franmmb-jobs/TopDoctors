package domain

import "time"

// Medical Domain - Repository Interfaces (Driven Ports - Outbound)

// PatientRepository defines operations for patient persistence
type PatientRepository interface {
	CreatePatient(patient *Patient) error
	GetPatientByID(id string) (*Patient, error)
	GetPatientByDNI(dni string) (*Patient, error)
	CreateDiagnosis(diagnosis *Diagnosis) error
	GetDiagnosisByPatientID(patientID string) ([]Diagnosis, error)
	GetByDiagnosisDateRange(startDate, endDate time.Time) ([]Diagnosis, error)
	GetDiagnosisByPatientName(name string) ([]Diagnosis, error)
	SearchDiagnosis(patientName *string, dateStart, dateEnd *time.Time) ([]Diagnosis, error)
}

// PatientService defines patient business operations
type PatientService interface {
	CreatePatient(patient *Patient) error
	GetPatient(dni string) (*Patient, error)
	CreateDiagnosis(diagnosis *Diagnosis) error
	GetDiagnostics(patientName *string, dateStart, dateEnd *time.Time) ([]Diagnosis, error)
}

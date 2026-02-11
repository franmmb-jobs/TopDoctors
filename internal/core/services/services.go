package services

import (
	"time"
	"topdoctors/internal/core/domain"
	"topdoctors/internal/core/ports"
)

// PatientService implementation
type PatientService struct {
	repo ports.PatientRepository
}

func NewPatientService(repo ports.PatientRepository) *PatientService {
	return &PatientService{repo: repo}
}

func (s *PatientService) CreatePatient(patient *domain.Patient) error {
	patient.CreatedAt = time.Now()
	patient.UpdatedAt = time.Now()
	return s.repo.CreatePatient(patient)
}

func (s *PatientService) GetPatient(dni string) (*domain.Patient, error) {
	return s.repo.GetByDNI(dni)
}

// DiagnosisService implementation
type DiagnosisService struct {
	repo        ports.DiagnosisRepository
	patientRepo ports.PatientRepository
}

func NewDiagnosisService(repo ports.DiagnosisRepository, patientRepo ports.PatientRepository) *DiagnosisService {
	return &DiagnosisService{
		repo:        repo,
		patientRepo: patientRepo,
	}
}

func (s *DiagnosisService) CreateDiagnosis(diagnosis *domain.Diagnosis) error {
	// Validate patient exists
	_, err := s.patientRepo.GetByID(diagnosis.PatientID)
	if err != nil {
		return err
	}
	diagnosis.CreatedAt = time.Now()
	return s.repo.CreateDiagnosis(diagnosis)
}

func (s *DiagnosisService) GetDiagnostics(patientName string, date string) ([]domain.Diagnosis, error) {
	var parsedDate *time.Time
	if date != "" {
		d, err := time.Parse("2006-01-02", date)
		if err == nil {
			parsedDate = &d
		}
	}
	return s.repo.Search(patientName, parsedDate)
}

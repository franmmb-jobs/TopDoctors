package application

import (
	"time"
	"topdoctors/internal/domain"
)

type PatientService struct {
	repo    domain.PatientRepository
	support domain.Support
}

func NewPatientService(repo domain.PatientRepository, support domain.Support) *PatientService {
	return &PatientService{repo: repo, support: support}
}

func (s *PatientService) CreatePatient(patient *domain.Patient) error {
	id, errCreateID := s.support.CreateNewID()
	if errCreateID != nil {
		return errCreateID
	}
	patient.ID = id

	// Enforce domain invariants
	if errValidate := patient.Validate(); errValidate != nil {
		return errValidate
	}

	return s.repo.CreatePatient(patient)
}

func (s *PatientService) GetPatient(dni string) (*domain.Patient, error) {
	return s.repo.GetPatientByDNI(dni)
}

func (s *PatientService) CreateDiagnosis(diagnosis *domain.Diagnosis) error {
	id, errCreateID := s.support.CreateNewID()
	if errCreateID != nil {
		return errCreateID
	}
	diagnosis.ID = id

	// Enforce domain invariants
	if errValidate := diagnosis.Validate(); errValidate != nil {
		return errValidate
	}

	// Internal logic: Validate patient exists in DB
	_, errGetPatient := s.repo.GetPatientByID(diagnosis.PatientID)
	if errGetPatient != nil {
		return errGetPatient
	}

	return s.repo.CreateDiagnosis(diagnosis)
}

func (s *PatientService) GetDiagnostics(patientName *string, dateStart, dateEnd *time.Time) ([]domain.Diagnosis, error) {
	return s.repo.SearchDiagnosis(patientName, dateStart, dateEnd)
}

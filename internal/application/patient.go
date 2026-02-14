package application

import (
	"log/slog"
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
		slog.Error("ID creation failed for patient", "error", errCreateID)
		return errCreateID
	}
	patient.ID = id

	// Enforce domain invariants
	if errValidate := patient.Validate(); errValidate != nil {
		slog.Warn("Patient validation failed", "error", errValidate)
		return errValidate
	}

	err := s.repo.CreatePatient(patient)
	if err != nil {
		slog.Error("Patient creation in repository failed", "error", err)
		return err
	}

	slog.Info("Patient created successfully", "patient_id", patient.ID)
	return nil
}

func (s *PatientService) GetPatient(dni string) (*domain.Patient, error) {
	return s.repo.GetPatientByDNI(dni)
}

func (s *PatientService) CreateDiagnosis(diagnosis *domain.Diagnosis) error {
	id, errCreateID := s.support.CreateNewID()
	if errCreateID != nil {
		slog.Error("ID creation failed for diagnosis", "error", errCreateID)
		return errCreateID
	}
	diagnosis.ID = id

	// Enforce domain invariants
	if errValidate := diagnosis.Validate(); errValidate != nil {
		slog.Warn("Diagnosis validation failed", "error", errValidate)
		return errValidate
	}

	// Internal logic: Validate patient exists in DB
	_, errGetPatient := s.repo.GetPatientByID(diagnosis.PatientID)
	if errGetPatient != nil {
		slog.Warn("Diagnosis creation failed: patient not found", "patient_id", diagnosis.PatientID)
		return errGetPatient
	}

	err := s.repo.CreateDiagnosis(diagnosis)
	if err != nil {
		slog.Error("Diagnosis creation in repository failed", "error", err)
		return err
	}

	slog.Info("Diagnosis created successfully", "diagnosis_id", diagnosis.ID, "patient_id", diagnosis.PatientID)
	return nil
}

func (s *PatientService) GetDiagnostics(patientName *string, dateStart, dateEnd *time.Time) ([]domain.Diagnosis, error) {
	return s.repo.SearchDiagnosis(patientName, dateStart, dateEnd)
}

package services_test

import (
	"errors"
	"testing"
	"time"
	"topdoctors/internal/core/domain"
	"topdoctors/internal/core/ports/mocks"
	"topdoctors/internal/core/services"

	"go.uber.org/mock/gomock"
)

func TestDiagnosisService_CreateDiagnosis(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDiagnosisRepo := mocks.NewMockDiagnosisRepository(ctrl)
	mockPatientRepo := mocks.NewMockPatientRepository(ctrl)

	// Assuming DiagnosisService implements ports.DiagnosisService
	// But struct is specific. We use the struct directly or update testing to use interface if needed.
	// Here we test the struct methods.
	diagnosisService := services.NewDiagnosisService(mockDiagnosisRepo, mockPatientRepo)

	diagnosis := &domain.Diagnosis{
		PatientID:    1,
		Diagnosis:    "Flu",
		Prescription: "Rest",
		Date:         time.Now(),
	}

	t.Run("Success", func(t *testing.T) {
		// Expect check for patient existence
		mockPatientRepo.EXPECT().GetByID(uint(1)).Return(&domain.Patient{ID: 1}, nil)

		// Expect create call
		mockDiagnosisRepo.EXPECT().CreateDiagnosis(gomock.Any()).Return(nil)

		err := diagnosisService.CreateDiagnosis(diagnosis)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Patient Not Found", func(t *testing.T) {
		mockPatientRepo.EXPECT().GetByID(uint(1)).Return(nil, errors.New("patient not found"))

		err := diagnosisService.CreateDiagnosis(diagnosis)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestDiagnosisService_GetDiagnostics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDiagnosisRepo := mocks.NewMockDiagnosisRepository(ctrl)
	mockPatientRepo := mocks.NewMockPatientRepository(ctrl)
	diagnosisService := services.NewDiagnosisService(mockDiagnosisRepo, mockPatientRepo)

	t.Run("Filter by Name", func(t *testing.T) {
		expectedDiagnostics := []domain.Diagnosis{
			{ID: 1, Diagnosis: "Flu"},
		}

		mockDiagnosisRepo.EXPECT().Search("John", nil).Return(expectedDiagnostics, nil)

		result, err := diagnosisService.GetDiagnostics("John", "")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(result) != 1 {
			t.Errorf("expected 1 diagnosis, got %d", len(result))
		}
	})
}

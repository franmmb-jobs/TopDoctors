package application

import (
	"errors"
	"testing"
	"time"
	"topdoctors/internal/domain"
	"topdoctors/internal/mocks"

	"go.uber.org/mock/gomock"
)

func TestPatientService_CreatePatient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPatientRepository(ctrl)
	mockSupport := mocks.NewMockSupport(ctrl)
	service := NewPatientService(mockRepo, mockSupport)

	patient := &domain.Patient{
		Name:  "Maria Garcia",
		DNI:   "12345678Z",
		Email: "maria@example.com",
	}

	t.Run("successful creation", func(t *testing.T) {
		mockSupport.EXPECT().CreateNewID().Return("01HMGNBPJNX0G2BZXJ7XW1RHPR", nil)
		mockRepo.EXPECT().CreatePatient(gomock.Any()).Return(nil)

		err := service.CreatePatient(patient)
		if err != nil {
			t.Errorf("CreatePatient() unexpected error = %v", err)
		}
		if patient.ID != "01HMGNBPJNX0G2BZXJ7XW1RHPR" {
			t.Errorf("CreatePatient() expected ID to be set, got %s", patient.ID)
		}
	})

	t.Run("ID creation failure", func(t *testing.T) {
		mockSupport.EXPECT().CreateNewID().Return("", errors.New("id error"))

		err := service.CreatePatient(patient)
		if err == nil {
			t.Error("CreatePatient() expected error, got nil")
		}
	})

	t.Run("validation failure", func(t *testing.T) {
		invalidPatient := &domain.Patient{Name: ""} // Missing ID, Name, DNI, etc.
		mockSupport.EXPECT().CreateNewID().Return("valid-id", nil)

		err := service.CreatePatient(invalidPatient)
		if err == nil {
			t.Error("CreatePatient() expected validation error, got nil")
		}
	})
}

func TestPatientService_CreateDiagnosis(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPatientRepository(ctrl)
	mockSupport := mocks.NewMockSupport(ctrl)
	service := NewPatientService(mockRepo, mockSupport)

	diagnosis := &domain.Diagnosis{
		PatientID: "01HMGNBPJNX0G2BZXJ7XW1RHPR",
		Diagnosis: "Fever",
		Date:      time.Now(),
	}

	t.Run("successful creation", func(t *testing.T) {
		mockSupport.EXPECT().CreateNewID().Return("diag-id", nil)
		mockRepo.EXPECT().GetPatientByID(diagnosis.PatientID).Return(&domain.Patient{}, nil)
		mockRepo.EXPECT().CreateDiagnosis(diagnosis).Return(nil)

		err := service.CreateDiagnosis(diagnosis)
		if err != nil {
			t.Errorf("CreateDiagnosis() unexpected error = %v", err)
		}
	})

	t.Run("patient not found", func(t *testing.T) {
		mockSupport.EXPECT().CreateNewID().Return("diag-id", nil)
		mockRepo.EXPECT().GetPatientByID(diagnosis.PatientID).Return(nil, errors.New("not found"))

		err := service.CreateDiagnosis(diagnosis)
		if err == nil {
			t.Error("CreateDiagnosis() expected error, got nil")
		}
	})
}

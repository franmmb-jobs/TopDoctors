package services_test

import (
	"fmt"
	"topdoctors/internal/core/domain"
)

// ExampleDiagnosisService_CreateDiagnosis demonstrates how to define a diagnosis
// struct that would be passed to the service.
func ExampleDiagnosisService_CreateDiagnosis() {
	diagnosis := &domain.Diagnosis{
		PatientID:    1,
		Diagnosis:    "Common Cold",
		Prescription: "Vitamin C and Rest",
	}

	// In a real scenario, you would call:
	// err := diagnosisService.CreateDiagnosis(diagnosis)

	fmt.Printf("Diagnosis for Patient %d: %s\n", diagnosis.PatientID, diagnosis.Diagnosis)
	// Output: Diagnosis for Patient 1: Common Cold
}

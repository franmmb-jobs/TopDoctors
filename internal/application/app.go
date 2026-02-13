package application

import (
	"topdoctors/internal/domain"
	"topdoctors/internal/infrastructure/config"
)

// Application is the container for all application services
type Application struct {
	auth    domain.UserService
	patient domain.PatientService
	support domain.Support
}

// NewApplication creates a new application instance with all services
func NewApplication(
	userRepo domain.UserRepository,
	patientRepo domain.PatientRepository,
	support domain.Support,
	cfg *config.Config,
) *Application {

	return &Application{
		auth:    NewAuthService(userRepo, support, cfg),
		patient: NewPatientService(patientRepo, support),
	}
}

// Auth returns the authentication service
func (a *Application) Auth() domain.UserService {
	return a.auth
}

// Patient returns the patient service
func (a *Application) Patient() domain.PatientService {
	return a.patient
}

package domain

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
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
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrInvalidDNI         = errors.New("invalid DNI format")
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

	//Validate DNI format
	if p.DNI == "" {
		return ErrEmptyDNI
	}
	_, err := ValidarDNI(p.DNI)
	if err != nil {
		return err
	}

	if p.Email == "" {
		return ErrEmptyEmail
	}
	if !ValidarEmail(p.Email) {
		return ErrInvalidEmail
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

// ValidarDNI verifica si un DNI español es matemáticamente consistente.
func ValidarDNI(dni string) (bool, error) {
	dni = strings.ToUpper(strings.TrimSpace(dni))

	// 1. Validar formato básico con Regex (8 números + 1 letra)
	re := regexp.MustCompile(`^[0-9]{8}[TRWAGMYFPDXBNJZSQVHLCKE]$`)
	if !re.MatchString(dni) {
		return false, ErrInvalidDNI
	}

	// 2. Separar número y letra
	numeroParte := dni[:8]
	letraProporcionada := string(dni[8])

	// 3. Calcular la letra que debería tener
	numero, _ := strconv.Atoi(numeroParte)
	letras := "TRWAGMYFPDXBNJZSQVHLCKE"
	letraCorrecta := string(letras[numero%23])

	// 4. Comparar
	if letraProporcionada != letraCorrecta {
		return false, ErrInvalidDNI
	}

	return true, nil
}

// ValidarEmail verifica si un email tiene un formato válido.
func ValidarEmail(email string) bool {
	// Regex simple para validar formato de email
	// No es perfecto (RFC 5322), pero es suficiente para la mayoría de los casos de uso.
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(emailRegex, email)
	return match
}

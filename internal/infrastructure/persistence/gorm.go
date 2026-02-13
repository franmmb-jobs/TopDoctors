package persistence

import (
	"time"
	"topdoctors/internal/domain"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// GormRepository implements all repository interfaces
type GormRepository struct {
	db  *gorm.DB
	cfg Config
}
type Config struct {
	DSN string
}

func NewGormRepository(cfg Config) (*GormRepository, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate
	err = db.AutoMigrate(&PatientDB{}, &DiagnosisDB{}, &UserDB{}, &UserTokenDB{})
	if err != nil {
		return nil, err
	}

	return &GormRepository{db: db, cfg: cfg}, nil
}

// Patient Repository Implementation
func (r *GormRepository) CreatePatient(patient *domain.Patient) error {

	dbPatient := toPatientDB(patient)
	err := r.db.Create(dbPatient).Error
	if err == nil {
		patient.ID = dbPatient.ULID
	}
	return err
}

func (r *GormRepository) GetPatientByID(id string) (*domain.Patient, error) {
	var patient PatientDB
	err := r.db.Where("ulid = ?", id).First(&patient).Error
	if err != nil {
		return nil, err
	}
	return toPatientDomain(&patient), nil
}

func (r *GormRepository) GetPatientByDNI(dni string) (*domain.Patient, error) {
	var patient PatientDB
	err := r.db.Where("dni = ?", dni).First(&patient).Error
	if err != nil {
		return nil, err
	}
	return toPatientDomain(&patient), nil
}

// Diagnosis Repository Implementation
func (r *GormRepository) CreateDiagnosis(diagnosis *domain.Diagnosis) error {
	dbDiagnosis := toDiagnosisDB(diagnosis)
	err := r.db.Create(dbDiagnosis).Error
	if err == nil {
		diagnosis.ID = dbDiagnosis.ULID
	}
	return err
}

func (r *GormRepository) GetDiagnosisByPatientID(patientID string) ([]domain.Diagnosis, error) {
	var diagnostics []DiagnosisDB
	err := r.db.Where("patient_ulid = ?", patientID).Find(&diagnostics).Error
	if err != nil {
		return nil, err
	}

	result := make([]domain.Diagnosis, len(diagnostics))
	for i, d := range diagnostics {
		result[i] = *toDiagnosisDomain(&d)
	}
	return result, nil
}

func (r *GormRepository) GetByDiagnosisDateRange(startDate, endDate time.Time) ([]domain.Diagnosis, error) {
	var diagnostics []DiagnosisDB
	err := r.db.Where("date BETWEEN ? AND ?", startDate, endDate).Find(&diagnostics).Error
	if err != nil {
		return nil, err
	}

	result := make([]domain.Diagnosis, len(diagnostics))
	for i, d := range diagnostics {
		result[i] = *toDiagnosisDomain(&d)
	}
	return result, nil
}

func (r *GormRepository) GetDiagnosisByPatientName(name string) ([]domain.Diagnosis, error) {
	var diagnostics []DiagnosisDB
	err := r.db.Joins("Patient").Where("Patient.name LIKE ?", "%"+name+"%").Find(&diagnostics).Error
	if err != nil {
		return nil, err
	}

	result := make([]domain.Diagnosis, len(diagnostics))
	for i, d := range diagnostics {
		result[i] = *toDiagnosisDomain(&d)
	}
	return result, nil
}

func (r *GormRepository) SearchDiagnosis(patientName *string, dateStart, dateEnd *time.Time) ([]domain.Diagnosis, error) {
	query := r.db.Model(&DiagnosisDB{}).Preload("Patient").Joins("Patient")

	if patientName != nil {
		query = query.Where("Patient.name LIKE ?", "%"+*patientName+"%")
	}

	if dateStart != nil {
		startOfDay := time.Date(dateStart.Year(), dateStart.Month(), dateStart.Day(), 0, 0, 0, 0, dateStart.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)
		query = query.Where("diagnoses.date >= ? AND diagnoses.date < ?", startOfDay, endOfDay)
	}

	if dateEnd != nil && dateEnd != dateStart {
		startOfDay := time.Date(dateEnd.Year(), dateEnd.Month(), dateEnd.Day(), 0, 0, 0, 0, dateEnd.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)
		query = query.Where("diagnoses.date <= ? AND diagnoses.date < ?", startOfDay, endOfDay)
	}

	var diagnostics []DiagnosisDB
	err := query.Find(&diagnostics).Error
	if err != nil {
		return nil, err
	}

	result := make([]domain.Diagnosis, len(diagnostics))
	for i, d := range diagnostics {
		result[i] = *toDiagnosisDomain(&d)
	}
	return result, nil
}

// User Repository Implementation
func (r *GormRepository) GetByUsername(username string) (*domain.User, error) {
	var user UserDB
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return toUserDomain(&user), nil
}

func (r *GormRepository) CreateUser(user *domain.User) error {
	dbUser := toUserDB(user)
	err := r.db.Create(dbUser).Error
	if err == nil {
		user.ID = dbUser.ULID
	}
	return err
}

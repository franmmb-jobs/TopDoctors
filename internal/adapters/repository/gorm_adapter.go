package repository

import (
	"time"
	"topdoctors/internal/core/domain"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository() (*GormRepository, error) {
	db, err := gorm.Open(sqlite.Open("diagnostics.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate
	err = db.AutoMigrate(&domain.Patient{}, &domain.Diagnosis{}, &domain.User{})
	if err != nil {
		return nil, err
	}

	return &GormRepository{db: db}, nil
}

func NewGormRepositoryWithDSN(dsn string) (*GormRepository, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate
	err = db.AutoMigrate(&domain.Patient{}, &domain.Diagnosis{}, &domain.User{})
	if err != nil {
		return nil, err
	}

	return &GormRepository{db: db}, nil
}

// Patient Repository Implementation
func (r *GormRepository) CreatePatient(patient *domain.Patient) error {
	return r.db.Create(patient).Error
}

func (r *GormRepository) GetByID(id uint) (*domain.Patient, error) {
	var patient domain.Patient
	err := r.db.First(&patient, id).Error
	return &patient, err
}

func (r *GormRepository) GetByDNI(dni string) (*domain.Patient, error) {
	var patient domain.Patient
	err := r.db.Where("dni = ?", dni).First(&patient).Error
	return &patient, err
}

// Diagnosis Repository Implementation
func (r *GormRepository) CreateDiagnosis(diagnosis *domain.Diagnosis) error {
	return r.db.Create(diagnosis).Error
}

func (r *GormRepository) GetByPatientID(patientID uint) ([]domain.Diagnosis, error) {
	var diagnostics []domain.Diagnosis
	err := r.db.Where("patient_id = ?", patientID).Find(&diagnostics).Error
	return diagnostics, err
}

func (r *GormRepository) GetByDateRange(startDate, endDate time.Time) ([]domain.Diagnosis, error) {
	var diagnostics []domain.Diagnosis
	err := r.db.Where("date BETWEEN ? AND ?", startDate, endDate).Find(&diagnostics).Error
	return diagnostics, err
}

func (r *GormRepository) GetByPatientName(name string) ([]domain.Diagnosis, error) {
	var diagnostics []domain.Diagnosis
	// Join with patient table
	err := r.db.Joins("Patient").Where("Patient.name LIKE ?", "%"+name+"%").Find(&diagnostics).Error
	return diagnostics, err
}

func (r *GormRepository) Search(patientName string, date *time.Time) ([]domain.Diagnosis, error) {
	query := r.db.Model(&domain.Diagnosis{}).Preload("Patient").Joins("Patient")

	if patientName != "" {
		query = query.Where("Patient.name LIKE ?", "%"+patientName+"%")
	}

	if date != nil {
		// Filter by day, ignoring time
		startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)
		query = query.Where("diagnoses.date >= ? AND diagnoses.date < ?", startOfDay, endOfDay)
	}

	var diagnostics []domain.Diagnosis
	err := query.Find(&diagnostics).Error
	return diagnostics, err
}

// User Repository Implementation
func (r *GormRepository) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *GormRepository) CreateUser(user *domain.User) error {
	return r.db.Create(user).Error
}

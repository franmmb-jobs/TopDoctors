package persistence

import (
	"time"
	"topdoctors/internal/domain"
)

// GORM models with tags (infrastructure concern)
type PatientDB struct {
	ID        uint   `gorm:"primaryKey,autoIncrement"`
	ULID      string `gorm:"column:ulid;unique"`
	Name      string
	DNI       string `gorm:"unique"`
	Email     string
	Phone     string
	Address   string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (PatientDB) TableName() string {
	return "patients"
}

type DiagnosisDB struct {
	ID           uint   `gorm:"primaryKey,autoIncrement"`
	ULID         string `gorm:"column:ulid;unique"`
	PatientULID  string `gorm:"column:patient_ulid"`
	PatientID    uint
	Patient      PatientDB `gorm:"foreignKey:PatientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Diagnosis    string
	Prescription string
	Date         time.Time
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func (DiagnosisDB) TableName() string {
	return "diagnoses"
}

type UserDB struct {
	ID       uint   `gorm:"primaryKey,autoIncrement"`
	ULID     string `gorm:"column:ulid;unique"`
	Username string `gorm:"unique"`
	Password string
}

func (UserDB) TableName() string {
	return "users"
}

type UserTokenDB struct {
	ID        uint      `gorm:"primaryKey,autoIncrement"`
	UserID    uint      `gorm:"column:user_id"`
	User      UserDB    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Token     string    `gorm:"unique"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (UserTokenDB) TableName() string {
	return "user_tokens"
}

// Mappers from domain to DB
func toPatientDB(p *domain.Patient) *PatientDB {
	return &PatientDB{
		ULID:    p.ID,
		Name:    p.Name,
		DNI:     p.DNI,
		Email:   p.Email,
		Phone:   p.Phone,
		Address: p.Address,
	}
}

func toPatientDomain(p *PatientDB) *domain.Patient {
	return &domain.Patient{
		ID:      p.ULID,
		Name:    p.Name,
		DNI:     p.DNI,
		Email:   p.Email,
		Phone:   p.Phone,
		Address: p.Address,
	}
}

func toDiagnosisDB(d *domain.Diagnosis) *DiagnosisDB {
	return &DiagnosisDB{
		ULID:         d.ID,
		PatientULID:  d.PatientID,
		Diagnosis:    d.Diagnosis,
		Prescription: d.Prescription,
		Date:         d.Date,
	}
}

func toDiagnosisDomain(d *DiagnosisDB) *domain.Diagnosis {
	diagnosis := &domain.Diagnosis{
		ID:           d.ULID,
		PatientID:    d.PatientULID,
		Diagnosis:    d.Diagnosis,
		Prescription: d.Prescription,
		Date:         d.Date,
	}

	// Only map patient if it was preloaded
	if d.Patient.ULID != "" {
		diagnosis.Patient = *toPatientDomain(&d.Patient)
	}

	return diagnosis
}

func toUserDB(u *domain.User) *UserDB {
	return &UserDB{
		ULID:     u.ID,
		Username: u.Username,
		Password: u.Password,
	}
}

func toUserDomain(u *UserDB) *domain.User {
	return &domain.User{
		ID:       u.ULID,
		Username: u.Username,
		Password: u.Password,
	}
}

func toUserTokenDB(t *domain.UserToken) *UserTokenDB {
	return &UserTokenDB{
		Token: t.Token,
	}
}

func toUserTokenDomain(t *UserTokenDB) *domain.UserToken {
	return &domain.UserToken{
		Token: t.Token,
	}
}

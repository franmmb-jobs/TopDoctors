package persistence

import (
	"time"
	"topdoctors/internal/domain"
)

// GORM models with tags (infrastructure concern)
type PatientDB struct {
	ID        uint      `gorm:"primaryKey,autoIncrement"`
	ULID      string    `gorm:"column:ulid;unique"`
	Name      string    `json:"name"`
	DNI       string    `gorm:"unique"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (PatientDB) TableName() string {
	return "patients"
}

type DiagnosisDB struct {
	ID           uint      `gorm:"primaryKey,autoIncrement"`
	ULID         string    `gorm:"column:ulid;unique"`
	PatientULID  string    `gorm:"column:patient_ulid"`
	PatientID    uint      `json:"patient_id"`
	Patient      PatientDB `gorm:"foreignKey:PatientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Diagnosis    string    `json:"diagnosis"`
	Prescription string    `json:"prescription"`
	Date         time.Time `json:"date"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
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
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	return &domain.Diagnosis{
		ID:           d.ULID,
		PatientID:    d.PatientULID,
		Patient:      *toPatientDomain(&d.Patient),
		Diagnosis:    d.Diagnosis,
		Prescription: d.Prescription,
		Date:         d.Date,
	}
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

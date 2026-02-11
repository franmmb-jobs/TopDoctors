package domain

import "time"

// Patient structure
type Patient struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	DNI       string    `json:"dni" gorm:"unique"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Diagnosis structure
type Diagnosis struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	PatientID    uint      `json:"patient_id"`
	Patient      Patient   `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Diagnosis    string    `json:"diagnosis"`
	Prescription string    `json:"prescription"`
	Date         time.Time `json:"date"`
	CreatedAt    time.Time `json:"created_at"`
}

// User structure for authentication
type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"` // Stored as hash
}

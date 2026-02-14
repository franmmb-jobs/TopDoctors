package domain

import (
	"testing"
)

func TestPatient_Validate(t *testing.T) {
	tests := []struct {
		name    string
		patient Patient
		wantErr error
	}{
		{
			name: "valid patient",
			patient: Patient{
				ID:    "01HMGNBPJNX0G2BZXJ7XW1RHPR",
				Name:  "Maria Garcia",
				DNI:   "12345678Z",
				Email: "maria@example.com",
			},
			wantErr: nil,
		},
		{
			name: "missing ID",
			patient: Patient{
				Name:  "Maria Garcia",
				DNI:   "12345678Z",
				Email: "maria@example.com",
			},
			wantErr: ErrEmptyPatientID,
		},
		{
			name: "missing Name",
			patient: Patient{
				ID:    "01HMGNBPJNX0G2BZXJ7XW1RHPR",
				DNI:   "12345678Z",
				Email: "maria@example.com",
			},
			wantErr: ErrEmptyName,
		},
		{
			name: "missing DNI",
			patient: Patient{
				ID:    "01HMGNBPJNX0G2BZXJ7XW1RHPR",
				Name:  "Maria Garcia",
				Email: "maria@example.com",
			},
			wantErr: ErrEmptyDNI,
		},
		{
			name: "invalid DNI format",
			patient: Patient{
				ID:    "01HMGNBPJNX0G2BZXJ7XW1RHPR",
				Name:  "Maria Garcia",
				DNI:   "12345678A", // Wrong letter
				Email: "maria@example.com",
			},
			wantErr: ErrInvalidDNI,
		},
		{
			name: "missing Email",
			patient: Patient{
				ID:   "01HMGNBPJNX0G2BZXJ7XW1RHPR",
				Name: "Maria Garcia",
				DNI:  "12345678Z",
			},
			wantErr: ErrEmptyEmail,
		},
		{
			name: "invalid Email format",
			patient: Patient{
				ID:    "01HMGNBPJNX0G2BZXJ7XW1RHPR",
				Name:  "Maria Garcia",
				DNI:   "12345678Z",
				Email: "invalid-email",
			},
			wantErr: ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.patient.Validate(); err != tt.wantErr {
				t.Errorf("Patient.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidarDNI(t *testing.T) {
	tests := []struct {
		name    string
		dni     string
		want    bool
		wantErr error
	}{
		{"valid Z", "12345678Z", true, nil}, // 12345678 % 23 = 14 -> Z
		{"valid S", "11111111H", true, nil}, // 11111111 % 23 = 21 -> H
		{"valid Q", "87654321X", true, nil}, // 87654321 % 23 = 10 -> X
		{"invalid length", "1234567Z", false, ErrInvalidDNI},
		{"invalid character", "123456789", false, ErrInvalidDNI},
		{"wrong letter", "12345678A", false, ErrInvalidDNI},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidarDNI(tt.dni)
			if got != tt.want {
				t.Errorf("ValidarDNI() got = %v, want %v", got, tt.want)
			}
			if err != tt.wantErr {
				t.Errorf("ValidarDNI() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidarEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{"valid", "test@example.com", true},
		{"valid with dot", "test.name@example.co.uk", true},
		{"no @", "testexample.com", false},
		{"no dot", "test@example", false},
		{"empty", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidarEmail(tt.email); got != tt.want {
				t.Errorf("ValidarEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

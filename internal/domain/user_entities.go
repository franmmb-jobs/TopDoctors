package domain

import "errors"

var (
	ErrEmptyUsername      = errors.New("username cannot be empty")
	ErrEmptyPassword      = errors.New("password cannot be empty")
	ErrEmptyUserID        = errors.New("user ID is required")
	ErrEmptyToken         = errors.New("token cannot be empty")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

// User represents an authenticated user
type User struct {
	ID       string
	Username string
	Password string // Stored as hash
	Token    *UserToken
}

// Validate ensures the user's domain invariants are met
func (p *User) Validate() error {
	if p.ID == "" {
		return ErrEmptyUserID
	}
	if p.Username == "" {
		return ErrEmptyUsername
	}
	if p.Password == "" {
		return ErrEmptyPassword
	}
	if p.Token != nil {
		return p.Token.Validate()
	}
	return nil
}

type UserToken struct {
	Token string
}

// Validate
func (p *UserToken) Validate() error {
	//Maybe add a validation for the token in the future
	return nil
}

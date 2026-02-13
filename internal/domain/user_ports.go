package domain

// Authentication Domain - Repository Interfaces (Driven Ports - Outbound)

// UserRepository defines operations for user persistence
type UserRepository interface {
	GetByUsername(username string) (*User, error)
	CreateUser(user *User) error
}

// Authentication Domain - Service Interfaces (Driving Ports - Inbound)

// UserService defines authentication operations
type UserService interface {
	Login(username, password string) (string, error)
	Register(username, password string) error
	ValidateToken(token string) error
}

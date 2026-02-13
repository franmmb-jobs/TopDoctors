package domain

type Support interface {
	CreateNewID() (string, error)
	GenerateHashPassword(password string) (string, error)
	CompareHashPassword(password string, hash string) error
	GenerateToken(user *User, secret string) (string, error)
	ValidateToken(token string, secret string) error
}

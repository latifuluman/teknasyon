package data

type UsersRepository interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id string) (*User, error)
	UpdateUser(user User) error
	DeleteUserByID(id string) error
	InsertUser(user User) (string, error)
	ResetPassword(password string, user User) error
	PasswordMatches(plainText string, user User) (bool, error)
}
type DatabaseRepository interface {
	UsersRepository
}

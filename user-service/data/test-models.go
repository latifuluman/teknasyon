package data

import (
	"database/sql"
	"time"
)

type PostgresTestRepository struct {
	Conn *sql.DB
}

func NewPostgresTestRepository(db *sql.DB) *PostgresTestRepository {
	return &PostgresTestRepository{
		Conn: db,
	}
}

// GetByEmail returns one user by email
func (u *PostgresTestRepository) GetUserByEmail(email string) (*User, error) {
	user := User{
		ID:        "",
		FirstName: "First",
		LastName:  "Last",
		Email:     "me@here.com",
		Password:  "",
		Active:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &user, nil
}

// GetOne returns one user by id
func (u *PostgresTestRepository) GetUserByID(id string) (*User, error) {
	user := User{
		ID:        "",
		FirstName: "Latif",
		LastName:  "Uluman",
		Email:     "latifuluman@gmail.com",
		Password:  "",
		Active:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &user, nil
}

// Update updates one user in the database, using the information
// stored in the receiver u
func (u *PostgresTestRepository) UpdateUser(user User) error {
	return nil
}

// DeleteByID deletes one user from the database, by ID
func (u *PostgresTestRepository) DeleteUserByID(id string) error {
	return nil
}

// Insert inserts a new user into the database, and returns the ID of the newly inserted row
func (u *PostgresTestRepository) InsertUser(user User) (string, error) {
	return "2", nil
}

// ResetPassword is the method we will use to change a user's password.
func (u *PostgresTestRepository) ResetPassword(password string, user User) error {
	return nil
}

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
// with the hash we have stored for a given user in the database. If the password
// and hash match, we return true; otherwise, we return false.
func (u *PostgresTestRepository) PasswordMatches(plainText string, user User) (bool, error) {
	return true, nil
}

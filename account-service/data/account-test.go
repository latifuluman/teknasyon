package data

import "time"

// GetAll returns a slice of all users, sorted by last name
func (u *PostgresTestRepository) GetAllAccounts(userID string) ([]*Account, error) {
	users := []*Account{}

	return users, nil
}

func (u *PostgresTestRepository) GetAccountByID(id string) (*Account, error) {
	user := Account{
		ID:          "",
		UserID:      "",
		AccountName: "",
		AccountType: "",
		Balance:     0,
		Active:      1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return &user, nil
}

// Update updates one user in the database, using the information
// stored in the receiver u
func (u *PostgresTestRepository) UpdateAccount(user *Account) error {
	return nil
}

// DeleteByID deletes one user from the database, by ID
func (u *PostgresTestRepository) DeleteAccountByID(id string) error {
	return nil
}

// Insert inserts a new user into the database, and returns the ID of the newly inserted row
func (u *PostgresTestRepository) InsertAccount(user Account) (string, error) {
	return "", nil
}

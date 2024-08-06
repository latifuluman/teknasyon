package data

import "database/sql"

func (u *PostgresTestRepository) Begin() (*sql.Tx, error) {
	return nil, nil
}

// GetAll returns a slice of all users, sorted by last name
func (u *PostgresTestRepository) GetAllTransactionsBySenderID(senderID string) ([]*Transaction, error) {
	transactions := []*Transaction{}

	return transactions, nil
}

// GetOne returns one user by id
func (u *PostgresTestRepository) GetTransactionByID(id string) (*Transaction, error) {
	transaction := &Transaction{}

	return transaction, nil
}

// Insert inserts a new user into the database, and returns the ID of the newly inserted row
func (u *PostgresTestRepository) InsertTransaction(t Transaction) (string, error) {

	return "", nil
}

package data

import (
	"context"
	"log"
	"time"
)

// User is the structure which holds one user from the database.
type Account struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	AccountName string    `json:"account_name"`
	AccountType string    `json:"account_type"`
	Balance     float32   `json:"balance"`
	Active      int       `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetAll returns a slice of all users, sorted by last name
func (u *PostgresRepository) GetAllAccounts(userID string) ([]*Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, account_name, account_type, balance, active, created_at, updated_at
	from accounts.accounts where user_id = $1 order by updated_at`

	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*Account

	for rows.Next() {
		var user Account
		err := rows.Scan(
			&user.ID,
			&user.AccountName,
			&user.AccountType,
			&user.Balance,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

// GetOne returns one user by id
func (u *PostgresRepository) GetAccountByID(id string) (*Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, user_id, account_name, account_type, balance, active, created_at, updated_at from accounts.accounts where id = $1`

	var user Account
	row := db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.UserID,
		&user.AccountName,
		&user.AccountType,
		&user.Balance,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update updates one user in the database, using the information
// stored in the receiver u
func (u *PostgresRepository) UpdateAccount(user *Account) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update accounts.accounts set
		account_name = $1,
		account_type = $2,
		balance = $3,
		active = $4,
		updated_at= $5
		where id = $6
	`

	_, err := db.ExecContext(ctx, stmt,
		user.AccountName,
		user.AccountType,
		user.Balance,
		user.Active,
		time.Now(),
		user.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// DeleteByID deletes one user from the database, by ID
func (u *PostgresRepository) DeleteAccountByID(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from accounts.accounts where id = $1`

	_, err := db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

// Insert inserts a new user into the database, and returns the ID of the newly inserted row
func (u *PostgresRepository) InsertAccount(a Account) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID string
	stmt := `insert into accounts.accounts (user_id, account_name, account_type, balance, active, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7) returning id`

	err := db.QueryRowContext(ctx, stmt,
		a.UserID,
		a.AccountName,
		a.AccountType,
		a.Balance,
		1,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return "", err
	}

	return newID, nil
}

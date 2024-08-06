package data

import (
	"context"
	"database/sql"
	"log"
	"time"
)

// User is the structure which holds one user from the database.
type Transaction struct {
	ID         string    `json:"id"`
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	Amount     float32   `json:"amount"`
	Type       string    `json:"type"`
	CreatedAt  time.Time `json:"created_at"`
}

func (u *PostgresRepository) Begin() (*sql.Tx, error) {
	return u.Conn.Begin()
}

// GetAll returns a slice of all users, sorted by last name
func (u *PostgresRepository) GetAllTransactionsBySenderID(senderID string) ([]*Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, sender_id, receiver_id, amount, type, created_at
	from accounts.transactions where user_id = $1 order by created_at`

	rows, err := db.QueryContext(ctx, query, senderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*Transaction

	for rows.Next() {
		var transaction Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.SenderID,
			&transaction.ReceiverID,
			&transaction.Amount,
			&transaction.Type,
			&transaction.CreatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

// GetOne returns one user by id
func (u *PostgresRepository) GetTransactionByID(id string) (*Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, sender_id, receiver_id, amount, type, created_at from accounts.transactions where id = $1`

	var transaction Transaction
	row := db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&transaction.ID,
		&transaction.SenderID,
		&transaction.ReceiverID,
		&transaction.Amount,
		&transaction.Type,
		&transaction.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

// Insert inserts a new user into the database, and returns the ID of the newly inserted row
func (u *PostgresRepository) InsertTransaction(t Transaction) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID string
	stmt := `insert into accounts.transactions (sender_id, receiver_id, amount, type, created_at)
		values ($1, $2, $3, $4, $5) returning id`

	err := db.QueryRowContext(ctx, stmt,
		t.SenderID,
		t.ReceiverID,
		t.Amount,
		t.Type,
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return "", err
	}

	return newID, nil
}

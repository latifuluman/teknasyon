package data

import "database/sql"

type DatabaseRepository interface {
	AccountsRepository
	TransactionsRepository
}

type TransactionsRepository interface {
	Begin() (*sql.Tx, error)
	GetAllTransactionsBySenderID(senderID string) ([]*Transaction, error)
	GetTransactionByID(id string) (*Transaction, error)
	InsertTransaction(t Transaction) (string, error)
}

type AccountsRepository interface {
	GetAccountByID(id string) (*Account, error)
	GetAllAccounts(userID string) ([]*Account, error)
	UpdateAccount(user *Account) error
	DeleteAccountByID(id string) error
	InsertAccount(user Account) (string, error)
}

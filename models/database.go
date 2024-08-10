package models

type Database interface {
	CreateAccount(a *Account) error
	GetAllAccounts() ([]*Account, error)
	GetAccountByID(id string) (*Account, error)
	DeleteAccount(id string) error
	UpdateAccount(id string, a *Account) error
}

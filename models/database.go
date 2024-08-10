package models

type Database interface {
	CreateAccount(a *Account) error
	GetAllAccounts() ([]*Account, error)
	GetAccountByID(id int) (*Account, error)
	DeleteAccount(id int) error
	UpdateAccount(id int, a *Account) error
}

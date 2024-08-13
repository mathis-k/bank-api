package models

type Database interface {
	CreateAccount(req *AccountRequest) (*Account, error)
	GetAllAccounts(maxResult int) ([]*Account, error)
	GetAccountByID(id string) (*Account, error)
	DeleteAccount(id string) (*Account, error)
	UpdateAccount(id string, req *AccountRequest) (*Account, error)
}

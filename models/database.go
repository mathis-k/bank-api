package models

type Database interface {
	Connect() error
	Disconnect() error
	CreateAccount(req *AccountRequest) (*Account, error)
	GetAllAccounts(maxResult uint64) ([]*Account, error)
	GetAccountByID(id string) (*Account, error)
	DeleteAccount(id string) (*Account, error)
	UpdateAccount(id string, req *AccountRequest) (*Account, error)
	Transfer(id string, req *TransferRequest) error
}

package main

type Storage interface {
	CreateAccount(a *Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(id int) (*Account, error)
	DeleteAccount(id int) error
	UpdateAccount(id int, a *Account) error
}

type MongoDbStorage struct {
}

package main

import (
	"math/rand"
)

type APIServer struct {
	listenAddress string
	storage       Storage
}
type APIResponse struct {
	Message string `json:"message"`
}

const welcomeMessage = `Welcome to the Bank JSON API Server! :)

Available endpoints:
GET /account - get all accounts
POST /account - create a new account
GET /account/{id} - get account by ID
PUT /account/{id} - update account by ID
DELETE /account/{id} - delete account by ID`

type Account struct {
	ID        int     `json:"id"` //unique account ID
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     string  `json:"email"`
	Number    int64   `json:"number"` //unique account number
	Balance   float64 `json:"balance"`
}

func NewAccount(firstName, lastName, email string, balance ...float64) *Account {
	var initialBalance float64
	if len(balance) > 0 {
		initialBalance = balance[0]
	} else {
		initialBalance = 0
	}
	return &Account{
		ID:        rand.Intn(1000),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Number:    int64(rand.Intn(1000000)),
		Balance:   initialBalance,
	}
}

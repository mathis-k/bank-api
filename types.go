package main

import (
	"math/rand"
	"net/http"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error
type APIServer struct {
	listenAddress string
}
type APIResponse struct {
	Message string `json:"message"`
}
type Account struct {
	ID        int     `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     string  `json:"email"`
	Number    int64   `json:"number"`
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

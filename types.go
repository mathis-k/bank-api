package main

import "math/rand"

type Account struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
	Balance   int64  `json:"balance"`
	Number    int64  `json:"number"`
}

func NewAccount(email, lastName, firstName string, balance ...int64) *Account {
	var initialbalance int64
	if len(balance) == 0 {
		initialbalance = 0
	} else {
		initialbalance = balance[0]
	}
	return &Account{
		ID:        rand.Intn(1000),
		Email:     email,
		LastName:  lastName,
		FirstName: firstName,
		Balance:   initialbalance,
		Number:    int64(rand.Intn(1000000)),
	}
}

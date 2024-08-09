package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"time"
)

type Account struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	AccountNumber int64              `bson:"account_number" json:"account_number"`
	Balance       float64            `bson:"balance" json:"balance"`
	FirstName     string             `bson:"first_name" json:"first_name"`
	LastName      string             `bson:"last_name" json:"last_name"`
	Email         string             `bson:"email" json:"email"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
}

func NewAccount(firstName, lastName, email string, balance ...float64) *Account {
	var initialBalance float64
	if len(balance) > 0 {
		initialBalance = balance[0]
	} else {
		initialBalance = 0
	}
	return &Account{
		FirstName:     firstName,
		LastName:      lastName,
		Email:         email,
		AccountNumber: int64(rand.Intn(1000000)),
		Balance:       initialBalance,
		CreatedAt:     time.Now(),
	}
}

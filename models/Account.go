package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Account struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	AccountNumber uint64             `bson:"account_number" json:"account_number"`
	User          primitive.ObjectID `bson:"user" json:"user"`
	Balance       float64            `bson:"balance" json:"balance"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
}

package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Transaction struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FromAccount primitive.ObjectID `bson:"from_account" json:"from_account"`
	ToAccount   primitive.ObjectID `bson:"to_account" json:"to_account"`
	Amount      float64            `bson:"amount" json:"amount"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

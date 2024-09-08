package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName string               `bson:"first_name" json:"first_name"`
	LastName  string               `bson:"last_name" json:"last_name"`
	Email     string               `bson:"email" json:"email"`
	Accounts  []primitive.ObjectID `bson:"accounts" json:"accounts"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
}

type UserRequest struct {
	FirstName string `bson:"first_name" json:"first_name"`
	LastName  string `bson:"last_name" json:"last_name"`
	Email     string `bson:"email" json:"email"`
}
